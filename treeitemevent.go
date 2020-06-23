// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

type treeItemEventHandlerInfo struct {
	handler TreeItemEventHandler
	once    bool
}

type TreeItemEventHandler func(item TreeItem)

type TreeItemEvent struct {
	handlers []treeItemEventHandlerInfo
}

func (e *TreeItemEvent) Attach(handler TreeItemEventHandler) int {
	handlerInfo := treeItemEventHandlerInfo{handler, false}

	for i, h := range e.handlers {
		if h.handler == nil {
			e.handlers[i] = handlerInfo
			return i
		}
	}

	e.handlers = append(e.handlers, handlerInfo)

	return len(e.handlers) - 1
}

func (e *TreeItemEvent) Detach(handle int) {
	e.handlers[handle].handler = nil
}

func (e *TreeItemEvent) Once(handler TreeItemEventHandler) {
	i := e.Attach(handler)
	e.handlers[i].once = true
}

type TreeItemEventPublisher struct {
	event TreeItemEvent
}

func (p *TreeItemEventPublisher) Event() *TreeItemEvent {
	return &p.event
}

func (p *TreeItemEventPublisher) Publish(item TreeItem) {
	for i, h := range p.event.handlers {
		if h.handler != nil {
			h.handler(item)

			if h.once {
				p.event.Detach(i)
			}
		}
	}
}
