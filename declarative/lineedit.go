// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type LineEdit struct {
	Widget        **walk.LineEdit
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
	Text          string
	ReadOnly      bool
	MaxLength     int
}

func (le LineEdit) Create(parent walk.Container) error {
	w, err := walk.NewLineEdit(parent)
	if err != nil {
		return err
	}

	return InitWidget(le, w, func() error {
		w.SetName(le.Name)

		if err := w.SetText(le.Text); err != nil {
			return err
		}

		w.SetReadOnly(le.ReadOnly)
		w.SetMaxLength(le.MaxLength)

		if le.Widget != nil {
			*le.Widget = w
		}

		return nil
	})
}

func (le LineEdit) LayoutParams() (stretchFactor, row, rowSpan, column, columnSpan int) {
	return le.StretchFactor, le.Row, le.RowSpan, le.Column, le.ColumnSpan
}
