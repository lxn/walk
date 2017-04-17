// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type TextEdit struct {
	// Window

	ContextMenuItems []MenuItem
	Enabled          Property
	Font             Font
	MaxSize          Size
	MinSize          Size
	Name             string
	OnKeyDown        walk.KeyEventHandler
	OnKeyPress       walk.KeyEventHandler
	OnKeyUp          walk.KeyEventHandler
	OnMouseDown      walk.MouseEventHandler
	OnMouseMove      walk.MouseEventHandler
	OnMouseUp        walk.MouseEventHandler
	OnSizeChanged    walk.EventHandler
	Persistent       bool
	ToolTipText      Property
	Visible          Property

	// Widget

	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	Row                int
	RowSpan            int
	StretchFactor      int

	// TextEdit

	AssignTo  **walk.TextEdit
	MaxLength int
	ReadOnly  Property
	Text      Property
}

func (te TextEdit) Create(builder *Builder) error {
	w, err := walk.NewTextEdit(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(te, w, func() error {
		if te.AssignTo != nil {
			*te.AssignTo = w
		}

		if te.MaxLength > 0 {
			w.SetMaxLength(te.MaxLength)
		}

		return nil
	})
}
