// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

type TextEdit struct {
	WidgetBase
	readOnlyChangedPublisher EventPublisher
	textChangedPublisher     EventPublisher
}

func NewTextEdit(parent Container) (*TextEdit, error) {
	te := new(TextEdit)

	if err := InitWidget(
		te,
		parent,
		"EDIT",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|ES_MULTILINE|ES_WANTRETURN,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	te.MustRegisterProperty("ReadOnly", NewProperty(
		func() interface{} {
			return te.ReadOnly()
		},
		func(v interface{}) error {
			return te.SetReadOnly(v.(bool))
		},
		te.readOnlyChangedPublisher.Event()))

	te.MustRegisterProperty("Text", NewProperty(
		func() interface{} {
			return te.Text()
		},
		func(v interface{}) error {
			return te.SetText(v.(string))
		},
		te.textChangedPublisher.Event()))

	return te, nil
}

func (*TextEdit) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (te *TextEdit) MinSizeHint() Size {
	return te.dialogBaseUnitsToPixels(Size{20, 12})
}

func (te *TextEdit) SizeHint() Size {
	return Size{100, 100}
}

func (te *TextEdit) Text() string {
	return windowText(te.hWnd)
}

func (te *TextEdit) TextLength() int {
	return int(SendMessage(te.hWnd, WM_GETTEXTLENGTH, 0, 0))
}

func (te *TextEdit) SetText(value string) error {
	return setWindowText(te.hWnd, value)
}

func (te *TextEdit) TextSelection() (start, end int) {
	te.SendMessage(EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (te *TextEdit) SetTextSelection(start, end int) {
	te.SendMessage(EM_SETSEL, uintptr(start), uintptr(end))
}

func (te *TextEdit) ReplaceSelectedText(text string, canUndo bool) {
	te.SendMessage(EM_REPLACESEL,
		uintptr(BoolToBOOL(canUndo)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))))
}

func (te *TextEdit) AppendText(value string) {
	s, e := te.TextSelection()
	l := te.TextLength()
	te.SetTextSelection(l, l)
	te.ReplaceSelectedText(value, false)
	te.SetTextSelection(s, e)
}

func (te *TextEdit) ReadOnly() bool {
	return te.hasStyleBits(ES_READONLY)
}

func (te *TextEdit) SetReadOnly(readOnly bool) error {
	if 0 == te.SendMessage(EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return newError("SendMessage(EM_SETREADONLY)")
	}

	te.readOnlyChangedPublisher.Publish()

	return nil
}

func (te *TextEdit) TextChanged() *Event {
	return te.textChangedPublisher.Event()
}

func (te *TextEdit) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case EN_CHANGE:
			te.textChangedPublisher.Publish()
		}

	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

		return DLGC_HASSETSEL | DLGC_WANTARROWS | DLGC_WANTCHARS

	case WM_KEYDOWN:
		if Key(wParam) == KeyA && ControlDown() {
			te.SetTextSelection(0, -1)
		}
	}

	return te.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
