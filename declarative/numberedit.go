// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type NumberEdit struct {
	AssignTo           **walk.NumberEdit
	Name               string
	Disabled           bool
	Hidden             bool
	Font               Font
	ToolTipText        string
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	OnKeyDown          walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	BindTo             string
	Validator          Validator
	Decimals           int
	Increment          float64
	MinValue           float64
	MaxValue           float64
	Value              float64
	OnValueChanged     walk.EventHandler
}

func (ne NumberEdit) Create(builder *Builder) error {
	w, err := walk.NewNumberEdit(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(ne, w, func() error {
		if err := w.SetBindingMember(ne.BindTo); err != nil {
			return err
		}

		if err := w.SetDecimals(ne.Decimals); err != nil {
			return err
		}

		inc := ne.Increment
		if inc <= 0 {
			inc = 1
		}

		if err := w.SetIncrement(inc); err != nil {
			return err
		}

		if ne.MinValue != 0 || ne.MaxValue != 0 {
			if err := w.SetRange(ne.MinValue, ne.MaxValue); err != nil {
				return err
			}
		}

		if err := w.SetValue(ne.Value); err != nil {
			return err
		}

		if ne.OnValueChanged != nil {
			w.ValueChanged().Attach(ne.OnValueChanged)
		}

		if ne.AssignTo != nil {
			*ne.AssignTo = w
		}

		return nil
	})
}

func (w NumberEdit) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action, OnKeyDown walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, w.Disabled, w.Hidden, &w.Font, w.ToolTipText, w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuActions, w.OnKeyDown, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}

func (ne NumberEdit) ValidatableInfo() Validator {
	return ne.Validator
}
