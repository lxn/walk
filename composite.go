// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

const compositeWindowClass = `\o/ Walk_Composite_Class \o/`

func init() {
	MustRegisterWindowClass(compositeWindowClass)
}

type Composite struct {
	ContainerBase
}

func newCompositeWithStyle(parent Widget, style uint32) (*Composite, error) {
	c := &Composite{}
	c.children = newWidgetList(c)
	c.SetPersistent(true)

	if err := InitChildWidget(
		c,
		parent,
		compositeWindowClass,
		WS_CHILD|WS_VISIBLE|style,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	return c, nil
}

func NewComposite(parent Container) (*Composite, error) {
	return newCompositeWithStyle(parent, 0)
}
