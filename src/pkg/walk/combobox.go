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
	. "walk/winapi/user32"
)

type ComboBox struct {
	WidgetBase
	items                         *ComboBoxItemList
	prevSelIndex                  int
	selectedIndexChangedPublisher EventPublisher
}

func NewComboBox(parent Container) (*ComboBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("COMBOBOX"), nil,
		CBS_DROPDOWNLIST|WS_CHILD|WS_TABSTOP|WS_VISIBLE|WS_VSCROLL,
		0, 0, 0, 0, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	cb := &ComboBox{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}, prevSelIndex: -1}

	succeeded := false
	defer func() {
		if !succeeded {
			cb.Dispose()
		}
	}()

	cb.items = newComboBoxItemList(cb)

	cb.SetFont(defaultFont)

	if err := parent.Children().Add(cb); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = cb

	succeeded = true
	return cb, nil
}

func (*ComboBox) LayoutFlagsMask() LayoutFlags {
	return HShrink | HGrow
}

func (cb *ComboBox) PreferredSize() Size {
	return cb.dialogBaseUnitsToPixels(Size{50, 14})
}

func (cb *ComboBox) Items() *ComboBoxItemList {
	return cb.items
}

func (cb *ComboBox) SelectedIndex() int {
	return int(SendMessage(cb.hWnd, CB_GETCURSEL, 0, 0))
}

func (cb *ComboBox) SetSelectedIndex(value int) os.Error {
	index := int(SendMessage(cb.hWnd, CB_SETCURSEL, uintptr(value), 0))

	if index != value {
		return newError("invalid index")
	}

	if value != cb.prevSelIndex {
		cb.prevSelIndex = value
		cb.selectedIndexChangedPublisher.Publish()
	}

	return nil
}

func (cb *ComboBox) SelectedIndexChanged() *Event {
	return cb.selectedIndexChangedPublisher.Event()
}

func (cb *ComboBox) TextSelection() (start, end int) {
	SendMessage(cb.hWnd, CB_GETEDITSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (cb *ComboBox) SetTextSelection(start, end int) {
	SendMessage(cb.hWnd, CB_SETEDITSEL, 0, uintptr(MAKELONG(uint16(start), uint16(end))))
}

func (cb *ComboBox) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint(wParam)) {
		case CBN_SELENDOK:
			if selIndex := cb.SelectedIndex(); selIndex != cb.prevSelIndex {
				cb.selectedIndexChangedPublisher.Publish()
				cb.prevSelIndex = selIndex
				return 0
			}
		}
	}

	return cb.WidgetBase.wndProc(hwnd, msg, wParam, lParam, origWndProcPtr)
}

func (cb *ComboBox) onInsertingComboBoxItem(index int, item *ComboBoxItem) (err os.Error) {
	if CB_ERR == SendMessage(cb.hWnd, CB_INSERTSTRING, uintptr(index), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item.text)))) {
		err = newError("CB_INSERTSTRING failed")
	}

	return
}

func (cb *ComboBox) onRemovingComboBoxItem(index int, item *ComboBoxItem) (err os.Error) {
	if CB_ERR == SendMessage(cb.hWnd, CB_DELETESTRING, uintptr(index), 0) {
		err = newError("CB_DELETESTRING failed")
	}

	if index == cb.prevSelIndex {
		cb.prevSelIndex = -1
	}

	return
}

func (cb *ComboBox) onClearingComboBoxItems() (err os.Error) {
	SendMessage(cb.hWnd, CB_RESETCONTENT, 0, 0)

	cb.prevSelIndex = -1

	return
}
