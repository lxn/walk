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

func NewDialog(owner *MainWindow) (*Dialog, os.Error) {
	ensureRegisteredWindowClass(dialogWindowClass, dialogWndProc, &dialogWndProcCallback)

	var ownerHWnd HWND
	if owner != nil {
		ownerHWnd = owner.hWnd
	}

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT|WS_EX_DLGMODALFRAME, syscall.StringToUTF16Ptr(dialogWindowClass), nil,
		WS_CAPTION|WS_POPUPWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 300, ownerHWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	dlg := &Dialog{
		TopLevelWindow: TopLevelWindow{
			Container: Container{
				Widget: Widget{
					hWnd: hWnd,
				},
			},
			owner: owner,
		},
	}

	dlg.children = newObservedWidgetList(dlg)

	widgetsByHWnd[hWnd] = dlg

	// This forces display of focus rectangles, as soon as the user starts to type.
	SendMessage(hWnd, WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	return dlg, nil
}

func (dlg *Dialog) RunMessageLoop() os.Error {
	if dlg.owner != nil {
		dlg.owner.SetEnabled(false)
		defer func() {
			dlg.owner.SetEnabled(true)
			SetWindowPos(dlg.owner.hWnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE)
		}()
	}

	return dlg.TopLevelWindow.RunMessageLoop()
}
