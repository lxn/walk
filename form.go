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

import (
	. "github.com/lxn/go-winapi"
)

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

type Form interface {
	Container
	AsFormBase() *FormBase
	Run() int
}

type FormBase struct {
	WindowBase
	form                  Form
	clientComposite       *Composite
	owner                 Form
	closingPublisher      CloseEventPublisher
	startingPublisher     EventPublisher
	titleChangedPublisher EventPublisher
	progressIndicator     *ProgressIndicator
	icon                  *Icon
	prevFocusHWnd         HWND
	isInRestoreState      bool
	closeReason           CloseReason
}

func (fb *FormBase) init(form Form) error {
	fb.form = form

	var err error
	if fb.clientComposite, err = NewComposite(form); err != nil {
		return err
	}
	fb.clientComposite.SetName("clientComposite")

	fb.clientComposite.children.observer = form.AsFormBase()

	fb.MustRegisterProperty("Title", NewProperty(
		func() interface{} {
			return fb.Title()
		},
		func(v interface{}) error {
			return fb.SetTitle(v.(string))
		},
		fb.titleChangedPublisher.Event()))

	return nil
}

func (fb *FormBase) AsContainerBase() *ContainerBase {
	return fb.clientComposite.AsContainerBase()
}

func (fb *FormBase) AsFormBase() *FormBase {
	return fb
}

func (fb *FormBase) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (fb *FormBase) SizeHint() Size {
	return fb.dialogBaseUnitsToPixels(Size{252, 218})
}

func (fb *FormBase) Children() *WidgetList {
	if fb.clientComposite == nil {
		return nil
	}

	return fb.clientComposite.Children()
}

func (fb *FormBase) Layout() Layout {
	if fb.clientComposite == nil {
		return nil
	}

	return fb.clientComposite.Layout()
}

func (fb *FormBase) SetLayout(value Layout) error {
	if fb.clientComposite == nil {
		return newError("clientComposite not initialized")
	}

	return fb.clientComposite.SetLayout(value)
}

func (fb *FormBase) DataBinder() *DataBinder {
	return fb.clientComposite.DataBinder()
}

func (fb *FormBase) SetDataBinder(db *DataBinder) {
	fb.clientComposite.SetDataBinder(db)
}

func (fb *FormBase) onInsertingWidget(index int, widget Widget) error {
	return fb.clientComposite.onInsertingWidget(index, widget)
}

func (fb *FormBase) onInsertedWidget(index int, widget Widget) error {
	err := fb.clientComposite.onInsertedWidget(index, widget)
	if err == nil {
		if layout := fb.Layout(); layout != nil {
			minClientSize := fb.Layout().MinSize()
			clientSize := fb.clientComposite.Size()

			if clientSize.Width < minClientSize.Width || clientSize.Height < minClientSize.Height {
				fb.SetClientSize(minClientSize)
			}
		}
	}

	return err
}

func (fb *FormBase) onRemovingWidget(index int, widget Widget) error {
	return fb.clientComposite.onRemovingWidget(index, widget)
}

func (fb *FormBase) onRemovedWidget(index int, widget Widget) error {
	return fb.clientComposite.onRemovedWidget(index, widget)
}

func (fb *FormBase) onClearingWidgets() error {
	return fb.clientComposite.onClearingWidgets()
}

func (fb *FormBase) onClearedWidgets() error {
	return fb.clientComposite.onClearedWidgets()
}

func (fb *FormBase) ContextMenu() *Menu {
	return fb.clientComposite.ContextMenu()
}

func (fb *FormBase) SetContextMenu(contextMenu *Menu) {
	fb.clientComposite.SetContextMenu(contextMenu)
}

func (fb *FormBase) Enabled() bool {
	return fb.enabled
}

func (fb *FormBase) SetEnabled(enabled bool) {
	fb.WindowBase.SetEnabled(enabled)
}

func (fb *FormBase) Font() *Font {
	if fb.font != nil {
		return fb.font
	}

	return defaultFont
}

func (fb *FormBase) Title() string {
	return windowText(fb.hWnd)
}

func (fb *FormBase) SetTitle(value string) error {
	return setWindowText(fb.hWnd, value)
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

func (fb *FormBase) Owner() Form {
	return fb.owner
}

func (fb *FormBase) SetOwner(value Form) error {
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
	fb.window.SetVisible(false)
}

func (fb *FormBase) Show() {
	if p, ok := fb.window.(Persistable); ok && p.Persistent() && appSingleton.settings != nil {
		p.RestoreState()
	}

	fb.window.SetVisible(true)
}

func (fb *FormBase) close() error {
	if p, ok := fb.window.(Persistable); ok && p.Persistent() && appSingleton.settings != nil {
		p.SaveState()
	}

	fb.window.Dispose()

	return nil
}

func (fb *FormBase) Close() error {
	fb.SendMessage(WM_CLOSE, 0, 0)

	return nil
}

func (fb *FormBase) Persistent() bool {
	return fb.clientComposite.persistent
}

func (fb *FormBase) SetPersistent(value bool) {
	fb.clientComposite.persistent = value
}

func (fb *FormBase) SaveState() error {
	if err := fb.clientComposite.SaveState(); err != nil {
		return err
	}

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

	return nil
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

	return fb.clientComposite.RestoreState()
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

	case WM_COMMAND:
		return fb.clientComposite.WndProc(hwnd, msg, wParam, lParam)

	case WM_GETMINMAXINFO:
		mmi := (*MINMAXINFO)(unsafe.Pointer(lParam))

		layout := fb.clientComposite.Layout()

		var min Size
		if layout != nil {
			min = fb.sizeFromClientSize(layout.MinSize())
		}

		mmi.PtMinTrackSize = POINT{
			int32(maxi(min.Width, fb.minSize.Width)),
			int32(maxi(min.Height, fb.minSize.Height)),
		}
		return 0

	case WM_NOTIFY:
		return fb.clientComposite.WndProc(hwnd, msg, wParam, lParam)

	case WM_SETTEXT:
		fb.titleChangedPublisher.Publish()

	case WM_SIZE, WM_SIZING:
		fb.clientComposite.SetBounds(fb.ClientBounds())

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

	return fb.WindowBase.WndProc(hwnd, msg, wParam, lParam)
}
