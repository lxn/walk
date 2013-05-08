// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ToolBar struct {
	AssignTo         **walk.ToolBar
	Name             string
	Enabled          Property
	Visible          Property
	Font             Font
	ToolTipText      Property
	MinSize          Size
	MaxSize          Size
	StretchFactor    int
	Row              int
	RowSpan          int
	Column           int
	ColumnSpan       int
	ContextMenuItems []MenuItem
	OnKeyDown        walk.KeyEventHandler
	OnKeyUp          walk.KeyEventHandler
	OnMouseDown      walk.MouseEventHandler
	OnMouseMove      walk.MouseEventHandler
	OnMouseUp        walk.MouseEventHandler
	OnSizeChanged    walk.EventHandler
	Actions          []*walk.Action
	MaxTextRows      int
	Orientation      Orientation
}

func (tb ToolBar) Create(builder *Builder) (err error) {
	var w *walk.ToolBar
	if tb.Orientation == Vertical {
		w, err = walk.NewVerticalToolBar(builder.Parent())
	} else {
		w, err = walk.NewToolBar(builder.Parent())
	}
	if err != nil {
		return
	}

	return builder.InitWidget(tb, w, func() error {
		imageList, err := walk.NewImageList(walk.Size{16, 16}, 0)
		if err != nil {
			return err
		}
		w.SetImageList(imageList)

		mtr := tb.MaxTextRows
		if mtr < 1 {
			mtr = 1
		}
		if err := w.SetMaxTextRows(mtr); err != nil {
			return err
		}

		if err := addToActionList(w.Actions(), tb.Actions); err != nil {
			return err
		}

		if tb.AssignTo != nil {
			*tb.AssignTo = w
		}

		return nil
	})
}

func (w ToolBar) WidgetInfo() (name string, disabled, hidden bool, font *Font, toolTipText string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuItems []MenuItem, OnKeyDown walk.KeyEventHandler, OnKeyUp walk.KeyEventHandler, OnMouseDown walk.MouseEventHandler, OnMouseMove walk.MouseEventHandler, OnMouseUp walk.MouseEventHandler, OnSizeChanged walk.EventHandler) {
	return w.Name, false, false, &w.Font, "", w.MinSize, w.MaxSize, w.StretchFactor, w.Row, w.RowSpan, w.Column, w.ColumnSpan, w.ContextMenuItems, w.OnKeyDown, w.OnKeyUp, w.OnMouseDown, w.OnMouseMove, w.OnMouseUp, w.OnSizeChanged
}
