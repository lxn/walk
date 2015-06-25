// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type CustomWidget struct {
	AssignTo            **walk.CustomWidget
	Name                string
	Enabled             Property
	Visible             Property
	Font                Font
	ToolTipText         Property
	MinSize             Size
	MaxSize             Size
	StretchFactor       int
	Row                 int
	RowSpan             int
	Column              int
	ColumnSpan          int
	AlwaysConsumeSpace  bool
	ContextMenuItems    []MenuItem
	OnKeyDown           walk.KeyEventHandler
	OnKeyPress          walk.KeyEventHandler
	OnKeyUp             walk.KeyEventHandler
	OnMouseDown         walk.MouseEventHandler
	OnMouseMove         walk.MouseEventHandler
	OnMouseUp           walk.MouseEventHandler
	OnSizeChanged       walk.EventHandler
	Style               uint32
	Paint               walk.PaintFunc
	ClearsBackground    bool
	InvalidatesOnResize bool
}

func (cw CustomWidget) Create(builder *Builder) error {
	w, err := walk.NewCustomWidget(builder.Parent(), uint(cw.Style), cw.Paint)
	if err != nil {
		return err
	}

	return builder.InitWidget(cw, w, func() error {
		w.SetClearsBackground(cw.ClearsBackground)
		w.SetInvalidatesOnResize(cw.InvalidatesOnResize)

		if cw.AssignTo != nil {
			*cw.AssignTo = w
		}

		return nil
	})
}

func (w CustomWidget) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, alwaysConsumeSpace bool, contextMenuItems []MenuItem, OnKeyDown walk.KeyEventHandler, OnKeyPress walk.KeyEventHandler, OnKeyUp walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, false, false, &w.Font, "", w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.AlwaysConsumeSpace, w.ContextMenuItems, w.OnKeyDown, w.OnKeyPress, w.OnKeyUp, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}
