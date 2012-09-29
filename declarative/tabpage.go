// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TabPage struct {
	Widget        **walk.TabPage
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Title         string
	Layout        Layout
	Children      []Widget
}

func (tp TabPage) Create(parent walk.Container) error {
	w, err := walk.NewTabPage()
	if err != nil {
		return err
	}

	return InitWidget(tp, w, func() error {
		if err := w.SetTitle(tp.Title); err != nil {
			return err
		}

		if tp.Widget != nil {
			*tp.Widget = w
		}

		return nil
	})
}

func (tp TabPage) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return tp.Name, tp.StretchFactor, tp.Row, tp.RowSpan, tp.Column, tp.ColumnSpan
}

func (tp TabPage) ContainerInfo() (Layout, []Widget) {
	return tp.Layout, tp.Children
}
