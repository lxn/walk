// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/user32"
)

type Margins struct {
	Left, Top, Right, Bottom int
}

type Layout interface {
	Container() IContainer
	SetContainer(value IContainer)
	Margins() *Margins
	SetMargins(value *Margins) os.Error
	Spacing() int
	SetSpacing(value int) os.Error
	Update(reset bool) os.Error
}

type IContainer interface {
	IWidget
	Children() *ObservedWidgetList
	Layout() Layout
	SetLayout(value Layout) os.Error
}

type RootWidget interface {
	IContainer
	RunMessageLoop() (int, os.Error)
}

type Container struct {
	Widget
	layout   Layout
	children *ObservedWidgetList
}

func (c *Container) Children() *ObservedWidgetList {
	return c.children
}

func (c *Container) Layout() Layout {
	return c.layout
}

func (c *Container) SetLayout(value Layout) os.Error {
	if c.layout != value {
		if c.layout != nil {
			c.layout.SetContainer(nil)
		}

		c.layout = value

		if value != nil && value.Container() != IContainer(c) {
			value.SetContainer(c)
		}
	}

	return nil
}

func (c *Container) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_COMMAND:
		switch HIWORD(uint(msg.WParam)) {
		case 0:
			hWnd := HWND(msg.LParam)
			if hWnd != 0 {
				if widget, ok := widgetsByHWnd[hWnd]; ok {
					if clickableWidget, ok := widget.(clickable); ok {
						clickableWidget.raiseClicked()
						return 0
					}
				}
			}

			// Menu
			actionId := uint16(LOWORD(uint(msg.WParam)))
			if action, ok := actionsById[actionId]; ok {
				action.raiseTriggered()
				return 0
			}

		case 1:
			// Accelerator
		}

	case WM_NOTIFY:
		nmh := (*NMHDR)(unsafe.Pointer(msg.LParam))
		if widget, ok := widgetsByHWnd[nmh.HwndFrom]; ok {
			// The widget that sent the notification shall handle it itself.
			widget.wndProc(msg, 0)
		}

	case WM_SIZE, WM_SIZING:
		if c.layout != nil {
			c.layout.Update(false)
		}
	}

	return c.Widget.wndProc(msg, origWndProcPtr)
}

func (c *Container) onInsertingWidget(index int, widget IWidget) (err os.Error) {
	return nil
}

func (c *Container) onInsertedWidget(index int, widget IWidget) (err os.Error) {
	if widget.Parent().Handle() != c.hWnd {
		err = widget.SetParent(widgetsByHWnd[c.hWnd].(IContainer))
		if err != nil {
			return
		}
	}

	if c.layout != nil {
		c.layout.Update(true)
	}

	return
}

func (c *Container) onRemovingWidget(index int, widget IWidget) (err os.Error) {
	if widget.Parent().Handle() == c.hWnd {
		err = widget.SetParent(nil)
	}

	return
}

func (c *Container) onRemovedWidget(index int, widget IWidget) (err os.Error) {
	if c.layout != nil {
		c.layout.Update(true)
	}

	return
}

func (c *Container) onClearingWidgets() (err os.Error) {
	panic("not implemented")
}

func (c *Container) onClearedWidgets() (err os.Error) {
	panic("not implemented")
}
