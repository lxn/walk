// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type GroupBox struct {
	AssignTo      **walk.GroupBox
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
	Title         string
	Layout        Layout
	Children      []Widget
}

func (gb GroupBox) Create(parent walk.Container) error {
	w, err := walk.NewGroupBox(parent)
	if err != nil {
		return err
	}

	return InitWidget(gb, w, func() error {
		if err := w.SetTitle(gb.Title); err != nil {
			return err
		}

		if gb.AssignTo != nil {
			*gb.AssignTo = w
		}

		return nil
	})
}

func (gb GroupBox) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return gb.Name, gb.StretchFactor, gb.Row, gb.RowSpan, gb.Column, gb.ColumnSpan
}

func (gb GroupBox) Font_() *Font {
	return &gb.Font
}

func (gb GroupBox) ContainerInfo() (Layout, []Widget) {
	return gb.Layout, gb.Children
}
