// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type Composite struct {
	// Window

	Background       Brush
	ContextMenuItems []MenuItem
	Enabled          Property
	Font             Font
	MaxSize          Size
	MinSize          Size
	Name             string
	OnKeyDown        walk.KeyEventHandler
	OnKeyPress       walk.KeyEventHandler
	OnKeyUp          walk.KeyEventHandler
	OnMouseDown      walk.MouseEventHandler
	OnMouseMove      walk.MouseEventHandler
	OnMouseUp        walk.MouseEventHandler
	OnSizeChanged    walk.EventHandler
	Persistent       bool
	ToolTipText      Property
	Visible          Property

	// Widget

	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	Row                int
	RowSpan            int
	StretchFactor      int

	// Container

	Children   []Widget
	DataBinder DataBinder
	Layout     Layout

	// Composite

	AssignTo    **walk.Composite
	Expressions func() map[string]walk.Expression
	Functions   map[string]func(args ...interface{}) (interface{}, error)
}

func (c Composite) Create(builder *Builder) error {
	w, err := walk.NewComposite(builder.Parent())
	if err != nil {
		return err
	}

	w.SetSuspended(true)
	builder.Defer(func() error {
		w.SetSuspended(false)
		return nil
	})

	return builder.InitWidget(c, w, func() error {
		if c.AssignTo != nil {
			*c.AssignTo = w
		}

		if c.Expressions != nil {
			for name, expr := range c.Expressions() {
				builder.expressions[name] = expr
			}
		}
		if c.Functions != nil {
			for name, fn := range c.Functions {
				builder.functions[name] = fn
			}
		}

		return nil
	})
}
