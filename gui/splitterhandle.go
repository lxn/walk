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

const splitterHandleWindowClass = `\o/ Walk_SplitterHandle_Class \o/`

var splitterHandleWndProcCallback *syscall.Callback

func splitterHandleWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	s, ok := widgetsByHWnd[msg.HWnd].(*splitterHandle)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return s.wndProc(msg, 0)
}

type splitterHandle struct {
	Widget
}

func newSplitterHandle(splitter *Splitter) (*splitterHandle, os.Error) {
	if splitter == nil {
		return nil, newError("splitter cannot be nil")
	}

	ensureRegisteredWindowClass(splitterHandleWindowClass, splitterHandleWndProc, &splitterHandleWndProcCallback)

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr(splitterHandleWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 0, 0, splitter.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	sh := &splitterHandle{Widget: Widget{hWnd: hWnd, parent: splitter}}

	sh.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = sh

	return sh, nil
}

func (sh *splitterHandle) LayoutFlags() LayoutFlags {
	splitter := sh.Parent().(*Splitter)

	if splitter.Orientation() == Horizontal {
		return GrowVert
	}

	return GrowHorz
}

func (sh *splitterHandle) PreferredSize() drawing.Size {
	splitter := sh.Parent().(*Splitter)
	handleWidth := splitter.HandleWidth()
	var size drawing.Size

	if splitter.Orientation() == Horizontal {
		size.Width = handleWidth
	} else {
		size.Height = handleWidth
	}

	return size
}
