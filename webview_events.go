// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type WebViewNavigatingArg struct {
	pDisp           *win.IDispatch
	url             *win.VARIANT
	flags           *win.VARIANT
	targetFrameName *win.VARIANT
	postData        *win.VARIANT
	headers         *win.VARIANT
	cancel          *win.VARIANT_BOOL
}

func (arg *WebViewNavigatingArg) Url() string {
	url := arg.url
	if url != nil && url.MustBSTR() != nil {
		return win.BSTRToString(url.MustBSTR())
	}
	return ""
}

func (arg *WebViewNavigatingArg) Flags() int32 {
	flags := arg.flags
	if flags != nil {
		return flags.MustLong()
	}
	return 0
}

func (arg *WebViewNavigatingArg) Headers() string {
	headers := arg.headers
	if headers != nil && headers.MustBSTR() != nil {
		return win.BSTRToString(headers.MustBSTR())
	}
	return ""
}

func (arg *WebViewNavigatingArg) TargetFrameName() string {
	targetFrameName := arg.targetFrameName
	if targetFrameName != nil && targetFrameName.MustBSTR() != nil {
		return win.BSTRToString(targetFrameName.MustBSTR())
	}
	return ""
}

func (arg *WebViewNavigatingArg) Canceled() bool {
	cancel := arg.cancel
	if cancel != nil {
		if *cancel != win.VARIANT_FALSE {
			return true
		} else {
			return false
		}
	}
	return false
}

func (arg *WebViewNavigatingArg) SetCanceled(value bool) {
	cancel := arg.cancel
	if cancel != nil {
		if value {
			*cancel = win.VARIANT_TRUE
		} else {
			*cancel = win.VARIANT_FALSE
		}
	}
}

type WebViewNavigatingEventHandler func(arg *WebViewNavigatingArg)

type WebViewNavigatingEvent struct {
	handlers []WebViewNavigatingEventHandler
}

