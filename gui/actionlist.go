// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"os"
)

type actionListObserver interface {
	onInsertingAction(index int, action *Action) (err os.Error)
	onRemovingAction(index int, action *Action) (err os.Error)
	onClearingActions() (err os.Error)
}

type ActionList struct {
	actions  vector.Vector
	observer actionListObserver
}

func newActionList(observer actionListObserver) *ActionList {
	return &ActionList{observer: observer}
}

func (l *ActionList) Add(action *Action) (index int, err os.Error) {
	index = l.actions.Len()
	err = l.Insert(index, action)
	if err != nil {
		return
	}

	return
}

func (l *ActionList) AddMenu(menu *Menu) (index int, action *Action, err os.Error) {
	index = l.actions.Len()
	action, err = l.InsertMenu(index, menu)
	if err != nil {
		return
	}

	return
}

func (l *ActionList) At(index int) *Action {
	return l.actions[index].(*Action)
}

func (l *ActionList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingActions()
		if err != nil {
			return
		}
	}

	l.actions.Resize(0, 8)

	return
}

func (l *ActionList) IndexOf(action *Action) int {
	for i, a := range l.actions {
		if a.(*Action) == action {
			return i
		}
	}

	return -1
}

func (l *ActionList) Insert(index int, action *Action) (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onInsertingAction(index, action)
		if err != nil {
			return
		}
	}

	l.actions.Insert(index, action)

	return
}

func (l *ActionList) InsertMenu(index int, menu *Menu) (action *Action, err os.Error) {
	action = NewAction()
	action.menu = menu

	observer := l.observer
	if observer != nil {
		err = observer.onInsertingAction(index, action)
		if err != nil {
			return
		}
	}

	l.actions.Insert(index, action)

	return
}

func (l *ActionList) Len() int {
	return l.actions.Len()
}

func (l *ActionList) Remove(action *Action) (err os.Error) {
	index := l.IndexOf(action)
	if index == -1 {
		return
	}

	return l.RemoveAt(index)
}

func (l *ActionList) RemoveAt(index int) (err os.Error) {
	observer := l.observer
	if observer != nil {
		action := l.actions[index].(*Action)
		err = observer.onRemovingAction(index, action)
		if err != nil {
			return
		}
	}

	l.actions.Delete(index)

	return
}
