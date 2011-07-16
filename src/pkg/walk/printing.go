// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"unsafe"
)

import . "walk/winapi"

func defaultPrinterNamePtr() *uint16 {
	var bufLen uint32

	GetDefaultPrinter(nil, &bufLen)

	buf := make([]uint16, bufLen)

	if !GetDefaultPrinter(&buf[0], &bufLen) {
		panic("failed to retrieve default printer name")
	}

	return &buf[0]
}

func defaultPrinterName() string {
	return UTF16PtrToString(defaultPrinterNamePtr())
}

// DefaultPrinterName returns the name of the default printer.
func DefaultPrinterName() (name string, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return defaultPrinterName(), nil
}

// PrinterNames returns the names of the installed printers.
func PrinterNames() ([]string, os.Error) {
	var bufLen, count uint32

	EnumPrinters(PRINTER_ENUM_LOCAL, nil, 4, nil, 0, &bufLen, &count)

	if bufLen == 0 {
		return make([]string, 0), nil
	}

	buf := make([]byte, int(bufLen))

	if !EnumPrinters(PRINTER_ENUM_LOCAL, nil, 4, &buf[0], bufLen, &bufLen, &count) {
		return nil, newError("EnumPrinters failed")
	}

	printers := (*[1 << 24]PRINTER_INFO_4)(unsafe.Pointer(&buf[0]))
	printerNames := make([]string, count)

	for i := 0; i < int(count); i++ {
		printerNames[i] = UTF16PtrToString(printers[i].PPrinterName)
	}

	return printerNames, nil
}
