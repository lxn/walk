// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type Size struct {
	Width, Height int
}

func (s Size) From96DPI(dpi int) Size {
	return scaleSize(s, float64(dpi)/96.0)
}

func (s Size) To96DPI(dpi int) Size {
	return scaleSize(s, 96.0/float64(dpi))
}

func scaleSize(value Size, scale float64) Size {
	return Size{
		Width:  scaleInt(value.Width, scale),
		Height: scaleInt(value.Height, scale),
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
