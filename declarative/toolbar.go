// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ToolBar struct {
	AssignTo      **walk.ToolBar
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	ContextMenu   Menu
	Font          Font
	Orientation   walk.Orientation
	Actions       []*walk.Action
}

func (tb ToolBar) Create(parent walk.Container) (err error) {
	var w *walk.ToolBar
	if tb.Orientation == walk.Vertical {
		w, err = walk.NewVerticalToolBar(parent)
	} else {
		w, err = walk.NewToolBar(parent)
	}
	if err != nil {
		return
	}

	return InitWidget(tb, w, func() error {
		imageList, err := walk.NewImageList(walk.Size{16, 16}, 0)
		if err != nil {
			return err
		}
		w.SetImageList(imageList)

		if err := tb.initActions(w); err != nil {
			return err
		}

		if tb.AssignTo != nil {
			*tb.AssignTo = w
		}

		return nil
	})
}

func (tb ToolBar) initActions(w *walk.ToolBar) error {
	actions := w.Actions()

	for _, a := range tb.Actions {
		if err := actions.Add(a); err != nil {
			return err
		}
	}

	return nil
}

func (tb ToolBar) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return tb.Name, tb.StretchFactor, tb.Row, tb.RowSpan, tb.Column, tb.ColumnSpan, &tb.ContextMenu
}

func (tb ToolBar) Font_() *Font {
	return &tb.Font
}
