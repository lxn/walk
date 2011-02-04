// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/user32"
)

type ComboBox struct {
	Widget
	items                         *ComboBoxItemList
	prevSelIndex                  int
	selectedIndexChangedPublisher EventPublisher
}

func NewComboBox(parent IContainer) (*ComboBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("COMBOBOX"), nil,
		CBS_DROPDOWNLIST|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	cb := &ComboBox{Widget: Widget{hWnd: hWnd, parent: parent}, prevSelIndex: -1}

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

func (*ComboBox) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (cb *ComboBox) PreferredSize() drawing.Size {
	return cb.dialogBaseUnitsToPixels(drawing.Size{50, 14})
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

	return nil
}

func (cb *ComboBox) SelectedIndexChanged() *Event {
	return cb.selectedIndexChangedPublisher.Event()
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

	return cb.Widget.wndProc(hwnd, msg, wParam, lParam, origWndProcPtr)
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
