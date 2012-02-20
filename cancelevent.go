// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type CancelEventHandler func(canceled *bool)

type CancelEvent struct {
	handlers []CancelEventHandler
}

func (e *CancelEvent) Attach(handler CancelEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *CancelEvent) Detach(handler CancelEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
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
