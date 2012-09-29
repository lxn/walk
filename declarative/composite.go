// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Composite struct {
	Widget        **walk.Composite
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Layout        Layout
	Children      []Widget
}

func (c Composite) Create(parent walk.Container) error {
	w, err := walk.NewComposite(parent)
	if err != nil {
		return err
	}

	return InitWidget(c, w, func() error {
		if c.Widget != nil {
			*c.Widget = w
		}

		return nil
	})
}

func (c Composite) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return c.Name, c.StretchFactor, c.Row, c.RowSpan, c.Column, c.ColumnSpan
}

func (c Composite) Layout_() Layout {
	return c.Layout
}

func (c Composite) Children_() []Widget {
	return c.Children
}
