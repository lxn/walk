// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"errors"
)

import (
	"github.com/lxn/walk"
)

type ComboBox struct {
	AssignTo              **walk.ComboBox
	Name                  string
	Disabled              bool
	Hidden                bool
	Font                  Font
	MinSize               Size
	MaxSize               Size
	StretchFactor         int
	Row                   int
	RowSpan               int
	Column                int
	ColumnSpan            int
	ContextMenuActions    []*walk.Action
	BindTo                string
	Optional              bool
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
		if _, ok := cb.Model.(walk.BindingValueProvider); !ok && cb.BindTo != "" {
			return errors.New("declarative.ComboBox: Data binding is only supported using a model that implements BindingValueProvider.")
		}

		if err := w.SetBindingMember(cb.BindTo); err != nil {
			return err
		}

		if !cb.Optional {
			w.SetValidator(walk.SelectionRequiredValidator())
		}

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

func (cb ComboBox) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return cb.Name, cb.Disabled, cb.Hidden, &cb.Font, cb.MinSize, cb.MaxSize, cb.StretchFactor, cb.Row, cb.RowSpan, cb.Column, cb.ColumnSpan, cb.ContextMenuActions
}
