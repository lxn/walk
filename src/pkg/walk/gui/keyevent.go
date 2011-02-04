// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type KeyEventHandler func(key int)

type KeyEvent struct {
	handlers []KeyEventHandler
}

func (e *KeyEvent) Attach(handler KeyEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *KeyEvent) Detach(handler KeyEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type KeyEventPublisher struct {
	event KeyEvent
}

func (p *KeyEventPublisher) Event() *KeyEvent {
	return &p.event
}

func (p *KeyEventPublisher) Publish(key int) {
	for _, handler := range p.event.handlers {
		handler(key)
	}
}
