// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"syscall"
)

import (
	. "walk/winapi/user32"
)

func MsgBox(owner RootWidget, title, message string, style uint) int {
	var ownerHWnd HWND

	if owner != nil {
		ownerHWnd = owner.Handle()
	}

	return MessageBox(ownerHWnd, syscall.StringToUTF16Ptr(message), syscall.StringToUTF16Ptr(title), style)
}
