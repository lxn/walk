// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/miu200521358/win"
)

type ProgressBar struct {
	WidgetBase
}

func NewProgressBar(parent Container) (*ProgressBar, error) {
	pb := new(ProgressBar)

	if err := InitWidget(
		pb,
		parent,
		"msctls_progress32",
		win.WS_VISIBLE,
		0); err != nil {
		return nil, err
	}

	return pb, nil
}

func (pb *ProgressBar) MinValue() int {
	return int(pb.SendMessage(win.PBM_GETRANGE, 1, 0))
}

func (pb *ProgressBar) MaxValue() int {
	return int(pb.SendMessage(win.PBM_GETRANGE, 0, 0))
}

func (pb *ProgressBar) SetRange(min, max int) {
	pb.SendMessage(win.PBM_SETRANGE32, uintptr(min), uintptr(max))
}

func (pb *ProgressBar) Value() int {
	return int(pb.SendMessage(win.PBM_GETPOS, 0, 0))
}

func (pb *ProgressBar) SetValue(value int) {
	pb.SendMessage(win.PBM_SETPOS, uintptr(value), 0)
}

func (pb *ProgressBar) MarqueeMode() bool {
	return pb.hasStyleBits(win.PBS_MARQUEE)
}

func (pb *ProgressBar) SetMarqueeMode(marqueeMode bool) error {
	if err := pb.ensureStyleBits(win.PBS_MARQUEE, marqueeMode); err != nil {
		return err
	}

	pb.SendMessage(win.PBM_SETMARQUEE, uintptr(win.BoolToBOOL(marqueeMode)), 0)

	return nil
}

func (pb *ProgressBar) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	return &progressBarLayoutItem{
		idealSize: pb.dialogBaseUnitsToPixels(Size{50, 14}),
		minSize:   pb.dialogBaseUnitsToPixels(Size{10, 14}),
	}
}

type progressBarLayoutItem struct {
	LayoutItemBase
	idealSize Size // in native pixels
	minSize   Size // in native pixels
}

func (*progressBarLayoutItem) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz | GreedyHorz
}

func (li *progressBarLayoutItem) IdealSize() Size {
	return li.idealSize
}

func (li *progressBarLayoutItem) MinSize() Size {
	return li.minSize
}
