// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"github.com/lxn/win"
)

type CheckBox struct {
	Button
}

func NewCheckBox(parent Container) (*CheckBox, error) {
	cb := new(CheckBox)

	if err := InitWidget(
		cb,
		parent,
		"BUTTON",
		win.WS_TABSTOP|win.WS_VISIBLE|win.BS_AUTOCHECKBOX,
		0); err != nil {
		return nil, err
	}

	cb.Button.init()

	return cb, nil
}

func (*CheckBox) LayoutFlags() LayoutFlags {
	return 0
}

func (cb *CheckBox) MinSizeHint() Size {
	defaultSize := cb.dialogBaseUnitsToPixels(Size{50, 10})
	textSize := cb.calculateTextSizeImpl("n" + windowText(cb.hWnd))

	// FIXME: Use GetThemePartSize instead of GetSystemMetrics?
	w := textSize.Width + int(win.GetSystemMetrics(win.SM_CXMENUCHECK))
	h := maxi(defaultSize.Height, textSize.Height)

	return Size{w, h}
}

func (cb *CheckBox) SizeHint() Size {
	return cb.MinSizeHint()
}

func (cb *CheckBox) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		switch win.HIWORD(uint32(wParam)) {
		case win.BN_CLICKED:
			cb.checkedChangedPublisher.Publish()
		}
	}

	return cb.Button.WndProc(hwnd, msg, wParam, lParam)
}
