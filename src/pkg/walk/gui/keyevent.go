// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type KeyEventArgs struct {
	EventArgs
	key int
}

func NewKeyEventArgs(sender interface{}, key int) *KeyEventArgs {
	return &KeyEventArgs{
		EventArgs: EventArgs{
			sender: sender,
		},
		key: key,
	}
}

func (a *KeyEventArgs) Key() int {
	return a.key
}

type KeyEventHandler func(args *KeyEventArgs)

type KeyEvent struct {
	handlers []KeyEventHandler
}

func (e *KeyEvent) Subscribe(handler KeyEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *KeyEvent) Unsubscribe(handler KeyEventHandler) {
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

func (p *KeyEventPublisher) Publish(args *KeyEventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
