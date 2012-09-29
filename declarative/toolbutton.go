// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ToolButton struct {
	Widget        **walk.ToolButton
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
	Text          string
	OnClicked     walk.EventHandler
}

func (tb ToolButton) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewToolButton(parent)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	if err := initWidget(tb, w); err != nil {
		return nil, err
	}

	w.SetName(tb.Name)

	if err := w.SetText(tb.Text); err != nil {
		return nil, err
	}

	if tb.OnClicked != nil {
		w.Clicked().Attach(tb.OnClicked)
	}

	if tb.Widget != nil {
		*tb.Widget = w
	}

	succeeded = true

	return w, nil
}

func (tb ToolButton) LayoutParams() (stretchFactor, row, rowSpan, column, columnSpan int) {
	return tb.StretchFactor, tb.Row, tb.RowSpan, tb.Column, tb.ColumnSpan
}
