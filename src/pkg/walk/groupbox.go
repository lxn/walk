// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import . "walk/winapi"

const groupBoxWindowClass = `\o/ Walk_GroupBox_Class \o/`

var groupBoxWindowClassRegistered bool

type GroupBox struct {
	WidgetBase
	hWndGroupBox HWND
	composite    *Composite
}

func NewGroupBox(parent Container) (*GroupBox, os.Error) {
	ensureRegisteredWindowClass(groupBoxWindowClass, &groupBoxWindowClassRegistered)

	gb := &GroupBox{}

	if err := initChildWidget(
		gb,
		parent,
		groupBoxWindowClass,
		WS_VISIBLE,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			gb.Dispose()
		}
	}()

	gb.hWndGroupBox = CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		WS_CHILD|WS_VISIBLE|BS_GROUPBOX,
		0, 0, 80, 24, gb.hWnd, 0, 0, nil)
	if gb.hWndGroupBox == 0 {
		return nil, lastError("CreateWindowEx(BUTTON)")
	}

	var err os.Error
	gb.composite, err = NewComposite(gb)
	if err != nil {
		return nil, err
	}

	// Set font to nil first to outsmart SetFont.
	gb.font = nil
	gb.SetFont(defaultFont)

	succeeded = true

	return gb, nil
}

func (gb *GroupBox) LayoutFlags() LayoutFlags {
	if gb.composite == nil {
		return 0
	}

	return gb.composite.LayoutFlags()
}

func (gb *GroupBox) MinSizeHint() Size {
	if gb.composite == nil {
		return gb.SizeHint()
	}

	cps := gb.composite.MinSizeHint()
	wbcb := gb.WidgetBase.ClientBounds()
	gbcb := gb.ClientBounds()

	return Size{cps.Width + wbcb.Width - gbcb.Width, cps.Height + wbcb.Height - gbcb.Height}
}

func (gb *GroupBox) SizeHint() Size {
	return Size{100, 100}
}

func (gb *GroupBox) ClientBounds() Rectangle {
	cb := widgetClientBounds(gb.hWndGroupBox)

	if gb.Layout() == nil {
		return cb
	}

	// FIXME: Use appropriate margins
	return Rectangle{cb.X + 1, cb.Y + 14, cb.Width - 2, cb.Height - 9}
}

func (gb *GroupBox) SetFont(value *Font) {
	if value != gb.font {
		setWidgetFont(gb.hWndGroupBox, value)

		gb.font = value
	}
}

func (gb *GroupBox) Title() string {
	return widgetText(gb.hWndGroupBox)
}

func (gb *GroupBox) SetTitle(value string) os.Error {
	return setWidgetText(gb.hWndGroupBox, value)
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

func (gb *GroupBox) SetLayout(value Layout) os.Error {
	return gb.composite.SetLayout(value)
}

func (gb *GroupBox) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if gb.composite != nil {
		switch msg {
		case WM_COMMAND, WM_NOTIFY:
			gb.composite.wndProc(hwnd, msg, wParam, lParam)

		case WM_SIZE, WM_SIZING:
			wbcb := gb.WidgetBase.ClientBounds()
			if !MoveWindow(
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

	return gb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
