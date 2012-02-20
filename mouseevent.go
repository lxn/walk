// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type MouseButton int

const (
	LeftButton MouseButton = iota
	RightButton
	MiddleButton
)

type MouseEventHandler func(x, y int, button MouseButton)

type MouseEvent struct {
	handlers []MouseEventHandler
}

func (e *MouseEvent) Attach(handler MouseEventHandler) int {
	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *MouseEvent) Detach(handle int) {
	e.handlers = append(e.handlers[:handle], e.handlers[handle+1:]...)
}

type MouseEventPublisher struct {
	event MouseEvent
}

func (p *MouseEventPublisher) Event() *MouseEvent {
	return &p.event
}

func (p *MouseEventPublisher) Publish(x, y int, button MouseButton) {
	for _, handler := range p.event.handlers {
		handler(x, y, button)
	}
}
