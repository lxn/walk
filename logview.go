// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"unsafe"
	"syscall"
	"fmt"
)

import . "github.com/lxn/go-winapi"

var logViewOrigWndProcPtr uintptr
var _ subclassedWidget = &LogView{}

type LogView struct {
	WidgetBase
}

func NewLogView(parent Container) (*LogView, error) {
	te := &LogView{}

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
	fmt.Println("AppendText=", value)
	textLength := te.textLength()
	fmt.Println("TextLength=", textLength)
	te.setTextSelection(textLength, textLength)
	r := SendMessage(te.hWnd, EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))))
	fmt.Println("Ret=", r)
}

func (te *LogView) setReadOnly(readOnly bool) error {
	if 0 == SendMessage(te.hWnd, EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return newError("SendMessage(EM_SETREADONLY)")
	}

	return nil
}

func (te *LogView) PostAppendText(value string){
	PostMessage(te.hWnd, TEM_APPENDTEXT, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))), 0)
}

func (te *LogView)  Write(p []byte)(int, error){
	te.PostAppendText(string(p))
	return 0, nil
}

func (te *LogView) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

		return DLGC_HASSETSEL | DLGC_WANTARROWS | DLGC_WANTCHARS
	case TEM_APPENDTEXT:
		fmt.Println("Received APPEND_TEXT", wParam)
		te.AppendText(UTF16PtrToString((*uint16)(unsafe.Pointer(wParam))))
	}

	return te.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
