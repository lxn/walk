// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ComboBox struct {
	AssignTo              **walk.ComboBox
	Name                  string
	Disabled              bool
	Hidden                bool
	MinSize               Size
	MaxSize               Size
	StretchFactor         int
	Row                   int
	RowSpan               int
	Column                int
	ColumnSpan            int
	ContextMenuActions    []*walk.Action
	Font                  Font
	Format                string
	Precision             int
	Model                 walk.ListModel
	OnCurrentIndexChanged walk.EventHandler
}

func (cb ComboBox) Create(parent walk.Container) error {
	w, err := walk.NewComboBox(parent)
	if err != nil {
		return err
	}

	return InitWidget(cb, w, func() error {
		w.SetFormat(cb.Format)
		w.SetPrecision(cb.Precision)

		if err := w.SetModel(cb.Model); err != nil {
			return err
		}

		if cb.OnCurrentIndexChanged != nil {
			w.CurrentIndexChanged().Attach(cb.OnCurrentIndexChanged)
		}

		if cb.AssignTo != nil {
			*cb.AssignTo = w
		}

		return nil
	})
}

func (cb ComboBox) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return cb.Name, cb.Disabled, cb.Hidden, cb.MinSize, cb.MaxSize, cb.StretchFactor, cb.Row, cb.RowSpan, cb.Column, cb.ColumnSpan, cb.ContextMenuActions
}

func (cb ComboBox) Font_() *Font {
	return &cb.Font
}
