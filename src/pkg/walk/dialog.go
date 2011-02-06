// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"log"
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
		ownerHWnd = owner.BaseWidget().hWnd
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
			ContainerBase: ContainerBase{
				WidgetBase: WidgetBase{
					hWnd: hWnd,
				},
			},
			owner: owner,
		},
	}

	dlg.children = newWidgetList(dlg)

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

func firstFocusableDescendant(container Container) Widget {
	for _, widget := range container.Children().items {
		if !widget.Visible() || !widget.Enabled() {
			continue
		}

		if c, ok := widget.(Container); ok {
			if w := firstFocusableDescendant(c); w != nil {
				return w
			}
		} else {
			style := uint(GetWindowLong(widget.BaseWidget().hWnd, GWL_STYLE))
			// FIXME: Ugly workaround for NumberEdit
			_, isTextSelectable := widget.(textSelectable)
			if style&WS_TABSTOP > 0 || isTextSelectable {
				return widget
			}
		}
	}

	return nil
}

type textSelectable interface {
	SetTextSelection(start, end int)
}

func (dlg *Dialog) Show() {
	dlg.TopLevelWindow.Show()

	widget := firstFocusableDescendant(dlg)
	if widget == nil {
		return
	}

	if err := widget.SetFocus(); err != nil {
		log.Print(err)
		return
	}

	if textSel, ok := widget.(textSelectable); ok {
		textSel.SetTextSelection(0, -1)
	}
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
