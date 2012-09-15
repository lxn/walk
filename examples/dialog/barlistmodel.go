// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
)

var barModel *BarListModel = NewPopulatedBarListModel()

type BarListModel struct {
	walk.ListModelBase
	items []*Bar
}

func NewPopulatedBarListModel() *BarListModel {
	return &BarListModel{
		items: []*Bar{
			{"Some", "S"},
			{"Example", "E"},
			{"Items", "I"},
		},
	}
}

func (blm *BarListModel) Index(bar *Bar) int {
	for i, b := range blm.items {
		if b == bar {
			return i
		}
	}

	return -1
}

func (blm *BarListModel) ItemCount() int {
	return len(blm.items)
}

func (blm *BarListModel) Value(index int) interface{} {
	return blm.items[index].Name
}
