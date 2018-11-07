// Copyright 2018 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type FlowLayout struct {
	container          Container
	hwnd2StretchFactor map[win.HWND]int
	margins            Margins
	spacing            int
	resetNeeded        bool
}

type flowLayoutSection struct {
	items            []flowLayoutSectionItem
	primarySpaceLeft int
	secondaryMinSize int
}

type flowLayoutSectionItem struct {
	widget  Widget
	minSize Size
}

func NewFlowLayout() *FlowLayout {
	return new(FlowLayout)
}

func (l *FlowLayout) Container() Container {
	return l.container
}

func (l *FlowLayout) SetContainer(value Container) {
	if value != l.container {
		if l.container != nil {
			l.container.SetLayout(nil)
		}

		l.container = value

		if value != nil && value.Layout() != Layout(l) {
			value.SetLayout(l)

			l.Update(true)
		}
	}
}

func (l *FlowLayout) Margins() Margins {
	return l.margins
}

func (l *FlowLayout) SetMargins(value Margins) error {
	if value.HNear < 0 || value.VNear < 0 || value.HFar < 0 || value.VFar < 0 {
		return newError("margins must be positive")
	}

	l.margins = value

	return nil
}

func (l *FlowLayout) Spacing() int {
	return l.spacing
}

func (l *FlowLayout) SetSpacing(value int) error {
	if value != l.spacing {
		if value < 0 {
			return newError("spacing cannot be negative")
		}

		l.spacing = value

		l.Update(false)
	}

	return nil
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

		l.Update(false)
	}

	return nil
}

func (l *FlowLayout) cleanupStretchFactors() {
	widgets := l.container.Children()

	for handle, _ := range l.hwnd2StretchFactor {
		if !widgets.containsHandle(handle) {
			delete(l.hwnd2StretchFactor, handle)
		}
	}
}

func (l *FlowLayout) LayoutFlags() LayoutFlags {
	if l.container == nil {
		return 0
	}

	var flags LayoutFlags

	children := l.container.Children()
	count := children.Len()
	if count == 0 {
		return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert
	} else {
		for i := 0; i < count; i++ {
			widget := children.At(i)

			if shouldLayoutWidget(widget) {
				wf := widget.LayoutFlags()

				if wf&GreedyHorz != 0 && widget.MaxSize().Width > 0 {
					wf &^= GreedyHorz
				}
				if wf&GreedyVert != 0 && widget.MaxSize().Height > 0 {
					wf &^= GreedyVert
				}

				flags |= wf
			}
		}
	}

	return flags
}

func (l *FlowLayout) MinSize() Size {
	if l.container == nil {
		return Size{}
	}

	children := l.container.Children()

	var width, height int
	for i := children.Len() - 1; i >= 0; i-- {
		widget := children.At(i)
		if !shouldLayoutWidget(widget) {
			continue
		}

		minSize := minSizeEffective(widget)

		width = maxi(minSize.Width, width)
		height = maxi(minSize.Height, height)
	}

	return Size{width + l.margins.HNear + l.margins.HFar, height + l.margins.VNear + l.margins.VFar}
}

func (l *FlowLayout) minSizeForWidth(w int) Size {
	if l.container == nil {
		return Size{}
	}

	sections := l.sectionsForPrimarySize(w)

	var width, height int
	for i, section := range sections {
		width = maxi(width, w-section.primarySpaceLeft)

		if i > 0 {
			height += l.spacing
		}
		height += section.secondaryMinSize
	}

	return Size{width + l.margins.HNear + l.margins.HFar, height + l.margins.VNear + l.margins.VFar}
}

func (l *FlowLayout) Update(reset bool) error {
	if l.container == nil {
		return nil
	}

	if reset {
		l.resetNeeded = true
	}

	if l.container.Suspended() {
		return nil
	}

	if !performingScheduledLayouts && scheduleLayout(l) {
		return nil
	}

	if l.resetNeeded {
		l.resetNeeded = false

		l.cleanupStretchFactors()
	}

	ifContainerIsScrollViewDoCoolSpecialLayoutStuff(l)

	bounds := l.container.ClientBounds()
	sections := l.sectionsForPrimarySize(bounds.Width)

	for _, section := range sections {
		var widgets []Widget
		for _, item := range section.items {
			widgets = append(widgets, item.widget)
		}

		bounds.Height = section.secondaryMinSize

		if err := performBoxLayout(widgets, Horizontal, bounds, l.Margins(), l.Spacing(), l.hwnd2StretchFactor); err != nil {
			return err
		}

		bounds.Y += bounds.Height + l.spacing
	}

	return nil
}

func (l *FlowLayout) sectionsForPrimarySize(primarySize int) []flowLayoutSection {
	children := l.container.Children()
	count := children.Len()

	var sections []flowLayoutSection

	section := flowLayoutSection{
		primarySpaceLeft: primarySize - l.margins.HNear - l.margins.HFar,
	}

	addSection := func() {
		sections = append(sections, section)
		section.items = nil
		section.primarySpaceLeft = primarySize - l.margins.HNear - l.margins.HFar
		section.secondaryMinSize = 0
	}

	for i := 0; i < count; i++ {
		var item flowLayoutSectionItem

		item.widget = children.At(i)

		if !shouldLayoutWidget(item.widget) {
			continue
		}

		item.minSize = minSizeEffective(item.widget)

		addItem := func() {
			section.items = append(section.items, item)
			section.primarySpaceLeft -= item.minSize.Width
			section.secondaryMinSize = maxi(section.secondaryMinSize, item.minSize.Height)
		}

		if section.primarySpaceLeft < item.minSize.Width && len(section.items) == 0 {
			addItem()
			addSection()
		} else if section.primarySpaceLeft < l.spacing+item.minSize.Width {
			addSection()
			addItem()
		} else {
			addItem()
		}
	}

	if len(section.items) > 0 {
		addSection()
	}

	return sections
}
