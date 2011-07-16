// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winapi

import (
	"syscall"
	"unsafe"
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
	NM_FIRST           = 0
	NM_OUTOFMEMORY     = NM_FIRST - 1
	NM_CLICK           = NM_FIRST - 2
	NM_DBLCLK          = NM_FIRST - 3
	NM_RETURN          = NM_FIRST - 4
	NM_RCLICK          = NM_FIRST - 5
	NM_RDBLCLK         = NM_FIRST - 6
	NM_SETFOCUS        = NM_FIRST - 7
	NM_KILLFOCUS       = NM_FIRST - 8
	NM_CUSTOMDRAW      = NM_FIRST - 12
	NM_HOVER           = NM_FIRST - 13
	NM_NCHITTEST       = NM_FIRST - 14
	NM_KEYDOWN         = NM_FIRST - 15
	NM_RELEASEDCAPTURE = NM_FIRST - 16
	NM_SETCURSOR       = NM_FIRST - 17
	NM_CHAR            = NM_FIRST - 18
	NM_TOOLTIPSCREATED = NM_FIRST - 19
	NM_LAST            = NM_FIRST - 99
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

// ImageList creation flags
const (
	ILC_MASK          = 0x00000001
	ILC_COLOR         = 0x00000000
	ILC_COLORDDB      = 0x000000FE
	ILC_COLOR4        = 0x00000004
	ILC_COLOR8        = 0x00000008
	ILC_COLOR16       = 0x00000010
	ILC_COLOR24       = 0x00000018
	ILC_COLOR32       = 0x00000020
	ILC_PALETTE       = 0x00000800
	ILC_MIRROR        = 0x00002000
	ILC_PERITEMMIRROR = 0x00008000
)

type HIMAGELIST HANDLE

type INITCOMMONCONTROLSEX struct {
	DwSize, DwICC uint32
}

var (
	// Library
	libcomctl32 uintptr

	// Functions
	imageList_Add        uintptr
	imageList_AddMasked  uintptr
	imageList_Create     uintptr
	imageList_Destroy    uintptr
	initCommonControlsEx uintptr
)

func init() {
	// Library
	libcomctl32 = MustLoadLibrary("comctl32.dll")

	// Functions
	imageList_Add = MustGetProcAddress(libcomctl32, "ImageList_Add")
	imageList_AddMasked = MustGetProcAddress(libcomctl32, "ImageList_AddMasked")
	imageList_Create = MustGetProcAddress(libcomctl32, "ImageList_Create")
	imageList_Destroy = MustGetProcAddress(libcomctl32, "ImageList_Destroy")
	initCommonControlsEx = MustGetProcAddress(libcomctl32, "InitCommonControlsEx")

	// Initialize the common controls we support
	var initCtrls INITCOMMONCONTROLSEX
	initCtrls.DwSize = uint32(unsafe.Sizeof(initCtrls))
	initCtrls.DwICC = ICC_LISTVIEW_CLASSES | ICC_PROGRESS_CLASS | ICC_TAB_CLASSES | ICC_TREEVIEW_CLASSES

	InitCommonControlsEx(&initCtrls)
}

func ImageList_Add(himl HIMAGELIST, hbmImage, hbmMask HBITMAP) int32 {
	ret, _, _ := syscall.Syscall(imageList_Add, 3,
		uintptr(himl),
		uintptr(hbmImage),
		uintptr(hbmMask))

	return int32(ret)
}

func ImageList_AddMasked(himl HIMAGELIST, hbmImage HBITMAP, crMask COLORREF) int32 {
	ret, _, _ := syscall.Syscall(imageList_AddMasked, 3,
		uintptr(himl),
		uintptr(hbmImage),
		uintptr(crMask))

	return int32(ret)
}

func ImageList_Create(cx, cy int32, flags uint32, cInitial, cGrow int32) HIMAGELIST {
	ret, _, _ := syscall.Syscall6(imageList_Create, 5,
		uintptr(cx),
		uintptr(cy),
		uintptr(flags),
		uintptr(cInitial),
		uintptr(cGrow),
		0)

	return HIMAGELIST(ret)
}

func ImageList_Destroy(hIml HIMAGELIST) bool {
	ret, _, _ := syscall.Syscall(imageList_Destroy, 1,
		uintptr(hIml),
		0,
		0)

	return ret != 0
}

func InitCommonControlsEx(lpInitCtrls *INITCOMMONCONTROLSEX) bool {
	ret, _, _ := syscall.Syscall(initCommonControlsEx, 1,
		uintptr(unsafe.Pointer(lpInitCtrls)),
		0,
		0)

	return ret != 0
}
