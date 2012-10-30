// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type clickable interface {
	raiseClicked()
}

type Button struct {
	WidgetBase
	checkedProperty  *Property
	clickedPublisher EventPublisher
	textProperty     *Property
}

func (b *Button) init() {
	b.checkedProperty = NewProperty(
		"Checked",
		func() interface{} {
			return b.Checked()
		},
		func(v interface{}) error {
			b.SetChecked(v.(bool))
			return nil
		},
		b.Clicked())

	b.textProperty = NewProperty(
		"Text",
		func() interface{} {
			return b.Text()
		},
		func(v interface{}) error {
			return b.SetText(v.(string))
		},
		nil)

	b.MustRegisterProperties(b.checkedProperty, b.textProperty)
}

func (b *Button) Text() string {
	return widgetText(b.hWnd)
}

func (b *Button) SetText(value string) error {
	if value == b.Text() {
		return nil
	}

	if err := setWidgetText(b.hWnd, value); err != nil {
		return err
	}

	return b.updateParentLayout()
}

func (b *Button) Checked() bool {
	return b.SendMessage(BM_GETCHECK, 0, 0) == BST_CHECKED
}

func (b *Button) SetChecked(value bool) {
	var chk uintptr

	if value {
		chk = BST_CHECKED
	} else {
		chk = BST_UNCHECKED
	}

	b.SendMessage(BM_SETCHECK, chk, 0)

	b.clickedPublisher.Publish()
}

func (b *Button) Clicked() *Event {
	return b.clickedPublisher.Event()
}

func (b *Button) raiseClicked() {
	b.clickedPublisher.Publish()
}

func (b *Button) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case BN_CLICKED:
			b.raiseClicked()
		}
	}

	return b.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
