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
	Font               Font
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
	Title              string
	Size               Size
	DataBinder         DataBinder
	Layout             Layout
	Children           []Widget
	DefaultButton      **walk.PushButton
	CancelButton       **walk.PushButton
}

func (d Dialog) Create(owner walk.RootWidget) error {
	w, err := walk.NewDialog(owner)
	if err != nil {
		return err
	}

	tlwi := topLevelWindowInfo{
		Name:               d.Name,
		Disabled:           d.Disabled,
		Hidden:             d.Hidden,
		Font:               d.Font,
		MinSize:            d.MinSize,
		MaxSize:            d.MaxSize,
		ContextMenuActions: d.ContextMenuActions,
		DataBinder:         d.DataBinder,
		Layout:             d.Layout,
		Children:           d.Children,
	}

	var db *walk.DataBinder
	if d.DataBinder.AssignTo == nil {
		d.DataBinder.AssignTo = &db
	}

	return InitWidget(tlwi, w, func() error {
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

			db := *d.DataBinder.AssignTo
			if db != nil {
				db.CanSubmitChanged().Attach(func() {
					(*d.DefaultButton).SetEnabled(db.CanSubmit())
				})
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
