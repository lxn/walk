// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"sort"
)

import . "github.com/lxn/go-winapi"

type gridLayoutCell struct {
	row    int
	column int
	widget Widget
}

type gridLayoutSection struct {
	greedyNonSpacerCount int
	greedySpacerCount    int
}

type GridLayout struct {
	container            Container
	margins              Margins
	spacing              int
	resetNeeded          bool
	rowStretchFactors    []int
	columnStretchFactors []int
	widget2Cell          map[Widget]*gridLayoutCell
	cells                [][]gridLayoutCell
}

func NewGridLayout() *GridLayout {
	l := &GridLayout{
		widget2Cell: make(map[Widget]*gridLayoutCell),
	}

	return l
}

func (l *GridLayout) Container() Container {
	return l.container
}

func (l *GridLayout) SetContainer(value Container) {
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

func (l *GridLayout) Margins() Margins {
	return l.margins
}

func (l *GridLayout) SetMargins(value Margins) os.Error {
	if value.HNear < 0 || value.VNear < 0 || value.HFar < 0 || value.VFar < 0 {
		return newError("margins must be positive")
	}

	l.margins = value

	return nil
}

func (l *GridLayout) Spacing() int {
	return l.spacing
}

func (l *GridLayout) SetSpacing(value int) os.Error {
	if value != l.spacing {
		if value < 0 {
			return newError("spacing cannot be negative")
		}

		l.spacing = value

		l.Update(false)
	}

	return nil
}

func (l *GridLayout) sufficientStretchFactors(stretchFactors []int, required int) []int {
	oldLen := len(stretchFactors)
	if oldLen < required {
		if cap(stretchFactors) < required {
			temp := make([]int, required, maxi(required, len(stretchFactors)*2))
			copy(temp, stretchFactors)
			stretchFactors = temp
		} else {
			stretchFactors = stretchFactors[:required]
		}

		for i := oldLen; i < len(stretchFactors); i++ {
			stretchFactors[i] = 1
		}
	}

	return stretchFactors
}

func (l *GridLayout) ensureSufficientSize(rows, columns int) {
	l.rowStretchFactors = l.sufficientStretchFactors(l.rowStretchFactors, rows)
	l.columnStretchFactors = l.sufficientStretchFactors(l.columnStretchFactors, columns)

	if len(l.cells) < len(l.rowStretchFactors) {
		if cap(l.cells) < cap(l.rowStretchFactors) {
			temp := make([][]gridLayoutCell, len(l.rowStretchFactors), cap(l.rowStretchFactors))
			copy(temp, l.cells)
			l.cells = temp
		} else {
			l.cells = l.cells[:len(l.rowStretchFactors)]
		}
	}

	for i := 0; i < len(l.cells); i++ {
		if len(l.cells[i]) < len(l.columnStretchFactors) {
			if cap(l.cells[i]) < cap(l.columnStretchFactors) {
				temp := make([]gridLayoutCell, len(l.columnStretchFactors))
				copy(temp, l.cells[i])
				l.cells[i] = temp
			} else {
				l.cells[i] = l.cells[i][:len(l.columnStretchFactors)]
			}
		}
	}

	// FIXME: Not sure if this works.
	for widget, cell := range l.widget2Cell {
		l.widget2Cell[widget] = &l.cells[cell.row][cell.column]
	}
}

func (l *GridLayout) RowStretchFactor(row int) int {
	if row < 0 {
		// FIXME: Should we rather return an error?
		return -1
	}

	if row >= len(l.rowStretchFactors) {
		return 1
	}

	return l.rowStretchFactors[row]
}

func (l *GridLayout) SetRowStretchFactor(row, factor int) os.Error {
	if row < 0 {
		return newError("row must be >= 0")
	}

	if factor != l.RowStretchFactor(row) {
		if l.container == nil {
			return newError("container required")
		}
		if factor < 1 {
			return newError("factor must be >= 1")
		}

		l.ensureSufficientSize(row+1, len(l.columnStretchFactors))

		l.rowStretchFactors[row] = factor

		l.Update(false)
	}

	return nil
}

func (l *GridLayout) ColumnStretchFactor(column int) int {
	if column < 0 {
		// FIXME: Should we rather return an error?
		return -1
	}

	if column >= len(l.columnStretchFactors) {
		return 1
	}

	return l.columnStretchFactors[column]
}

func (l *GridLayout) SetColumnStretchFactor(column, factor int) os.Error {
	if column < 0 {
		return newError("column must be >= 0")
	}

	if factor != l.ColumnStretchFactor(column) {
		if l.container == nil {
			return newError("container required")
		}
		if factor < 1 {
			return newError("factor must be >= 1")
		}

		l.ensureSufficientSize(len(l.rowStretchFactors), column+1)

		l.columnStretchFactors[column] = factor

		l.Update(false)
	}

	return nil
}

func (l *GridLayout) Location(widget Widget) (row, column int, ok bool) {
	if widget == nil {
		return 0, 0, false
	}

	cell := l.widget2Cell[widget]

	if cell == nil ||
		l.container == nil ||
		!l.container.Children().containsHandle(widget.BaseWidget().hWnd) {
		return 0, 0, false
	}

	return cell.row, cell.column, true
}

func (l *GridLayout) SetLocation(widget Widget, row, column int) os.Error {
	if widget == nil {
		return newError("widget required")
	}
	if l.container == nil {
		return newError("container required")
	}
	if !l.container.Children().containsHandle(widget.BaseWidget().hWnd) {
		return newError("widget must be child of container")
	}
	if row < 0 || column < 0 {
		return newError("location must be positive")
	}

	l.ensureSufficientSize(row+1, column+1)

	cell := &l.cells[row][column]
	cell.row = row
	cell.column = column
	cell.widget = widget

	l.widget2Cell[widget] = cell

	return nil
}

func (l *GridLayout) cleanup() {
	// Make sure only children of our container occupy the precious cells.
	children := l.container.Children()
	for widget, cell := range l.widget2Cell {
		if !children.containsHandle(widget.BaseWidget().hWnd) {
			l.widget2Cell[widget] = nil, false
			cell.widget = nil
		}
	}
}

func (l *GridLayout) stretchFactorsTotal(stretchFactors []int) int {
	total := 0

	for _, v := range stretchFactors {
		total += maxi(1, v)
	}

	return total
}

func (l *GridLayout) LayoutFlags() LayoutFlags {
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
				flags |= widget.LayoutFlags()
			}
		}
	}

	return flags
}

