// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type RadioButton struct {
	AssignTo           **walk.RadioButton
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
	Text               string
	OnClicked          walk.EventHandler
}

func (rb RadioButton) Create(parent walk.Container) error {
	w, err := walk.NewRadioButton(parent)
	if err != nil {
		return err
	}

	return InitWidget(rb, w, func() error {
		if err := w.SetText(rb.Text); err != nil {
			return err
		}

		if rb.OnClicked != nil {
			w.Clicked().Attach(rb.OnClicked)
		}

		if rb.AssignTo != nil {
			*rb.AssignTo = w
		}

		return nil
	})
}

func (rb RadioButton) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return rb.Name, rb.Disabled, rb.Hidden, &rb.Font, rb.MinSize, rb.MaxSize, rb.StretchFactor, rb.Row, rb.RowSpan, rb.Column, rb.ColumnSpan, rb.ContextMenuActions
}
