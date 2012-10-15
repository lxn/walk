// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type LineErrorPresenter struct {
	AssignTo           *walk.ErrorPresenter
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
}

func (lep LineErrorPresenter) Create(parent walk.Container) error {
	w, err := walk.NewLineErrorPresenter(parent)
	if err != nil {
		return err
	}

	return InitWidget(lep, w, func() error {
		if lep.AssignTo != nil {
			*lep.AssignTo = w
		}

		return nil
	})
}

func (lep LineErrorPresenter) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return lep.Name, lep.Disabled, lep.Hidden, &lep.Font, lep.MinSize, lep.MaxSize, lep.StretchFactor, lep.Row, lep.RowSpan, lep.Column, lep.ColumnSpan, lep.ContextMenuActions
}
