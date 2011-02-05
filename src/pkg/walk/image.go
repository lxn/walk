// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"strings"
)

import (
	. "walk/winapi/gdi32"
)

type Image interface {
	draw(hdc HDC, location Point) os.Error
	drawStretched(hdc HDC, bounds Rectangle) os.Error
	Dispose()
	Size() Size
}

func NewImageFromFile(filePath string) (Image, os.Error) {
	if strings.HasSuffix(filePath, ".emf") {
		return NewMetafileFromFile(filePath)
	}

	return NewBitmapFromFile(filePath)
}
