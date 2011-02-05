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

const tabPageWindowClass = `\o/ Walk_TabPage_Class \o/`

var tabPageWndProcPtr uintptr

func tabPageWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	tp, ok := widgetsByHWnd[hwnd].(*TabPage)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return tp.wndProc(hwnd, msg, wParam, lParam, 0)
}

type TabPage struct {
	Container
}

func NewTabPage() (*TabPage, os.Error) {
	ensureRegisteredWindowClass(tabPageWindowClass, tabPageWndProc, &tabPageWndProcPtr)

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

func (tp *TabPage) PreferredSize() Size {
	return tp.dialogBaseUnitsToPixels(Size{100, 100})
}
