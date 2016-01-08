// Copyright 2016 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"log"    // TODO remove
	"unsafe" // TODO remove
)

import (
	"github.com/lxn/win"
)

type Slider struct {
	WidgetBase
	valueChangedPublisher EventPublisher
	layoutFlags           LayoutFlags
}

func NewSlider(parent Container) (*Slider, error) {
	return NewSliderWithOrientation(parent, Horizontal)
}

func NewSliderWithOrientation(parent Container, orientation Orientation) (*Slider, error) {
	sl := new(Slider)

	var style uint32 = win.WS_VISIBLE | win.TBS_TOOLTIPS
	if orientation == Vertical {
		style |= win.TBS_VERT
		sl.layoutFlags = ShrinkableVert | GrowableVert | GreedyVert
	} else {
		sl.layoutFlags = ShrinkableHorz | GrowableHorz | GreedyHorz
	}

	if err := InitWidget(
		sl,
		parent,
		"msctls_trackbar32",
		style,
		0); err != nil {
		return nil, err
	}

	sl.MustRegisterProperty("Value", NewProperty(
		func() interface{} {
			return sl.Value()
		},
		func(v interface{}) error {
			sl.SetValue(v.(int))
			return nil
		},
		sl.valueChangedPublisher.Event()))

	return sl, nil
}

func (sl *Slider) LayoutFlags() LayoutFlags {
	return sl.layoutFlags
}

func (sl *Slider) SizeHint() Size {
	return sl.MinSizeHint()
}

func (sl *Slider) MinSizeHint() Size {
	return sl.dialogBaseUnitsToPixels(Size{20, 20})
}

func (sl *Slider) MinValue() int {
	return int(sl.SendMessage(win.TBM_GETRANGEMIN, 0, 0))
}

func (sl *Slider) MaxValue() int {
	return int(sl.SendMessage(win.TBM_GETRANGEMAX, 0, 0))
}

func (sl *Slider) SetRange(min, max int) {
	sl.SendMessage(win.TBM_SETRANGEMIN, 0, uintptr(min))
	sl.SendMessage(win.TBM_SETRANGEMAX, 1, uintptr(max))
}

func (sl *Slider) Value() int {
	return int(sl.SendMessage(win.TBM_GETPOS, 0, 0))
}

func (sl *Slider) SetValue(value int) {
	sl.SendMessage(win.TBM_SETPOS, 1, uintptr(value))
	sl.valueChangedPublisher.Publish()
}

// ValueChanged returns an Event that can be used to track changes to Value.
func (sl *Slider) ValueChanged() *Event {
	return sl.valueChangedPublisher.Event()
}

// I'm not sure anymore, if this is the right way!
// It seems the proper way is to hook on WM_HSCROLL and WM_VSCROLL.
// But they are not received here!
// I read, a trackbar does not receive them, it sends them to it's parents.
// Any Idea on what to do then?
func (sl *Slider) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_NOTIFY: // 0x4E
		code := ((*win.NMHDR)(unsafe.Pointer(lParam))).Code
		switch code {
		case win.TRBN_THUMBPOSCHANGING:
			// TRBN_THUMBPOSCHANGING is documented to exist only on Vista and above.
			// I run windows 10, still I never recieve it.
			log.Fatal("RECEIVED TRBN_THUMBPOSCHANGING") // This never happens.
			sl.valueChangedPublisher.Publish()
		case win.NM_CUSTOMDRAW:
			// This is received however, show we do anything with it?
			nmcd := (*win.NMCUSTOMDRAW)(unsafe.Pointer(lParam))
			log.Printf("WM_NOTIFY NM_CUSTOMDRAW: %+v code=%x\n", nmcd, nmcd.Hdr.Code)
		default:
			// I don't receive anything else, when I change the slider with an arrow key.
			// There are some others however, when I drag with the mouse.
			log.Printf("Other WM_NOFITY NMHDR CODE: 0x%x\n", code)
		}
	case win.WM_GETDLGCODE:
		// I recieve this too. lParm is 0x8fca8
		log.Printf("WM_GETDLGCODE: 0x%x 0x%x", wParam, lParam)
	default:
		log.Printf("Other msg: %v 0x%x\n", msg, msg)
		// I receive these other messages, when I press an arrow key on a slider:
		// 0x101  WM_KEYUP
		// 0x100  WM_KEYFIRST
		// 0xf    WM_PAINT
		// There is no WM_VSROLL (0x114) or WM_HSCROLL (0x115).
	}
	return sl.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
