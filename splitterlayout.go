// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type splitterLayout struct {
	container   Container
	orientation Orientation
	fractions   []float64
	resetNeeded bool
}

func newSplitterLayout(orientation Orientation) *splitterLayout {
	return &splitterLayout{
		orientation: orientation,
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
	return Margins{}
}

func (l *splitterLayout) SetMargins(value Margins) error {
	return newError("not supported")
}

func (l *splitterLayout) Spacing() int {
	return 0
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

func (l *splitterLayout) Fractions() []float64 {
	return l.fractions
}

func (l *splitterLayout) SetFractions(fractions []float64) error {
	l.fractions = fractions

	return l.Update(false)
}

func (l *splitterLayout) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (l *splitterLayout) MinSize() Size {
	var s Size

	for _, widget := range l.container.Children().items {
		cur := widget.BaseWidget().minSizeEffective()

		if l.orientation == Horizontal {
			s.Width += cur.Width
			s.Height = maxi(s.Height, cur.Height)
		} else {
			s.Height += cur.Height
			s.Width = maxi(s.Width, cur.Width)
		}
	}

	return s
}

func (l *splitterLayout) spaceForRegularWidgets() int {
	splitter := l.container.(*Splitter)
	cb := splitter.ClientBounds().Size()

	var space int
	if l.orientation == Horizontal {
		space = cb.Width
	} else {
		space = cb.Height
	}

	return space - (splitter.Children().Len()/2)*splitter.handleWidth
}

func (l *splitterLayout) reset() {
	children := l.container.Children()
	regularCount := children.Len()/2 + children.Len()%2

	if cap(l.fractions) < regularCount {
		temp := make([]float64, regularCount)
		copy(temp, l.fractions)
		l.fractions = temp
	}

	l.fractions = l.fractions[:regularCount]

	if regularCount == 0 {
		return
	}

	fraction := 1 / float64(regularCount)

	for i := 0; i < regularCount; i++ {
		l.fractions[i] = fraction
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
		l.resetNeeded = false

		l.reset()
	}

	widgets := l.container.Children().items
	splitter := l.container.(*Splitter)
	handleWidth := splitter.HandleWidth()
	sizes := make([]int, len(widgets))
	cb := splitter.ClientBounds()
	space1 := l.spaceForRegularWidgets()

	var space2 int
	if l.orientation == Horizontal {
		space2 = cb.Height
	} else {
		space2 = cb.Width
	}

	for i := range widgets {
		j := i/2 + i%2

		if i%2 == 0 {
			sizes[i] = int(float64(space1) * l.fractions[j])
		} else {
			sizes[i] = handleWidth
		}
	}

	hdwp := BeginDeferWindowPos(int32(len(widgets)))
	if hdwp == 0 {
		return lastError("BeginDeferWindowPos")
	}

	p1 := 0
	for i, widget := range widgets {
		var s1 int
		if i == len(widgets)-1 {
			s1 = space1 - p1
		} else {
			s1 = sizes[i]
		}

		var x, y, w, h int
		if l.orientation == Horizontal {
			x, y, w, h = p1, 0, s1, space2
		} else {
			x, y, w, h = 0, p1, space2, s1
		}

		if hdwp = DeferWindowPos(
			hdwp,
			widget.BaseWidget().hWnd,
			0,
			int32(x),
			int32(y),
			int32(w),
			int32(h),
			SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_NOZORDER); hdwp == 0 {

			return lastError("DeferWindowPos")
		}

		p1 += s1
	}

	if !EndDeferWindowPos(hdwp) {
		return lastError("EndDeferWindowPos")
	}

	return nil
}
