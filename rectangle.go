// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

// Rectangle96DPI defines upper left corner with width and height region in 1/96".
type Rectangle96DPI struct {
	X, Y, Width, Height Pixel96DPI
}

// ForDPI converts from 1/96" units to native pixels.
func (r Rectangle96DPI) ForDPI(dpi int) Rectangle {
	return r.scale(float64(dpi) / 96.0)
}

func (r Rectangle96DPI) scale(scale float64) Rectangle {
	return Rectangle{
		X:      r.X.scale(scale),
		Y:      r.Y.scale(scale),
		Width:  r.Width.scale(scale),
		Height: r.Height.scale(scale),
	}
}

// RectangleFrom96DPI converts from 1/96" units to native pixels.
func RectangleFrom96DPI(value Rectangle96DPI, dpi int) Rectangle {
	return value.ForDPI(dpi)
}

// Rectangle defines upper left corner with width and height region in native pixels.
type Rectangle struct {
	X, Y, Width, Height Pixel
}

func rectangleFromRECT(r win.RECT) Rectangle {
	return Rectangle{
		X:      Pixel(r.Left),
		Y:      Pixel(r.Top),
		Width:  Pixel(r.Right - r.Left),
		Height: Pixel(r.Bottom - r.Top),
	}
}

func (r Rectangle) Left() Pixel {
	return r.X
}

func (r Rectangle) Top() Pixel {
	return r.Y
}

func (r Rectangle) Right() Pixel {
	return r.X + r.Width - 1
}

func (r Rectangle) Bottom() Pixel {
	return r.Y + r.Height - 1
}

func (r Rectangle) Location() Point {
	return Point{r.X, r.Y}
}

func (r *Rectangle) SetLocation(p Point) Rectangle {
	r.X = p.X
	r.Y = p.Y

	return *r
}

func (r Rectangle) Size() Size {
	return Size{r.Width, r.Height}
}

func (r *Rectangle) SetSize(s Size) Rectangle {
	r.Width = s.Width
	r.Height = s.Height

	return *r
}

func (r Rectangle) toRECT() win.RECT {
	return win.RECT{
		int32(r.X),
		int32(r.Y),
		int32(r.X + r.Width),
		int32(r.Y + r.Height),
	}
}

// To96DPI converts from native pixels to 1/96" units.
func (r Rectangle) To96DPI(dpi int) Rectangle96DPI {
	return r.scale(96.0 / float64(dpi))
}

func (r Rectangle) scale(scale float64) Rectangle96DPI {
	return Rectangle96DPI{
		X:      r.X.scale(scale),
		Y:      r.Y.scale(scale),
		Width:  r.Width.scale(scale),
		Height: r.Height.scale(scale),
	}
}

// RectangleTo96DPI converts from native pixels to 1/96" units.
func RectangleTo96DPI(value Rectangle, dpi int) Rectangle96DPI {
	return value.To96DPI(dpi)
}

// RectangleGrid measures X and Y of the rectangular area upper left corner, width and height in grid rows and columns.
type RectangleGrid struct {
	X, Y, Width, Height int
}
