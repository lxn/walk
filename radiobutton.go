// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "github.com/lxn/go-winapi"
)

type RadioButton struct {
	Button
	group                *radioButtonGroup
	checkedDiscriminator interface{}
}

type radioButtonGroup struct {
	checkedButton *RadioButton
}

type radioButtonish interface {
	radioButton() *RadioButton
}

func NewRadioButton(parent Container) (*RadioButton, error) {
	rb := new(RadioButton)

	if err := InitChildWidget(
		rb,
		parent,
		"BUTTON",
		WS_TABSTOP|WS_VISIBLE|BS_AUTORADIOBUTTON,
		0); err != nil {
		return nil, err
	}

	rb.Button.init()

	if count := parent.Children().Len(); count > 1 {
		if prevRB, ok := parent.Children().At(count - 2).(radioButtonish); ok {
			rb.group = prevRB.radioButton().group
		}
	}
	if rb.group == nil {
		rb.group = new(radioButtonGroup)
	}

	rb.MustRegisterProperty("CheckedValue", NewProperty(
		func() interface{} {
			if rb.Checked() {
				return rb.checkedDiscriminator
			}

			return nil
		},
		func(v interface{}) error {
			checked := v == rb.checkedDiscriminator
			if checked {
				rb.group.checkedButton = rb
			}
			rb.SetChecked(checked)

			return nil
		},
		rb.clickedPublisher.Event()))

	return rb, nil
}

func (rb *RadioButton) radioButton() *RadioButton {
	return rb
}

func (*RadioButton) LayoutFlags() LayoutFlags {
	return 0
}

func (rb *RadioButton) MinSizeHint() Size {
	defaultSize := rb.dialogBaseUnitsToPixels(Size{50, 10})
	textSize := rb.calculateTextSizeImpl("n" + widgetText(rb.hWnd))

	// FIXME: Use GetThemePartSize instead of GetSystemMetrics?
	w := textSize.Width + int(GetSystemMetrics(SM_CXMENUCHECK))
	h := maxi(defaultSize.Height, textSize.Height)

	return Size{w, h}
}

func (rb *RadioButton) SizeHint() Size {
	return rb.MinSizeHint()
}

func (rb *RadioButton) CheckedDiscriminator() interface{} {
	return rb.checkedDiscriminator
}

func (rb *RadioButton) SetCheckedDiscriminator(checkedDiscriminator interface{}) {
	rb.checkedDiscriminator = checkedDiscriminator
}

func (rb *RadioButton) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case BN_CLICKED:
			prevChecked := rb.group.checkedButton
			rb.group.checkedButton = rb

			if prevChecked != rb {
				if prevChecked != nil {
					prevChecked.setChecked(false)
				}

				rb.setChecked(true)
			}
		}
	}

	return rb.Button.WndProc(hwnd, msg, wParam, lParam)
}
