// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type HSpacer struct {
	Name          string
	MinSize       Size
	MaxSize       Size
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

func (hs HSpacer) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return hs.Name, false, false, nil, hs.MinSize, hs.MaxSize, hs.StretchFactor, hs.Row, hs.RowSpan, hs.Column, hs.ColumnSpan, nil
}

type VSpacer struct {
	Name          string
	MinSize       Size
	MaxSize       Size
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

func (vs VSpacer) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return vs.Name, false, false, nil, vs.MinSize, vs.MaxSize, vs.StretchFactor, vs.Row, vs.RowSpan, vs.Column, vs.ColumnSpan, nil
}
