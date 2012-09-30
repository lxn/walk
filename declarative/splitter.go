// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Splitter struct {
	AssignTo      **walk.Splitter
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	ContextMenu   Menu
	Layout        Layout
	Children      []Widget
	HandleWidth   int
	Orientation   walk.Orientation
}

func (s Splitter) Create(parent walk.Container) error {
	w, err := walk.NewSplitter(parent)
	if err != nil {
		return err
	}

	return InitWidget(s, w, func() error {
		if s.HandleWidth > 0 {
			if err := w.SetHandleWidth(s.HandleWidth); err != nil {
				return err
			}
		}
		if err := w.SetOrientation(s.Orientation); err != nil {
			return err
		}

		if s.AssignTo != nil {
			*s.AssignTo = w
		}

		return nil
	})
}

func (s Splitter) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return s.Name, s.StretchFactor, s.Row, s.RowSpan, s.Column, s.ColumnSpan, &s.ContextMenu
}

func (s Splitter) ContainerInfo() (Layout, []Widget) {
	return s.Layout, s.Children
}
