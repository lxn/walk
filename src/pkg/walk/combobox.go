// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/gdi32"
	. "walk/winapi/user32"
)

var comboBoxOrigWndProcPtr uintptr
var _ subclassedWidget = &ComboBox{}

type ComboBox struct {
	WidgetBase
	items                        *ComboBoxItemList
	maxItemTextWidth             int
	prevCurIndex                 int
	currentIndexChangedPublisher EventPublisher
}

func NewComboBox(parent Container) (*ComboBox, os.Error) {
	cb := &ComboBox{prevCurIndex: -1}

	if err := initChildWidget(
		cb,
		parent,
		"COMBOBOX",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|CBS_DROPDOWNLIST,
		0); err != nil {
		return nil, err
	}

	cb.items = newComboBoxItemList(cb)

	return cb, nil
}

func (*ComboBox) origWndProcPtr() uintptr {
	return comboBoxOrigWndProcPtr
}

func (*ComboBox) setOrigWndProcPtr(ptr uintptr) {
	comboBoxOrigWndProcPtr = ptr
}

func (*ComboBox) LayoutFlags() LayoutFlags {
	return GrowableHorz
}

func (cb *ComboBox) SizeHint() Size {
	defaultSize := cb.dialogBaseUnitsToPixels(Size{50, 12})

	if cb.items != nil && cb.maxItemTextWidth <= 0 {
		cb.maxItemTextWidth = cb.calculateMaxItemTextWidth()
	}

	// FIXME: Use GetThemePartSize instead of guessing
	w := maxi(defaultSize.Width, cb.maxItemTextWidth+24)
	h := defaultSize.Height + 1

	return Size{w, h}
}

func (cb *ComboBox) calculateMaxItemTextWidth() int {
	hdc := GetDC(cb.hWnd)
	if hdc == 0 {
		newError("GetDC failed")
		return -1
	}
	defer ReleaseDC(cb.hWnd, hdc)

	hFontOld := SelectObject(hdc, HGDIOBJ(cb.Font().handleForDPI(0)))
	defer SelectObject(hdc, hFontOld)

	var maxWidth int

	for _, item := range cb.items.items {
		var s SIZE
		str := syscall.StringToUTF16(item.Text())

		if !GetTextExtentPoint32(hdc, &str[0], len(str)-1, &s) {
			newError("GetTextExtentPoint32 failed")
			return -1
		}

		maxWidth = maxi(maxWidth, s.CX)
	}

	return maxWidth
}

func (cb *ComboBox) Items() *ComboBoxItemList {
	return cb.items
}

func (cb *ComboBox) CurrentIndex() int {
	return int(SendMessage(cb.hWnd, CB_GETCURSEL, 0, 0))
}

func (cb *ComboBox) SetCurrentIndex(value int) os.Error {
	index := int(SendMessage(cb.hWnd, CB_SETCURSEL, uintptr(value), 0))

	if index != value {
		return newError("invalid index")
	}

	if value != cb.prevCurIndex {
		cb.prevCurIndex = value
		cb.currentIndexChangedPublisher.Publish()
	}

	return nil
}

func (cb *ComboBox) CurrentIndexChanged() *Event {
	return cb.currentIndexChangedPublisher.Event()
}

func (cb *ComboBox) Text() string {
	return widgetText(cb.hWnd)
}

func (cb *ComboBox) SetText(value string) os.Error {
	return setWidgetText(cb.hWnd, value)
}

func (cb *ComboBox) TextSelection() (start, end int) {
	SendMessage(cb.hWnd, CB_GETEDITSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (cb *ComboBox) SetTextSelection(start, end int) {
	SendMessage(cb.hWnd, CB_SETEDITSEL, 0, uintptr(MAKELONG(uint16(start), uint16(end))))
}

func (cb *ComboBox) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint(wParam)) {
		case CBN_SELENDOK:
			if selIndex := cb.CurrentIndex(); selIndex != cb.prevCurIndex {
				cb.currentIndexChangedPublisher.Publish()
				cb.prevCurIndex = selIndex
				return 0
			}
		}
	}

	return cb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}

func (cb *ComboBox) onInsertingComboBoxItem(index int, item *ComboBoxItem) os.Error {
	if CB_ERR == SendMessage(cb.hWnd, CB_INSERTSTRING, uintptr(index), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item.text)))) {
		return newError("CB_INSERTSTRING failed")
	}

	cb.maxItemTextWidth = 0

	return nil
}

func (cb *ComboBox) onRemovingComboBoxItem(index int, item *ComboBoxItem) os.Error {
	if CB_ERR == SendMessage(cb.hWnd, CB_DELETESTRING, uintptr(index), 0) {
		return newError("CB_DELETESTRING failed")
	}

	cb.maxItemTextWidth = 0
	if index == cb.prevCurIndex {
		cb.prevCurIndex = -1
	}

	return nil
}

func (cb *ComboBox) onClearingComboBoxItems() os.Error {
	SendMessage(cb.hWnd, CB_RESETCONTENT, 0, 0)

	cb.maxItemTextWidth = 0
	cb.prevCurIndex = -1

	return nil
}
