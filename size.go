// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "github.com/lxn/win"

// Size96DPI defines width and height in 1/96" units.
type Size96DPI struct {
	Width, Height Pixel96DPI
}

func minSize96DPI(a, b Size96DPI) Size96DPI {
	var s Size96DPI

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

func maxSize96DPI(a, b Size96DPI) Size96DPI {
	var s Size96DPI

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

// ForDPI converts from 1/96" units to native pixels.
func (s Size96DPI) ForDPI(dpi int) Size {
	return s.scale(float64(dpi) / 96.0)
}

func (s Size96DPI) scale(scale float64) Size {
	return Size{
		Width:  s.Width.scale(scale),
		Height: s.Height.scale(scale),
	}
}

// SizeFrom96DPI converts from 1/96" units to native pixels.
func SizeFrom96DPI(value Size96DPI, dpi int) Size {
	return value.ForDPI(dpi)
}

// Size defines width and height in native pixels.
type Size struct {
	Width, Height Pixel
}

func sizeFromSIZE(s win.SIZE) Size {
	return Size{
		Width:  Pixel(s.CX),
		Height: Pixel(s.CY),
	}
}

func sizeFromRECT(r win.RECT) Size {
	return Size{
		Width:  Pixel(r.Right - r.Left),
		Height: Pixel(r.Bottom - r.Top),
	}
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

// To96DPI converts from native pixels to 1/96" units.
func (s Size) To96DPI(dpi int) Size96DPI {
	return s.scale(96.0 / float64(dpi))
}

func (s Size) scale(scale float64) Size96DPI {
	return Size96DPI{
		Width:  s.Width.scale(scale),
		Height: s.Height.scale(scale),
	}
}

// SizeTo96DPI converts from native pixels to 1/96" units.
func SizeTo96DPI(value Size, dpi int) Size96DPI {
	return value.To96DPI(dpi)
}

// SizeDBU defines width and height in dialog base units.
type SizeDBU struct {
	Width, Height PixelDBU
}
