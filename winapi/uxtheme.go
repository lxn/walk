// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winapi

import (
	"syscall"
	"unsafe"
)

var (
	// Library
	libuxtheme uintptr

	// Functions
	setWindowTheme uintptr
)

func init() {
	// Library
	libuxtheme = MustLoadLibrary("uxtheme.dll")

	// Functions
	setWindowTheme = MustGetProcAddress(libuxtheme, "SetWindowTheme")
}

func SetWindowTheme(hwnd HWND, pszSubAppName, pszSubIdList *uint16) HRESULT {
	ret, _, _ := syscall.Syscall(setWindowTheme, 3,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pszSubAppName)),
		uintptr(unsafe.Pointer(pszSubIdList)))

	return HRESULT(ret)
}
