// Copyright 2012 The Walk Authors. All rights reserved.
// Use of lb source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
	"math/big"
	"syscall"
	"time"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

type ListBox struct {
	WidgetBase
	model                        ListModel
	format                       string
	precision                    int
	prevCurIndex                 int
	itemsResetHandlerHandle      int
	itemChangedHandlerHandle     int
	maxItemTextWidth             int
	currentIndexChangedPublisher EventPublisher
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

func (*ListBox) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (lb *ListBox) itemString(index int) string {
	switch val := lb.model.Value(index).(type) {
	case string:
		return val

	case time.Time:
		return val.Format(lb.format)

	case *big.Rat:
		return val.FloatString(lb.precision)

	default:
		return fmt.Sprintf(lb.format, val)
	}

	panic("unreachable")
}

//insert one item from list model
func (lb *ListBox) insertItemAt(index int) error {
	str := lb.itemString(index)
	lp := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str)))
	ret := int(SendMessage(lb.hWnd, LB_INSERTSTRING, uintptr(index), lp))
	if ret == LB_ERRSPACE || ret == LB_ERR {
		return newError("SendMessage(LB_INSERTSTRING)")
	}
	return nil
}

// reread all the items from list model
func (lb *ListBox) resetItems() error {
	lb.SetSuspended(true)
	defer lb.SetSuspended(false)

	SendMessage(lb.hWnd, LB_RESETCONTENT, 0, 0)

	lb.maxItemTextWidth = 0

	lb.SetCurrentIndex(-1)

	if lb.model == nil {
		return nil
	}

	count := lb.model.ItemCount()

	for i := 0; i < count; i++ {
		if err := lb.insertItemAt(i); err != nil {
			return err
		}
	}

	return nil
}

func (lb *ListBox) attachModel() {
	itemsResetHandler := func() {
		lb.resetItems()
	}
	lb.itemsResetHandlerHandle = lb.model.ItemsReset().Attach(itemsResetHandler)

	itemChangedHandler := func(index int) {
		if CB_ERR == SendMessage(lb.hWnd, LB_DELETESTRING, uintptr(index), 0) {
			newError("SendMessage(CB_DELETESTRING)")
		}

		lb.insertItemAt(index)

		lb.SetCurrentIndex(lb.prevCurIndex)
	}
	lb.itemChangedHandlerHandle = lb.model.ItemChanged().Attach(itemChangedHandler)
}

func (lb *ListBox) detachModel() {
	lb.model.ItemsReset().Detach(lb.itemsResetHandlerHandle)
	lb.model.ItemChanged().Detach(lb.itemChangedHandlerHandle)
}

func (lb *ListBox) Model() ListModel {
	return lb.model
}

func (lb *ListBox) SetModel(model ListModel) error {
	if lb.model != nil {
		lb.detachModel()
	}

	lb.model = model

	if model != nil {
		lb.attachModel()

		return lb.resetItems()
	}

	return nil
}

func (lb *ListBox) Format() string {
	return lb.format
}

func (lb *ListBox) SetFormat(value string) {
	lb.format = value
}

func (lb *ListBox) Precision() int {
	return lb.precision
}

func (lb *ListBox) SetPrecision(value int) {
	lb.precision = value
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

	if lb.model == nil {
		return -1
	}
	count := lb.model.ItemCount()
	for i := 0; i < count; i++ {
		item := lb.itemString(i)
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
	if value < 0 {
		return nil
	}
	ret := int(SendMessage(lb.hWnd, LB_SETCURSEL, uintptr(value), 0))
	if ret == LB_ERR {
		return newError("Invalid index or ensure lb is single-selection listbox")
	}

	if value != lb.prevCurIndex {
		lb.prevCurIndex = value
		lb.currentIndexChangedPublisher.Publish()
	}
	return nil
}

func (lb *ListBox) CurrentIndexChanged() *Event {
	return lb.currentIndexChangedPublisher.Event()
}

func (lb *ListBox) DblClicked() *Event {
	return lb.dblClickedPublisher.Event()
}

func (lb *ListBox) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case LBN_SELCHANGE:
			lb.currentIndexChangedPublisher.Publish()

		case LBN_DBLCLK:
			lb.dblClickedPublisher.Publish()
		}
	}

	return lb.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
