// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

type listViewItemListObserver interface {
	onInsertingListViewItem(index int, item *ListViewItem) (err os.Error)
	onRemovingListViewItem(index int, item *ListViewItem) (err os.Error)
	onClearingListViewItems() (err os.Error)
}

type ListViewItemList struct {
	items    []*ListViewItem
	observer listViewItemListObserver
}

func newListViewItemList(observer listViewItemListObserver) *ListViewItemList {
	return &ListViewItemList{observer: observer}
}

func (l *ListViewItemList) Add(item *ListViewItem) (index int, err os.Error) {
	index = len(l.items)
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *ListViewItemList) At(index int) *ListViewItem {
	return l.items[index]
}

func (l *ListViewItemList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingListViewItems()
		if err != nil {
			return
		}
	}

	l.items = l.items[:0]

	return
}

func (l *ListViewItemList) IndexOf(item *ListViewItem) int {
	for i, lvi := range l.items {
		if lvi == item {
			return i
		}
	}

	return -1
}

func (l *ListViewItemList) Insert(index int, item *ListViewItem) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingListViewItem(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item

	return
}

func (l *ListViewItemList) Len() int {
	return len(l.items)
}

func (l *ListViewItemList) Remove(item *ListViewItem) (err os.Error) {
	index := l.IndexOf(item)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *ListViewItemList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	if observer != nil {
		item := l.items[index]
		err = observer.onRemovingListViewItem(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	return
}
