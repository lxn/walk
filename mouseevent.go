// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type MouseButton int

const (
	LeftButton   MouseButton = win.MK_LBUTTON
	RightButton  MouseButton = win.MK_RBUTTON
	MiddleButton MouseButton = win.MK_MBUTTON
)

type mouseEventHandlerInfo struct {
	handler MouseEventHandler
	once    bool
}

// MouseEventHandler is called for mouse events. x and y are measured in native pixels.
type MouseEventHandler func(x, y int, button MouseButton)

type MouseEvent struct {
	handlers []mouseEventHandlerInfo
}

func (e *MouseEvent) Attach(handler MouseEventHandler) int {
	handlerInfo := mouseEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *MouseEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *MouseEvent) Once(handler MouseEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type MouseEventPublisher struct {
	event MouseEvent
}

func (p *MouseEventPublisher) Event() *MouseEvent {
	return &p.event
}

// Publish publishes mouse event. x and y are measured in native pixels.
func (p *MouseEventPublisher) Publish(x, y int, button MouseButton) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(x, y, button)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}

func MouseWheelEventDelta(button MouseButton) int {
	return int(int32(button) >> 16)
}

func MouseWheelEventKeyState(button MouseButton) int {
	return int(int32(button) & 0xFFFF)
}
