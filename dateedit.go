// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"time"
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

type DateEdit struct {
	WidgetBase
	dateChangedPublisher EventPublisher
}

func newDateEdit(parent Container, style uint32) (*DateEdit, error) {
	de := new(DateEdit)

	if err := InitWidget(
		de,
		parent,
		"SysDateTimePick32",
		WS_TABSTOP|WS_VISIBLE|DTS_SHORTDATEFORMAT|style,
		0); err != nil {
		return nil, err
	}

	de.MustRegisterProperty("Date", NewProperty(
		func() interface{} {
			return de.Date()
		},
		func(v interface{}) error {
			return de.SetDate(v.(time.Time))
		},
		de.dateChangedPublisher.Event()))

	return de, nil
}

func NewDateEdit(parent Container) (*DateEdit, error) {
	return newDateEdit(parent, 0)
}

func NewDateEditWithNoneOption(parent Container) (*DateEdit, error) {
	return newDateEdit(parent, DTS_SHOWNONE)
}

func (*DateEdit) LayoutFlags() LayoutFlags {
	return GrowableHorz
}

func (de *DateEdit) MinSizeHint() Size {
	return de.dialogBaseUnitsToPixels(Size{64, 12})
}

func (de *DateEdit) SizeHint() Size {
	return de.MinSizeHint()
}

func (de *DateEdit) systemTimeToTime(st *SYSTEMTIME) time.Time {
	if st == nil || !de.hasStyleBits(DTS_SHOWNONE) && st.WYear == 1601 && st.WMonth == 1 && st.WDay == 1 {
		return time.Time{}
	}

	return time.Date(int(st.WYear), time.Month(st.WMonth), int(st.WDay), 0, 0, 0, 0, time.Local)
}

func (de *DateEdit) timeToSystemTime(t time.Time) *SYSTEMTIME {
	if t.IsZero() {
		if de.hasStyleBits(DTS_SHOWNONE) {
			return nil
		} else {
			return &SYSTEMTIME{
				WYear:  uint16(1601),
				WMonth: uint16(1),
				WDay:   uint16(1),
			}
		}
	}

	return &SYSTEMTIME{
		WYear:  uint16(t.Year()),
		WMonth: uint16(t.Month()),
		WDay:   uint16(t.Day()),
	}
}

func (de *DateEdit) systemTime() (*SYSTEMTIME, error) {
	var st SYSTEMTIME

	switch de.SendMessage(DTM_GETSYSTEMTIME, 0, uintptr(unsafe.Pointer(&st))) {
	case GDT_VALID:
		return &st, nil

	case GDT_NONE:
		return nil, nil
	}

	return nil, newError("SendMessage(DTM_GETSYSTEMTIME)")
}

func (de *DateEdit) setSystemTime(st *SYSTEMTIME) error {
	var wParam uintptr

	if st != nil {
		wParam = GDT_VALID
	} else {
		wParam = GDT_NONE
	}

	if 0 == de.SendMessage(DTM_SETSYSTEMTIME, wParam, uintptr(unsafe.Pointer(st))) {
		return newError("SendMessage(DTM_SETSYSTEMTIME)")
	}

	de.dateChangedPublisher.Publish()

	return nil
}

func (de *DateEdit) Range() (min, max time.Time) {
	var st [2]SYSTEMTIME

	ret := de.SendMessage(DTM_GETRANGE, 0, uintptr(unsafe.Pointer(&st[0])))

	if ret&GDTR_MIN > 0 {
		min = de.systemTimeToTime(&st[0])
	}

	if ret&GDTR_MAX > 0 {
		max = de.systemTimeToTime(&st[1])
	}

	return
}

func (de *DateEdit) SetRange(min, max time.Time) error {
	if !min.IsZero() && !max.IsZero() {
		if min.Year() > max.Year() ||
			min.Year() == max.Year() && min.Month() > max.Month() ||
			min.Year() == max.Year() && min.Month() == max.Month() && min.Day() > max.Day() {
			return newError("invalid range")
		}
	}

	var st [2]SYSTEMTIME
	var wParam uintptr

	if !min.IsZero() {
		wParam |= GDTR_MIN
		st[0] = *de.timeToSystemTime(min)
	}

	if !max.IsZero() {
		wParam |= GDTR_MAX
		st[1] = *de.timeToSystemTime(max)
	}

	if 0 == de.SendMessage(DTM_SETRANGE, wParam, uintptr(unsafe.Pointer(&st[0]))) {
		return newError("SendMessage(DTM_SETRANGE)")
	}

	return nil
}

func (de *DateEdit) Date() time.Time {
	st, err := de.systemTime()
	if err != nil {
		return time.Time{}
	}

	if st == nil {
		return time.Time{}
	}

	return time.Unix(de.systemTimeToTime(st).Unix(), 0)
}

func (de *DateEdit) SetDate(date time.Time) error {
	return de.setSystemTime(de.timeToSystemTime(date))
}

func (de *DateEdit) DateChanged() *Event {
	return de.dateChangedPublisher.Event()
}

func (de *DateEdit) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_NOTIFY:
		switch uint32(((*NMHDR)(unsafe.Pointer(lParam))).Code) {
		case DTN_DATETIMECHANGE:
			de.dateChangedPublisher.Publish()
		}
	}

	return de.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
