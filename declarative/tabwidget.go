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
	Margins            Margins
	MarginsZero        bool
	PageTitles         []string
	Pages              []Widget
}

func (tw TabWidget) Create(parent walk.Container) error {
	w, err := walk.NewTabWidget(parent)
	if err != nil {
		return err
	}

	return InitWidget(tw, w, func() error {
		for i, page := range tw.Pages {
			wp, err := walk.NewTabPage()
			if err != nil {
				return err
			}

			if len(tw.PageTitles) > i {
				if err := wp.SetTitle(tw.PageTitles[i]); err != nil {
					return err
				}
			}

			if err := w.Pages().Add(wp); err != nil {
				return err
			}

			l := walk.NewHBoxLayout()
			m := tw.Margins
			if !tw.MarginsZero && m.isZero() {
				m = Margins{9, 9, 9, 9}
			}

			if err := l.SetMargins(m.toW()); err != nil {
				return err
			}

			if err := wp.SetLayout(l); err != nil {
				return err
			}

			if err := page.Create(wp); err != nil {
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
