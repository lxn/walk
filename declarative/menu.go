// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Menu struct {
	Actions []*walk.Action
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
