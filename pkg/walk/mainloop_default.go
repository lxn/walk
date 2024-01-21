// Copyright 2019 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows,!walk_use_cgo

package walk

import (
	"unsafe"

	"github.com/miu200521358/win"
)

func (fb *FormBase) mainLoop() int {
	msg := (*win.MSG)(unsafe.Pointer(win.GlobalAlloc(0, unsafe.Sizeof(win.MSG{}))))
	defer win.GlobalFree(win.HGLOBAL(unsafe.Pointer(msg)))

	for fb.hWnd != 0 {
		switch win.GetMessage(msg, 0, 0, 0) {
		case 0:
			return int(msg.WParam)

		case -1:
			return -1
		}

		switch msg.Message {
		case win.WM_KEYDOWN:
			if fb.handleKeyDown(msg) {
				continue
			}
		}

		if !win.IsDialogMessage(fb.hWnd, msg) {
			win.TranslateMessage(msg)
			win.DispatchMessage(msg)
		}

		fb.group.RunSynchronized()
	}

	return 0
}
