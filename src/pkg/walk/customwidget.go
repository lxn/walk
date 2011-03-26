// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

const customWidgetWindowClass = `\o/ Walk_CustomWidget_Class \o/`

var customWidgetWindowClassRegistered bool

type PaintFunc func(canvas *Canvas, updateBounds Rectangle) os.Error

type CustomWidget struct {
	WidgetBase
	paint               PaintFunc
	clearsBackground    bool
	invalidatesOnResize bool
}

func NewCustomWidget(parent Container, style uint, paint PaintFunc) (*CustomWidget, os.Error) {
	ensureRegisteredWindowClass(customWidgetWindowClass, &customWidgetWindowClassRegistered)

	cw := &CustomWidget{paint: paint}

	if err := initChildWidget(
		cw,
		parent,
		customWidgetWindowClass,
		WS_VISIBLE|style,
		0); err != nil {
		return nil, err
	}

	return cw, nil
}

func (*CustomWidget) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (cw *CustomWidget) SizeHint() Size {
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

func (cw *CustomWidget) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		if cw.paint == nil {
			newError("paint func is nil")
			break
		}

		var ps PAINTSTRUCT

		hdc := BeginPaint(cw.hWnd, &ps)
		if hdc == 0 {
			newError("BeginPaint failed")
			break
		}
		defer EndPaint(cw.hWnd, &ps)

		canvas, err := newCanvasFromHDC(hdc)
		if err != nil {
			newError("newCanvasFromHDC failed")
			break
		}
		defer canvas.Dispose()

		r := &ps.RcPaint
		err = cw.paint(canvas, Rectangle{r.Left, r.Top, r.Right - r.Left, r.Bottom - r.Top})
		if err != nil {
			newError("paint failed")
			break
		}

		return 0

	case WM_ERASEBKGND:
		if !cw.clearsBackground {
			return 1
		}

	case WM_SIZE, WM_SIZING:
		if cw.invalidatesOnResize {
			cw.Invalidate()
		}
	}

	return cw.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
