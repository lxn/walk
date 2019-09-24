// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "math"

// IntFrom96DPI converts from 1/96" units to native pixels.
func IntFrom96DPI(value int, dpi int) Pixel {
	return scaleInt(value, float64(dpi)/96.0)
}

func scaleInt(value int, scale float64) Pixel {
	return Pixel(math.Round(float64(value) * scale))
}

// Pixel defines distance in native pixels.
type Pixel int

func maxPixel(a, b Pixel) Pixel {
	if a > b {
		return a
	}

	return b
}

func minPixel(a, b Pixel) Pixel {
	if a < b {
		return a
	}

	return b
}

func scalePixel(value Pixel, scale float64) int {
	return int(math.Round(float64(value) * scale))
}

// IntTo96DPI converts from native pixels to 1/96" units.
func IntTo96DPI(value Pixel, dpi int) int {
	return scalePixel(value, 96.0/float64(dpi))
}

// PixelDBU defines distance in dialog base units.
type PixelDBU int
