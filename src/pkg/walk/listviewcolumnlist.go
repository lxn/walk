// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

type listViewColumnListObserver interface {
	onInsertingListViewColumn(index int, column *ListViewColumn) os.Error
	onRemovingListViewColumn(index int, column *ListViewColumn) os.Error
	onClearingListViewColumns() os.Error
}

type ListViewColumnList struct {
	columns  []*ListViewColumn
	observer listViewColumnListObserver
}

func newListViewColumnList(observer listViewColumnListObserver) *ListViewColumnList {
	return &ListViewColumnList{observer: observer}
}

func (l *ListViewColumnList) Add(column *ListViewColumn) os.Error {
	return l.Insert(len(l.columns), column)
}

func (l *ListViewColumnList) At(index int) *ListViewColumn {
	return l.columns[index]
}

func (l *ListViewColumnList) Clear() os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingListViewColumns(); err != nil {
			return err
		}
	}

	l.columns = l.columns[:0]

	return nil
}

func (l *ListViewColumnList) Contains(column *ListViewColumn) bool {
	return l.Index(column) > -1
}

func (l *ListViewColumnList) Index(column *ListViewColumn) int {
	for i, c := range l.columns {
		if c == column {
			return i
		}
	}

	return -1
}

func (l *ListViewColumnList) Insert(index int, column *ListViewColumn) os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingListViewColumn(index, column); err != nil {
			return err
		}
	}

	l.columns = append(l.columns, nil)
	copy(l.columns[index+1:], l.columns[index:])
	l.columns[index] = column

	return nil
}

func (l *ListViewColumnList) Len() int {
	return len(l.columns)
}

func (l *ListViewColumnList) Remove(column *ListViewColumn) os.Error {
	index := l.Index(column)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *ListViewColumnList) RemoveAt(index int) os.Error {
	observer := l.observer
	if observer != nil {
		column := l.columns[index]
		if err := observer.onRemovingListViewColumn(index, column); err != nil {
			return err
		}
	}

	l.columns = append(l.columns[:index], l.columns[index+1:]...)

	return nil
}
