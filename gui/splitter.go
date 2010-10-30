// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

type Orientation byte

const (
	Horizontal Orientation = iota
	Vertical
)

const splitterWindowClass = `\o/ Walk_Splitter_Class \o/`

var splitterWndProcCallback *syscall.Callback

func splitterWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	s, ok := widgetsByHWnd[msg.HWnd].(*Splitter)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return s.wndProc(msg, 0)
}

type Splitter struct {
	Container
	orientation Orientation
}

func NewSplitter(parent IContainer) (*Splitter, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(splitterWindowClass, splitterWndProc, &splitterWndProcCallback)

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr(splitterWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 200, 100, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	s := &Splitter{Container: Container{Widget: Widget{hWnd: hWnd, parent: parent}}}

	s.children = newObservedWidgetList(s)

	s.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = s

	parent.Children().Add(s)

	return s, nil
}

func (s *Splitter) onInsertingWidget(index int, widget IWidget) (err os.Error) {
	return nil
}

func (s *Splitter) onInsertedWidget(index int, widget IWidget) (err os.Error) {
	panic("not implemented")
}

func (s *Splitter) onRemovingWidget(index int, widget IWidget) (err os.Error) {
	return s.Container.onRemovingWidget(index, widget)
}

func (s *Splitter) onRemovedWidget(index int, widget IWidget) (err os.Error) {
	panic("not implemented")
}

func (s *Splitter) onClearingWidgets() (err os.Error) {
	panic("not implemented")
}

func (s *Splitter) onClearedWidgets() (err os.Error) {
	panic("not implemented")
}
