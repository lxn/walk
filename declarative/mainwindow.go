// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type MainWindow struct {
	AssignTo           **walk.MainWindow
	Name               string
	Disabled           bool
	Hidden             bool
	Font               Font
	ToolTipText        string
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
	OnKeyDown          walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	Title              string
	Size               Size
	DataBinder         DataBinder
	Layout             Layout
	Children           []Widget
	MenuActions        []*walk.Action
	ToolBarActions     []*walk.Action
}

func (mw MainWindow) Create() error {
	w, err := walk.NewMainWindow()
	if err != nil {
		return err
	}

	tlwi := topLevelWindowInfo{
		Name:               mw.Name,
		Disabled:           mw.Disabled,
		Hidden:             mw.Hidden,
		Font:               mw.Font,
		ToolTipText:        mw.ToolTipText,
		MinSize:            mw.MinSize,
		MaxSize:            mw.MaxSize,
		ContextMenuActions: mw.ContextMenuActions,
		OnKeyDown:          mw.OnKeyDown,
		OnMouseDown:        mw.OnMouseDown,
		OnMouseMove:        mw.OnMouseMove,
		OnMouseUp:          mw.OnMouseUp,
		OnSizeChanged:      mw.OnSizeChanged,
		DataBinder:         mw.DataBinder,
		Layout:             mw.Layout,
		Children:           mw.Children,
	}

	builder := NewBuilder(nil)

	return builder.InitWidget(tlwi, w, func() error {
		if err := w.SetTitle(mw.Title); err != nil {
			return err
		}

		if err := w.SetSize(mw.Size.toW()); err != nil {
			return err
		}

		if err := addToActionList(w.Menu().Actions(), mw.MenuActions); err != nil {
			return err
		}

		imageList, err := walk.NewImageList(walk.Size{16, 16}, 0)
		if err != nil {
			return err
		}
		w.ToolBar().SetImageList(imageList)

		if err := addToActionList(w.ToolBar().Actions(), mw.ToolBarActions); err != nil {
			return err
		}

		if mw.AssignTo != nil {
			*mw.AssignTo = w
		}

		return nil
	})
}

func (mw MainWindow) Run() (int, error) {
	var w *walk.MainWindow

	if mw.AssignTo == nil {
		mw.AssignTo = &w
	}

	if err := mw.Create(); err != nil {
		return 0, err
	}

	return (*mw.AssignTo).Run(), nil
}
