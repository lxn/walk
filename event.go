// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type EventHandler func()

type Event struct {
	handlers []EventHandler
}

func (e *Event) Attach(handler EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *Event) Detach(handle int) {
	e.handlers[handle] = nil
}

type EventPublisher struct {
	event Event
}

func (p *EventPublisher) Event() *Event {
	return &p.event
}

func (p *EventPublisher) Publish() {
	events := inProgressEventsByForm[appSingleton.activeForm]
	events = append(events, &p.event)
	inProgressEventsByForm[appSingleton.activeForm] = events

	defer func() {
		events = events[:len(events)-1]
		if len(events) == 0 {
			delete(inProgressEventsByForm, appSingleton.activeForm)
		} else {
			inProgressEventsByForm[appSingleton.activeForm] = events
			return
		}

		layouts := scheduledLayoutsByForm[appSingleton.activeForm]
		delete(scheduledLayoutsByForm, appSingleton.activeForm)
		if len(layouts) == 0 {
			return
		}

		old := performingScheduledLayouts
		performingScheduledLayouts = true
		defer func() {
			performingScheduledLayouts = old
		}()

		if formResizeScheduled {
			formResizeScheduled = false

			bounds := appSingleton.activeForm.Bounds()

			if appSingleton.activeForm.AsFormBase().fixedSize() {
				bounds.Width, bounds.Height = 0, 0
			}

			appSingleton.activeForm.SetBounds(bounds)
		}

		for _, layout := range layouts {
			if widget, ok := layout.Container().(Widget); ok && widget.Form() != appSingleton.activeForm {
				continue
			}

			layout.Update(false)
		}
	}()

	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}
