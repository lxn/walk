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
	. "walk/winapi/user32"
)

type TextEdit struct {
	WidgetBase
}

func NewTextEdit(parent Container) (*TextEdit, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("EDIT"), nil,
		ES_MULTILINE|ES_WANTRETURN|WS_CHILD|WS_TABSTOP|WS_VISIBLE|WS_VSCROLL,
		0, 0, 160, 80, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	te := &TextEdit{
		WidgetBase: WidgetBase{
			hWnd:   hWnd,
			parent: parent,
		},
	}

	succeeded := false
	defer func() {
		if !succeeded {
			te.Dispose()
		}
	}()

	te.layoutFlags = te.LayoutFlagsMask()

	te.SetFont(defaultFont)

	if err := parent.Children().Add(te); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = te

	succeeded = true

	return te, nil
}

func (*TextEdit) LayoutFlagsMask() LayoutFlags {
	return HShrink | HGrow | VShrink | VGrow
}

func (te *TextEdit) PreferredSize() Size {
	return te.dialogBaseUnitsToPixels(Size{100, 100})
}

func (te *TextEdit) TextSelection() (start, end int) {
	SendMessage(te.hWnd, EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (te *TextEdit) SetTextSelection(start, end int) {
	SendMessage(te.hWnd, EM_SETSEL, uintptr(start), uintptr(end))
}
