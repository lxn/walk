// Copyright 2018 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type FlowLayout struct {
	LayoutBase
	hwnd2StretchFactor map[win.HWND]int
}

func NewFlowLayout() *FlowLayout {
	l := &FlowLayout{
		LayoutBase: LayoutBase{
			margins96dpi: Margins{9, 9, 9, 9},
			spacing96dpi: 6,
		},
		hwnd2StretchFactor: make(map[win.HWND]int),
	}
	l.layout = l

	return l
}

func (l *FlowLayout) StretchFactor(widget Widget) int {
	if factor, ok := l.hwnd2StretchFactor[widget.Handle()]; ok {
		return factor
	}

	return 1
}

func (l *FlowLayout) SetStretchFactor(widget Widget, factor int) error {
	if factor != l.StretchFactor(widget) {
		if l.container == nil {
			return newError("container required")
		}

		handle := widget.Handle()

		if !l.container.Children().containsHandle(handle) {
			return newError("unknown widget")
		}
		if factor < 1 {
			return newError("factor must be >= 1")
		}

		l.hwnd2StretchFactor[handle] = factor

		l.container.RequestLayout()
	}

	return nil
}

func (l *FlowLayout) CreateLayoutItem(ctx *LayoutContext) ContainerLayoutItem {
	li := &flowLayoutItem{
		size2MinSize:       make(map[Size]Size),
		hwnd2StretchFactor: make(map[win.HWND]int),
	}

	for hwnd, sf := range l.hwnd2StretchFactor {
		li.hwnd2StretchFactor[hwnd] = sf
	}

	return li
}

type flowLayoutItem struct {
	ContainerLayoutItemBase
	size2MinSize       map[Size]Size // in native pixels
	hwnd2StretchFactor map[win.HWND]int
}

type flowLayoutSection struct {
	items            []flowLayoutSectionItem
	primarySpaceLeft int // in native pixels
	secondaryMinSize int // in native pixels
}

type flowLayoutSectionItem struct {
	item    LayoutItem
	minSize Size // in native pixels
}

func (*flowLayoutItem) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (li *flowLayoutItem) MinSize() Size {
	return li.MinSizeForSize(li.geometry.ClientSize)
}

func (li *flowLayoutItem) HeightForWidth(width int) int {
	return li.MinSizeForSize(Size{width, li.geometry.ClientSize.Height}).Height
}

func (li *flowLayoutItem) MinSizeForSize(size Size) Size {
	if min, ok := li.size2MinSize[size]; ok {
		return min
	}

	spacing := IntFrom96DPI(li.spacing96dpi, li.ctx.dpi)
	margins := MarginsFrom96DPI(li.margins96dpi, li.ctx.dpi)

	bounds := Rectangle{Width: size.Width}

	sections := li.sectionsForPrimarySize(size.Width)

	var s Size
	var maxPrimary int

	for i, section := range sections {
		var items []LayoutItem
		var sectionMinWidth int
		for _, sectionItem := range section.items {
			items = append(items, sectionItem.item)

			sectionMinWidth += sectionItem.minSize.Width
		}
		sectionMinWidth += (len(section.items) - 1) * spacing
		maxPrimary = maxi(maxPrimary, sectionMinWidth)

		bounds.Height = section.secondaryMinSize

		margins96dpi := li.margins96dpi
		if i > 0 {
			margins96dpi.VNear = 0
		}
		if i < len(sections)-1 {
			margins96dpi.VFar = 0
		}

		layoutItems := boxLayoutItems(li, items, Horizontal, li.alignment, bounds, margins96dpi, li.spacing96dpi, li.hwnd2StretchFactor)

		var maxSecondary int

		for _, item := range layoutItems {
			if hfw, ok := item.Item.(HeightForWidther); ok && hfw.HasHeightForWidth() {
				item.Bounds.Height = hfw.HeightForWidth(item.Bounds.Width)
			} else {
				min := li.MinSizeEffectiveForChild(item.Item)
				item.Bounds.Height = min.Height
			}

			maxSecondary = maxi(maxSecondary, item.Bounds.Height)
		}

		s.Height += maxSecondary

		bounds.Y += maxSecondary + spacing
	}

	s.Width = maxPrimary

	s.Width += margins.HNear + margins.HFar
	s.Height += margins.VNear + margins.VFar + (len(sections)-1)*spacing

	if s.Width > 0 && s.Height > 0 {
		li.size2MinSize[size] = s
	}

	return s
}

