// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drawing

import (
	"os"
)

import (
	. "walk/winapi/gdi32"
)

type HatchStyle int

const (
	HatchHorizontal       HatchStyle = HS_HORIZONTAL
	HatchVertical         HatchStyle = HS_VERTICAL
	HatchForwardDiagonal  HatchStyle = HS_FDIAGONAL
	HatchBackwardDiagonal HatchStyle = HS_BDIAGONAL
	HatchCross            HatchStyle = HS_CROSS
	HatchDiagonalCross    HatchStyle = HS_DIAGCROSS
)

type Brush interface {
	Dispose()
	handle() HBRUSH
	logbrush() *LOGBRUSH
}

type nullBrush struct {
	hBrush HBRUSH
}

func newNullBrush() *nullBrush {
	lb := &LOGBRUSH{LbStyle: BS_NULL}

	hBrush := CreateBrushIndirect(lb)
	if hBrush == 0 {
		panic("failed to create null brush")
	}

	return &nullBrush{hBrush: hBrush}
}

func (b *nullBrush) Dispose() {
	if b.hBrush != 0 {
		DeleteObject(HGDIOBJ(b.hBrush))

		b.hBrush = 0
	}
}

func (b *nullBrush) handle() HBRUSH {
	return b.hBrush
}

func (b *nullBrush) logbrush() *LOGBRUSH {
	return &LOGBRUSH{LbStyle: BS_NULL}
}

var nullBrushSingleton Brush = newNullBrush()

func NullBrush() Brush {
	return nullBrushSingleton
}

type SolidColorBrush struct {
	hBrush HBRUSH
	color  Color
}

func NewSolidColorBrush(color Color) (*SolidColorBrush, os.Error) {
	lb := &LOGBRUSH{LbStyle: BS_SOLID, LbColor: COLORREF(color)}

	hBrush := CreateBrushIndirect(lb)
	if hBrush == 0 {
		return nil, newError("CreateBrushIndirect failed")
	}

	return &SolidColorBrush{hBrush: hBrush, color: color}, nil
}

func (b *SolidColorBrush) Color() Color {
	return b.color
}

func (b *SolidColorBrush) Dispose() {
	if b.hBrush != 0 {
		DeleteObject(HGDIOBJ(b.hBrush))

		b.hBrush = 0
	}
}

func (b *SolidColorBrush) handle() HBRUSH {
	return b.hBrush
}

func (b *SolidColorBrush) logbrush() *LOGBRUSH {
	return &LOGBRUSH{LbStyle: BS_SOLID, LbColor: COLORREF(b.color)}
}

type HatchBrush struct {
	hBrush HBRUSH
	color  Color
	style  HatchStyle
}

func NewHatchBrush(color Color, style HatchStyle) (*HatchBrush, os.Error) {
	lb := &LOGBRUSH{LbStyle: BS_HATCHED, LbColor: COLORREF(color), LbHatch: uintptr(style)}

	hBrush := CreateBrushIndirect(lb)
	if hBrush == 0 {
		return nil, newError("CreateBrushIndirect failed")
	}

	return &HatchBrush{hBrush: hBrush, color: color, style: style}, nil
}

func (b *HatchBrush) Color() Color {
	return b.color
}

func (b *HatchBrush) Dispose() {
	if b.hBrush != 0 {
		DeleteObject(HGDIOBJ(b.hBrush))

		b.hBrush = 0
	}
}

func (b *HatchBrush) handle() HBRUSH {
	return b.hBrush
}

func (b *HatchBrush) logbrush() *LOGBRUSH {
	return &LOGBRUSH{LbStyle: BS_HATCHED, LbColor: COLORREF(b.color), LbHatch: uintptr(b.style)}
}

func (b *HatchBrush) Style() HatchStyle {
	return b.style
}
