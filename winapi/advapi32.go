// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winapi

import (
	"syscall"
	"unsafe"
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
	ACCESS_MASK uint32
	HKEY        HANDLE
	REGSAM      ACCESS_MASK
)

var (
	// Library
	libadvapi32 uintptr

	// Functions
	regCloseKey     uintptr
	regOpenKeyEx    uintptr
	regQueryValueEx uintptr
)

func init() {
	// Library
	libadvapi32 = MustLoadLibrary("advapi32.dll")

	// Functions
	regCloseKey = MustGetProcAddress(libadvapi32, "RegCloseKey")
	regOpenKeyEx = MustGetProcAddress(libadvapi32, "RegOpenKeyExW")
	regQueryValueEx = MustGetProcAddress(libadvapi32, "RegQueryValueExW")
}

func RegCloseKey(hKey HKEY) int32 {
	ret, _, _ := syscall.Syscall(regCloseKey, 1,
		uintptr(hKey),
		0,
		0)

	return int32(ret)
}

func RegOpenKeyEx(hKey HKEY, lpSubKey *uint16, ulOptions uint32, samDesired REGSAM, phkResult *HKEY) int32 {
	ret, _, _ := syscall.Syscall6(regOpenKeyEx, 5,
		uintptr(hKey),
		uintptr(unsafe.Pointer(lpSubKey)),
		uintptr(ulOptions),
		uintptr(samDesired),
		uintptr(unsafe.Pointer(phkResult)),
		0)

	return int32(ret)
}

func RegQueryValueEx(hKey HKEY, lpValueName *uint16, lpReserved, lpType *uint32, lpData *byte, lpcbData *uint32) int32 {
	ret, _, _ := syscall.Syscall6(regQueryValueEx, 6,
		uintptr(hKey),
		uintptr(unsafe.Pointer(lpValueName)),
		uintptr(unsafe.Pointer(lpReserved)),
		uintptr(unsafe.Pointer(lpType)),
		uintptr(unsafe.Pointer(lpData)),
		uintptr(unsafe.Pointer(lpcbData)))

	return int32(ret)
}
