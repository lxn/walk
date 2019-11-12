// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "github.com/lxn/win"

// Point defines 2D coordinate in 1/96" units ot native pixels.
type Point struct {
	X, Y int
}

func (p Point) toPOINT() win.POINT {
	return win.POINT{
		X: int32(p.X),
		Y: int32(p.Y),
	}
}

func pointPixelsFromPOINT(p win.POINT) Point {
	return Point{
		X: int(p.X),
		Y: int(p.Y),
	}
}
