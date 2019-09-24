// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "github.com/lxn/win"

// Point96DPI defines 2D coordinate in 1/96" units.
type Point96DPI struct {
	X, Y Pixel96DPI
}

// ForDPI converts from 1/96" units to native pixels.
func (p Point96DPI) ForDPI(dpi int) Point {
	return p.scale(float64(dpi) / 96.0)
}

func (p Point96DPI) scale(scale float64) Point {
	return Point{
		X: p.X.scale(scale),
		Y: p.Y.scale(scale),
	}
}

// PointFrom96DPI converts from 1/96" units to native pixels.
func PointFrom96DPI(value Point96DPI, dpi int) Point {
	return value.ForDPI(dpi)
}

// Point defines 2D coordinate in native pixels.
type Point struct {
	X, Y Pixel
}

func (p Point) toPOINT() win.POINT {
	return win.POINT{
		X: int32(p.X),
		Y: int32(p.Y),
	}
}

// To96DPI converts from native pixels to 1/96" units.
func (p Point) To96DPI(dpi int) Point96DPI {
	return p.scale(96.0 / float64(dpi))
}

func (p Point) scale(scale float64) Point96DPI {
	return Point96DPI{
		X: p.X.scale(scale),
		Y: p.Y.scale(scale),
	}
}

// PointTo96DPI converts from native pixels to 1/96" units.
func PointTo96DPI(value Point, dpi int) Point96DPI {
	return value.To96DPI(dpi)
}
