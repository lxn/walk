// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type ImageView struct {
	// Window

	Background         Brush
	ContextMenuItems   []MenuItem
	Enabled            Property
	Font               Font
	MaxSize            Size
	MinSize            Size
	Name               string
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
	Row                int
	RowSpan            int
	StretchFactor      int

	// ImageView

	AssignTo **walk.ImageView
	Image    Property
	Margin   Property
}

func (iv ImageView) Create(builder *Builder) error {
	w, err := walk.NewImageView(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(iv, w, func() error {
		if iv.AssignTo != nil {
			*iv.AssignTo = w
		}

		return nil
	})
}
