// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type TreeViewItemEventHandler func(item *TreeViewItem)

type TreeViewItemEvent struct {
	handlers []TreeViewItemEventHandler
}

func (e *TreeViewItemEvent) Attach(handler TreeViewItemEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *TreeViewItemEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type TreeViewItemEventPublisher struct {
	event TreeViewItemEvent
}

func (p *TreeViewItemEventPublisher) Event() *TreeViewItemEvent {
	return &p.event
}

func (p *TreeViewItemEventPublisher) Publish(item *TreeViewItem) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(item)
		}
	}
}
