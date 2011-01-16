// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

type actionListObserver interface {
	onInsertingAction(index int, action *Action) (err os.Error)
	onRemovingAction(index int, action *Action) (err os.Error)
	onClearingActions() (err os.Error)
}

type ActionList struct {
	actions  []*Action
	observer actionListObserver
}

func newActionList(observer actionListObserver) *ActionList {
	return &ActionList{observer: observer}
}

func (l *ActionList) Add(action *Action) (index int, err os.Error) {
	index = len(l.actions)
	err = l.Insert(index, action)
	if err != nil {
		return
	}

	return
}

func (l *ActionList) AddMenu(menu *Menu) (index int, action *Action, err os.Error) {
	index = len(l.actions)
	action, err = l.InsertMenu(index, menu)
	if err != nil {
		return
	}

	return
}

func (l *ActionList) At(index int) *Action {
	return l.actions[index]
}

func (l *ActionList) Clear() (err os.Error) {
	observer := l.observer
	if observer != nil {
		err = observer.onClearingActions()
		if err != nil {
			return
		}
	}

	l.actions = l.actions[:0]

	return
}

func (l *ActionList) IndexOf(action *Action) int {
	for i, a := range l.actions {
		if a == action {
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

	l.actions = append(l.actions, nil)
	copy(l.actions[index+1:], l.actions[index:])
	l.actions[index] = action

	return
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
		action := l.actions[index]
		err = observer.onRemovingAction(index, action)
		if err != nil {
			return
		}
	}

	l.actions = append(l.actions[:index], l.actions[index+1:]...)

	return
}
