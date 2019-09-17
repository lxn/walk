// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type EventHandler func()

type Event struct {
	handlers []EventHandler
}

func (e *Event) Attach(handler EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *Event) Detach(handle int) {
	e.handlers[handle] = nil
}

type EventPublisher struct {
	event Event
}

func (p *EventPublisher) Event() *Event {
	return &p.event
}

func (p *EventPublisher) Publish() {
	// This is a kludge to find the form that the event publisher is
	// affiliated with. It's only necessary because the event publisher
	// doesn't keep a pointer to the form on its own, and the call
	// to Publish isn't providing it either.
	if form := App().ActiveForm(); form != nil {
		fb := form.AsFormBase()
		fb.inProgressEventCount++
		defer func() {
			fb.inProgressEventCount--
			if fb.inProgressEventCount == 0 && fb.layoutScheduled {
				fb.layoutScheduled = false
				fb.startLayout()
			}
		}()
	}

	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}
