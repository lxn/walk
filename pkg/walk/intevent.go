// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type intEventHandlerInfo struct {
	handler IntEventHandler
	once    bool
}

type IntEventHandler func(n int)

type IntEvent struct {
	handlers []intEventHandlerInfo
}

func (e *IntEvent) Attach(handler IntEventHandler) int {
	handlerInfo := intEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *IntEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *IntEvent) Once(handler IntEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type IntEventPublisher struct {
	event IntEvent
}

func (p *IntEventPublisher) Event() *IntEvent {
	return &p.event
}

func (p *IntEventPublisher) Publish(n int) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(n)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
