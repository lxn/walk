// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type TabPage struct {
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

	// Container

	Children   []Widget
	DataBinder DataBinder
	Layout     Layout

	// TabPage

	AssignTo **walk.TabPage
	Content  Widget
	Image    Property
	Title    Property
}

func (tp TabPage) Create(builder *Builder) error {
	w, err := walk.NewTabPage()
	if err != nil {
		return err
	}

	return builder.InitWidget(tp, w, func() error {
		if tp.Content != nil && len(tp.Children) == 0 {
			if err := tp.Content.Create(builder); err != nil {
				return err
			}
		}

		if tp.AssignTo != nil {
			*tp.AssignTo = w
		}

		return nil
	})
}
