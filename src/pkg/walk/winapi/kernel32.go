// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winapi

import (
	"syscall"
	"unsafe"
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
	libkernel32 uintptr

	// Functions
	fileTimeToSystemTime   uintptr
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
	systemTimeToFileTime   uintptr
)

type (
	ATOM      uint16
	HANDLE    uintptr
	HGLOBAL   HANDLE
	HINSTANCE HANDLE
	LCID      uint32
)

type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

type SYSTEMTIME struct {
	WYear         uint16
	WMonth        uint16
	WDayOfWeek    uint16
	WDay          uint16
	WHour         uint16
	WMinute       uint16
	WSecond       uint16
	WMilliseconds uint16
}

func init() {
	// Library
	libkernel32 = MustLoadLibrary("kernel32.dll")

	// Functions
	fileTimeToSystemTime = MustGetProcAddress(libkernel32, "FileTimeToSystemTime")
	getLastError = MustGetProcAddress(libkernel32, "GetLastError")
	getLogicalDriveStrings = MustGetProcAddress(libkernel32, "GetLogicalDriveStringsW")
	getModuleHandle = MustGetProcAddress(libkernel32, "GetModuleHandleW")
	getThreadLocale = MustGetProcAddress(libkernel32, "GetThreadLocale")
	globalAlloc = MustGetProcAddress(libkernel32, "GlobalAlloc")
	globalFree = MustGetProcAddress(libkernel32, "GlobalFree")
	globalLock = MustGetProcAddress(libkernel32, "GlobalLock")
	globalUnlock = MustGetProcAddress(libkernel32, "GlobalUnlock")
	moveMemory = MustGetProcAddress(libkernel32, "RtlMoveMemory")
	mulDiv = MustGetProcAddress(libkernel32, "MulDiv")
	setLastError = MustGetProcAddress(libkernel32, "SetLastError")
	systemTimeToFileTime = MustGetProcAddress(libkernel32, "SystemTimeToFileTime")
}

func FileTimeToSystemTime(lpFileTime *FILETIME, lpSystemTime *SYSTEMTIME) bool {
	ret, _, _ := syscall.Syscall(fileTimeToSystemTime, 2,
		uintptr(unsafe.Pointer(lpFileTime)),
		uintptr(unsafe.Pointer(lpSystemTime)),
		0)

	return ret != 0
}

func GetLastError() uint32 {
	ret, _, _ := syscall.Syscall(getLastError, 0,
		0,
		0,
		0)

	return uint32(ret)
}

func GetLogicalDriveStrings(nBufferLength uint32, lpBuffer *uint16) uint32 {
	ret, _, _ := syscall.Syscall(getLogicalDriveStrings, 2,
		uintptr(nBufferLength),
		uintptr(unsafe.Pointer(lpBuffer)),
		0)

	return uint32(ret)
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

func GlobalAlloc(uFlags uint32, dwBytes uintptr) HGLOBAL {
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

func MulDiv(nNumber, nNumerator, nDenominator int32) int32 {
	ret, _, _ := syscall.Syscall(mulDiv, 3,
		uintptr(nNumber),
		uintptr(nNumerator),
		uintptr(nDenominator))

	return int32(ret)
}

func SetLastError(dwErrorCode uint32) {
	syscall.Syscall(setLastError, 1,
		uintptr(dwErrorCode),
		0,
		0)
}

func SystemTimeToFileTime(lpSystemTime *SYSTEMTIME, lpFileTime *FILETIME) bool {
	ret, _, _ := syscall.Syscall(systemTimeToFileTime, 2,
		uintptr(unsafe.Pointer(lpSystemTime)),
		uintptr(unsafe.Pointer(lpFileTime)),
		0)

	return ret != 0
}
