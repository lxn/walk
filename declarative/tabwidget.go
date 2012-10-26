// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TabWidget struct {
	AssignTo           **walk.TabWidget
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
	ContentMargins     Margins
	ContentMarginsZero bool
	Pages              []TabPage
}

func (tw TabWidget) Create(parent walk.Container) error {
	w, err := walk.NewTabWidget(parent)
	if err != nil {
		return err
	}

	return InitWidget(tw, w, func() error {
		for _, tp := range tw.Pages {
			var wp *walk.TabPage
			if tp.AssignTo == nil {
				tp.AssignTo = &wp
			}

			if tp.Content != nil && len(tp.Children) == 0 {
				tp.Layout = HBox{Margins: tw.ContentMargins, MarginsZero: tw.ContentMarginsZero}
			}

			if err := tp.Create(nil); err != nil {
				return err
			}

			if err := w.Pages().Add(*tp.AssignTo); err != nil {
				return err
			}
		}

		if tw.AssignTo != nil {
			*tw.AssignTo = w
		}

		return nil
	})
}

func (w TabWidget) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action, OnKeyDown walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, w.Disabled, w.Hidden, &w.Font, w.ToolTipText, w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuActions, w.OnKeyDown, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}
