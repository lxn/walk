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
	Enabled           Property
	Visible           Property
	Font              Font
	ToolTipText       Property
	MinSize           Size
	MaxSize           Size
	StretchFactor     int
	Row               int
	RowSpan           int
	Column            int
	ColumnSpan        int
	ContextMenuItems  []MenuItem
	OnKeyDown         walk.KeyEventHandler
	OnKeyUp           walk.KeyEventHandler
	OnMouseDown       walk.MouseEventHandler
	OnMouseMove       walk.MouseEventHandler
	OnMouseUp         walk.MouseEventHandler
	OnSizeChanged     walk.EventHandler
	Text              Property
	ReadOnly          Property
	CueBanner         string
	MaxLength         int
	PasswordMode      bool
	OnEditingFinished walk.EventHandler
	OnTextChanged     walk.EventHandler
}

func (le LineEdit) Create(builder *Builder) error {
	w, err := walk.NewLineEdit(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(le, w, func() error {
		if err := w.SetCueBanner(le.CueBanner); err != nil {
			return err
		}
		w.SetMaxLength(le.MaxLength)
		w.SetPasswordMode(le.PasswordMode)

		if le.OnEditingFinished != nil {
			w.EditingFinished().Attach(le.OnEditingFinished)
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

func (w LineEdit) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuItems []MenuItem, OnKeyDown walk.KeyEventHandler, OnKeyUp walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, false, false, &w.Font, "", w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuItems, w.OnKeyDown, w.OnKeyUp, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}
