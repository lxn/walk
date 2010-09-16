// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"os"
)

type treeViewItemListObserver interface {
	onInsertingTreeViewItem(parent *TreeViewItem, index int, item *TreeViewItem) (err os.Error)
	onRemovingTreeViewItem(index int, item *TreeViewItem) (err os.Error)
	onClearingTreeViewItems() (err os.Error)
}

type TreeViewItemList struct {
	items    vector.Vector
	observer treeViewItemListObserver
	parent   *TreeViewItem
}

func newTreeViewItemList(observer treeViewItemListObserver) *TreeViewItemList {
	return &TreeViewItemList{observer: observer}
}

func (l *TreeViewItemList) Add(item *TreeViewItem) (index int, err os.Error) {
	index = l.items.Len()
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *TreeViewItemList) At(index int) *TreeViewItem {
	return l.items[index].(*TreeViewItem)
}

func (l *TreeViewItemList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingTreeViewItems()
		if err != nil {
			return
		}
	}

	l.items.Resize(0, 8)

	return
}

func (l *TreeViewItemList) IndexOf(item *TreeViewItem) int {
	for i, tvi := range l.items {
		if tvi.(*TreeViewItem) == item {
			return i
		}
	}

	return -1
}

func (l *TreeViewItemList) Insert(index int, item *TreeViewItem) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingTreeViewItem(l.parent, index, item)
		if err != nil {
			return
		}
	}

	l.items.Insert(index, item)

	return
}

func (l *TreeViewItemList) Len() int {
	return l.items.Len()
}

func (l *TreeViewItemList) Remove(item *TreeViewItem) (err os.Error) {
	index := l.IndexOf(item)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *TreeViewItemList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	if observer != nil {
		item := l.items[index].(*TreeViewItem)
		err = observer.onRemovingTreeViewItem(index, item)
		if err != nil {
			return
		}
	}

	l.items.Delete(index)

	return
}
