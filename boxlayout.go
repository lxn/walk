// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"math"
	"sort"
	"sync"

	"github.com/lxn/win"
)

type Orientation byte

const (
	NoOrientation Orientation = 0
	Horizontal                = 1 << 0
	Vertical                  = 1 << 1
)

type BoxLayout struct {
	LayoutBase
	orientation        Orientation
	hwnd2StretchFactor map[win.HWND]int
}

func newBoxLayout(orientation Orientation) *BoxLayout {
	l := &BoxLayout{
		LayoutBase: LayoutBase{
			margins96dpi: Margins{9, 9, 9, 9},
			spacing96dpi: 6,
		},
		orientation:        orientation,
		hwnd2StretchFactor: make(map[win.HWND]int),
	}
	l.layout = l

	return l
}

func NewHBoxLayout() *BoxLayout {
	return newBoxLayout(Horizontal)
}

func NewVBoxLayout() *BoxLayout {
	return newBoxLayout(Vertical)
}

func (l *BoxLayout) Orientation() Orientation {
	return l.orientation
}

func (l *BoxLayout) SetOrientation(value Orientation) error {
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

func (l *BoxLayout) StretchFactor(widget Widget) int {
	if factor, ok := l.hwnd2StretchFactor[widget.Handle()]; ok {
		return factor
	}

	return 1
}

func (l *BoxLayout) SetStretchFactor(widget Widget, factor int) error {
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

func (l *BoxLayout) CreateLayoutItem(ctx *LayoutContext) ContainerLayoutItem {
	li := &boxLayoutItem{
		size2MinSize:       make(map[Size]Size),
		orientation:        l.orientation,
		hwnd2StretchFactor: make(map[win.HWND]int),
	}

	for hwnd, sf := range l.hwnd2StretchFactor {
		li.hwnd2StretchFactor[hwnd] = sf
	}

	return li
}

type boxLayoutItemInfo struct {
	item     LayoutItem
	index    int
	prefSize int // in native pixels
	minSize  int // in native pixels
	maxSize  int // in native pixels
	stretch  int
	greedy   bool
}

type boxLayoutItemInfoList []boxLayoutItemInfo

func (l boxLayoutItemInfoList) Len() int {
	return len(l)
}

func (l boxLayoutItemInfoList) Less(i, j int) bool {
	_, iIsSpacer := l[i].item.(*spacerLayoutItem)
	_, jIsSpacer := l[j].item.(*spacerLayoutItem)

	if l[i].greedy == l[j].greedy {
		if iIsSpacer == jIsSpacer {
			minDiff := l[i].minSize - l[j].minSize

			if minDiff == 0 {
				return l[i].maxSize/l[i].stretch < l[j].maxSize/l[j].stretch
			}

			return minDiff > 0
		}

		return jIsSpacer
	}

	return l[i].greedy
}

func (l boxLayoutItemInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type boxLayoutItem struct {
	ContainerLayoutItemBase
	mutex              sync.Mutex
	size2MinSize       map[Size]Size // in native pixels
	orientation        Orientation
	hwnd2StretchFactor map[win.HWND]int
}

func (li *boxLayoutItem) LayoutFlags() LayoutFlags {
	return boxLayoutFlags(li.orientation, li.children)
}

func (li *boxLayoutItem) IdealSize() Size {
	return li.MinSize()
}

func (li *boxLayoutItem) MinSize() Size {
	return li.MinSizeForSize(li.geometry.ClientSize)
}

func (li *boxLayoutItem) HeightForWidth(width int) int {
	return li.MinSizeForSize(Size{width, li.geometry.ClientSize.Height}).Height
}

func (li *boxLayoutItem) MinSizeForSize(size Size) Size {
	li.mutex.Lock()
	defer li.mutex.Unlock()

	if min, ok := li.size2MinSize[size]; ok {
		return min
	}

	bounds := Rectangle{Width: size.Width, Height: size.Height}

	items := boxLayoutItems(li, itemsToLayout(li.children), li.orientation, li.alignment, bounds, li.margins96dpi, li.spacing96dpi, li.hwnd2StretchFactor)

	margins := MarginsFrom96DPI(li.margins96dpi, li.ctx.dpi)
	spacing := IntFrom96DPI(li.spacing96dpi, li.ctx.dpi)
	s := Size{margins.HNear + margins.HFar, margins.VNear + margins.VFar}

	var maxSecondary int
	for _, item := range items {
		min := li.MinSizeEffectiveForChild(item.Item)

		if hfw, ok := item.Item.(HeightForWidther); ok && hfw.HasHeightForWidth() {
			item.Bounds.Height = hfw.HeightForWidth(item.Bounds.Width)
		} else {
			item.Bounds.Height = min.Height
		}
		item.Bounds.Width = min.Width

		if li.orientation == Horizontal {
			maxSecondary = maxi(maxSecondary, item.Bounds.Height)

			s.Width += item.Bounds.Width
		} else {
			maxSecondary = maxi(maxSecondary, item.Bounds.Width)

			s.Height += item.Bounds.Height
		}
	}

	if li.orientation == Horizontal {
		s.Width += (len(items) - 1) * spacing
		s.Height += maxSecondary
	} else {
		s.Height += (len(items) - 1) * spacing
		s.Width += maxSecondary
	}

	if s.Width > 0 && s.Height > 0 {
		li.size2MinSize[size] = s
	}

	return s
}

func (li *boxLayoutItem) PerformLayout() []LayoutResultItem {
	cb := Rectangle{Width: li.geometry.ClientSize.Width, Height: li.geometry.ClientSize.Height}
	return boxLayoutItems(li, itemsToLayout(li.children), li.orientation, li.alignment, cb, li.margins96dpi, li.spacing96dpi, li.hwnd2StretchFactor)
}

func boxLayoutFlags(orientation Orientation, children []LayoutItem) LayoutFlags {
	var flags LayoutFlags
	var hasNonShrinkableHorz bool
	var hasNonShrinkableVert bool

	if len(children) == 0 {
		return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert
	} else {
		for i := 0; i < len(children); i++ {
			item := children[i]

			if _, ok := item.(*splitterHandleLayoutItem); ok || !shouldLayoutItem(item) {
				continue
			}

			if s, ok := item.(*spacerLayoutItem); ok {
				if s.greedyLocallyOnly {
					continue
				}
			}

			f := item.LayoutFlags()
			flags |= f
			if f&ShrinkableHorz == 0 {
				hasNonShrinkableHorz = true
			}
			if f&ShrinkableVert == 0 {
				hasNonShrinkableVert = true
			}
		}
	}

	if orientation == Horizontal {
		flags |= GrowableHorz

		if hasNonShrinkableVert {
			flags &^= ShrinkableVert
		}
	} else {
		flags |= GrowableVert

		if hasNonShrinkableHorz {
			flags &^= ShrinkableHorz
		}
	}

	return flags
}

// boxLayoutItems lays out items. bounds parameter is in native pixels.
func boxLayoutItems(container ContainerLayoutItem, items []LayoutItem, orientation Orientation, alignment Alignment2D, bounds Rectangle, margins96dpi Margins, spacing96dpi int, hwnd2StretchFactor map[win.HWND]int) []LayoutResultItem {
	if len(items) == 0 {
		return nil
	}

	dpi := container.Context().dpi
	margins := MarginsFrom96DPI(margins96dpi, dpi)
	spacing := IntFrom96DPI(spacing96dpi, dpi)

	var greedyNonSpacerCount int
	var greedySpacerCount int
	var stretchFactorsTotal [3]int
	stretchFactors := make([]int, len(items))
	var minSizesRemaining int
	minSizes := make([]int, len(items))
	maxSizes := make([]int, len(items))
	sizes := make([]int, len(items))
	prefSizes2 := make([]int, len(items))
	var shrinkableAmount1Total int
	shrinkableAmount1 := make([]int, len(items))
	shrinkable2 := make([]bool, len(items))
	growable2 := make([]bool, len(items))
	sortedItemInfo := boxLayoutItemInfoList(make([]boxLayoutItemInfo, len(items)))

	for i, item := range items {
		sf := hwnd2StretchFactor[item.Handle()]
		if sf == 0 {
			sf = 1
		}
		stretchFactors[i] = sf

		geometry := item.Geometry()

		flags := item.LayoutFlags()

		max := geometry.MaxSize
		var pref Size
		if hfw, ok := item.(HeightForWidther); !ok || !hfw.HasHeightForWidth() {
			if is, ok := item.(IdealSizer); ok {
				pref = is.IdealSize()
			}
		}

		if orientation == Horizontal {
			growable2[i] = flags&GrowableVert > 0

			minSizes[i] = container.MinSizeEffectiveForChild(item).Width

			if max.Width > 0 {
				maxSizes[i] = max.Width
			} else if pref.Width > 0 && flags&GrowableHorz == 0 {
				maxSizes[i] = pref.Width
			} else {
				maxSizes[i] = 32768
			}

			prefSizes2[i] = pref.Height

			sortedItemInfo[i].prefSize = pref.Width
			sortedItemInfo[i].greedy = flags&GreedyHorz > 0
		} else {
			growable2[i] = flags&GrowableHorz > 0

			if hfw, ok := item.(HeightForWidther); ok && hfw.HasHeightForWidth() {
				minSizes[i] = hfw.HeightForWidth(bounds.Width - margins.HNear - margins.HFar)
			} else {
				minSizes[i] = container.MinSizeEffectiveForChild(item).Height
			}

			if max.Height > 0 {
				maxSizes[i] = max.Height
			} else if hfw, ok := item.(HeightForWidther); ok && flags&GrowableVert == 0 && hfw.HasHeightForWidth() {
				maxSizes[i] = minSizes[i]
			} else if pref.Height > 0 && flags&GrowableVert == 0 {
				maxSizes[i] = pref.Height
			} else {
				maxSizes[i] = 32768
			}

			prefSizes2[i] = pref.Width

			sortedItemInfo[i].prefSize = pref.Height
			sortedItemInfo[i].greedy = flags&GreedyVert > 0
		}

		sortedItemInfo[i].index = i
		sortedItemInfo[i].minSize = minSizes[i]
		sortedItemInfo[i].maxSize = maxSizes[i]
		sortedItemInfo[i].stretch = sf
		sortedItemInfo[i].item = item

		if orientation == Horizontal && flags&(ShrinkableHorz|GrowableHorz|GreedyHorz) == ShrinkableHorz ||
			orientation == Vertical && flags&(ShrinkableVert|GrowableVert|GreedyVert) == ShrinkableVert {
			if amount := sortedItemInfo[i].prefSize - minSizes[i]; amount > 0 {
				shrinkableAmount1[i] = amount
				shrinkableAmount1Total += amount
			}
		}
		shrinkable2[i] = orientation == Horizontal && flags&ShrinkableVert != 0 || orientation == Vertical && flags&ShrinkableHorz != 0

		if shrinkableAmount1[i] > 0 {
			minSizesRemaining += sortedItemInfo[i].prefSize
		} else {
			minSizesRemaining += minSizes[i]
		}

		if sortedItemInfo[i].greedy {
			if _, isSpacer := item.(*spacerLayoutItem); !isSpacer {
				greedyNonSpacerCount++
				stretchFactorsTotal[0] += sf
			} else {
				greedySpacerCount++
				stretchFactorsTotal[1] += sf
			}
		} else {
			stretchFactorsTotal[2] += sf
		}
	}

	sort.Stable(sortedItemInfo)

	var start1, start2, space1, space2 int
	if orientation == Horizontal {
		start1 = bounds.X + margins.HNear
		start2 = bounds.Y + margins.VNear
		space1 = bounds.Width - margins.HNear - margins.HFar
		space2 = bounds.Height - margins.VNear - margins.VFar
	} else {
		start1 = bounds.Y + margins.VNear
		start2 = bounds.X + margins.HNear
		space1 = bounds.Height - margins.VNear - margins.VFar
		space2 = bounds.Width - margins.HNear - margins.HFar
	}

	spacingRemaining := spacing * (len(items) - 1)
	excess := float64(space1 - minSizesRemaining - spacingRemaining)

	offsets := [3]int{0, greedyNonSpacerCount, greedyNonSpacerCount + greedySpacerCount}
	counts := [3]int{greedyNonSpacerCount, greedySpacerCount, len(items) - greedyNonSpacerCount - greedySpacerCount}

	for i := 0; i < 3; i++ {
		stretchFactorsRemaining := stretchFactorsTotal[i]

		for j := 0; j < counts[i]; j++ {
			info := sortedItemInfo[offsets[i]+j]
			k := info.index

			stretch := stretchFactors[k]
			min := info.minSize
			max := info.maxSize
			var size int
			var corrected bool
			if shrinkableAmount1[k] > 0 {
				size = info.prefSize
				if excess < 0.0 {
					size -= mini(shrinkableAmount1[k], int(math.Round(-excess/float64(shrinkableAmount1Total)*float64(shrinkableAmount1[k]))))
					corrected = true
				}
			} else {
				size = min
			}

			if !corrected && min < max {
				excessSpace := float64(space1 - minSizesRemaining - spacingRemaining)
				size += int(math.Round(excessSpace * float64(stretch) / float64(stretchFactorsRemaining)))
				if size < min {
					size = min
				} else if size > max {
					size = max
				}
			}

			sizes[k] = size

			if shrinkableAmount1[k] > 0 {
				minSizesRemaining -= info.prefSize
			} else {
				minSizesRemaining -= min
			}
			stretchFactorsRemaining -= stretch
			space1 -= (size + spacing)
			spacingRemaining -= spacing
		}
	}

	results := make([]LayoutResultItem, 0, len(items))

	excessTotal := space1 - minSizesRemaining - spacingRemaining
	excessShare := excessTotal / len(items)
	halfExcessShare := excessTotal / (len(items) * 2)
	p1 := start1
	for i, item := range items {
		s1 := sizes[i]

		var s2 int
		if hfw, ok := item.(HeightForWidther); ok && orientation == Horizontal && hfw.HasHeightForWidth() {
			s2 = hfw.HeightForWidth(s1)
		} else if shrinkable2[i] || growable2[i] {
			s2 = space2
		} else {
			s2 = prefSizes2[i]
		}

		align := item.Geometry().Alignment
		if align == AlignHVDefault {
			align = alignment
		}

		var x, y, w, h, p2 int
		if orientation == Horizontal {
			switch align {
			case AlignHNearVNear, AlignHNearVCenter, AlignHNearVFar:
				// nop

			case AlignHFarVNear, AlignHFarVCenter, AlignHFarVFar:
				p1 += excessShare

			default:
				p1 += halfExcessShare
			}

			switch align {
			case AlignHNearVNear, AlignHCenterVNear, AlignHFarVNear:
				p2 = start2

			case AlignHNearVFar, AlignHCenterVFar, AlignHFarVFar:
				p2 = start2 + space2 - s2

			default:
				p2 = start2 + (space2-s2)/2
			}

			x, y, w, h = p1, p2, s1, s2
		} else {
			switch align {
			case AlignHNearVNear, AlignHCenterVNear, AlignHFarVNear:
				// nop

			case AlignHNearVFar, AlignHCenterVFar, AlignHFarVFar:
				p1 += excessShare

			default:
				p1 += halfExcessShare
			}

			switch align {
			case AlignHNearVNear, AlignHNearVCenter, AlignHNearVFar:
				p2 = start2

			case AlignHFarVNear, AlignHFarVCenter, AlignHFarVFar:
				p2 = start2 + space2 - s2

			default:
				p2 = start2 + (space2-s2)/2
			}

			x, y, w, h = p2, p1, s2, s1
		}

		if orientation == Horizontal {
			switch align {
			case AlignHNearVNear, AlignHNearVCenter, AlignHNearVFar:
				p1 += excessShare

			case AlignHFarVNear, AlignHFarVCenter, AlignHFarVFar:
				// nop

			default:
				p1 += halfExcessShare
			}

		} else {
			switch align {
			case AlignHNearVNear, AlignHCenterVNear, AlignHFarVNear:
				p1 += excessShare

			case AlignHNearVFar, AlignHCenterVFar, AlignHFarVFar:
				// nop

			default:
				p1 += halfExcessShare
			}
		}

		p1 += s1 + spacing

		results = append(results, LayoutResultItem{Item: item, Bounds: Rectangle{X: x, Y: y, Width: w, Height: h}})
	}

	return results
}
