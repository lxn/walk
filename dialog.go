// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"unsafe"

	"github.com/lxn/win"
)

const (
	DlgCmdNone     = 0
	DlgCmdOK       = win.IDOK
	DlgCmdCancel   = win.IDCANCEL
	DlgCmdAbort    = win.IDABORT
	DlgCmdRetry    = win.IDRETRY
	DlgCmdIgnore   = win.IDIGNORE
	DlgCmdYes      = win.IDYES
	DlgCmdNo       = win.IDNO
	DlgCmdClose    = win.IDCLOSE
	DlgCmdHelp     = win.IDHELP
	DlgCmdTryAgain = win.IDTRYAGAIN
	DlgCmdContinue = win.IDCONTINUE
	DlgCmdTimeout  = win.IDTIMEOUT
)

const dialogWindowClass = `\o/ Walk_Dialog_Class \o/`

func init() {
	AppendToWalkInit(func() {
		MustRegisterWindowClass(dialogWindowClass)
	})
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
	return newDialogWithStyle(owner, win.WS_THICKFRAME)
}

func NewDialogWithFixedSize(owner Form) (*Dialog, error) {
	return newDialogWithStyle(owner, 0)
}

func newDialogWithStyle(owner Form, style uint32) (*Dialog, error) {
	dlg := &Dialog{
		FormBase: FormBase{
			owner: owner,
		},
	}

	if err := InitWindow(
		dlg,
		owner,
		dialogWindowClass,
		win.WS_CAPTION|win.WS_SYSMENU|style,
		0); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			dlg.Dispose()
		}
	}()

	dlg.centerInOwnerWhenRun = owner != nil

	dlg.result = DlgCmdNone

	succeeded = true

	return dlg, nil
}

func (dlg *Dialog) DefaultButton() *PushButton {
	return dlg.defaultButton
}

func (dlg *Dialog) SetDefaultButton(button *PushButton) error {
	if button != nil && !win.IsChild(dlg.hWnd, button.hWnd) {
		return newError("not a descendant of the dialog")
	}

	succeeded := false
	if dlg.defaultButton != nil {
		if err := dlg.defaultButton.setAndClearStyleBits(win.BS_PUSHBUTTON, win.BS_DEFPUSHBUTTON); err != nil {
			return err
		}
		defer func() {
			if !succeeded {
				dlg.defaultButton.setAndClearStyleBits(win.BS_DEFPUSHBUTTON, win.BS_PUSHBUTTON)
			}
		}()
	}

	if button != nil {
		if err := button.setAndClearStyleBits(win.BS_DEFPUSHBUTTON, win.BS_PUSHBUTTON); err != nil {
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
	if button != nil && !win.IsChild(dlg.hWnd, button.hWnd) {
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

func (dlg *Dialog) Show() {
	var willRestore bool
	if dlg.Persistent() {
		state, _ := dlg.ReadState()
		willRestore = state != ""
	}

	if !willRestore {
		var size Size
		if layout := dlg.Layout(); layout != nil {
			size = maxSize(dlg.clientComposite.MinSizeHint(), dlg.MinSizePixels())
		} else {
			size = dlg.SizePixels()
		}

		if dlg.owner != nil {
			ob := dlg.owner.BoundsPixels()

			if dlg.centerInOwnerWhenRun {
				dlg.SetBoundsPixels(fitRectToScreen(dlg.hWnd, Rectangle{
					ob.X + (ob.Width-size.Width)/2,
					ob.Y + (ob.Height-size.Height)/2,
					size.Width,
					size.Height,
				}))
			}
		} else {
			b := dlg.BoundsPixels()

			dlg.SetBoundsPixels(Rectangle{b.X, b.Y, size.Width, size.Height})
		}
	}

	dlg.FormBase.Show()

	dlg.startLayout()
}

// fitRectToScreen fits rectangle to screen. Input and output rectangles are in native pixels.
func fitRectToScreen(hWnd win.HWND, r Rectangle) Rectangle {
	var mi win.MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))

	if !win.GetMonitorInfo(win.MonitorFromWindow(
		hWnd, win.MONITOR_DEFAULTTOPRIMARY), &mi) {

		return r
	}

	mon := rectangleFromRECT(mi.RcWork)

	dpi := win.GetDpiForWindow(hWnd)
	mon.Height -= int(win.GetSystemMetricsForDpi(win.SM_CYCAPTION, dpi))

	if r.Width <= mon.Width {
		switch {
		case r.X < mon.X:
			r.X = mon.X
		case r.X+r.Width > mon.X+mon.Width:
			r.X = mon.X + mon.Width - r.Width
		}
	}

	if r.Height <= mon.Height {
		switch {
		case r.Y < mon.Y:
			r.Y = mon.Y
		case r.Y+r.Height > mon.Y+mon.Height:
			r.Y = mon.Y + mon.Height - r.Height
		}
	}

	return r
}

func (dlg *Dialog) Run() int {
	dlg.Show()

	dlg.FormBase.Run()

	return dlg.result
}

func (dlg *Dialog) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		if win.HIWORD(uint32(wParam)) == 0 {
			switch win.LOWORD(uint32(wParam)) {
			case DlgCmdOK:
				if dlg.defaultButton != nil {
					dlg.defaultButton.raiseClicked()
				}

			case DlgCmdCancel:
				if dlg.cancelButton != nil {
					dlg.cancelButton.raiseClicked()
				}
			}
		}
	}

	return dlg.FormBase.WndProc(hwnd, msg, wParam, lParam)
}
