// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crutches

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

// Registered message keys
const (
	_WM_RESIZE_KEY       = 0
	_WM_COMMAND_KEY      = 1
	_WM_CONTEXTMENU_KEY  = 2
	_WM_ITEMCHANGED_KEY  = 3
	_WM_ITEMACTIVATE_KEY = 4
	_WM_CLOSE_KEY        = 5
)

// Library
var (
	crutches uint32
)

// Functions
var (
	_getCustomMessage       uint32
	_getRegisteredMessageId uint32
	_registerWindowClass    uint32
)

// Registered message ids
var (
	closeMsgId        uint
	commandMsgId      uint
	contextMenuMsgId  uint
	itemActivateMsgId uint
	itemChangedMsgId  uint
	resizeMsgId       uint
)

func CloseMsgId() uint {
    return closeMsgId
}

func CommandMsgId() uint {
    return commandMsgId
}

func ContextMenuMsgId() uint {
    return contextMenuMsgId
}

func ItemActivateMsgId() uint {
    return itemActivateMsgId
}

func ItemChangedMsgId() uint {
    return itemChangedMsgId
}

func ResizeMsgId() uint {
    return resizeMsgId
}

type Message struct {
	Hwnd   HWND
	Msg    uint
	WParam uintptr
	LParam uintptr
}

func init() {
	// Library
	crutches = MustLoadLibrary("crutches.dll")

	// Functions
	_getCustomMessage = MustGetProcAddress(crutches, "GetCustomMessage@4")
	_getRegisteredMessageId = MustGetProcAddress(crutches, "GetRegisteredMessageId@4")
	_registerWindowClass = MustGetProcAddress(crutches, "RegisterWindowClass@4")

	resizeMsgId = getRegisteredMessageId(_WM_RESIZE_KEY)
	commandMsgId = getRegisteredMessageId(_WM_COMMAND_KEY)
	contextMenuMsgId = getRegisteredMessageId(_WM_CONTEXTMENU_KEY)
	itemChangedMsgId = getRegisteredMessageId(_WM_ITEMCHANGED_KEY)
	itemActivateMsgId = getRegisteredMessageId(_WM_ITEMACTIVATE_KEY)
	closeMsgId = getRegisteredMessageId(_WM_CLOSE_KEY)
}

func GetCustomMessage(msg *Message) int {
	ret, _, _ := syscall.Syscall(uintptr(_getCustomMessage),
		uintptr(unsafe.Pointer(msg)),
		0,
		0)

	return int(ret)
}

func getRegisteredMessageId(key uint) uint {
	ret, _, _ := syscall.Syscall(uintptr(_getRegisteredMessageId),
		uintptr(key),
		0,
		0)

	return uint(ret)
}

func RegisterWindowClass(hInstance HINSTANCE) ATOM {
	ret, _, _ := syscall.Syscall(uintptr(_registerWindowClass),
		uintptr(hInstance),
		0,
		0)

	return ATOM(ret)
}
