// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"fmt"
	"os"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

type CloseReason int

const (
	CloseReasonUnknown CloseReason = iota
	CloseReasonUser
)

type ClosingEventArgs interface {
	CancelEventArgs
	Reason() CloseReason
}

type closingEventArgs struct {
	cancelEventArgs
	reason CloseReason
}

func (a *closingEventArgs) Reason() CloseReason {
	return a.reason
}

type ClosingEventHandler func(args ClosingEventArgs)


type TopLevelWindow struct {
	Container
	owner           *MainWindow
	clientArea      *Composite
	closingHandlers vector.Vector
	closeReason     CloseReason
}

func (tlw *TopLevelWindow) ClientArea() *Composite {
	return tlw.clientArea
}

func (tlw *TopLevelWindow) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz | ShrinkVert | GrowVert
}

func (tlw *TopLevelWindow) PreferredSize() drawing.Size {
	return tlw.dialogBaseUnitsToPixels(drawing.Size{252, 218})
}

func (tlw *TopLevelWindow) RunMessageLoop() os.Error {
	return tlw.runMessageLoop()
}

func (tlw *TopLevelWindow) Owner() *MainWindow {
	return tlw.owner
}

func (tlw *TopLevelWindow) SetOwner(value *MainWindow) os.Error {
	tlw.owner = value

	var ownerHWnd HWND
	if value != nil {
		ownerHWnd = value.hWnd
	}

	SetLastError(0)
	if 0 == SetWindowLong(tlw.hWnd, GWL_HWNDPARENT, int(ownerHWnd)) {
		return lastError("SetWindowLong")
	}

	return nil
}

func (tlw *TopLevelWindow) Hide() {
	ShowWindow(tlw.hWnd, SW_HIDE)
}

func (tlw *TopLevelWindow) Show() {
	ShowWindow(tlw.hWnd, SW_SHOW)
}

func (tlw *TopLevelWindow) close() os.Error {
	// FIXME: Remove this and children from widgetsByHWnd
	tlw.Dispose()

	return nil
}

func (tlw *TopLevelWindow) Close() os.Error {
	SendMessage(tlw.hWnd, WM_CLOSE, 0, 0)

	return nil
}

func (tlw *TopLevelWindow) SaveState() (string, os.Error) {
	var wp WINDOWPLACEMENT

	wp.Length = uint(unsafe.Sizeof(wp))

	if !GetWindowPlacement(tlw.hWnd, &wp) {
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

func (tlw *TopLevelWindow) RestoreState(s string) os.Error {
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

	if !SetWindowPlacement(tlw.hWnd, &wp) {
		return lastError("SetWindowPlacement")
	}

	return nil
}

func (tlw *TopLevelWindow) AddClosingHandler(handler ClosingEventHandler) {
	tlw.closingHandlers.Push(handler)
}

func (tlw *TopLevelWindow) RemoveClosingHandler(handler ClosingEventHandler) {
	for i, h := range tlw.closingHandlers {
		if h.(ClosingEventHandler) == handler {
			tlw.closingHandlers.Delete(i)
			break
		}
	}
}

func (tlw *TopLevelWindow) raiseClosing(args *closingEventArgs) {
	for _, handlerIface := range tlw.closingHandlers {
		handler := handlerIface.(ClosingEventHandler)
		handler(args)
	}
}

func (tlw *TopLevelWindow) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_CLOSE:
		args := &closingEventArgs{
			cancelEventArgs: cancelEventArgs{
				eventArgs: eventArgs{
					widgetsByHWnd[tlw.hWnd],
				},
			},
			reason: tlw.closeReason,
		}
		tlw.closeReason = CloseReasonUnknown
		tlw.raiseClosing(args)
		if !args.Canceled() {
			tlw.close()
		}
		return 0

	case WM_SYSCOMMAND:
		if msg.WParam == SC_CLOSE {
			tlw.closeReason = CloseReasonUser
		}
	}

	return tlw.Container.wndProc(msg, origWndProcPtr)
}
