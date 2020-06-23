// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type closeEventHandlerInfo struct {
	handler CloseEventHandler
	once    bool
}

type CloseEventHandler func(canceled *bool, reason CloseReason)

type CloseEvent struct {
	handlers []closeEventHandlerInfo
}

func (e *CloseEvent) Attach(handler CloseEventHandler) int {
	handlerInfo := closeEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *CloseEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *CloseEvent) Once(handler CloseEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type CloseEventPublisher struct {
	event CloseEvent
}

func (p *CloseEventPublisher) Event() *CloseEvent {
	return &p.event
}

func (p *CloseEventPublisher) Publish(canceled *bool, reason CloseReason) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(canceled, reason)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
