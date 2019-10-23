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
	container    Container
	orientation  Orientation
	margins96dpi Margins
	hwnd2Item    map[win.HWND]*splitterLayoutItem
	resetNeeded  bool
	suspended    bool
}

type splitterLayoutItem struct {
	size                 int // in native pixels
	oldExplicitSize      int // in native pixels
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
	return l.margins96dpi
}

func (l *splitterLayout) SetMargins(value Margins) error {
	l.margins96dpi = value

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

// spaceUnavailableToRegularWidgets returns amount of space unavailable to regular widgets in native pixels.
func (l *splitterLayout) spaceUnavailableToRegularWidgets() int {
	splitter := l.container.(*Splitter)

	var space int

	for _, widget := range l.container.Children().items {
		if _, isHandle := widget.window.(*splitterHandle); isHandle && widget.visible {
			space += splitter.handleWidth
		}
	}

	return IntFrom96DPI(space, splitter.DPI())
}

func (l *splitterLayout) CreateLayoutItem(ctx *LayoutContext) ContainerLayoutItem {
	splitter := l.container.(*Splitter)

	hwnd2Item := make(map[win.HWND]*splitterLayoutItem, len(l.hwnd2Item))
	for hwnd, sli := range l.hwnd2Item {
		hwnd2Item[hwnd] = sli
	}

	li := &splitterContainerLayoutItem{
		orientation:                    l.orientation,
		hwnd2Item:                      hwnd2Item,
		spaceUnavailableToRegularItems: l.spaceUnavailableToRegularWidgets(),
		handleWidth96dpi:               splitter.HandleWidth(),
		anyNonFixed:                    l.anyNonFixed(),
		resetNeeded:                    l.resetNeeded,
	}

	li.margins96dpi = l.margins96dpi

	return li
}

type splitterContainerLayoutItem struct {
	ContainerLayoutItemBase
	orientation                    Orientation
	hwnd2Item                      map[win.HWND]*splitterLayoutItem
	spaceUnavailableToRegularItems int // in native pixels
	handleWidth96dpi               int
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
	return li.MinSizeForSize(li.geometry.ClientSize)
}

func (li *splitterContainerLayoutItem) HeightForWidth(width int) int {
	return li.MinSizeForSize(Size{width, li.geometry.ClientSize.Height}).Height
}

func (li *splitterContainerLayoutItem) MinSizeForSize(size Size) Size {
	marginsPixels := MarginsFrom96DPI(li.margins96dpi, li.ctx.dpi)
	margins := Size{marginsPixels.HNear + marginsPixels.HFar, marginsPixels.VNear + marginsPixels.VFar}
	s := margins

	for _, item := range li.children {
		if !anyVisibleItemInHierarchy(item) {
			continue
		}

		var cur Size

		if sli, ok := li.hwnd2Item[item.Handle()]; ok && li.anyNonFixed && sli.fixed {
			cur = item.Geometry().Size

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

	margins := MarginsFrom96DPI(li.margins96dpi, li.ctx.dpi)
	handleWidthPixels := IntFrom96DPI(li.handleWidth96dpi, li.ctx.dpi)
	sizes := make([]int, len(li.children))
	cb := Rectangle{Width: li.geometry.ClientSize.Width, Height: li.geometry.ClientSize.Height}
	cb.X += margins.HNear
	cb.Y += margins.HFar
	cb.Width -= margins.HNear + margins.HFar
	cb.Height -= margins.VNear + margins.VFar

	var space1, space2 int
	if li.orientation == Horizontal {
		space1 = cb.Width - li.spaceUnavailableToRegularItems
		space2 = cb.Height
	} else {
		space1 = cb.Height - li.spaceUnavailableToRegularItems
		space2 = cb.Width
	}

	type WidgetItem struct {
		item       *splitterLayoutItem
		index      int
		min        int // in native pixels
		max        int // in native pixels
		shrinkable bool
		growable   bool
	}

	var wis []WidgetItem

	anyNonFixed := li.anyNonFixed
	var totalRegularSize int
	for i, item := range li.children {
		if !anyVisibleItemInHierarchy(item) {
			continue
		}

		if i%2 == 0 {
			slItem := li.hwnd2Item[item.Handle()]

			var wi *WidgetItem

			if !anyNonFixed || !slItem.fixed {
				var min, max int

				minSize := li.MinSizeEffectiveForChild(item)
				maxSize := item.Geometry().MaxSize

				if li.orientation == Horizontal {
					min = minSize.Width
					max = maxSize.Width
				} else {
					min = minSize.Height
					max = maxSize.Height
				}

				wis = append(wis, WidgetItem{item: slItem, index: i, min: min, max: max})

				wi = &wis[len(wis)-1]
			}

			size := slItem.size
			var idealSize Size
			if hfw, ok := item.(HeightForWidther); ok && li.orientation == Vertical && hfw.HasHeightForWidth() {
				idealSize.Height = hfw.HeightForWidth(space2)
			} else {
				switch sizer := item.(type) {
				case IdealSizer:
					idealSize = sizer.IdealSize()

				case MinSizer:
					idealSize = sizer.MinSize()
				}
			}

			if flags := item.LayoutFlags(); li.orientation == Horizontal {
				if flags&ShrinkableHorz == 0 {
					size = maxi(size, idealSize.Width)
					if wi != nil {
						wi.min = maxi(wi.min, size)
					}
				} else if wi != nil {
					wi.shrinkable = true
				}
				if flags&GrowableHorz == 0 {
					size = mini(size, idealSize.Width)
					if wi != nil {
						wi.max = mini(wi.max, size)
					}
				} else if wi != nil {
					wi.growable = true
				}
			} else {
				if flags&ShrinkableVert == 0 {
					size = maxi(size, idealSize.Height)
					if wi != nil {
						wi.min = maxi(wi.min, size)
					}
				} else if wi != nil {
					wi.shrinkable = true
				}
				if flags&GrowableVert == 0 {
					size = mini(size, idealSize.Height)
					if wi != nil {
						wi.max = mini(wi.max, size)
					}
				} else if wi != nil {
					wi.growable = true
				}
			}

			totalRegularSize += size
			sizes[i] = size
		} else {
			sizes[i] = handleWidthPixels
		}
	}

	var resultItems []LayoutResultItem

	diff := space1 - totalRegularSize

	if diff != 0 && len(sizes) > 1 {
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
				if !wItem.item.keepSize && (diff < 0 && wItem.item.size > wItem.min || diff > 0 && (wItem.item.size < wItem.max || wItem.max == 0)) {
					wi = &wItem
					break
				}
			}
			if wi == nil {
				break
			}

			if diff > 0 {
				sizes[wi.index]++
				wi.item.size++
				wi.item.growth++
				diff--
			} else {
				sizes[wi.index]--
				wi.item.size--
				wi.item.growth--
				diff++
			}
		}
	}

	var p1 int
	if li.orientation == Horizontal {
		p1 = margins.HNear
	} else {
		p1 = margins.VNear
	}
	for i, item := range li.children {
		if !anyVisibleItemInHierarchy(item) {
			continue
		}

		s1 := sizes[i]

		var x, y, w, h int
		if li.orientation == Horizontal {
			x, y, w, h = p1, margins.VNear, s1, space2
		} else {
			x, y, w, h = margins.HNear, p1, space2, s1
		}

		resultItems = append(resultItems, LayoutResultItem{Item: item, Bounds: Rectangle{x, y, w, h}})

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
		regularSpace = li.Geometry().ClientSize.Width - li.spaceUnavailableToRegularItems
	} else {
		regularSpace = li.Geometry().ClientSize.Height - li.spaceUnavailableToRegularItems
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
