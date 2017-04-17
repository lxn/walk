// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type GradientComposite struct {
	// Window

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
	Layout     Layout
	DataBinder DataBinder

	// GradientComposite

	AssignTo    **walk.GradientComposite
	Color1      Property
	Color2      Property
	Expressions func() map[string]walk.Expression
	Functions   map[string]func(args ...interface{}) (interface{}, error)
	Vertical    Property
}

func (gc GradientComposite) Create(builder *Builder) error {
	w, err := walk.NewGradientComposite(builder.Parent())
	if err != nil {
		return err
	}

	w.SetSuspended(true)
	builder.Defer(func() error {
		w.SetSuspended(false)
		return nil
	})

	return builder.InitWidget(gc, w, func() error {
		if gc.AssignTo != nil {
			*gc.AssignTo = w
		}

		if gc.Expressions != nil {
			for name, expr := range gc.Expressions() {
				builder.expressions[name] = expr
			}
		}
		if gc.Functions != nil {
			for name, fn := range gc.Functions {
				builder.functions[name] = fn
			}
		}

		return nil
	})
}
