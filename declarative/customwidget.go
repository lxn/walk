// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type CustomWidget struct {
	Widget              **walk.CustomWidget
	Name                string
	StretchFactor       int
	Row                 int
	RowSpan             int
	Column              int
	ColumnSpan          int
	Style               uint32
	Paint               walk.PaintFunc
	ClearsBackground    bool
	InvalidatesOnResize bool
}

func (cw CustomWidget) Create(parent walk.Container) error {
	w, err := walk.NewCustomWidget(parent, uint(cw.Style), cw.Paint)
	if err != nil {
		return err
	}

	return InitWidget(cw, w, func() error {
		w.SetClearsBackground(cw.ClearsBackground)
		w.SetInvalidatesOnResize(cw.InvalidatesOnResize)

		if cw.Widget != nil {
			*cw.Widget = w
		}

		return nil
	})
}

func (cw CustomWidget) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return cw.Name, cw.StretchFactor, cw.Row, cw.RowSpan, cw.Column, cw.ColumnSpan
}
