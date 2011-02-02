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

func NewMainWindow() (*MainWindow, os.Error) {
	ensureRegisteredWindowClass(mainWindowWindowClass, mainWindowWndProc, &mainWindowWndProcCallback)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(mainWindowWindowClass), nil,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 300, 0, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	mw := &MainWindow{
		TopLevelWindow: TopLevelWindow{
			Container: Container{
				Widget: Widget{
					hWnd: hWnd,
				},
			},
		},
	}

	succeeded := false
	defer func() {
		if !succeeded {
			mw.Dispose()
		}
	}()

	mw.SetPersistent(true)

	mw.children = newObservedWidgetList(mw)

	err := mw.SetLayout(NewVBoxLayout())
	if err != nil {
		return nil, err
	}

	if mw.menu, err = newMenuBar(); err != nil {
		return nil, err
	}
	SetMenu(mw.hWnd, mw.menu.hMenu)

	if mw.toolBar, err = NewToolBar(mw); err != nil {
		return nil, err
	}

	if mw.clientArea, err = NewComposite(mw); err != nil {
		return nil, err
	}
	mw.clientArea.SetName("clientArea")

	widgetsByHWnd[hWnd] = mw

	// This forces display of focus rectangles, as soon as the user starts to type.
	SendMessage(hWnd, WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	succeeded = true
	return mw, nil
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