func (l *GridLayout) MinSize() Size {
	if l.container == nil {
		return Size{}
	}

	widths := make([]int, len(l.cells[0]))
	heights := make([]int, len(l.cells))

	for row := 0; row < len(heights); row++ {
		for col := 0; col < len(widths); col++ {
			widget := l.cells[row][col].widget

			if !shouldLayoutWidget(widget) {
				continue
			}

			min := widget.MinSize()
			minHint := widget.MinSizeHint()

			widths[col] = maxi(widths[col], maxi(minHint.Width, min.Width))
			heights[row] = maxi(heights[row], maxi(minHint.Height, min.Height))
		}
	}

	width := l.margins.HNear + l.spacing*(len(widths)-1) + l.margins.HFar
	height := l.margins.VNear + l.spacing*(len(heights)-1) + l.margins.VFar

	for _, w := range widths {
		width += w
	}
	for _, h := range heights {
		height += h
	}

	return Size{width, height}
}

type gridLayoutSectionInfo struct {
	index              int
	minSize            int
	maxSize            int
	stretch            int
	hasGreedyNonSpacer bool
	hasGreedySpacer    bool
}

type gridLayoutSectionInfoList []gridLayoutSectionInfo

func (l gridLayoutSectionInfoList) Len() int {
	return len(l)
}

func (l gridLayoutSectionInfoList) Less(i, j int) bool {
	if l[i].hasGreedyNonSpacer == l[j].hasGreedyNonSpacer {
		if l[i].hasGreedySpacer == l[j].hasGreedySpacer {
			minDiff := l[i].minSize - l[j].minSize

			if minDiff == 0 {
				return l[i].maxSize/l[i].stretch < l[j].maxSize/l[j].stretch
			}

			return minDiff > 0
		}

		return l[i].hasGreedySpacer
	}

	return l[i].hasGreedyNonSpacer
}

func (l gridLayoutSectionInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l *GridLayout) Update(reset bool) os.Error {
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

		l.cleanup()
	}

	widths := l.sectionSizes(Horizontal)
	heights := l.sectionSizes(Vertical)

	hdwp := BeginDeferWindowPos(int32(l.container.Children().Len()))
	if hdwp == 0 {
		return lastError("BeginDeferWindowPos")
	}

	y := l.margins.VNear
	for row := 0; row < len(heights); row++ {
		h := heights[row]

		if h == 0 {
			continue
		}

		x := l.margins.HNear
		for col := 0; col < len(widths); col++ {
			w := widths[col]

			if w == 0 {
				continue
			}

			widget := l.cells[row][col].widget

			if widget != nil {
				// FIXME: This currently assumes all widgets can grow.
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
			}

			x += w + l.spacing
		}

		y += h + l.spacing
	}

	if !EndDeferWindowPos(hdwp) {
		return lastError("EndDeferWindowPos")
	}

	return nil
}

