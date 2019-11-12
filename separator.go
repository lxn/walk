// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type Separator struct {
	WidgetBase
	vertical bool
}

func NewHSeparator(parent Container) (*Separator, error) {
	return newSeparator(parent, false)
}

func NewVSeparator(parent Container) (*Separator, error) {
	return newSeparator(parent, true)
}

func newSeparator(parent Container, vertical bool) (*Separator, error) {
	s := &Separator{vertical: vertical}

	if err := InitWidget(
		s,
		parent,
		"STATIC",
		win.WS_VISIBLE|win.SS_ETCHEDHORZ,
		0); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Separator) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	var layoutFlags LayoutFlags
	if s.vertical {
		layoutFlags = GrowableHorz | GreedyHorz
	} else {
		layoutFlags = GrowableVert | GreedyVert
	}

	return &separatorLayoutItem{
		layoutFlags: layoutFlags,
	}
}

type separatorLayoutItem struct {
	LayoutItemBase
	layoutFlags LayoutFlags
}

func (li *separatorLayoutItem) LayoutFlags() LayoutFlags {
	return li.layoutFlags
}

func (li *separatorLayoutItem) IdealSize() Size {
	return li.MinSize()
}

func (li *separatorLayoutItem) MinSize() Size {
	return SizeFrom96DPI(Size{2, 2}, li.ctx.dpi)
}
