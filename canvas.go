// Copyright 2010 The Walk Authorc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/lxn/win"
)

// DrawText format flags
type DrawTextFormat uint

const (
	TextTop                  DrawTextFormat = win.DT_TOP
	TextLeft                 DrawTextFormat = win.DT_LEFT
	TextCenter               DrawTextFormat = win.DT_CENTER
	TextRight                DrawTextFormat = win.DT_RIGHT
	TextVCenter              DrawTextFormat = win.DT_VCENTER
	TextBottom               DrawTextFormat = win.DT_BOTTOM
	TextWordbreak            DrawTextFormat = win.DT_WORDBREAK
	TextSingleLine           DrawTextFormat = win.DT_SINGLELINE
	TextExpandTabs           DrawTextFormat = win.DT_EXPANDTABS
	TextTabstop              DrawTextFormat = win.DT_TABSTOP
	TextNoClip               DrawTextFormat = win.DT_NOCLIP
	TextExternalLeading      DrawTextFormat = win.DT_EXTERNALLEADING
	TextCalcRect             DrawTextFormat = win.DT_CALCRECT
	TextNoPrefix             DrawTextFormat = win.DT_NOPREFIX
	TextInternal             DrawTextFormat = win.DT_INTERNAL
	TextEditControl          DrawTextFormat = win.DT_EDITCONTROL
	TextPathEllipsis         DrawTextFormat = win.DT_PATH_ELLIPSIS
	TextEndEllipsis          DrawTextFormat = win.DT_END_ELLIPSIS
	TextModifyString         DrawTextFormat = win.DT_MODIFYSTRING
	TextRTLReading           DrawTextFormat = win.DT_RTLREADING
	TextWordEllipsis         DrawTextFormat = win.DT_WORD_ELLIPSIS
	TextNoFullWidthCharBreak DrawTextFormat = win.DT_NOFULLWIDTHCHARBREAK
	TextHidePrefix           DrawTextFormat = win.DT_HIDEPREFIX
	TextPrefixOnly           DrawTextFormat = win.DT_PREFIXONLY
)

var gM *uint16

func init() {
	AppendToWalkInit(func() {
		gM = syscall.StringToUTF16Ptr("gM")
	})
}

type Canvas struct {
	hdc                 win.HDC
	hBmpStock           win.HBITMAP
	window              Window
	dpi                 int
	bitmap              *Bitmap
	recordingMetafile   *Metafile
	measureTextMetafile *Metafile
	doNotDispose        bool
}

