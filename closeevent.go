// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type CloseEventHandler func(canceled *bool, reason CloseReason)

type CloseEvent struct {
	handlers []CloseEventHandler
}

func (e *CloseEvent) Attach(handler CloseEventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *CloseEvent) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
}

type CloseEventPublisher struct {
	event CloseEvent
}

func (p *CloseEventPublisher) Event() *CloseEvent {
	return &p.event
}

func (p *CloseEventPublisher) Publish(canceled *bool, reason CloseReason) {
	for _, handler := range p.event.handlers {
		handler(canceled, reason)
	}
}
