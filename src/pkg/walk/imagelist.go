// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/comctl32"
	. "walk/winapi/gdi32"
)

type ImageList struct {
	hIml      HIMAGELIST
	maskColor Color
}

func NewImageList(imageSize Size, maskColor Color) (*ImageList, os.Error) {
	hIml := ImageList_Create(imageSize.Width, imageSize.Height, ILC_MASK|ILC_COLOR24, 8, 8)
	if hIml == 0 {
		return nil, newError("ImageList_Create failed")
	}

	return &ImageList{hIml: hIml, maskColor: maskColor}, nil
}

func (il *ImageList) Add(bitmap, maskBitmap *Bitmap) (int, os.Error) {
	if bitmap == nil {
		return 0, newError("bitmap cannot be nil")
	}

	var maskHandle HBITMAP
	if maskBitmap != nil {
		maskHandle = maskBitmap.handle()
	}

	index := ImageList_Add(il.hIml, bitmap.handle(), maskHandle)
	if index == -1 {
		return 0, newError("ImageList_Add failed")
	}

	return index, nil
}

func (il *ImageList) AddMasked(bitmap *Bitmap) (int, os.Error) {
	if bitmap == nil {
		return 0, newError("bitmap cannot be nil")
	}

	index := ImageList_AddMasked(il.hIml, bitmap.handle(), COLORREF(il.maskColor))
	if index == -1 {
		return 0, newError("ImageList_AddMasked failed")
	}

	return index, nil
}

func (il *ImageList) Dispose() {
	if il.hIml != 0 {
		ImageList_Destroy(il.hIml)
		il.hIml = 0
	}
}

func (il *ImageList) MaskColor() Color {
	return il.maskColor
}
