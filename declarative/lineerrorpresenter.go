// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type LineErrorPresenter struct {
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

	// LineErrorPresenter

	AssignTo *walk.ErrorPresenter
}

func (lep LineErrorPresenter) Create(builder *Builder) error {
	w, err := walk.NewLineErrorPresenter(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(lep, w, func() error {
		if lep.AssignTo != nil {
			*lep.AssignTo = w
		}

		return nil
	})
}
