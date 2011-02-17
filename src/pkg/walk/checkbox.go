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

type CheckBox struct {
	Button
}

func NewCheckBox(parent Container) (*CheckBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_AUTOCHECKBOX|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	cb := &CheckBox{
		Button: Button{
			WidgetBase: WidgetBase{
				hWnd:   hWnd,
				parent: parent,
			},
		},
	}

	succeeded := false
	defer func() {
		if !succeeded {
			cb.Dispose()
		}
	}()

	cb.SetFont(defaultFont)

	if err := parent.Children().Add(cb); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = cb

	succeeded = true

	return cb, nil
}

func (*CheckBox) LayoutFlags() LayoutFlags {
	return HShrink | HGrow
}

func (cb *CheckBox) PreferredSize() Size {
	return cb.dialogBaseUnitsToPixels(Size{50, 10})
}
