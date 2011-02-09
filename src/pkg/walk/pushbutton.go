// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

var pushButtonSubclassWndProcPtr uintptr
var pushButtonOrigWndProcPtr uintptr

func pushButtonSubclassWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	pb, ok := widgetsByHWnd[hwnd].(*PushButton)
	if !ok {
		return CallWindowProc(pushButtonOrigWndProcPtr, hwnd, msg, wParam, lParam)
	}

	return pb.wndProc(hwnd, msg, wParam, lParam, pushButtonOrigWndProcPtr)
}

type PushButton struct {
	Button
}

func NewPushButton(parent Container) (*PushButton, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	if pushButtonSubclassWndProcPtr == 0 {
		pushButtonSubclassWndProcPtr = syscall.NewCallback(pushButtonSubclassWndProc)
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_PUSHBUTTON|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	pb := &PushButton{Button: Button{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}}}

	succeeded := false
	defer func() {
		if !succeeded {
			pb.Dispose()
		}
	}()

	pushButtonOrigWndProcPtr = uintptr(SetWindowLong(hWnd, GWL_WNDPROC, int(pushButtonSubclassWndProcPtr)))
	if pushButtonOrigWndProcPtr == 0 {
		return nil, lastError("SetWindowLong")
	}

	pb.SetFont(defaultFont)

	if err := parent.Children().Add(pb); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = pb

	succeeded = true

	return pb, nil
}

func (*PushButton) LayoutFlags() LayoutFlags {
	return 0
}

func (pb *PushButton) PreferredSize() Size {
	var s Size

	SendMessage(pb.hWnd, BCM_GETIDEALSIZE, 0, uintptr(unsafe.Pointer(&s)))

	return s
}

func (pb *PushButton) ensureProperDialogDefaultButton(hwndFocus HWND) {
	widget, ok := widgetsByHWnd[hwndFocus]
	if !ok {
		return
	}

	if _, ok = widget.(*PushButton); ok {
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
		log.Print(err)
		return
	}

	if err := defBtn.Invalidate(); err != nil {
		log.Print(err)
		return
	}
}

func (pb *PushButton) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
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
				if err := pb.setAndClearStyleBits(BS_DEFPUSHBUTTON, BS_PUSHBUTTON); err != nil {
					log.Print(err)
				}
				return DLGC_BUTTON | DLGC_DEFPUSHBUTTON
			}

			break
		}

		pb.ensureProperDialogDefaultButton(hwndFocus)

	case WM_KILLFOCUS:
		pb.ensureProperDialogDefaultButton(HWND(wParam))
	}

	return pb.Button.wndProc(hwnd, msg, wParam, lParam, pushButtonOrigWndProcPtr)
}
