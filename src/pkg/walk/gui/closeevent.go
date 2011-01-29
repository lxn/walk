// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type CloseEventArgs struct {
	CancelEventArgs
	reason CloseReason
}

func NewCloseEventArgs(sender interface{}, reason CloseReason) *CloseEventArgs {
	return &CloseEventArgs{
		CancelEventArgs: CancelEventArgs{
			EventArgs: EventArgs{
				sender: sender,
			},
		},
		reason: reason,
	}
}

func (a *CloseEventArgs) Reason() CloseReason {
	return a.reason
}

type CloseEventHandler func(args *CloseEventArgs)

type CloseEvent struct {
	handlers []CloseEventHandler
}

func (e *CloseEvent) Subscribe(handler CloseEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *CloseEvent) Unsubscribe(handler CloseEventHandler) {
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

func (p *CloseEventPublisher) Publish(args *CloseEventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
