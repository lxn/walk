// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type stringEventHandlerInfo struct {
	handler StringEventHandler
	once    bool
}

type StringEventHandler func(s string)

type StringEvent struct {
	handlers []stringEventHandlerInfo
}

func (e *StringEvent) Attach(handler StringEventHandler) int {
	handlerInfo := stringEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *StringEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *StringEvent) Once(handler StringEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type StringEventPublisher struct {
	event StringEvent
}

func (p *StringEventPublisher) Event() *StringEvent {
	return &p.event
}

func (p *StringEventPublisher) Publish(s string) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(s)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
