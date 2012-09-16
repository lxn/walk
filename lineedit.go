// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

type LineEdit struct {
	WidgetBase
	validator                Validator
	editingFinishedPublisher EventPublisher
	returnPressedPublisher   EventPublisher
	textChanged              EventPublisher
}

func newLineEdit(parent Widget) (*LineEdit, error) {
	le := &LineEdit{}

	if err := InitWidget(
		le,
		parent,
		"EDIT",
		WS_CHILD|WS_TABSTOP|WS_VISIBLE|ES_AUTOHSCROLL,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	return le, nil
}

func NewLineEdit(parent Container) (*LineEdit, error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	le, err := newLineEdit(parent)
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
	if FALSE == le.SendMessage(EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf))) {
		newError("EM_GETCUEBANNER failed")
		return ""
	}

	return syscall.UTF16ToString(buf)
}

func (le *LineEdit) SetCueBanner(value string) error {
	if FALSE == le.SendMessage(EM_SETCUEBANNER, FALSE, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value)))) {
		return newError("EM_SETCUEBANNER failed")
	}

	return nil
}

func (le *LineEdit) MaxLength() int {
	return int(le.SendMessage(EM_GETLIMITTEXT, 0, 0))
}

func (le *LineEdit) SetMaxLength(value int) {
	le.SendMessage(EM_LIMITTEXT, uintptr(value), 0)
}

func (le *LineEdit) Text() string {
	return widgetText(le.hWnd)
}

func (le *LineEdit) SetText(value string) error {
	return setWidgetText(le.hWnd, value)
}

func (le *LineEdit) TextSelection() (start, end int) {
	le.SendMessage(EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (le *LineEdit) SetTextSelection(start, end int) {
	le.SendMessage(EM_SETSEL, uintptr(start), uintptr(end))
}

func (le *LineEdit) PasswordMode() bool {
	return le.SendMessage(EM_GETPASSWORDCHAR, 0, 0) != 0
}

func (le *LineEdit) SetPasswordMode(value bool) {
	le.SendMessage(EM_SETPASSWORDCHAR, uintptr('*'), 0)
}

func (le *LineEdit) ReadOnly() bool {
	return le.hasStyleBits(ES_READONLY)
}

func (le *LineEdit) SetReadOnly(readOnly bool) error {
	if 0 == le.SendMessage(EM_SETREADONLY, uintptr(BoolToBOOL(readOnly)), 0) {
		return newError("SendMessage(EM_SETREADONLY)")
	}

	return nil
}

func (le *LineEdit) Validator() Validator {
	return le.validator
}

func (le *LineEdit) SetValidator(validator Validator) {
	le.validator = validator
}

func (*LineEdit) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz | GreedyHorz
}

func (le *LineEdit) MinSizeHint() Size {
	return le.dialogBaseUnitsToPixels(Size{20, 12})
}

func (le *LineEdit) SizeHint() Size {
	return le.dialogBaseUnitsToPixels(Size{50, 12})
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

func (le *LineEdit) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	/*	case WM_CHAR:
		if le.validator == nil {
			break
		}

		s := []uint16{uint16(wParam), 0}
		str := le.Text() + UTF16PtrToString(&s[0])

		if le.validator.Validate(str) == Invalid {
			return 0
		}*/

	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
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

	return le.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
