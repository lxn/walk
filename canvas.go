// Copyright 2010 The Walk Authorc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import . "walk/winapi"

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

type Canvas struct {
	hdc                 HDC
	hwnd                HWND
	dpix                int
	dpiy                int
	doNotDispose        bool
	recordingMetafile   *Metafile
	measureTextMetafile *Metafile
}

func NewCanvasFromImage(image Image) (*Canvas, os.Error) {
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

		return (&Canvas{hdc: hdc}).init()

	case *Metafile:
		c, err := newCanvasFromHDC(img.hdc)
		if err != nil {
			return nil, err
		}

		c.recordingMetafile = img

		return c, nil
	}

	return nil, newError("unsupported image type")
}

func newCanvasFromHWND(hwnd HWND) (*Canvas, os.Error) {
	hdc := GetDC(hwnd)
	if hdc == 0 {
		return nil, newError("GetDC failed")
	}

	return (&Canvas{hdc: hdc, hwnd: hwnd}).init()
}

func newCanvasFromHDC(hdc HDC) (*Canvas, os.Error) {
	if hdc == 0 {
		return nil, newError("invalid hdc")
	}

	return (&Canvas{hdc: hdc, doNotDispose: true}).init()
}

func (c *Canvas) init() (*Canvas, os.Error) {
	c.dpix = int(GetDeviceCaps(c.hdc, LOGPIXELSX))
	c.dpiy = int(GetDeviceCaps(c.hdc, LOGPIXELSY))

	if SetBkMode(c.hdc, TRANSPARENT) == 0 {
		return nil, newError("SetBkMode failed")
	}

	switch SetStretchBltMode(c.hdc, HALFTONE) {
	case 0, ERROR_INVALID_PARAMETER:
		return nil, newError("SetStretchBltMode failed")
	}

	if !SetBrushOrgEx(c.hdc, 0, 0, nil) {
		return nil, newError("SetBrushOrgEx failed")
	}

	return c, nil
}

func (c *Canvas) Dispose() {
	if !c.doNotDispose && c.hdc != 0 {
		if c.hwnd == 0 {
			DeleteDC(c.hdc)
		} else {
			ReleaseDC(c.hwnd, c.hdc)
		}

		c.hdc = 0
	}

	if c.recordingMetafile != nil {
		c.recordingMetafile.ensureFinished()
		c.recordingMetafile = nil
	}

	if c.measureTextMetafile != nil {
		c.measureTextMetafile.Dispose()
		c.measureTextMetafile = nil
	}
}

func (c *Canvas) withGdiObj(handle HGDIOBJ, f func() os.Error) os.Error {
	oldHandle := SelectObject(c.hdc, handle)
	if oldHandle == 0 {
		return newError("SelectObject failed")
	}
	defer SelectObject(c.hdc, oldHandle)

	return f()
}

func (c *Canvas) withBrush(brush Brush, f func() os.Error) os.Error {
	return c.withGdiObj(HGDIOBJ(brush.handle()), f)
}

func (c *Canvas) withFontAndTextColor(font *Font, color Color, f func() os.Error) os.Error {
	return c.withGdiObj(HGDIOBJ(font.handleForDPI(c.dpiy)), func() os.Error {
		oldColor := SetTextColor(c.hdc, COLORREF(color))
		if oldColor == CLR_INVALID {
			return newError("SetTextColor failed")
		}
		defer func() {
			SetTextColor(c.hdc, oldColor)
		}()

		return f()
	})
}

func (c *Canvas) Bounds() Rectangle {
	return Rectangle{
		Width:  int(GetDeviceCaps(c.hdc, HORZRES)),
		Height: int(GetDeviceCaps(c.hdc, VERTRES)),
	}
}

func (c *Canvas) withPen(pen Pen, f func() os.Error) os.Error {
	return c.withGdiObj(HGDIOBJ(pen.handle()), f)
}

func (c *Canvas) withBrushAndPen(brush Brush, pen Pen, f func() os.Error) os.Error {
	return c.withBrush(brush, func() os.Error {
		return c.withPen(pen, f)
	})
}

func (c *Canvas) ellipse(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) os.Error {
	return c.withBrushAndPen(brush, pen, func() os.Error {
		if !Ellipse(
			c.hdc,
			int32(bounds.X),
			int32(bounds.Y),
			int32(bounds.X+bounds.Width+sizeCorrection),
			int32(bounds.Y+bounds.Height+sizeCorrection)) {

			return newError("Ellipse failed")
		}

		return nil
	})
}

func (c *Canvas) DrawEllipse(pen Pen, bounds Rectangle) os.Error {
	return c.ellipse(nullBrushSingleton, pen, bounds, 0)
}

