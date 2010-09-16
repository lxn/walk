// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"os"
)

type comboBoxItemListObserver interface {
	onInsertingComboBoxItem(index int, item *ComboBoxItem) (err os.Error)
	onRemovingComboBoxItem(index int, item *ComboBoxItem) (err os.Error)
	onClearingComboBoxItems() (err os.Error)
}

type ComboBoxItemList struct {
	items    vector.Vector
	observer comboBoxItemListObserver
}

func newComboBoxItemList(observer comboBoxItemListObserver) *ComboBoxItemList {
	return &ComboBoxItemList{observer: observer}
}

func (l *ComboBoxItemList) Add(item *ComboBoxItem) (index int, err os.Error) {
	index = l.items.Len()
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *ComboBoxItemList) At(index int) *ComboBoxItem {
	return l.items[index].(*ComboBoxItem)
}

func (l *ComboBoxItemList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingComboBoxItems()
		if err != nil {
			return
		}
	}

	l.items.Resize(0, 8)

	return
}

func (l *ComboBoxItemList) IndexOf(item *ComboBoxItem) int {
	for i, tvi := range l.items {
		if tvi.(*ComboBoxItem) == item {
			return i
		}
	}

	return -1
}

func (l *ComboBoxItemList) Insert(index int, item *ComboBoxItem) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingComboBoxItem(index, item)
		if err != nil {
			return
		}
	}

	l.items.Insert(index, item)

	return
}

func (l *ComboBoxItemList) Len() int {
	return l.items.Len()
}

func (l *ComboBoxItemList) Remove(item *ComboBoxItem) (err os.Error) {
	index := l.IndexOf(item)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *ComboBoxItemList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	if observer != nil {
		item := l.items[index].(*ComboBoxItem)
		err = observer.onRemovingComboBoxItem(index, item)
		if err != nil {
			return
		}
	}

	l.items.Delete(index)

	return
}
