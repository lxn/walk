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
	draw(hdc win.HDC, location Point) error
	drawStretched(hdc win.HDC, bounds Rectangle) error
	Dispose()
	Size() Size96DPI
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

func NewImageFromFile(filePath string) (Image, error) {
	if strings.HasSuffix(filePath, ".ico") {
		return NewIconFromFile(filePath)
	} else if strings.HasSuffix(filePath, ".emf") {
		return NewMetafileFromFile(filePath)
	}

	return NewBitmapFromFile(filePath)
}

type PaintFuncImage struct {
	size96dpi Size96DPI
	paint     func(canvas *Canvas, bounds Rectangle) error
	dispose   func()
}

func NewPaintFuncImage(size Size96DPI, paint func(canvas *Canvas, bounds Rectangle) error) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paint: paint}
}

func NewPaintFuncImageWithDispose(size Size96DPI, paint func(canvas *Canvas, bounds Rectangle) error, dispose func()) *PaintFuncImage {
	return &PaintFuncImage{size96dpi: size, paint: paint, dispose: dispose}
}

func (pfi *PaintFuncImage) draw(hdc win.HDC, location Point) error {
	dpi := dpiForHDC(hdc)
	size := pfi.size96dpi.ForDPI(dpi)

	return pfi.drawStretched(hdc, Rectangle{location.X, location.Y, size.Width, size.Height})
}

func (pfi *PaintFuncImage) drawStretched(hdc win.HDC, bounds Rectangle) error {
	canvas, err := newCanvasFromHDC(hdc)
	if err != nil {
		return err
	}
	defer canvas.Dispose()

	return pfi.drawStretchedOnCanvas(canvas, bounds)
}

func (pfi *PaintFuncImage) drawStretchedOnCanvas(canvas *Canvas, bounds Rectangle) error {
	return pfi.paint(canvas, bounds)
}

func (pfi *PaintFuncImage) Dispose() {
	if pfi.dispose != nil {
		pfi.dispose()
		pfi.dispose = nil
	}
}

func (pfi *PaintFuncImage) Size() Size96DPI {
	return pfi.size96dpi
}
