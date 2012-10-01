// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ListBox struct {
	AssignTo              **walk.ListBox
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
	Format                string
	Precision             int
	Model                 walk.ListModel
	OnCurrentIndexChanged walk.EventHandler
	OnItemActivated       walk.EventHandler
}

func (lb ListBox) Create(parent walk.Container) error {
	w, err := walk.NewListBox(parent)
	if err != nil {
		return err
	}

	return InitWidget(lb, w, func() error {
		w.SetFormat(lb.Format)
		w.SetPrecision(lb.Precision)

		if err := w.SetModel(lb.Model); err != nil {
			return err
		}

		if lb.OnCurrentIndexChanged != nil {
			w.CurrentIndexChanged().Attach(lb.OnCurrentIndexChanged)
		}
		if lb.OnItemActivated != nil {
			w.DblClicked().Attach(lb.OnItemActivated)
		}

		if lb.AssignTo != nil {
			*lb.AssignTo = w
		}

		return nil
	})
}

func (lb ListBox) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return lb.Name, lb.Disabled, lb.Hidden, &lb.Font, lb.MinSize, lb.MaxSize, lb.StretchFactor, lb.Row, lb.RowSpan, lb.Column, lb.ColumnSpan, lb.ContextMenuActions
}
