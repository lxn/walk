// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

type listViewItemListObserver interface {
	onInsertingListViewItem(index int, item *ListViewItem) os.Error
	onRemovingListViewItem(index int, item *ListViewItem) os.Error
	onClearingListViewItems() os.Error
}

type ListViewItemList struct {
	items    []*ListViewItem
	observer listViewItemListObserver
}

func newListViewItemList(observer listViewItemListObserver) *ListViewItemList {
	return &ListViewItemList{observer: observer}
}

func (l *ListViewItemList) Add(item *ListViewItem) os.Error {
	return l.Insert(len(l.items), item)
}

func (l *ListViewItemList) At(index int) *ListViewItem {
	return l.items[index]
}

func (l *ListViewItemList) Clear() os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingListViewItems(); err != nil {
			return err
		}
	}

	l.items = l.items[:0]

	return nil
}

func (l *ListViewItemList) Index(item *ListViewItem) int {
	for i, lvi := range l.items {
		if lvi == item {
			return i
		}
	}

	return -1
}

func (l *ListViewItemList) Insert(index int, item *ListViewItem) os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingListViewItem(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item

	return nil
}

func (l *ListViewItemList) Len() int {
	return len(l.items)
}

func (l *ListViewItemList) Remove(item *ListViewItem) os.Error {
	index := l.Index(item)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *ListViewItemList) RemoveAt(index int) os.Error {
	observer := l.observer
	if observer != nil {
		item := l.items[index]
		if err := observer.onRemovingListViewItem(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	return nil
}
