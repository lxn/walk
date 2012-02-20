// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type TreeViewItemSelectionEventHandler func(old, new *TreeViewItem)

type TreeViewItemSelectionEvent struct {
	handlers []TreeViewItemSelectionEventHandler
}

func (e *TreeViewItemSelectionEvent) Attach(handler TreeViewItemSelectionEventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *TreeViewItemSelectionEvent) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
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
