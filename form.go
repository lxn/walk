// Copyright 2012 The Walk Authors. All rights reserved.
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

type CloseReason byte

const (
	CloseReasonUnknown CloseReason = iota
	CloseReasonUser
)

var syncFuncs struct {
	m     sync.Mutex
	funcs []func()
}

var syncMsgId uint32
var taskbarButtonCreatedMsgId uint32

func init() {
	syncMsgId = RegisterWindowMessage(syscall.StringToUTF16Ptr("WalkSync"))
	taskbarButtonCreatedMsgId = RegisterWindowMessage(syscall.StringToUTF16Ptr("TaskbarButtonCreated"))
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

type FormBase struct {
	ContainerBase
	owner                 RootWidget
	closingPublisher      CloseEventPublisher
	startingPublisher     EventPublisher
	titleChangedPublisher EventPublisher
	progressIndicator     *ProgressIndicator
	icon                  *Icon
	prevFocusHWnd         HWND
	isInRestoreState      bool
	closeReason           CloseReason
}

func (fb *FormBase) init() {
	fb.MustRegisterProperty("Title", NewProperty(
		func() interface{} {
			return fb.Title()
		},
		func(v interface{}) error {
			return fb.SetTitle(v.(string))
		},
		fb.titleChangedPublisher.Event()))
}

func (fb *FormBase) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (fb *FormBase) SizeHint() Size {
	return fb.dialogBaseUnitsToPixels(Size{252, 218})
}

func (fb *FormBase) Enabled() bool {
	return fb.enabled
}

func (fb *FormBase) SetEnabled(enabled bool) {
	fb.WidgetBase.SetEnabled(enabled)
}

func (fb *FormBase) Font() *Font {
	if fb.font != nil {
		return fb.font
	}

	return defaultFont
}

func (fb *FormBase) Title() string {
	return widgetText(fb.hWnd)
}

func (fb *FormBase) SetTitle(value string) error {
	return setWidgetText(fb.hWnd, value)
}

func (fb *FormBase) Run() int {
	fb.startingPublisher.Publish()

	var msg MSG

	for fb.hWnd != 0 {
		switch GetMessage(&msg, 0, 0, 0) {
		case 0:
			return int(msg.WParam)

		case -1:
			return -1
		}

		if !IsDialogMessage(fb.hWnd, &msg) {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}

		runSynchronized()
	}

	return 0
}

func (fb *FormBase) Starting() *Event {
	return fb.startingPublisher.Event()
}

func (fb *FormBase) Owner() RootWidget {
	return fb.owner
}

func (fb *FormBase) SetOwner(value RootWidget) error {
	fb.owner = value

	var ownerHWnd HWND
	if value != nil {
		ownerHWnd = value.Handle()
	}

	SetLastError(0)
	if 0 == SetWindowLong(
		fb.hWnd,
		GWL_HWNDPARENT,
		int32(ownerHWnd)) && GetLastError() != 0 {

		return lastError("SetWindowLong")
	}

	return nil
}

func (fb *FormBase) Icon() *Icon {
	return fb.icon
}

func (fb *FormBase) SetIcon(icon *Icon) {
	fb.icon = icon

	var hIcon uintptr
	if icon != nil {
		hIcon = uintptr(icon.hIcon)
	}

	fb.SendMessage(WM_SETICON, 0, hIcon)
	fb.SendMessage(WM_SETICON, 1, hIcon)
}

func (fb *FormBase) Hide() {
	fb.widget.SetVisible(false)
}

func (fb *FormBase) Show() {
	if p, ok := fb.widget.(Persistable); ok && p.Persistent() && appSingleton.settings != nil {
		p.RestoreState()
	}

	fb.widget.SetVisible(true)
}

func (fb *FormBase) close() error {
	if p, ok := fb.widget.(Persistable); ok && p.Persistent() && appSingleton.settings != nil {
		p.SaveState()
	}

	fb.widget.Dispose()

	return nil
}

func (fb *FormBase) Close() error {
	fb.SendMessage(WM_CLOSE, 0, 0)

	return nil
}

func (fb *FormBase) SaveState() error {
	var wp WINDOWPLACEMENT

	wp.Length = uint32(unsafe.Sizeof(wp))

	if !GetWindowPlacement(fb.hWnd, &wp) {
		return lastError("GetWindowPlacement")
	}

	state := fmt.Sprint(
		wp.Flags, wp.ShowCmd,
		wp.PtMinPosition.X, wp.PtMinPosition.Y,
		wp.PtMaxPosition.X, wp.PtMaxPosition.Y,
		wp.RcNormalPosition.Left, wp.RcNormalPosition.Top,
		wp.RcNormalPosition.Right, wp.RcNormalPosition.Bottom)

	if err := fb.putState(state); err != nil {
		return err
	}

	return fb.ContainerBase.SaveState()
}

func (fb *FormBase) RestoreState() error {
	if fb.isInRestoreState {
		return nil
	}
	fb.isInRestoreState = true
	defer func() {
		fb.isInRestoreState = false
	}()

	state, err := fb.getState()
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

	if !SetWindowPlacement(fb.hWnd, &wp) {
		return lastError("SetWindowPlacement")
	}

	if err := fb.ContainerBase.RestoreState(); err != nil {
		return err
	}

	return nil
}

func (fb *FormBase) Closing() *CloseEvent {
	return fb.closingPublisher.Event()
}

func (fb *FormBase) ProgressIndicator() *ProgressIndicator {
	return fb.progressIndicator
}

func (fb *FormBase) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_ACTIVATE:
		switch LOWORD(uint32(wParam)) {
		case WA_ACTIVE, WA_CLICKACTIVE:
			if fb.prevFocusHWnd != 0 {
				SetFocus(fb.prevFocusHWnd)
			}

		case WA_INACTIVE:
			fb.prevFocusHWnd = GetFocus()
		}
		return 0

	case WM_CLOSE:
		fb.closeReason = CloseReasonUnknown
		var canceled bool
		fb.closingPublisher.Publish(&canceled, fb.closeReason)
		if !canceled {
			if fb.owner != nil {
				fb.owner.SetEnabled(true)
				if !SetWindowPos(fb.owner.Handle(), HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_SHOWWINDOW) {
					lastError("SetWindowPos")
				}
			}

			fb.close()
		}
		return 0

	case WM_GETMINMAXINFO:
		mmi := (*MINMAXINFO)(unsafe.Pointer(lParam))

		var layout Layout
		if container, ok := fb.widget.(Container); ok {
			layout = container.Layout()
		}

		var min Size
		if layout != nil {
			min = fb.sizeFromClientSize(layout.MinSize())
		}

		mmi.PtMinTrackSize = POINT{
			int32(maxi(min.Width, fb.minSize.Width)),
			int32(maxi(min.Height, fb.minSize.Height)),
		}
		return 0

	case WM_SETTEXT:
		fb.titleChangedPublisher.Publish()

	case WM_SYSCOMMAND:
		if wParam == SC_CLOSE {
			fb.closeReason = CloseReasonUser
		}

	case taskbarButtonCreatedMsgId:
		version := GetVersion()
		major := version & 0xFF
		minor := version & 0xFF00 >> 8
		// Check that the OS is Win 7 or later (Win 7 is v6.1).
		if major > 6 || (major == 6 && minor > 0) {
			fb.progressIndicator, _ = newTaskbarList3(fb.hWnd)
		}
	}

	return fb.ContainerBase.WndProc(hwnd, msg, wParam, lParam)
}
