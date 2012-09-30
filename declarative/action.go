// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Action struct {
	AssignTo    **walk.Action
	Text        string
	Image       *walk.Bitmap
	OnTriggered walk.EventHandler
}

func (a Action) createMenuAction(menu *walk.Menu) (*walk.Action, error) {
	action := walk.NewAction()

	if _, err := a.initAction(action); err != nil {
		return nil, err
	}

	if menu != nil {
		if err := menu.Actions().Add(action); err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a Action) createToolBarAction() (*walk.Action, error) {
	return a.initAction(walk.NewAction())
}

func (a Action) initAction(wa *walk.Action) (*walk.Action, error) {
	text := a.Text
	if text == "" {
		text = "-"
	}
	if err := wa.SetText(text); err != nil {
		return nil, err
	}
	if err := wa.SetImage(a.Image); err != nil {
		return nil, err
	}

	if a.OnTriggered != nil {
		wa.Triggered().Attach(a.OnTriggered)
	}

	if a.AssignTo != nil {
		*a.AssignTo = wa
	}

	return wa, nil
}

type ActionRef struct {
	Action *walk.Action
}

func (ar ActionRef) createMenuAction(menu *walk.Menu) (*walk.Action, error) {
	if menu != nil {
		if err := menu.Actions().Add(ar.Action); err != nil {
			return nil, err
		}
	}

	return ar.Action, nil
}

func (ar ActionRef) createToolBarAction() (*walk.Action, error) {
	return ar.Action, nil
}

type SubMenu struct {
	Text  string
	Items []MenuItem
}

func (sm SubMenu) createMenuAction(menu *walk.Menu) (*walk.Action, error) {
	if menu == nil {
		var err error
		if menu, err = walk.NewMenu(); err != nil {
			return nil, err
		}
	}

	subMenu, err := walk.NewMenu()
	if err != nil {
		return nil, err
	}

	action, err := menu.Actions().AddMenu(subMenu)
	if err != nil {
		return nil, err
	}

	if err := action.SetText(sm.Text); err != nil {
		return nil, err
	}

	for _, item := range sm.Items {
		if _, err := item.createMenuAction(subMenu); err != nil {
			return nil, err
		}
	}

	return action, nil
}

func addToActionList(list *walk.ActionList, actions []*walk.Action) error {
	for _, a := range actions {
		if err := list.Add(a); err != nil {
			return err
		}
	}

	return nil
}

func CreateMenuActions(items ...MenuItem) ([]*walk.Action, error) {
	var actions []*walk.Action

	for _, item := range items {
		action, err := item.createMenuAction(nil)
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func CreateToolBarActions(items ...ToolBarItem) ([]*walk.Action, error) {
	var actions []*walk.Action

	for _, item := range items {
		action, err := item.createToolBarAction()
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}
