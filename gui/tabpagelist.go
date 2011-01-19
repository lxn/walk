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
	onInsertingPage(index int, page *TabPage) (err os.Error)
	onInsertedPage(index int, page *TabPage) (err os.Error)
	onRemovingPage(index int, page *TabPage) (err os.Error)
	onRemovedPage(index int, page *TabPage) (err os.Error)
	onClearingPages() (err os.Error)
	onClearedPages() (err os.Error)
}

type TabPageList struct {
	items    []*TabPage
	observer tabPageListObserver
}

func newTabPageList(observer tabPageListObserver) *TabPageList {
	return &TabPageList{observer: observer}
}

func (l *TabPageList) Add(item *TabPage) (index int, err os.Error) {
	index = len(l.items)
	return index, l.Insert(index, item)
}

func (l *TabPageList) At(index int) *TabPage {
	return l.items[index]
}

func (l *TabPageList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingPages()
		if err != nil {
			return
		}
	}

	oldItems := l.items
	l.items = l.items[:0]

	if observer != nil {
		err = observer.onClearedPages()
		if err != nil {
			l.items = oldItems
			return
		}
	}

	return
}

func (l *TabPageList) IndexOf(item *TabPage) int {
	for i, lvi := range l.items {
		if lvi == item {
			return i
		}
	}

	return -1
}

func (l *TabPageList) Contains(item *TabPage) bool {
	return l.IndexOf(item) > -1
}

func (l *TabPageList) IndexOfHandle(handle HWND) int {
	for i, lvi := range l.items {
		if lvi.Handle() == handle {
			return i
		}
	}

	return -1
}

func (l *TabPageList) ContainsHandle(handle HWND) bool {
	return l.IndexOfHandle(handle) > -1
}

func (l *TabPageList) insertIntoSlice(index int, item *TabPage) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
}

func (l *TabPageList) Insert(index int, item *TabPage) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingPage(index, item)
		if err != nil {
			return
		}
	}

	l.insertIntoSlice(index, item)

	if observer != nil {
		err = observer.onInsertedPage(index, item)
		if err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return
		}
	}

	return
}

func (l *TabPageList) Len() int {
	return len(l.items)
}

func (l *TabPageList) Remove(item *TabPage) (err os.Error) {
	index := l.IndexOf(item)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *TabPageList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	item := l.items[index]
	if observer != nil {
		err = observer.onRemovingPage(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		err = observer.onRemovedPage(index, item)
		if err != nil {
			l.insertIntoSlice(index, item)
			return
		}
	}

	return
}
