// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import . "walk/winapi"

var checkBoxOrigWndProcPtr uintptr
var _ subclassedWidget = &CheckBox{}

type CheckBox struct {
	Button
}

func NewCheckBox(parent Container) (*CheckBox, os.Error) {
	cb := &CheckBox{}

	if err := initChildWidget(
		cb,
		parent,
		"BUTTON",
		WS_TABSTOP|WS_VISIBLE|BS_AUTOCHECKBOX,
		0); err != nil {
		return nil, err
	}

	return cb, nil
}

func (*CheckBox) origWndProcPtr() uintptr {
	return checkBoxOrigWndProcPtr
}

func (*CheckBox) setOrigWndProcPtr(ptr uintptr) {
	checkBoxOrigWndProcPtr = ptr
}

func (*CheckBox) LayoutFlags() LayoutFlags {
	return 0
}

func (cb *CheckBox) SizeHint() Size {
	defaultSize := cb.dialogBaseUnitsToPixels(Size{50, 10})
	textSize := cb.calculateTextSize()

	// FIXME: Use GetThemePartSize instead of GetSystemMetrics?
	w := textSize.Width + int(GetSystemMetrics(SM_CXMENUCHECK))
	h := maxi(defaultSize.Height, textSize.Height)

	return Size{w, h}
}
