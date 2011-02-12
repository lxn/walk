// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"log"
	"math"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

const numberEditWindowClass = `\o/ Walk_NumberEdit_Class \o/`

var numberEditWndProcPtr uintptr

func numberEditWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	ne, ok := widgetsByHWnd[hwnd]
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return ne.wndProc(hwnd, msg, wParam, lParam, 0)
}

type NumberEdit struct {
	WidgetBase
	edit                  *LineEdit
	hWndUpDown            HWND
	decimals              int
	increment             float64
	minValue              float64
	maxValue              float64
	valueChangedPublisher EventPublisher
}

func NewNumberEdit(parent Container) (*NumberEdit, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(numberEditWindowClass, numberEditWndProc, &numberEditWndProcPtr)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(numberEditWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 0, 0, parent.BaseWidget().hWnd, 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	ne := &NumberEdit{WidgetBase: WidgetBase{hWnd: hWnd, parent: parent}}

	var succeeded bool
	defer func() {
		if !succeeded {
			ne.Dispose()
		}
	}()

	var err os.Error
	ne.edit, err = newLineEdit(hWnd)
	if err != nil {
		return nil, err
	}
	if err = ne.edit.setAndClearStyleBits(ES_RIGHT, ES_LEFT|ES_CENTER); err != nil {
		return nil, err
	}

	ne.hWndUpDown = CreateWindowEx(
		0, syscall.StringToUTF16Ptr("msctls_updown32"), nil,
		UDS_ALIGNRIGHT|UDS_ARROWKEYS|UDS_HOTTRACK|WS_CHILD|WS_VISIBLE,
		0, 0, 16, 20, hWnd, 0, 0, nil)
	if ne.hWndUpDown == 0 {
		return nil, lastError("CreateWindowEx")
	}

	SendMessage(ne.hWndUpDown, UDM_SETBUDDY, uintptr(ne.edit.hWnd), 0)

	if err = parent.Children().Add(ne); err != nil {
		return nil, err
	}

	if err = ne.SetValue(0); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = ne

	succeeded = true

	return ne, nil
}

func (ne *NumberEdit) Enabled() bool {
	return ne.edit.Enabled()
}

func (ne *NumberEdit) SetEnabled(value bool) {
	ne.edit.SetEnabled(value)
}

func (ne *NumberEdit) Font() *Font {
	return ne.edit.Font()
}

func (ne *NumberEdit) SetFont(value *Font) {
	ne.edit.SetFont(value)
}

func (*NumberEdit) LayoutFlagsMask() LayoutFlags {
	return HShrink | HGrow
}

func (ne *NumberEdit) PreferredSize() Size {
	return ne.dialogBaseUnitsToPixels(Size{50, 14})
}

func (ne *NumberEdit) Decimals() int {
	return ne.decimals
}

func (ne *NumberEdit) SetDecimals(value int) os.Error {
	if value < 0 {
		return newError("invalid value")
	}

	ne.decimals = value

	return nil
}

func (ne *NumberEdit) Increment() float64 {
	return ne.increment
}

func (ne *NumberEdit) SetIncrement(value float64) os.Error {
	ne.increment = value

	return nil
}

func (ne *NumberEdit) MinValue() float64 {
	return ne.minValue
}

func (ne *NumberEdit) MaxValue() float64 {
	return ne.maxValue
}

func (ne *NumberEdit) SetRange(min, max float64) os.Error {
	if min > max {
		return newError("invalid range")
	}

	ne.minValue = min
	ne.maxValue = max

	return nil
}

func (ne *NumberEdit) Value() float64 {
	val, _ := strconv.Atof64(ne.edit.Text())

	return val
}

func (ne *NumberEdit) SetValue(value float64) os.Error {
	text := strconv.Ftoa64(value, 'f', ne.decimals)
	desiredPrecVal, _ := strconv.Atof64(text)
	if math.Fabs(desiredPrecVal-ne.Value()) < math.SmallestNonzeroFloat64 {
		return nil
	}

	if err := ne.edit.SetText(text); err != nil {
		return err
	}

	return nil
}

func (ne *NumberEdit) ValueChanged() *Event {
	return ne.valueChangedPublisher.Event()
}

func (ne *NumberEdit) SetFocus() os.Error {
	if SetFocus(ne.edit.hWnd) == 0 {
		return lastError("SetFocus")
	}

	return nil
}

func (ne *NumberEdit) TextSelection() (start, end int) {
	return ne.edit.TextSelection()
}

func (ne *NumberEdit) SetTextSelection(start, end int) {
	ne.edit.SetTextSelection(start, end)
}

func (ne *NumberEdit) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint(wParam)) {
		case EN_CHANGE:
			ne.valueChangedPublisher.Publish()
		}

	case WM_NOTIFY:
		switch ((*NMHDR)(unsafe.Pointer(lParam))).Code {
		case UDN_DELTAPOS:
			nmud := (*NMUPDOWN)(unsafe.Pointer(lParam))
			val := ne.Value()
			val -= float64(nmud.IDelta) * ne.increment
			if err := ne.SetValue(val); err != nil {
				log.Println(err)
			}
		}

	case WM_SIZE, WM_SIZING:
		cb := ne.ClientBounds()
		if err := ne.edit.SetBounds(cb); err != nil {
			log.Println(err)
			break
		}
		SendMessage(ne.hWndUpDown, UDM_SETBUDDY, uintptr(ne.edit.hWnd), 0)
	}

	return ne.WidgetBase.wndProc(hwnd, msg, wParam, lParam, origWndProcPtr)
}
