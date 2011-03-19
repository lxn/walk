// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

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

func (l *splitterLayout) SetMargins(value Margins) os.Error {
	return newError("not supported")
}

func (l *splitterLayout) Spacing() int {
	return 0
}

func (l *splitterLayout) SetSpacing(value int) os.Error {
	return newError("not supported")
}

func (l *splitterLayout) Orientation() Orientation {
	return l.orientation
}

func (l *splitterLayout) SetOrientation(value Orientation) os.Error {
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

func (l *splitterLayout) SetFractions(fractions []float64) os.Error {
	l.fractions = fractions

	return l.Update(false)
}

func (l *splitterLayout) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (l *splitterLayout) MinSize() Size {
	return Size{10, 10}
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
	space := l.spaceForRegularWidgets()

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

	splitter := l.container.(*Splitter)
	handleWidth := splitter.HandleWidth()
	handlesFraction := float64((regularCount-1)*handleWidth) / float64(space+(regularCount-1)*handleWidth)
	regularWidgetsFraction := (1 - handlesFraction) / float64(regularCount)

	for i := 0; i < regularCount; i++ {
		l.fractions[i] = regularWidgetsFraction
	}
}

func (l *splitterLayout) Update(reset bool) os.Error {
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
	cb := l.container.ClientBounds()

	var space1 int
	var space2 int
	if l.orientation == Horizontal {
		space1 = cb.Width
		space2 = cb.Height
	} else {
		space1 = cb.Height
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

	p1 := 0
	for i, widget := range widgets {
		var s1 int
		if i == len(widgets)-1 {
			s1 = space1 - p1
		} else {
			s1 = sizes[i]
		}

		if l.orientation == Horizontal {
			widget.SetBounds(Rectangle{p1, 0, s1, space2})
		} else {
			widget.SetBounds(Rectangle{0, p1, space2, s1})
		}

		p1 += s1
	}

	return nil
}
