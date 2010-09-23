// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	FALSE = 0
	TRUE  = 1
)

type (
	BOOL    int
	HRESULT int32
)

func MustLoadLibrary(name string) uint32 {
	lib, errno := syscall.LoadLibrary(name)
	if errno != 0 {
		panic(fmt.Sprintf(`syscall.LoadLibrary("%s") failed: %s`, name, syscall.Errstr(errno)))
	}

	return lib
}

func MustGetProcAddress(lib uint32, name string) uint32 {
	addr, errno := syscall.GetProcAddress(lib, name)
	if errno != 0 {
		panic(fmt.Sprintf(`syscall.GetProcAddress(%d, "%s") failed: %s`, lib, name, syscall.Errstr(errno)))
	}

	return addr
}

func SUCCEEDED(hr HRESULT) bool {
	return hr >= 0
}

func FAILED(hr HRESULT) bool {
	return hr < 0
}

func MAKELONG(lo, hi uint16) uint {
	return uint(uint(lo) | ((uint(hi)) << 16))
}

func LOWORD(dw uint) uint16 {
	return uint16(dw)
}

func HIWORD(dw uint) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func UTF16PtrToString(s *uint16) string {
	return syscall.UTF16ToString((*[1 << 30]uint16)(unsafe.Pointer(s))[0:])
}

func BoolToBOOL(value bool) BOOL {
	if value {
		return 1
	}

	return 0
}
