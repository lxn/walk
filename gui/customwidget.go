// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	"walk/drawing"
	. "walk/winapi/user32"
)

const customWidgetWindowClass = `\o/ Walk_CustomWidget_Class \o/`

var customWidgetWndProcCallback *syscall.Callback

func customWidgetWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	cw, ok := customWidgetsByHWND[msg.HWnd]
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return cw.wndProc(msg)
}

type PaintFunc func(surface *drawing.Surface, updateBounds drawing.Rectangle) os.Error

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

	ensureRegisteredWindowClass(customWidgetWindowClass, customWidgetWndProc, &customWidgetWndProcCallback)

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

func (cw *CustomWidget) PreferredSize() drawing.Size {
	return drawing.Size{100, 100}
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

func (cw *CustomWidget) wndProc(msg *MSG) uintptr {
	switch msg.Message {
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

		surface, err := drawing.NewSurfaceFromHDC(hdc)
		if err != nil {
			// TODO: log?
			break
		}
		defer surface.Dispose()

		r := &ps.RcPaint
		err = cw.paint(surface, drawing.Rectangle{r.Left, r.Top, r.Right - r.Left, r.Bottom - r.Top})
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

	return cw.Widget.wndProc(msg)
}
