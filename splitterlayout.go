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
	suspended   bool
}

type splitterLayoutItem struct {
	size                 int
	oldExplicitSize      int
	stretchFactor        int
	growth               int
	visibleChangedHandle int
	fixed                bool
	keepSize             bool
	wasVisible           bool
}

func newSplitterLayout(orientation Orientation) *splitterLayout {
	return &splitterLayout{
		orientation: orientation,
		hwnd2Item:   make(map[win.HWND]*splitterLayoutItem),
	}
}

func (l *splitterLayout) asLayoutBase() *LayoutBase {
	return nil
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

			l.container.RequestLayout()
		}
	}
}

func (l *splitterLayout) Margins() Margins {
	return l.margins
}

func (l *splitterLayout) SetMargins(value Margins) error {
	l.margins = value

	l.container.RequestLayout()

	return nil
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

		l.container.RequestLayout()
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

		l.container.RequestLayout()
	}

	return nil
}

func (l *splitterLayout) anyNonFixed() bool {
	for i, widget := range l.container.Children().items {
		if i%2 == 0 && widget.visible && !l.Fixed(widget.window.(Widget)) {
			return true
		}
	}

	return false
}

func (l *splitterLayout) spaceUnavailableToRegularWidgets() int {
	splitter := l.container.(*Splitter)

	var space int
	if l.orientation == Horizontal {
		space = l.margins.HNear + l.margins.HFar
	} else {
		space = l.margins.VNear + l.margins.VFar
	}

	for _, widget := range l.container.Children().items {
		if _, isHandle := widget.window.(*splitterHandle); isHandle && widget.visible {
			space += splitter.handleWidth
		}
	}

	return space
}

func (l splitterLayout) CreateLayoutItem(ctx *LayoutContext) ContainerLayoutItem {
	splitter := l.container.(*Splitter)

	hwnd2Item := make(map[win.HWND]*splitterLayoutItem, len(l.hwnd2Item))
	for hwnd, sli := range l.hwnd2Item {
		hwnd2Item[hwnd] = sli
	}

	li := &splitterContainerLayoutItem{
		orientation:                    l.orientation,
		hwnd2Item:                      hwnd2Item,
		spaceUnavailableToRegularItems: l.spaceUnavailableToRegularWidgets(),
		handleWidth:                    splitter.HandleWidth(),
		anyNonFixed:                    l.anyNonFixed(),
		resetNeeded:                    l.resetNeeded,
	}

	l.resetNeeded = false

	return li
}

type splitterContainerLayoutItem struct {
	ContainerLayoutItemBase
	orientation                    Orientation
	hwnd2Item                      map[win.HWND]*splitterLayoutItem
	spaceUnavailableToRegularItems int
	handleWidth                    int
	anyNonFixed                    bool
	resetNeeded                    bool
}

func (li *splitterContainerLayoutItem) StretchFactor(item LayoutItem) int {
	sli := li.hwnd2Item[item.Handle()]
	if sli == nil || sli.stretchFactor == 0 {
		return 1
	}

	return sli.stretchFactor
}

func (li *splitterContainerLayoutItem) LayoutFlags() LayoutFlags {
	return boxLayoutFlags(li.orientation, li.children)
}

func (li *splitterContainerLayoutItem) MinSize() Size {
	return li.MinSizeForSize(li.geometry.clientSize)
}

func (li *splitterContainerLayoutItem) HeightForWidth(width int) int {
	return li.MinSizeForSize(Size{width, li.geometry.clientSize.Height}).Height
}

func (li *splitterContainerLayoutItem) MinSizeForSize(size Size) Size {
	margins := Size{li.margins.HNear + li.margins.HFar, li.margins.VNear + li.margins.VFar}
	s := margins

	for _, item := range li.children {
		if !anyVisibleItemInHierarchy(item) {
			continue
		}

		var cur Size

		if sli, ok := li.hwnd2Item[item.Handle()]; ok && li.anyNonFixed && sli.fixed {
			cur = item.Geometry().size

			if li.orientation == Horizontal {
				cur.Height = 0
			} else {
				cur.Width = 0
			}
		} else {
			cur = li.MinSizeEffectiveForChild(item)
		}

		if li.orientation == Horizontal {
			s.Width += cur.Width
			s.Height = maxi(s.Height, margins.Height+cur.Height)
		} else {
			s.Height += cur.Height
			s.Width = maxi(s.Width, margins.Width+cur.Width)
		}
	}

	return s
}

