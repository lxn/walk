// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type ResolutionType int16

const (
	ResCustom ResolutionType = 0
	ResDraft  ResolutionType = DMRES_DRAFT
	ResLow    ResolutionType = DMRES_LOW
	ResMedium ResolutionType = DMRES_MEDIUM
	ResHigh   ResolutionType = DMRES_HIGH
)

type Resolution struct {
	typ ResolutionType
	x   int
	y   int
}

func (r *Resolution) Type() ResolutionType {
	return r.typ
}

func (r *Resolution) X() int {
	return r.x
}

func (r *Resolution) Y() int {
	return r.y
}
