// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kernel32

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
)

const MAX_PATH = 260

// Error codes
const (
	ERROR_SUCCESS             = 0
	ERROR_FILE_NOT_FOUND      = 2
	ERROR_INVALID_PARAMETER   = 87
	ERROR_INSUFFICIENT_BUFFER = 122
	ERROR_MORE_DATA           = 234
)

// GlobalAlloc flags
const (
	GHND          = 0x0042
	GMEM_FIXED    = 0x0000
	GMEM_MOVEABLE = 0x0002
	GMEM_ZEROINIT = 0x0040
	GPTR          = 0x004
)

var (
	// Library
	lib uintptr

	// Functions
	getLastError           uintptr
	getLogicalDriveStrings uintptr
	getModuleHandle        uintptr
	getThreadLocale        uintptr
	globalAlloc            uintptr
	globalFree             uintptr
	globalLock             uintptr
	globalUnlock           uintptr
	moveMemory             uintptr
	mulDiv                 uintptr
	setLastError           uintptr
)

type (
	ATOM      uint16
	HANDLE    uintptr
	HGLOBAL   HANDLE
	HINSTANCE HANDLE
	LCID      uint
)

func init() {
	// Library
	lib = MustLoadLibrary("kernel32.dll")

	// Functions
	getLastError = MustGetProcAddress(lib, "GetLastError")
	getLogicalDriveStrings = MustGetProcAddress(lib, "GetLogicalDriveStringsW")
	getModuleHandle = MustGetProcAddress(lib, "GetModuleHandleW")
	getThreadLocale = MustGetProcAddress(lib, "GetThreadLocale")
	globalAlloc = MustGetProcAddress(lib, "GlobalAlloc")
	globalFree = MustGetProcAddress(lib, "GlobalFree")
	globalLock = MustGetProcAddress(lib, "GlobalLock")
	globalUnlock = MustGetProcAddress(lib, "GlobalUnlock")
	moveMemory = MustGetProcAddress(lib, "RtlMoveMemory")
	mulDiv = MustGetProcAddress(lib, "MulDiv")
	setLastError = MustGetProcAddress(lib, "SetLastError")
}

func GetLastError() uint {
	ret, _, _ := syscall.Syscall(getLastError, 0,
		0,
		0,
		0)

	return uint(ret)
}

func GetLogicalDriveStrings(nBufferLength uint, lpBuffer *uint16) uint {
	ret, _, _ := syscall.Syscall(getLogicalDriveStrings, 2,
		uintptr(nBufferLength),
		uintptr(unsafe.Pointer(lpBuffer)),
		0)

	return uint(ret)
}

func GetModuleHandle(lpModuleName *uint16) HINSTANCE {
	ret, _, _ := syscall.Syscall(getModuleHandle, 1,
		uintptr(unsafe.Pointer(lpModuleName)),
		0,
		0)

	return HINSTANCE(ret)
}

func GetThreadLocale() LCID {
	ret, _, _ := syscall.Syscall(getThreadLocale, 0,
		0,
		0,
		0)

	return LCID(ret)
}

func GlobalAlloc(uFlags uint, dwBytes uintptr) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalAlloc, 2,
		uintptr(uFlags),
		dwBytes,
		0)

	return HGLOBAL(ret)
}

func GlobalFree(hMem HGLOBAL) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalFree, 1,
		uintptr(hMem),
		0,
		0)

	return HGLOBAL(ret)
}

func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalLock, 1,
		uintptr(hMem),
		0,
		0)

	return unsafe.Pointer(ret)
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := syscall.Syscall(globalUnlock, 1,
		uintptr(hMem),
		0,
		0)

	return ret != 0
}

func MoveMemory(destination, source unsafe.Pointer, length uintptr) {
	syscall.Syscall(moveMemory, 3,
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

func MulDiv(nNumber, nNumerator, nDenominator int) int {
	ret, _, _ := syscall.Syscall(mulDiv, 3,
		uintptr(nNumber),
		uintptr(nNumerator),
		uintptr(nDenominator))

	return int(ret)
}

func SetLastError(dwErrorCode uint) {
	syscall.Syscall(setLastError, 1,
		uintptr(dwErrorCode),
		0,
		0)
}
