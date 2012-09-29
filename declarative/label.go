// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Label struct {
	Widget        **walk.Label
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

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	if err := initWidget(l, w); err != nil {
		return err
	}

	w.SetName(l.Name)

	if err := w.SetText(l.Text); err != nil {
		return err
	}

	if l.Widget != nil {
		*l.Widget = w
	}

	succeeded = true

	return nil
}

func (l Label) LayoutParams() (stretchFactor, row, rowSpan, column, columnSpan int) {
	return l.StretchFactor, l.Row, l.RowSpan, l.Column, l.ColumnSpan
}
