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
	. "walk/winapi"
	. "walk/winapi/user32"
)

var lineEditSubclassWndProcPtr uintptr
var lineEditOrigWndProcPtr uintptr

func lineEditSubclassWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	le, ok := widgetsByHWnd[hwnd].(*LineEdit)
	if !ok {
		return CallWindowProc(lineEditOrigWndProcPtr, hwnd, msg, wParam, lParam)
	}

	return le.wndProc(hwnd, msg, wParam, lParam, lineEditOrigWndProcPtr)
}

type LineEdit struct {
	WidgetBase
	editingFinishedPublisher EventPublisher
	returnPressedPublisher   EventPublisher
	textChanged              EventPublisher
}

func newLineEdit(parentHWND HWND) (*LineEdit, os.Error) {
	if lineEditSubclassWndProcPtr == 0 {
		lineEditSubclassWndProcPtr = syscall.NewCallback(lineEditSubclassWndProc)
	}

	hWnd := CreateWindowEx(
		WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("EDIT"), nil,
		ES_AUTOHSCROLL|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parentHWND, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	le := &LineEdit{WidgetBase: WidgetBase{hWnd: hWnd}}

	var succeeded bool
	defer func() {
		if !succeeded {
			le.Dispose()
		}
	}()

	lineEditOrigWndProcPtr = uintptr(SetWindowLong(hWnd, GWL_WNDPROC, int(lineEditSubclassWndProcPtr)))
	if lineEditOrigWndProcPtr == 0 {
		return nil, lastError("SetWindowLong")
	}

	le.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = le

	succeeded = true

	return le, nil
}

func NewLineEdit(parent Container) (*LineEdit, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	le, err := newLineEdit(parent.BaseWidget().hWnd)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			le.Dispose()
		}
	}()

	le.parent = parent
	if err = parent.Children().Add(le); err != nil {
		return nil, err
	}

	succeeded = true

	return le, nil
}

func (le *LineEdit) CueBanner() string {
	buf := make([]uint16, 128)
	if FALSE == SendMessage(le.hWnd, EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf))) {
		log.Print(newError("EM_GETCUEBANNER failed"))
		return ""
	}

	return syscall.UTF16ToString(buf)
}

func (le *LineEdit) SetCueBanner(value string) os.Error {
	if FALSE == SendMessage(le.hWnd, EM_SETCUEBANNER, FALSE, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value)))) {
		return newError("EM_SETCUEBANNER failed")
	}

	return nil
}

func (le *LineEdit) MaxLength() int {
	return int(SendMessage(le.hWnd, EM_GETLIMITTEXT, 0, 0))
}

func (le *LineEdit) SetMaxLength(value int) {
	SendMessage(le.hWnd, EM_LIMITTEXT, uintptr(value), 0)
}

func (le *LineEdit) TextSelection() (start, end int) {
	SendMessage(le.hWnd, EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (le *LineEdit) SetTextSelection(start, end int) {
	SendMessage(le.hWnd, EM_SETSEL, uintptr(start), uintptr(end))
}

func (le *LineEdit) PasswordMode() bool {
	return SendMessage(le.hWnd, EM_GETPASSWORDCHAR, 0, 0) != 0
}

func (le *LineEdit) SetPasswordMode(value bool) {
	SendMessage(le.hWnd, EM_SETPASSWORDCHAR, uintptr('*'), 0)
}

func (*LineEdit) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (le *LineEdit) PreferredSize() Size {
	return le.dialogBaseUnitsToPixels(Size{50, 14})
}

func (le *LineEdit) EditingFinished() *Event {
	return le.editingFinishedPublisher.Event()
}

func (le *LineEdit) ReturnPressed() *Event {
	return le.returnPressedPublisher.Event()
}

func (le *LineEdit) TextChanged() *Event {
	return le.textChanged.Event()
}

func (le *LineEdit) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint(wParam)) {
		case EN_CHANGE:
			le.textChanged.Publish()
		}

	case WM_GETDLGCODE:
		if root := rootWidget(le); root != nil {
			if dlg, ok := root.(dialogish); ok {
				if dlg.DefaultButton() != nil {
					// If the LineEdit lives in a Dialog that has a DefaultButton,
					// we won't swallow the return key. 
					break
				}
			}
		}

		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

	case WM_KEYDOWN:
		if wParam == VK_RETURN {
			le.returnPressedPublisher.Publish()
			le.editingFinishedPublisher.Publish()
		}

	case WM_KILLFOCUS:
		// FIXME: This may be dangerous, see remarks section:
		// http://msdn.microsoft.com/en-us/library/ms646282(v=vs.85).aspx
		le.editingFinishedPublisher.Publish()
	}

	return le.WidgetBase.wndProc(hwnd, msg, wParam, lParam, lineEditOrigWndProcPtr)
}
