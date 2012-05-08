// Copyright 2012 The Walk Authors. All rights reserved.
// Use of lb source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

var listBoxOrigWndProcPtr uintptr
var _ subclassedWidget = &ListBox{}

type ListBox struct {
	WidgetBase
	maxItemTextWidth             int
	CurrentIndexChangedPublisher EventPublisher
	dblClickedPublisher          EventPublisher
}

func NewListBox(parent Container) (*ListBox, error) {
	lb := &ListBox{}
	err := initChildWidget(
		lb,
		parent,
		"LISTBOX",
		WS_TABSTOP|WS_VISIBLE|LBS_STANDARD,
		0)
	if err != nil {
		return nil, err
	}
	return lb, nil
}

func (*ListBox) origWndProcPtr() uintptr {
	return listBoxOrigWndProcPtr
}

func (*ListBox) setOrigWndProcPtr(ptr uintptr) {
	listBoxOrigWndProcPtr = ptr
}

func (*ListBox) LayoutFlags() LayoutFlags {
	return GrowableHorz | GrowableVert
}

func (lb *ListBox) AddString(item string) {
	SendMessage(lb.hWnd, LB_ADDSTRING, 0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item))))
}

//If lb parameter is -1, the string is added to the end of the list.
func (lb *ListBox) InsertString(index int, item string) error {
	if index < -1 {
		return newError("Invalid index")
	}

	ret := int(SendMessage(lb.hWnd, LB_INSERTSTRING, uintptr(index), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item)))))
	if ret == LB_ERRSPACE || ret == LB_ERR {
		return newError("Fail to insert string")
	}
	return nil
}

func (lb *ListBox) DeleteString(index uint) error {
	ret := int(SendMessage(lb.hWnd, LB_DELETESTRING, uintptr(index), 0))
	if ret == LB_ERR {
		return newError("Fail to delete string")
	}
	return nil
}

func (lb *ListBox) GetString(index uint) string {
	len := int(SendMessage(lb.hWnd, LB_GETTEXTLEN, uintptr(index), 0))
	if len == LB_ERR {
		return ""
	}

	buf := make([]uint16, len+1)
	_ = SendMessage(lb.hWnd, LB_GETTEXT, uintptr(index), uintptr(unsafe.Pointer(&buf[0])))

	if len == LB_ERR {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func (lb *ListBox) ResetContent() {
	SendMessage(lb.hWnd, LB_RESETCONTENT, 0, 0)
}

//The return value is the number of items in the list box,
//or LB_ERR (-1) if an error occurs.
func (lb *ListBox) GetCount() (uint, error) {
	retPtr := SendMessage(lb.hWnd, LB_GETCOUNT, 0, 0)
	ret := int(retPtr)
	if ret == LB_ERR {
		return 0, newError("Fail to get count")
	}
	return uint(ret), nil
}

func (lb *ListBox) calculateMaxItemTextWidth() int {
	hdc := GetDC(lb.hWnd)
	if hdc == 0 {
		newError("GetDC failed")
		return -1
	}
	defer ReleaseDC(lb.hWnd, hdc)

	hFontOld := SelectObject(hdc, HGDIOBJ(lb.Font().handleForDPI(0)))
	defer SelectObject(hdc, hFontOld)

	var maxWidth int

	count, _ := lb.GetCount()
	var i uint
	for i = 0; i < count; i++ {
		item := lb.GetString(i)
		var s SIZE
		str := syscall.StringToUTF16(item)

		if !GetTextExtentPoint32(hdc, &str[0], int32(len(str)-1), &s) {
			newError("GetTextExtentPoint32 failed")
			return -1
		}

		maxWidth = maxi(maxWidth, int(s.CX))
	}

	return maxWidth
}

func (lb *ListBox) SizeHint() Size {

	defaultSize := lb.dialogBaseUnitsToPixels(Size{50, 12})

	if lb.maxItemTextWidth <= 0 {
		lb.maxItemTextWidth = lb.calculateMaxItemTextWidth()
	}

	// FIXME: Use GetThemePartSize instead of guessing
	w := maxi(defaultSize.Width, lb.maxItemTextWidth+24)
	h := defaultSize.Height + 1

	return Size{w, h}

}

func (lb *ListBox) CurrentIndex() int {
	return int(SendMessage(lb.hWnd, LB_GETCURSEL, 0, 0))
}

func (lb *ListBox) SetCurrentIndex(value int) error {
	ret := int(SendMessage(lb.hWnd, LB_SETCURSEL, uintptr(value), 0))
	if ret == LB_ERR {
		return newError("Invalid index or ensure lb is single-selection listbox")
	}
	return nil
}

func (lb *ListBox) CurrentString() string {
	index := lb.CurrentIndex()
	length := int(SendMessage(lb.hWnd, LB_GETTEXTLEN, uintptr(index), 0)) + 1
	buffer := make([]uint16, length+1)
	SendMessage(lb.hWnd, LB_GETTEXT, uintptr(index), uintptr(unsafe.Pointer(&buffer[0])))
	return syscall.UTF16ToString(buffer)
}

func (lb *ListBox) CurrentIndexChanged() *Event {
	return lb.CurrentIndexChangedPublisher.Event()
}

func (lb *ListBox) DblClicked() *Event {
	return lb.dblClickedPublisher.Event()
}

func (lb *ListBox) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case LBN_SELCHANGE:
			lb.CurrentIndexChangedPublisher.Publish()
		case LBN_DBLCLK:
			lb.dblClickedPublisher.Publish()
		}
	}

	return lb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
