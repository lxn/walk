// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

import (
	"github.com/lxn/go-winapi"
	"github.com/lxn/walk"
)

const myWidgetWindowClass = "MyWidget Class"

func init() {
	walk.MustRegisterWindowClass(myWidgetWindowClass)
}

type MyWidget struct {
	walk.WidgetBase
}

func NewMyWidget(parent walk.Container) (*MyWidget, error) {
	w := new(MyWidget)

	if err := walk.InitChildWidget(
		w,
		parent,
		myWidgetWindowClass,
		winapi.WS_VISIBLE,
		0); err != nil {

		return nil, err
	}

	bg, err := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
	if err != nil {
		return nil, err
	}
	w.SetBackground(bg)

	return w, nil
}

func (*MyWidget) MinSizeHint() walk.Size {
	return walk.Size{50, 50}
}

func (w *MyWidget) WndProc(hwnd winapi.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case winapi.WM_LBUTTONDOWN:
		log.Printf("%s: WM_LBUTTONDOWN", w.Name())
	}

	return w.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

type MyPushButton struct {
	*walk.PushButton
}

func NewMyPushButton(parent walk.Container) (*MyPushButton, error) {
	pb, err := walk.NewPushButton(parent)
	if err != nil {
		return nil, err
	}

	mpb := &MyPushButton{pb}

	if err := walk.InitWrapperWidget(mpb); err != nil {
		return nil, err
	}

	return mpb, nil
}

func (mpb *MyPushButton) WndProc(hwnd winapi.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case winapi.WM_LBUTTONDOWN:
		log.Printf("%s: WM_LBUTTONDOWN", mpb.Text())
	}

	return mpb.PushButton.WndProc(hwnd, msg, wParam, lParam)
}

func main() {
	walk.SetPanicOnError(true)

	mw, _ := walk.NewMainWindow()

	mw.SetTitle("Walk External Widgets Example")
	mw.SetLayout(walk.NewHBoxLayout())

	a, _ := NewMyWidget(mw)
	a.SetName("a")

	b, _ := NewMyWidget(mw)
	b.SetName("b")

	c, _ := NewMyWidget(mw)
	c.SetName("c")

	mpb, _ := NewMyPushButton(mw)
	mpb.SetText("MyPushButton")

	mw.SetSize(walk.Size{400, 300})
	mw.Show()

	mw.Run()
}
