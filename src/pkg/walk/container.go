// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

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
	Container() Container
	SetContainer(value Container)
	Margins() *Margins
	SetMargins(value *Margins) os.Error
	Spacing() int
	SetSpacing(value int) os.Error
	Update(reset bool) os.Error
}

type Container interface {
	Widget
	Children() *WidgetList
	Layout() Layout
	SetLayout(value Layout) os.Error
}

type RootWidget interface {
	Container
	Run() int
}

type ContainerBase struct {
	WidgetBase
	layout     Layout
	children   *WidgetList
	persistent bool
}

func (c *ContainerBase) Children() *WidgetList {
	return c.children
}

func (c *ContainerBase) Layout() Layout {
	return c.layout
}

func (c *ContainerBase) SetLayout(value Layout) os.Error {
	if c.layout != value {
		if c.layout != nil {
			c.layout.SetContainer(nil)
		}

		c.layout = value

		if value != nil && value.Container() != Container(c) {
			value.SetContainer(c)
		}
	}

	return nil
}

func (c *ContainerBase) forEachPersistableChild(f func(p Persistable) os.Error) os.Error {
	for _, child := range c.children.items {
		if persistable, ok := child.(Persistable); ok && persistable.Persistent() {
			if err := f(persistable); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ContainerBase) Persistent() bool {
	return c.persistent
}

func (c *ContainerBase) SetPersistent(value bool) {
	c.persistent = value
}

func (c *ContainerBase) SaveState() os.Error {
	return c.forEachPersistableChild(func(p Persistable) os.Error {
		return p.SaveState()
	})
}

func (c *ContainerBase) RestoreState() os.Error {
	return c.forEachPersistableChild(func(p Persistable) os.Error {
		return p.RestoreState()
	})
}

func (c *ContainerBase) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint(wParam)) {
		case 0:
			hWnd := HWND(lParam)
			if hWnd != 0 {
				if widget, ok := widgetsByHWnd[hWnd]; ok {
					if clickableWidget, ok := widget.(clickable); ok {
						clickableWidget.raiseClicked()
						return 0
					}
				}
			}

			// Menu
			actionId := uint16(LOWORD(uint(wParam)))
			if action, ok := actionsById[actionId]; ok {
				action.raiseTriggered()
				return 0
			}

		case 1:
			// Accelerator

		default:
			// The widget that sent the notification shall handle it itself.
			hWnd := HWND(lParam)
			if widget, ok := widgetsByHWnd[hWnd]; ok {
				widget.wndProc(hwnd, msg, wParam, lParam, 0)
				return 0
			}
		}

	case WM_NOTIFY:
		nmh := (*NMHDR)(unsafe.Pointer(lParam))
		if widget, ok := widgetsByHWnd[nmh.HwndFrom]; ok {
			// The widget that sent the notification shall handle it itself.
			widget.wndProc(hwnd, msg, wParam, lParam, 0)
		}

	case WM_SIZE, WM_SIZING:
		if c.layout != nil {
			c.layout.Update(false)
		}
	}

	return c.WidgetBase.wndProc(hwnd, msg, wParam, lParam, origWndProcPtr)
}

func (c *ContainerBase) onInsertingWidget(index int, widget Widget) (err os.Error) {
	return nil
}

func (c *ContainerBase) onInsertedWidget(index int, widget Widget) (err os.Error) {
	if widget.Parent().BaseWidget().hWnd != c.hWnd {
		err = widget.SetParent(widgetsByHWnd[c.hWnd].(Container))
		if err != nil {
			return
		}
	}

	if c.layout != nil {
		c.layout.Update(true)
	}

	return
}

func (c *ContainerBase) onRemovingWidget(index int, widget Widget) (err os.Error) {
	if widget.Parent().BaseWidget().hWnd == c.hWnd {
		err = widget.SetParent(nil)
	}

	return
}

func (c *ContainerBase) onRemovedWidget(index int, widget Widget) (err os.Error) {
	if c.layout != nil {
		c.layout.Update(true)
	}

	return
}

func (c *ContainerBase) onClearingWidgets() (err os.Error) {
	panic("not implemented")
}

func (c *ContainerBase) onClearedWidgets() (err os.Error) {
	panic("not implemented")
}
