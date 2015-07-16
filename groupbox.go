// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
)

import (
	"github.com/lxn/win"
)

const groupBoxWindowClass = `\o/ Walk_GroupBox_Class \o/`

func init() {
	MustRegisterWindowClass(groupBoxWindowClass)
}

type GroupBox struct {
	WidgetBase
	hWndGroupBox          win.HWND
	composite             *Composite
	titleChangedPublisher EventPublisher
}

func NewGroupBox(parent Container) (*GroupBox, error) {
	gb := new(GroupBox)

	if err := InitWidget(
		gb,
		parent,
		groupBoxWindowClass,
		win.WS_VISIBLE,
		win.WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			gb.Dispose()
		}
	}()

	var err error
	gb.composite, err = NewComposite(gb)
	if err != nil {
		return nil, err
	}

	gb.hWndGroupBox, _ = win.CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		win.WS_CHILD|win.WS_VISIBLE|win.BS_GROUPBOX,
		0, 0, 80, 24, gb.hWnd, 0, 0, nil)
	if gb.hWndGroupBox == 0 {
		return nil, lastError("CreateWindowEx(BUTTON)")
	}

	// Set font to nil first to outsmart SetFont.
	gb.font = nil
	gb.SetFont(defaultFont)

	gb.MustRegisterProperty("Title", NewProperty(
		func() interface{} {
			return gb.Title()
		},
		func(v interface{}) error {
			return gb.SetTitle(v.(string))
		},
		gb.titleChangedPublisher.Event()))

	succeeded = true

	return gb, nil
}

func (gb *GroupBox) AsContainerBase() *ContainerBase {
	return gb.composite.AsContainerBase()
}

func (gb *GroupBox) LayoutFlags() LayoutFlags {
	if gb.composite == nil {
		return 0
	}

	return gb.composite.LayoutFlags()
}

func (gb *GroupBox) MinSizeHint() Size {
	if gb.composite == nil {
		return Size{100, 100}
	}

	cmsh := gb.composite.MinSizeHint()

	return Size{cmsh.Width + 2, cmsh.Height + 9}
}

func (gb *GroupBox) SizeHint() Size {
	return gb.MinSizeHint()
}

func (gb *GroupBox) ClientBounds() Rectangle {
	cb := windowClientBounds(gb.hWndGroupBox)

	if gb.Layout() == nil {
		return cb
	}

	// FIXME: Use appropriate margins
	return Rectangle{cb.X + 1, cb.Y + 14, cb.Width - 2, cb.Height - 9}
}

func (gb *GroupBox) SetFont(value *Font) {
	if value != gb.font {
		setWindowFont(gb.hWndGroupBox, value)

		gb.font = value

		gb.composite.SetFont(value)
	}
}

func (gb *GroupBox) SetSuspended(suspend bool) {
	gb.composite.SetSuspended(suspend)
	gb.WidgetBase.SetSuspended(suspend)
	gb.Invalidate()
}

func (gb *GroupBox) DataBinder() *DataBinder {
	return gb.composite.dataBinder
}

func (gb *GroupBox) SetDataBinder(dataBinder *DataBinder) {
	gb.composite.SetDataBinder(dataBinder)
}

func (gb *GroupBox) Title() string {
	return windowText(gb.hWndGroupBox)
}

func (gb *GroupBox) SetTitle(value string) error {
	return setWindowText(gb.hWndGroupBox, value)
}

func (gb *GroupBox) Children() *WidgetList {
	if gb.composite == nil {
		// Without this we would get into trouble in NewComposite.
		return nil
	}

	return gb.composite.Children()
}

func (gb *GroupBox) Layout() Layout {
	return gb.composite.Layout()
}

func (gb *GroupBox) SetLayout(value Layout) error {
	return gb.composite.SetLayout(value)
}

func (gb *GroupBox) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if gb.composite != nil {
		switch msg {
		case win.WM_COMMAND, win.WM_NOTIFY:
			gb.composite.WndProc(hwnd, msg, wParam, lParam)

		case win.WM_SETTEXT:
			gb.titleChangedPublisher.Publish()

		case win.WM_SIZE, win.WM_SIZING:
			wbcb := gb.WidgetBase.ClientBounds()
			if !win.MoveWindow(
				gb.hWndGroupBox,
				int32(wbcb.X),
				int32(wbcb.Y),
				int32(wbcb.Width),
				int32(wbcb.Height),
				true) {

				lastError("MoveWindow")
				break
			}

			gbcb := gb.ClientBounds()
			gb.composite.SetBounds(gbcb)
		}
	}

	return gb.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
