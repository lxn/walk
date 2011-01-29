// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	"walk/drawing"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

const mainWindowWindowClass = `\o/ Walk_MainWindow_Class \o/`

var mainWindowWndProcCallback *syscall.Callback

func mainWindowWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	mw, ok := widgetsByHWnd[msg.HWnd].(*MainWindow)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return mw.wndProc(msg, 0)
}

type MainWindow struct {
	TopLevelWindow
	menu    *Menu
	toolBar *ToolBar
}

func NewMainWindow() (mw *MainWindow, err os.Error) {
	ensureRegisteredWindowClass(mainWindowWindowClass, mainWindowWndProc, &mainWindowWndProcCallback)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(mainWindowWindowClass), nil,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 300, 0, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	wnd := &MainWindow{TopLevelWindow: TopLevelWindow{Container: Container{Widget: Widget{hWnd: hWnd}}}}

	defer func() {
		if x := recover(); x != nil {
			wnd.Dispose()

			err = toError(x)
		}
	}()

	wnd.children = newObservedWidgetList(wnd)

	widgetsByHWnd[hWnd] = wnd

	wnd.SetLayout(NewVBoxLayout())

	wnd.menu, err = newMenuBar()
	if err != nil {
		panic(err)
	}
	SetMenu(wnd.hWnd, wnd.menu.hMenu)

	wnd.toolBar, err = NewToolBar(wnd)
	if err != nil {
		panic(err)
	}

	wnd.clientArea, err = NewComposite(wnd)
	if err != nil {
		panic(err)
	}

	// This forces display of focus rectangles, as soon as the user starts to type.
	SendMessage(hWnd, WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	mw = wnd

	return
}

func (mw *MainWindow) Menu() *Menu {
	return mw.menu
}

func (mw *MainWindow) ToolBar() *ToolBar {
	return mw.toolBar
}

func (mw *MainWindow) ClientBounds() (bounds drawing.Rectangle, err os.Error) {
	bounds, err = mw.Widget.ClientBounds()
	if err != nil {
		return
	}

	if mw.toolBar.Actions().Len() > 0 {
		tlbBounds, err := mw.toolBar.Bounds()
		if err != nil {
			return
		}

		bounds.Y += tlbBounds.Height
		bounds.Height -= tlbBounds.Height
	}

	return
}

func (mw *MainWindow) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_SIZE, WM_SIZING:
		SendMessage(mw.toolBar.hWnd, TB_AUTOSIZE, 0, 0)
	}

	return mw.TopLevelWindow.wndProc(msg, origWndProcPtr)
}
