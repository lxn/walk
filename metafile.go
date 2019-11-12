// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"math"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const milimeterPerMeter float64 = 1000.0

type Metafile struct {
	hdc  win.HDC
	hemf win.HENHMETAFILE
	size Size // in native pixels
	dpi  Size
}

func NewMetafile(referenceCanvas *Canvas) (*Metafile, error) {
	hdc := win.CreateEnhMetaFile(referenceCanvas.hdc, nil, nil, nil)
	if hdc == 0 {
		return nil, newError("CreateEnhMetaFile failed")
	}

	return &Metafile{hdc: hdc}, nil
}

func NewMetafileFromFile(filePath string) (*Metafile, error) {
	hemf := win.GetEnhMetaFile(syscall.StringToUTF16Ptr(filePath))
	if hemf == 0 {
		return nil, newError("GetEnhMetaFile failed")
	}

	mf := &Metafile{hemf: hemf}

	err := mf.readSizeFromHeader()
	if err != nil {
		return nil, err
	}

	return mf, nil
}

func (mf *Metafile) Dispose() {
	mf.ensureFinished()

	if mf.hemf != 0 {
		win.DeleteEnhMetaFile(mf.hemf)

		mf.hemf = 0
	}
}

func (mf *Metafile) Save(filePath string) error {
	hemf := win.CopyEnhMetaFile(mf.hemf, syscall.StringToUTF16Ptr(filePath))
	if hemf == 0 {
		return newError("CopyEnhMetaFile failed")
	}

	win.DeleteEnhMetaFile(hemf)

	return nil
}

func (mf *Metafile) readSizeFromHeader() error {
	var hdr win.ENHMETAHEADER

	if win.GetEnhMetaFileHeader(mf.hemf, uint32(unsafe.Sizeof(hdr)), &hdr) == 0 {
		return newError("GetEnhMetaFileHeader failed")
	}

	mf.size = sizeFromRECT(hdr.RclBounds)
	scale := milimeterPerMeter / inchesPerMeter
	mf.dpi = Size{
		int(math.Round(float64(hdr.SzlDevice.CX) / float64(hdr.SzlMillimeters.CX) * scale)),
		int(math.Round(float64(hdr.SzlDevice.CY) / float64(hdr.SzlMillimeters.CY) * scale)),
	}

	return nil
}

func (mf *Metafile) ensureFinished() error {
	if mf.hdc == 0 {
		if mf.hemf == 0 {
			return newError("already disposed")
		} else {
			return nil
		}
	}

	mf.hemf = win.CloseEnhMetaFile(mf.hdc)
	if mf.hemf == 0 {
		return newError("CloseEnhMetaFile failed")
	}

	mf.hdc = 0

	return mf.readSizeFromHeader()
}

// Size returns image size in 1/96" units.
func (mf *Metafile) Size() Size {
	return Size{
		Width:  scaleInt(mf.size.Width, 96.0/float64(mf.dpi.Width)),
		Height: scaleInt(mf.size.Height, 96.0/float64(mf.dpi.Height)),
	}
}

func (mf *Metafile) draw(hdc win.HDC, location Point) error {
	return mf.drawStretched(hdc, Rectangle{location.X, location.Y, mf.size.Width, mf.size.Height})
}

func (mf *Metafile) drawStretched(hdc win.HDC, bounds Rectangle) error {
	rc := bounds.toRECT()

	if !win.PlayEnhMetaFile(hdc, mf.hemf, &rc) {
		return newError("PlayEnhMetaFile failed")
	}

	return nil
}
