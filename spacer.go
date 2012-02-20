// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

const spacerWindowClass = `\o/ Walk_Spacer_Class \o/`

var spacerWindowClassRegistered bool

type Spacer struct {
	WidgetBase
	sizeHint    Size
	layoutFlags LayoutFlags
}

func newSpacer(parent Container, layoutFlags LayoutFlags, sizeHint Size) (*Spacer, os.Error) {
	ensureRegisteredWindowClass(spacerWindowClass, &spacerWindowClassRegistered)

	s := &Spacer{
		layoutFlags: layoutFlags,
		sizeHint:    sizeHint,
	}

	if err := initChildWidget(
		s,
		parent,
		spacerWindowClass,
		0,
		0); err != nil {
		return nil, err
	}

	return s, nil
}

func NewHSpacer(parent Container) (*Spacer, os.Error) {
	return newSpacer(parent, ShrinkableHorz|ShrinkableVert|GrowableHorz|GreedyHorz, Size{})
}

func NewHSpacerFixed(parent Container, width int) (*Spacer, os.Error) {
	return newSpacer(parent, 0, Size{width, 0})
}

func NewVSpacer(parent Container) (*Spacer, os.Error) {
	return newSpacer(parent, ShrinkableHorz|ShrinkableVert|GrowableVert|GreedyVert, Size{})
}

func NewVSpacerFixed(parent Container, height int) (*Spacer, os.Error) {
	return newSpacer(parent, 0, Size{0, height})
}

func (s *Spacer) LayoutFlags() LayoutFlags {
	return s.layoutFlags
}

func (s *Spacer) SizeHint() Size {
	return s.sizeHint
}
