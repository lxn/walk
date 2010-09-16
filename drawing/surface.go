// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drawing

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/gdi32"
	. "walk/winapi/user32"
)

// DrawText format flags
type DrawTextFormat uint

const (
	TextTop                  DrawTextFormat = DT_TOP
	TextLeft                 DrawTextFormat = DT_LEFT
	TextCenter               DrawTextFormat = DT_CENTER
	TextRight                DrawTextFormat = DT_RIGHT
	TextVCenter              DrawTextFormat = DT_VCENTER
	TextBottom               DrawTextFormat = DT_BOTTOM
	TextWordbreak            DrawTextFormat = DT_WORDBREAK
	TextSingleLine           DrawTextFormat = DT_SINGLELINE
	TextExpandTabs           DrawTextFormat = DT_EXPANDTABS
	TextTabstop              DrawTextFormat = DT_TABSTOP
	TextNoClip               DrawTextFormat = DT_NOCLIP
	TextExternalLeading      DrawTextFormat = DT_EXTERNALLEADING
	TextCalcRect             DrawTextFormat = DT_CALCRECT
	TextNoPrefix             DrawTextFormat = DT_NOPREFIX
	TextInternal             DrawTextFormat = DT_INTERNAL
	TextEditControl          DrawTextFormat = DT_EDITCONTROL
	TextPathEllipsis         DrawTextFormat = DT_PATH_ELLIPSIS
	TextEndEllipsis          DrawTextFormat = DT_END_ELLIPSIS
	TextModifyString         DrawTextFormat = DT_MODIFYSTRING
	TextRTLReading           DrawTextFormat = DT_RTLREADING
	TextWordEllipsis         DrawTextFormat = DT_WORD_ELLIPSIS
	TextNoFullWidthCharBreak DrawTextFormat = DT_NOFULLWIDTHCHARBREAK
	TextHidePrefix           DrawTextFormat = DT_HIDEPREFIX
	TextPrefixOnly           DrawTextFormat = DT_PREFIXONLY
)

type Surface struct {
	hdc  HDC
	hwnd HWND
}

func NewCompatibleBitmapSurface(size Size) (*Surface, os.Error) {
	return nil, newError("not implemented")
}

func NewDeviceSurface(driver, device string, devMode *DEVMODE) (*Surface, os.Error) {
	hdc := CreateDC(syscall.StringToUTF16Ptr(driver), syscall.StringToUTF16Ptr(device), nil, devMode)
	if hdc == 0 {
		return nil, newError("CreateDC failed")
	}

	if SetBkMode(hdc, TRANSPARENT) == 0 {
		return nil, newError("SetBkMode failed")
	}

	return &Surface{hdc: hdc}, nil
}

func NewWidgetSurface(hwnd HWND) (*Surface, os.Error) {
	hdc := GetDC(hwnd)
	if hdc == 0 {
		return nil, newError("GetDC failed")
	}

	if SetBkMode(hdc, TRANSPARENT) == 0 {
		return nil, newError("SetBkMode failed")
	}

	return &Surface{hdc: hdc, hwnd: hwnd}, nil
}

func (s *Surface) Dispose() {
	if s.hdc != 0 {
		if s.hwnd == 0 {
			DeleteDC(s.hdc)
		} else {
			ReleaseDC(s.hwnd, s.hdc)
		}

		s.hdc = 0
	}
}

func (s *Surface) withGdiObj(handle HGDIOBJ, f func() os.Error) os.Error {
	oldHandle := SelectObject(s.hdc, handle)
	if oldHandle == 0 {
		return newError("SelectObject failed")
	}
	defer SelectObject(s.hdc, oldHandle)

	return f()
}

func (s *Surface) withBrush(brush Brush, f func() os.Error) os.Error {
	return s.withGdiObj(HGDIOBJ(brush.handle()), f)
}

func (s *Surface) withFontAndTextColor(font *Font, color Color, f func() os.Error) os.Error {
	return s.withGdiObj(HGDIOBJ(font.hFont), func() os.Error {
		oldColor := SetTextColor(s.hdc, COLORREF(color))
		if oldColor == CLR_INVALID {
			return newError("SetTextColor failed")
		}
		defer func() {
			SetTextColor(s.hdc, oldColor)
		}()

		return f()
	})
}

func (s *Surface) withPen(pen Pen, f func() os.Error) os.Error {
	return s.withGdiObj(HGDIOBJ(pen.handle()), f)
}

func (s *Surface) withBrushAndPen(brush Brush, pen Pen, f func() os.Error) os.Error {
	return s.withBrush(brush, func() os.Error {
		return s.withPen(pen, f)
	})
}

func (s *Surface) ellipse(brush Brush, pen Pen, bounds *Rectangle) os.Error {
	return s.withBrushAndPen(brush, pen, func() os.Error {
		if !Ellipse(s.hdc, bounds.X, bounds.Y, bounds.X+bounds.Width, bounds.Y+bounds.Height) {
			return newError("Ellipse failed")
		}

		return nil
	})
}

func (s *Surface) DrawEllipse(pen Pen, bounds *Rectangle) os.Error {
	return s.ellipse(nullBrushSingleton, pen, bounds)
}

func (s *Surface) FillEllipse(brush Brush, bounds *Rectangle) os.Error {
	return s.ellipse(brush, nullPenSingleton, bounds)
}

func (s *Surface) DrawLine(pen Pen, from, to Point) os.Error {
	if !MoveToEx(s.hdc, from.X, from.Y, nil) {
		return newError("MoveToEx failed")
	}

	return s.withPen(pen, func() os.Error {
		if !LineTo(s.hdc, to.X, to.Y) {
			return newError("LineTo failed")
		}

		return nil
	})
}

func (s *Surface) rectangle(brush Brush, pen Pen, bounds *Rectangle) os.Error {
	return s.withBrushAndPen(brush, pen, func() os.Error {
		if !Rectangle_(s.hdc, bounds.X, bounds.Y, bounds.X+bounds.Width, bounds.Y+bounds.Height) {
			return newError("Rectangle_ failed")
		}

		return nil
	})
}

func (s *Surface) DrawRectangle(pen Pen, bounds *Rectangle) os.Error {
	return s.rectangle(nullBrushSingleton, pen, bounds)
}

func (s *Surface) FillRectangle(brush Brush, bounds *Rectangle) os.Error {
	return s.rectangle(brush, nullPenSingleton, bounds)
}

func (s *Surface) DrawText(text string, font *Font, color Color, bounds *Rectangle, format DrawTextFormat) os.Error {
	return s.withFontAndTextColor(font, color, func() os.Error {
		ret := DrawTextEx(s.hdc, syscall.StringToUTF16Ptr(text), -1, bounds.toRECT(), uint(format), nil)
		if ret == 0 {
			return newError("DrawTextEx failed")
		}

		return nil
	})
}
