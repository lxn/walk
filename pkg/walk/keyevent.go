// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type keyEventHandlerInfo struct {
	handler KeyEventHandler
	once    bool
}

type KeyEventHandler func(key Key)

type KeyEvent struct {
	handlers []keyEventHandlerInfo
}

func (e *KeyEvent) Attach(handler KeyEventHandler) int {
	handlerInfo := keyEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *KeyEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *KeyEvent) Once(handler KeyEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type KeyEventPublisher struct {
	event KeyEvent
}

func (p *KeyEventPublisher) Event() *KeyEvent {
	return &p.event
}

func (p *KeyEventPublisher) Publish(key Key) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(key)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
