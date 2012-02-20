// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type EventHandler func()

type Event struct {
	handlers []EventHandler
}

func (e *Event) Attach(handler EventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *Event) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
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
