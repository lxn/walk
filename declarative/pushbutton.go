// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type PushButton struct {
	AssignTo           **walk.PushButton
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
	Font               Font
	Text               string
	OnClicked          walk.EventHandler
}

func (pb PushButton) Create(parent walk.Container) error {
	w, err := walk.NewPushButton(parent)
	if err != nil {
		return err
	}

	return InitWidget(pb, w, func() error {
		if err := w.SetText(pb.Text); err != nil {
			return err
		}

		if pb.OnClicked != nil {
			w.Clicked().Attach(pb.OnClicked)
		}

		if pb.AssignTo != nil {
			*pb.AssignTo = w
		}

		return nil
	})
}

func (pb PushButton) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return pb.Name, pb.Disabled, pb.Hidden, pb.MinSize, pb.MaxSize, pb.StretchFactor, pb.Row, pb.RowSpan, pb.Column, pb.ColumnSpan, pb.ContextMenuActions
}

func (pb PushButton) Font_() *Font {
	return &pb.Font
}
