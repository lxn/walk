// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Dialog struct {
	AssignTo           **walk.Dialog
	Name               string
	Disabled           bool
	Hidden             bool
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
	Font               Font
	Title              string
	Size               Size
	Layout             Layout
	Children           []Widget
	DefaultButton      **walk.PushButton
	CancelButton       **walk.PushButton
}

func (d Dialog) Create(parent walk.Container) error {
	var owner walk.RootWidget
	if o, ok := parent.(walk.RootWidget); ok {
		owner = o
	}

	w, err := walk.NewDialog(owner)
	if err != nil {
		return err
	}

	return InitWidget(d, w, func() error {
		if err := w.SetTitle(d.Title); err != nil {
			return err
		}

		if err := w.SetSize(d.Size.toW()); err != nil {
			return err
		}

		if d.DefaultButton != nil {
			if err := w.SetDefaultButton(*d.DefaultButton); err != nil {
				return err
			}
		}
		if d.CancelButton != nil {
			if err := w.SetCancelButton(*d.CancelButton); err != nil {
				return err
			}
		}

		if d.AssignTo != nil {
			*d.AssignTo = w
		}

		return nil
	})
}

func (d Dialog) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return d.Name, d.Disabled, d.Hidden, d.MinSize, d.MaxSize, 0, 0, 0, 0, 0, d.ContextMenuActions
}

func (d Dialog) Font_() *Font {
	return &d.Font
}

func (d Dialog) ContainerInfo() (Layout, []Widget) {
	return d.Layout, d.Children
}
