// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"unsafe"
)

import (
	. "walk/winapi/user32"
)

var textEditOrigWndProcPtr uintptr
var _ subclassedWidget = &TextEdit{}

type TextEdit struct {
	WidgetBase
}

func NewTextEdit(parent Container) (*TextEdit, os.Error) {
	te := &TextEdit{}

	if err := initChildWidget(
		te,
		parent,
		"EDIT",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|ES_MULTILINE|ES_WANTRETURN,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	return te, nil
}

func (*TextEdit) origWndProcPtr() uintptr {
	return textEditOrigWndProcPtr
}

func (*TextEdit) setOrigWndProcPtr(ptr uintptr) {
	textEditOrigWndProcPtr = ptr
}

func (*TextEdit) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (te *TextEdit) PreferredSize() Size {
	return Size{100, 100}
}

func (te *TextEdit) TextSelection() (start, end int) {
	SendMessage(te.hWnd, EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (te *TextEdit) SetTextSelection(start, end int) {
	SendMessage(te.hWnd, EM_SETSEL, uintptr(start), uintptr(end))
}

func (te *TextEdit) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_GETDLGCODE:
		result := CallWindowProc(textEditOrigWndProcPtr, hwnd, msg, wParam, lParam)
		return result &^ DLGC_HASSETSEL
	}

	return te.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
