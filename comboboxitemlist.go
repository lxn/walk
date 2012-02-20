// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type comboBoxItemListObserver interface {
	onInsertingComboBoxItem(index int, item *ComboBoxItem) error
	onRemovingComboBoxItem(index int, item *ComboBoxItem) error
	onClearingComboBoxItems() error
}

type ComboBoxItemList struct {
	items    []*ComboBoxItem
	observer comboBoxItemListObserver
}

func newComboBoxItemList(observer comboBoxItemListObserver) *ComboBoxItemList {
	return &ComboBoxItemList{observer: observer}
}

func (l *ComboBoxItemList) Add(item *ComboBoxItem) error {
	return l.Insert(len(l.items), item)
}

func (l *ComboBoxItemList) At(index int) *ComboBoxItem {
	return l.items[index]
}

func (l *ComboBoxItemList) Clear() error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingComboBoxItems(); err != nil {
			return err
		}
	}

	l.items = l.items[:0]

	return nil
}

func (l *ComboBoxItemList) Contains(item *ComboBoxItem) bool {
	return l.Index(item) > -1
}

func (l *ComboBoxItemList) Index(item *ComboBoxItem) int {
	for i, tvi := range l.items {
		if tvi == item {
			return i
		}
	}

	return -1
}

func (l *ComboBoxItemList) Insert(index int, item *ComboBoxItem) error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingComboBoxItem(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item

	return nil
}

func (l *ComboBoxItemList) Len() int {
	return len(l.items)
}

func (l *ComboBoxItemList) Remove(item *ComboBoxItem) error {
	index := l.Index(item)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *ComboBoxItemList) RemoveAt(index int) error {
	observer := l.observer
	if observer != nil {
		item := l.items[index]
		if err := observer.onRemovingComboBoxItem(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	return nil
}