func (l *GridLayout) sectionSizes(orientation Orientation) []int {
	var stretchFactors []int
	if orientation == Horizontal {
		stretchFactors = l.columnStretchFactors
	} else {
		stretchFactors = l.rowStretchFactors
	}

	var sectionCountWithGreedyNonSpacer int
	var sectionCountWithGreedySpacer int
	var stretchFactorsTotal [3]int
	var minSizesRemaining int
	minSizes := make([]int, len(stretchFactors))
	maxSizes := make([]int, len(stretchFactors))
	sizes := make([]int, len(stretchFactors))
	sortedSections := gridLayoutSectionInfoList(make([]gridLayoutSectionInfo, len(stretchFactors)))

	for i := 0; i < len(stretchFactors); i++ {
		var otherAxisCount int
		if orientation == Horizontal {
			otherAxisCount = len(l.rowStretchFactors)
		} else {
			otherAxisCount = len(l.columnStretchFactors)
		}

		for j := 0; j < otherAxisCount; j++ {
			var widget Widget
			if orientation == Horizontal {
				widget = l.cells[j][i].widget
			} else {
				widget = l.cells[i][j].widget
			}

			if !shouldLayoutWidget(widget) {
				continue
			}

			flags := widget.LayoutFlags()

			min := widget.MinSize()
			max := widget.MaxSize()
			minHint := widget.MinSizeHint()
			pref := widget.SizeHint()

			if orientation == Horizontal {
				minSizes[i] = maxi(minSizes[i], maxi(min.Width, minHint.Width))

				if max.Width > 0 {
					maxSizes[i] = maxi(maxSizes[i], max.Width)
				} else if pref.Width > 0 && flags&GrowableHorz == 0 {
					maxSizes[i] = maxi(maxSizes[i], pref.Width)
				} else {
					maxSizes[i] = 32768
				}

				if flags&GreedyHorz > 0 {
					if _, isSpacer := widget.(*Spacer); isSpacer {
						sortedSections[i].hasGreedySpacer = true
					} else {
						sortedSections[i].hasGreedyNonSpacer = true
					}
				}
			} else {
				minSizes[i] = maxi(minSizes[i], maxi(min.Height, minHint.Height))

				if max.Height > 0 {
					maxSizes[i] = maxi(maxSizes[i], max.Height)
				} else if pref.Height > 0 && flags&GrowableVert == 0 {
					maxSizes[i] = maxi(maxSizes[i], pref.Height)
				} else {
					maxSizes[i] = 32768
				}

				if flags&GreedyVert > 0 {
					if _, isSpacer := widget.(*Spacer); isSpacer {
						sortedSections[i].hasGreedySpacer = true
					} else {
						sortedSections[i].hasGreedyNonSpacer = true
					}
				}
			}
		}

		sortedSections[i].index = i
		sortedSections[i].minSize = minSizes[i]
		sortedSections[i].maxSize = maxSizes[i]
		sortedSections[i].stretch = maxi(1, stretchFactors[i])

		minSizesRemaining += minSizes[i]

		if sortedSections[i].hasGreedyNonSpacer {
			sectionCountWithGreedyNonSpacer++
			stretchFactorsTotal[0] += stretchFactors[i]
		} else if sortedSections[i].hasGreedySpacer {
			sectionCountWithGreedySpacer++
			stretchFactorsTotal[1] += stretchFactors[i]
		} else {
			stretchFactorsTotal[2] += stretchFactors[i]
		}
	}

	sort.Sort(sortedSections)

	cb := l.container.ClientBounds()
	var space int
	if orientation == Horizontal {
		space = cb.Width - l.margins.HNear - l.margins.HFar
	} else {
		space = cb.Height - l.margins.VNear - l.margins.VFar
	}

	spacingRemaining := l.spacing * (len(stretchFactors) - 1)

	offsets := [3]int{0, sectionCountWithGreedyNonSpacer, sectionCountWithGreedyNonSpacer + sectionCountWithGreedySpacer}
	counts := [3]int{sectionCountWithGreedyNonSpacer, sectionCountWithGreedySpacer, len(stretchFactors) - sectionCountWithGreedyNonSpacer - sectionCountWithGreedySpacer}

	for i := 0; i < 3; i++ {
		stretchFactorsRemaining := stretchFactorsTotal[i]

		for j := 0; j < counts[i]; j++ {
			info := sortedSections[offsets[i]+j]
			k := info.index

			stretch := stretchFactors[k]
			min := info.minSize
			max := info.maxSize
			size := min

			if min < max {
				excessSpace := float64(space - minSizesRemaining - spacingRemaining)

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
			space -= (size + l.spacing)
			spacingRemaining -= l.spacing
		}
	}

	return sizes
}
