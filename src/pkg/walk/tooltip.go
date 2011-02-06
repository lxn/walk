// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

type ToolTip struct {
	WidgetBase
}

func NewToolTip(parent Container) (*ToolTip, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		WS_EX_TOPMOST, syscall.StringToUTF16Ptr("tooltips_class32"), nil,
		TTS_ALWAYSTIP|TTS_BALLOON|WS_POPUP,
		CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT,
		parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	tt := &ToolTip{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}}
	tt.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = tt

	parent.Children().Add(tt)

	SetWindowPos(hWnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)

	return tt, nil
}

func (*ToolTip) LayoutFlags() LayoutFlags {
	return 0
}

func (tt *ToolTip) PreferredSize() Size {
	return Size{0, 0}
}

func (tt *ToolTip) Title() string {
	var gt TTGETTITLE

	buf := make([]uint16, 128)

	gt.DwSize = uint(unsafe.Sizeof(gt))
	gt.Cch = uint(len(buf))
	gt.PszTitle = &buf[0]

	SendMessage(tt.hWnd, TTM_GETTITLE, 0, uintptr(unsafe.Pointer(&gt)))

	return syscall.UTF16ToString(buf)
}

func (tt *ToolTip) SetTitle(value string) os.Error {
	if FALSE == SendMessage(tt.hWnd, TTM_SETTITLE, uintptr(TTI_INFO), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value)))) {
		return newError("TTM_SETTITLE failed")
	}

	return nil
}

func (tt *ToolTip) AddWidget(widget Widget, text string) os.Error {
	var ti TOOLINFO

	ti.CbSize = uint(unsafe.Sizeof(ti))
	parent := widget.Parent()
	if parent != nil {
		ti.Hwnd = parent.BaseWidget().hWnd
	}
	ti.UFlags = TTF_IDISHWND | TTF_SUBCLASS
	ti.UId = uintptr(widget.BaseWidget().hWnd)
	ti.LpszText = syscall.StringToUTF16Ptr(text)

	if FALSE == SendMessage(tt.hWnd, TTM_ADDTOOL, 0, uintptr(unsafe.Pointer(&ti))) {
		return newError("TTM_ADDTOOL failed")
	}

	return nil
}

func (tt *ToolTip) RemoveWidget(widget Widget) os.Error {
	panic("not implemented")
}
