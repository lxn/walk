// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "walk/winapi/user32"
)

type clickable interface {
	raiseClicked()
}

type Button struct {
	Widget
	clickedPublisher EventPublisher
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

func (b *Button) Clicked() *Event {
	return b.clickedPublisher.Event()
}

func (b *Button) raiseClicked() {
	b.clickedPublisher.Publish()
}
