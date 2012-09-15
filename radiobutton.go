// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type RadioButton struct {
	Button
}

func NewRadioButton(parent Container) (*RadioButton, error) {
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

func (*RadioButton) LayoutFlags() LayoutFlags {
	return 0
}

func (rb *RadioButton) MinSizeHint() Size {
	defaultSize := rb.dialogBaseUnitsToPixels(Size{50, 10})
	textSize := rb.calculateTextSize()

	// FIXME: Use GetThemePartSize instead of GetSystemMetrics?
	w := textSize.Width + int(GetSystemMetrics(SM_CXMENUCHECK))
	h := maxi(defaultSize.Height, textSize.Height)

	return Size{w, h}
}

func (rb *RadioButton) SizeHint() Size {
	return rb.MinSizeHint()
}
