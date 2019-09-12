// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type Point struct {
	X, Y int
}

func (p Point) From96DPI(dpi int) Point {
	return scalePoint(p, float64(dpi)/96.0)
}

func (p Point) To96DPI(dpi int) Point {
	return scalePoint(p, 96.0/float64(dpi))
}

func scalePoint(value Point, scale float64) Point {
	return Point{
		X: scaleInt(value.X, scale),
		Y: scaleInt(value.Y, scale),
	}
}
