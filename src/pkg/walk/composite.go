// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

const compositeWindowClass = `\o/ Walk_Composite_Class \o/`

var compositeWindowWndProcPtr uintptr

func compositeWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	c, ok := widgetsByHWnd[hwnd].(*Composite)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return c.wndProc(hwnd, msg, wParam, lParam, 0)
}

type Composite struct {
	ContainerBase
}

func newCompositeWithStyle(parent IContainer, style uint) (*Composite, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(compositeWindowClass, compositeWndProc, &compositeWindowWndProcPtr)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(compositeWindowClass), nil,
		WS_CHILD|WS_VISIBLE|style,
		0, 0, 0, 0, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	c := &Composite{ContainerBase: ContainerBase{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}}}

	succeeded := false
	defer func() {
		if !succeeded {
			c.Dispose()
		}
	}()

	c.SetPersistent(true)

	c.children = newWidgetList(c)

	c.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = c

	if err := parent.Children().Add(c); err != nil {
		return nil, err
	}

	succeeded = true
	return c, nil
}

func NewComposite(parent IContainer) (*Composite, os.Error) {
	return newCompositeWithStyle(parent, 0)
}

func (c *Composite) LayoutFlags() LayoutFlags {
	var flags LayoutFlags

	count := c.children.Len()
	if count == 0 {
		return ShrinkHorz | ShrinkVert
	} else {
		for i := 0; i < count; i++ {
			flags |= c.children.At(i).LayoutFlags()
		}
	}

	return flags
}

func (c *Composite) PreferredSize() Size {
	var maxW, maxH int

	count := c.children.Len()
	for i := 0; i < count; i++ {
		prefSize := c.children.At(i).PreferredSize()
		if prefSize.Width > maxW {
			maxW = prefSize.Width
		}
		if prefSize.Height > maxH {
			maxH = prefSize.Height
		}
	}

	if c.layout != nil {
		marg := c.layout.Margins()
		maxW += marg.Left + marg.Right
		maxH += marg.Top + marg.Bottom
	}

	return Size{maxW, maxH}
}
