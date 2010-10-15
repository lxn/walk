// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type HorizontalAlignment int

const (
	LeftAlignment HorizontalAlignment = iota
	RightAlignment
	CenterAlignment
)


type EventArgs interface {
	Sender() interface{}
}

type eventArgs struct {
	sender interface{}
}

func (a *eventArgs) Sender() interface{} {
	return a.sender
}

type EventHandler func(args EventArgs)


type KeyEventArgs interface {
	EventArgs
	Key() int
}

type keyEventArgs struct {
	eventArgs
	key int
}

func (a *keyEventArgs) Key() int {
	return a.key
}

type KeyEventHandler func(args KeyEventArgs)


type MouseButton int

const (
	LeftButton MouseButton = iota
	RightButton
	MiddleButton
)

type MouseEventArgs interface {
	EventArgs
	X() int
	Y() int
	Button() MouseButton
}

type mouseEventArgs struct {
	eventArgs
	x, y   int
	button MouseButton
}

func (a *mouseEventArgs) X() int {
	return a.x
}

func (a *mouseEventArgs) Y() int {
	return a.y
}

func (a *mouseEventArgs) Button() MouseButton {
	return a.button
}

type MouseEventHandler func(args MouseEventArgs)


type CancelEventArgs interface {
	EventArgs
	Canceled() bool
	SetCanceled(value bool)
}

type cancelEventArgs struct {
	eventArgs
	canceled bool
}

func (a *cancelEventArgs) Canceled() bool {
	return a.canceled
}

func (a *cancelEventArgs) SetCanceled(value bool) {
	a.canceled = value
}

type CancelEventHandler func(args CancelEventArgs)
