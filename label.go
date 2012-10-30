// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type Label struct {
	WidgetBase
	bindingMember string
	textProperty  *Property
}

func NewLabel(parent Container) (*Label, error) {
	l := &Label{}

	if err := InitChildWidget(
		l,
		parent,
		"STATIC",
		WS_VISIBLE|SS_CENTERIMAGE,
		0); err != nil {
		return nil, err
	}

	l.textProperty = NewProperty(
		"Text",
		func() interface{} {
			return l.Text()
		},
		func(v interface{}) error {
			return l.SetText(v.(string))
		},
		nil)

	l.MustRegisterProperties(l.textProperty)

	return l, nil
}

func (*Label) LayoutFlags() LayoutFlags {
	return GrowableVert
}

func (l *Label) MinSizeHint() Size {
	return l.calculateTextSize()
}

func (l *Label) SizeHint() Size {
	return l.MinSizeHint()
}

func (l *Label) BindingMember() string {
	return l.bindingMember
}

func (l *Label) SetBindingMember(member string) error {
	if err := validateBindingMemberSyntax(member); err != nil {
		return err
	}

	l.bindingMember = member

	return nil
}

func (l *Label) BindingValue() interface{} {
	return l.Text()
}

func (l *Label) SetBindingValue(value interface{}) error {
	return l.SetText(value.(string))
}

func (l *Label) BindingValueChanged() *Event {
	return nil
}

func (l *Label) Text() string {
	return widgetText(l.hWnd)
}

func (l *Label) SetText(value string) error {
	if value == l.Text() {
		return nil
	}

	if err := setWidgetText(l.hWnd, value); err != nil {
		return err
	}

	return l.updateParentLayout()
}

func (l *Label) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_SIZE, WM_SIZING:
		l.Invalidate()
	}

	return l.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
