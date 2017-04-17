// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type CaseMode uint32

const (
	CaseModeMixed CaseMode = CaseMode(walk.CaseModeMixed)
	CaseModeUpper CaseMode = CaseMode(walk.CaseModeUpper)
	CaseModeLower CaseMode = CaseMode(walk.CaseModeLower)
)

type LineEdit struct {
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

	// LineEdit

	AssignTo          **walk.LineEdit
	CaseMode          CaseMode
	CueBanner         string
	MaxLength         int
	OnEditingFinished walk.EventHandler
	OnTextChanged     walk.EventHandler
	PasswordMode      bool
	ReadOnly          Property
	Text              Property
}

func (le LineEdit) Create(builder *Builder) error {
	w, err := walk.NewLineEdit(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(le, w, func() error {
		if le.CueBanner != "" {
			if err := w.SetCueBanner(le.CueBanner); err != nil {
				return err
			}
		}
		w.SetMaxLength(le.MaxLength)
		w.SetPasswordMode(le.PasswordMode)

		if err := w.SetCaseMode(walk.CaseMode(le.CaseMode)); err != nil {
			return err
		}

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
