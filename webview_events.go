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

func (arg *WebViewNavigatingArg) Cancel() bool {
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

func (arg *WebViewNavigatingArg) SetCancel(value bool) {
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

type WebViewNavigatedEventArg struct {
	pDisp *win.IDispatch
	url   *win.VARIANT
}

func (arg *WebViewNavigatedEventArg) Url() string {
	url := arg.url
	if url != nil && url.MustBSTR() != nil {
		return win.BSTRToString(url.MustBSTR())
	}
	return ""
}

type WebViewNavigatedEventHandler func(arg *WebViewNavigatedEventArg)

type WebViewNavigatedEvent struct {
	handlers []WebViewNavigatedEventHandler
}

func (e *WebViewNavigatedEvent) Attach(handler WebViewNavigatedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewNavigatedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewNavigatedEventPublisher struct {
	event WebViewNavigatedEvent
}

func (p *WebViewNavigatedEventPublisher) Event() *WebViewNavigatedEvent {
	return &p.event
}

func (p *WebViewNavigatedEventPublisher) Publish(arg *WebViewNavigatedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewDownloadingEvent struct {
	handlers []EventHandler
}

func (e *WebViewDownloadingEvent) Attach(handler EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewDownloadingEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewDownloadingEventPublisher struct {
	event WebViewDownloadingEvent
}

func (p *WebViewDownloadingEventPublisher) Event() *WebViewDownloadingEvent {
	return &p.event
}

func (p *WebViewDownloadingEventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}

type WebViewDownloadedEvent struct {
	handlers []EventHandler
}

func (e *WebViewDownloadedEvent) Attach(handler EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewDownloadedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewDownloadedEventPublisher struct {
	event WebViewDownloadedEvent
}

func (p *WebViewDownloadedEventPublisher) Event() *WebViewDownloadedEvent {
	return &p.event
}

func (p *WebViewDownloadedEventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
		}
	}
}

type WebViewDocumentCompletedEventArg struct {
	pDisp *win.IDispatch
	url   *win.VARIANT
}

func (arg *WebViewDocumentCompletedEventArg) Url() string {
	url := arg.url
	if url != nil && url.MustBSTR() != nil {
		return win.BSTRToString(url.MustBSTR())
	}
	return ""
}

type WebViewDocumentCompletedEventHandler func(arg *WebViewDocumentCompletedEventArg)

type WebViewDocumentCompletedEvent struct {
	handlers []WebViewDocumentCompletedEventHandler
}

func (e *WebViewDocumentCompletedEvent) Attach(handler WebViewDocumentCompletedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewDocumentCompletedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewDocumentCompletedEventPublisher struct {
	event WebViewDocumentCompletedEvent
}

func (p WebViewDocumentCompletedEventPublisher) Event() *WebViewDocumentCompletedEvent {
	return &p.event
}

func (p *WebViewDocumentCompletedEventPublisher) Publish(arg *WebViewDocumentCompletedEventArg) {
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

func (arg *WebViewNavigatedErrorEventArg) Cancel() bool {
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

func (arg *WebViewNavigatedErrorEventArg) SetCancel(value bool) {
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

func (arg *WebViewNewWindowEventArg) Cancel() bool {
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

func (arg *WebViewNewWindowEventArg) SetCancel(value bool) {
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

type WebViewQuittingEvent struct {
	handlers []EventHandler
}

func (e *WebViewQuittingEvent) Attach(handler EventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewQuittingEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewQuittingEventPublisher struct {
	event WebViewQuittingEvent
}

func (p *WebViewQuittingEventPublisher) Event() *WebViewQuittingEvent {
	return &p.event
}

func (p *WebViewQuittingEventPublisher) Publish() {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler()
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

func (arg *WebViewWindowClosingEventArg) Cancel() bool {
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

func (arg *WebViewWindowClosingEventArg) SetCancel(value bool) {
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

type WebViewStatusBarVisibleChangedEventArg struct {
	statusBar win.VARIANT_BOOL
}

func (arg *WebViewStatusBarVisibleChangedEventArg) Visible() bool {
	statusBar := arg.statusBar
	if statusBar != win.VARIANT_FALSE {
		return true
	} else {
		return false
	}
	return false
}

type WebViewStatusBarVisibleChangedEventHandler func(arg *WebViewStatusBarVisibleChangedEventArg)

type WebViewStatusBarVisibleChangedEvent struct {
	handlers []WebViewStatusBarVisibleChangedEventHandler
}

func (e *WebViewStatusBarVisibleChangedEvent) Attach(handler WebViewStatusBarVisibleChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewStatusBarVisibleChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewStatusBarVisibleChangedEventPublisher struct {
	event WebViewStatusBarVisibleChangedEvent
}

func (p *WebViewStatusBarVisibleChangedEventPublisher) Event() *WebViewStatusBarVisibleChangedEvent {
	return &p.event
}

func (p *WebViewStatusBarVisibleChangedEventPublisher) Publish(arg *WebViewStatusBarVisibleChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewTheaterModeChangedEventArg struct {
	theaterMode win.VARIANT_BOOL
}

func (arg *WebViewTheaterModeChangedEventArg) IsTheaterMode() bool {
	theaterMode := arg.theaterMode
	if theaterMode != win.VARIANT_FALSE {
		return true
	} else {
		return false
	}
	return false
}

type WebViewTheaterModeChangedEventHandler func(arg *WebViewTheaterModeChangedEventArg)

type WebViewTheaterModeChangedEvent struct {
	handlers []WebViewTheaterModeChangedEventHandler
}

func (e *WebViewTheaterModeChangedEvent) Attach(handler WebViewTheaterModeChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewTheaterModeChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewTheaterModeChangedEventPublisher struct {
	event WebViewTheaterModeChangedEvent
}

func (p *WebViewTheaterModeChangedEventPublisher) Event() *WebViewTheaterModeChangedEvent {
	return &p.event
}

func (p *WebViewTheaterModeChangedEventPublisher) Publish(arg *WebViewTheaterModeChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewToolBarVisibleChangedEventArg struct {
	toolBar win.VARIANT_BOOL
}

func (arg *WebViewToolBarVisibleChangedEventArg) Visible() bool {
	toolBar := arg.toolBar
	if toolBar != win.VARIANT_FALSE {
		return true
	} else {
		return false
	}
	return false
}

type WebViewToolBarVisibleChangedEventHandler func(arg *WebViewToolBarVisibleChangedEventArg)

type WebViewToolBarVisibleChangedEvent struct {
	handlers []WebViewToolBarVisibleChangedEventHandler
}

func (e *WebViewToolBarVisibleChangedEvent) Attach(handler WebViewToolBarVisibleChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewToolBarVisibleChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewToolBarVisibleChangedEventPublisher struct {
	event WebViewToolBarVisibleChangedEvent
}

func (p *WebViewToolBarVisibleChangedEventPublisher) Event() *WebViewToolBarVisibleChangedEvent {
	return &p.event
}

func (p *WebViewToolBarVisibleChangedEventPublisher) Publish(arg *WebViewToolBarVisibleChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewBrowserVisibleChangedEventArg struct {
	vVisible win.VARIANT_BOOL
}

func (arg *WebViewBrowserVisibleChangedEventArg) Visible() bool {
	vVisible := arg.vVisible
	if vVisible != win.VARIANT_FALSE {
		return true
	} else {
		return false
	}
	return false
}

type WebViewBrowserVisibleChangedEventHandler func(arg *WebViewBrowserVisibleChangedEventArg)

type WebViewBrowserVisibleChangedEvent struct {
	handlers []WebViewBrowserVisibleChangedEventHandler
}

func (e *WebViewBrowserVisibleChangedEvent) Attach(handler WebViewBrowserVisibleChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewBrowserVisibleChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewBrowserVisibleChangedEventPublisher struct {
	event WebViewBrowserVisibleChangedEvent
}

func (p *WebViewBrowserVisibleChangedEventPublisher) Event() *WebViewBrowserVisibleChangedEvent {
	return &p.event
}

func (p *WebViewBrowserVisibleChangedEventPublisher) Publish(arg *WebViewBrowserVisibleChangedEventArg) {
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

type WebViewProgressChangedEventArg struct {
	nProgress    int32
	nProgressMax int32
}

func (arg *WebViewProgressChangedEventArg) Progress() int32 {
	return arg.nProgress
}

func (arg *WebViewProgressChangedEventArg) ProgressMax() int32 {
	return arg.nProgressMax
}

type WebViewProgressChangedEventHandler func(arg *WebViewProgressChangedEventArg)

type WebViewProgressChangedEvent struct {
	handlers []WebViewProgressChangedEventHandler
}

func (e *WebViewProgressChangedEvent) Attach(handler WebViewProgressChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewProgressChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewProgressChangedEventPublisher struct {
	event WebViewProgressChangedEvent
}

func (p *WebViewProgressChangedEventPublisher) Event() *WebViewProgressChangedEvent {
	return &p.event
}

func (p *WebViewProgressChangedEventPublisher) Publish(arg *WebViewProgressChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewStatusTextChangedEventArg struct {
	sText *uint16
}

func (arg *WebViewStatusTextChangedEventArg) StatusText() string {
	sText := arg.sText
	if sText != nil {
		return win.BSTRToString(sText)
	}
	return ""
}

type WebViewStatusTextChangedEventHandler func(arg *WebViewStatusTextChangedEventArg)

type WebViewStatusTextChangedEvent struct {
	handlers []WebViewStatusTextChangedEventHandler
}

func (e *WebViewStatusTextChangedEvent) Attach(handler WebViewStatusTextChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewStatusTextChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewStatusTextChangedEventPublisher struct {
	event WebViewStatusTextChangedEvent
}

func (p *WebViewStatusTextChangedEventPublisher) Event() *WebViewStatusTextChangedEvent {
	return &p.event
}

func (p *WebViewStatusTextChangedEventPublisher) Publish(arg *WebViewStatusTextChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}

type WebViewTitleChangedEventArg struct {
	sText *uint16
}

func (arg *WebViewTitleChangedEventArg) Title() string {
	sText := arg.sText
	if sText != nil {
		return win.BSTRToString(sText)
	}
	return ""
}

type WebViewTitleChangedEventHandler func(arg *WebViewTitleChangedEventArg)

type WebViewTitleChangedEvent struct {
	handlers []WebViewTitleChangedEventHandler
}

func (e *WebViewTitleChangedEvent) Attach(handler WebViewTitleChangedEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *WebViewTitleChangedEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type WebViewTitleChangedEventPublisher struct {
	event WebViewTitleChangedEvent
}

func (p *WebViewTitleChangedEventPublisher) Event() *WebViewTitleChangedEvent {
	return &p.event
}

func (p *WebViewTitleChangedEventPublisher) Publish(arg *WebViewTitleChangedEventArg) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(arg)
		}
	}
}
