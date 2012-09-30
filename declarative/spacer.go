// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type HSpacer struct {
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Size          int
}

func (hs HSpacer) Create(parent walk.Container) (err error) {
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

	return InitWidget(hs, w, nil)
}

func (hs HSpacer) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return hs.Name, hs.StretchFactor, hs.Row, hs.RowSpan, hs.Column, hs.ColumnSpan, nil
}

type VSpacer struct {
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	Size          int
}

func (vs VSpacer) Create(parent walk.Container) (err error) {
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

	return InitWidget(vs, w, nil)
}

func (vs VSpacer) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int, contextMenu *Menu) {
	return vs.Name, vs.StretchFactor, vs.Row, vs.RowSpan, vs.Column, vs.ColumnSpan, nil
}
