// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type CloseEventHandler func(canceled *bool, reason CloseReason)

type CloseEvent struct {
	handlers []CloseEventHandler
}

func (e *CloseEvent) Attach(handler CloseEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *CloseEvent) Detach(handler CloseEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type CloseEventPublisher struct {
	event CloseEvent
}

func (p *CloseEventPublisher) Event() *CloseEvent {
	return &p.event
}

func (p *CloseEventPublisher) Publish(canceled *bool, reason CloseReason) {
	for _, handler := range p.event.handlers {
		handler(canceled, reason)
	}
}
