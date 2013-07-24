// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"errors"
	"syscall"
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
	"github.com/lxn/walk"
)

type LogView struct {
	walk.WidgetBase
	logChan chan string
}

const (
	TEM_APPENDTEXT = WM_USER + 6
)

func NewLogView(parent walk.Container) (*LogView, error) {
	lc := make(chan string, 1024)
	lv := &LogView{logChan: lc}

	if err := walk.InitWidget(
		lv,
		parent,
		"EDIT",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|ES_MULTILINE|ES_WANTRETURN,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}
	lv.setReadOnly(true)
	lv.SendMessage(EM_SETLIMITTEXT, 4294967295, 0)
	return lv, nil
}

func (*LogView) LayoutFlags() walk.LayoutFlags {
	return walk.ShrinkableHorz | walk.ShrinkableVert | walk.GrowableHorz | walk.GrowableVert | walk.GreedyHorz | walk.GreedyVert
}

func (*LogView) MinSizeHint() walk.Size {
	return walk.Size{20, 12}
}

func (*LogView) SizeHint() walk.Size {
	return walk.Size{100, 100}
}

func (lv *LogView) setTextSelection(start, end int) {
	lv.SendMessage(EM_SETSEL, uintptr(start), uintptr(end))
}

func (lv *LogView) textLength() int {
	return int(lv.SendMessage(0x000E, uintptr(0), uintptr(0)))
}

func (lv *LogView) AppendText(value string) {
	textLength := lv.textLength()
	lv.setTextSelection(textLength, textLength)
	lv.SendMessage(EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))))
}

func (lv *LogView) setReadOnly(readOnly bool) error {
	if 0 == lv.SendMessage(EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return errors.New("fail to call EM_SETREADONLY")
	}

	return nil
}

func (lv *LogView) PostAppendText(value string) {
	lv.logChan <- value
	PostMessage(lv.Handle(), TEM_APPENDTEXT, 0, 0)
}

func (lv *LogView) Write(p []byte) (int, error) {
	lv.PostAppendText(string(p) + "\r\n")
	return len(p), nil
}

func (lv *LogView) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

		return DLGC_HASSETSEL | DLGC_WANTARROWS | DLGC_WANTCHARS
	case TEM_APPENDTEXT:
		select {
		case value := <-lv.logChan:
			lv.AppendText(value)
		default:
			return 0
		}
	}

	return lv.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
