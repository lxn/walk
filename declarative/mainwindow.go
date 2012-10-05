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
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
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
		MinSize:            mw.MinSize,
		MaxSize:            mw.MaxSize,
		ContextMenuActions: mw.ContextMenuActions,
		DataBinder:         mw.DataBinder,
		Layout:             mw.Layout,
		Children:           mw.Children,
	}

	return InitWidget(tlwi, w, func() error {
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
