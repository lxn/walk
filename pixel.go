// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "math"

// Pixel96DPI defines distance in 1/96" units.
type Pixel96DPI int

// ForDPI converts from 1/96" units to native pixels.
func (p Pixel96DPI) ForDPI(dpi int) Pixel {
	return p.scale(float64(dpi) / 96.0)
}

func (p Pixel96DPI) scale(scale float64) Pixel {
	return Pixel(math.Round(float64(p) * scale))
}

func assertPixel96DPIOr(value interface{}, defaultValue Pixel96DPI) Pixel96DPI {
	if n, ok := value.(Pixel96DPI); ok {
		return n
	}
	if n, ok := value.(int); ok {
		return Pixel96DPI(n)
	}

	return defaultValue
}

// IntFrom96DPI converts from 1/96" units to native pixels.
func IntFrom96DPI(value Pixel96DPI, dpi int) Pixel {
	return value.ForDPI(dpi)
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

// To96DPI converts from native pixels to 1/96" units.
func (p Pixel) To96DPI(dpi int) Pixel96DPI {
	return p.scale(96.0 / float64(dpi))
}

func (p Pixel) scale(scale float64) Pixel96DPI {
	return Pixel96DPI(math.Round(float64(p) * scale))
}

// IntTo96DPI converts from native pixels to 1/96" units.
func IntTo96DPI(value Pixel, dpi int) Pixel96DPI {
	return value.To96DPI(dpi)
}

// PixelDBU defines distance in dialog base units.
type PixelDBU int
