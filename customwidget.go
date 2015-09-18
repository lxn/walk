// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

const customWidgetWindowClass = `\o/ Walk_CustomWidget_Class \o/`

func init() {
	MustRegisterWindowClass(customWidgetWindowClass)
}

type PaintFunc func(canvas *Canvas, updateBounds Rectangle) error

type CustomWidget struct {
	WidgetBase
	paint               PaintFunc
	clearsBackground    bool
	invalidatesOnResize bool
}

func NewCustomWidget(parent Container, style uint, paint PaintFunc) (*CustomWidget, error) {
	cw := &CustomWidget{paint: paint}

	if err := InitWidget(
		cw,
		parent,
		customWidgetWindowClass,
		win.WS_VISIBLE|uint32(style),
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

func (cw *CustomWidget) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_PAINT:
		if cw.paint == nil {
			newError("paint func is nil")
			break
		}

		var ps win.PAINTSTRUCT

		hdc := win.BeginPaint(cw.hWnd, &ps)
		if hdc == 0 {
			newError("BeginPaint failed")
			break
		}
		defer win.EndPaint(cw.hWnd, &ps)

		canvas, err := newCanvasFromHDC(hdc)
		if err != nil {
			newError("newCanvasFromHDC failed")
			break
		}
		defer canvas.Dispose()

		r := &ps.RcPaint
		err = cw.paint(
			canvas,
			Rectangle{
				int(r.Left),
				int(r.Top),
				int(r.Right - r.Left),
				int(r.Bottom - r.Top),
			})
		if err != nil {
			newError("paint failed")
			break
		}

		return 0

	case win.WM_ERASEBKGND:
		if !cw.clearsBackground {
			return 1
		}

	case win.WM_SIZE, win.WM_SIZING:
		if cw.invalidatesOnResize {
			cw.Invalidate()
		}
	}

	return cw.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
