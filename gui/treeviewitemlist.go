// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

type treeViewItemListObserver interface {
	onInsertingTreeViewItem(parent *TreeViewItem, index int, item *TreeViewItem) (err os.Error)
	onRemovingTreeViewItem(index int, item *TreeViewItem) (err os.Error)
	onClearingTreeViewItems(parent *TreeViewItem) (err os.Error)
}

type TreeViewItemList struct {
	items    []*TreeViewItem
	observer treeViewItemListObserver
	parent   *TreeViewItem
}

func newTreeViewItemList(observer treeViewItemListObserver) *TreeViewItemList {
	return &TreeViewItemList{observer: observer}
}

func (l *TreeViewItemList) Add(item *TreeViewItem) (index int, err os.Error) {
	index = len(l.items)
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *TreeViewItemList) At(index int) *TreeViewItem {
	return l.items[index]
}

func (l *TreeViewItemList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingTreeViewItems(l.parent)
		if err != nil {
			return
		}
	}

	l.items = l.items[:0]

	return
}

func (l *TreeViewItemList) IndexOf(item *TreeViewItem) int {
	for i, tvi := range l.items {
		if tvi == item {
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

	item.parent = l.parent

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item

	return
}

func (l *TreeViewItemList) Len() int {
	return len(l.items)
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
		item := l.items[index]
		err = observer.onRemovingTreeViewItem(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	return
}
