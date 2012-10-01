// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ImageView struct {
	AssignTo           **walk.ImageView
	Name               string
	Disabled           bool
	Hidden             bool
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	Image              walk.Image
}

func (iv ImageView) Create(parent walk.Container) error {
	w, err := walk.NewImageView(parent)
	if err != nil {
		return err
	}

	return InitWidget(iv, w, func() error {
		if err := w.SetImage(iv.Image); err != nil {
			return err
		}

		if iv.AssignTo != nil {
			*iv.AssignTo = w
		}

		return nil
	})
}

func (iv ImageView) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return iv.Name, iv.Disabled, iv.Hidden, iv.MinSize, iv.MaxSize, iv.StretchFactor, iv.Row, iv.RowSpan, iv.Column, iv.ColumnSpan, iv.ContextMenuActions
}
