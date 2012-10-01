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
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	Font               Font
	Decimals           int
	Increment          float64
	MinValue           float64
	MaxValue           float64
	Value              float64
	OnValueChanged     walk.EventHandler
}

func (ne NumberEdit) Create(parent walk.Container) error {
	w, err := walk.NewNumberEdit(parent)
	if err != nil {
		return err
	}

	return InitWidget(ne, w, func() error {
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

func (ne NumberEdit) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return ne.Name, ne.Disabled, ne.Hidden, ne.MinSize, ne.MaxSize, ne.StretchFactor, ne.Row, ne.RowSpan, ne.Column, ne.ColumnSpan, ne.ContextMenuActions
}

func (ne NumberEdit) Font_() *Font {
	return &ne.Font
}
