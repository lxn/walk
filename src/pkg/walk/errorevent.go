// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

type ErrorEventHandler func(err os.Error)

type ErrorEvent struct {
	handlers []ErrorEventHandler
}

func (e *ErrorEvent) Attach(handler ErrorEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *ErrorEvent) Detach(handler ErrorEventHandler) {
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

type ErrorEventPublisher struct {
	event ErrorEvent
}

func (p *ErrorEventPublisher) Event() *ErrorEvent {
	return &p.event
}

func (p *ErrorEventPublisher) Publish(err os.Error) {
	for _, handler := range p.event.handlers {
		handler(err)
	}
}
