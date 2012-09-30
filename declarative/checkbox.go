// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type CheckBox struct {
	AssignTo      **walk.CheckBox
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
	OnClicked     walk.EventHandler
}

func (cb CheckBox) Create(parent walk.Container) error {
	w, err := walk.NewCheckBox(parent)
	if err != nil {
		return err
	}

	return InitWidget(cb, w, func() error {
		if err := w.SetText(cb.Text); err != nil {
			return err
		}

		if cb.OnClicked != nil {
			w.Clicked().Attach(cb.OnClicked)
		}

		if cb.AssignTo != nil {
			*cb.AssignTo = w
		}

		return nil
	})
}

func (cb CheckBox) CommonInfo() (name string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return cb.Name, cb.MinSize, cb.MaxSize, cb.StretchFactor, cb.Row, cb.RowSpan, cb.Column, cb.ColumnSpan, &cb.ContextMenu
}

func (cb CheckBox) Font_() *Font {
	return &cb.Font
}
