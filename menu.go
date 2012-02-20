// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

type Menu struct {
	hMenu   HMENU
	hWnd    HWND
	actions *ActionList
}

func newMenuBar() (*Menu, error) {
	hMenu := CreateMenu()
	if hMenu == 0 {
		return nil, lastError("CreateMenu")
	}

	m := &Menu{hMenu: hMenu}
	m.actions = newActionList(m)

	return m, nil
}

func NewMenu() (*Menu, error) {
	hMenu := CreatePopupMenu()
	if hMenu == 0 {
		return nil, lastError("CreatePopupMenu")
	}

	var mi MENUINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))

	if !GetMenuInfo(hMenu, &mi) {
		return nil, lastError("GetMenuInfo")
	}

	mi.FMask |= MIM_STYLE
	mi.DwStyle = MNS_CHECKORBMP

	if !SetMenuInfo(hMenu, &mi) {
		return nil, lastError("SetMenuInfo")
	}

	m := &Menu{hMenu: hMenu}
	m.actions = newActionList(m)

	return m, nil
}

func (m *Menu) Dispose() {
	if m.hMenu != 0 {
		DestroyMenu(m.hMenu)
		m.hMenu = 0
	}
}

func (m *Menu) IsDisposed() bool {
	return m.hMenu == 0
}

func (m *Menu) Actions() *ActionList {
	return m.actions
}

func (m *Menu) initMenuItemInfoFromAction(mii *MENUITEMINFO, action *Action) {
	mii.CbSize = uint32(unsafe.Sizeof(*mii))
	mii.FMask = MIIM_FTYPE | MIIM_ID | MIIM_STATE | MIIM_STRING
	if action.image != nil {
		mii.FMask |= MIIM_BITMAP
		mii.HbmpItem = action.image.handle()
	}
	if action.text == "-" {
		mii.FType = MFT_SEPARATOR
	} else {
		mii.FType = MFT_STRING
		mii.DwTypeData = syscall.StringToUTF16Ptr(action.text)
		mii.Cch = uint32(len([]rune(action.text)))
	}
	mii.WID = uint32(action.id)

	if action.Enabled() {
		mii.FState &^= MFS_DISABLED
	} else {
		mii.FState |= MFS_DISABLED
	}

	menu := action.menu
	if menu != nil {
		mii.FMask |= MIIM_SUBMENU
		mii.HSubMenu = menu.hMenu
	}
}

func (m *Menu) onActionChanged(action *Action) error {
	var mii MENUITEMINFO

	m.initMenuItemInfoFromAction(&mii, action)

	if !SetMenuItemInfo(m.hMenu, uint32(m.actions.Index(action)), true, &mii) {
		return newError("SetMenuItemInfo failed")
	}

	return nil
}

func (m *Menu) onInsertingAction(index int, action *Action) error {
	var mii MENUITEMINFO

	m.initMenuItemInfoFromAction(&mii, action)

	if !InsertMenuItem(m.hMenu, uint32(index), true, &mii) {
		return newError("InsertMenuItem failed")
	}

	action.addChangedHandler(m)

	menu := action.menu
	if menu != nil {
		menu.hWnd = m.hWnd
	}

	if m.hWnd != 0 {
		DrawMenuBar(m.hWnd)
	}

	return nil
}

func (m *Menu) onRemovingAction(index int, action *Action) error {
	if !RemoveMenu(m.hMenu, uint32(index), MF_BYPOSITION) {
		return lastError("RemoveMenu")
	}

	action.removeChangedHandler(m)

	if m.hWnd != 0 {
		DrawMenuBar(m.hWnd)
	}

	return nil
}

func (m *Menu) onClearingActions() error {
	for i := m.actions.Len() - 1; i >= 0; i-- {
		if !RemoveMenu(m.hMenu, uint32(i), MF_BYPOSITION) {
			return lastError("RemoveMenu")
		}

		m.actions.At(i).removeChangedHandler(m)
	}

	if m.hWnd != 0 {
		DrawMenuBar(m.hWnd)
	}

	return nil
}
