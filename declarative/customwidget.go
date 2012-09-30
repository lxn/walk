// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type CustomWidget struct {
	AssignTo            **walk.CustomWidget
	Name                string
	MinSize             Size
	MaxSize             Size
	StretchFactor       int
	Row                 int
	RowSpan             int
	Column              int
	ColumnSpan          int
	ContextMenuActions  []*walk.Action
	Style               uint32
	Paint               walk.PaintFunc
	ClearsBackground    bool
	InvalidatesOnResize bool
}

func (cw CustomWidget) Create(parent walk.Container) error {
	w, err := walk.NewCustomWidget(parent, uint(cw.Style), cw.Paint)
	if err != nil {
		return err
	}

	return InitWidget(cw, w, func() error {
		w.SetClearsBackground(cw.ClearsBackground)
		w.SetInvalidatesOnResize(cw.InvalidatesOnResize)

		if cw.AssignTo != nil {
			*cw.AssignTo = w
		}

		return nil
	})
}

func (cw CustomWidget) WidgetInfo() (name string, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return cw.Name, cw.MinSize, cw.MaxSize, cw.StretchFactor, cw.Row, cw.RowSpan, cw.Column, cw.ColumnSpan, cw.ContextMenuActions
}
