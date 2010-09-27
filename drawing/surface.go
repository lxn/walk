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
	. "walk/winapi/kernel32"
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
	hdc          HDC
	hwnd         HWND
	doNotDispose bool
}

func initHDC(hdc HDC) os.Error {
	if SetBkMode(hdc, TRANSPARENT) == 0 {
		return newError("SetBkMode failed")
	}

	switch SetStretchBltMode(hdc, HALFTONE) {
	case 0, ERROR_INVALID_PARAMETER:
		return newError("SetStretchBltMode failed")
	}

	if !SetBrushOrgEx(hdc, 0, 0, nil) {
		return newError("SetBrushOrgEx failed")
	}

	return nil
}

func NewSurfaceFromBitmap(bmp *Bitmap) (*Surface, os.Error) {
	hdc := CreateCompatibleDC(0)
	if hdc == 0 {
		return nil, newError("CreateCompatibleDC failed")
	}
	succeeded := false

	defer func() {
		if !succeeded {
			DeleteDC(hdc)
		}
	}()

	if SelectObject(hdc, HGDIOBJ(bmp.hBmp)) == 0 {
		return nil, newError("SelectObject failed")
	}

	if err := initHDC(hdc); err != nil {
		return nil, err
	}

	succeeded = true

	return &Surface{hdc: hdc}, nil
}

func NewSurfaceFromDevice(driver, device string, devMode *DEVMODE) (*Surface, os.Error) {
	hdc := CreateDC(syscall.StringToUTF16Ptr(driver), syscall.StringToUTF16Ptr(device), nil, devMode)
	if hdc == 0 {
		return nil, newError("CreateDC failed")
	}

	if err := initHDC(hdc); err != nil {
		return nil, err
	}

	return &Surface{hdc: hdc}, nil
}

func NewSurfaceFromWidget(hwnd HWND) (*Surface, os.Error) {
	hdc := GetDC(hwnd)
	if hdc == 0 {
		return nil, newError("GetDC failed")
	}

	if err := initHDC(hdc); err != nil {
		return nil, err
	}

	return &Surface{hdc: hdc, hwnd: hwnd}, nil
}

func NewSurfaceFromHDC(hdc HDC) (*Surface, os.Error) {
	if hdc == 0 {
		return nil, newError("invalid hdc")
	}

	if err := initHDC(hdc); err != nil {
		return nil, err
	}

	return &Surface{hdc: hdc, doNotDispose: true}, nil
}

func (s *Surface) Dispose() {
	if !s.doNotDispose && s.hdc != 0 {
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

func (s *Surface) ellipse(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) os.Error {
	return s.withBrushAndPen(brush, pen, func() os.Error {
		if !Ellipse(s.hdc, bounds.X, bounds.Y, bounds.X+bounds.Width+sizeCorrection, bounds.Y+bounds.Height+sizeCorrection) {
			return newError("Ellipse failed")
		}

		return nil
	})
}

func (s *Surface) DrawEllipse(pen Pen, bounds Rectangle) os.Error {
	return s.ellipse(nullBrushSingleton, pen, bounds, 0)
}

func (s *Surface) FillEllipse(brush Brush, bounds Rectangle) os.Error {
	return s.ellipse(brush, nullPenSingleton, bounds, 1)
}

func (s *Surface) DrawImage(image Image, location Point) os.Error {
	if image == nil {
		return newError("image cannot be nil")
	}

	return image.draw(s.hdc, location)
}

func (s *Surface) DrawImageStretched(image Image, bounds Rectangle) os.Error {
	if image == nil {
		return newError("image cannot be nil")
	}

	return image.drawStretched(s.hdc, bounds)
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

func (s *Surface) rectangle(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) os.Error {
	return s.withBrushAndPen(brush, pen, func() os.Error {
		if !Rectangle_(s.hdc, bounds.X, bounds.Y, bounds.X+bounds.Width+sizeCorrection, bounds.Y+bounds.Height+sizeCorrection) {
			return newError("Rectangle_ failed")
		}

		return nil
	})
}

func (s *Surface) DrawRectangle(pen Pen, bounds Rectangle) os.Error {
	return s.rectangle(nullBrushSingleton, pen, bounds, 0)
}

func (s *Surface) FillRectangle(brush Brush, bounds Rectangle) os.Error {
	return s.rectangle(brush, nullPenSingleton, bounds, 1)
}

func (s *Surface) DrawText(text string, font *Font, color Color, bounds Rectangle, format DrawTextFormat) os.Error {
	return s.withFontAndTextColor(font, color, func() os.Error {
		rect := bounds.toRECT()
		ret := DrawTextEx(s.hdc, syscall.StringToUTF16Ptr(text), -1, &rect, uint(format), nil)
		if ret == 0 {
			return newError("DrawTextEx failed")
		}

		return nil
	})
}
