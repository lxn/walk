// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type KeyEventHandler func(key int)

type KeyEvent struct {
	handlers []KeyEventHandler
}

func (e *KeyEvent) Attach(handler KeyEventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *KeyEvent) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
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
