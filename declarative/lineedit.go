// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type LineEdit struct {
	AssignTo          **walk.LineEdit
	Name              string
	StretchFactor     int
	Row               int
	RowSpan           int
	Column            int
	ColumnSpan        int
	ContextMenu       Menu
	Font              Font
	Text              string
	ReadOnly          bool
	CueBanner         string
	MaxLength         int
	PasswordMode      bool
	OnEditingFinished walk.EventHandler
	OnReturnPressed   walk.EventHandler
	OnTextChanged     walk.EventHandler
}

func (le LineEdit) Create(parent walk.Container) error {
	w, err := walk.NewLineEdit(parent)
	if err != nil {
		return err
	}

	return InitWidget(le, w, func() error {
		if err := w.SetText(le.Text); err != nil {
			return err
		}

		w.SetReadOnly(le.ReadOnly)

		if err := w.SetCueBanner(le.CueBanner); err != nil {
			return err
		}
		w.SetMaxLength(le.MaxLength)
		w.SetPasswordMode(le.PasswordMode)

		if le.OnEditingFinished != nil {
			w.EditingFinished().Attach(le.OnEditingFinished)
		}
		if le.OnReturnPressed != nil {
			w.ReturnPressed().Attach(le.OnReturnPressed)
		}
		if le.OnTextChanged != nil {
			w.TextChanged().Attach(le.OnTextChanged)
		}

		if le.AssignTo != nil {
			*le.AssignTo = w
		}

		return nil
	})
}

func (le LineEdit) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return le.Name, le.StretchFactor, le.Row, le.RowSpan, le.Column, le.ColumnSpan, &le.ContextMenu
}

func (le LineEdit) Font_() *Font {
	return &le.Font
}
