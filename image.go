// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"strings"
)

import (
	. "github.com/lxn/go-winapi"
)

type Image interface {
	draw(hdc HDC, location Point) error
	drawStretched(hdc HDC, bounds Rectangle) error
	Dispose()
	Size() Size
}

func NewImageFromFile(filePath string) (Image, error) {
	if strings.HasSuffix(filePath, ".emf") {
		return NewMetafileFromFile(filePath)
	}

	return NewBitmapFromFile(filePath)
}
