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
	draw(hdc win.HDC, location PointPixels) error
	drawStretched(hdc win.HDC, bounds RectanglePixels) error
	Dispose()
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

// NewImageFromFile loads image from file with 96dpi. Supported types are .ico, .emf, .bmp, .png...
//
// Deprecated: Newer applications should use DPI-aware variant.
func NewImageFromFile(filePath string) (Image, error) {
	return NewImageFromFileForDPI(filePath, 96)
}

// NewImageFromFileForDPI loads image from file with given DPI. Supported types are .ico, .emf,
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
	size96dpi Size
	paint     func(canvas *Canvas, bounds RectanglePixels) error
	dispose   func()
}

func NewPaintFuncImage(size Size, paint func(canvas *Canvas, bounds RectanglePixels) error) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paint: paint}
}

func NewPaintFuncImageWithDispose(size Size, paint func(canvas *Canvas, bounds RectanglePixels) error, dispose func()) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paint: paint, dispose: dispose}
}

func (pfi *PaintFuncImage) draw(hdc win.HDC, location PointPixels) error {
	dpi := dpiForHDC(hdc)
	size := SizeFrom96DPI(pfi.size96dpi, dpi)

	return pfi.drawStretched(hdc, RectanglePixels{location.X, location.Y, size.Width, size.Height})
}

func (pfi *PaintFuncImage) drawStretched(hdc win.HDC, bounds RectanglePixels) error {
	canvas, err := newCanvasFromHDC(hdc)
	if err != nil {
		return err
	}
	defer canvas.Dispose()

	return pfi.drawStretchedOnCanvas(canvas, bounds)
}

func (pfi *PaintFuncImage) drawStretchedOnCanvas(canvas *Canvas, bounds RectanglePixels) error {
	return pfi.paint(canvas, bounds)
}

func (pfi *PaintFuncImage) Dispose() {
	if pfi.dispose != nil {
		pfi.dispose()
		pfi.dispose = nil
	}
}

func (pfi *PaintFuncImage) Size() Size {
	return pfi.size96dpi
}
