// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

const splitterHandleWindowClass = `\o/ Walk_SplitterHandle_Class \o/`

var splitterHandleWndProcPtr uintptr

func splitterHandleWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	s, ok := widgetsByHWnd[hwnd].(*splitterHandle)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return s.wndProc(hwnd, msg, wParam, lParam, 0)
}

type splitterHandle struct {
	Widget
}

func newSplitterHandle(splitter *Splitter) (*splitterHandle, os.Error) {
	if splitter == nil {
		return nil, newError("splitter cannot be nil")
	}

	ensureRegisteredWindowClass(splitterHandleWindowClass, splitterHandleWndProc, &splitterHandleWndProcPtr)

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr(splitterHandleWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 0, 0, splitter.BaseWidget().hWnd, 0, 0, nil)
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

func (sh *splitterHandle) PreferredSize() Size {
	splitter := sh.Parent().(*Splitter)
	handleWidth := splitter.HandleWidth()
	var size Size

	if splitter.Orientation() == Horizontal {
		size.Width = handleWidth
	} else {
		size.Height = handleWidth
	}

	return size
}
