// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winapi


import (
	"syscall"
	"unsafe"
)

// EnumPrinters flags
const (
	PRINTER_ENUM_DEFAULT     = 0x00000001
	PRINTER_ENUM_LOCAL       = 0x00000002
	PRINTER_ENUM_CONNECTIONS = 0x00000004
	PRINTER_ENUM_FAVORITE    = 0x00000004
	PRINTER_ENUM_NAME        = 0x00000008
	PRINTER_ENUM_REMOTE      = 0x00000010
	PRINTER_ENUM_SHARED      = 0x00000020
	PRINTER_ENUM_NETWORK     = 0x00000040
)

type PRINTER_INFO_4 struct {
	PPrinterName *uint16
	PServerName  *uint16
	Attributes   uint
}

var (
	// Library
	libwinspool uintptr

	// Functions
	deviceCapabilities uintptr
	documentProperties uintptr
	enumPrinters       uintptr
	getDefaultPrinter  uintptr
)

func init() {
	// Library
	libwinspool = MustLoadLibrary("winspool.drv")

	// Functions
	deviceCapabilities = MustGetProcAddress(libwinspool, "DeviceCapabilitiesW")
	documentProperties = MustGetProcAddress(libwinspool, "DocumentPropertiesW")
	enumPrinters = MustGetProcAddress(libwinspool, "EnumPrintersW")
	getDefaultPrinter = MustGetProcAddress(libwinspool, "GetDefaultPrinterW")
}

func DeviceCapabilities(pDevice, pPort *uint16, fwCapability uint16, pOutput *uint16, pDevMode *DEVMODE) uint {
	ret, _, _ := syscall.Syscall6(deviceCapabilities, 5,
		uintptr(unsafe.Pointer(pDevice)),
		uintptr(unsafe.Pointer(pPort)),
		uintptr(fwCapability),
		uintptr(unsafe.Pointer(pOutput)),
		uintptr(unsafe.Pointer(pDevMode)),
		0)

	return uint(ret)
}

func DocumentProperties(hWnd HWND, hPrinter HANDLE, pDeviceName *uint16, pDevModeOutput, pDevModeInput *DEVMODE, fMode uint) int {
	ret, _, _ := syscall.Syscall6(documentProperties, 6,
		uintptr(hWnd),
		uintptr(hPrinter),
		uintptr(unsafe.Pointer(pDeviceName)),
		uintptr(unsafe.Pointer(pDevModeOutput)),
		uintptr(unsafe.Pointer(pDevModeInput)),
		uintptr(fMode))

	return int(ret)
}

func EnumPrinters(Flags uint, Name *uint16, Level uint, pPrinterEnum *byte, cbBuf uint, pcbNeeded, pcReturned *uint) bool {
	ret, _, _ := syscall.Syscall9(enumPrinters, 7,
		uintptr(Flags),
		uintptr(unsafe.Pointer(Name)),
		uintptr(Level),
		uintptr(unsafe.Pointer(pPrinterEnum)),
		uintptr(cbBuf),
		uintptr(unsafe.Pointer(pcbNeeded)),
		uintptr(unsafe.Pointer(pcReturned)),
		0,
		0)

	return ret != 0
}

func GetDefaultPrinter(pszBuffer *uint16, pcchBuffer *uint) bool {
	ret, _, _ := syscall.Syscall(getDefaultPrinter, 2,
		uintptr(unsafe.Pointer(pszBuffer)),
		uintptr(unsafe.Pointer(pcchBuffer)),
		0)

	return ret != 0
}
