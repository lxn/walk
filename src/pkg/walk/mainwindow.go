// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import . "walk/winapi"

const mainWindowWindowClass = `\o/ Walk_MainWindow_Class \o/`

var mainWindowWindowClassRegistered bool

type MainWindow struct {
	TopLevelWindow
	menu            *Menu
	toolBar         *ToolBar
	clientComposite *Composite
}

func NewMainWindow() (*MainWindow, os.Error) {
	ensureRegisteredWindowClass(mainWindowWindowClass, &mainWindowWindowClassRegistered)

	mw := &MainWindow{}

	if err := initWidget(
		mw,
		nil,
		mainWindowWindowClass,
		WS_OVERLAPPEDWINDOW,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			mw.Dispose()
		}
	}()

	mw.SetPersistent(true)

	var err os.Error

	if mw.menu, err = newMenuBar(); err != nil {
		return nil, err
	}
	SetMenu(mw.hWnd, mw.menu.hMenu)

	if mw.toolBar, err = NewToolBar(mw); err != nil {
		return nil, err
	}

	if mw.clientComposite, err = NewComposite(mw); err != nil {
		return nil, err
	}
	mw.clientComposite.SetName("clientComposite")

	// This forces display of focus rectangles, as soon as the user starts to type.
	SendMessage(mw.hWnd, WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	succeeded = true

	return mw, nil
}

func (mw *MainWindow) Children() *WidgetList {
	if mw.clientComposite == nil {
		return nil
	}

	return mw.clientComposite.Children()
}

func (mw *MainWindow) Layout() Layout {
	if mw.clientComposite == nil {
		return nil
	}

	return mw.clientComposite.Layout()
}

func (mw *MainWindow) SetLayout(value Layout) os.Error {
	if mw.clientComposite == nil {
		return newError("clientComposite not initialized")
	}

	return mw.clientComposite.SetLayout(value)
}

func (mw *MainWindow) Menu() *Menu {
	return mw.menu
}

func (mw *MainWindow) ToolBar() *ToolBar {
	return mw.toolBar
}

func (mw *MainWindow) ClientBounds() Rectangle {
	bounds := mw.WidgetBase.ClientBounds()

	if mw.toolBar.Actions().Len() > 0 {
		tlbBounds := mw.toolBar.Bounds()

		bounds.Y += tlbBounds.Height
		bounds.Height -= tlbBounds.Height
	}

	return bounds
}

func (mw *MainWindow) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_SIZE, WM_SIZING:
		SendMessage(mw.toolBar.hWnd, TB_AUTOSIZE, 0, 0)

		mw.clientComposite.SetBounds(mw.ClientBounds())
	}

	return mw.TopLevelWindow.wndProc(hwnd, msg, wParam, lParam)
}
