// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "github.com/lxn/go-winapi"
)

type PenStyle int

// Pen styles
const (
	PenSolid       PenStyle = PS_SOLID
	PenDash        PenStyle = PS_DASH
	PenDot         PenStyle = PS_DOT
	PenDashDot     PenStyle = PS_DASHDOT
	PenDashDotDot  PenStyle = PS_DASHDOTDOT
	PenNull        PenStyle = PS_NULL
	PenInsideFrame PenStyle = PS_INSIDEFRAME
	PenUserStyle   PenStyle = PS_USERSTYLE
	PenAlternate   PenStyle = PS_ALTERNATE
)

// Pen cap styles (geometric pens only)
const (
	PenCapRound  PenStyle = PS_ENDCAP_ROUND
	PenCapSquare PenStyle = PS_ENDCAP_SQUARE
	PenCapFlat   PenStyle = PS_ENDCAP_FLAT
)

// Pen join styles (geometric pens only)
const (
	PenJoinBevel PenStyle = PS_JOIN_BEVEL
	PenJoinMiter PenStyle = PS_JOIN_MITER
	PenJoinRound PenStyle = PS_JOIN_ROUND
)

type Pen interface {
	handle() HPEN
	Dispose()
	Style() PenStyle
	Width() int
}

type nullPen struct {
	hPen HPEN
}

func newNullPen() *nullPen {
	lb := &LOGBRUSH{LbStyle: BS_NULL}

	hPen := ExtCreatePen(PS_COSMETIC|PS_NULL, 1, lb, 0, nil)
	if hPen == 0 {
		panic("failed to create null brush")
	}

	return &nullPen{hPen: hPen}
}

func (p *nullPen) Dispose() {
	if p.hPen != 0 {
		DeleteObject(HGDIOBJ(p.hPen))

		p.hPen = 0
	}
}

func (p *nullPen) handle() HPEN {
	return p.hPen
}

func (p *nullPen) Style() PenStyle {
	return PenNull
}

func (p *nullPen) Width() int {
	return 0
}

var nullPenSingleton Pen = newNullPen()

func NullPen() Pen {
	return nullPenSingleton
}

type CosmeticPen struct {
	hPen  HPEN
	style PenStyle
	color Color
}

func NewCosmeticPen(style PenStyle, color Color) (*CosmeticPen, error) {
	lb := &LOGBRUSH{LbStyle: BS_SOLID, LbColor: COLORREF(color)}

	style |= PS_COSMETIC

	hPen := ExtCreatePen(uint32(style), 1, lb, 0, nil)
	if hPen == 0 {
		return nil, newError("ExtCreatePen failed")
	}

	return &CosmeticPen{hPen: hPen, style: style, color: color}, nil
}

func (p *CosmeticPen) Dispose() {
	if p.hPen != 0 {
		DeleteObject(HGDIOBJ(p.hPen))

		p.hPen = 0
	}
}

func (p *CosmeticPen) handle() HPEN {
	return p.hPen
}

func (p *CosmeticPen) Style() PenStyle {
	return p.style
}

func (p *CosmeticPen) Color() Color {
	return p.color
}

func (p *CosmeticPen) Width() int {
	return 1
}

type GeometricPen struct {
	hPen  HPEN
	style PenStyle
	brush Brush
	width int
}

func NewGeometricPen(style PenStyle, width int, brush Brush) (*GeometricPen, error) {
	if brush == nil {
		return nil, newError("brush cannot be nil")
	}

	style |= PS_GEOMETRIC

	hPen := ExtCreatePen(uint32(style), uint32(width), brush.logbrush(), 0, nil)
	if hPen == 0 {
		return nil, newError("ExtCreatePen failed")
	}

	return &GeometricPen{
		hPen:  hPen,
		style: style,
		width: width,
		brush: brush,
	}, nil
}

func (p *GeometricPen) Dispose() {
	if p.hPen != 0 {
		DeleteObject(HGDIOBJ(p.hPen))

		p.hPen = 0
	}
}

func (p *GeometricPen) handle() HPEN {
	return p.hPen
}

func (p *GeometricPen) Style() PenStyle {
	return p.style
}

func (p *GeometricPen) Width() int {
	return p.width
}

func (p *GeometricPen) Brush() Brush {
	return p.brush
}
