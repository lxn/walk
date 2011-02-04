// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type TreeViewItemEventHandler func(item *TreeViewItem)

type TreeViewItemEvent struct {
	handlers []TreeViewItemEventHandler
}

func (e *TreeViewItemEvent) Attach(handler TreeViewItemEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *TreeViewItemEvent) Detach(handler TreeViewItemEventHandler) {
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

func (p *TreeViewItemEventPublisher) Publish(item *TreeViewItem) {
	for _, handler := range p.event.handlers {
		handler(item)
	}
}
