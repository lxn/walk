// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

type listViewColumnListObserver interface {
	onInsertingListViewColumn(index int, column *ListViewColumn) (err os.Error)
	onRemovingListViewColumn(index int, column *ListViewColumn) (err os.Error)
	onClearingListViewColumns() (err os.Error)
}

type ListViewColumnList struct {
	columns  []*ListViewColumn
	observer listViewColumnListObserver
}

func newListViewColumnList(observer listViewColumnListObserver) *ListViewColumnList {
	return &ListViewColumnList{observer: observer}
}

func (l *ListViewColumnList) Add(column *ListViewColumn) (index int, err os.Error) {
	index = len(l.columns)
	err = l.Insert(index, column)
	if err != nil {
		return
	}

	return
}

func (l *ListViewColumnList) At(index int) *ListViewColumn {
	return l.columns[index]
}

func (l *ListViewColumnList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingListViewColumns()
		if err != nil {
			return
		}
	}

	l.columns = l.columns[:0]

	return
}

func (l *ListViewColumnList) IndexOf(column *ListViewColumn) int {
	for i, c := range l.columns {
		if c == column {
			return i
		}
	}

	return -1
}

func (l *ListViewColumnList) Insert(index int, column *ListViewColumn) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingListViewColumn(index, column)
		if err != nil {
			return
		}
	}

	l.columns = append(append(l.columns[:index], column), l.columns[index:]...)

	return
}

func (l *ListViewColumnList) Len() int {
	return len(l.columns)
}

func (l *ListViewColumnList) Remove(column *ListViewColumn) (err os.Error) {
	index := l.IndexOf(column)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *ListViewColumnList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	if observer != nil {
		column := l.columns[index]
		err = observer.onRemovingListViewColumn(index, column)
		if err != nil {
			return
		}
	}

	l.columns = append(l.columns[:index], l.columns[index+1:]...)

	return
}
