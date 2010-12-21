// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oleaut32

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
)

type VARTYPE uint16

const (
	VT_EMPTY            VARTYPE = 0
	VT_NULL             VARTYPE = 1
	VT_I2               VARTYPE = 2
	VT_I4               VARTYPE = 3
	VT_R4               VARTYPE = 4
	VT_R8               VARTYPE = 5
	VT_CY               VARTYPE = 6
	VT_DATE             VARTYPE = 7
	VT_BSTR             VARTYPE = 8
	VT_DISPATCH         VARTYPE = 9
	VT_ERROR            VARTYPE = 10
	VT_BOOL             VARTYPE = 11
	VT_VARIANT          VARTYPE = 12
	VT_UNKNOWN          VARTYPE = 13
	VT_DECIMAL          VARTYPE = 14
	VT_I1               VARTYPE = 16
	VT_UI1              VARTYPE = 17
	VT_UI2              VARTYPE = 18
	VT_UI4              VARTYPE = 19
	VT_I8               VARTYPE = 20
	VT_UI8              VARTYPE = 21
	VT_INT              VARTYPE = 22
	VT_UINT             VARTYPE = 23
	VT_VOID             VARTYPE = 24
	VT_HRESULT          VARTYPE = 25
	VT_PTR              VARTYPE = 26
	VT_SAFEARRAY        VARTYPE = 27
	VT_CARRAY           VARTYPE = 28
	VT_USERDEFINED      VARTYPE = 29
	VT_LPSTR            VARTYPE = 30
	VT_LPWSTR           VARTYPE = 31
	VT_RECORD           VARTYPE = 36
	VT_INT_PTR          VARTYPE = 37
	VT_UINT_PTR         VARTYPE = 38
	VT_FILETIME         VARTYPE = 64
	VT_BLOB             VARTYPE = 65
	VT_STREAM           VARTYPE = 66
	VT_STORAGE          VARTYPE = 67
	VT_STREAMED_OBJECT  VARTYPE = 68
	VT_STORED_OBJECT    VARTYPE = 69
	VT_BLOB_OBJECT      VARTYPE = 70
	VT_CF               VARTYPE = 71
	VT_CLSID            VARTYPE = 72
	VT_VERSIONED_STREAM VARTYPE = 73
	VT_BSTR_BLOB        VARTYPE = 0xfff
	VT_VECTOR           VARTYPE = 0x1000
	VT_ARRAY            VARTYPE = 0x2000
	VT_BYREF            VARTYPE = 0x4000
	VT_RESERVED         VARTYPE = 0x8000
	VT_ILLEGAL          VARTYPE = 0xffff
	VT_ILLEGALMASKED    VARTYPE = 0xfff
	VT_TYPEMASK         VARTYPE = 0xfff
)

type VARIANT_BOOL int16
//type BSTR *uint16

func StringToBSTR(value string) *uint16 /*BSTR*/ {
	// IMPORTANT: Don't forget to free the BSTR value when no longer needed!
	return SysAllocString(value)
}

func BSTRToString(value *uint16 /*BSTR*/ ) string {
	// ISSUE: Is this really ok?
	bstrArrPtr := (*[2000000000]uint16)(unsafe.Pointer(value))

	bstrSlice := make([]uint16, SysStringLen(value))
	copy(bstrSlice, bstrArrPtr[:])

	return syscall.UTF16ToString(bstrSlice)
}

type VAR_I4 struct {
	vt        VARTYPE
	reserved1 [6]byte
	lVal      int
	reserved2 [4]byte
}

func IntToVariantI4(value int) *VAR_I4 {
	return &VAR_I4{vt: VT_I4, lVal: value}
}

func VariantI4ToInt(value *VAR_I4) int {
	return value.lVal
}

type VAR_BOOL struct {
	vt        VARTYPE
	reserved1 [6]byte
	boolVal   VARIANT_BOOL
	reserved2 [6]byte
}

func BoolToVariantBool(value bool) *VAR_BOOL {
	return &VAR_BOOL{vt: VT_BOOL, boolVal: VARIANT_BOOL(BoolToBOOL(value))}
}

func VariantBoolToBool(value *VAR_BOOL) bool {
	return value.boolVal != 0
}

type VAR_BSTR struct {
	vt        VARTYPE
	reserved1 [6]byte
	bstrVal   *uint16 /*BSTR*/
	reserved2 [4]byte // 32-bit specific
}

func StringToVariantBSTR(value string) *VAR_BSTR {
	// IMPORTANT: Don't forget to free the BSTR value when no longer needed!
	return &VAR_BSTR{vt: VT_BSTR, bstrVal: StringToBSTR(value)}
}

func VariantBSTRToString(value *VAR_BSTR) string {
	return BSTRToString(value.bstrVal)
}

var (
	// Library
	lib uint32

	// Functions
	sysAllocString uint32
	sysFreeString  uint32
	sysStringLen   uint32
)

func init() {
	// Library
	lib = MustLoadLibrary("oleaut32.dll")

	// Functions
	sysAllocString = MustGetProcAddress(lib, "SysAllocString")
	sysFreeString = MustGetProcAddress(lib, "SysFreeString")
	sysStringLen = MustGetProcAddress(lib, "SysStringLen")
}

func SysAllocString(s string) *uint16 /*BSTR*/ {
	ret, _, _ := syscall.Syscall(uintptr(sysAllocString),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s))),
		0,
		0)

	return (*uint16) /*BSTR*/ (unsafe.Pointer(ret))
}

func SysFreeString(bstr *uint16 /*BSTR*/ ) {
	syscall.Syscall(uintptr(sysFreeString),
		uintptr(unsafe.Pointer(bstr)),
		0,
		0)
}

func SysStringLen(bstr *uint16 /*BSTR*/ ) uint {
	ret, _, _ := syscall.Syscall(uintptr(sysStringLen),
		uintptr(unsafe.Pointer(bstr)),
		0,
		0)

	return uint(ret)
}
