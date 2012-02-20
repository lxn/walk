// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type treeViewItemListObserver interface {
	onInsertingTreeViewItem(parent *TreeViewItem, index int, item *TreeViewItem) error
	onRemovingTreeViewItem(index int, item *TreeViewItem) error
	onClearingTreeViewItems(parent *TreeViewItem) error
}

type TreeViewItemList struct {
	items    []*TreeViewItem
	observer treeViewItemListObserver
	parent   *TreeViewItem
}

func newTreeViewItemList(observer treeViewItemListObserver) *TreeViewItemList {
	return &TreeViewItemList{observer: observer}
}

func (l *TreeViewItemList) Add(item *TreeViewItem) error {
	return l.Insert(len(l.items), item)
}

func (l *TreeViewItemList) At(index int) *TreeViewItem {
	return l.items[index]
}

func (l *TreeViewItemList) Clear() error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingTreeViewItems(l.parent); err != nil {
			return err
		}
	}

	l.items = l.items[:0]

	return nil
}

func (l *TreeViewItemList) Contains(item *TreeViewItem) bool {
	return l.Index(item) > -1
}

func (l *TreeViewItemList) Index(item *TreeViewItem) int {
	for i, tvi := range l.items {
		if tvi == item {
			return i
		}
	}

	return -1
}

func (l *TreeViewItemList) Insert(index int, item *TreeViewItem) error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingTreeViewItem(l.parent, index, item); err != nil {
			return err
		}
	}

	item.parent = l.parent

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item

	return nil
}

func (l *TreeViewItemList) Len() int {
	return len(l.items)
}

func (l *TreeViewItemList) Remove(item *TreeViewItem) error {
	index := l.Index(item)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *TreeViewItemList) RemoveAt(index int) error {
	observer := l.observer
	if observer != nil {
		item := l.items[index]
		if err := observer.onRemovingTreeViewItem(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	return nil
}
