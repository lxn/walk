// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
)

type listViewItemChangedHandler interface {
	onListViewItemChanged(item *ListViewItem)
}

type ListViewItem struct {
	texts           []string
	changedHandlers vector.Vector
}

func NewListViewItem() *ListViewItem {
	return &ListViewItem{}
}

func (lvi *ListViewItem) Texts() []string {
	return lvi.texts
}

func (lvi *ListViewItem) SetTexts(value []string) {
	lvi.texts = value

	lvi.raiseChanged()
}

func (lvi *ListViewItem) addChangedHandler(handler listViewItemChangedHandler) {
	lvi.changedHandlers.Push(handler)
}

func (lvi *ListViewItem) removeChangedHandler(handler listViewItemChangedHandler) {
	for i, h := range lvi.changedHandlers {
		if h.(listViewItemChangedHandler) == handler {
			lvi.changedHandlers.Delete(i)
			break
		}
	}
}

func (lvi *ListViewItem) raiseChanged() {
	for _, handlerIface := range lvi.changedHandlers {
		handler := handlerIface.(listViewItemChangedHandler)
		handler.onListViewItemChanged(lvi)
	}
}
