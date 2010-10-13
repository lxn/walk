// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printing

import (
	"os"
	"utf8"
)

import (
	"walk/drawing"
)

type simpleTextPart struct {
	item   *simpleTextItem
	offset int
	length int
	bounds drawing.Rectangle
}

func (part *simpleTextPart) Bounds() drawing.Rectangle {
	return part.bounds
}

func (part *simpleTextPart) Draw(surface *drawing.Surface) os.Error {
	item := part.item
	text := item.text.Slice(part.offset, part.offset+part.length)

	return surface.DrawText(text, item.font, item.color, part.bounds, item.format)
}

type simpleTextItem struct {
	text          *utf8.String
	font          *drawing.Font
	color         drawing.Color
	preferredSize drawing.Size
	format        drawing.DrawTextFormat
	parts         []*simpleTextPart
}

func (item *simpleTextItem) Dispose() {
	item.font.Dispose()
}

func (item *simpleTextItem) PartCount() int {
	return len(item.parts)
}

func (item *simpleTextItem) Part(i int) part {
	return item.parts[i]
}

func (item *simpleTextItem) AddNewPart(surface *drawing.Surface, bounds drawing.Rectangle) (part part, more bool, err os.Error) {
	partCount := len(item.parts)
	var offset int
	if partCount > 0 {
		prevPart := item.parts[len(item.parts)-1]
		offset = prevPart.offset + prevPart.length
	}

	runeCount := item.text.RuneCount()
	text := item.text.Slice(offset, runeCount)

	fontHeight, err := surface.FontHeight(item.font)
	if err != nil {
		return
	}

	bounds.Height = (bounds.Height / fontHeight) * fontHeight
	if bounds.Height == 0 {
		more = true
		return
	}

	boundsMeasured, runesFitted, err := surface.MeasureText(text, item.font, bounds, item.format)
	if err != nil {
		return
	}

	p := &simpleTextPart{
		item:   item,
		offset: offset,
		length: runesFitted,
		bounds: boundsMeasured,
	}

	if partCount == cap(item.parts) {
		parts := make([]*simpleTextPart, partCount, partCount*2)
		copy(parts, item.parts)
		item.parts = parts
	}

	item.parts = item.parts[0 : partCount+1]
	item.parts[partCount] = p

	part = p
	more = p.offset+p.length < runeCount

	return
}

func (item *simpleTextItem) NextPartMinSize() drawing.Size {
	return drawing.Size{}
}

func (item *simpleTextItem) PreferredSize() drawing.Size {
	return item.preferredSize
}
