// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package advapi32

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/kernel32"
)

const KEY_READ REGSAM = 0x20019

const (
	HKEY_CLASSES_ROOT     HKEY = 0x80000000
	HKEY_CURRENT_USER     HKEY = 0x80000001
	HKEY_LOCAL_MACHINE    HKEY = 0x80000002
	HKEY_USERS            HKEY = 0x80000003
	HKEY_PERFORMANCE_DATA HKEY = 0x80000004
	HKEY_CURRENT_CONFIG   HKEY = 0x80000005
	HKEY_DYN_DATA         HKEY = 0x80000006
)

type (
	ACCESS_MASK uint
	HKEY        HANDLE
	REGSAM      ACCESS_MASK
)

var (
	// Library
	lib uintptr

	// Functions
	regCloseKey     uintptr
	regOpenKeyEx    uintptr
	regQueryValueEx uintptr
)

func init() {
	// Library
	lib = MustLoadLibrary("advapi32.dll")

	// Functions
	regCloseKey = MustGetProcAddress(lib, "RegCloseKey")
	regOpenKeyEx = MustGetProcAddress(lib, "RegOpenKeyExW")
	regQueryValueEx = MustGetProcAddress(lib, "RegQueryValueExW")
}

func RegCloseKey(hKey HKEY) int {
	ret, _, _ := syscall.Syscall(regCloseKey,
		uintptr(hKey),
		0,
		0)

	return int(ret)
}

func RegOpenKeyEx(hKey HKEY, lpSubKey *uint16, ulOptions uint, samDesired REGSAM, phkResult *HKEY) int {
	ret, _, _ := syscall.Syscall6(regOpenKeyEx,
		uintptr(hKey),
		uintptr(unsafe.Pointer(lpSubKey)),
		uintptr(ulOptions),
		uintptr(samDesired),
		uintptr(unsafe.Pointer(phkResult)),
		0)

	return int(ret)
}

func RegQueryValueEx(hKey HKEY, lpValueName *uint16, lpReserved, lpType *uint, lpData *byte, lpcbData *uint) int {
	ret, _, _ := syscall.Syscall6(regQueryValueEx,
		uintptr(hKey),
		uintptr(unsafe.Pointer(lpValueName)),
		uintptr(unsafe.Pointer(lpReserved)),
		uintptr(unsafe.Pointer(lpType)),
		uintptr(unsafe.Pointer(lpData)),
		uintptr(unsafe.Pointer(lpcbData)))

	return int(ret)
}
