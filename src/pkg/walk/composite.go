// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

const compositeWindowClass = `\o/ Walk_Composite_Class \o/`

var compositeWindowClassRegistered bool

type Composite struct {
	ContainerBase
}

func newCompositeWithStyle(parent Container, style uint) (*Composite, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(compositeWindowClass, &compositeWindowClassRegistered)

	c := &Composite{}

	if err := initWidget(
		c,
		parent,
		compositeWindowClass,
		WS_CHILD|WS_VISIBLE|style,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			c.Dispose()
		}
	}()

	c.SetPersistent(true)

	c.children = newWidgetList(c)

	if parent.Children() == nil {
		// This may happen if the composite is (ab)used to implement 
		// Container semantics for some other widgets like GroupBox.
		if SetParent(c.hWnd, parent.BaseWidget().hWnd) == 0 {
			return nil, lastError("SetParent")
		}
	} else {
		if err := parent.Children().Add(c); err != nil {
			return nil, err
		}
	}

	succeeded = true

	return c, nil
}

func NewComposite(parent Container) (*Composite, os.Error) {
	return newCompositeWithStyle(parent, 0)
}

func (c *Composite) LayoutFlags() LayoutFlags {
	var flags LayoutFlags

	count := c.children.Len()
	if count == 0 {
		return HShrink | VShrink
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
		maxW += marg.HNear + marg.HFar
		maxH += marg.VNear + marg.VFar
	}

	return Size{maxW, maxH}
}
