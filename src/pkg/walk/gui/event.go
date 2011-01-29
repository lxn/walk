// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type EventArgs struct {
	sender interface{}
}

func NewEventArgs(sender interface{}) *EventArgs {
	return &EventArgs{
		sender: sender,
	}
}

func (a *EventArgs) Sender() interface{} {
	return a.sender
}

type EventHandler func(args *EventArgs)

type Event struct {
	handlers []EventHandler
}

func (e *Event) Subscribe(handler EventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *Event) Unsubscribe(handler EventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type EventPublisher struct {
	event Event
}

func (p *EventPublisher) Event() *Event {
	return &p.event
}

func (p *EventPublisher) Publish(args *EventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
