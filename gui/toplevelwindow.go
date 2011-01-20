// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"os"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

type CloseReason int

const (
	CloseReasonUnknown CloseReason = iota
	CloseReasonUser
)

type TopLevelWindow struct {
	Container
	owner            RootWidget
	clientArea       *Composite
	closingPublisher CloseEventPublisher
	closeReason      CloseReason
	prevFocusHWnd    HWND
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

func (tlw *TopLevelWindow) RunMessageLoop() (int, os.Error) {
	return tlw.runMessageLoop()
}

func (tlw *TopLevelWindow) Owner() RootWidget {
	return tlw.owner
}

func (tlw *TopLevelWindow) SetOwner(value RootWidget) os.Error {
	tlw.owner = value

	var ownerHWnd HWND
	if value != nil {
		ownerHWnd = value.Handle()
	}

	SetLastError(0)
	if 0 == SetWindowLong(tlw.hWnd, GWL_HWNDPARENT, int(ownerHWnd)) && GetLastError() != 0 {
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

func (tlw *TopLevelWindow) Closing() *CloseEvent {
	return tlw.closingPublisher.Event()
}

func (tlw *TopLevelWindow) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_ACTIVATE:
		switch LOWORD(uint(msg.WParam)) {
		case WA_ACTIVE, WA_CLICKACTIVE:
			if tlw.prevFocusHWnd != 0 {
				SetFocus(tlw.prevFocusHWnd)
			}

		case WA_INACTIVE:
			tlw.prevFocusHWnd = GetFocus()
		}
		return 0

	case WM_CLOSE:
		args := NewCloseEventArgs(widgetsByHWnd[tlw.hWnd], tlw.closeReason)
		tlw.closeReason = CloseReasonUnknown
		tlw.closingPublisher.Publish(args)
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
