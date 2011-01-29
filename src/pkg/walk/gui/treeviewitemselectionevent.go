// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type TreeViewItemSelectionEventArgs struct {
	EventArgs
	old *TreeViewItem
	new *TreeViewItem
}

func NewTreeViewItemSelectionEventArgs(sender interface{}, old, new *TreeViewItem) *TreeViewItemSelectionEventArgs {
	return &TreeViewItemSelectionEventArgs{
		EventArgs: EventArgs{
			sender: sender,
		},
		old: old,
		new: new,
	}
}

func (a *TreeViewItemSelectionEventArgs) Old() *TreeViewItem {
	return a.old
}

func (a *TreeViewItemSelectionEventArgs) New() *TreeViewItem {
	return a.new
}

type TreeViewItemSelectionEventHandler func(args *TreeViewItemSelectionEventArgs)

type TreeViewItemSelectionEvent struct {
	handlers []TreeViewItemSelectionEventHandler
}

func (e *TreeViewItemSelectionEvent) Subscribe(handler TreeViewItemSelectionEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *TreeViewItemSelectionEvent) Unsubscribe(handler TreeViewItemSelectionEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type TreeViewItemSelectionEventPublisher struct {
	event TreeViewItemSelectionEvent
}

func (p *TreeViewItemSelectionEventPublisher) Event() *TreeViewItemSelectionEvent {
	return &p.event
}

func (p *TreeViewItemSelectionEventPublisher) Publish(args *TreeViewItemSelectionEventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
