// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Margins struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

type HBox struct {
	Margins *Margins
	Spacing int
}

func (hb HBox) Create() (walk.Layout, error) {
	l := walk.NewHBoxLayout()

	var m walk.Margins
	if hb.Margins == nil {
		m = walk.Margins{9, 9, 9, 9}
	} else {
		hbm := hb.Margins

		m = walk.Margins{hbm.Left, hbm.Top, hbm.Right, hbm.Bottom}
	}

	if err := l.SetMargins(m); err != nil {
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

	var m walk.Margins
	if vb.Margins == nil {
		m = walk.Margins{9, 9, 9, 9}
	} else {
		vbm := vb.Margins

		m = walk.Margins{vbm.Left, vbm.Top, vbm.Right, vbm.Bottom}
	}

	if err := l.SetMargins(m); err != nil {
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

	var m walk.Margins
	if g.Margins == nil {
		m = walk.Margins{9, 9, 9, 9}
	} else {
		gm := g.Margins

		m = walk.Margins{gm.Left, gm.Top, gm.Right, gm.Bottom}
	}

	if err := l.SetMargins(m); err != nil {
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
