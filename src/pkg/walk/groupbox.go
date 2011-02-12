// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"log"
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

const groupBoxWindowClass = `\o/ Walk_GroupBox_Class \o/`

var groupBoxWindowWndProcPtr uintptr

func groupBoxWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	c, ok := widgetsByHWnd[hwnd].(*GroupBox)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return c.wndProc(hwnd, msg, wParam, lParam, 0)
}

type GroupBox struct {
	WidgetBase
	hWndGroupBox HWND
	composite    *Composite
}

func NewGroupBox(parent Container) (*GroupBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(groupBoxWindowClass, groupBoxWndProc, &groupBoxWindowWndProcPtr)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(groupBoxWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 0, 0, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx(groupBoxWindowClass)")
	}

	gb := &GroupBox{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}}

	succeeded := false
	defer func() {
		if !succeeded {
			gb.Dispose()
		}
	}()

	gb.layoutFlags = HShrink | HGrow | VShrink | VGrow

	gb.hWndGroupBox = CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_GROUPBOX|WS_CHILD|WS_VISIBLE,
		0, 0, 80, 24, hWnd, 0, 0, nil)
	if gb.hWndGroupBox == 0 {
		return nil, lastError("CreateWindowEx(BUTTON)")
	}

	gb.SetFont(defaultFont)

	var err os.Error
	gb.composite, err = NewComposite(gb)
	if err != nil {
		return nil, err
	}

	if err := parent.Children().Add(gb); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = gb

	succeeded = true

	return gb, nil
}

func (gb *GroupBox) LayoutFlagsMask() LayoutFlags {
	return gb.composite.LayoutFlagsMask()
}

func (gb *GroupBox) PreferredSize() Size {
	cps := gb.composite.PreferredSize()
	wbcb := gb.WidgetBase.ClientBounds()
	gbcb := gb.ClientBounds()

	return Size{cps.Width + wbcb.Width - gbcb.Width, cps.Height + wbcb.Height - gbcb.Height}
}

func (gb *GroupBox) ClientBounds() Rectangle {
	cb := widgetClientBounds(gb.hWndGroupBox)

	// FIXME: Use appropriate margins
	return Rectangle{cb.X + 8, cb.Y + 24, cb.Width - 16, cb.Height - 32}
}

func (gb *GroupBox) SetFont(value *Font) {
	if value != gb.font {
		setWidgetFont(gb.hWndGroupBox, value)

		gb.font = value
	}
}

func (gb *GroupBox) Text() string {
	return widgetText(gb.hWndGroupBox)
}

func (gb *GroupBox) SetText(value string) os.Error {
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

func (gb *GroupBox) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_SIZE, WM_SIZING:
		wbcb := gb.WidgetBase.ClientBounds()
		if !MoveWindow(gb.hWndGroupBox, wbcb.X, wbcb.Y, wbcb.Width, wbcb.Height, true) {
			log.Print(lastError("MoveWindow"))
			break
		}

		gbcb := gb.ClientBounds()
		if err := gb.composite.SetBounds(gbcb); err != nil {
			log.Print(err)
		}
	}

	return gb.WidgetBase.wndProc(hwnd, msg, wParam, lParam, origWndProcPtr)
}
