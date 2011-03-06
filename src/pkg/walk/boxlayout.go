// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"sort"
)

type Orientation byte

const (
	Horizontal Orientation = iota
	Vertical
)

type BoxLayout struct {
	container            Container
	margins              Margins
	spacing              int
	orientation          Orientation
	widget2StretchFactor map[*WidgetBase]int
	resetNeeded          bool
}

func newBoxLayout(orientation Orientation) *BoxLayout {
	return &BoxLayout{
		orientation:          orientation,
		widget2StretchFactor: make(map[*WidgetBase]int),
	}
}

func NewHBoxLayout() *BoxLayout {
	return newBoxLayout(Horizontal)
}

func NewVBoxLayout() *BoxLayout {
	return newBoxLayout(Vertical)
}

func (l *BoxLayout) Container() Container {
	return l.container
}

func (l *BoxLayout) SetContainer(value Container) {
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

func (l *BoxLayout) Margins() Margins {
	return l.margins
}

func (l *BoxLayout) SetMargins(value Margins) os.Error {
	if value.HNear < 0 || value.VNear < 0 || value.HFar < 0 || value.VFar < 0 {
		return newError("margins must be positive")
	}

	l.margins = value

	return nil
}

func (l *BoxLayout) Orientation() Orientation {
	return l.orientation
}

func (l *BoxLayout) SetOrientation(value Orientation) os.Error {
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

func (l *BoxLayout) Spacing() int {
	return l.spacing
}

func (l *BoxLayout) SetSpacing(value int) os.Error {
	if value != l.spacing {
		if value < 0 {
			return newError("spacing cannot be negative")
		}

		l.spacing = value

		l.Update(false)
	}

	return nil
}

func (l *BoxLayout) StretchFactor(widget Widget) int {
	if factor, ok := l.widget2StretchFactor[widget.BaseWidget()]; ok {
		return factor
	}

	return 1
}

func (l *BoxLayout) SetStretchFactor(widget Widget, factor int) os.Error {
	if factor != l.StretchFactor(widget) {
		if l.container == nil {
			return newError("container required")
		}
		if !l.container.Children().containsHandle(widget.BaseWidget().hWnd) {
			return newError("unknown widget")
		}
		if factor < 1 {
			return newError("factor must be >= 1")
		}

		l.widget2StretchFactor[widget.BaseWidget()] = factor

		l.Update(false)
	}

	return nil
}

func (l *BoxLayout) cleanupStretchFactors() {
	widgets := l.container.Children()

	for widget, _ := range l.widget2StretchFactor {
		if !widgets.containsHandle(widget.BaseWidget().hWnd) {
			l.widget2StretchFactor[widget.BaseWidget()] = 0, false
		}
	}
}

type widgetInfo struct {
	index   int
	minSize int
	maxSize int
	stretch int
	widget  Widget
}

type widgetInfoList []widgetInfo

func (l widgetInfoList) Len() int {
	return len(l)
}

func (l widgetInfoList) Less(i, j int) bool {
	_, iIsSpacer := l[i].widget.(*Spacer)
	_, jIsSpacer := l[j].widget.(*Spacer)

	if iIsSpacer == jIsSpacer {
		minDiff := l[i].minSize - l[j].minSize

		if minDiff == 0 {
			return l[i].maxSize/l[i].stretch < l[j].maxSize/l[j].stretch
		}

		return minDiff > 0
	}

	return jIsSpacer
}

func (l widgetInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l *BoxLayout) MinSize() Size {
	if l.container == nil {
		return Size{}
	}

	// Begin by finding out which widgets we care about.
	children := l.container.Children()
	widgets := make([]Widget, 0, children.Len())

	for i := 0; i < cap(widgets); i++ {
		widget := children.At(i)

		ps := widget.PreferredSize()
		if ps.Width == 0 && ps.Height == 0 && widget.LayoutFlags() == 0 {
			continue
		}

		widgets = append(widgets, widget)
	}

	// Prepare some useful data.
	sizes := make([]int, len(widgets))
	var s2 int

	for i, widget := range widgets {
		max := widget.MaxSize()
		pref := widget.PreferredSize()

		if l.orientation == Horizontal {
			if max.Width > 0 {
				sizes[i] = mini(pref.Width, max.Width)
			} else {
				sizes[i] = pref.Width
			}

			if pref.Height > s2 {
				s2 = pref.Height
			}
		} else {
			if max.Height > 0 {
				sizes[i] = mini(pref.Height, max.Height)
			} else {
				sizes[i] = pref.Height
			}

			if pref.Width > s2 {
				s2 = pref.Width
			}
		}
	}

	s1 := l.spacing * (len(widgets) - 1)

	if l.orientation == Horizontal {
		s1 += l.margins.HNear + l.margins.HFar
		s2 += l.margins.VNear + l.margins.VFar
	} else {
		s1 += l.margins.VNear + l.margins.VFar
		s2 += l.margins.HNear + l.margins.HFar
	}

	for _, s := range sizes {
		s1 += s
	}

	if l.orientation == Horizontal {
		return Size{s1, s2}
	}

	return Size{s2, s1}
}

func (l *BoxLayout) Update(reset bool) os.Error {
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
		l.resetNeeded = false

		// Make GC happy.
		l.cleanupStretchFactors()
	}

	// Begin by finding out which widgets we care about.
	children := l.container.Children()
	widgets := make([]Widget, 0, children.Len())

	for i := 0; i < cap(widgets); i++ {
		widget := children.At(i)

		ps := widget.PreferredSize()
		if ps.Width == 0 && ps.Height == 0 && widget.LayoutFlags() == 0 {
			continue
		}

		widgets = append(widgets, widget)
	}

	// Prepare some useful data.
	var nonSpacerCount int
	var stretchFactorsTotal [2]int
	stretchFactors := make([]int, len(widgets))
	var minSizesRemaining int
	minSizes := make([]int, len(widgets))
	maxSizes := make([]int, len(widgets))
	sizes := make([]int, len(widgets))
	prefSizes2 := make([]int, len(widgets))
	canGrow2 := make([]bool, len(widgets))
	sortedWidgetInfo := widgetInfoList(make([]widgetInfo, len(widgets)))

	for i, widget := range widgets {
		sf := l.widget2StretchFactor[widget.BaseWidget()]
		if sf == 0 {
			sf = 1
		}
		stretchFactors[i] = sf

		flags := widget.LayoutFlags()

		min := widget.MinSize()
		max := widget.MaxSize()
		pref := widget.PreferredSize()

		if l.orientation == Horizontal {
			canGrow2[i] = flags&VGrow > 0

			if min.Width > 0 {
				minSizes[i] = min.Width
			} else if pref.Width > 0 && flags&HShrink == 0 {
				minSizes[i] = pref.Width
			}

			if max.Width > 0 {
				maxSizes[i] = max.Width
			} else if pref.Width > 0 && flags&HGrow == 0 {
				maxSizes[i] = pref.Width
			} else {
				maxSizes[i] = 32768
			}

			prefSizes2[i] = pref.Height
		} else {
			canGrow2[i] = flags&HGrow > 0

			if min.Height > 0 {
				minSizes[i] = min.Height
			} else if pref.Height > 0 && flags&VShrink == 0 {
				minSizes[i] = pref.Height
			}

			if max.Height > 0 {
				maxSizes[i] = max.Height
			} else if pref.Height > 0 && flags&VGrow == 0 {
				maxSizes[i] = pref.Height
			} else {
				maxSizes[i] = 32768
			}

			prefSizes2[i] = pref.Width
		}

		sortedWidgetInfo[i].index = i
		sortedWidgetInfo[i].minSize = minSizes[i]
		sortedWidgetInfo[i].maxSize = maxSizes[i]
		sortedWidgetInfo[i].stretch = sf
		sortedWidgetInfo[i].widget = widget

		minSizesRemaining += minSizes[i]

		if _, isSpacer := widget.(*Spacer); !isSpacer {
			nonSpacerCount++
			stretchFactorsTotal[0] += sf
		} else {
			stretchFactorsTotal[1] += sf
		}
	}

	sort.Sort(sortedWidgetInfo)

	cb := l.container.ClientBounds()
	var start1, start2, space1, space2 int
	if l.orientation == Horizontal {
		start1 = cb.X + l.margins.HNear
		start2 = cb.Y + l.margins.VNear
		space1 = cb.Width - l.margins.HNear - l.margins.HFar
		space2 = cb.Height - l.margins.VNear - l.margins.VFar
	} else {
		start1 = cb.Y + l.margins.VNear
		start2 = cb.X + l.margins.HNear
		space1 = cb.Height - l.margins.VNear - l.margins.VFar
		space2 = cb.Width - l.margins.HNear - l.margins.HFar
	}

	// Now calculate widget primary axis sizes.
	spacingRemaining := l.spacing * (len(widgets) - 1)

	offsets := [2]int{0, nonSpacerCount}
	counts := [2]int{nonSpacerCount, len(widgets) - nonSpacerCount}

	for i := 0; i < 2; i++ {
		stretchFactorsRemaining := stretchFactorsTotal[i]

		for j := 0; j < counts[i]; j++ {
			info := sortedWidgetInfo[offsets[i]+j]
			k := info.index

			stretch := stretchFactors[k]
			min := info.minSize
			max := info.maxSize
			size := min

			if min < max {
				excessSpace := float64(space1 - minSizesRemaining - spacingRemaining)
				size += int(excessSpace * float64(stretch) / float64(stretchFactorsRemaining))
				if size < min {
					size = min
				} else if size > max {
					size = max
				}
			}

			sizes[k] = size

			minSizesRemaining -= min
			stretchFactorsRemaining -= stretch
			space1 -= (size + l.spacing)
			spacingRemaining -= l.spacing
		}
	}

	// Finally position widgets.
	excessTotal := space1 - minSizesRemaining - spacingRemaining
	excessShare := excessTotal / (len(widgets) + 1)
	p1 := start1
	for i, widget := range widgets {
		p1 += excessShare
		s1 := sizes[i]

		var s2 int
		if canGrow2[i] {
			s2 = space2
		} else {
			s2 = prefSizes2[i]
		}

		p2 := start2 + (space2-s2)/2

		if l.orientation == Horizontal {
			widget.SetBounds(Rectangle{p1, p2, s1, s2})
		} else {
			widget.SetBounds(Rectangle{p2, p1, s2, s1})
		}

		p1 += s1 + l.spacing
	}

	return nil
}
