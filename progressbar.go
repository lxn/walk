// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

var progressBarOrigWndProcPtr uintptr
var _ subclassedWidget = &ProgressBar{}

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

func (*ProgressBar) origWndProcPtr() uintptr {
	return progressBarOrigWndProcPtr
}

func (*ProgressBar) setOrigWndProcPtr(ptr uintptr) {
	progressBarOrigWndProcPtr = ptr
}

func (*ProgressBar) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz | GreedyHorz
}

func (pb *ProgressBar) SizeHint() Size {
	return pb.dialogBaseUnitsToPixels(Size{50, 14})
}

func (pb *ProgressBar) ProgressPercent() int {
	return int(SendMessage(pb.hWnd, PBM_GETPOS, 0, 0))
}

func (pb *ProgressBar) SetProgressPercent(value int) {
	SendMessage(pb.hWnd, PBM_SETPOS, uintptr(value), 0)
}
