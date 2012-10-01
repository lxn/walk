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
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	Pages              []TabPage
}

func (tw TabWidget) Create(parent walk.Container) error {
	w, err := walk.NewTabWidget(parent)
	if err != nil {
		return err
	}

	return InitWidget(tw, w, func() error {
		var p *walk.TabPage

		for _, page := range tw.Pages {
			if page.AssignTo == nil {
				page.AssignTo = &p
			}

			if err := page.Create(nil); err != nil {
				return err
			}

			if err := w.Pages().Add(p); err != nil {
				return err
			}
		}

		if tw.AssignTo != nil {
			*tw.AssignTo = w
		}

		return nil
	})
}

func (tw TabWidget) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return tw.Name, tw.Disabled, tw.Hidden, &tw.Font, tw.MinSize, tw.MaxSize, tw.StretchFactor, tw.Row, tw.RowSpan, tw.Column, tw.ColumnSpan, tw.ContextMenuActions
}
