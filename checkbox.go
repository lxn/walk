// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type CheckBox struct {
	Button
}

func NewCheckBox(parent Container) (*CheckBox, error) {
	cb := &CheckBox{}

	if err := InitChildWidget(
		cb,
		parent,
		"BUTTON",
		WS_TABSTOP|WS_VISIBLE|BS_AUTOCHECKBOX,
		0); err != nil {
		return nil, err
	}

	return cb, nil
}

func (*CheckBox) LayoutFlags() LayoutFlags {
	return 0
}

func (cb *CheckBox) MinSizeHint() Size {
	defaultSize := cb.dialogBaseUnitsToPixels(Size{50, 10})
	textSize := cb.calculateTextSize()

	// FIXME: Use GetThemePartSize instead of GetSystemMetrics?
	w := textSize.Width + int(GetSystemMetrics(SM_CXMENUCHECK))
	h := maxi(defaultSize.Height, textSize.Height)

	return Size{w, h}
}

func (cb *CheckBox) SizeHint() Size {
	return cb.MinSizeHint()
}
