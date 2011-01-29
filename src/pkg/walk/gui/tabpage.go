// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	"walk/drawing"
)

import (
	. "walk/winapi/user32"
)

const tabPageWindowClass = `\o/ Walk_TabPage_Class \o/`

var tabPageWndProcCallback *syscall.Callback

func tabPageWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	tp, ok := widgetsByHWnd[msg.HWnd].(*TabPage)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return tp.wndProc(msg, 0)
}

type TabPage struct {
	Container
}

func NewTabPage() (*TabPage, os.Error) {
	ensureRegisteredWindowClass(tabPageWindowClass, tabPageWndProc, &tabPageWndProcCallback)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(tabPageWindowClass), nil,
		WS_POPUP,
		0, 0, 0, 0, 0, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	tp := &TabPage{
		Container: Container{
			Widget: Widget{
				hWnd: hWnd,
			},
		},
	}

	succeeded := false
	defer func() {
		if !succeeded {
			tp.Dispose()
		}
	}()

	tp.children = newObservedWidgetList(tp)

	tp.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = tp

	succeeded = true

	return tp, nil
}

func (*TabPage) LayoutFlags() LayoutFlags {
	return GrowHorz | GrowVert | ShrinkHorz | ShrinkVert
}

func (tp *TabPage) PreferredSize() drawing.Size {
	return tp.dialogBaseUnitsToPixels(drawing.Size{100, 100})
}
