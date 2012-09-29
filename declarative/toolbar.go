// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ToolBar struct {
	Widget        **walk.ToolBar
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
	Orientation   walk.Orientation
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
		if tb.Widget != nil {
			*tb.Widget = w
		}

		return nil
	})
}

func (tb ToolBar) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return tb.Name, tb.StretchFactor, tb.Row, tb.RowSpan, tb.Column, tb.ColumnSpan
}

func (tb ToolBar) Font_() *Font {
	return &tb.Font
}
