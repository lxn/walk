// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

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

var dialogWndProcPtr uintptr

func dialogWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	dlg, ok := widgetsByHWnd[hwnd].(*Dialog)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return dlg.wndProc(hwnd, msg, wParam, lParam, 0)
}

type Dialog struct {
	TopLevelWindow
	result        DialogCommandId
	defaultButton *PushButton
}

func NewDialog(owner RootWidget) (*Dialog, os.Error) {
	ensureRegisteredWindowClass(dialogWindowClass, dialogWndProc, &dialogWndProcPtr)

	var ownerHWnd HWND
	if owner != nil {
		ownerHWnd = owner.Handle()
	}

	hWnd := CreateWindowEx(
		WS_EX_DLGMODALFRAME, syscall.StringToUTF16Ptr(dialogWindowClass), nil,
		WS_CAPTION|WS_SYSMENU,
		-12345, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, ownerHWnd, 0, 0, nil)
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

	dlg.result = DlgCmdClose

	return dlg, nil
}

func (dlg *Dialog) DefaultButton() *PushButton {
	return dlg.defaultButton
}

func (dlg *Dialog) SetDefaultButton(button *PushButton) os.Error {
	if button != nil && !IsChild(dlg.hWnd, button.hWnd) {
		return newError("not a descendant of the dialog")
	}

	succeeded := false
	if dlg.defaultButton != nil {
		if err := dlg.defaultButton.setAndClearStyleBits(BS_PUSHBUTTON, BS_DEFPUSHBUTTON); err != nil {
			return err
		}
		defer func() {
			if !succeeded {
				dlg.defaultButton.setAndClearStyleBits(BS_DEFPUSHBUTTON, BS_PUSHBUTTON)
			}
		}()
	}

	if button != nil {
		if err := button.setAndClearStyleBits(BS_DEFPUSHBUTTON, BS_PUSHBUTTON); err != nil {
			return err
		}
	}

	dlg.defaultButton = button

	succeeded = true

	return nil
}

func (dlg *Dialog) Accept() {
	dlg.Close(DlgCmdOK)
}

func (dlg *Dialog) Cancel() {
	dlg.Close(DlgCmdCancel)
}

func (dlg *Dialog) Close(result DialogCommandId) {
	dlg.result = result

	dlg.TopLevelWindow.Close()
}

func (dlg *Dialog) Run() DialogCommandId {
	if dlg.owner != nil {
		ob := dlg.owner.Bounds()
		b := dlg.Bounds()
		if b.X == -12345 {
			dlg.SetBounds(Rectangle{
				ob.X + (ob.Width-b.Width)/2,
				ob.Y + (ob.Height-b.Height)/2,
				b.Width,
				b.Height,
			})
		}
	}

	dlg.Show()

	if dlg.owner != nil {
		dlg.owner.SetEnabled(false)
	}

	dlg.TopLevelWindow.Run()

	return dlg.result
}