func (c *Canvas) FillEllipse(brush Brush, bounds Rectangle) os.Error {
	return c.ellipse(brush, nullPenSingleton, bounds, 1)
}

func (c *Canvas) DrawImage(image Image, location Point) os.Error {
	if image == nil {
		return newError("image cannot be nil")
	}

	return image.draw(c.hdc, location)
}

func (c *Canvas) DrawImageStretched(image Image, bounds Rectangle) os.Error {
	if image == nil {
		return newError("image cannot be nil")
	}

	return image.drawStretched(c.hdc, bounds)
}

func (c *Canvas) DrawLine(pen Pen, from, to Point) os.Error {
	if !MoveToEx(c.hdc, from.X, from.Y, nil) {
		return newError("MoveToEx failed")
	}

	return c.withPen(pen, func() os.Error {
		if !LineTo(c.hdc, int32(to.X), int32(to.Y)) {
			return newError("LineTo failed")
		}

		return nil
	})
}

func (c *Canvas) rectangle(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) os.Error {
	return c.withBrushAndPen(brush, pen, func() os.Error {
		if !Rectangle_(
			c.hdc,
			int32(bounds.X),
			int32(bounds.Y),
			int32(bounds.X+bounds.Width+sizeCorrection),
			int32(bounds.Y+bounds.Height+sizeCorrection)) {

			return newError("Rectangle_ failed")
		}

		return nil
	})
}

func (c *Canvas) DrawRectangle(pen Pen, bounds Rectangle) os.Error {
	return c.rectangle(nullBrushSingleton, pen, bounds, 0)
}

func (c *Canvas) FillRectangle(brush Brush, bounds Rectangle) os.Error {
	return c.rectangle(brush, nullPenSingleton, bounds, 1)
}

func (c *Canvas) DrawText(text string, font *Font, color Color, bounds Rectangle, format DrawTextFormat) os.Error {
	return c.withFontAndTextColor(font, color, func() os.Error {
		rect := bounds.toRECT()
		ret := DrawTextEx(
			c.hdc,
			syscall.StringToUTF16Ptr(text),
			-1,
			&rect,
			uint32(format)|DT_EDITCONTROL,
			nil)
		if ret == 0 {
			return newError("DrawTextEx failed")
		}

		return nil
	})
}

func (c *Canvas) fontHeight(font *Font) (height int, err os.Error) {
	err = c.withFontAndTextColor(font, 0, func() os.Error {
		var size SIZE
		if !GetTextExtentPoint32(c.hdc, gM, 2, &size) {
			return newError("GetTextExtentPoint32 failed")
		}

		height = int(size.CY)
		if height == 0 {
			return newError("invalid font height")
		}

		return nil
	})

	return
}

func (c *Canvas) MeasureText(text string, font *Font, bounds Rectangle, format DrawTextFormat) (boundsMeasured Rectangle, runesFitted int, err os.Error) {
	// HACK: We don't want to actually draw on the Canvas here, but if we use
	// the DT_CALCRECT flag to avoid drawing, DRAWTEXTPARAMc.UiLengthDrawn will
	// not contain a useful value. To work around this, we create an in-memory
	// metafile and draw into that instead.
	if c.measureTextMetafile == nil {
		c.measureTextMetafile, err = NewMetafile(c)
		if err != nil {
			return
		}
	}

	hFont := HGDIOBJ(font.handleForDPI(c.dpiy))
	oldHandle := SelectObject(c.measureTextMetafile.hdc, hFont)
	if oldHandle == 0 {
		err = newError("SelectObject failed")
		return
	}
	defer SelectObject(c.measureTextMetafile.hdc, oldHandle)

	rect := &RECT{
		int32(bounds.X),
		int32(bounds.Y),
		int32(bounds.X + bounds.Width),
		int32(bounds.Y + bounds.Height),
	}
	var params DRAWTEXTPARAMS
	params.CbSize = uint32(unsafe.Sizeof(params))

	strPtr := syscall.StringToUTF16Ptr(text)
	dtfmt := uint32(format) | DT_EDITCONTROL | DT_WORDBREAK

	height := DrawTextEx(
		c.measureTextMetafile.hdc, strPtr, -1, rect, dtfmt, &params)
	if height == 0 {
		err = newError("DrawTextEx failed")
		return
	}

	boundsMeasured = Rectangle{
		int(rect.Left),
		int(rect.Top),
		int(rect.Right - rect.Left),
		int(height),
	}
	runesFitted = int(params.UiLengthDrawn)

	return
}
