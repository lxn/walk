// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type CheckBox struct {
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

	// Button

	Checked          Property
	OnClicked        walk.EventHandler
	OnCheckedChanged walk.EventHandler
	Text             Property

	// CheckBox

	AssignTo            **walk.CheckBox
	CheckState          Property
	OnCheckStateChanged walk.EventHandler
	Tristate            bool
}

func (cb CheckBox) Create(builder *Builder) error {
	w, err := walk.NewCheckBox(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(cb, w, func() error {
		w.SetPersistent(cb.Persistent)

		if err := w.SetTristate(cb.Tristate); err != nil {
			return err
		}

		if cb.OnClicked != nil {
			w.Clicked().Attach(cb.OnClicked)
		}

		if cb.OnCheckedChanged != nil {
			w.CheckedChanged().Attach(cb.OnCheckedChanged)
		}

		if cb.OnCheckStateChanged != nil {
			w.CheckStateChanged().Attach(cb.OnCheckStateChanged)
		}

		if cb.AssignTo != nil {
			*cb.AssignTo = w
		}

		return nil
	})
}
