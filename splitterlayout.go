// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"sort"

	"github.com/lxn/win"
)

type splitterLayout struct {
	container   Container
	orientation Orientation
	margins     Margins
	hwnd2Item   map[win.HWND]*splitterLayoutItem
	resetNeeded bool
}

type splitterLayoutItem struct {
	size            int
	oldExplicitSize int
	stretchFactor   int
	growth          int
	fixed           bool
	keepSize        bool
}

func newSplitterLayout(orientation Orientation) *splitterLayout {
	return &splitterLayout{
		orientation: orientation,
		hwnd2Item:   make(map[win.HWND]*splitterLayoutItem),
	}
}

func (l *splitterLayout) Container() Container {
	return l.container
}

func (l *splitterLayout) SetContainer(value Container) {
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

func (l *splitterLayout) Margins() Margins {
	return l.margins
}

func (l *splitterLayout) SetMargins(value Margins) error {
	l.margins = value

	return l.Update(false)
}

func (l *splitterLayout) Spacing() int {
	return l.container.(*Splitter).handleWidth
}

func (l *splitterLayout) SetSpacing(value int) error {
	return newError("not supported")
}

func (l *splitterLayout) Orientation() Orientation {
	return l.orientation
}

func (l *splitterLayout) SetOrientation(value Orientation) error {
	if value != l.orientation {
		switch value {
		case Horizontal, Vertical:

		default:
			return newError("invalid Orientation value")
		}

		l.orientation = value

		l.Update(false)
	}

	return nil
}

func (l *splitterLayout) Fixed(widget Widget) bool {
	item := l.hwnd2Item[widget.Handle()]
	return item != nil && item.fixed
}

func (l *splitterLayout) StretchFactor(widget Widget) int {
	item := l.hwnd2Item[widget.Handle()]
	if item == nil || item.stretchFactor == 0 {
		return 1
	}

	return item.stretchFactor
}

func (l *splitterLayout) SetStretchFactor(widget Widget, factor int) error {
	if factor != l.StretchFactor(widget) {
		if factor < 1 {
			return newError("factor must be >= 1")
		}

		if l.container == nil {
			return newError("container required")
		}

		item := l.hwnd2Item[widget.Handle()]
		if item == nil {
			item = new(splitterLayoutItem)
			l.hwnd2Item[widget.Handle()] = item
		}

		item.stretchFactor = factor

		l.Update(false)
	}

	return nil
}

func (l *splitterLayout) cleanupItems() {
	widgets := l.container.Children()

	for handle, _ := range l.hwnd2Item {
		if !widgets.containsHandle(handle) {
			delete(l.hwnd2Item, handle)
		}
	}
}

func (l *splitterLayout) LayoutFlags() LayoutFlags {
	if l.container == nil {
		return 0
	}

	return boxLayoutFlags(l.orientation, l.container.Children())
}

func (l *splitterLayout) MinSize() Size {
	if l.container == nil {
		return Size{}
	}

	return l.MinSizeForSize(l.container.ClientBounds().Size())
}

func (l *splitterLayout) MinSizeForSize(size Size) Size {
	margins := Size{l.margins.HNear + l.margins.HFar, l.margins.VNear + l.margins.VFar}
	s := margins

	anyNonFixed := l.anyNonFixed()

	for _, widget := range l.container.Children().items {
		var cur Size

		if anyNonFixed && l.Fixed(widget) {
			cur = widget.Size()

			if l.orientation == Horizontal {
				cur.Height = 0
			} else {
				cur.Width = 0
			}
		} else {
			cur = minSizeEffective(widget)
		}

		if l.orientation == Horizontal {
			s.Width += cur.Width
			s.Height = maxi(s.Height, margins.Height+cur.Height)
		} else {
			s.Height += cur.Height
			s.Width = maxi(s.Width, margins.Width+cur.Width)
		}
	}

	return s
}

func (l *splitterLayout) anyNonFixed() bool {
	for i, widget := range l.container.Children().items {
		if i%2 == 0 && !l.Fixed(widget) {
			return true
		}
	}

	return false
}

func (l *splitterLayout) spaceForRegularWidgets() int {
	splitter := l.container.(*Splitter)
	s := splitter.ClientBounds().Size()

	var space int
	if l.orientation == Horizontal {
		space = s.Width - l.margins.HNear - l.margins.HFar
	} else {
		space = s.Height - l.margins.VNear - l.margins.VFar
	}

	return space - (splitter.Children().Len()/2)*splitter.handleWidth
}

func (l *splitterLayout) reset() {
	l.cleanupItems()

	children := l.container.Children()
	minSizes := make([]int, children.Len())
	var minSizesTotal int
	for i, w := range children.items {
		if i%2 == 1 {
			continue
		}

		min := minSizeEffective(w)
		if l.orientation == Horizontal {
			minSizes[i] = min.Width
			minSizesTotal += min.Width
		} else {
			minSizes[i] = min.Height
			minSizesTotal += min.Height
		}
	}
	regularSpace := l.spaceForRegularWidgets()

	stretchTotal := 0
	for i := children.Len() - 1; i >= 0; i-- {
		if i%2 == 1 {
			continue
		}

		child := children.At(i)

		if item := l.hwnd2Item[child.Handle()]; item == nil {
			l.hwnd2Item[child.Handle()] = &splitterLayoutItem{stretchFactor: 1}
		}

		stretchTotal += l.StretchFactor(child)
	}

	for i := children.Len() - 1; i >= 0; i-- {
		if i%2 == 1 {
			continue
		}

		child := children.At(i)

		item := l.hwnd2Item[child.Handle()]
		item.growth = 0
		item.keepSize = false
		if item.oldExplicitSize > 0 {
			item.size = item.oldExplicitSize
		} else {
			item.size = int(float64(l.StretchFactor(child)) / float64(stretchTotal) * float64(regularSpace))
		}

		min := minSizes[i]
		if minSizesTotal <= regularSpace {
			if item.size < min {
				item.size = min
			}
		}

		if item.size >= min {
			flags := child.LayoutFlags()

			if l.orientation == Horizontal && flags&GrowableHorz == 0 || l.orientation == Vertical && flags&GrowableVert == 0 {
				item.size = min
				item.keepSize = true
			}
		}
	}
}

func (l *splitterLayout) Update(reset bool) error {
	if l.container == nil {
		return newError("container required")
	}

	if reset {
		l.resetNeeded = true
	}

	if l.container.Suspended() {
		return nil
	}

	if l.resetNeeded {
		reset = true
		l.resetNeeded = false

		l.reset()
	}

	widgets := l.container.Children().items
	splitter := l.container.(*Splitter)
	handleWidth := splitter.HandleWidth()
	sizes := make([]int, len(widgets))
	cb := splitter.ClientBounds()
	cb.X += l.margins.HNear
	cb.Y += l.margins.HFar
	cb.Width -= l.margins.HNear + l.margins.HFar
	cb.Height -= l.margins.VNear + l.margins.VFar
	space1 := l.spaceForRegularWidgets()

	var space2 int
	if l.orientation == Horizontal {
		space2 = cb.Height
	} else {
		space2 = cb.Width
	}

	anyNonFixed := l.anyNonFixed()
	var totalRegularSize int
	for i, widget := range widgets {
		if i%2 == 0 {
			size := l.hwnd2Item[widget.Handle()].size
			totalRegularSize += size
			sizes[i] = size
		} else {
			sizes[i] = handleWidth
		}
	}

	diff := space1 - totalRegularSize

	if diff != 0 && len(sizes) > 1 {
		type WidgetItem struct {
			item  *splitterLayoutItem
			index int
			min   int
			max   int
		}

		var wis []WidgetItem

		for i, widget := range widgets {
			if i%2 == 0 {
				if item := l.hwnd2Item[widget.Handle()]; !anyNonFixed || !item.fixed {
					var min, max int

					minSize := minSizeEffective(widget)
					maxSize := widget.MaxSize()

					if l.orientation == Horizontal {
						min = minSize.Width
						max = maxSize.Width
					} else {
						min = minSize.Height
						max = maxSize.Height
					}

					wis = append(wis, WidgetItem{item: item, index: i, min: min, max: max})
				}
			}
		}

		for diff != 0 {
			sort.SliceStable(wis, func(i, j int) bool {
				a := wis[i]
				b := wis[j]

				x := float64(a.item.growth) / float64(a.item.stretchFactor)
				y := float64(b.item.growth) / float64(b.item.stretchFactor)

				if diff > 0 {
					return x < y && (a.max == 0 || a.max > a.item.size)
				} else {
					return x > y && a.min < a.item.size
				}
			})

			var wi *WidgetItem
			for _, wItem := range wis {
				if !wItem.item.keepSize {
					wi = &wItem
					break
				}
			}
			if wi == nil {
				break
			}

			if diff > 0 {
				if wi.max > 0 && wi.item.size >= wi.max {
					break
				}

				sizes[wi.index]++
				wi.item.size++
				wi.item.growth++
				diff--
			} else {
				if wi.item.size <= wi.min {
					break
				}

				sizes[wi.index]--
				wi.item.size--
				wi.item.growth--
				diff++
			}
		}
	}

	maybeInvalidate := l.container.AsContainerBase().hasComplexBackground()

	hdwp := win.BeginDeferWindowPos(int32(len(widgets)))
	if hdwp == 0 {
		return lastError("BeginDeferWindowPos")
	}

	var p1 int
	if l.orientation == Horizontal {
		p1 = l.margins.HNear
	} else {
		p1 = l.margins.VNear
	}
	for i, widget := range widgets {
		s1 := sizes[i]

		var x, y, w, h int
		if l.orientation == Horizontal {
			x, y, w, h = p1, l.margins.VNear, s1, space2
		} else {
			x, y, w, h = l.margins.HNear, p1, space2, s1
		}

		if maybeInvalidate {
			if b := widget.Bounds(); w == b.Width && h == b.Height && (x != b.X || y != b.Y) {
				widget.Invalidate()
			}
		}

		if hdwp = win.DeferWindowPos(
			hdwp,
			widget.Handle(),
			0,
			int32(x),
			int32(y),
			int32(w),
			int32(h),
			win.SWP_NOACTIVATE|win.SWP_NOOWNERZORDER|win.SWP_NOZORDER); hdwp == 0 {

			return lastError("DeferWindowPos")
		}

		p1 += s1
	}

	if !win.EndDeferWindowPos(hdwp) {
		return lastError("EndDeferWindowPos")
	}

	return nil
}
