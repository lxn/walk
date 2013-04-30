// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Alignment1D uint

const (
	AlignNear Alignment1D = iota
	AlignCenter
	AlignFar
)

type TableViewColumn struct {
	Alignment  Alignment1D
	DataMember string
	Format     string
	Precision  int
	Title      string
	Width      int
}

func (tvc TableViewColumn) Create(tv *walk.TableView) error {
	w := walk.NewTableViewColumn()

	if err := w.SetAlignment(walk.Alignment1D(tvc.Alignment)); err != nil {
		return err
	}
	w.SetDataMember(tvc.DataMember)
	if tvc.Format != "" {
		if err := w.SetFormat(tvc.Format); err != nil {
			return err
		}
	}
	if err := w.SetPrecision(tvc.Precision); err != nil {
		return err
	}
	if err := w.SetTitle(tvc.Title); err != nil {
		return err
	}
	if err := w.SetWidth(tvc.Width); err != nil {
		return err
	}

	return tv.Columns().Add(w)
}
