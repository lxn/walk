// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
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

var gM = syscall.StringToUTF16Ptr("gM")

type Surface struct {
	hdc                 HDC
	hwnd                HWND
	dpix                int
	dpiy                int
	doNotDispose        bool
	recordingMetafile   *Metafile
	measureTextMetafile *Metafile
}

func NewSurfaceFromImage(image Image) (*Surface, os.Error) {
	switch img := image.(type) {
	case *Bitmap:
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

		if SelectObject(hdc, HGDIOBJ(img.hBmp)) == 0 {
			return nil, newError("SelectObject failed")
		}

		succeeded = true

		return (&Surface{hdc: hdc}).init()

	case *Metafile:
		surface, err := newSurfaceFromHDC(img.hdc)
		if err != nil {
			return nil, err
		}

		surface.recordingMetafile = img

		return surface, nil
	}

	return nil, newError("unsupported image type")
}

func newSurfaceFromHWND(hwnd HWND) (*Surface, os.Error) {
	hdc := GetDC(hwnd)
	if hdc == 0 {
		return nil, newError("GetDC failed")
	}

	return (&Surface{hdc: hdc, hwnd: hwnd}).init()
}

func newSurfaceFromHDC(hdc HDC) (*Surface, os.Error) {
	if hdc == 0 {
		return nil, newError("invalid hdc")
	}

	return (&Surface{hdc: hdc, doNotDispose: true}).init()
}

func (s *Surface) init() (*Surface, os.Error) {
	s.dpix = GetDeviceCaps(s.hdc, LOGPIXELSX)
	s.dpiy = GetDeviceCaps(s.hdc, LOGPIXELSY)

	if SetBkMode(s.hdc, TRANSPARENT) == 0 {
		return nil, newError("SetBkMode failed")
	}

	switch SetStretchBltMode(s.hdc, HALFTONE) {
	case 0, ERROR_INVALID_PARAMETER:
		return nil, newError("SetStretchBltMode failed")
	}

	if !SetBrushOrgEx(s.hdc, 0, 0, nil) {
		return nil, newError("SetBrushOrgEx failed")
	}

	return s, nil
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

	if s.recordingMetafile != nil {
		s.recordingMetafile.ensureFinished()
		s.recordingMetafile = nil
	}

	if s.measureTextMetafile != nil {
		s.measureTextMetafile.Dispose()
		s.measureTextMetafile = nil
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
	return s.withGdiObj(HGDIOBJ(font.HandleForDPI(s.dpiy)), func() os.Error {
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

func (s *Surface) Bounds() Rectangle {
	return Rectangle{
		Width:  GetDeviceCaps(s.hdc, HORZRES),
		Height: GetDeviceCaps(s.hdc, VERTRES),
	}
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
		ret := DrawTextEx(s.hdc, syscall.StringToUTF16Ptr(text), -1, &rect, uint(format)|DT_EDITCONTROL, nil)
		if ret == 0 {
			return newError("DrawTextEx failed")
		}

		return nil
	})
}

func (s *Surface) FontHeight(font *Font) (height int, err os.Error) {
	err = s.withFontAndTextColor(font, 0, func() os.Error {
		var size SIZE
		if !GetTextExtentPoint32(s.hdc, gM, 2, &size) {
			return newError("GetTextExtentPoint32 failed")
		}

		height = size.CY
		if height == 0 {
			return newError("invalid font height")
		}

		return nil
	})

	return
}

func (s *Surface) MeasureText(text string, font *Font, bounds Rectangle, format DrawTextFormat) (boundsMeasured Rectangle, runesFitted int, err os.Error) {
	// HACK: We don't want to actually draw on the surface here, but if we use
	// the DT_CALCRECT flag to avoid drawing, DRAWTEXTPARAMS.UiLengthDrawn will
	// not contain a useful value. To work around this, we create an in-memory
	// metafile and draw into that instead.
	if s.measureTextMetafile == nil {
		s.measureTextMetafile, err = NewMetafile(s)
		if err != nil {
			return
		}
	}

	hFont := HGDIOBJ(font.HandleForDPI(s.dpiy))
	oldHandle := SelectObject(s.measureTextMetafile.hdc, hFont)
	if oldHandle == 0 {
		err = newError("SelectObject failed")
		return
	}
	defer SelectObject(s.measureTextMetafile.hdc, oldHandle)

	rect := &RECT{bounds.X, bounds.Y, bounds.X + bounds.Width, bounds.Y + bounds.Height}
	var params DRAWTEXTPARAMS
	params.CbSize = uint(unsafe.Sizeof(params))

	strPtr := syscall.StringToUTF16Ptr(text)
	dtfmt := uint(format) | DT_EDITCONTROL | DT_WORDBREAK

	height := DrawTextEx(s.measureTextMetafile.hdc, strPtr, -1, rect, dtfmt, &params)
	if height == 0 {
		err = newError("DrawTextEx failed")
		return
	}

	boundsMeasured = Rectangle{rect.Left, rect.Top, rect.Right - rect.Left, height} //rect.Bottom - rect.Top}
	runesFitted = int(params.UiLengthDrawn)

	return
}
