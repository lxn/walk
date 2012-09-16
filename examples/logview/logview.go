// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"errors"
	"unsafe"
	"syscall"
)

import . "github.com/lxn/go-winapi"
import (
	"github.com/lxn/walk"
)

type LogView struct {
	walk.WidgetBase
	logChan chan string
}

func NewLogView(parent walk.Container) (*LogView, error) {
	lc := make(chan string, 1024)
	te := &LogView{logChan:lc}

	if err := walk.InitChildWidget(
		te,
		parent,
		"EDIT",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|ES_MULTILINE|ES_WANTRETURN,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}
	te.setReadOnly(true)
	SendMessage(te.Handle(), EM_SETLIMITTEXT, 4294967295, 0)
	return te, nil
}

func (*LogView) LayoutFlags() walk.LayoutFlags {
	return walk.ShrinkableHorz | walk.ShrinkableVert | walk.GrowableHorz | walk.GrowableVert | walk.GreedyHorz | walk.GreedyVert
}

func (te *LogView) MinSizeHint() walk.Size {
	return walk.Size{20, 12}
}

func (te *LogView) SizeHint() walk.Size {
	return walk.Size{100, 100}
}

func (te *LogView) setTextSelection(start, end int) {
	SendMessage(te.Handle(), EM_SETSEL, uintptr(start), uintptr(end))
}

func (te *LogView) textLength() int{
	return int(SendMessage(te.Handle(), 0x000E, uintptr(0), uintptr(0)))
}

func (te *LogView) AppendText(value string) {
	textLength := te.textLength()
	te.setTextSelection(textLength, textLength)
	SendMessage(te.Handle(), EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))))
}

func (te *LogView) setReadOnly(readOnly bool) error {
	if 0 == SendMessage(te.Handle(), EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return errors.New("fail to call EM_SETREADONLY")
	}

	return nil
}

func (te *LogView) PostAppendText(value string){
	te.logChan <- value
	PostMessage(te.Handle(), TEM_APPENDTEXT, 0, 0)
}

func (te *LogView)  Write(p []byte)(int, error){
	te.PostAppendText(string(p) + "\r\n")
	return len(p), nil
}

func (te *LogView) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
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

	return te.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
