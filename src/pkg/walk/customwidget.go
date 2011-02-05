// Copyright 2010 The Walk Authors. All rights reserved.
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

const customWidgetWindowClass = `\o/ Walk_CustomWidget_Class \o/`

var customWidgetWndProcPtr uintptr

func customWidgetWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	cw, ok := customWidgetsByHWND[hwnd]
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return cw.wndProc(hwnd, msg, wParam, lParam, 0)
}

type PaintFunc func(surface *Surface, updateBounds Rectangle) os.Error

type CustomWidget struct {
	Widget
	paint               PaintFunc
	clearsBackground    bool
	invalidatesOnResize bool
}

var customWidgetsByHWND map[HWND]*CustomWidget

func NewCustomWidget(parent IContainer, style uint, paint PaintFunc) (*CustomWidget, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	if customWidgetsByHWND == nil {
		customWidgetsByHWND = make(map[HWND]*CustomWidget)
	}

	ensureRegisteredWindowClass(customWidgetWindowClass, customWidgetWndProc, &customWidgetWndProcPtr)

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr(customWidgetWindowClass), nil,
		WS_CHILD|WS_VISIBLE|style,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	cw := &CustomWidget{Widget: Widget{hWnd: hWnd, parent: parent}, paint: paint}

	cw.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = cw
	customWidgetsByHWND[hWnd] = cw

	parent.Children().Add(cw)

	return cw, nil
}

func (*CustomWidget) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz | ShrinkVert | GrowVert
}

func (cw *CustomWidget) PreferredSize() Size {
	return Size{100, 100}
}

func (cw *CustomWidget) ClearsBackground() bool {
	return cw.clearsBackground
}

func (cw *CustomWidget) SetClearsBackground(value bool) {
	cw.clearsBackground = value
}

func (cw *CustomWidget) InvalidatesOnResize() bool {
	return cw.invalidatesOnResize
}

func (cw *CustomWidget) SetInvalidatesOnResize(value bool) {
	cw.invalidatesOnResize = value
}

func (cw *CustomWidget) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		if cw.paint == nil {
			// TODO: log?
			break
		}

		var ps PAINTSTRUCT

		hdc := BeginPaint(cw.hWnd, &ps)
		if hdc == 0 {
			// TODO: log?
			break
		}
		defer EndPaint(cw.hWnd, &ps)

		surface, err := NewSurfaceFromHDC(hdc)
		if err != nil {
			// TODO: log?
			break
		}
		defer surface.Dispose()

		r := &ps.RcPaint
		err = cw.paint(surface, Rectangle{r.Left, r.Top, r.Right - r.Left, r.Bottom - r.Top})
		if err != nil {
			// TODO: log?
			break
		}

		return 0

	case WM_ERASEBKGND:
		if !cw.ClearsBackground() {
			return 1
		}

	case WM_SIZE, WM_SIZING:
		if cw.InvalidatesOnResize() {
			cw.Invalidate()
		}
	}

	return cw.Widget.wndProc(hwnd, msg, wParam, lParam, origWndProcPtr)
}
