// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type GradientComposite struct {
	AssignTo           **walk.GradientComposite
	Name               string
	Enabled            Property
	Visible            Property
	Font               Font
	ToolTipText        Property
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	AlwaysConsumeSpace bool
	ContextMenuItems   []MenuItem
	OnKeyDown          walk.KeyEventHandler
	OnKeyPress         walk.KeyEventHandler
	OnKeyUp            walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	DataBinder         DataBinder
	Layout             Layout
	Children           []Widget
	Expressions        func() map[string]walk.Expression
	Functions          map[string]func(args ...interface{}) (interface{}, error)
	Vertical           Property
	Color1             Property
	Color2             Property
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

func (w GradientComposite) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, alwaysConsumeSpace bool, contextMenuItems []MenuItem, OnKeyDown walk.KeyEventHandler, OnKeyPress walk.KeyEventHandler, OnKeyUp walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, false, false, &w.Font, "", w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.AlwaysConsumeSpace, w.ContextMenuItems, w.OnKeyDown, w.OnKeyPress, w.OnKeyUp, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}

func (gc GradientComposite) ContainerInfo() (DataBinder, Layout, []Widget) {
	return gc.DataBinder, gc.Layout, gc.Children
}
