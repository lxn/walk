// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TreeView struct {
	Widget        **walk.TreeView
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
}

func (tv TreeView) Create(parent walk.Container) error {
	w, err := walk.NewTreeView(parent)
	if err != nil {
		return err
	}

	return InitWidget(tv, w, func() error {
		if tv.Widget != nil {
			*tv.Widget = w
		}

		return nil
	})
}

func (tv TreeView) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return tv.Name, tv.StretchFactor, tv.Row, tv.RowSpan, tv.Column, tv.ColumnSpan
}

func (tv TreeView) Font_() *Font {
	return &tv.Font
}
