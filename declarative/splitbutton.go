// Copyright 2016 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type SplitButton struct {
	// Window

	Background         Brush
	ContextMenuItems   []MenuItem
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

	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	GraphicsEffects    []walk.WidgetGraphicsEffect
	Row                int
	RowSpan            int
	StretchFactor      int

	// Button

	Image     Property
	Text      Property
	OnClicked walk.EventHandler

	// SplitButton

	AssignTo       **walk.SplitButton
	ImageAboveText bool
	MenuItems      []MenuItem
}

func (sb SplitButton) Create(builder *Builder) error {
	w, err := walk.NewSplitButton(builder.Parent())
	if err != nil {
		return err
	}

	if sb.AssignTo != nil {
		*sb.AssignTo = w
	}

	builder.deferBuildMenuActions(w.Menu(), sb.MenuItems)

	return builder.InitWidget(sb, w, func() error {
		if err := w.SetImageAboveText(sb.ImageAboveText); err != nil {
			return err
		}

		if sb.OnClicked != nil {
			w.Clicked().Attach(sb.OnClicked)
		}

		return nil
	})
}