func (e *WebViewNavigatingEvent) Attach(handler WebViewNavigatingEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewNavigatingEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewNavigatingEventPublisher struct {
	event WebViewNavigatingEvent
}

func (p *WebViewNavigatingEventPublisher) Event() *WebViewNavigatingEvent {
	return &p.event
}

func (p *WebViewNavigatingEventPublisher) Publish(arg *WebViewNavigatingArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewNavigatedErrorEventArg struct {
	pDisp           *win.IDispatch
	url             *win.VARIANT
	targetFrameName *win.VARIANT
	statusCode      *win.VARIANT
	cancel          *win.VARIANT_BOOL
}

func (arg *WebViewNavigatedErrorEventArg) Url() string {
	url := arg.url
	if url != nil && url.MustBSTR() != nil {
		return win.BSTRToString(url.MustBSTR())
	}
	return ""
}

func (arg *WebViewNavigatedErrorEventArg) TargetFrameName() string {
	targetFrameName := arg.targetFrameName
	if targetFrameName != nil && targetFrameName.MustBSTR() != nil {
		return win.BSTRToString(targetFrameName.MustBSTR())
	}
	return ""
}

func (arg *WebViewNavigatedErrorEventArg) StatusCode() int32 {
	statusCode := arg.statusCode
	if statusCode != nil {
		return statusCode.MustLong()
	}
	return 0
}

func (arg *WebViewNavigatedErrorEventArg) Canceled() bool {
	cancel := arg.cancel
	if cancel != nil {
		if *cancel != win.VARIANT_FALSE {
			return true
		} else {
			return false
		}
	}
	return false
}

func (arg *WebViewNavigatedErrorEventArg) SetCanceled(value bool) {
	cancel := arg.cancel
	if cancel != nil {
		if value {
			*cancel = win.VARIANT_TRUE
		} else {
			*cancel = win.VARIANT_FALSE
		}
	}
}

type WebViewNavigatedErrorEventHandler func(arg *WebViewNavigatedErrorEventArg)

type WebViewNavigatedErrorEvent struct {
	handlers []WebViewNavigatedErrorEventHandler
}

func (e *WebViewNavigatedErrorEvent) Attach(handler WebViewNavigatedErrorEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewNavigatedErrorEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewNavigatedErrorEventPublisher struct {
	event WebViewNavigatedErrorEvent
}

func (p *WebViewNavigatedErrorEventPublisher) Event() *WebViewNavigatedErrorEvent {
	return &p.event
}

func (p *WebViewNavigatedErrorEventPublisher) Publish(arg *WebViewNavigatedErrorEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewNewWindowEventArg struct {
	ppDisp         **win.IDispatch
	cancel         *win.VARIANT_BOOL
	dwFlags        uint32
	bstrUrlContext *uint16
	bstrUrl        *uint16
}

func (arg *WebViewNewWindowEventArg) Canceled() bool {
	cancel := arg.cancel
	if cancel != nil {
		if *cancel != win.VARIANT_FALSE {
			return true
		} else {
			return false
		}
	}
	return false
}

func (arg *WebViewNewWindowEventArg) SetCanceled(value bool) {
	cancel := arg.cancel
	if cancel != nil {
		if value {
			*cancel = win.VARIANT_TRUE
		} else {
			*cancel = win.VARIANT_FALSE
		}
	}
}

func (arg *WebViewNewWindowEventArg) Flags() uint32 {
	return arg.dwFlags
}

func (arg *WebViewNewWindowEventArg) UrlContext() string {
	bstrUrlContext := arg.bstrUrlContext
	if bstrUrlContext != nil {
		return win.BSTRToString(bstrUrlContext)
	}
	return ""
}

func (arg *WebViewNewWindowEventArg) Url() string {
	bstrUrl := arg.bstrUrl
	if bstrUrl != nil {
		return win.BSTRToString(bstrUrl)
	}
	return ""
}

type WebViewNewWindowEventHandler func(arg *WebViewNewWindowEventArg)

type WebViewNewWindowEvent struct {
	handlers []WebViewNewWindowEventHandler
}

func (e *WebViewNewWindowEvent) Attach(handler WebViewNewWindowEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewNewWindowEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewNewWindowEventPublisher struct {
	event WebViewNewWindowEvent
}

func (p *WebViewNewWindowEventPublisher) Event() *WebViewNewWindowEvent {
	return &p.event
}

func (p *WebViewNewWindowEventPublisher) Publish(arg *WebViewNewWindowEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewWindowClosingEventArg struct {
	bIsChildWindow win.VARIANT_BOOL
	cancel         *win.VARIANT_BOOL
}

func (arg *WebViewWindowClosingEventArg) IsChildWindow() bool {
	bIsChildWindow := arg.bIsChildWindow
	if bIsChildWindow != win.VARIANT_FALSE {
		return true
	} else {
		return false
	}
	return false
}

func (arg *WebViewWindowClosingEventArg) Canceled() bool {
	cancel := arg.cancel
	if cancel != nil {
		if *cancel != win.VARIANT_FALSE {
			return true
		} else {
			return false
		}
	}
	return false
}

func (arg *WebViewWindowClosingEventArg) SetCanceled(value bool) {
	cancel := arg.cancel
	if cancel != nil {
		if value {
			*cancel = win.VARIANT_TRUE
		} else {
			*cancel = win.VARIANT_FALSE
		}
	}
}

type WebViewWindowClosingEventHandler func(arg *WebViewWindowClosingEventArg)

type WebViewWindowClosingEvent struct {
	handlers []WebViewWindowClosingEventHandler
}

func (e *WebViewWindowClosingEvent) Attach(handler WebViewWindowClosingEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewWindowClosingEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewWindowClosingEventPublisher struct {
	event WebViewWindowClosingEvent
}

func (p *WebViewWindowClosingEventPublisher) Event() *WebViewWindowClosingEvent {
	return &p.event
}

func (p *WebViewWindowClosingEventPublisher) Publish(arg *WebViewWindowClosingEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewCommandStateChangedEventArg struct {
	command int32
	enable  win.VARIANT_BOOL
}

func (arg *WebViewCommandStateChangedEventArg) Command() int32 {
	return arg.command
}

func (arg *WebViewCommandStateChangedEventArg) Enabled() bool {
	enable := arg.enable
	if enable != win.VARIANT_FALSE {
		return true
	} else {
		return false
	}
	return false
}

type WebViewCommandStateChangedEventHandler func(arg *WebViewCommandStateChangedEventArg)

type WebViewCommandStateChangedEvent struct {
	handlers []WebViewCommandStateChangedEventHandler
}

func (e *WebViewCommandStateChangedEvent) Attach(handler WebViewCommandStateChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewCommandStateChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewCommandStateChangedEventPublisher struct {
	event WebViewCommandStateChangedEvent
}

func (p *WebViewCommandStateChangedEventPublisher) Event() *WebViewCommandStateChangedEvent {
	return &p.event
}

func (p *WebViewCommandStateChangedEventPublisher) Publish(arg *WebViewCommandStateChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}
