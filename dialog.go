// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

const (
	DlgCmdNone     = 0
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
	MustRegisterWindowClass(dialogWindowClass)
}

type dialogish interface {
	DefaultButton() *PushButton
	CancelButton() *PushButton
}

type Dialog struct {
	FormBase
	result               int
	defaultButton        *PushButton
	cancelButton         *PushButton
	centerInOwnerWhenRun bool
}

func NewDialog(owner Form) (*Dialog, error) {
	dlg := &Dialog{
		FormBase: FormBase{
			owner: owner,
		},
	}

	if err := InitWindow(
		dlg,
		owner,
		dialogWindowClass,
		WS_CAPTION|WS_SYSMENU|WS_THICKFRAME,
		WS_EX_DLGMODALFRAME); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			dlg.Dispose()
		}
	}()

	dlg.centerInOwnerWhenRun = owner != nil

	// This forces display of focus rectangles, as soon as the user starts to type.
	dlg.SendMessage(WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	dlg.result = DlgCmdNone

	succeeded = true

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

func (dlg *Dialog) Result() int {
	return dlg.result
}

func (dlg *Dialog) Accept() {
	dlg.Close(DlgCmdOK)
}

func (dlg *Dialog) Cancel() {
	dlg.Close(DlgCmdCancel)
}

func (dlg *Dialog) Close(result int) {
	dlg.result = result

	dlg.FormBase.Close()
}

func firstFocusableDescendantCallback(hwnd HWND, lParam uintptr) uintptr {
	widget := windowFromHandle(hwnd)

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

	EnumChildWindows(container.Handle(), firstFocusableDescendantCallbackPtr, uintptr(unsafe.Pointer(&hwnd)))

	return windowFromHandle(hwnd).(Widget)
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
		var size Size
		if layout := dlg.Layout(); layout != nil {
			size = layout.MinSize()
			min := dlg.MinSize()
			size.Width = maxi(size.Width, min.Width)
			size.Height = maxi(size.Height, min.Height)
		} else {
			size = dlg.Size()
		}

		ob := dlg.owner.Bounds()

		if dlg.centerInOwnerWhenRun {
			dlg.SetBounds(Rectangle{
				ob.X + (ob.Width-size.Width)/2,
				ob.Y + (ob.Height-size.Height)/2,
				size.Width,
				size.Height,
			})
		}
	} else {
		dlg.SetBounds(dlg.Bounds())
	}

	dlg.FormBase.Show()

	dlg.focusFirstCandidateDescendant()
}

func (dlg *Dialog) Run() int {
	dlg.Show()

	if dlg.owner != nil {
		dlg.owner.SetEnabled(false)
	}

	dlg.FormBase.Run()

	return dlg.result
}
