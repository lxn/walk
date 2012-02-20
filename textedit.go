// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import "unsafe"

import . "github.com/lxn/go-winapi"

var textEditOrigWndProcPtr uintptr
var _ subclassedWidget = &TextEdit{}

type TextEdit struct {
	WidgetBase
}

func NewTextEdit(parent Container) (*TextEdit, error) {
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

func (te *TextEdit) SizeHint() Size {
	return Size{100, 100}
}

func (te *TextEdit) Text() string {
	return widgetText(te.hWnd)
}

func (te *TextEdit) SetText(value string) error {
	return setWidgetText(te.hWnd, value)
}

func (te *TextEdit) TextSelection() (start, end int) {
	SendMessage(te.hWnd, EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (te *TextEdit) SetTextSelection(start, end int) {
	SendMessage(te.hWnd, EM_SETSEL, uintptr(start), uintptr(end))
}

func (te *TextEdit) ReadOnly() bool {
	return te.hasStyleBits(ES_READONLY)
}

func (te *TextEdit) SetReadOnly(readOnly bool) error {
	if 0 == SendMessage(te.hWnd, EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return newError("SendMessage(EM_SETREADONLY)")
	}

	return nil
}

func (te *TextEdit) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

		return DLGC_HASSETSEL | DLGC_WANTARROWS | DLGC_WANTCHARS
	}

	return te.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
