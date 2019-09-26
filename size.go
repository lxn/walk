// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "github.com/lxn/win"

// Size defines width and height in 1/96" units.
type Size struct {
	Width, Height int
}

func minSize(a, b Size) Size {
	var s Size

	if a.Width < b.Width {
		s.Width = a.Width
	} else {
		s.Width = b.Width
	}

	if a.Height < b.Height {
		s.Height = a.Height
	} else {
		s.Height = b.Height
	}

	return s
}

func maxSize(a, b Size) Size {
	var s Size

	if a.Width > b.Width {
		s.Width = a.Width
	} else {
		s.Width = b.Width
	}

	if a.Height > b.Height {
		s.Height = a.Height
	} else {
		s.Height = b.Height
	}

	return s
}

func scaleSize(value Size, scale float64) SizePixels {
	return SizePixels{
		Width:  scaleInt(value.Width, scale),
		Height: scaleInt(value.Height, scale),
	}
}

// SizeFrom96DPI converts from 1/96" units to native pixels.
func SizeFrom96DPI(value Size, dpi int) SizePixels {
	return scaleSize(value, float64(dpi)/96.0)
}

// SizePixels defines width and height in native pixels.
type SizePixels struct {
	Width, Height int
}

func sizeFromSIZE(s win.SIZE) SizePixels {
	return SizePixels{
		Width:  int(s.CX),
		Height: int(s.CY),
	}
}

func sizeFromRECT(r win.RECT) SizePixels {
	return SizePixels{
		Width:  int(r.Right - r.Left),
		Height: int(r.Bottom - r.Top),
	}
}

func minSizePixels(a, b SizePixels) SizePixels {
	var s SizePixels

	if a.Width < b.Width {
		s.Width = a.Width
	} else {
		s.Width = b.Width
	}

	if a.Height < b.Height {
		s.Height = a.Height
	} else {
		s.Height = b.Height
	}

	return s
}

func maxSizePixels(a, b SizePixels) SizePixels {
	var s SizePixels

	if a.Width > b.Width {
		s.Width = a.Width
	} else {
		s.Width = b.Width
	}

	if a.Height > b.Height {
		s.Height = a.Height
	} else {
		s.Height = b.Height
	}

	return s
}

func scaleSizePixels(value SizePixels, scale float64) Size {
	return Size{
		Width:  scaleInt(value.Width, scale),
		Height: scaleInt(value.Height, scale),
	}
}

// SizeTo96DPI converts from native pixels to 1/96" units.
func SizeTo96DPI(value SizePixels, dpi int) Size {
	return scaleSizePixels(value, 96.0/float64(dpi))
}

// SizeDBU defines width and height in dialog base units.
type SizeDBU struct {
	Width, Height int
}
