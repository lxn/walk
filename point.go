// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "github.com/lxn/win"

// Point defines 2D coordinate in 1/96" units.
type Point struct {
	X, Y int
}

func scalePoint(value Point, scale float64) PointPixels {
	return PointPixels{
		X: scaleInt(value.X, scale),
		Y: scaleInt(value.Y, scale),
	}
}

// PointFrom96DPI converts from 1/96" units to native pixels.
func PointFrom96DPI(value Point, dpi int) PointPixels {
	return scalePoint(value, float64(dpi)/96.0)
}

// PointPixels defines 2D coordinate in native pixels.
type PointPixels struct {
	X, Y int
}

func (p PointPixels) toPOINT() win.POINT {
	return win.POINT{
		X: int32(p.X),
		Y: int32(p.Y),
	}
}

func pointPixelsFromPOINT(p win.POINT) PointPixels {
	return PointPixels{
		X: int(p.X),
		Y: int(p.Y),
	}
}

func scalePointPixels(value PointPixels, scale float64) Point {
	return Point{
		X: scaleInt(value.X, scale),
		Y: scaleInt(value.Y, scale),
	}
}

// PointTo96DPI converts from native pixels to 1/96" units.
func PointTo96DPI(value PointPixels, dpi int) Point {
	return scalePointPixels(value, 96.0/float64(dpi))
}
