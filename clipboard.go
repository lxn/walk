// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

const clipboardWindowClass = `\o/ Walk_Clipboard_Class \o/`

func init() {
	MustRegisterWindowClass(clipboardWindowClass)

	hwnd := CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr(clipboardWindowClass),
		nil,
		0,
		0,
		0,
		0,
		0,
		HWND_MESSAGE,
		0,
		0,
		nil)

	if hwnd == 0 {
		panic("failed to create clipboard window")
	}

	clipboard.hwnd = hwnd
}

var clipboard ClipboardService

// Clipboard returns an object that provides access to the system clipboard.
func Clipboard() *ClipboardService {
	return &clipboard
}

// ClipboardService provides access to the system clipboard.
type ClipboardService struct {
	hwnd HWND
}

// Clear clears the contents of the clipboard.
func (c *ClipboardService) Clear() error {
	return c.withOpenClipboard(func() error {
		if !EmptyClipboard() {
			return lastError("EmptyClipboard")
		}

		return nil
	})
}

// ContainsText returns whether the clipboard currently contains text data.
func (c *ClipboardService) ContainsText() (available bool, err error) {
	err = c.withOpenClipboard(func() error {
		available = IsClipboardFormatAvailable(CF_UNICODETEXT)

		return nil
	})

	return
}

// Text returns the current text data of the clipboard.
func (c *ClipboardService) Text() (text string, err error) {
	err = c.withOpenClipboard(func() error {
		hMem := HGLOBAL(GetClipboardData(CF_UNICODETEXT))
		if hMem == 0 {
			return lastError("GetClipboardData")
		}

		p := GlobalLock(hMem)
		if p == nil {
			return lastError("GlobalLock()")
		}
		defer GlobalUnlock(hMem)

		text = UTF16PtrToString((*uint16)(p))

		return nil
	})

	return
}

// SetText sets the current text data of the clipboard.
func (c *ClipboardService) SetText(s string) error {
	return c.withOpenClipboard(func() error {
		utf16, err := syscall.UTF16FromString(s)
		if err != nil {
			return err
		}

		hMem := GlobalAlloc(GMEM_MOVEABLE, uintptr(len(utf16)*2))
		if hMem == 0 {
			return lastError("GlobalAlloc")
		}

		p := GlobalLock(hMem)
		if p == nil {
			return lastError("GlobalLock()")
		}

		MoveMemory(p, unsafe.Pointer(&utf16[0]), uintptr(len(utf16)*2))

		GlobalUnlock(hMem)

		if 0 == SetClipboardData(CF_UNICODETEXT, HANDLE(hMem)) {
			// We need to free hMem.
			defer GlobalFree(hMem)

			return lastError("SetClipboardData")
		}

		// The system now owns the memory referred to by hMem.

		return nil
	})
}

func (c *ClipboardService) withOpenClipboard(f func() error) error {
	if !OpenClipboard(c.hwnd) {
		return lastError("OpenClipboard")
	}
	defer CloseClipboard()

	return f()
}
