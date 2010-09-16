// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	"walk/drawing"
	. "walk/winapi/comctl32"
	. "walk/winapi/gdi32"
	. "walk/winapi/user32"
)

type ImageList struct {
	hIml HIMAGELIST
}

func loadImageList(filePath string, imageWidth int, transparentColor COLORREF) HIMAGELIST {
	hIml := ImageList_LoadImage(
		0,
		syscall.StringToUTF16Ptr(filePath),
		imageWidth,
		8,
		transparentColor,
		IMAGE_BITMAP,
		LR_LOADFROMFILE)

	return hIml
}

func NewImageList(fileName string, imageWidth int, color drawing.Color) (*ImageList, os.Error) {
	hIml := loadImageList(fileName, imageWidth, COLORREF(color))
	if hIml == 0 {
		return nil, newError("ImageList_LoadImage failed")
	}

	return &ImageList{hIml: hIml}, nil
}

func (il *ImageList) Dispose() {
	if il.hIml != 0 {
		ImageList_Destroy(il.hIml)
		il.hIml = 0
	}
}

func (il *ImageList) IsDisposed() bool {
	return il.hIml == 0
}
