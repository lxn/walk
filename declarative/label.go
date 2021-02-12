// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type EllipsisMode int

const (
	EllipsisNone = EllipsisMode(walk.EllipsisNone)
	EllipsisEnd  = EllipsisMode(walk.EllipsisEnd)
	EllipsisPath = EllipsisMode(walk.EllipsisPath)
)

type Label struct {
	// Window

	Accessibility      Accessibility
	Background         Brush
	ContextMenuItems   []MenuItem
	DoubleBuffering    bool
	Enabled            Property
	Font               Font
	MaxSize            Size
	MinSize            Size
	Name               string
	OnBoundsChanged    walk.EventHandler
	OnKeyDown          walk.KeyEventHandler
	OnKeyPress         walk.KeyEventHandler
	OnKeyUp            walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	Persistent         bool
	RightToLeftReading bool
	ToolTipText        Property
	Visible            Property

	// Widget

	Alignment          Alignment2D
	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	GraphicsEffects    []walk.WidgetGraphicsEffect
	Row                int
	RowSpan            int
	StretchFactor      int

	// Label

	AssignTo      **walk.Label
	EllipsisMode  EllipsisMode
	NoPrefix      bool
	Text          Property
	TextAlignment Alignment1D
	TextColor     walk.Color
}

func (l Label) Create(builder *Builder) error {
	var style uint32
	if l.NoPrefix {
		style |= win.SS_NOPREFIX
	}

	w, err := walk.NewLabelWithStyle(builder.Parent(), style)
	if err != nil {
		return err
	}

	if l.AssignTo != nil {
		*l.AssignTo = w
	}

	return builder.InitWidget(l, w, func() error {
		if err := w.SetEllipsisMode(walk.EllipsisMode(l.EllipsisMode)); err != nil {
			return err
		}

		if err := w.SetTextAlignment(walk.Alignment1D(l.TextAlignment)); err != nil {
			return err
		}

		w.SetTextColor(l.TextColor)

		return nil
	})
}
