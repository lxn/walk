// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type SysCommandEventHandler func(wParam uint32)

type SysCommandEvent struct {
	handlers []SysCommandEventHandler
}

func (e *SysCommandEvent) Attach(handler SysCommandEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *SysCommandEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type SysCommandEventPublisher struct {
	event SysCommandEvent
}

func (p *SysCommandEventPublisher) Event() *SysCommandEvent {
	return &p.event
}

func (p *SysCommandEventPublisher) Publish(wParam uint32) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(wParam)
		}
	}
}
