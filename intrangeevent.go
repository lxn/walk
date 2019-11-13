// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type intRangeEventHandlerInfo struct {
	handler IntRangeEventHandler
	once    bool
}

type IntRangeEventHandler func(from, to int)

type IntRangeEvent struct {
	handlers []intRangeEventHandlerInfo
}

func (e *IntRangeEvent) Attach(handler IntRangeEventHandler) int {
	handlerInfo := intRangeEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *IntRangeEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *IntRangeEvent) Once(handler IntRangeEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type IntRangeEventPublisher struct {
	event IntRangeEvent
}

func (p *IntRangeEventPublisher) Event() *IntRangeEvent {
	return &p.event
}

func (p *IntRangeEventPublisher) Publish(from, to int) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(from, to)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
