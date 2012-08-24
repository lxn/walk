// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"unsafe"
	"syscall"
)

import . "github.com/lxn/go-winapi"

var logViewOrigWndProcPtr uintptr
var _ subclassedWidget = &LogView{}

type LogView struct {
	WidgetBase
	logChan chan string
}

func NewLogView(parent Container) (*LogView, error) {
	lc := make(chan string, 1024)
	te := &LogView{logChan:lc}

	if err := initChildWidget(
		te,
		parent,
		"EDIT",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|ES_MULTILINE|ES_WANTRETURN,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}
	te.setReadOnly(true)
	SendMessage(te.hWnd, EM_SETLIMITTEXT, 4294967295, 0)
	return te, nil
}

func (*LogView) origWndProcPtr() uintptr {
	return textEditOrigWndProcPtr
}

func (*LogView) setOrigWndProcPtr(ptr uintptr) {
	textEditOrigWndProcPtr = ptr
}

func (*LogView) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (te *LogView) MinSizeHint() Size {
	return te.dialogBaseUnitsToPixels(Size{20, 12})
}

func (te *LogView) SizeHint() Size {
	return Size{100, 100}
}

func (te *LogView) setTextSelection(start, end int) {
	SendMessage(te.hWnd, EM_SETSEL, uintptr(start), uintptr(end))
}

func (te *LogView) textLength() int{
	return int(SendMessage(te.hWnd, 0x000E, uintptr(0), uintptr(0)))
}

func (te *LogView) AppendText(value string) {
	textLength := te.textLength()
	te.setTextSelection(textLength, textLength)
	SendMessage(te.hWnd, EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))))
}

func (te *LogView) setReadOnly(readOnly bool) error {
	if 0 == SendMessage(te.hWnd, EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return newError("SendMessage(EM_SETREADONLY)")
	}

	return nil
}

func (te *LogView) PostAppendText(value string){
	te.logChan <- value
	PostMessage(te.hWnd, TEM_APPENDTEXT, 0, 0)
}

func (te *LogView)  Write(p []byte)(int, error){
	te.PostAppendText(string(p) + "\r\n")
	return len(p), nil
}

func (te *LogView) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

		return DLGC_HASSETSEL | DLGC_WANTARROWS | DLGC_WANTCHARS
	case TEM_APPENDTEXT:
		select {
		case value := <- te.logChan:
			te.AppendText(value)
		default:
			return 0
		}
	}

	return te.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
