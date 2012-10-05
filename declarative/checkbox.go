// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type CheckBox struct {
	AssignTo           **walk.CheckBox
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
	BindTo             string
	Text               string
	OnClicked          walk.EventHandler
}

func (cb CheckBox) Create(parent walk.Container) error {
	w, err := walk.NewCheckBox(parent)
	if err != nil {
		return err
	}

	return InitWidget(cb, w, func() error {
		if err := w.SetBindingMember(cb.BindTo); err != nil {
			return err
		}

		if err := w.SetText(cb.Text); err != nil {
			return err
		}

		if cb.OnClicked != nil {
			w.Clicked().Attach(cb.OnClicked)
		}

		if cb.AssignTo != nil {
			*cb.AssignTo = w
		}

		return nil
	})
}

func (cb CheckBox) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return cb.Name, cb.Disabled, cb.Hidden, &cb.Font, cb.MinSize, cb.MaxSize, cb.StretchFactor, cb.Row, cb.RowSpan, cb.Column, cb.ColumnSpan, cb.ContextMenuActions
}
