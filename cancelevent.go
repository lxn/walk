// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type CancelEventHandler func(canceled *bool)

type CancelEvent struct {
	handlers []CancelEventHandler
}

func (e *CancelEvent) Attach(handler CancelEventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *CancelEvent) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
}

type CancelEventPublisher struct {
	event CancelEvent
}

func (p *CancelEventPublisher) Event() *CancelEvent {
	return &p.event
}

func (p *CancelEventPublisher) Publish(canceled *bool) {
	for _, handler := range p.event.handlers {
		handler(canceled)
	}
}
