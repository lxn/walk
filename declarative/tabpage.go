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
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
	DataBinder         DataBinder
	Layout             Layout
	Children           []Widget
	Title              string
	Content            Widget
}

func (tp TabPage) Create(parent walk.Container) error {
	w, err := walk.NewTabPage()
	if err != nil {
		return err
	}

	return InitWidget(tp, w, func() error {
		if err := w.SetTitle(tp.Title); err != nil {
			return err
		}

		if tp.Content != nil && len(tp.Children) == 0 {
			if err := tp.Content.Create(w); err != nil {
				return err
			}
		}

		if tp.AssignTo != nil {
			*tp.AssignTo = w
		}

		return nil
	})
}

func (tp TabPage) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return tp.Name, tp.Disabled, tp.Hidden, &tp.Font, tp.MinSize, tp.MaxSize, 0, 0, 0, 0, 0, tp.ContextMenuActions
}

func (tp TabPage) ContainerInfo() (DataBinder, Layout, []Widget) {
	return tp.DataBinder, tp.Layout, tp.Children
}
