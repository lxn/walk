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
	. "walk/winapi"
	. "walk/winapi/user32"
)

var lineEditSubclassWndProcCallback *syscall.Callback
var lineEditSubclassWndProcPtr uintptr
var lineEditOrigWndProcPtr uintptr

func lineEditSubclassWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	le, ok := widgetsByHWnd[msg.HWnd].(*LineEdit)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return CallWindowProc(lineEditOrigWndProcPtr, msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return le.wndProc(msg, lineEditOrigWndProcPtr)
}

type LineEdit struct {
	Widget
}

func NewLineEdit(parent IContainer) (*LineEdit, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	if lineEditSubclassWndProcCallback == nil {
		lineEditSubclassWndProcCallback = syscall.NewCallback(lineEditSubclassWndProc, 4+4+4+4)
		lineEditSubclassWndProcPtr = uintptr(lineEditSubclassWndProcCallback.ExtFnEntry())
	}

	hWnd := CreateWindowEx(
		WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("EDIT"), nil,
		ES_AUTOHSCROLL|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	lineEditOrigWndProcPtr = uintptr(SetWindowLong(hWnd, GWL_WNDPROC, int(lineEditSubclassWndProcPtr)))
	if lineEditOrigWndProcPtr == 0 {
		return nil, lastError("SetWindowLong")
	}

	le := &LineEdit{Widget: Widget{hWnd: hWnd, parent: parent}}
	le.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = le

	parent.Children().Add(le)

	return le, nil
}

func (le *LineEdit) CueBanner() (string, os.Error) {
	buf := make([]uint16, 128)
	if FALSE == SendMessage(le.hWnd, EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf))) {
		return "", newError("EM_GETCUEBANNER failed")
	}

	return syscall.UTF16ToString(buf), nil
}

func (le *LineEdit) SetCueBanner(value string) os.Error {
	if FALSE == SendMessage(le.hWnd, EM_SETCUEBANNER, FALSE, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value)))) {
		return newError("EM_SETCUEBANNER failed")
	}

	return nil
}

func (*LineEdit) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (le *LineEdit) PreferredSize() drawing.Size {
	return le.dialogBaseUnitsToPixels(drawing.Size{50, 14})
}

func (le *LineEdit) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_GETDLGCODE:
		if msg.WParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}
	}

	return le.Widget.wndProc(msg, lineEditOrigWndProcPtr)
}