func NewCanvasFromImage(image Image) (*Canvas, error) {
	switch img := image.(type) {
	case *Bitmap:
		hdc := win.CreateCompatibleDC(0)
		if hdc == 0 {
			return nil, newError("CreateCompatibleDC failed")
		}
		succeeded := false

		defer func() {
			if !succeeded {
				win.DeleteDC(hdc)
			}
		}()

		var hBmpStock win.HBITMAP
		if hBmpStock = win.HBITMAP(win.SelectObject(hdc, win.HGDIOBJ(img.hBmp))); hBmpStock == 0 {
			return nil, newError("SelectObject failed")
		}

		succeeded = true

		return (&Canvas{hdc: hdc, hBmpStock: hBmpStock, bitmap: img, dpi: img.dpi}).init()

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

func newCanvasFromWindow(window Window) (*Canvas, error) {
	hdc := win.GetDC(window.Handle())
	if hdc == 0 {
		return nil, newError("GetDC failed")
	}

	return (&Canvas{hdc: hdc, window: window}).init()
}

func newCanvasFromHDC(hdc win.HDC) (*Canvas, error) {
	if hdc == 0 {
		return nil, newError("invalid hdc")
	}

	return (&Canvas{hdc: hdc, doNotDispose: true}).init()
}

func (c *Canvas) init() (*Canvas, error) {
	if c.dpi == 0 {
		c.dpi = dpiForHDC(c.hdc)
	}

	if win.SetBkMode(c.hdc, win.TRANSPARENT) == 0 {
		return nil, newError("SetBkMode failed")
	}

	switch win.SetStretchBltMode(c.hdc, win.HALFTONE) {
	case 0, win.ERROR_INVALID_PARAMETER:
		return nil, newError("SetStretchBltMode failed")
	}

	if !win.SetBrushOrgEx(c.hdc, 0, 0, nil) {
		return nil, newError("SetBrushOrgEx failed")
	}

	return c, nil
}

func (c *Canvas) Dispose() {
	if !c.doNotDispose && c.hdc != 0 {
		if c.bitmap != nil {
			win.SelectObject(c.hdc, win.HGDIOBJ(c.hBmpStock))
			win.DeleteDC(c.hdc)
			c.bitmap.postProcess()
		} else {
			win.ReleaseDC(c.window.Handle(), c.hdc)
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

func (c *Canvas) DPI() int {
	if c.window != nil {
		return c.window.DPI()
	}

	return c.dpi
}

func (c *Canvas) withGdiObj(handle win.HGDIOBJ, f func() error) error {
	oldHandle := win.SelectObject(c.hdc, handle)
	if oldHandle == 0 {
		return newError("SelectObject failed")
	}
	defer win.SelectObject(c.hdc, oldHandle)

	return f()
}

func (c *Canvas) withBrush(brush Brush, f func() error) error {
	return c.withGdiObj(win.HGDIOBJ(brush.handle()), f)
}

func (c *Canvas) withFontAndTextColor(font *Font, color Color, f func() error) error {
	return c.withGdiObj(win.HGDIOBJ(font.handleForDPI(c.DPI())), func() error {
		oldColor := win.SetTextColor(c.hdc, win.COLORREF(color))
		if oldColor == win.CLR_INVALID {
			return newError("SetTextColor failed")
		}
		defer func() {
			win.SetTextColor(c.hdc, oldColor)
		}()

		return f()
	})
}

func (c *Canvas) HDC() win.HDC {
	return c.hdc
}

func (c *Canvas) Bounds() Rectangle {
	return RectangleTo96DPI(c.BoundsPixels(), c.DPI())
}

func (c *Canvas) BoundsPixels() Rectangle {
	return Rectangle{
		Width:  int(win.GetDeviceCaps(c.hdc, win.HORZRES)),
		Height: int(win.GetDeviceCaps(c.hdc, win.VERTRES)),
	}
}

func (c *Canvas) withPen(pen Pen, f func() error) error {
	return c.withGdiObj(win.HGDIOBJ(pen.handleForDPI(c.dpi)), f)
}

func (c *Canvas) withBrushAndPen(brush Brush, pen Pen, f func() error) error {
	return c.withBrush(brush, func() error {
		return c.withPen(pen, f)
	})
}

// ellipse draws an ellipse in 1/96" units. sizeCorrection parameter is in native pixels.
//
// Deprecated: Newer applications should use ellipsePixels.
func (c *Canvas) ellipse(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) error {
	return c.ellipsePixels(brush, pen, RectangleFrom96DPI(bounds, c.DPI()), sizeCorrection)
}

// ellipsePixels draws an ellipse in native pixels.
func (c *Canvas) ellipsePixels(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) error {
	return c.withBrushAndPen(brush, pen, func() error {
		if !win.Ellipse(
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

// DrawEllipse draws an ellipse in 1/96" units.
//
// Deprecated: Newer applications should use DrawEllipsePixels.
func (c *Canvas) DrawEllipse(pen Pen, bounds Rectangle) error {
	return c.ellipse(nullBrushSingleton, pen, bounds, 0)
}

// DrawEllipsePixels draws an ellipse in native pixels.
func (c *Canvas) DrawEllipsePixels(pen Pen, bounds Rectangle) error {
	return c.ellipsePixels(nullBrushSingleton, pen, bounds, 0)
}

// FillEllipse draws a filled ellipse in 1/96" units.
//
// Deprecated: Newer applications should use FillEllipsePixels.
func (c *Canvas) FillEllipse(brush Brush, bounds Rectangle) error {
	return c.ellipse(brush, nullPenSingleton, bounds, 1)
}

// FillEllipsePixels draws a filled in native pixels.
func (c *Canvas) FillEllipsePixels(brush Brush, bounds Rectangle) error {
	return c.ellipsePixels(brush, nullPenSingleton, bounds, 1)
}

// DrawImage draws image at given location (upper left) in 1/96" units unstretched.
//
// Deprecated: Newer applications should use DrawImagePixels.
func (c *Canvas) DrawImage(image Image, location Point) error {
	return c.DrawImagePixels(image, PointFrom96DPI(location, c.DPI()))
}

// DrawImagePixels draws image at given location (upper left) in native pixels unstretched.
func (c *Canvas) DrawImagePixels(image Image, location Point) error {
	if image == nil {
		return newError("image cannot be nil")
	}

	return image.draw(c.hdc, location)
}

// DrawImageStretched draws image at given location in 1/96" units stretched.
//
// Deprecated: Newer applications should use DrawImageStretchedPixels.
func (c *Canvas) DrawImageStretched(image Image, bounds Rectangle) error {
	return c.DrawImageStretchedPixels(image, RectangleFrom96DPI(bounds, c.DPI()))
}

// DrawImageStretchedPixels draws image at given location in native pixels stretched.
func (c *Canvas) DrawImageStretchedPixels(image Image, bounds Rectangle) error {
	if image == nil {
		return newError("image cannot be nil")
	}

	if dsoc, ok := image.(interface {
		drawStretchedOnCanvasPixels(canvas *Canvas, bounds Rectangle) error
	}); ok {
		return dsoc.drawStretchedOnCanvasPixels(c, bounds)
	}

	return image.drawStretched(c.hdc, bounds)
}

// DrawBitmapWithOpacity draws bitmap with opacity at given location in 1/96" units stretched.
//
// Deprecated: Newer applications should use DrawBitmapWithOpacityPixels.
func (c *Canvas) DrawBitmapWithOpacity(bmp *Bitmap, bounds Rectangle, opacity byte) error {
	return c.DrawBitmapWithOpacityPixels(bmp, RectangleFrom96DPI(bounds, c.DPI()), opacity)
}

// DrawBitmapWithOpacityPixels draws bitmap with opacity at given location in native pixels
// stretched.
func (c *Canvas) DrawBitmapWithOpacityPixels(bmp *Bitmap, bounds Rectangle, opacity byte) error {
	if bmp == nil {
		return newError("bmp cannot be nil")
	}

	return bmp.alphaBlend(c.hdc, bounds, opacity)
}

// DrawBitmapPart draws bitmap at given location in native pixels.
func (c *Canvas) DrawBitmapPart(bmp *Bitmap, dst, src Rectangle) error {
	return c.DrawBitmapPartWithOpacityPixels(bmp, dst, src, 0xff)
}

// DrawBitmapPartWithOpacity draws bitmap at given location in 1/96" units.
//
// Deprecated: Newer applications should use DrawBitmapPartWithOpacityPixels.
func (c *Canvas) DrawBitmapPartWithOpacity(bmp *Bitmap, dst, src Rectangle, opacity byte) error {
	dpi := c.DPI()
	return c.DrawBitmapPartWithOpacityPixels(bmp, RectangleFrom96DPI(dst, dpi), RectangleFrom96DPI(src, dpi), opacity)
}

// DrawBitmapPartWithOpacityPixels draws bitmap at given location in native pixels.
func (c *Canvas) DrawBitmapPartWithOpacityPixels(bmp *Bitmap, dst, src Rectangle, opacity byte) error {
	if bmp == nil {
		return newError("bmp cannot be nil")
	}

	return bmp.alphaBlendPart(c.hdc, dst, src, opacity)
}

// DrawLine draws a line between two points in 1/96" units.
//
// Deprecated: Newer applications should use DrawLinePixels.
func (c *Canvas) DrawLine(pen Pen, from, to Point) error {
	dpi := c.DPI()
	return c.DrawLinePixels(pen, PointFrom96DPI(from, dpi), PointFrom96DPI(to, dpi))
}

// DrawLinePixels draws a line between two points in native pixels.
func (c *Canvas) DrawLinePixels(pen Pen, from, to Point) error {
	if !win.MoveToEx(c.hdc, int(from.X), int(from.Y), nil) {
		return newError("MoveToEx failed")
	}

	return c.withPen(pen, func() error {
		if !win.LineTo(c.hdc, int32(to.X), int32(to.Y)) {
			return newError("LineTo failed")
		}

		return nil
	})
}

// DrawLine draws a line between given points in 1/96" units.
//
// Deprecated: Newer applications should use DrawLinePixels.
func (c *Canvas) DrawPolyline(pen Pen, points []Point) error {
	if len(points) < 1 {
		return nil
	}

	dpi := c.DPI()

	pts := make([]win.POINT, len(points))
	for i, p := range points {
		pts[i] = PointFrom96DPI(p, dpi).toPOINT()
	}

	return c.withPen(pen, func() error {
		if !win.Polyline(c.hdc, unsafe.Pointer(&pts[0].X), int32(len(pts))) {
			return newError("Polyline failed")
		}

		return nil
	})
}

// DrawPolylinePixels draws a line between given points in native pixels.
func (c *Canvas) DrawPolylinePixels(pen Pen, points []Point) error {
	if len(points) < 1 {
		return nil
	}

	pts := make([]win.POINT, len(points))
	for i, p := range points {
		pts[i] = p.toPOINT()
	}

	return c.withPen(pen, func() error {
		if !win.Polyline(c.hdc, unsafe.Pointer(&pts[0].X), int32(len(pts))) {
			return newError("Polyline failed")
		}

		return nil
	})
}

// rectangle draws a rectangle in 1/96" units. sizeCorrection parameter is in native pixels.
//
// Deprecated: Newer applications should use rectanglePixels.
func (c *Canvas) rectangle(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) error {
	return c.rectanglePixels(brush, pen, RectangleFrom96DPI(bounds, c.DPI()), sizeCorrection)
}

// rectanglePixels draws a rectangle in native pixels.
func (c *Canvas) rectanglePixels(brush Brush, pen Pen, bounds Rectangle, sizeCorrection int) error {
	return c.withBrushAndPen(brush, pen, func() error {
		if !win.Rectangle_(
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

// DrawRectangle draws a rectangle in 1/96" units.
//
// Deprecated: Newer applications should use DrawRectanglePixels.
func (c *Canvas) DrawRectangle(pen Pen, bounds Rectangle) error {
	return c.rectangle(nullBrushSingleton, pen, bounds, 0)
}

// DrawRectanglePixels draws a rectangle in native pixels.
func (c *Canvas) DrawRectanglePixels(pen Pen, bounds Rectangle) error {
	return c.rectanglePixels(nullBrushSingleton, pen, bounds, 0)
}

// FillRectangle draws a filled rectangle in 1/96" units.
//
// Deprecated: Newer applications should use FillRectanglePixels.
func (c *Canvas) FillRectangle(brush Brush, bounds Rectangle) error {
	return c.rectangle(brush, nullPenSingleton, bounds, 1)
}

// FillRectanglePixels draws a filled rectangle in native pixels.
func (c *Canvas) FillRectanglePixels(brush Brush, bounds Rectangle) error {
	return c.rectanglePixels(brush, nullPenSingleton, bounds, 1)
}

// roundedRectangle draws a rounded rectangle in 1/96" units. sizeCorrection parameter is in native
// pixels.
//
// Deprecated: Newer applications should use roundedRectanglePixels.
func (c *Canvas) roundedRectangle(brush Brush, pen Pen, bounds Rectangle, ellipseSize Size, sizeCorrection int) error {
	dpi := c.DPI()
	return c.roundedRectanglePixels(brush, pen, RectangleFrom96DPI(bounds, dpi), SizeFrom96DPI(ellipseSize, dpi), sizeCorrection)
}

// roundedRectanglePixels draws a rounded rectangle in native pixels.
func (c *Canvas) roundedRectanglePixels(brush Brush, pen Pen, bounds Rectangle, ellipseSize Size, sizeCorrection int) error {
	return c.withBrushAndPen(brush, pen, func() error {
		if !win.RoundRect(
			c.hdc,
			int32(bounds.X),
			int32(bounds.Y),
			int32(bounds.X+bounds.Width+sizeCorrection),
			int32(bounds.Y+bounds.Height+sizeCorrection),
			int32(ellipseSize.Width),
			int32(ellipseSize.Height)) {

			return newError("RoundRect failed")
		}

		return nil
	})
}

// DrawRoundedRectangle draws a rounded rectangle in 1/96" units. sizeCorrection parameter is in native
// pixels.
//
// Deprecated: Newer applications should use DrawRoundedRectanglePixels.
func (c *Canvas) DrawRoundedRectangle(pen Pen, bounds Rectangle, ellipseSize Size) error {
	return c.roundedRectangle(nullBrushSingleton, pen, bounds, ellipseSize, 0)
}

// DrawRoundedRectanglePixels draws a rounded rectangle in native pixels.
func (c *Canvas) DrawRoundedRectanglePixels(pen Pen, bounds Rectangle, ellipseSize Size) error {
	return c.roundedRectanglePixels(nullBrushSingleton, pen, bounds, ellipseSize, 0)
}

// FillRoundedRectangle draws a filled rounded rectangle in 1/96" units. sizeCorrection parameter
// is in native
// pixels.
//
// Deprecated: Newer applications should use FillRoundedRectanglePixels.
func (c *Canvas) FillRoundedRectangle(brush Brush, bounds Rectangle, ellipseSize Size) error {
	return c.roundedRectangle(brush, nullPenSingleton, bounds, ellipseSize, 1)
}

// FillRoundedRectanglePixels draws a filled rounded rectangle in native pixels.
func (c *Canvas) FillRoundedRectanglePixels(brush Brush, bounds Rectangle, ellipseSize Size) error {
	return c.roundedRectanglePixels(brush, nullPenSingleton, bounds, ellipseSize, 1)
}

// GradientFillRectangle draws a gradient filled rectangle in 1/96" units.
//
// Deprecated: Newer applications should use GradientFillRectanglePixels.
func (c *Canvas) GradientFillRectangle(color1, color2 Color, orientation Orientation, bounds Rectangle) error {
	return c.GradientFillRectanglePixels(color1, color2, orientation, RectangleFrom96DPI(bounds, c.DPI()))
}

// GradientFillRectanglePixels draws a gradient filled rectangle in native pixels.
func (c *Canvas) GradientFillRectanglePixels(color1, color2 Color, orientation Orientation, bounds Rectangle) error {
	vertices := [2]win.TRIVERTEX{
		{
			X:     int32(bounds.X),
			Y:     int32(bounds.Y),
			Red:   uint16(color1.R()) * 256,
			Green: uint16(color1.G()) * 256,
			Blue:  uint16(color1.B()) * 256,
			Alpha: 0,
		}, {
			X:     int32(bounds.X + bounds.Width),
			Y:     int32(bounds.Y + bounds.Height),
			Red:   uint16(color2.R()) * 256,
			Green: uint16(color2.G()) * 256,
			Blue:  uint16(color2.B()) * 256,
			Alpha: 0,
		},
	}

	indices := win.GRADIENT_RECT{
		UpperLeft:  0,
		LowerRight: 1,
	}

	var o uint32
	if orientation == Vertical {
		o = 1
	}

	if !win.GradientFill(c.hdc, &vertices[0], 2, unsafe.Pointer(&indices), 1, o) {
		return newError("GradientFill failed")
	}

	return nil
}

// DrawText draws text at given location in 1/96" units.
//
// Deprecated: Newer applications should use DrawTextPixels.
func (c *Canvas) DrawText(text string, font *Font, color Color, bounds Rectangle, format DrawTextFormat) error {
	return c.DrawTextPixels(text, font, color, RectangleFrom96DPI(bounds, c.DPI()), format)
}

// DrawTextPixels draws text at given location in native pixels.
func (c *Canvas) DrawTextPixels(text string, font *Font, color Color, bounds Rectangle, format DrawTextFormat) error {
	return c.withFontAndTextColor(font, color, func() error {
		rect := bounds.toRECT()
		ret := win.DrawTextEx(
			c.hdc,
			syscall.StringToUTF16Ptr(text),
			-1,
			&rect,
			uint32(format)|win.DT_EDITCONTROL,
			nil)
		if ret == 0 {
			return newError("DrawTextEx failed")
		}

		return nil
	})
}

// fontHeight returns font height in native pixels.
func (c *Canvas) fontHeight(font *Font) (height int, err error) {
	err = c.withFontAndTextColor(font, 0, func() error {
		var size win.SIZE
		if !win.GetTextExtentPoint32(c.hdc, gM, 2, &size) {
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

// measureTextForDPI measures text for given DPI. Input and output bounds are in native pixels.
func (c *Canvas) measureTextForDPI(text string, font *Font, bounds Rectangle, format DrawTextFormat, dpi int) (boundsMeasured Rectangle, err error) {
	hFont := win.HGDIOBJ(font.handleForDPI(dpi))
	oldHandle := win.SelectObject(c.hdc, hFont)
	if oldHandle == 0 {
		err = newError("SelectObject failed")
		return
	}
	defer win.SelectObject(c.hdc, oldHandle)

	rect := &win.RECT{
		int32(bounds.X),
		int32(bounds.Y),
		int32(bounds.X + bounds.Width),
		int32(bounds.Y + bounds.Height),
	}
	var params win.DRAWTEXTPARAMS
	params.CbSize = uint32(unsafe.Sizeof(params))

	strPtr := syscall.StringToUTF16Ptr(text)
	dtfmt := uint32(format) | win.DT_CALCRECT | win.DT_EDITCONTROL | win.DT_NOPREFIX | win.DT_WORDBREAK

	height := win.DrawTextEx(
		c.hdc, strPtr, -1, rect, dtfmt, &params)
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

	return
}

// MeasureText measures text size. Input and output bounds are in 1/96" units.
//
// Deprecated: Newer applications should use MeasureTextPixels.
func (c *Canvas) MeasureText(text string, font *Font, bounds Rectangle, format DrawTextFormat) (boundsMeasured Rectangle, runesFitted int, err error) {
	dpi := c.DPI()
	var boundsMeasuredPixels Rectangle
	boundsMeasuredPixels, runesFitted, err = c.MeasureTextPixels(text, font, RectangleFrom96DPI(bounds, dpi), format)
	if err != nil {
		return
	}
	boundsMeasured = RectangleTo96DPI(boundsMeasuredPixels, dpi)
	return
}

// MeasureTextPixels measures text size. Input and output bounds are in native pixels.
func (c *Canvas) MeasureTextPixels(text string, font *Font, bounds Rectangle, format DrawTextFormat) (boundsMeasured Rectangle, runesFitted int, err error) {
	boundsMeasured, _, runesFitted, err = c.measureAndModifyTextPixels(text, font, bounds, format)
	return
}

// MeasureAndModifyTextPixels measures text size and also supports modification
// of the text which occurs if it does not fit into the specified bounds.
//
// Input and output bounds are in native pixels.
func (c *Canvas) MeasureAndModifyTextPixels(text string, font *Font, bounds Rectangle, format DrawTextFormat) (boundsMeasured Rectangle, textDisplayed string, err error) {
	var textPtr *uint16
	var runesFitted int
	if boundsMeasured, textPtr, runesFitted, err = c.measureAndModifyTextPixels(text, font, bounds, format|TextModifyString); err != nil {
		return
	}

	if runesFitted == utf8.RuneCountInString(text) {
		textDisplayed = text
	} else {
		if format&(TextEndEllipsis|TextPathEllipsis) != 0 {
			textDisplayed = win.UTF16PtrToString(textPtr)
		} else {
			textDisplayed = string(([]rune)(text)[:runesFitted])
		}
	}

	return
}

func (c *Canvas) measureAndModifyTextPixels(text string, font *Font, bounds Rectangle, format DrawTextFormat) (boundsMeasured Rectangle, textPtr *uint16, runesFitted int, err error) {
	// HACK: We don't want to actually draw on the Canvas here, but if we use
	// the DT_CALCRECT flag to avoid drawing, params.UiLengthDrawn will
	// not contain a useful value. To work around this, we create an in-memory
	// metafile and draw into that instead.
	if c.measureTextMetafile == nil {
		c.measureTextMetafile, err = NewMetafile(c)
		if err != nil {
			return
		}
	}

	hFont := win.HGDIOBJ(font.handleForDPI(c.DPI()))
	oldHandle := win.SelectObject(c.measureTextMetafile.hdc, hFont)
	if oldHandle == 0 {
		err = newError("SelectObject failed")
		return
	}
	defer win.SelectObject(c.measureTextMetafile.hdc, oldHandle)

	rect := &win.RECT{
		int32(bounds.X),
		int32(bounds.Y),
		int32(bounds.X + bounds.Width),
		int32(bounds.Y + bounds.Height),
	}
	var params win.DRAWTEXTPARAMS
	params.CbSize = uint32(unsafe.Sizeof(params))

	strPtr := syscall.StringToUTF16Ptr(text)
	dtfmt := uint32(format) | win.DT_EDITCONTROL | win.DT_WORDBREAK

	height := win.DrawTextEx(
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
	textPtr = strPtr
	runesFitted = int(params.UiLengthDrawn)

	return
}
