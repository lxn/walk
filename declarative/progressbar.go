// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ProgressBar struct {
	AssignTo           **walk.ProgressBar
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
	MinValue           int
	MaxValue           int
	Value              int
}

func (pb ProgressBar) Create(parent walk.Container) error {
	w, err := walk.NewProgressBar(parent)
	if err != nil {
		return err
	}

	return InitWidget(pb, w, func() error {
		w.SetRange(pb.MinValue, pb.MaxValue)
		w.SetValue(pb.Value)

		if pb.AssignTo != nil {
			*pb.AssignTo = w
		}

		return nil
	})
}

func (pb ProgressBar) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return pb.Name, pb.Disabled, pb.Hidden, &pb.Font, pb.MinSize, pb.MaxSize, pb.StretchFactor, pb.Row, pb.RowSpan, pb.Column, pb.ColumnSpan, pb.ContextMenuActions
}
