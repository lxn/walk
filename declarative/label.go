// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Label struct {
	AssignTo           **walk.Label
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
	BindTo             string
	Text               string
}

func (l Label) Create(parent walk.Container) error {
	w, err := walk.NewLabel(parent)
	if err != nil {
		return err
	}

	return InitWidget(l, w, func() error {
		if err := w.SetBindingMember(l.BindTo); err != nil {
			return err
		}

		if err := w.SetText(l.Text); err != nil {
			return err
		}

		if l.AssignTo != nil {
			*l.AssignTo = w
		}

		return nil
	})
}

func (l Label) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return l.Name, l.Disabled, l.Hidden, &l.Font, l.MinSize, l.MaxSize, l.StretchFactor, l.Row, l.RowSpan, l.Column, l.ColumnSpan, l.ContextMenuActions
}
