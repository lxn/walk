// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"os"
)

import (
	. "walk/winapi/user32"
)

type Button struct {
	Widget
	clickedHandlers vector.Vector
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
	b.clickedHandlers.Push(handler)
}

func (b *Button) RemoveClickedHandler(handler EventHandler) {
	for i, h := range b.clickedHandlers {
		if h.(EventHandler) == handler {
			b.clickedHandlers.Delete(i)
			break
		}
	}
}

func (b *Button) raiseClicked() {
	for _, handlerIface := range b.clickedHandlers {
		handler := handlerIface.(EventHandler)
		handler(&eventArgs{widgetsByHWnd[b.hWnd]})
	}
}

func (b *Button) raiseEvent(msg *MSG) os.Error {
	b.Widget.raiseEvent(msg)

	switch msg.Message {
	case WM_KEYUP:
		if msg.WParam != VK_SPACE {
			break
		}

		b.raiseClicked()

	case WM_LBUTTONUP:
		if trackedMouseDownHWnd != b.Widget.hWnd {
			break
		}
		bounds, err := b.ClientBounds()
		if err != nil {
			return err
		}

		x, y := GET_X_LPARAM(msg.LParam), GET_Y_LPARAM(msg.LParam)
		if x < 0 || x >= bounds.Width || y < 0 || y >= bounds.Height {
			break
		}

		b.raiseClicked()
	}

	return nil
}
