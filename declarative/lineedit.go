// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type LineEdit struct {
	AssignTo           **walk.LineEdit
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
	CueBanner          string
	MaxLength          int
	PasswordMode       bool
	OnEditingFinished  walk.EventHandler
	OnReturnPressed    walk.EventHandler
	OnTextChanged      walk.EventHandler
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

func (le LineEdit) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return le.Name, le.Disabled, le.Hidden, &le.Font, le.MinSize, le.MaxSize, le.StretchFactor, le.Row, le.RowSpan, le.Column, le.ColumnSpan, le.ContextMenuActions
}
