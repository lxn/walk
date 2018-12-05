// Copyright 2018 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"unsafe"

	"github.com/lxn/win"
)

type static struct {
	WidgetBase
	textAlignment Alignment2D
	textColor     Color
}

func (s *static) init(widget Widget, parent Container) error {
	if err := InitWidget(
		widget,
		parent,
		"STATIC",
		win.WS_VISIBLE|win.SS_OWNERDRAW,
		0); err != nil {
		return err
	}

	s.SetBackground(nullBrushSingleton)

	return nil
}

func (*static) LayoutFlags() LayoutFlags {
	return GrowableHorz | GrowableVert
}

func (s *static) MinSizeHint() Size {
	return s.calculateTextSizeForWidth(0)
}

func (s *static) SizeHint() Size {
	return s.MinSizeHint()
}

func (s *static) HeightForWidth(width int) int {
	return s.MinSizeHint().Height
}

func (s *static) textAlignment1D() Alignment1D {
	switch s.textAlignment {
	case AlignHCenterVCenter:
		return AlignCenter

	case AlignHFarVCenter:
		return AlignFar

	default:
		return AlignNear
	}
}

func (s *static) setTextAlignment1D(alignment Alignment1D) error {
	var align Alignment2D

	switch alignment {
	case AlignCenter:
		align = AlignHCenterVCenter

	case AlignFar:
		align = AlignHFarVCenter

	default:
		align = AlignHNearVCenter
	}

	return s.setTextAlignment(align)
}

func (s *static) setTextAlignment(alignment Alignment2D) error {
	if alignment == s.textAlignment {
		return nil
	}

	s.textAlignment = alignment

	s.Invalidate()

	return nil
}

func (s *static) setText(value string) (changed bool, err error) {
	if value == s.text() {
		return false, nil
	}

	if err := s.WidgetBase.setText(value); err != nil {
		return false, err
	}

	return true, s.updateParentLayout()
}

func (s *static) TextColor() Color {
	return s.textColor
}

func (s *static) SetTextColor(c Color) {
	s.textColor = c

	s.Invalidate()
}

func (s *static) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_NCHITTEST:
		return win.HTCLIENT

	case win.WM_SIZE, win.WM_SIZING:
		s.Invalidate()

	case win.WM_DRAWITEM:
		dis := (*win.DRAWITEMSTRUCT)(unsafe.Pointer(lParam))

		canvas, err := newCanvasFromHDC(dis.HDC)
		if err != nil {
			break
		}
		canvas.Dispose()

		format := TextWordbreak

		switch s.textAlignment {
		case AlignHNearVNear, AlignHNearVCenter, AlignHNearVFar:
			format |= TextLeft

		case AlignHCenterVNear, AlignHCenterVCenter, AlignHCenterVFar:
			format |= TextCenter

		case AlignHFarVNear, AlignHFarVCenter, AlignHFarVFar:
			format |= TextRight
		}

		switch s.textAlignment {
		case AlignHNearVNear, AlignHCenterVNear, AlignHFarVNear:
			format |= TextTop

		case AlignHNearVCenter, AlignHCenterVCenter, AlignHFarVCenter:
			format |= TextVCenter

		case AlignHNearVFar, AlignHCenterVFar, AlignHFarVFar:
			format |= TextBottom
		}

		bounds := rectangleFromRECT(dis.RcItem)

		if format&TextVCenter != 0 || format&TextBottom != 0 {
			size := s.calculateTextSizeForWidth(bounds.Width)

			if format&TextVCenter != 0 {
				bounds.Y += (bounds.Height - size.Height) / 2
			} else {
				bounds.Y += bounds.Height - size.Height
			}

			bounds.Height = size.Height
		}

		bg, wnd := s.backgroundEffective()
		if bg == nil {
			bg = sysColorBtnFaceBrushSingleton
		}

		s.prepareDCForBackground(dis.HDC, s.hWnd, wnd)

		if err := canvas.FillRectangle(bg, s.ClientBounds()); err != nil {
			break
		}

		if err := canvas.DrawText(s.text(), s.Font(), s.textColor, bounds, format); err != nil {
			break
		}

		return 1
	}

	return s.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
