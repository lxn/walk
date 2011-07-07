// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import . "walk/winapi"

type Metafile struct {
	hdc  HDC
	hemf HENHMETAFILE
	size Size
}

func NewMetafile(referenceCanvas *Canvas) (*Metafile, os.Error) {
	hdc := CreateEnhMetaFile(referenceCanvas.hdc, nil, nil, nil)
	if hdc == 0 {
		return nil, newError("CreateEnhMetaFile failed")
	}

	return &Metafile{hdc: hdc}, nil
}

func NewMetafileFromFile(filePath string) (*Metafile, os.Error) {
	hemf := GetEnhMetaFile(syscall.StringToUTF16Ptr(filePath))
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
		DeleteEnhMetaFile(mf.hemf)

		mf.hemf = 0
	}
}

func (mf *Metafile) Save(filePath string) os.Error {
	hemf := CopyEnhMetaFile(mf.hemf, syscall.StringToUTF16Ptr(filePath))
	if hemf == 0 {
		return newError("CopyEnhMetaFile failed")
	}

	DeleteEnhMetaFile(hemf)

	return nil
}

func (mf *Metafile) readSizeFromHeader() os.Error {
	var hdr ENHMETAHEADER

	if GetEnhMetaFileHeader(mf.hemf, uint(unsafe.Sizeof(hdr)), &hdr) == 0 {
		return newError("GetEnhMetaFileHeader failed")
	}

	mf.size = Size{
		hdr.RclBounds.Right - hdr.RclBounds.Left,
		hdr.RclBounds.Bottom - hdr.RclBounds.Top,
	}

	return nil
}

func (mf *Metafile) ensureFinished() os.Error {
	if mf.hdc == 0 {
		if mf.hemf == 0 {
			return newError("already disposed")
		} else {
			return nil
		}
	}

	mf.hemf = CloseEnhMetaFile(mf.hdc)
	if mf.hemf == 0 {
		return newError("CloseEnhMetaFile failed")
	}

	mf.hdc = 0

	return mf.readSizeFromHeader()
}

func (mf *Metafile) Size() Size {
	return mf.size
}

func (mf *Metafile) draw(hdc HDC, location Point) os.Error {
	return mf.drawStretched(hdc, Rectangle{location.X, location.Y, mf.size.Width, mf.size.Height})
}

func (mf *Metafile) drawStretched(hdc HDC, bounds Rectangle) os.Error {
	rc := bounds.toRECT()

	if !PlayEnhMetaFile(hdc, mf.hemf, &rc) {
		return newError("PlayEnhMetaFile failed")
	}

	return nil
}
