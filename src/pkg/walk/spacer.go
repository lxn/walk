// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

const spacerWindowClass = `\o/ Walk_Spacer_Class \o/`

var spacerWndProcPtr uintptr

func spacerWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	s, ok := widgetsByHWnd[hwnd].(*Spacer)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return s.wndProc(hwnd, msg, wParam, lParam, 0)
}

type Spacer struct {
	WidgetBase
	preferredSize   Size
	layoutFlagsMask LayoutFlags
}

func newSpacer(parent Container, layoutFlagsMask LayoutFlags, prefSize Size) (*Spacer, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(spacerWindowClass, spacerWndProc, &spacerWndProcPtr)

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr(spacerWindowClass), nil,
		WS_CHILD,
		0, 0, 0, 0, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	s := &Spacer{
		WidgetBase: WidgetBase{
			hWnd:        hWnd,
			parent:      parent,
			layoutFlags: layoutFlagsMask,
		},
		layoutFlagsMask: layoutFlagsMask,
		preferredSize:   prefSize,
	}

	succeeded := false
	defer func() {
		if !succeeded {
			s.Dispose()
		}
	}()

	s.SetFont(defaultFont)

	if err := parent.Children().Add(s); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = s

	succeeded = true

	return s, nil
}

func NewHSpacer(parent Container) (*Spacer, os.Error) {
	return newSpacer(parent, HGrow|HShrink|VShrink, Size{})
}

func NewHSpacerFixed(parent Container, width int) (*Spacer, os.Error) {
	return newSpacer(parent, 0, Size{width, 0})
}

func NewVSpacer(parent Container) (*Spacer, os.Error) {
	return newSpacer(parent, HShrink|VGrow|VShrink, Size{})
}

func NewVSpacerFixed(parent Container, height int) (*Spacer, os.Error) {
	return newSpacer(parent, HShrink, Size{0, height})
}

func (s *Spacer) LayoutFlagsMask() LayoutFlags {
	return s.layoutFlagsMask
}

func (s *Spacer) PreferredSize() Size {
	return s.preferredSize
}
