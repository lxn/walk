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
	OnKeyDown             walk.KeyEventHandler
	OnMouseDown           walk.MouseEventHandler
	OnMouseMove           walk.MouseEventHandler
	OnMouseUp             walk.MouseEventHandler
	OnSizeChanged         walk.EventHandler
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

func (w ListBox) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action, OnKeyDown walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, w.Disabled, w.Hidden, &w.Font, w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuActions, w.OnKeyDown, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}
