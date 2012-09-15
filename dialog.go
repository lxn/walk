// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

const (
	DlgCmdOK       = IDOK
	DlgCmdCancel   = IDCANCEL
	DlgCmdAbort    = IDABORT
	DlgCmdRetry    = IDRETRY
	DlgCmdIgnore   = IDIGNORE
	DlgCmdYes      = IDYES
	DlgCmdNo       = IDNO
	DlgCmdClose    = IDCLOSE
	DlgCmdHelp     = IDHELP
	DlgCmdTryAgain = IDTRYAGAIN
	DlgCmdContinue = IDCONTINUE
	DlgCmdTimeout  = IDTIMEOUT
)

const dialogWindowClass = `\o/ Walk_Dialog_Class \o/`

func init() {
	mustRegisterWindowClass(dialogWindowClass)
}

type dialogish interface {
	DefaultButton() *PushButton
	CancelButton() *PushButton
}

type Dialog struct {
	TopLevelWindow
	result               int
	defaultButton        *PushButton
	cancelButton         *PushButton
	centerInOwnerWhenRun bool
}

func NewDialog(owner RootWidget) (*Dialog, error) {
	dlg := &Dialog{
		TopLevelWindow: TopLevelWindow{
			owner: owner,
		},
	}

	if err := initWidget(
		dlg,
		owner,
		dialogWindowClass,
		WS_CAPTION|WS_SYSMENU|WS_THICKFRAME,
		WS_EX_DLGMODALFRAME); err != nil {
		return nil, err
	}

	dlg.centerInOwnerWhenRun = owner != nil

	dlg.children = newWidgetList(dlg)

	// This forces display of focus rectangles, as soon as the user starts to type.
	SendMessage(dlg.hWnd, WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	dlg.result = DlgCmdClose

	return dlg, nil
}

func (dlg *Dialog) DefaultButton() *PushButton {
	return dlg.defaultButton
}

func (dlg *Dialog) SetDefaultButton(button *PushButton) error {
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

func (dlg *Dialog) CancelButton() *PushButton {
	return dlg.cancelButton
}

func (dlg *Dialog) SetCancelButton(button *PushButton) error {
	if button != nil && !IsChild(dlg.hWnd, button.hWnd) {
		return newError("not a descendant of the dialog")
	}

	dlg.cancelButton = button

	return nil
}

func (dlg *Dialog) Accept() {
	dlg.Close(DlgCmdOK)
}

func (dlg *Dialog) Cancel() {
	dlg.Close(DlgCmdCancel)
}

func (dlg *Dialog) Close(result int) {
	dlg.result = result

	dlg.TopLevelWindow.Close()
}

func firstFocusableDescendantCallback(hwnd HWND, lParam uintptr) uintptr {
	widget := widgetFromHWND(hwnd)

	if widget == nil || !widget.Visible() || !widget.Enabled() {
		return 1
	}

	style := uint(GetWindowLong(hwnd, GWL_STYLE))
	// FIXME: Ugly workaround for NumberEdit
	_, isTextSelectable := widget.(textSelectable)
	if style&WS_TABSTOP > 0 || isTextSelectable {
		hwndPtr := (*HWND)(unsafe.Pointer(lParam))
		*hwndPtr = hwnd
		return 0
	}

	return 1
}

var firstFocusableDescendantCallbackPtr = syscall.NewCallback(firstFocusableDescendantCallback)

func firstFocusableDescendant(container Container) Widget {
	var hwnd HWND

	EnumChildWindows(container.BaseWidget().hWnd, firstFocusableDescendantCallbackPtr, uintptr(unsafe.Pointer(&hwnd)))

	return widgetFromHWND(hwnd)
}

type textSelectable interface {
	SetTextSelection(start, end int)
}

func (dlg *Dialog) focusFirstCandidateDescendant() {
	widget := firstFocusableDescendant(dlg)
	if widget == nil {
		return
	}

	if err := widget.SetFocus(); err != nil {
		return
	}

	if textSel, ok := widget.(textSelectable); ok {
		textSel.SetTextSelection(0, -1)
	}
}

func (dlg *Dialog) Show() {
	if dlg.owner != nil {
		ob := dlg.owner.Bounds()
		b := dlg.Bounds()
		if dlg.centerInOwnerWhenRun {
			dlg.SetBounds(Rectangle{
				ob.X + (ob.Width-b.Width)/2,
				ob.Y + (ob.Height-b.Height)/2,
				b.Width,
				b.Height,
			})
		}
	} else {
		dlg.SetBounds(dlg.Bounds())
	}

	dlg.TopLevelWindow.Show()

	dlg.focusFirstCandidateDescendant()
}

func (dlg *Dialog) Run() int {
	dlg.Show()

	if dlg.owner != nil {
		dlg.owner.SetEnabled(false)
	}

	dlg.TopLevelWindow.Run()

	return dlg.result
}