func (li *splitterContainerLayoutItem) PerformLayout() []LayoutResultItem {
	if li.resetNeeded {
		li.reset()
	}

	sizes := make([]int, len(li.children))
	cb := Rectangle{Width: li.geometry.clientSize.Width, Height: li.geometry.clientSize.Height}
	cb.X += li.margins.HNear
	cb.Y += li.margins.HFar
	cb.Width -= li.margins.HNear + li.margins.HFar
	cb.Height -= li.margins.VNear + li.margins.VFar

	var space1, space2 int
	if li.orientation == Horizontal {
		space1 = cb.Width - li.spaceUnavailableToRegularItems
		space2 = cb.Height
	} else {
		space1 = cb.Height - li.spaceUnavailableToRegularItems
		space2 = cb.Width
	}

	anyNonFixed := li.anyNonFixed
	var totalRegularSize int
	for i, item := range li.children {
		if !anyVisibleItemInHierarchy(item) {
			continue
		}

		if i%2 == 0 {
			size := li.hwnd2Item[item.Handle()].size
			totalRegularSize += size
			sizes[i] = size
		} else {
			sizes[i] = li.handleWidth
		}
	}

	var resultItems []LayoutResultItem

	diff := space1 - totalRegularSize

	if diff != 0 && len(sizes) > 1 {
		type WidgetItem struct {
			item  *splitterLayoutItem
			index int
			min   int
			max   int
		}

		var wis []WidgetItem

		for i, item := range li.children {
			if !anyVisibleItemInHierarchy(item) {
				continue
			}

			if i%2 == 0 {
				if it := li.hwnd2Item[item.Handle()]; !anyNonFixed || !it.fixed {
					var min, max int

					minSize := li.MinSizeEffectiveForChild(item)
					maxSize := item.Geometry().maxSize

					if li.orientation == Horizontal {
						min = minSize.Width
						max = maxSize.Width
					} else {
						min = minSize.Height
						max = maxSize.Height
					}

					wis = append(wis, WidgetItem{item: it, index: i, min: min, max: max})
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

	var p1 int
	if li.orientation == Horizontal {
		p1 = li.margins.HNear
	} else {
		p1 = li.margins.VNear
	}
	for i, item := range li.children {
		if !anyVisibleItemInHierarchy(item) {
			continue
		}

		s1 := sizes[i]

		var x, y, w, h int
		if li.orientation == Horizontal {
			x, y, w, h = p1, li.margins.VNear, s1, space2
		} else {
			x, y, w, h = li.margins.HNear, p1, space2, s1
		}

		resultItems = append(resultItems, LayoutResultItem{item: item, bounds: Rectangle{x, y, w, h}})

		p1 += s1
	}

	return resultItems
}

func (li *splitterContainerLayoutItem) reset() {
	var anyVisible bool

	for i, item := range li.children {
		sli := li.hwnd2Item[item.Handle()]

		visible := anyVisibleItemInHierarchy(item)
		if !anyVisible && visible {
			anyVisible = true
		}

		if sli == nil || visible == sli.wasVisible {
			continue
		}

		sli.wasVisible = visible

		if _, isHandle := item.(*splitterHandleLayoutItem); !isHandle {
			var handleIndex int

			if i == 0 {
				if len(li.children) > 1 {
					handleIndex = 1
				} else {
					handleIndex = -1
				}
			} else {
				handleIndex = i - 1
			}

			if handleIndex > -1 {
				li.children[handleIndex].AsLayoutItemBase().visible = visible
			}
		}
	}

	if li.Visible() != anyVisible {
		li.AsLayoutItemBase().visible = anyVisible
	}

	minSizes := make([]int, len(li.children))
	var minSizesTotal int
	for i, item := range li.children {
		if i%2 == 1 || !anyVisibleItemInHierarchy(item) {
			continue
		}

		min := li.MinSizeEffectiveForChild(item)
		if li.orientation == Horizontal {
			minSizes[i] = min.Width
			minSizesTotal += min.Width
		} else {
			minSizes[i] = min.Height
			minSizesTotal += min.Height
		}
	}

	var regularSpace int
	if li.orientation == Horizontal {
		regularSpace = li.Geometry().clientSize.Width - li.spaceUnavailableToRegularItems
	} else {
		regularSpace = li.Geometry().clientSize.Height - li.spaceUnavailableToRegularItems
	}

	stretchTotal := 0
	for i, item := range li.children {
		if i%2 == 1 || !anyVisibleItemInHierarchy(item) {
			continue
		}

		if sli := li.hwnd2Item[item.Handle()]; sli == nil {
			li.hwnd2Item[item.Handle()] = &splitterLayoutItem{stretchFactor: 1}
		}

		stretchTotal += li.StretchFactor(item)
	}

	for i, item := range li.children {
		if i%2 == 1 || !anyVisibleItemInHierarchy(item) {
			continue
		}

		sli := li.hwnd2Item[item.Handle()]
		sli.growth = 0
		sli.keepSize = false
		if sli.oldExplicitSize > 0 {
			sli.size = sli.oldExplicitSize
		} else {
			sli.size = int(float64(li.StretchFactor(item)) / float64(stretchTotal) * float64(regularSpace))
		}

		min := minSizes[i]
		if minSizesTotal <= regularSpace {
			if sli.size < min {
				sli.size = min
			}
		}

		if sli.size >= min {
			flags := item.LayoutFlags()

			if li.orientation == Horizontal && flags&GrowableHorz == 0 || li.orientation == Vertical && flags&GrowableVert == 0 {
				sli.size = min
				sli.keepSize = true
			}
		}
	}
}
