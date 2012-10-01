// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TextEdit struct {
	AssignTo           **walk.TextEdit
	Name               string
	Disabled           bool
	Hidden             bool
	Font               Font
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	Text               string
	ReadOnly           bool
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

func (te TextEdit) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return te.Name, te.Disabled, te.Hidden, &te.Font, te.MinSize, te.MaxSize, te.StretchFactor, te.Row, te.RowSpan, te.Column, te.ColumnSpan, te.ContextMenuActions
}
