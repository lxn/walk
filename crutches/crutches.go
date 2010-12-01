// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crutches

import (
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

// Functions
var (
	_getRegisteredMessageId uint32
	_registerWindowMessage  uint32
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

// Internal functions defined in .c files
func initqueue()
func initcrutch()
func getCustomMessage(msgPointer uintptr) uintptr
func registerWindowClass(hInstance uintptr) uintptr
func getRegisteredMessageId(key uint) uint

func init() {
	initcrutch()
	initqueue()

	// Library
	user32 := MustLoadLibrary("user32.dll")

	// Functions
	_registerWindowMessage = MustGetProcAddress(user32, "RegisterWindowMessageW")

	resizeMsgId = getRegisteredMessageId(_WM_RESIZE_KEY)
	// resizeMsgId = getRegisteredMessage2("resize_0b0f95e6-7ef7-4767-b484-940e7a3cf4f1")
	commandMsgId = getRegisteredMessageId(_WM_COMMAND_KEY)
	contextMenuMsgId = getRegisteredMessageId(_WM_CONTEXTMENU_KEY)
	itemChangedMsgId = getRegisteredMessageId(_WM_ITEMCHANGED_KEY)
	itemActivateMsgId = getRegisteredMessageId(_WM_ITEMACTIVATE_KEY)
	closeMsgId = getRegisteredMessageId(_WM_CLOSE_KEY)
}

func GetCustomMessage(msg *Message) int {
	ret := getCustomMessage(uintptr(unsafe.Pointer(msg)))
	return int(ret)
}
/*
func getRegisteredMessage2(msgid string) uint {
    ret, _, _ := syscall.Syscall(uintptr(_registerWindowMessage),
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msgid))),
        0,
        0)
    return uint(ret)
}
*/

func RegisterWindowClass(hInstance HINSTANCE) ATOM {
	ret := registerWindowClass(uintptr(hInstance))
	return ATOM(ret)
}
