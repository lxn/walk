// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ProgressBar struct {
	AssignTo      **walk.ProgressBar
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	ContextMenu   Menu
	MinValue      int
	MaxValue      int
	Value         int
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

func (pb ProgressBar) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return pb.Name, pb.StretchFactor, pb.Row, pb.RowSpan, pb.Column, pb.ColumnSpan, &pb.ContextMenu
}
