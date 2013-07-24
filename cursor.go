// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"image"
)

import (
	. "github.com/lxn/go-winapi"
)

type Cursor interface {
	Dispose()
	handle() HCURSOR
}

type stockCursor struct {
	hCursor HCURSOR
}

func (sc stockCursor) Dispose() {
	// nop
}

func (sc stockCursor) handle() HCURSOR {
	return sc.hCursor
}

func CursorArrow() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_ARROW))}
}

func CursorIBeam() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_IBEAM))}
}

func CursorWait() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_WAIT))}
}

func CursorCross() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_CROSS))}
}

func CursorUpArrow() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_UPARROW))}
}

func CursorSizeNWSE() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_SIZENWSE))}
}

func CursorSizeNESW() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_SIZENESW))}
}

func CursorSizeWE() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_SIZEWE))}
}

func CursorSizeNS() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_SIZENS))}
}

func CursorSizeAll() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_SIZEALL))}
}

func CursorNo() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_NO))}
}

func CursorHand() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_HAND))}
}

func CursorAppStarting() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_APPSTARTING))}
}

func CursorHelp() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_HELP))}
}

func CursorIcon() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_ICON))}
}

func CursorSize() Cursor {
	return stockCursor{LoadCursor(0, MAKEINTRESOURCE(IDC_SIZE))}
}

type customCursor struct {
	hCursor HCURSOR
}

func NewCursorFromImage(im image.Image, hotspot image.Point) (Cursor, error) {
	i, err := createAlphaCursorOrIconFromImage(im, hotspot, false)
	if err != nil {
		return nil, err
	}
	return customCursor{HCURSOR(i)}, nil
}

func (cc customCursor) Dispose() {
	DestroyIcon(HICON(cc.hCursor))
}

func (cc customCursor) handle() HCURSOR {
	return cc.hCursor
}
