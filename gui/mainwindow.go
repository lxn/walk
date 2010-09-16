// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

import (
    "walk/crutches"
	"walk/drawing"
	. "walk/winapi/comctl32"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

var (
	containerWindowClassAtom ATOM
)

type MainWindow struct {
	Container
	owner      *MainWindow
	menu       *Menu
	toolBar    *ToolBar
	clientArea *Composite
}

func ensureMainWindowInitialized() {
	if containerWindowClassAtom != 0 {
		return
	}

	hInst := GetHInstance()

	containerWindowClassAtom = crutches.RegisterWindowClass(hInst)
	if containerWindowClassAtom == 0 {
		panic("registerWindowClass for MainWindow window class failed.")
	}
}

func NewMainWindow() (mw *MainWindow, err os.Error) {
	ensureMainWindowInitialized()

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr("Container_WindowClass"), nil,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 300, 0, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	wnd := &MainWindow{Container: Container{Widget: Widget{hWnd: hWnd}}}

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

func (mw *MainWindow) ClientArea() *Composite {
	return mw.clientArea
}

func (*MainWindow) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz | ShrinkVert | GrowVert
}

func (mw *MainWindow) PreferredSize() drawing.Size {
	return mw.dialogBaseUnitsToPixels(drawing.Size{252, 218})
}

func (mw *MainWindow) RunMessageLoop() os.Error {
	return mw.runMessageLoop()
}

func (mw *MainWindow) Owner() *MainWindow {
	return mw.owner
}

func (mw *MainWindow) SetOwner(value *MainWindow) os.Error {
	mw.owner = value

	var ownerHWnd HWND
	if value != nil {
		ownerHWnd = value.hWnd
	}

	SetLastError(0)
	if 0 == SetWindowLong(mw.hWnd, GWL_HWNDPARENT, int(ownerHWnd)) {
		return lastError("SetWindowLong")
	}

	return nil
}

func (mw *MainWindow) Menu() *Menu {
	return mw.menu
}

func (mw *MainWindow) ToolBar() *ToolBar {
	return mw.toolBar
}

func (mw *MainWindow) ClientBounds() (bounds *drawing.Rectangle, err os.Error) {
	tlbBounds, err := mw.toolBar.Bounds()
	if err != nil {
		return
	}
	toolBarHeight := tlbBounds.Height

	bounds, err = mw.Widget.ClientBounds()
	if err != nil {
		return
	}
	bounds.Y += toolBarHeight
	bounds.Height -= toolBarHeight

	return
}

func (mw *MainWindow) Hide() {
	ShowWindow(mw.hWnd, SW_HIDE)
}

func (mw *MainWindow) Show() {
	ShowWindow(mw.hWnd, SW_SHOW)
}

func (mw *MainWindow) Close() (err os.Error) {
	// FIXME: Remove this and children from widgetsByHWnd
	mw.Dispose()

	return
}

func (mw *MainWindow) SaveState() (string, os.Error) {
	var wp WINDOWPLACEMENT

	wp.Length = uint(unsafe.Sizeof(wp))

	if !GetWindowPlacement(mw.hWnd, &wp) {
		return "", lastError("GetWindowPlacement")
	}

	return fmt.Sprint(
		wp.Flags, wp.ShowCmd,
		wp.PtMinPosition.X, wp.PtMinPosition.Y,
		wp.PtMaxPosition.X, wp.PtMaxPosition.Y,
		wp.RcNormalPosition.Left, wp.RcNormalPosition.Top,
		wp.RcNormalPosition.Right, wp.RcNormalPosition.Bottom),
		nil
}

func (mw *MainWindow) RestoreState(s string) os.Error {
	var wp WINDOWPLACEMENT

	_, err := fmt.Sscan(s,
		&wp.Flags, &wp.ShowCmd,
		&wp.PtMinPosition.X, &wp.PtMinPosition.Y,
		&wp.PtMaxPosition.X, &wp.PtMaxPosition.Y,
		&wp.RcNormalPosition.Left, &wp.RcNormalPosition.Top,
		&wp.RcNormalPosition.Right, &wp.RcNormalPosition.Bottom)
	if err != nil {
		return err
	}

	wp.Length = uint(unsafe.Sizeof(wp))

	if !SetWindowPlacement(mw.hWnd, &wp) {
		return lastError("SetWindowPlacement")
	}

	return nil
}

func (mw *MainWindow) raiseEvent(msg *MSG) (err os.Error) {
	switch msg.Message {
	case crutches.CloseMsgId():
		mw.Close()

		/*	case commandMsgId:
			switch HIWORD(DWORD(msg.WParam)) {
			case 0:
				// menu
				actionId := uint16(LOWORD(DWORD(msg.WParam)))
				if action, ok := actionsById[actionId]; ok {
					action.raiseTriggered()
				}
			}*/

	case crutches.ResizeMsgId():
		SendMessage(mw.toolBar.hWnd, TB_AUTOSIZE, 0, 0)
	}

	return mw.Container.raiseEvent(msg)
}
