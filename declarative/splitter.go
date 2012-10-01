// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Splitter struct {
	AssignTo           **walk.Splitter
	Name               string
	Disabled           bool
	Hidden             bool
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	Layout             Layout
	Children           []Widget
	HandleWidth        int
	Orientation        Orientation
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
		if err := w.SetOrientation(walk.Orientation(s.Orientation)); err != nil {
			return err
		}

		if s.AssignTo != nil {
			*s.AssignTo = w
		}

		return nil
	})
}

func (s Splitter) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return s.Name, s.Disabled, s.Hidden, s.MinSize, s.MaxSize, s.StretchFactor, s.Row, s.RowSpan, s.Column, s.ColumnSpan, s.ContextMenuActions
}

func (s Splitter) ContainerInfo() (Layout, []Widget) {
	return s.Layout, s.Children
}
