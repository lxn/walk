// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

type PushButton struct {
	Button
}

func NewPushButton(parent IContainer) (*PushButton, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		/*BS_NOTIFY|*/ BS_PUSHBUTTON|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	pb := &PushButton{Button: Button{Widget: Widget{hWnd: hWnd, parent: parent}}}
	pb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = pb

	parent.Children().Add(pb)

	return pb, nil
}

func (*PushButton) LayoutFlags() LayoutFlags {
	return 0
}

func (pb *PushButton) PreferredSize() drawing.Size {
	var s drawing.Size

	SendMessage(pb.hWnd, BCM_GETIDEALSIZE, 0, uintptr(unsafe.Pointer(&s)))

	return s
}
