// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ProgressBar struct {
	Widget        **walk.ProgressBar
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
}

func (pb ProgressBar) Create(parent walk.Container) error {
	w, err := walk.NewProgressBar(parent)
	if err != nil {
		return err
	}

	return InitWidget(pb, w, func() error {
		if pb.Widget != nil {
			*pb.Widget = w
		}

		return nil
	})
}

func (pb ProgressBar) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return pb.Name, pb.StretchFactor, pb.Row, pb.RowSpan, pb.Column, pb.ColumnSpan
}
