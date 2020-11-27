// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/xlplbo/walk"
	"github.com/xlplbo/win"
)

type LogView struct {
	walk.TextEdit
	logChan chan string
	maxSize int
}

const (
	TEM_APPENDTEXT = win.WM_USER + 6
)

func NewLogView(parent walk.Container) (*LogView, error) {
	lv := &LogView{
		logChan: make(chan string, 1024),
		maxSize: 0x7FFF0000,
	}
	if err := walk.InitWidget(lv,
		parent,
		"EDIT",
		win.WS_TABSTOP|win.WS_VISIBLE|win.WS_VSCROLL|win.ES_MULTILINE|win.ES_WANTRETURN,
		win.WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}
	lv.TextEdit.SetReadOnly(true)
	lv.TextEdit.SetMaxLength(lv.maxSize)
	return lv, nil
}

func (lv *LogView) AppendText(value string) {
	textLength := lv.TextEdit.TextLength()
	if textLength+len(value) < lv.maxSize {
		lv.TextEdit.AppendText(value)
	} else {
		lv.TextEdit.SetText(value)
	}
}

func (lv *LogView) PostAppendText(value string) {
	lv.logChan <- value
	win.PostMessage(lv.Handle(), TEM_APPENDTEXT, 0, 0)
}

func (lv *LogView) Write(p []byte) (int, error) {
	lv.PostAppendText(string(p) + "\r\n")
	return len(p), nil
}

func (lv *LogView) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case TEM_APPENDTEXT:
		select {
		case value := <-lv.logChan:
			lv.AppendText(value)
		default:
		}
	}
	return lv.TextEdit.WndProc(hwnd, msg, wParam, lParam)
}
