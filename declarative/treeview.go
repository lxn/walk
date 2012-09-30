// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TreeView struct {
	AssignTo      **walk.TreeView
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
}

func (tv TreeView) Create(parent walk.Container) error {
	w, err := walk.NewTreeView(parent)
	if err != nil {
		return err
	}

	return InitWidget(tv, w, func() error {
		if tv.AssignTo != nil {
			*tv.AssignTo = w
		}

		return nil
	})
}

func (tv TreeView) WidgetInfo() (name string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return tv.Name, tv.MinSize, tv.MaxSize, tv.StretchFactor, tv.Row, tv.RowSpan, tv.Column, tv.ColumnSpan, &tv.ContextMenu
}

func (tv TreeView) Font_() *Font {
	return &tv.Font
}
