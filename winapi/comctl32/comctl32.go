// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comctl32

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/gdi32"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

// Button control messages
const (
	BCM_FIRST            = 0x1600
	BCM_GETIDEALSIZE     = BCM_FIRST + 0x0001
	BCM_SETIMAGELIST     = BCM_FIRST + 0x0002
	BCM_GETIMAGELIST     = BCM_FIRST + 0x0003
	BCM_SETTEXTMARGIN    = BCM_FIRST + 0x0004
	BCM_GETTEXTMARGIN    = BCM_FIRST + 0x0005
	BCM_SETDROPDOWNSTATE = BCM_FIRST + 0x0006
	BCM_SETSPLITINFO     = BCM_FIRST + 0x0007
	BCM_GETSPLITINFO     = BCM_FIRST + 0x0008
	BCM_SETNOTE          = BCM_FIRST + 0x0009
	BCM_GETNOTE          = BCM_FIRST + 0x000A
	BCM_GETNOTELENGTH    = BCM_FIRST + 0x000B
)

const (
	CCM_FIRST            = 0x2000
	CCM_LAST             = CCM_FIRST + 0x200
	CCM_SETBKCOLOR       = 8193
	CCM_SETCOLORSCHEME   = 8194
	CCM_GETCOLORSCHEME   = 8195
	CCM_GETDROPTARGET    = 8196
	CCM_SETUNICODEFORMAT = 8197
	CCM_GETUNICODEFORMAT = 8198
	CCM_SETVERSION       = 0x2007
	CCM_GETVERSION       = 0x2008
	CCM_SETNOTIFYWINDOW  = 0x2009
	CCM_SETWINDOWTHEME   = 0x200b
	CCM_DPISCALE         = 0x200c
)

// Common controls styles
const (
	CCS_TOP           = 1
	CCS_NOMOVEY       = 2
	CCS_BOTTOM        = 3
	CCS_NORESIZE      = 4
	CCS_NOPARENTALIGN = 8
	CCS_ADJUSTABLE    = 32
	CCS_NODIVIDER     = 64
	CCS_VERT          = 128
	CCS_LEFT          = 129
	CCS_NOMOVEX       = 130
	CCS_RIGHT         = 131
)

// InitCommonControlsEx flags
const (
	ICC_LISTVIEW_CLASSES   = 1
	ICC_TREEVIEW_CLASSES   = 2
	ICC_BAR_CLASSES        = 4
	ICC_TAB_CLASSES        = 8
	ICC_UPDOWN_CLASS       = 16
	ICC_PROGRESS_CLASS     = 32
	ICC_HOTKEY_CLASS       = 64
	ICC_ANIMATE_CLASS      = 128
	ICC_WIN95_CLASSES      = 255
	ICC_DATE_CLASSES       = 256
	ICC_USEREX_CLASSES     = 512
	ICC_COOL_CLASSES       = 1024
	ICC_INTERNET_CLASSES   = 2048
	ICC_PAGESCROLLER_CLASS = 4096
	ICC_NATIVEFNTCTL_CLASS = 8192
	INFOTIPSIZE            = 1024
	ICC_STANDARD_CLASSES   = 0x00004000
	ICC_LINK_CLASS         = 0x00008000
)

// WM_NOTITY messages
const (
	NM_FIRST           = 4294967295
	NM_OUTOFMEMORY     = NM_FIRST - 0
	NM_CLICK           = NM_FIRST - 1
	NM_DBLCLK          = NM_FIRST - 2
	NM_RETURN          = NM_FIRST - 3
	NM_RCLICK          = NM_FIRST - 4
	NM_RDBLCLK         = NM_FIRST - 5
	NM_SETFOCUS        = NM_FIRST - 6
	NM_KILLFOCUS       = NM_FIRST - 7
	NM_CUSTOMDRAW      = NM_FIRST - 11
	NM_HOVER           = NM_FIRST - 12
	NM_NCHITTEST       = NM_FIRST - 13
	NM_KEYDOWN         = NM_FIRST - 14
	NM_RELEASEDCAPTURE = NM_FIRST - 15
	NM_SETCURSOR       = NM_FIRST - 16
	NM_CHAR            = NM_FIRST - 17
	NM_TOOLTIPSCREATED = NM_FIRST - 18
	NM_LAST            = NM_FIRST - 98
)

// ProgressBar messages
const (
	PBM_SETPOS      = WM_USER + 2
	PBM_DELTAPOS    = WM_USER + 3
	PBM_SETSTEP     = WM_USER + 4
	PBM_STEPIT      = WM_USER + 5
	PBM_SETRANGE32  = 1030
	PBM_GETRANGE    = 1031
	PBM_GETPOS      = 1032
	PBM_SETBARCOLOR = 1033
	PBM_SETBKCOLOR  = CCM_SETBKCOLOR
	PBS_SMOOTH      = 1
	PBS_VERTICAL    = 4
)

type HIMAGELIST HANDLE

type INITCOMMONCONTROLSEX struct {
	DwSize, DwICC uint
}

var (
	// Library
	lib uint32

	// Functions
	imageList_Destroy    uint32
	imageList_LoadImage  uint32
	initCommonControlsEx uint32
)

func init() {
	// Library
	lib = MustLoadLibrary("comctl32.dll")

	// Functions
	imageList_Destroy = MustGetProcAddress(lib, "ImageList_Destroy")
	imageList_LoadImage = MustGetProcAddress(lib, "ImageList_LoadImageW")
	initCommonControlsEx = MustGetProcAddress(lib, "InitCommonControlsEx")

	// Initialize the common controls we support
	var initCtrls INITCOMMONCONTROLSEX
	initCtrls.DwSize = uint(unsafe.Sizeof(initCtrls))
	initCtrls.DwICC = ICC_LISTVIEW_CLASSES | ICC_PROGRESS_CLASS

	InitCommonControlsEx(&initCtrls)
}

func ImageList_Destroy(hIml HIMAGELIST) bool {
	ret, _, _ := syscall.Syscall(uintptr(imageList_Destroy),
		uintptr(hIml),
		0,
		0)

	return ret != 0
}

func ImageList_LoadImage(hi HINSTANCE, lpbmp *uint16, cx, cGrow int, crMask COLORREF, uType, uFlags uint) HIMAGELIST {
	ret, _, _ := syscall.Syscall9(uintptr(imageList_LoadImage),
		uintptr(hi),
		uintptr(unsafe.Pointer(lpbmp)),
		uintptr(cx),
		uintptr(cGrow),
		uintptr(crMask),
		uintptr(uType),
		uintptr(uFlags),
		0,
		0)

	return HIMAGELIST(ret)
}

func InitCommonControlsEx(lpInitCtrls *INITCOMMONCONTROLSEX) bool {
	ret, _, _ := syscall.Syscall(uintptr(initCommonControlsEx),
		uintptr(unsafe.Pointer(lpInitCtrls)),
		0,
		0)

	return ret != 0
}
