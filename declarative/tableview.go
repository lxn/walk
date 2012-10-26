// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TableView struct {
	AssignTo                   **walk.TableView
	Name                       string
	Disabled                   bool
	Hidden                     bool
	Font                       Font
	ToolTipText                string
	MinSize                    Size
	MaxSize                    Size
	StretchFactor              int
	Row                        int
	RowSpan                    int
	Column                     int
	ColumnSpan                 int
	ContextMenuActions         []*walk.Action
	OnKeyDown                  walk.KeyEventHandler
	OnMouseDown                walk.MouseEventHandler
	OnMouseMove                walk.MouseEventHandler
	OnMouseUp                  walk.MouseEventHandler
	OnSizeChanged              walk.EventHandler
	Model                      walk.TableModel
	AlternatingRowBGColor      walk.Color
	CheckBoxes                 bool
	ItemStateChangedEventDelay int
	LastColumnStretched        bool
	ReorderColumnsEnabled      bool
	SingleItemSelection        bool
	OnCurrentIndexChanged      walk.EventHandler
	OnSelectedIndexesChanged   walk.EventHandler
	OnItemActivated            walk.EventHandler
}

func (tv TableView) Create(parent walk.Container) error {
	w, err := walk.NewTableView(parent)
	if err != nil {
		return err
	}

	return InitWidget(tv, w, func() error {
		if err := w.SetModel(tv.Model); err != nil {
			return err
		}

		if tv.AlternatingRowBGColor != 0 {
			w.SetAlternatingRowBGColor(tv.AlternatingRowBGColor)
		}
		w.SetCheckBoxes(tv.CheckBoxes)
		w.SetItemStateChangedEventDelay(tv.ItemStateChangedEventDelay)
		if err := w.SetLastColumnStretched(tv.LastColumnStretched); err != nil {
			return err
		}
		w.SetReorderColumnsEnabled(tv.ReorderColumnsEnabled)
		if err := w.SetSingleItemSelection(tv.SingleItemSelection); err != nil {
			return err
		}

		if tv.OnCurrentIndexChanged != nil {
			w.CurrentIndexChanged().Attach(tv.OnCurrentIndexChanged)
		}
		if tv.OnSelectedIndexesChanged != nil {
			w.SelectedIndexesChanged().Attach(tv.OnSelectedIndexesChanged)
		}
		if tv.OnItemActivated != nil {
			w.ItemActivated().Attach(tv.OnItemActivated)
		}

		if tv.AssignTo != nil {
			*tv.AssignTo = w
		}

		return nil
	})
}

func (w TableView) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action, OnKeyDown walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, w.Disabled, w.Hidden, &w.Font, w.ToolTipText, w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuActions, w.OnKeyDown, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}
