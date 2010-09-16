// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
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
	items    vector.Vector
	observer widgetListObserver
}

func newObservedWidgetList(observer widgetListObserver) *ObservedWidgetList {
	return &ObservedWidgetList{observer: observer}
}

func (l *ObservedWidgetList) Add(item IWidget) (index int, err os.Error) {
	index = l.items.Len()
	err = l.Insert(index, item)
	if err != nil {
		return
	}

	return
}

func (l *ObservedWidgetList) At(index int) IWidget {
	return l.items[index].(IWidget)
}

func (l *ObservedWidgetList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingWidgets()
		if err != nil {
			return
		}
	}

	oldLen := l.items.Len()
	l.items = vector.Vector(l.items[0:0])

	if observer != nil {
		err = observer.onClearedWidgets()
		if err != nil {
			l.items = vector.Vector(l.items[0:oldLen])
			return
		}
	}

	l.items.Resize(0, 8)

	return
}

func (l *ObservedWidgetList) IndexOf(item IWidget) int {
	for i, lvi := range l.items {
		if lvi.(IWidget) == item {
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
		if lvi.(IWidget).Handle() == handle {
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

	l.items.Insert(index, item)

	if observer != nil {
		err = observer.onInsertedWidget(index, item)
		if err != nil {
			l.items.Delete(index)
			return
		}
	}

	return
}

func (l *ObservedWidgetList) Len() int {
	return l.items.Len()
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
	item := l.items[index].(IWidget)
	if observer != nil {
		err = observer.onRemovingWidget(index, item)
		if err != nil {
			return
		}
	}

	l.items.Delete(index)

	if observer != nil {
		err = observer.onRemovedWidget(index, item)
		if err != nil {
			l.items.Insert(index, item)
			return
		}
	}

	return
}
