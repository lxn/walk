// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

func RegisterGlobalHotKey(owner Form, hkid int, hKey Shortcut) bool {
	var ownerHWnd win.HWND

	if owner != nil {
		ownerHWnd = owner.Handle()
	}

	var modifiers uint
	if hKey.Modifiers & ModAlt != 0 {
		modifiers |= 1
	}
	if hKey.Modifiers & ModControl != 0 {
		modifiers |= 2
	}
	if hKey.Modifiers & ModShift != 0 {
		modifiers |= 4
	}
	// https://msdn.microsoft.com/ru-ru/library/windows/desktop/ms646309.aspx
	return win.RegisterHotKey(ownerHWnd, hkid, modifiers, uint(hKey.Key))
}
