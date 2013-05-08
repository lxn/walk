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
	TopLevelWindow
	windowPlacement *WINDOWPLACEMENT
	menu            *Menu
	toolBar         *ToolBar
	clientComposite *Composite
}

func NewMainWindow() (*MainWindow, error) {
	mw := &MainWindow{}

	if err := InitWidget(
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
	SetMenu(mw.hWnd, mw.menu.hMenu)

	if mw.toolBar, err = NewToolBar(mw); err != nil {
		return nil, err
	}

	if mw.clientComposite, err = NewComposite(mw); err != nil {
		return nil, err
	}
	mw.clientComposite.SetName("clientComposite")

	mw.clientComposite.children.observer = mw

	// This forces display of focus rectangles, as soon as the user starts to type.
	mw.SendMessage(WM_CHANGEUISTATE, UIS_INITIALIZE, 0)

	mw.TopLevelWindow.init()

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

func (mw *MainWindow) SetLayout(value Layout) error {
	if mw.clientComposite == nil {
		return newError("clientComposite not initialized")
	}

	return mw.clientComposite.SetLayout(value)
}

func (mw *MainWindow) ContextMenu() *Menu {
	return mw.clientComposite.ContextMenu()
}

func (mw *MainWindow) SetContextMenu(contextMenu *Menu) {
	mw.clientComposite.SetContextMenu(contextMenu)
}

func (mw *MainWindow) SaveState() error {
	if err := mw.clientComposite.SaveState(); err != nil {
		return err
	}

	return mw.TopLevelWindow.SaveState()
}

func (mw *MainWindow) RestoreState() error {
	if err := mw.clientComposite.RestoreState(); err != nil {
		return err
	}

	return mw.TopLevelWindow.RestoreState()
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

func (mw *MainWindow) SetVisible(visible bool) {
	if visible {
		DrawMenuBar(mw.hWnd)

		if mw.clientComposite.layout != nil {
			mw.clientComposite.layout.Update(false)
		}
	}

	mw.TopLevelWindow.SetVisible(visible)
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

func (mw *MainWindow) onInsertingWidget(index int, widget Widget) error {
	return mw.clientComposite.onInsertingWidget(index, widget)
}

func (mw *MainWindow) onInsertedWidget(index int, widget Widget) error {
	err := mw.clientComposite.onInsertedWidget(index, widget)
	if err == nil {
		minClientSize := mw.Layout().MinSize()
		clientSize := mw.clientComposite.Size()

		if clientSize.Width < minClientSize.Width || clientSize.Height < minClientSize.Height {
			mw.SetClientSize(minClientSize)
		}
	}

	return err
}

func (mw *MainWindow) onRemovingWidget(index int, widget Widget) error {
	return mw.clientComposite.onRemovingWidget(index, widget)
}

func (mw *MainWindow) onRemovedWidget(index int, widget Widget) error {
	return mw.clientComposite.onRemovedWidget(index, widget)
}

func (mw *MainWindow) onClearingWidgets() error {
	return mw.clientComposite.onClearingWidgets()
}

func (mw *MainWindow) onClearedWidgets() error {
	return mw.clientComposite.onClearedWidgets()
}

func (mw *MainWindow) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_SIZE, WM_SIZING:
		mw.toolBar.SendMessage(TB_AUTOSIZE, 0, 0)

		mw.clientComposite.SetBounds(mw.ClientBounds())
	}

	return mw.TopLevelWindow.WndProc(hwnd, msg, wParam, lParam)
}
