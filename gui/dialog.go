// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

type DialogCommandId int

const (
	DlgCmdOK       DialogCommandId = IDOK
	DlgCmdCancel   DialogCommandId = IDCANCEL
	DlgCmdAbort    DialogCommandId = IDABORT
	DlgCmdRetry    DialogCommandId = IDRETRY
	DlgCmdIgnore   DialogCommandId = IDIGNORE
	DlgCmdYes      DialogCommandId = IDYES
	DlgCmdNo       DialogCommandId = IDNO
	DlgCmdClose    DialogCommandId = IDCLOSE
	DlgCmdHelp     DialogCommandId = IDHELP
	DlgCmdTryAgain DialogCommandId = IDTRYAGAIN
	DlgCmdContinue DialogCommandId = IDCONTINUE
	DlgCmdTimeout  DialogCommandId = IDTIMEOUT
)

const dialogWindowClass = `\o/ Walk_Dialog_Class \o/`

var dialogWndProcCallback *syscall.Callback

func dialogWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	dlg, ok := widgetsByHWnd[msg.HWnd].(*Dialog)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return dlg.wndProc(msg, 0)
}

type Dialog struct {
	TopLevelWindow
}

func NewDialog() (*Dialog, os.Error) {
	ensureRegisteredWindowClass(dialogWindowClass, dialogWndProc, &dialogWndProcCallback)

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr(dialogWindowClass), nil,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 300, 0, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	d := &Dialog{TopLevelWindow: TopLevelWindow{Container: Container{Widget: Widget{hWnd: hWnd}}}}

	d.children = newObservedWidgetList(d)

	widgetsByHWnd[hWnd] = d

	// This forces display of focus rectangles, as soon as the user starts to type.
	SendMessage(hWnd, WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	return d, nil
}
