// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

type ProgressBar struct {
	Widget
}

func NewProgressBar(parent IContainer) (*ProgressBar, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("msctls_progress32"), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 80, 24, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	pb := &ProgressBar{Widget: Widget{hWnd: hWnd, parent: parent}}
	pb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = pb

	parent.Children().Add(pb)

	return pb, nil
}

func (*ProgressBar) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (pb *ProgressBar) PreferredSize() Size {
	return pb.dialogBaseUnitsToPixels(Size{50, 14})
}

func (pb *ProgressBar) ProgressPercent() int {
	return int(SendMessage(pb.hWnd, PBM_GETPOS, 0, 0))
}

func (pb *ProgressBar) SetProgressPercent(value int) {
	SendMessage(pb.hWnd, PBM_SETPOS, uintptr(value), 0)
}
