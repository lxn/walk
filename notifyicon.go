// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

const notifyIconWindowClass = `\o/ Walk_NotifyIcon_Class \o/`

func init() {
	MustRegisterWindowClass(notifyIconWindowClass)
}

func notifyIconWndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	// Retrieve our *NotifyIcon from the message window.
	ptr := GetWindowLongPtr(hwnd, GWLP_USERDATA)
	ni := (*NotifyIcon)(unsafe.Pointer(ptr))

	switch lParam {
	case WM_LBUTTONDOWN:
		ni.publishMouseEvent(&ni.mouseDownPublisher, LeftButton)

	case WM_LBUTTONUP:
		ni.publishMouseEvent(&ni.mouseUpPublisher, LeftButton)

	case WM_RBUTTONDOWN:
		ni.publishMouseEvent(&ni.mouseDownPublisher, RightButton)

	case WM_RBUTTONUP:
		ni.publishMouseEvent(&ni.mouseUpPublisher, RightButton)

		SendMessage(hwnd, msg, wParam, WM_CONTEXTMENU)

	case WM_CONTEXTMENU:
		if ni.contextMenu.Actions().Len() == 0 {
			break
		}

		SetForegroundWindow(hwnd)

		var p POINT
		if !GetCursorPos(&p) {
			lastError("GetCursorPos")
		}

		actionId := uint16(TrackPopupMenuEx(
			ni.contextMenu.hMenu,
			TPM_NOANIMATION|TPM_RETURNCMD,
			p.X,
			p.Y,
			hwnd,
			nil))
		if actionId != 0 {
			if action, ok := actionsById[actionId]; ok {
				action.raiseTriggered()
			}
		}

		return 0
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

// NotifyIcon represents an icon in the taskbar notification area.
type NotifyIcon struct {
	id                 uint32
	hWnd               HWND
	contextMenu        *Menu
	icon               *Icon
	toolTip            string
	visible            bool
	mouseDownPublisher MouseEventPublisher
	mouseUpPublisher   MouseEventPublisher
}

// NewNotifyIcon creates and returns a new NotifyIcon.
//
// The NotifyIcon is initially not visible.
func NewNotifyIcon() (*NotifyIcon, error) {
	// Create the message-only window for the NotifyIcon.
	hWnd := CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr(notifyIconWindowClass),
		nil,
		0,
		0,
		0,
		0,
		0,
		HWND_MESSAGE,
		0,
		0,
		nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	// Add our notify icon to the status area and make sure it is hidden.
	nid := NOTIFYICONDATA{
		HWnd:             hWnd,
		UFlags:           NIF_MESSAGE | NIF_STATE,
		DwState:          NIS_HIDDEN,
		DwStateMask:      NIS_HIDDEN,
		UCallbackMessage: notifyIconMessageId,
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	if !Shell_NotifyIcon(NIM_ADD, &nid) {
		return nil, newError("Shell_NotifyIcon")
	}

	// We want XP-compatible message behavior.
	nid.UVersion = NOTIFYICON_VERSION

	if !Shell_NotifyIcon(NIM_SETVERSION, &nid) {
		return nil, newError("Shell_NotifyIcon")
	}

	// Create and initialize the NotifyIcon already.
	menu, err := NewMenu()
	if err != nil {
		return nil, err
	}

	ni := &NotifyIcon{
		id:          nid.UID,
		hWnd:        hWnd,
		contextMenu: menu,
	}

	// Set our *NotifyIcon as user data for the message window.
	SetWindowLongPtr(hWnd, GWLP_USERDATA, uintptr(unsafe.Pointer(ni)))

	return ni, nil
}

func (ni *NotifyIcon) notifyIconData() *NOTIFYICONDATA {
	nid := &NOTIFYICONDATA{
		UID:  ni.id,
		HWnd: ni.hWnd,
	}
	nid.CbSize = uint32(unsafe.Sizeof(*nid))

	return nid
}

// Dispose releases the operating system resources associated with the
// NotifyIcon.
//
// The associated Icon is not disposed of.
func (ni *NotifyIcon) Dispose() error {
	if ni.hWnd == 0 {
		return nil
	}

	nid := ni.notifyIconData()

	if !Shell_NotifyIcon(NIM_DELETE, nid) {
		return newError("Shell_NotifyIcon")
	}

	if !DestroyWindow(ni.hWnd) {
		return lastError("DestroyWindow")
	}
	ni.hWnd = 0

	return nil
}

func (ni *NotifyIcon) showMessage(title, info string, iconType uint32) error {
	nid := ni.notifyIconData()
	nid.UFlags = NIF_INFO
	nid.DwInfoFlags = iconType
	copy(nid.SzInfoTitle[:], syscall.StringToUTF16(title))
	copy(nid.SzInfo[:], syscall.StringToUTF16(info))

	if !Shell_NotifyIcon(NIM_MODIFY, nid) {
		return newError("Shell_NotifyIcon")
	}

	return nil
}

// ShowMessage displays a neutral message balloon above the NotifyIcon.
//
// The NotifyIcon must be visible before calling this method.
func (ni *NotifyIcon) ShowMessage(title, info string) error {
	return ni.showMessage(title, info, NIIF_NONE)
}

// ShowInfo displays an info message balloon above the NotifyIcon.
//
// The NotifyIcon must be visible before calling this method.
func (ni *NotifyIcon) ShowInfo(title, info string) error {
	return ni.showMessage(title, info, NIIF_INFO)
}

// ShowWarning displays a warning message balloon above the NotifyIcon.
//
// The NotifyIcon must be visible before calling this method.
func (ni *NotifyIcon) ShowWarning(title, info string) error {
	return ni.showMessage(title, info, NIIF_WARNING)
}

// ShowError displays an error message balloon above the NotifyIcon.
//
// The NotifyIcon must be visible before calling this method.
func (ni *NotifyIcon) ShowError(title, info string) error {
	return ni.showMessage(title, info, NIIF_ERROR)
}

// ShowCustom displays a custom icon message balloon above the NotifyIcon.
//
// The NotifyIcon must be visible before calling this method.
func (ni *NotifyIcon) ShowCustom(title, info string) error {
	return ni.showMessage(title, info, NIIF_USER)
}

// ContextMenu returns the context menu of the NotifyIcon.
func (ni *NotifyIcon) ContextMenu() *Menu {
	return ni.contextMenu
}

// Icon returns the Icon of the NotifyIcon.
func (ni *NotifyIcon) Icon() *Icon {
	return ni.icon
}

// SetIcon sets the Icon of the NotifyIcon.
func (ni *NotifyIcon) SetIcon(icon *Icon) error {
	if icon == ni.icon {
		return nil
	}

	nid := ni.notifyIconData()
	nid.UFlags = NIF_ICON
	if icon == nil {
		nid.HIcon = 0
	} else {
		nid.HIcon = icon.hIcon
	}

	if !Shell_NotifyIcon(NIM_MODIFY, nid) {
		return newError("Shell_NotifyIcon")
	}

	ni.icon = icon

	return nil
}

// ToolTip returns the tool tip text of the NotifyIcon.
func (ni *NotifyIcon) ToolTip() string {
	return ni.toolTip
}

// SetToolTip sets the tool tip text of the NotifyIcon.
func (ni *NotifyIcon) SetToolTip(toolTip string) error {
	if toolTip == ni.toolTip {
		return nil
	}

	nid := ni.notifyIconData()
	nid.UFlags = NIF_TIP
	copy(nid.SzTip[:], syscall.StringToUTF16(toolTip))

	if !Shell_NotifyIcon(NIM_MODIFY, nid) {
		return newError("Shell_NotifyIcon")
	}

	ni.toolTip = toolTip

	return nil
}

// Visible returns if the NotifyIcon is visible.
func (ni *NotifyIcon) Visible() bool {
	return ni.visible
}

// SetVisible sets if the NotifyIcon is visible.
func (ni *NotifyIcon) SetVisible(visible bool) error {
	if visible == ni.visible {
		return nil
	}

	nid := ni.notifyIconData()
	nid.UFlags = NIF_STATE
	nid.DwStateMask = NIS_HIDDEN
	if !visible {
		nid.DwState = NIS_HIDDEN
	}

	if !Shell_NotifyIcon(NIM_MODIFY, nid) {
		return newError("Shell_NotifyIcon")
	}

	ni.visible = visible

	return nil
}

func (ni *NotifyIcon) publishMouseEvent(publisher *MouseEventPublisher, button MouseButton) {
	var p POINT
	if !GetCursorPos(&p) {
		lastError("GetCursorPos")
	}

	publisher.Publish(int(p.X), int(p.Y), button)
}

// MouseDown returns the event that is published when a mouse button is pressed
// while the cursor is over the NotifyIcon.
func (ni *NotifyIcon) MouseDown() *MouseEvent {
	return ni.mouseDownPublisher.Event()
}

// MouseDown returns the event that is published when a mouse button is released
// while the cursor is over the NotifyIcon.
func (ni *NotifyIcon) MouseUp() *MouseEvent {
	return ni.mouseUpPublisher.Event()
}
