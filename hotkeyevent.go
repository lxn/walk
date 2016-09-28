// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type HotkeyEventHandler func(hkid int)

type HotkeyEvent struct {
	handlers []HotkeyEventHandler
}

func (e *HotkeyEvent) Attach(handler HotkeyEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *HotkeyEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type HotkeyEventPublisher struct {
	event HotkeyEvent
}

func (p *HotkeyEventPublisher) Event() *HotkeyEvent {
	return &p.event
}

func (p *HotkeyEventPublisher) Publish(hkid int) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(hkid)
		}
	}
}
