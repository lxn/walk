// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uxtheme

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/user32"
)

var (
	// Library
	lib uint32

	// Functions
	setWindowTheme uint32
)

func init() {
	// Library
	lib = MustLoadLibrary("uxtheme.dll")

	// Functions
	setWindowTheme = MustGetProcAddress(lib, "SetWindowTheme")
}

func SetWindowTheme(hwnd HWND, pszSubAppName, pszSubIdList *uint16) HRESULT {
	ret, _, _ := syscall.Syscall(uintptr(setWindowTheme),
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pszSubAppName)),
		uintptr(unsafe.Pointer(pszSubIdList)))

	return HRESULT(ret)
}
