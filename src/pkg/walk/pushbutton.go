// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"unsafe"
)

import . "walk/winapi"

var pushButtonOrigWndProcPtr uintptr
var _ subclassedWidget = &PushButton{}

type PushButton struct {
	Button
}

func NewPushButton(parent Container) (*PushButton, os.Error) {
	pb := &PushButton{}

	if err := initChildWidget(
		pb,
		parent,
		"BUTTON",
		WS_TABSTOP|WS_VISIBLE|BS_PUSHBUTTON,
		0); err != nil {
		return nil, err
	}

	return pb, nil
}

func (*PushButton) origWndProcPtr() uintptr {
	return pushButtonOrigWndProcPtr
}

func (*PushButton) setOrigWndProcPtr(ptr uintptr) {
	pushButtonOrigWndProcPtr = ptr
}

func (*PushButton) LayoutFlags() LayoutFlags {
	return GrowableHorz
}

func (pb *PushButton) SizeHint() Size {
	var s Size

	SendMessage(pb.hWnd, BCM_GETIDEALSIZE, 0, uintptr(unsafe.Pointer(&s)))

	minSize := pb.dialogBaseUnitsToPixels(Size{50, 14})

	s.Width = maxi(s.Width, minSize.Width)
	s.Height = maxi(s.Height, minSize.Height)

	return s
}

func (pb *PushButton) ensureProperDialogDefaultButton(hwndFocus HWND) {
	widget := widgetFromHWND(hwndFocus)
	if widget == nil {
		return
	}

	if _, ok := widget.(*PushButton); ok {
		return
	}

	root := rootWidget(pb)
	if root == nil {
		return
	}

	dlg, ok := root.(dialogish)
	if !ok {
		return
	}

	defBtn := dlg.DefaultButton()
	if defBtn == nil {
		return
	}

	if err := defBtn.setAndClearStyleBits(BS_DEFPUSHBUTTON, BS_PUSHBUTTON); err != nil {
		return
	}

	if err := defBtn.Invalidate(); err != nil {
		return
	}
}

func (pb *PushButton) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_GETDLGCODE:
		hwndFocus := GetFocus()
		if hwndFocus == pb.hWnd {
			root := rootWidget(pb)
			if root == nil {
				break
			}

			dlg, ok := root.(dialogish)
			if !ok {
				break
			}

			defBtn := dlg.DefaultButton()
			if defBtn == pb {
				pb.setAndClearStyleBits(BS_DEFPUSHBUTTON, BS_PUSHBUTTON)
				return DLGC_BUTTON | DLGC_DEFPUSHBUTTON
			}

			break
		}

		pb.ensureProperDialogDefaultButton(hwndFocus)

	case WM_KILLFOCUS:
		pb.ensureProperDialogDefaultButton(HWND(wParam))
	}

	return pb.Button.wndProc(hwnd, msg, wParam, lParam)
}
