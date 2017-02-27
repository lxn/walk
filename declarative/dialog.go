// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type Dialog struct {
	AssignTo         **walk.Dialog
	Name             string
	Enabled          Property
	Visible          Property
	Font             Font
	MinSize          Size
	MaxSize          Size
	ContextMenuItems []MenuItem
	OnKeyDown        walk.KeyEventHandler
	OnKeyPress       walk.KeyEventHandler
	OnKeyUp          walk.KeyEventHandler
	OnMouseDown      walk.MouseEventHandler
	OnMouseMove      walk.MouseEventHandler
	OnMouseUp        walk.MouseEventHandler
	OnSizeChanged    walk.EventHandler
	Icon             Property
	Title            Property
	Size             Size
	DataBinder       DataBinder
	Layout           Layout
	Children         []Widget
	DefaultButton    **walk.PushButton
	CancelButton     **walk.PushButton
	FixedSize        bool
	Expressions      func() map[string]walk.Expression
	Functions        map[string]func(args ...interface{}) (interface{}, error)
}

func (d Dialog) Create(owner walk.Form) error {
	var w *walk.Dialog
	var err error

	if d.FixedSize {
		w, err = walk.NewDialogWithFixedSize(owner)
	} else {
		w, err = walk.NewDialog(owner)
	}

	if err != nil {
		return err
	}

	tlwi := topLevelWindowInfo{
		Name:             d.Name,
		Enabled:          d.Enabled,
		Visible:          d.Visible,
		Font:             d.Font,
		ToolTipText:      "",
		MinSize:          d.MinSize,
		MaxSize:          d.MaxSize,
		ContextMenuItems: d.ContextMenuItems,
		DataBinder:       d.DataBinder,
		Layout:           d.Layout,
		Children:         d.Children,
		OnKeyDown:        d.OnKeyDown,
		OnKeyPress:       d.OnKeyPress,
		OnKeyUp:          d.OnKeyUp,
		OnMouseDown:      d.OnMouseDown,
		OnMouseMove:      d.OnMouseMove,
		OnMouseUp:        d.OnMouseUp,
		OnSizeChanged:    d.OnSizeChanged,
		Icon:             d.Icon,
		Title:            d.Title,
	}

	var db *walk.DataBinder
	if d.DataBinder.AssignTo == nil {
		d.DataBinder.AssignTo = &db
	}

	builder := NewBuilder(nil)

	w.SetSuspended(true)
	builder.Defer(func() error {
		w.SetSuspended(false)
		return nil
	})

	return builder.InitWidget(tlwi, w, func() error {
		if err := w.SetSize(d.Size.toW()); err != nil {
			return err
		}

		if d.DefaultButton != nil {
			if err := w.SetDefaultButton(*d.DefaultButton); err != nil {
				return err
			}

			if db := *d.DataBinder.AssignTo; db != nil {
				if db.DataSource() != nil {
					(*d.DefaultButton).SetEnabled(db.CanSubmit())
				}

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

		if d.Expressions != nil {
			for name, expr := range d.Expressions() {
				builder.expressions[name] = expr
			}
		}
		if d.Functions != nil {
			for name, fn := range d.Functions {
				builder.functions[name] = fn
			}
		}

		return nil
	})
}

func (d Dialog) Run(owner walk.Form) (int, error) {
	var w *walk.Dialog

	if d.AssignTo == nil {
		d.AssignTo = &w
	}

	if err := d.Create(owner); err != nil {
		return 0, err
	}

	return (*d.AssignTo).Run(), nil
}
