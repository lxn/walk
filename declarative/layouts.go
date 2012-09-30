// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Orientation byte

const (
	Horizontal Orientation = Orientation(walk.Horizontal)
	Vertical   Orientation = Orientation(walk.Vertical)
)

type Margins struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func (m *Margins) toW() walk.Margins {
	if m == nil {
		return walk.Margins{9, 9, 9, 9}
	}

	return walk.Margins{m.Left, m.Top, m.Right, m.Bottom}
}

type Size struct {
	Width  int
	Height int
}

func (s Size) toW() walk.Size {
	return walk.Size{s.Width, s.Height}
}

type HBox struct {
	Margins *Margins
	Spacing int
}

func (hb HBox) Create() (walk.Layout, error) {
	l := walk.NewHBoxLayout()

	if err := l.SetMargins(hb.Margins.toW()); err != nil {
		return nil, err
	}

	var s int
	switch hb.Spacing {
	case -1:
		s = 0

	case 0:
		s = 6

	default:
		s = hb.Spacing
	}
	if err := l.SetSpacing(s); err != nil {
		return nil, err
	}

	return l, nil
}

type VBox struct {
	Margins *Margins
	Spacing int
}

func (vb VBox) Create() (walk.Layout, error) {
	l := walk.NewVBoxLayout()

	if err := l.SetMargins(vb.Margins.toW()); err != nil {
		return nil, err
	}

	var s int
	switch vb.Spacing {
	case -1:
		s = 0

	case 0:
		s = 6

	default:
		s = vb.Spacing
	}
	if err := l.SetSpacing(s); err != nil {
		return nil, err
	}

	return l, nil
}

type Grid struct {
	Margins *Margins
	Spacing int
}

func (g Grid) Create() (walk.Layout, error) {
	l := walk.NewGridLayout()

	if err := l.SetMargins(g.Margins.toW()); err != nil {
		return nil, err
	}

	var s int
	switch g.Spacing {
	case -1:
		s = 0

	case 0:
		s = 6

	default:
		s = g.Spacing
	}
	if err := l.SetSpacing(s); err != nil {
		return nil, err
	}

	return l, nil
}
