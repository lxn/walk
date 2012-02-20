// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type ErrorEventHandler func(err error)

type ErrorEvent struct {
	handlers []ErrorEventHandler
}

func (e *ErrorEvent) Attach(handler ErrorEventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *ErrorEvent) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
}

type ErrorEventPublisher struct {
	event ErrorEvent
}

func (p *ErrorEventPublisher) Event() *ErrorEvent {
	return &p.event
}

func (p *ErrorEventPublisher) Publish(err error) {
	for _, handler := range p.event.handlers {
		handler(err)
	}
}
