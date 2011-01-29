// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

type tabPageListObserver interface {
	onInsertingPage(index int, page *TabPage) os.Error
	onInsertedPage(index int, page *TabPage) os.Error
	onRemovingPage(index int, page *TabPage) os.Error
	onRemovedPage(index int, page *TabPage) os.Error
	onClearingPages() os.Error
	onClearedPages() os.Error
}

type TabPageList struct {
	items    []*TabPage
	observer tabPageListObserver
}

func newTabPageList(observer tabPageListObserver) *TabPageList {
	return &TabPageList{observer: observer}
}

func (l *TabPageList) Add(item *TabPage) os.Error {
	return l.Insert(len(l.items), item)
}

func (l *TabPageList) At(index int) *TabPage {
	return l.items[index]
}

func (l *TabPageList) Clear() os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingPages(); err != nil {
			return err
		}
	}

	oldItems := l.items
	l.items = l.items[:0]

	if observer != nil {
		if err := observer.onClearedPages(); err != nil {
			l.items = oldItems
			return err
		}
	}

	return nil
}

func (l *TabPageList) Index(item *TabPage) int {
	for i, lvi := range l.items {
		if lvi == item {
			return i
		}
	}

	return -1
}

func (l *TabPageList) Contains(item *TabPage) bool {
	return l.Index(item) > -1
}

func (l *TabPageList) IndexHandle(handle HWND) int {
	for i, lvi := range l.items {
		if lvi.Handle() == handle {
			return i
		}
	}

	return -1
}

func (l *TabPageList) ContainsHandle(handle HWND) bool {
	return l.IndexHandle(handle) > -1
}

func (l *TabPageList) insertIntoSlice(index int, item *TabPage) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
}

func (l *TabPageList) Insert(index int, item *TabPage) os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingPage(index, item); err != nil {
			return err
		}
	}

	l.insertIntoSlice(index, item)

	if observer != nil {
		if err := observer.onInsertedPage(index, item); err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return err
		}
	}

	return nil
}

func (l *TabPageList) Len() int {
	return len(l.items)
}

func (l *TabPageList) Remove(item *TabPage) os.Error {
	index := l.Index(item)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *TabPageList) RemoveAt(index int) os.Error {
	observer := l.observer
	item := l.items[index]
	if observer != nil {
		if err := observer.onRemovingPage(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		if err := observer.onRemovedPage(index, item); err != nil {
			l.insertIntoSlice(index, item)
			return err
		}
	}

	return nil
}
