// Copyright 2010 The Walk Authors. All rights reserved.
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

func init() {
	var err error
	if globalToolTip, err = NewToolTip(); err != nil {
		panic(err)
	}
}

type ToolTip struct {
	WidgetBase
}

var globalToolTip *ToolTip

func NewToolTip() (*ToolTip, error) {
	tt := &ToolTip{}

	if err := InitWindow(
		tt,
		nil,
		"tooltips_class32",
		WS_POPUP|TTS_ALWAYSTIP,
		WS_EX_TOPMOST); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			tt.Dispose()
		}
	}()

	SetWindowPos(tt.hWnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)

	succeeded = true

	return tt, nil
}

func (*ToolTip) LayoutFlags() LayoutFlags {
	return 0
}

func (tt *ToolTip) SizeHint() Size {
	return Size{0, 0}
}

func (tt *ToolTip) Title() string {
	var gt TTGETTITLE

	buf := make([]uint16, 100)

	gt.DwSize = uint32(unsafe.Sizeof(gt))
	gt.Cch = uint32(len(buf))
	gt.PszTitle = &buf[0]

	tt.SendMessage(TTM_GETTITLE, 0, uintptr(unsafe.Pointer(&gt)))

	return syscall.UTF16ToString(buf)
}

func (tt *ToolTip) SetTitle(title string) error {
	return tt.setTitle(title, TTI_NONE)
}

func (tt *ToolTip) SetInfoTitle(title string) error {
	return tt.setTitle(title, TTI_INFO)
}

func (tt *ToolTip) SetWarningTitle(title string) error {
	return tt.setTitle(title, TTI_WARNING)
}

func (tt *ToolTip) SetErrorTitle(title string) error {
	return tt.setTitle(title, TTI_ERROR)
}

func (tt *ToolTip) setTitle(title string, icon uintptr) error {
	if len(title) > 99 {
		title = title[:99]
	}

	if FALSE == tt.SendMessage(TTM_SETTITLE, icon, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title)))) {
		return newError("TTM_SETTITLE failed")
	}

	return nil
}

func (tt *ToolTip) AddTool(tool Widget) error {
	hwnd := tool.Handle()

	var ti TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hwnd = hwnd
	ti.UFlags = TTF_IDISHWND | TTF_SUBCLASS
	ti.UId = uintptr(hwnd)

	if FALSE == tt.SendMessage(TTM_ADDTOOL, 0, uintptr(unsafe.Pointer(&ti))) {
		return newError("TTM_ADDTOOL failed")
	}

	return nil
}

func (tt *ToolTip) RemoveTool(tool Widget) error {
	hwnd := tool.Handle()

	var ti TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hwnd = hwnd
	ti.UId = uintptr(hwnd)

	tt.SendMessage(TTM_DELTOOL, 0, uintptr(unsafe.Pointer(&ti)))

	return nil
}

func (tt *ToolTip) Text(tool Widget) string {
	ti := tt.toolInfo(tool)
	if ti == nil {
		return ""
	}

	return UTF16PtrToString(ti.LpszText)
}

func (tt *ToolTip) SetText(tool Widget, text string) error {
	ti := tt.toolInfo(tool)
	if ti == nil {
		return newError("unknown tool")
	}

	if len(text) > 79 {
		text = text[:79]
	}

	ti.LpszText = syscall.StringToUTF16Ptr(text)

	tt.SendMessage(TTM_SETTOOLINFO, 0, uintptr(unsafe.Pointer(ti)))

	return nil
}

func (tt *ToolTip) toolInfo(tool Widget) *TOOLINFO {
	var ti TOOLINFO
	var buf [80]uint16

	hwnd := tool.Handle()

	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hwnd = hwnd
	ti.UId = uintptr(hwnd)
	ti.LpszText = &buf[0]

	if FALSE == tt.SendMessage(TTM_GETTOOLINFO, 0, uintptr(unsafe.Pointer(&ti))) {
		return nil
	}

	return &ti
}
