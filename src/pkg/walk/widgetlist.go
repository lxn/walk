// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import . "walk/winapi"

type widgetListObserver interface {
	onInsertingWidget(index int, widget Widget) os.Error
	onInsertedWidget(index int, widget Widget) os.Error
	onRemovingWidget(index int, widget Widget) os.Error
	onRemovedWidget(index int, widget Widget) os.Error
	onClearingWidgets() os.Error
	onClearedWidgets() os.Error
}

type WidgetList struct {
	items    []Widget
	observer widgetListObserver
}

func newWidgetList(observer widgetListObserver) *WidgetList {
	return &WidgetList{observer: observer}
}

func (l *WidgetList) Add(item Widget) os.Error {
	return l.Insert(len(l.items), item)
}

func (l *WidgetList) At(index int) Widget {
	return l.items[index]
}

func (l *WidgetList) Clear() os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingWidgets(); err != nil {
			return err
		}
	}

	oldItems := l.items
	l.items = l.items[:0]

	if observer != nil {
		if err := observer.onClearedWidgets(); err != nil {
			l.items = oldItems
			return err
		}
	}

	return nil
}

func (l *WidgetList) Index(item Widget) int {
	for i, widget := range l.items {
		if widget == item {
			return i
		}
	}

	return -1
}

func (l *WidgetList) Contains(item Widget) bool {
	return l.Index(item) > -1
}

func (l *WidgetList) indexHandle(handle HWND) int {
	for i, widget := range l.items {
		if widget.BaseWidget().hWnd == handle {
			return i
		}
	}

	return -1
}

func (l *WidgetList) containsHandle(handle HWND) bool {
	return l.indexHandle(handle) > -1
}

func (l *WidgetList) insertIntoSlice(index int, item Widget) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
}

func (l *WidgetList) Insert(index int, item Widget) os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingWidget(index, item); err != nil {
			return err
		}
	}

	l.insertIntoSlice(index, item)

	if observer != nil {
		if err := observer.onInsertedWidget(index, item); err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return err
		}
	}

	return nil
}

func (l *WidgetList) Len() int {
	return len(l.items)
}

func (l *WidgetList) Remove(item Widget) os.Error {
	index := l.Index(item)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *WidgetList) RemoveAt(index int) os.Error {
	observer := l.observer
	item := l.items[index]
	if observer != nil {
		if err := observer.onRemovingWidget(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		if err := observer.onRemovedWidget(index, item); err != nil {
			l.insertIntoSlice(index, item)
			return err
		}
	}

	return nil
}
