// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type ProgressBar struct {
	WidgetBase
}

func NewProgressBar(parent Container) (*ProgressBar, error) {
	pb := &ProgressBar{}

	if err := initChildWidget(
		pb,
		parent,
		"msctls_progress32",
		WS_VISIBLE,
		0); err != nil {
		return nil, err
	}

	return pb, nil
}

func (*ProgressBar) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz | GreedyHorz
}

func (pb *ProgressBar) MinSizeHint() Size {
	return pb.dialogBaseUnitsToPixels(Size{10, 14})
}

func (pb *ProgressBar) SizeHint() Size {
	return pb.dialogBaseUnitsToPixels(Size{50, 14})
}

func (pb *ProgressBar) MinValue() int {
	return int(SendMessage(pb.hWnd, PBM_GETRANGE, 1, 0))
}

func (pb *ProgressBar) MaxValue() int {
	return int(SendMessage(pb.hWnd, PBM_GETRANGE, 0, 0))
}

func (pb *ProgressBar) SetRange(min, max int) {
	SendMessage(pb.hWnd, PBM_SETRANGE32, uintptr(min), uintptr(max))
}

func (pb *ProgressBar) Value() int {
	return int(SendMessage(pb.hWnd, PBM_GETPOS, 0, 0))
}

func (pb *ProgressBar) SetValue(value int) {
	SendMessage(pb.hWnd, PBM_SETPOS, uintptr(value), 0)
}