func (li *flowLayoutItem) PerformLayout() []LayoutResultItem {
	spacing := IntFrom96DPI(li.spacing96dpi, li.ctx.dpi)
	bounds := Rectangle{Width: li.geometry.ClientSize.Width, Height: li.geometry.ClientSize.Height}

	sections := li.sectionsForPrimarySize(bounds.Width)

	var resultItems []LayoutResultItem

	for i, section := range sections {
		var items []LayoutItem
		for _, sectionItem := range section.items {
			items = append(items, sectionItem.item)
		}

		bounds.Height = section.secondaryMinSize

		margins96dpi := li.margins96dpi
		if i > 0 {
			margins96dpi.VNear = 0
		}
		if i < len(sections)-1 {
			margins96dpi.VFar = 0
		}

		layoutItems := boxLayoutItems(li, items, Horizontal, li.alignment, bounds, margins96dpi, li.spacing96dpi, li.hwnd2StretchFactor)

		margins := MarginsFrom96DPI(margins96dpi, li.ctx.dpi)

		var maxSecondary int

		for _, item := range layoutItems {
			if hfw, ok := item.Item.(HeightForWidther); ok && hfw.HasHeightForWidth() {
				item.Bounds.Height = hfw.HeightForWidth(item.Bounds.Width)
			} else {
				item.Bounds.Height = li.MinSizeEffectiveForChild(item.Item).Height
			}

			maxSecondary = maxi(maxSecondary, item.Bounds.Height)
		}

		bounds.Height = maxSecondary + margins.VNear + margins.VFar

		resultItems = append(resultItems, boxLayoutItems(li, items, Horizontal, li.alignment, bounds, margins96dpi, li.spacing96dpi, li.hwnd2StretchFactor)...)

		bounds.Y += bounds.Height + spacing
	}

	return resultItems
}

// sectionsForPrimarySize calculates sections for primary width in native pixels.
func (li *flowLayoutItem) sectionsForPrimarySize(primarySize int) []flowLayoutSection {
	margins := MarginsFrom96DPI(li.margins96dpi, li.ctx.dpi)
	spacing := IntFrom96DPI(li.spacing96dpi, li.ctx.dpi)

	var sections []flowLayoutSection

	section := flowLayoutSection{
		primarySpaceLeft: primarySize - margins.HNear - margins.HFar,
	}

	addSection := func() {
		sections = append(sections, section)
		section.items = nil
		section.primarySpaceLeft = primarySize - margins.HNear - margins.HFar
		section.secondaryMinSize = 0
	}

	for _, item := range li.children {
		var sectionItem flowLayoutSectionItem

		sectionItem.item = item

		if !shouldLayoutItem(item) {
			continue
		}

		sectionItem.minSize = li.MinSizeEffectiveForChild(item)

		addItem := func() {
			section.items = append(section.items, sectionItem)
			if len(section.items) > 1 {
				section.primarySpaceLeft -= spacing
			}
			section.primarySpaceLeft -= sectionItem.minSize.Width

			section.secondaryMinSize = maxi(section.secondaryMinSize, sectionItem.minSize.Height)
		}

		if section.primarySpaceLeft < sectionItem.minSize.Width && len(section.items) == 0 {
			addItem()
			addSection()
		} else if section.primarySpaceLeft < spacing+sectionItem.minSize.Width && len(section.items) > 0 {
			addSection()
			addItem()
		} else {
			addItem()
		}
	}

	if len(section.items) > 0 {
		addSection()
	}

	if len(sections) > 0 {
		sections[0].secondaryMinSize += margins.VNear
		sections[len(sections)-1].secondaryMinSize += margins.VFar
	}

	return sections
}
