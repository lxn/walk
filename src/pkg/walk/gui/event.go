// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type EventHandler func()

type Event struct {
	handlers []EventHandler
}

func (e *Event) Attach(handler EventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *Event) Detach(handler EventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type EventPublisher struct {
	event Event
}

func (p *EventPublisher) Event() *Event {
	return &p.event
}

func (p *EventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		handler()
	}
}
