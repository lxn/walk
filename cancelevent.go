// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type cancelEventHandlerInfo struct {
	handler CancelEventHandler
	once    bool
}

type CancelEventHandler func(canceled *bool)

type CancelEvent struct {
	handlers []cancelEventHandlerInfo
}

func (e *CancelEvent) Attach(handler CancelEventHandler) int {
	handlerInfo := cancelEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *CancelEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *CancelEvent) Once(handler CancelEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type CancelEventPublisher struct {
	event CancelEvent
}

func (p *CancelEventPublisher) Event() *CancelEvent {
	return &p.event
}

func (p *CancelEventPublisher) Publish(canceled *bool) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(canceled)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
