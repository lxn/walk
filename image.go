// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"strconv"
	"strings"

	"github.com/lxn/win"
)

type Image interface {
	// draw draws image at location (upper left) in native pixels unstreched.
	draw(hdc win.HDC, location Point) error

	// drawStretched draws image streched to given bounds in native pixels.
	drawStretched(hdc win.HDC, bounds Rectangle) error

	Dispose()

	// Size returns image size in 1/96" units.
	Size() Size
}

func ImageFrom(src interface{}) (img Image, err error) {
	switch src := src.(type) {
	case nil:
		// nop

	case Image:
		img = src

	case ExtractableIcon:
		img, err = NewIconExtractedFromFileWithSize(src.FilePath_(), src.Index_(), src.Size_())

	case int:
		img, err = Resources.Image(strconv.Itoa(src))

	case string:
		img, err = Resources.Image(src)

	default:
		err = ErrInvalidType
	}

	return
}

// NewImageFromFile loads image from file at 96dpi. Supported types are .ico, .emf, .bmp, .png...
//
// Deprecated: Newer applications should use NewImageFromFileForDPI.
func NewImageFromFile(filePath string) (Image, error) {
	return NewImageFromFileForDPI(filePath, 96)
}

// NewImageFromFileForDPI loads image from file at given DPI. Supported types are .ico, .emf,
// .bmp, .png...
func NewImageFromFileForDPI(filePath string, dpi int) (Image, error) {
	if strings.HasSuffix(filePath, ".ico") {
		return NewIconFromFile(filePath)
	} else if strings.HasSuffix(filePath, ".emf") {
		return NewMetafileFromFile(filePath)
	}

	return NewBitmapFromFileForDPI(filePath, dpi)
}

type PaintFuncImage struct {
	size96dpi   Size
	paint       PaintFunc // in 1/96" units
	paintPixels PaintFunc // in native pixels
	dispose     func()
}

// NewPaintFuncImage creates new PaintFuncImage struct. size parameter and paint function bounds
// parameter are specified in 1/96" units.
func NewPaintFuncImage(size Size, paint func(canvas *Canvas, bounds Rectangle) error) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paint: paint}
}

// NewPaintFuncImagePixels creates new PaintFuncImage struct. size parameter is specified in 1/96"
// units. paint function bounds parameter is specified in native pixels.
func NewPaintFuncImagePixels(size Size, paint func(canvas *Canvas, bounds Rectangle) error) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paintPixels: paint}
}

// NewPaintFuncImageWithDispose creates new PaintFuncImage struct. size parameter and paint
// function bounds parameter are specified in 1/96" units.
func NewPaintFuncImageWithDispose(size Size, paint func(canvas *Canvas, bounds Rectangle) error, dispose func()) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paint: paint, dispose: dispose}
}

// NewPaintFuncImagePixelsWithDispose creates new PaintFuncImage struct. size parameter is
// specified in 1/96" units. paint function bounds parameter is specified in native pixels.
func NewPaintFuncImagePixelsWithDispose(size Size, paint func(canvas *Canvas, bounds Rectangle) error, dispose func()) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paintPixels: paint, dispose: dispose}
}

func (pfi *PaintFuncImage) draw(hdc win.HDC, location Point) error {
	dpi := dpiForHDC(hdc)
	size := SizeFrom96DPI(pfi.size96dpi, dpi)

	return pfi.drawStretched(hdc, Rectangle{location.X, location.Y, size.Width, size.Height})
}

func (pfi *PaintFuncImage) drawStretched(hdc win.HDC, bounds Rectangle) error {
	canvas, err := newCanvasFromHDC(hdc)
	if err != nil {
		return err
	}
	defer canvas.Dispose()

	return pfi.drawStretchedOnCanvasPixels(canvas, bounds)
}

func (pfi *PaintFuncImage) drawStretchedOnCanvasPixels(canvas *Canvas, bounds Rectangle) error {
	if pfi.paintPixels != nil {
		return pfi.paintPixels(canvas, bounds)
	}
	if pfi.paint != nil {
		return pfi.paint(canvas, RectangleTo96DPI(bounds, canvas.DPI()))
	}

	return newError("paint(Pixels) func is nil")
}

func (pfi *PaintFuncImage) Dispose() {
	if pfi.dispose != nil {
		pfi.dispose()
		pfi.dispose = nil
	}
}

// Size returns image size in 1/96" units.
func (pfi *PaintFuncImage) Size() Size {
	return pfi.size96dpi
}
