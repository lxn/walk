// Copyright 2018 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"sync"
)

type TextLabel struct {
	static
	textChangedPublisher EventPublisher
}

func NewTextLabel(parent Container) (*TextLabel, error) {
	return NewTextLabelWithStyle(parent, 0)
}

func NewTextLabelWithStyle(parent Container, style uint32) (*TextLabel, error) {
	tl := new(TextLabel)

	if err := tl.init(tl, parent); err != nil {
		return nil, err
	}

	tl.textAlignment = AlignHNearVNear

	tl.MustRegisterProperty("Text", NewProperty(
		func() interface{} {
			return tl.Text()
		},
		func(v interface{}) error {
			return tl.SetText(assertStringOr(v, ""))
		},
		tl.textChangedPublisher.Event()))

	return tl, nil
}

func (tl *TextLabel) asStatic() *static {
	return &tl.static
}

func (tl *TextLabel) TextAlignment() Alignment2D {
	return tl.textAlignment
}

func (tl *TextLabel) SetTextAlignment(alignment Alignment2D) error {
	if alignment == AlignHVDefault {
		alignment = AlignHNearVNear
	}

	return tl.setTextAlignment(alignment)
}

func (tl *TextLabel) Text() string {
	return tl.text()
}

func (tl *TextLabel) SetText(text string) error {
	if changed, err := tl.setText(text); err != nil {
		return err
	} else if !changed {
		return nil
	}

	tl.textChangedPublisher.Publish()

	return nil
}

func (tl *TextLabel) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	return &textLabelLayoutItem{
		width2Height: make(map[int]int),
		text:         tl.Text(),
		font:         tl.Font(),
		minWidth:     tl.MinSizePixels().Width,
	}
}

type textLabelLayoutItem struct {
	LayoutItemBase
	mutex        sync.Mutex
	width2Height map[int]int // in native pixels
	text         string
	font         *Font
	minWidth     int // in native pixels
}

func (*textLabelLayoutItem) LayoutFlags() LayoutFlags {
	return GrowableHorz | GrowableVert
}

func (li *textLabelLayoutItem) IdealSize() Size {
	return li.MinSize()
}

func (li *textLabelLayoutItem) MinSize() Size {
	return calculateTextSize(li.text, li.font, li.ctx.dpi, li.minWidth, li.handle)
}

func (li *textLabelLayoutItem) HasHeightForWidth() bool {
	return true
}

func (li *textLabelLayoutItem) HeightForWidth(width int) int {
	li.mutex.Lock()
	defer li.mutex.Unlock()

	if height, ok := li.width2Height[width]; ok {
		return height
	}

	size := calculateTextSize(li.text, li.font, li.ctx.dpi, width, li.handle)

	li.width2Height[width] = size.Height

	return size.Height
}
