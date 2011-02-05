// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type listViewItemChangedHandler interface {
	onListViewItemChanged(item *ListViewItem)
}

type ListViewItem struct {
	texts           []string
	userData        interface{}
	changedHandlers []listViewItemChangedHandler
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

func (lvi *ListViewItem) UserData() interface{} {
	return lvi.userData
}

func (lvi *ListViewItem) SetUserData(value interface{}) {
	lvi.userData = value
}

func (lvi *ListViewItem) addChangedHandler(handler listViewItemChangedHandler) {
	lvi.changedHandlers = append(lvi.changedHandlers, handler)
}

func (lvi *ListViewItem) removeChangedHandler(handler listViewItemChangedHandler) {
	for i, h := range lvi.changedHandlers {
		if h == handler {
			lvi.changedHandlers = append(lvi.changedHandlers[:i], lvi.changedHandlers[i+1:]...)
			break
		}
	}
}

func (lvi *ListViewItem) raiseChanged() {
	for _, handler := range lvi.changedHandlers {
		handler.onListViewItemChanged(lvi)
	}
}
