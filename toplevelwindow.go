// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

type CloseReason int

const (
	CloseReasonUnknown CloseReason = iota
	CloseReasonUser
)

var syncFuncs struct {
	m     sync.Mutex
	funcs []func()
}

var syncMsgId uint32

func init() {
	syncMsgId = RegisterWindowMessage(syscall.StringToUTF16Ptr("WalkSync"))
}

func synchronize(f func()) {
	syncFuncs.m.Lock()
	defer syncFuncs.m.Unlock()
	syncFuncs.funcs = append(syncFuncs.funcs, f)
}

func runSynchronized() {
	// Clear the list of callbacks first to avoid deadlock
	// if a callback itself calls Synchronize()...
	syncFuncs.m.Lock()
	funcs := syncFuncs.funcs
	syncFuncs.funcs = nil
	syncFuncs.m.Unlock()
	for _, f := range funcs {
		f()
	}
}

type TopLevelWindow struct {
	ContainerBase
	owner             RootWidget
	closingPublisher  CloseEventPublisher
	closeReason       CloseReason
	prevFocusHWnd     HWND
	isInRestoreState  bool
	startingPublisher EventPublisher
	icon              *Icon
}

func (tlw *TopLevelWindow) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (tlw *TopLevelWindow) SizeHint() Size {
	return tlw.dialogBaseUnitsToPixels(Size{252, 218})
}

func (tlw *TopLevelWindow) Title() string {
	return widgetText(tlw.hWnd)
}

func (tlw *TopLevelWindow) SetTitle(value string) error {
	return setWidgetText(tlw.hWnd, value)
}

func (tlw *TopLevelWindow) Run() int {
	tlw.startingPublisher.Publish()

	var msg MSG

	for tlw.hWnd != 0 {
		switch GetMessage(&msg, 0, 0, 0) {
		case 0:
			return int(msg.WParam)

		case -1:
			return -1
		}

		if !IsDialogMessage(tlw.hWnd, &msg) {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}

		runSynchronized()
	}

	return 0
}

func (tlw *TopLevelWindow) Starting() *Event {
	return tlw.startingPublisher.Event()
}

func (tlw *TopLevelWindow) Owner() RootWidget {
	return tlw.owner
}

func (tlw *TopLevelWindow) SetOwner(value RootWidget) error {
	tlw.owner = value

	var ownerHWnd HWND
	if value != nil {
		ownerHWnd = value.BaseWidget().hWnd
	}

	SetLastError(0)
	if 0 == SetWindowLong(
		tlw.hWnd,
		GWL_HWNDPARENT,
		int32(ownerHWnd)) && GetLastError() != 0 {

		return lastError("SetWindowLong")
	}

	return nil
}

func (tlw *TopLevelWindow) Icon() *Icon {
	return tlw.icon
}

func (tlw *TopLevelWindow) SetIcon(icon *Icon) {
	tlw.icon = icon

	var hIcon uintptr
	if icon != nil {
		hIcon = uintptr(icon.hIcon)
	}

	SendMessage(tlw.hWnd, WM_SETICON, 0, hIcon)
	SendMessage(tlw.hWnd, WM_SETICON, 1, hIcon)
}

func (tlw *TopLevelWindow) Hide() {
	tlw.SetVisible(false)
}

func (tlw *TopLevelWindow) Show() {
	tlw.SetVisible(true)
}

func (tlw *TopLevelWindow) close() error {
	tlw.Dispose()

	return nil
}

func (tlw *TopLevelWindow) Close() error {
	SendMessage(tlw.hWnd, WM_CLOSE, 0, 0)

	return nil
}

func (tlw *TopLevelWindow) SaveState() error {
	var wp WINDOWPLACEMENT

	wp.Length = uint32(unsafe.Sizeof(wp))

	if !GetWindowPlacement(tlw.hWnd, &wp) {
		return lastError("GetWindowPlacement")
	}

	state := fmt.Sprint(
		wp.Flags, wp.ShowCmd,
		wp.PtMinPosition.X, wp.PtMinPosition.Y,
		wp.PtMaxPosition.X, wp.PtMaxPosition.Y,
		wp.RcNormalPosition.Left, wp.RcNormalPosition.Top,
		wp.RcNormalPosition.Right, wp.RcNormalPosition.Bottom)

	if err := tlw.putState(state); err != nil {
		return err
	}

	return tlw.ContainerBase.SaveState()
}

func (tlw *TopLevelWindow) RestoreState() error {
	if tlw.isInRestoreState {
		return nil
	}
	tlw.isInRestoreState = true
	defer func() {
		tlw.isInRestoreState = false
	}()

	state, err := tlw.getState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	var wp WINDOWPLACEMENT

	if _, err := fmt.Sscan(state,
		&wp.Flags, &wp.ShowCmd,
		&wp.PtMinPosition.X, &wp.PtMinPosition.Y,
		&wp.PtMaxPosition.X, &wp.PtMaxPosition.Y,
		&wp.RcNormalPosition.Left, &wp.RcNormalPosition.Top,
		&wp.RcNormalPosition.Right, &wp.RcNormalPosition.Bottom); err != nil {
		return err
	}

	wp.Length = uint32(unsafe.Sizeof(wp))

	if !SetWindowPlacement(tlw.hWnd, &wp) {
		return lastError("SetWindowPlacement")
	}

	if err := tlw.ContainerBase.RestoreState(); err != nil {
		return err
	}

	return nil
}

func (tlw *TopLevelWindow) Closing() *CloseEvent {
	return tlw.closingPublisher.Event()
}

func (tlw *TopLevelWindow) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_ACTIVATE:
		switch LOWORD(uint32(wParam)) {
		case WA_ACTIVE, WA_CLICKACTIVE:
			if tlw.prevFocusHWnd != 0 {
				SetFocus(tlw.prevFocusHWnd)
			}

		case WA_INACTIVE:
			tlw.prevFocusHWnd = GetFocus()
		}
		return 0

	case WM_CLOSE:
		tlw.closeReason = CloseReasonUnknown
		var canceled bool
		tlw.closingPublisher.Publish(&canceled, tlw.closeReason)
		if !canceled {
			if tlw.owner != nil {
				tlw.owner.SetEnabled(true)
				if !SetWindowPos(tlw.owner.BaseWidget().hWnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_SHOWWINDOW) {
					lastError("SetWindowPos")
				}
			}

			tlw.close()
		}
		return 0

	case WM_GETMINMAXINFO:
		mmi := (*MINMAXINFO)(unsafe.Pointer(lParam))

		var layout Layout
		if container, ok := tlw.widget.(Container); ok {
			layout = container.Layout()
		}

		var min Size
		if layout != nil {
			min = tlw.sizeFromClientSize(layout.MinSize())
		}

		mmi.PtMinTrackSize = POINT{
			int32(maxi(min.Width, tlw.minSize.Width)),
			int32(maxi(min.Height, tlw.minSize.Height)),
		}
		return 0

	case WM_SYSCOMMAND:
		if wParam == SC_CLOSE {
			tlw.closeReason = CloseReasonUser
		}
	}

	return tlw.ContainerBase.wndProc(hwnd, msg, wParam, lParam)
}
