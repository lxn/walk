// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TabPage struct {
	AssignTo           **walk.TabPage
	Name               string
	Disabled           bool
	Hidden             bool
	Font               Font
	ToolTipText        string
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
	OnKeyDown          walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	DataBinder         DataBinder
	Layout             Layout
	Children           []Widget
	Title              string
	Content            Widget
}

func (tp TabPage) Create(builder *Builder) error {
	w, err := walk.NewTabPage()
	if err != nil {
		return err
	}

	return builder.InitWidget(tp, w, func() error {
		if err := w.SetTitle(tp.Title); err != nil {
			return err
		}

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

func (w TabPage) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action, OnKeyDown walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, w.Disabled, w.Hidden, &w.Font, w.ToolTipText, w.MinSize, w.MaxSize, 0, 0, 0, 0, 0, w.ContextMenuActions, w.OnKeyDown, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}

func (tp TabPage) ContainerInfo() (DataBinder, Layout, []Widget) {
	return tp.DataBinder, tp.Layout, tp.Children
}
