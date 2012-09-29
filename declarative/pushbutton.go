// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type PushButton struct {
	Widget        **walk.PushButton
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Font          Font
	Text          string
	OnClicked     walk.EventHandler
}

func (pb PushButton) Create(parent walk.Container) error {
	w, err := walk.NewPushButton(parent)
	if err != nil {
		return err
	}

	return InitWidget(pb, w, func() error {
		w.SetName(pb.Name)

		if err := w.SetText(pb.Text); err != nil {
			return err
		}

		if pb.OnClicked != nil {
			w.Clicked().Attach(pb.OnClicked)
		}

		if pb.Widget != nil {
			*pb.Widget = w
		}

		return nil
	})
}

func (pb PushButton) LayoutParams() (stretchFactor, row, rowSpan, column, columnSpan int) {
	return pb.StretchFactor, pb.Row, pb.RowSpan, pb.Column, pb.ColumnSpan
}
