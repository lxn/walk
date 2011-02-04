// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type TreeViewItemSelectionEventHandler func(old, new *TreeViewItem)

type TreeViewItemSelectionEvent struct {
	handlers []TreeViewItemSelectionEventHandler
}

func (e *TreeViewItemSelectionEvent) Attach(handler TreeViewItemSelectionEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *TreeViewItemSelectionEvent) Detach(handler TreeViewItemSelectionEventHandler) {
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

func (p *TreeViewItemSelectionEventPublisher) Publish(old, new *TreeViewItem) {
	for _, handler := range p.event.handlers {
		handler(old, new)
	}
}
