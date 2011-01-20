// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type CancelEventArgs struct {
	EventArgs
	canceled bool
}

func NewCancelEventArgs(sender interface{}) *CancelEventArgs {
	return &CancelEventArgs{
		EventArgs: EventArgs{
			sender: sender,
		},
	}
}

func (a *CancelEventArgs) Canceled() bool {
	return a.canceled
}

func (a *CancelEventArgs) SetCanceled(value bool) {
	a.canceled = value
}

type CancelEventHandler func(args *CancelEventArgs)

type CancelEvent struct {
	handlers []CancelEventHandler
}

func (e *CancelEvent) Subscribe(handler CancelEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *CancelEvent) Unsubscribe(handler CancelEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type CancelEventPublisher struct {
	event CancelEvent
}

func (p *CancelEventPublisher) Event() *CancelEvent {
	return &p.event
}

func (p *CancelEventPublisher) Publish(args *CancelEventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
