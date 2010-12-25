// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	. "walk/winapi/user32"
)

type clickable interface {
	raiseClicked()
}

type Button struct {
	Widget
	clickedHandlers []EventHandler
}

func (b *Button) Checked() bool {
	return SendMessage(b.hWnd, BM_GETCHECK, 0, 0) == BST_CHECKED
}

func (b *Button) SetChecked(value bool) {
	var chk uintptr

	if value {
		chk = BST_CHECKED
	} else {
		chk = BST_UNCHECKED
	}

	SendMessage(b.hWnd, BM_SETCHECK, chk, 0)
}

func (b *Button) AddClickedHandler(handler EventHandler) {
	b.clickedHandlers = append(b.clickedHandlers, handler)
}

func (b *Button) RemoveClickedHandler(handler EventHandler) {
	for i, h := range b.clickedHandlers {
		if h == handler {
			b.clickedHandlers = append(b.clickedHandlers[:i], b.clickedHandlers[i+1:]...)
			break
		}
	}
}

func (b *Button) raiseClicked() {
	args := &eventArgs{widgetsByHWnd[b.hWnd]}
	for _, handler := range b.clickedHandlers {
		handler(args)
	}
}
