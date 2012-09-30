// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TextEdit struct {
	AssignTo      **walk.TextEdit
	Name          string
	MinSize       Size
	MaxSize       Size
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	ContextMenu   Menu
	Font          Font
	Text          string
	ReadOnly      bool
}

func (te TextEdit) Create(parent walk.Container) error {
	w, err := walk.NewTextEdit(parent)
	if err != nil {
		return err
	}

	return InitWidget(te, w, func() error {
		if err := w.SetText(te.Text); err != nil {
			return err
		}

		w.SetReadOnly(te.ReadOnly)

		if te.AssignTo != nil {
			*te.AssignTo = w
		}

		return nil
	})
}

func (te TextEdit) CommonInfo() (name string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return te.Name, te.MinSize, te.MaxSize, te.StretchFactor, te.Row, te.RowSpan, te.Column, te.ColumnSpan, &te.ContextMenu
}

func (te TextEdit) Font_() *Font {
	return &te.Font
}
