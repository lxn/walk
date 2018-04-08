// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type WvBeforeNavigate2EventHandler func(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	flags *win.VARIANT,
	targetFrameName *win.VARIANT,
	postData *win.VARIANT,
	headers *win.VARIANT,
	cancel *win.VARIANT_BOOL)

type WvBeforeNavigate2Event struct {
	handlers []WvBeforeNavigate2EventHandler
}

func (e *WvBeforeNavigate2Event) Attach(handler WvBeforeNavigate2EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvBeforeNavigate2Event) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvBeforeNavigate2EventPublisher struct {
	event WvBeforeNavigate2Event
}

func (p *WvBeforeNavigate2EventPublisher) Event() *WvBeforeNavigate2Event {
	return &p.event
}

func (p *WvBeforeNavigate2EventPublisher) Publish(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	flags *win.VARIANT,
	targetFrameName *win.VARIANT,
	postData *win.VARIANT,
	headers *win.VARIANT,
	cancel *win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(
				pDisp,
				url,
				flags,
				targetFrameName,
				postData,
				headers,
				cancel)
		}
	}
}

type WvNavigateComplete2EventHandler func(pDisp *win.IDispatch, url *win.VARIANT)

type WvNavigateComplete2Event struct {
	handlers []WvNavigateComplete2EventHandler
}

func (e *WvNavigateComplete2Event) Attach(handler WvNavigateComplete2EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvNavigateComplete2Event) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvNavigateComplete2EventPublisher struct {
	event WvNavigateComplete2Event
}

func (p *WvNavigateComplete2EventPublisher) Event() *WvNavigateComplete2Event {
	return &p.event
}

func (p *WvNavigateComplete2EventPublisher) Publish(pDisp *win.IDispatch, url *win.VARIANT) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(pDisp, url)
		}
	}
}

type WvDownloadBeginEventHandler func()

type WvDownloadBeginEvent struct {
	handlers []WvDownloadBeginEventHandler
}

func (e *WvDownloadBeginEvent) Attach(handler WvDownloadBeginEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvDownloadBeginEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvDownloadBeginEventPublisher struct {
	event WvDownloadBeginEvent
}

func (p *WvDownloadBeginEventPublisher) Event() *WvDownloadBeginEvent {
	return &p.event
}

func (p *WvDownloadBeginEventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}

type WvDownloadCompleteEventHandler func()

type WvDownloadCompleteEvent struct {
	handlers []WvDownloadCompleteEventHandler
}

func (e *WvDownloadCompleteEvent) Attach(handler WvDownloadCompleteEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvDownloadCompleteEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvDownloadCompleteEventPublisher struct {
	event WvDownloadCompleteEvent
}

func (p *WvDownloadCompleteEventPublisher) Event() *WvDownloadCompleteEvent {
	return &p.event
}

func (p *WvDownloadCompleteEventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}

type WvDocumentCompleteEventHandler func(pDisp *win.IDispatch, url *win.VARIANT)

type WvDocumentCompleteEvent struct {
	handlers []WvDocumentCompleteEventHandler
}

func (e *WvDocumentCompleteEvent) Attach(handler WvDocumentCompleteEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvDocumentCompleteEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvDocumentCompleteEventPublisher struct {
	event WvDocumentCompleteEvent
}

func (p *WvDocumentCompleteEventPublisher) Event() *WvDocumentCompleteEvent {
	return &p.event
}

func (p *WvDocumentCompleteEventPublisher) Publish(pDisp *win.IDispatch, url *win.VARIANT) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(pDisp, url)
		}
	}
}

type WvNavigateErrorEventHandler func(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	targetFrameName *win.VARIANT,
	statusCode *win.VARIANT,
	cancel *win.VARIANT_BOOL)

type WvNavigateErrorEvent struct {
	handlers []WvNavigateErrorEventHandler
}

func (e *WvNavigateErrorEvent) Attach(handler WvNavigateErrorEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvNavigateErrorEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvNavigateErrorEventPublisher struct {
	event WvNavigateErrorEvent
}

func (p *WvNavigateErrorEventPublisher) Event() *WvNavigateErrorEvent {
	return &p.event
}

func (p *WvNavigateErrorEventPublisher) Publish(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	targetFrameName *win.VARIANT,
	statusCode *win.VARIANT,
	cancel *win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(
				pDisp,
				url,
				targetFrameName,
				statusCode,
				cancel)
		}
	}
}

type WvNewWindow3EventHandler func(
	ppDisp **win.IDispatch,
	cancel *win.VARIANT_BOOL,
	dwFlags uint32,
	bstrUrlContext *uint16,
	bstrUrl *uint16)

type WvNewWindow3Event struct {
	handlers []WvNewWindow3EventHandler
}

func (e *WvNewWindow3Event) Attach(handler WvNewWindow3EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvNewWindow3Event) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvNewWindow3EventPublisher struct {
	event WvNewWindow3Event
}

func (p *WvNewWindow3EventPublisher) Event() *WvNewWindow3Event {
	return &p.event
}

func (p *WvNewWindow3EventPublisher) Publish(
	ppDisp **win.IDispatch,
	cancel *win.VARIANT_BOOL,
	dwFlags uint32,
	bstrUrlContext *uint16,
	bstrUrl *uint16) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(
				ppDisp,
				cancel,
				dwFlags,
				bstrUrlContext,
				bstrUrl)
		}
	}
}

type WvOnQuitEventHandler func()

type WvOnQuitEvent struct {
	handlers []WvOnQuitEventHandler
}

