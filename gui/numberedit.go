// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"log"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

const numberEditWindowClass = `\o/ Walk_NumberEdit_Class \o/`

var numberEditWndProcCallback *syscall.Callback

func numberEditWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	ne, ok := widgetsByHWnd[msg.HWnd]
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return ne.wndProc(msg, 0)
}

type NumberEdit struct {
	Widget
	edit       *LineEdit
	hWndUpDown HWND
	decimals   int
	increment  float64
	minValue   float64
	maxValue   float64
}

func NewNumberEdit(parent IContainer) (*NumberEdit, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(numberEditWindowClass, numberEditWndProc, &numberEditWndProcCallback)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(numberEditWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	ne := &NumberEdit{Widget: Widget{hWnd: hWnd, parent: parent}}

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
	if err = ne.edit.setAndUnsetStyleBits(ES_RIGHT, ES_LEFT|ES_CENTER); err != nil {
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

func (ne *NumberEdit) Font() *drawing.Font {
	return ne.edit.Font()
}

func (ne *NumberEdit) SetFont(value *drawing.Font) {
	ne.edit.SetFont(value)
}

func (*NumberEdit) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (ne *NumberEdit) PreferredSize() drawing.Size {
	return ne.dialogBaseUnitsToPixels(drawing.Size{50, 14})
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
	return ne.edit.SetText(strconv.Ftoa64(value, 'f', ne.decimals))
}

func (ne *NumberEdit) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_NOTIFY:
		switch ((*NMHDR)(unsafe.Pointer(msg.LParam))).Code {
		case UDN_DELTAPOS:
			nmud := (*NMUPDOWN)(unsafe.Pointer(msg.LParam))
			val := ne.Value()
			val -= float64(nmud.IDelta) * ne.increment
			if err := ne.SetValue(val); err != nil {
				log.Println(err)
			}
		}

	case WM_SIZE, WM_SIZING:
		cb, err := ne.ClientBounds()
		if err != nil {
			log.Println(err)
			break
		}
		if err = ne.edit.SetBounds(cb); err != nil {
			log.Println(err)
			break
		}
		SendMessage(ne.hWndUpDown, UDM_SETBUDDY, uintptr(ne.edit.hWnd), 0)
	}

	return ne.Widget.wndProc(msg, origWndProcPtr)
}
