// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type TreeViewItemEventArgs struct {
	EventArgs
	item *TreeViewItem
}

func NewTreeViewItemEventArgs(sender interface{}, item *TreeViewItem) *TreeViewItemEventArgs {
	return &TreeViewItemEventArgs{
		EventArgs: EventArgs{
			sender: sender,
		},
		item: item,
	}
}

func (a *TreeViewItemEventArgs) Item() *TreeViewItem {
	return a.item
}

type TreeViewItemEventHandler func(args *TreeViewItemEventArgs)

type TreeViewItemEvent struct {
	handlers []TreeViewItemEventHandler
}

func (e *TreeViewItemEvent) Subscribe(handler TreeViewItemEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *TreeViewItemEvent) Unsubscribe(handler TreeViewItemEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type TreeViewItemEventPublisher struct {
	event TreeViewItemEvent
}

func (p *TreeViewItemEventPublisher) Event() *TreeViewItemEvent {
	return &p.event
}

func (p *TreeViewItemEventPublisher) Publish(args *TreeViewItemEventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
