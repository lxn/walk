// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type GroupBox struct {
	AssignTo           **walk.GroupBox
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
	Title              string
	Layout             Layout
	Children           []Widget
}

func (gb GroupBox) Create(parent walk.Container) error {
	w, err := walk.NewGroupBox(parent)
	if err != nil {
		return err
	}

	return InitWidget(gb, w, func() error {
		if err := w.SetTitle(gb.Title); err != nil {
			return err
		}

		if gb.AssignTo != nil {
			*gb.AssignTo = w
		}

		return nil
	})
}

func (gb GroupBox) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return gb.Name, gb.Disabled, gb.Hidden, &gb.Font, gb.MinSize, gb.MaxSize, gb.StretchFactor, gb.Row, gb.RowSpan, gb.Column, gb.ColumnSpan, gb.ContextMenuActions
}

func (gb GroupBox) ContainerInfo() (Layout, []Widget) {
	return gb.Layout, gb.Children
}