func (e *WvOnQuitEvent) Attach(handler WvOnQuitEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvOnQuitEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvOnQuitEventPublisher struct {
	event WvOnQuitEvent
}

func (p *WvOnQuitEventPublisher) Event() *WvOnQuitEvent {
	return &p.event
}

func (p *WvOnQuitEventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}

type WvWindowClosingEventHandler func(bIsChildWindow win.VARIANT_BOOL, cancel *win.VARIANT_BOOL)

type WvWindowClosingEvent struct {
	handlers []WvWindowClosingEventHandler
}

func (e *WvWindowClosingEvent) Attach(handler WvWindowClosingEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvWindowClosingEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvWindowClosingEventPublisher struct {
	event WvWindowClosingEvent
}

func (p *WvWindowClosingEventPublisher) Event() *WvWindowClosingEvent {
	return &p.event
}

func (p *WvWindowClosingEventPublisher) Publish(bIsChildWindow win.VARIANT_BOOL, cancel *win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(bIsChildWindow, cancel)
		}
	}
}

type WvOnStatusBarEventHandler func(statusBar win.VARIANT_BOOL)

type WvOnStatusBarEvent struct {
	handlers []WvOnStatusBarEventHandler
}

func (e *WvOnStatusBarEvent) Attach(handler WvOnStatusBarEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvOnStatusBarEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvOnStatusBarEventPublisher struct {
	event WvOnStatusBarEvent
}

func (p *WvOnStatusBarEventPublisher) Event() *WvOnStatusBarEvent {
	return &p.event
}

func (p *WvOnStatusBarEventPublisher) Publish(statusBar win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(statusBar)
		}
	}
}

type WvOnTheaterModeEventHandler func(theaterMode win.VARIANT_BOOL)

type WvOnTheaterModeEvent struct {
	handlers []WvOnTheaterModeEventHandler
}

func (e *WvOnTheaterModeEvent) Attach(handler WvOnTheaterModeEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvOnTheaterModeEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvOnTheaterModeEventPublisher struct {
	event WvOnTheaterModeEvent
}

func (p *WvOnTheaterModeEventPublisher) Event() *WvOnTheaterModeEvent {
	return &p.event
}

func (p *WvOnTheaterModeEventPublisher) Publish(theaterMode win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(theaterMode)
		}
	}
}

type WvOnToolBarEventHandler func(toolBar win.VARIANT_BOOL)

type WvOnToolBarEvent struct {
	handlers []WvOnToolBarEventHandler
}

func (e *WvOnToolBarEvent) Attach(handler WvOnToolBarEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvOnToolBarEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvOnToolBarEventPublisher struct {
	event WvOnToolBarEvent
}

func (p *WvOnToolBarEventPublisher) Event() *WvOnToolBarEvent {
	return &p.event
}

func (p *WvOnToolBarEventPublisher) Publish(toolBar win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(toolBar)
		}
	}
}

type WvOnVisibleEventHandler func(vVisible win.VARIANT_BOOL)

type WvOnVisibleEvent struct {
	handlers []WvOnVisibleEventHandler
}

func (e *WvOnVisibleEvent) Attach(handler WvOnVisibleEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvOnVisibleEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvOnVisibleEventPublisher struct {
	event WvOnVisibleEvent
}

func (p *WvOnVisibleEventPublisher) Event() *WvOnVisibleEvent {
	return &p.event
}

func (p *WvOnVisibleEventPublisher) Publish(vVisible win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(vVisible)
		}
	}
}

type WvCommandStateChangeEventHandler func(command int32, enable win.VARIANT_BOOL)

type WvCommandStateChangeEvent struct {
	handlers []WvCommandStateChangeEventHandler
}

func (e *WvCommandStateChangeEvent) Attach(handler WvCommandStateChangeEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvCommandStateChangeEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvCommandStateChangeEventPublisher struct {
	event WvCommandStateChangeEvent
}

func (p *WvCommandStateChangeEventPublisher) Event() *WvCommandStateChangeEvent {
	return &p.event
}

func (p *WvCommandStateChangeEventPublisher) Publish(command int32, enable win.VARIANT_BOOL) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(command, enable)
		}
	}
}

type WvProgressChangeEventHandler func(nProgress int32, nProgressMax int32)

type WvProgressChangeEvent struct {
	handlers []WvProgressChangeEventHandler
}

func (e *WvProgressChangeEvent) Attach(handler WvProgressChangeEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvProgressChangeEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvProgressChangeEventPublisher struct {
	event WvProgressChangeEvent
}

func (p *WvProgressChangeEventPublisher) Event() *WvProgressChangeEvent {
	return &p.event
}

func (p *WvProgressChangeEventPublisher) Publish(nProgress int32, nProgressMax int32) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(nProgress, nProgressMax)
		}
	}
}

type WvStatusTextChangeEventHandler func(sText *uint16)

type WvStatusTextChangeEvent struct {
	handlers []WvStatusTextChangeEventHandler
}

func (e *WvStatusTextChangeEvent) Attach(handler WvStatusTextChangeEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvStatusTextChangeEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvStatusTextChangeEventPublisher struct {
	event WvStatusTextChangeEvent
}

func (p *WvStatusTextChangeEventPublisher) Event() *WvStatusTextChangeEvent {
	return &p.event
}

func (p *WvStatusTextChangeEventPublisher) Publish(sText *uint16) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(sText)
		}
	}
}

type WvTitleChangeEventHandler func(sText *uint16)

type WvTitleChangeEvent struct {
	handlers []WvTitleChangeEventHandler
}

func (e *WvTitleChangeEvent) Attach(handler WvTitleChangeEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WvTitleChangeEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WvTitleChangeEventPublisher struct {
	event WvTitleChangeEvent
}

func (p *WvTitleChangeEventPublisher) Event() *WvTitleChangeEvent {
	return &p.event
}

func (p *WvTitleChangeEventPublisher) Publish(sText *uint16) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(sText)
		}
	}
}
