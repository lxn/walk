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
	providedModel                interface{}
	dataMember                   string
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
	err := InitChildWidget(
		lb,
		parent,
		"LISTBOX",
		WS_TABSTOP|WS_VISIBLE|LBS_NOINTEGRALHEIGHT|LBS_STANDARD,
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
	ret := int(lb.SendMessage(LB_INSERTSTRING, uintptr(index), lp))
	if ret == LB_ERRSPACE || ret == LB_ERR {
		return newError("SendMessage(LB_INSERTSTRING)")
	}
	return nil
}

// reread all the items from list model
func (lb *ListBox) resetItems() error {
	lb.SetSuspended(true)
	defer lb.SetSuspended(false)

	lb.SendMessage(LB_RESETCONTENT, 0, 0)

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
		if CB_ERR == lb.SendMessage(LB_DELETESTRING, uintptr(index), 0) {
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

// Model returns the model of the ListBox.
func (lb *ListBox) Model() interface{} {
	return lb.providedModel
}

// SetModel sets the model of the ListBox.
//
// It is required that mdl either implements walk.ListModel or
// walk.ReflectListModel or be a slice of pointers to struct.
func (lb *ListBox) SetModel(mdl interface{}) error {
	model, ok := mdl.(ListModel)
	if !ok && mdl != nil {
		var err error
		if model, err = newReflectListModel(mdl); err != nil {
			return err
		}

		if badms, ok := model.(bindingAndDisplayMemberSetter); ok {
			badms.setDisplayMember(lb.dataMember)
		}
	}
	lb.providedModel = mdl

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

// DataMember returns the member from the model of the ListBox that is displayed
// in the ListBox.
//
// This is only applicable to walk.ReflectListModel models and simple slices of
// pointers to struct.
func (lb *ListBox) DataMember() string {
	return lb.dataMember
}

// SetDataMember sets the member from the model of the ListBox that is displayed
// in the ListBox.
//
// This is only applicable to walk.ReflectListModel models and simple slices of
// pointers to struct.
//
// For a model consisting of items of type S, the type of the specified member T
// and displayMember "Foo", this can be one of the following:
//
//	A field		Foo T
//	A method	func (s S) Foo() T
//	A method	func (s S) Foo() (T, error)
//
// If displayMember is not a simple member name like "Foo", but a path to a
// member like "A.B.Foo", members "A" and "B" both must be one of the options
// mentioned above, but with T having type pointer to struct.
func (lb *ListBox) SetDataMember(dataMember string) {
	lb.dataMember = dataMember

	if badms, ok := lb.model.(bindingAndDisplayMemberSetter); ok {
		badms.setDisplayMember(dataMember)
	}
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
	return int(lb.SendMessage(LB_GETCURSEL, 0, 0))
}

func (lb *ListBox) SetCurrentIndex(value int) error {
	if value < 0 {
		return nil
	}
	ret := int(lb.SendMessage(LB_SETCURSEL, uintptr(value), 0))
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
			lb.prevCurIndex = lb.CurrentIndex()
			lb.currentIndexChangedPublisher.Publish()

		case LBN_DBLCLK:
			lb.dblClickedPublisher.Publish()
		}
	}

	return lb.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
