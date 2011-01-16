// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

type comboBoxItemListObserver interface {
	onInsertingComboBoxItem(index int, item *ComboBoxItem) (err os.Error)
	onRemovingComboBoxItem(index int, item *ComboBoxItem) (err os.Error)
	onClearingComboBoxItems() (err os.Error)
}

type ComboBoxItemList struct {
	items    []*ComboBoxItem
	observer comboBoxItemListObserver
}

func newComboBoxItemList(observer comboBoxItemListObserver) *ComboBoxItemList {
	return &ComboBoxItemList{observer: observer}
}

func (l *ComboBoxItemList) Add(item *ComboBoxItem) (index int, err os.Error) {
	index = len(l.items)
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *ComboBoxItemList) At(index int) *ComboBoxItem {
	return l.items[index]
}

func (l *ComboBoxItemList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingComboBoxItems()
		if err != nil {
			return
		}
	}

	l.items = l.items[:0]

	return
}

func (l *ComboBoxItemList) IndexOf(item *ComboBoxItem) int {
	for i, tvi := range l.items {
		if tvi == item {
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

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item

	return
}

func (l *ComboBoxItemList) Len() int {
	return len(l.items)
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
		item := l.items[index]
		err = observer.onRemovingComboBoxItem(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	return
}
