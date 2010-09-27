// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
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
	SetLayout(value Layout)
}

type RootWidget interface {
	IContainer
	RunMessageLoop() os.Error
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

func (c *Container) SetLayout(value Layout) {
	if c.layout != value {
		if c.layout != nil {
			c.layout.SetContainer(nil)
		}

		c.layout = value

		if value != nil && value.Container() != IContainer(c) {
			value.SetContainer(c)
		}
	}
}

func (c *Container) wndProc(msg *MSG) uintptr {
	switch msg.Message {
	/*    case _WM_USER_NOTIFY:
	//        nmh := (*_NMHDR)(unsafe.Pointer(msg.lParam))
	        nmh := (*_NMITEMACTIVATE)(unsafe.Pointer(msg.lParam))
	        fmt.Printf("Container.raiseEvent: _WM_USER_NOTIFY: nmh: %+v\n", nmh)

	        if source, ok := eventSourcesByHWnd[nmh.hdr.hwndFrom]; ok {
	            // The widget that sent the message shall handle it itself.
	            source.raiseEvent(msg)
	        }*/

	case WM_COMMAND:
		hWnd := HWND(msg.LParam)
		if widget, ok := widgetsByHWnd[hWnd]; ok {
			if clickableWidget, ok := widget.(clickable); ok {
				clickableWidget.raiseClicked()
			}
			return 0
		} else {
			break
		}

	case WM_SIZE, WM_SIZING:
		if c.layout != nil {
			c.layout.Update(false)
		}
	}

	return c.Widget.wndProc(msg)
}

func (c *Container) onInsertingWidget(index int, widget IWidget) (err os.Error) {
	return nil
}

func (c *Container) onInsertedWidget(index int, widget IWidget) (err os.Error) {
	if widget.Parent().Handle() != c.hWnd {
		err = widget.SetParent(c)
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
	if widget.Parent() == IContainer(c) {
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
