// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import . "walk/winapi"

const compositeWindowClass = `\o/ Walk_Composite_Class \o/`

var compositeWindowClassRegistered bool

type Composite struct {
	ContainerBase
}

func newCompositeWithStyle(parent Container, style uint) (*Composite, os.Error) {
	ensureRegisteredWindowClass(compositeWindowClass, &compositeWindowClassRegistered)

	c := &Composite{}
	c.children = newWidgetList(c)
	c.SetPersistent(true)

	if err := initChildWidget(
		c,
		parent,
		compositeWindowClass,
		WS_CHILD|WS_VISIBLE|style,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	return c, nil
}

func NewComposite(parent Container) (*Composite, os.Error) {
	return newCompositeWithStyle(parent, 0)
}
