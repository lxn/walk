// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user32

import (
	. "walk/winapi/gdi32"
)

// Constants for MENUITEMINFO.fMask
const (
	MIIM_STATE      = 1
	MIIM_ID         = 2
	MIIM_SUBMENU    = 4
	MIIM_CHECKMARKS = 8
	MIIM_TYPE       = 16
	MIIM_DATA       = 32
	MIIM_STRING     = 64
	MIIM_BITMAP     = 128
	MIIM_FTYPE      = 256
)

// Constants for MENUITEMINFO.fType
const (
	MFT_BITMAP       = 4
	MFT_MENUBARBREAK = 32
	MFT_MENUBREAK    = 64
	MFT_OWNERDRAW    = 256
	MFT_RADIOCHECK   = 512
	MFT_RIGHTJUSTIFY = 0x4000
	MFT_SEPARATOR    = 0x800
	MFT_RIGHTORDER   = 0x2000
	MFT_STRING       = 0
)

// Constants for MENUITEMINFO.fState
const (
	MFS_CHECKED   = 8
	MFS_DEFAULT   = 4096
	MFS_DISABLED  = 3
	MFS_ENABLED   = 0
	MFS_GRAYED    = 3
	MFS_HILITE    = 128
	MFS_UNCHECKED = 0
	MFS_UNHILITE  = 0
)

// Constants for MENUITEMINFO.hbmp*
const (
	HBMMENU_CALLBACK        = -1
	HBMMENU_SYSTEM          = 1
	HBMMENU_MBAR_RESTORE    = 2
	HBMMENU_MBAR_MINIMIZE   = 3
	HBMMENU_MBAR_CLOSE      = 5
	HBMMENU_MBAR_CLOSE_D    = 6
	HBMMENU_MBAR_MINIMIZE_D = 7
	HBMMENU_POPUP_CLOSE     = 8
	HBMMENU_POPUP_RESTORE   = 9
	HBMMENU_POPUP_MAXIMIZE  = 10
	HBMMENU_POPUP_MINIMIZE  = 11
)

type MENUITEMINFO struct {
	CbSize        uint
	FMask         uint
	FType         uint
	FState        uint
	WID           uint
	HSubMenu      HMENU
	HbmpChecked   HBITMAP
	HbmpUnchecked HBITMAP
	DwItemData    uintptr
	DwTypeData    *uint16
	Cch           uint
	HbmpItem      HBITMAP
}
