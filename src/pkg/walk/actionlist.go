// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

type actionListObserver interface {
	onInsertingAction(index int, action *Action) os.Error
	onRemovingAction(index int, action *Action) os.Error
	onClearingActions() os.Error
}

type ActionList struct {
	actions  []*Action
	observer actionListObserver
}

func newActionList(observer actionListObserver) *ActionList {
	return &ActionList{observer: observer}
}

func (l *ActionList) Add(action *Action) os.Error {
	return l.Insert(len(l.actions), action)
}

func (l *ActionList) AddMenu(menu *Menu) (*Action, os.Error) {
	return l.InsertMenu(len(l.actions), menu)
}

func (l *ActionList) At(index int) *Action {
	return l.actions[index]
}

func (l *ActionList) Clear() os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingActions(); err != nil {
			return err
		}
	}

	l.actions = l.actions[:0]

	return nil
}

func (l *ActionList) Index(action *Action) int {
	for i, a := range l.actions {
		if a == action {
			return i
		}
	}

	return -1
}

func (l *ActionList) Insert(index int, action *Action) os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingAction(index, action); err != nil {
			return err
		}
	}

	l.actions = append(l.actions, nil)
	copy(l.actions[index+1:], l.actions[index:])
	l.actions[index] = action

	return nil
}

func (l *ActionList) InsertMenu(index int, menu *Menu) (*Action, os.Error) {
	action := NewAction()
	action.menu = menu

	if err := l.Insert(index, action); err != nil {
		return nil, err
	}

	return action, nil
}

func (l *ActionList) Len() int {
	return len(l.actions)
}

func (l *ActionList) Remove(action *Action) os.Error {
	index := l.Index(action)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *ActionList) RemoveAt(index int) os.Error {
	observer := l.observer
	if observer != nil {
		action := l.actions[index]
		if err := observer.onRemovingAction(index, action); err != nil {
			return err
		}
	}

	l.actions = append(l.actions[:index], l.actions[index+1:]...)

	return nil
}
