// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

var radioButtonOrigWndProcPtr uintptr
var _ subclassedWidget = &RadioButton{}

type RadioButton struct {
	Button
}

func NewRadioButton(parent Container) (*RadioButton, os.Error) {
	rb := &RadioButton{}

	if err := initChildWidget(
		rb,
		parent,
		"BUTTON",
		WS_TABSTOP|WS_VISIBLE|BS_AUTORADIOBUTTON,
		0); err != nil {
		return nil, err
	}

	return rb, nil
}

func (*RadioButton) origWndProcPtr() uintptr {
	return radioButtonOrigWndProcPtr
}

func (*RadioButton) setOrigWndProcPtr(ptr uintptr) {
	radioButtonOrigWndProcPtr = ptr
}

func (*RadioButton) LayoutFlags() LayoutFlags {
	return 0
}

func (rb *RadioButton) PreferredSize() Size {
	return rb.dialogBaseUnitsToPixels(Size{50, 10})
}
