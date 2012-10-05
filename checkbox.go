// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type CheckBox struct {
	Button
	bindingMember string
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

func (cb *CheckBox) BindingMember() string {
	return cb.bindingMember
}

func (cb *CheckBox) SetBindingMember(member string) error {
	if err := validateBindingMemberSyntax(member); err != nil {
		return err
	}

	cb.bindingMember = member

	return nil
}

func (cb *CheckBox) BindingValue() interface{} {
	return cb.Checked()
}

func (cb *CheckBox) SetBindingValue(value interface{}) error {
	cb.SetChecked(value.(bool))

	return nil
}

func (cb *CheckBox) BindingValueChanged() *Event {
	return cb.Clicked()
}
