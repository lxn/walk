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
	lib uint32

	// Functions
	getLastError    uint32
	getModuleHandle uint32
	globalAlloc     uint32
	globalFree      uint32
	globalLock      uint32
	globalUnlock    uint32
	moveMemory      uint32
	mulDiv          uint32
	setLastError    uint32
)

type (
	ATOM      uint16
	HANDLE    uintptr
	HGLOBAL   HANDLE
	HINSTANCE HANDLE
)

func init() {
	// Library
	lib = MustLoadLibrary("kernel32.dll")

	// Functions
	getLastError = MustGetProcAddress(lib, "GetLastError")
	getModuleHandle = MustGetProcAddress(lib, "GetModuleHandleW")
	globalAlloc = MustGetProcAddress(lib, "GlobalAlloc")
	globalFree = MustGetProcAddress(lib, "GlobalFree")
	globalLock = MustGetProcAddress(lib, "GlobalLock")
	globalUnlock = MustGetProcAddress(lib, "GlobalUnlock")
	moveMemory = MustGetProcAddress(lib, "RtlMoveMemory")
	mulDiv = MustGetProcAddress(lib, "MulDiv")
	setLastError = MustGetProcAddress(lib, "SetLastError")
}

func GetLastError() uint {
	ret, _, _ := syscall.Syscall(uintptr(setLastError),
		0,
		0,
		0)

	return uint(ret)
}

func GetHInstance() HINSTANCE {
	ret, _, _ := syscall.Syscall(uintptr(getModuleHandle),
		0,
		0,
		0)

	return HINSTANCE(ret)
}

func GlobalAlloc(uFlags uint, dwBytes uintptr) HGLOBAL {
	ret, _, _ := syscall.Syscall(uintptr(globalAlloc),
		uintptr(uFlags),
		dwBytes,
		0)

	return HGLOBAL(ret)
}

func GlobalFree(hMem HGLOBAL) HGLOBAL {
	ret, _, _ := syscall.Syscall(uintptr(globalFree),
		uintptr(hMem),
		0,
		0)

	return HGLOBAL(ret)
}

func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(uintptr(globalLock),
		uintptr(hMem),
		0,
		0)

	return unsafe.Pointer(ret)
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := syscall.Syscall(uintptr(globalUnlock),
		uintptr(hMem),
		0,
		0)

	return ret != 0
}

func MoveMemory(destination, source unsafe.Pointer, length uintptr) {
	syscall.Syscall(uintptr(moveMemory),
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

func MulDiv(nNumber, nNumerator, nDenominator int) int {
	ret, _, _ := syscall.Syscall(uintptr(mulDiv),
		uintptr(nNumber),
		uintptr(nNumerator),
		uintptr(nDenominator))

	return int(ret)
}

func SetLastError(dwErrorCode uint) {
	syscall.Syscall(uintptr(setLastError),
		uintptr(dwErrorCode),
		0,
		0)
}
