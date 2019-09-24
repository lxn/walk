// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

// Rectangle defines upper left corner with width and height region in 1/96".
type Rectangle struct {
	X, Y, Width, Height int
}

func (r Rectangle) Size() Size {
	return Size{r.Width, r.Height}
}

func scaleRectangle(value Rectangle, scale float64) RectanglePixels {
	return RectanglePixels{
		X:      scaleInt(value.X, scale),
		Y:      scaleInt(value.Y, scale),
		Width:  scaleInt(value.Width, scale),
		Height: scaleInt(value.Height, scale),
	}
}

// RectangleFrom96DPI converts from 1/96" units to native pixels.
func RectangleFrom96DPI(value Rectangle, dpi int) RectanglePixels {
	return scaleRectangle(value, float64(dpi)/96.0)
}

// RectanglePixels defines upper left corner with width and height region in native pixels.
type RectanglePixels struct {
	X, Y, Width, Height Pixel
}

func rectangleFromRECT(r win.RECT) RectanglePixels {
	return RectanglePixels{
		X:      Pixel(r.Left),
		Y:      Pixel(r.Top),
		Width:  Pixel(r.Right - r.Left),
		Height: Pixel(r.Bottom - r.Top),
	}
}

func (r RectanglePixels) Left() Pixel {
	return r.X
}

func (r RectanglePixels) Top() Pixel {
	return r.Y
}

func (r RectanglePixels) Right() Pixel {
	return r.X + r.Width - 1
}

func (r RectanglePixels) Bottom() Pixel {
	return r.Y + r.Height - 1
}

func (r RectanglePixels) Location() PointPixels {
	return PointPixels{r.X, r.Y}
}

func (r *RectanglePixels) SetLocation(p PointPixels) RectanglePixels {
	r.X = p.X
	r.Y = p.Y

	return *r
}

func (r RectanglePixels) Size() SizePixels {
	return SizePixels{r.Width, r.Height}
}

func (r *RectanglePixels) SetSize(s SizePixels) RectanglePixels {
	r.Width = s.Width
	r.Height = s.Height

	return *r
}

func (r RectanglePixels) toRECT() win.RECT {
	return win.RECT{
		int32(r.X),
		int32(r.Y),
		int32(r.X + r.Width),
		int32(r.Y + r.Height),
	}
}

func scaleRectanglePixels(value RectanglePixels, scale float64) Rectangle {
	return Rectangle{
		X:      scalePixel(value.X, scale),
		Y:      scalePixel(value.Y, scale),
		Width:  scalePixel(value.Width, scale),
		Height: scalePixel(value.Height, scale),
	}
}

// RectangleTo96DPI converts from native pixels to 1/96" units.
func RectangleTo96DPI(value RectanglePixels, dpi int) Rectangle {
	return scaleRectanglePixels(value, 96.0/float64(dpi))
}

// RectangleGrid measures X and Y of the rectangular area upper left corner, width and height in grid rows and columns.
type RectangleGrid struct {
	X, Y, Width, Height int
}
