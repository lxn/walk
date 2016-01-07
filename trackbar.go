// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

// TODO: should we call this TrackBar, Trackbar or Slider?
type TrackBar struct {
	WidgetBase
	valueChangedPublisher EventPublisher
}

func NewTrackBar(parent Container) (*TrackBar, error) {
	return NewTrackBarWithOrientation(parent, Horizontal)
}

func NewTrackBarWithOrientation(parent Container, orientation Orientation) (*TrackBar, error) {
	tb := new(TrackBar)

	var style uint32 = win.WS_VISIBLE | win.TBS_TOOLTIPS
	if orientation == Vertical {
		style |= win.TBS_VERT
	}

	if err := InitWidget(
		tb,
		parent,
		"msctls_trackbar32",
		style,
		0); err != nil {
		return nil, err
	}

	tb.MustRegisterProperty("Value", NewProperty(
		func() interface{} {
			return tb.Value()
		},
		func(v interface{}) error {
			tb.SetValue(v.(int))
			return nil
		},
		tb.valueChangedPublisher.Event()))

	return tb, nil
}

func (tb *TrackBar) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
	// TODO: is this all needed? Or should it depend on the orientation?
}

func (tb *TrackBar) SizeHint() Size {
	return Size{100, 100}
	// TODO: is there anything usefull to set?
}

func (tb *TrackBar) MinSizeHint() Size {
	return tb.dialogBaseUnitsToPixels(Size{20, 20})
}

// TODO: The following methods are similar to progressbar's.
// What should be the type of the value:
// int, uint[32], float64 or interface{}?

func (tb *TrackBar) MinValue() int {
	return int(tb.SendMessage(win.TBM_GETRANGEMIN, 0, 0))
}

func (tb *TrackBar) MaxValue() int {
	return int(tb.SendMessage(win.TBM_GETRANGEMAX, 0, 0))
}

func (tb *TrackBar) SetRange(min, max int) {
	tb.SendMessage(win.TBM_SETRANGEMIN, 0, uintptr(min))
	tb.SendMessage(win.TBM_SETRANGEMAX, 1, uintptr(max))
}

func (tb *TrackBar) Value() int {
	return int(tb.SendMessage(win.TBM_GETPOS, 0, 0))
}

func (tb *TrackBar) SetValue(value int) {
	tb.SendMessage(win.TBM_SETPOS, 1, uintptr(value))
	tb.valueChangedPublisher.Publish()
}

// ValueChanged returns an Event that can be used to track changes to Value.
func (tb *TrackBar) ValueChanged() *Event {
	return tb.valueChangedPublisher.Event()
}

func (tb *TrackBar) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	// TODO; Is this the right message to grab?
	case win.WM_NOTIFY:
		// TODO: That's all there is to do?
		tb.valueChangedPublisher.Publish()
	}
	return tb.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
