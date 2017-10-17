// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type HSpacer struct {
	// Window

	MaxSize Size
	MinSize Size
	Name    string

	// Widget

	Column        int
	ColumnSpan    int
	Row           int
	RowSpan       int
	StretchFactor int

	// Spacer

	Size int
}

func (hs HSpacer) Create(builder *Builder) (err error) {
	var w *walk.Spacer
	if hs.Size > 0 {
		if w, err = walk.NewHSpacerFixed(builder.Parent(), hs.Size); err != nil {
			return
		}
	} else {
		if w, err = walk.NewHSpacer(builder.Parent()); err != nil {
			return
		}
	}

	return builder.InitWidget(hs, w, nil)
}

type VSpacer struct {
	// Window

	MaxSize Size
	MinSize Size
	Name    string

	// Widget

	Column        int
	ColumnSpan    int
	Row           int
	RowSpan       int
	StretchFactor int

	// Spacer

	Size int
}

func (vs VSpacer) Create(builder *Builder) (err error) {
	var w *walk.Spacer
	if vs.Size > 0 {
		if w, err = walk.NewVSpacerFixed(builder.Parent(), vs.Size); err != nil {
			return
		}
	} else {
		if w, err = walk.NewVSpacer(builder.Parent()); err != nil {
			return
		}
	}

	return builder.InitWidget(vs, w, nil)
}
