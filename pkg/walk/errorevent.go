// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type errorEventHandlerInfo struct {
	handler ErrorEventHandler
	once    bool
}

type ErrorEventHandler func(err error)

type ErrorEvent struct {
	handlers []errorEventHandlerInfo
}

func (e *ErrorEvent) Attach(handler ErrorEventHandler) int {
	handlerInfo := errorEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *ErrorEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *ErrorEvent) Once(handler ErrorEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type ErrorEventPublisher struct {
	event ErrorEvent
}

func (p *ErrorEventPublisher) Event() *ErrorEvent {
	return &p.event
}

func (p *ErrorEventPublisher) Publish(err error) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(err)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
