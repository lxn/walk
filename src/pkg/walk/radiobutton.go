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

type RadioButton struct {
	Button
}

func NewRadioButton(parent Container) (*RadioButton, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_AUTORADIOBUTTON /*|BS_NOTIFY*/ |WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	rb := &RadioButton{Button: Button{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}}}
	rb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = rb

	parent.Children().Add(rb)

	return rb, nil
}

func (*RadioButton) LayoutFlagsMask() LayoutFlags {
	return HShrink | HGrow
}

func (rb *RadioButton) PreferredSize() Size {
	return rb.dialogBaseUnitsToPixels(Size{50, 10})
}
