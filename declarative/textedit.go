// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TextEdit struct {
	Widget     **walk.TextEdit
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Font       Font
	Text       string
	ReadOnly   bool
}

func (te TextEdit) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewTextEdit(parent)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	if err := initWidget(te, w); err != nil {
		return nil, err
	}

	w.SetName(te.Name)

	if err := w.SetText(te.Text); err != nil {
		return nil, err
	}

	w.SetReadOnly(te.ReadOnly)

	if te.Widget != nil {
		*te.Widget = w
	}

	succeeded = true

	return w, nil
}

func (te TextEdit) LayoutParams() (hStretch, vStretch, row, rowSpan, col, colSpan int) {
	return te.HStretch, te.VStretch, te.Row, te.RowSpan, te.Column, te.ColumnSpan
}
