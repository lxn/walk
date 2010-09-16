// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drawing

import (
	"walk/winapi/gdi32"
)

type Rectangle struct {
	X, Y, Width, Height int
}

func (r *Rectangle) Left() int {
	return r.X
}

func (r *Rectangle) Top() int {
	return r.Y
}

func (r *Rectangle) Right() int {
	return r.X + r.Width - 1
}

func (r *Rectangle) Bottom() int {
	return r.Y + r.Height - 1
}

func (r *Rectangle) HCenter() int {
	return r.X + r.Width/2
}

func (r *Rectangle) VCenter() int {
	return r.Y + r.Height/2
}

func (r *Rectangle) Indent(value int) *Rectangle {
	valueTwice := value * 2

	r.X += value
	r.Y += value
	r.Width -= valueTwice
	r.Height -= valueTwice

	return r
}

func (r *Rectangle) HIndent(value int) *Rectangle {
	r.X += value
	r.Width -= value * 2

	return r
}

func (r *Rectangle) VIndent(value int) *Rectangle {
	r.Y += value
	r.Height -= value * 2

	return r
}

func (r *Rectangle) Size() Size {
	return Size{r.Width, r.Height}
}

func (r *Rectangle) SetSize(s Size) *Rectangle {
	r.Width = s.Width
	r.Height = s.Height

	return r
}

func (r *Rectangle) toRECT() *gdi32.RECT {
	return &gdi32.RECT{r.X, r.Y, r.X + r.Width, r.Y + r.Height}
}
