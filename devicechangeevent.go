// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

const (
	DBT_DEVICEARRIVAL           = 0x8000
	DBT_DEVICEREMOVECOMPLETE    = 0x8004
)

type DeviceArrivalEventHandler func()
type DeviceRemoveEventHandler func()

type DeviceChangeEvent struct {
	hWnd     win.HWND
	arrivalhandlers []DeviceArrivalEventHandler
	removehandlers []DeviceRemoveEventHandler
}


func (e *DeviceChangeEvent) AttachArrival(handler DeviceArrivalEventHandler) int {
	for i, h := range e.arrivalhandlers {
		if h == nil {
			e.arrivalhandlers[i] = handler
			return i
		}
	}

	e.arrivalhandlers = append(e.arrivalhandlers, handler)
	return len(e.arrivalhandlers) - 1
}

func (e *DeviceChangeEvent) AttachRemove(handler DeviceRemoveEventHandler) int {
	for i, h := range e.removehandlers {
		if h == nil {
			e.removehandlers[i] = handler
			return i
		}
	}

	e.removehandlers = append(e.removehandlers, handler)
	return len(e.removehandlers) - 1
}

func (e *DeviceChangeEvent) DetachArrival(handle int) {
	e.arrivalhandlers[handle] = nil
	for _, h := range e.arrivalhandlers {
		if h != nil {
			return
		}
	}
}

func (e *DeviceChangeEvent) DetachRemove(handle int) {
	e.removehandlers[handle] = nil
	for _, h := range e.removehandlers {
		if h != nil {
			return
		}
	}
}

type DeviceChangeEventPublisher struct {
	event DeviceChangeEvent
}

func (p *DeviceChangeEventPublisher) Event(hWnd win.HWND) *DeviceChangeEvent {
	p.event.hWnd = hWnd
	return &p.event
}

func (p *DeviceChangeEventPublisher) Publish(wParam int) {
	if wParam == DBT_DEVICEARRIVAL {
		for _, handler := range p.event.arrivalhandlers {
			if handler != nil {
				handler()
			}
		}
	} else if wParam == DBT_DEVICEREMOVECOMPLETE {
		for _, handler := range p.event.removehandlers {
			if handler != nil {
				handler()
			}
		}
	}
}
