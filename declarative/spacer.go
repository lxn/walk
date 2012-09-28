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

func (hs HSpacer) Create(parent walk.Container) (widget walk.Widget, err error) {
	var w *walk.Spacer
	if hs.Size > 0 {
		if w, err = walk.NewHSpacerFixed(parent, hs.Size); err != nil {
			return
		}
	} else {
		if w, err = walk.NewHSpacer(parent); err != nil {
			return
		}
	}

	if err = initWidget(hs, w); err != nil {
		return
	}

	return w, nil
}

func (hs HSpacer) LayoutParams() (hStretch, vStretch, row, rowSpan, col, colSpan int) {
	return hs.HStretch, hs.VStretch, hs.Row, hs.RowSpan, hs.Column, hs.ColumnSpan
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

func (vs VSpacer) Create(parent walk.Container) (widget walk.Widget, err error) {
	var w *walk.Spacer
	if vs.Size > 0 {
		if w, err = walk.NewVSpacerFixed(parent, vs.Size); err != nil {
			return
		}
	} else {
		if w, err = walk.NewVSpacer(parent); err != nil {
			return
		}
	}

	if err = initWidget(vs, w); err != nil {
		return
	}

	return w, nil
}

func (vs VSpacer) LayoutParams() (hStretch, vStretch, row, rowSpan, col, colSpan int) {
	return vs.HStretch, vs.VStretch, vs.Row, vs.RowSpan, vs.Column, vs.ColumnSpan
}
