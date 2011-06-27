// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type IndexList struct {
	items []int
}

func NewIndexList(items []int) *IndexList {
	return &IndexList{items}
}

func (l *IndexList) At(index int) int {
	return l.items[index]
}

func (l *IndexList) Len() int {
	return len(l.items)
}
