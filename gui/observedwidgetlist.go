// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

type widgetListObserver interface {
	onInsertingWidget(index int, widget IWidget) (err os.Error)
	onInsertedWidget(index int, widget IWidget) (err os.Error)
	onRemovingWidget(index int, widget IWidget) (err os.Error)
	onRemovedWidget(index int, widget IWidget) (err os.Error)
	onClearingWidgets() (err os.Error)
	onClearedWidgets() (err os.Error)
}

type ObservedWidgetList struct {
	items    []IWidget
	observer widgetListObserver
}

func newObservedWidgetList(observer widgetListObserver) *ObservedWidgetList {
	return &ObservedWidgetList{observer: observer}
}

func (l *ObservedWidgetList) Add(item IWidget) (index int, err os.Error) {
	index = len(l.items)
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *ObservedWidgetList) At(index int) IWidget {
	return l.items[index]
}

func (l *ObservedWidgetList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingWidgets()
		if err != nil {
			return
		}
	}

	oldItems := l.items
	l.items = l.items[:0]

	if observer != nil {
		err = observer.onClearedWidgets()
		if err != nil {
			l.items = oldItems
			return
		}
	}

	return
}

func (l *ObservedWidgetList) IndexOf(item IWidget) int {
	for i, lvi := range l.items {
		if lvi == item {
			return i
		}
	}

	return -1
}

func (l *ObservedWidgetList) Contains(item IWidget) bool {
	return l.IndexOf(item) > -1
}

func (l *ObservedWidgetList) IndexOfHandle(handle HWND) int {
	for i, lvi := range l.items {
		if lvi.Handle() == handle {
			return i
		}
	}

	return -1
}

func (l *ObservedWidgetList) ContainsHandle(handle HWND) bool {
	return l.IndexOfHandle(handle) > -1
}

func (l *ObservedWidgetList) Insert(index int, item IWidget) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingWidget(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(append(l.items[:index], item), l.items[index:]...)

	if observer != nil {
		err = observer.onInsertedWidget(index, item)
		if err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return
		}
	}

	return
}

func (l *ObservedWidgetList) Len() int {
	return len(l.items)
}

func (l *ObservedWidgetList) Remove(item IWidget) (err os.Error) {
	index := l.IndexOf(item)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *ObservedWidgetList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	item := l.items[index]
	if observer != nil {
		err = observer.onRemovingWidget(index, item)
		if err != nil {
			return
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		err = observer.onRemovedWidget(index, item)
		if err != nil {
			l.items = append(append(l.items[:index], item), l.items[index:]...)
			return
		}
	}

	return
}
