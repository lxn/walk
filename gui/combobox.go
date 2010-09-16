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
	. "walk/winapi/user32"
)

type ComboBox struct {
	Widget
	items *ComboBoxItemList
}

func NewComboBox(parent IContainer) (*ComboBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("COMBOBOX"), nil,
		CBS_DROPDOWN|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	cb := &ComboBox{Widget: Widget{hWnd: hWnd, parent: parent}}

	cb.items = newComboBoxItemList(cb)

	cb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = cb

	parent.Children().Add(cb)

	return cb, nil
}

func (*ComboBox) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (cb *ComboBox) PreferredSize() drawing.Size {
	return cb.dialogBaseUnitsToPixels(drawing.Size{50, 14})
}

func (cb *ComboBox) raiseEvent(msg *MSG) os.Error {
	return cb.Widget.raiseEvent(msg)
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

	return
}

func (cb *ComboBox) onClearingComboBoxItems() (err os.Error) {
	SendMessage(cb.hWnd, CB_RESETCONTENT, 0, 0)

	return
}
