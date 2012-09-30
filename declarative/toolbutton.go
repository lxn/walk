// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ToolButton struct {
	AssignTo      **walk.ToolButton
	Name          string
	MinSize       Size
	MaxSize       Size
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	ContextMenu   Menu
	Font          Font
	Text          string
	OnClicked     walk.EventHandler
}

func (tb ToolButton) Create(parent walk.Container) error {
	w, err := walk.NewToolButton(parent)
	if err != nil {
		return err
	}

	return InitWidget(tb, w, func() error {
		if err := w.SetText(tb.Text); err != nil {
			return err
		}

		if tb.OnClicked != nil {
			w.Clicked().Attach(tb.OnClicked)
		}

		if tb.AssignTo != nil {
			*tb.AssignTo = w
		}

		return nil
	})
}

func (tb ToolButton) CommonInfo() (name string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return tb.Name, tb.MinSize, tb.MaxSize, tb.StretchFactor, tb.Row, tb.RowSpan, tb.Column, tb.ColumnSpan, &tb.ContextMenu
}

func (tb ToolButton) Font_() *Font {
	return &tb.Font
}
