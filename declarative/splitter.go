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
	Font               Font
	ToolTipText        string
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	OnKeyDown          walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	DataBinder         DataBinder
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

func (w Splitter) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action, OnKeyDown walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, w.Disabled, w.Hidden, &w.Font, w.ToolTipText, w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuActions, w.OnKeyDown, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}

func (s Splitter) ContainerInfo() (DataBinder, Layout, []Widget) {
	return s.DataBinder, nil, s.Children
}
