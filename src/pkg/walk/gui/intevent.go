// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type IntEventHandler func(n int)

type IntEvent struct {
	handlers []IntEventHandler
}

func (e *IntEvent) Attach(handler IntEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *IntEvent) Detach(handler IntEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type IntEventPublisher struct {
	event IntEvent
}

func (p *IntEventPublisher) Event() *IntEvent {
	return &p.event
}

func (p *IntEventPublisher) Publish(n int) {
	for _, handler := range p.event.handlers {
		handler(n)
	}
}
