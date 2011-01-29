// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type MouseButton int

const (
	LeftButton MouseButton = iota
	RightButton
	MiddleButton
)

type MouseEventArgs struct {
	EventArgs
	x      int
	y      int
	button MouseButton
}

func NewMouseEventArgs(sender interface{}, x, y int, button MouseButton) *MouseEventArgs {
	return &MouseEventArgs{
		EventArgs: EventArgs{
			sender: sender,
		},
		x:      x,
		y:      y,
		button: button,
	}
}

func (a *MouseEventArgs) X() int {
	return a.x
}

func (a *MouseEventArgs) Y() int {
	return a.y
}

func (a *MouseEventArgs) Button() MouseButton {
	return a.button
}

type MouseEventHandler func(args *MouseEventArgs)

type MouseEvent struct {
	handlers []MouseEventHandler
}

func (e *MouseEvent) Subscribe(handler MouseEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *MouseEvent) Unsubscribe(handler MouseEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type MouseEventPublisher struct {
	event MouseEvent
}

func (p *MouseEventPublisher) Event() *MouseEvent {
	return &p.event
}

func (p *MouseEventPublisher) Publish(args *MouseEventArgs) {
	for _, handler := range p.event.handlers {
		handler(args)
	}
}
