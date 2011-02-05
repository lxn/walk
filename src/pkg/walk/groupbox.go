// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

type GroupBox struct {
	Widget
}

func NewGroupBox(parent IContainer) (*GroupBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_GROUPBOX|WS_CHILD|WS_VISIBLE,
		0, 0, 80, 24, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	gb := &GroupBox{Widget: Widget{hWnd: hWnd, parent: parent}}
	gb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = gb

	parent.Children().Add(gb)

	return gb, nil
}

func (*GroupBox) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz | ShrinkVert | GrowVert
}

func (gb *GroupBox) PreferredSize() Size {
	return gb.dialogBaseUnitsToPixels(Size{100, 100})
}
