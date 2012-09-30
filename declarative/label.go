// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Label struct {
	AssignTo      **walk.Label
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
	Text          string
}

func (l Label) Create(parent walk.Container) error {
	w, err := walk.NewLabel(parent)
	if err != nil {
		return err
	}

	return InitWidget(l, w, func() error {
		if err := w.SetText(l.Text); err != nil {
			return err
		}

		if l.AssignTo != nil {
			*l.AssignTo = w
		}

		return nil
	})
}

func (l Label) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return l.Name, l.StretchFactor, l.Row, l.RowSpan, l.Column, l.ColumnSpan
}

func (l Label) Font_() *Font {
	return &l.Font
}
