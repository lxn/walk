// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type HSpacer struct {
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Size       int
}

func (hs HSpacer) Create(parent walk.Container) (walk.Widget, error) {
	if hs.Size > 0 {
		return walk.NewHSpacerFixed(parent, hs.Size)
	}

	return walk.NewHSpacer(parent)
}

type VSpacer struct {
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Size       int
}

func (vs VSpacer) Create(parent walk.Container) (walk.Widget, error) {
	if vs.Size > 0 {
		return walk.NewVSpacerFixed(parent, vs.Size)
	}

	return walk.NewVSpacer(parent)
}
