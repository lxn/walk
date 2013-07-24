// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

const mainWindowWindowClass = `\o/ Walk_MainWindow_Class \o/`

func init() {
	MustRegisterWindowClass(mainWindowWindowClass)
}

type MainWindow struct {
	FormBase
	windowPlacement *WINDOWPLACEMENT
	menu            *Menu
	toolBar         *ToolBar
	statusBar       *StatusBar
}

func NewMainWindow() (*MainWindow, error) {
	mw := new(MainWindow)

	if err := InitWindow(
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

	var err error

	if mw.menu, err = newMenuBar(); err != nil {
		return nil, err
	}
	if !SetMenu(mw.hWnd, mw.menu.hMenu) {
		return nil, lastError("SetMenu")
	}

	if mw.toolBar, err = NewToolBar(mw); err != nil {
		return nil, err
	}

	if mw.statusBar, err = NewStatusBar(mw); err != nil {
		return nil, err
	}

	// This forces display of focus rectangles, as soon as the user starts to type.
	mw.SendMessage(WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	succeeded = true

	return mw, nil
}

func (mw *MainWindow) Menu() *Menu {
	return mw.menu
}

func (mw *MainWindow) ToolBar() *ToolBar {
	return mw.toolBar
}

func (mw *MainWindow) StatusBar() *StatusBar {
	return mw.statusBar
}

func (mw *MainWindow) ClientBounds() Rectangle {
	bounds := mw.FormBase.ClientBounds()

	if mw.toolBar.Actions().Len() > 0 {
		tlbBounds := mw.toolBar.Bounds()

		bounds.Y += tlbBounds.Height
		bounds.Height -= tlbBounds.Height
	}

	if mw.statusBar.Visible() {
		bounds.Height -= mw.statusBar.Height()
	}

	return bounds
}

func (mw *MainWindow) SetVisible(visible bool) {
	if visible {
		DrawMenuBar(mw.hWnd)

		if mw.clientComposite.layout != nil {
			mw.clientComposite.layout.Update(false)
		}
	}

	mw.FormBase.SetVisible(visible)
}

func (mw *MainWindow) Fullscreen() bool {
	return GetWindowLong(mw.hWnd, GWL_STYLE)&WS_OVERLAPPEDWINDOW == 0
}

func (mw *MainWindow) SetFullscreen(fullscreen bool) error {
	if fullscreen == mw.Fullscreen() {
		return nil
	}

	if fullscreen {
		var mi MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		if mw.windowPlacement == nil {
			mw.windowPlacement = new(WINDOWPLACEMENT)
		}

		if !GetWindowPlacement(mw.hWnd, mw.windowPlacement) {
			return lastError("GetWindowPlacement")
		}
		if !GetMonitorInfo(MonitorFromWindow(
			mw.hWnd, MONITOR_DEFAULTTOPRIMARY), &mi) {

			return newError("GetMonitorInfo")
		}

		if err := mw.ensureStyleBits(WS_OVERLAPPEDWINDOW, false); err != nil {
			return err
		}

		if r := mi.RcMonitor; !SetWindowPos(
			mw.hWnd, HWND_TOP,
			r.Left, r.Top, r.Right-r.Left, r.Bottom-r.Top,
			SWP_FRAMECHANGED|SWP_NOOWNERZORDER) {

			return lastError("SetWindowPos")
		}
	} else {
		if err := mw.ensureStyleBits(WS_OVERLAPPEDWINDOW, true); err != nil {
			return err
		}

		if !SetWindowPlacement(mw.hWnd, mw.windowPlacement) {
			return lastError("SetWindowPlacement")
		}

		if !SetWindowPos(mw.hWnd, 0, 0, 0, 0, 0, SWP_FRAMECHANGED|SWP_NOMOVE|
			SWP_NOOWNERZORDER|SWP_NOSIZE|SWP_NOZORDER) {

			return lastError("SetWindowPos")
		}
	}

	return nil
}
