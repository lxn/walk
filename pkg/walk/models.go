// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"syscall"

	"github.com/miu200521358/win"
)

// BindingValueProvider is the interface that a model must implement to support
// data binding with widgets like ComboBox.
type BindingValueProvider interface {
	BindingValue(index int) interface{}
}

// ListModel is the interface that a model must implement to support widgets
// like ComboBox.
type ListModel interface {
	// ItemCount returns the number of items in the model.
	ItemCount() int

	// Value returns the value that should be displayed for the given index.
	Value(index int) interface{}

	// ItemsReset returns the event that the model should publish when the
	// number of its items changes.
	ItemsReset() *Event

	// ItemChanged returns the event that the model should publish when an item
	// was changed.
	ItemChanged() *IntEvent

	// ItemsInserted returns the event that the model should publish when a
	// contiguous range of items was inserted.
	ItemsInserted() *IntRangeEvent

	// ItemsRemoved returns the event that the model should publish when a
	// contiguous range of items was removed.
	ItemsRemoved() *IntRangeEvent
}

// ListModelBase implements the ItemsReset and ItemChanged methods of the
// ListModel interface.
type ListModelBase struct {
	itemsResetPublisher    EventPublisher
	itemChangedPublisher   IntEventPublisher
	itemsInsertedPublisher IntRangeEventPublisher
	itemsRemovedPublisher  IntRangeEventPublisher
}

func (lmb *ListModelBase) ItemsReset() *Event {
	return lmb.itemsResetPublisher.Event()
}

func (lmb *ListModelBase) ItemChanged() *IntEvent {
	return lmb.itemChangedPublisher.Event()
}

func (lmb *ListModelBase) ItemsInserted() *IntRangeEvent {
	return lmb.itemsInsertedPublisher.Event()
}

func (lmb *ListModelBase) ItemsRemoved() *IntRangeEvent {
	return lmb.itemsRemovedPublisher.Event()
}

func (lmb *ListModelBase) PublishItemsReset() {
	lmb.itemsResetPublisher.Publish()
}

func (lmb *ListModelBase) PublishItemChanged(index int) {
	lmb.itemChangedPublisher.Publish(index)
}

func (lmb *ListModelBase) PublishItemsInserted(from, to int) {
	lmb.itemsInsertedPublisher.Publish(from, to)
}

func (lmb *ListModelBase) PublishItemsRemoved(from, to int) {
	lmb.itemsRemovedPublisher.Publish(from, to)
}

// ReflectListModel provides an alternative to the ListModel interface. It
// uses reflection to obtain data.
type ReflectListModel interface {
	// Items returns the model data, which must be a slice of pointer to struct.
	Items() interface{}

	// ItemsReset returns the event that the model should publish when the
	// number of its items changes.
	ItemsReset() *Event

	// ItemChanged returns the event that the model should publish when an item
	// was changed.
	ItemChanged() *IntEvent

	// ItemsInserted returns the event that the model should publish when a
	// contiguous range of items was inserted.
	ItemsInserted() *IntRangeEvent

	// ItemsRemoved returns the event that the model should publish when a
	// contiguous range of items was removed.
	ItemsRemoved() *IntRangeEvent

	setValueFunc(value func(index int) interface{})
}

// ReflectListModelBase implements the ItemsReset and ItemChanged methods of
// the ReflectListModel interface.
type ReflectListModelBase struct {
	ListModelBase
	value func(index int) interface{}
}

func (rlmb *ReflectListModelBase) setValueFunc(value func(index int) interface{}) {
	rlmb.value = value
}

func (rlmb *ReflectListModelBase) Value(index int) interface{} {
	return rlmb.value(index)
}

// TableModel is the interface that a model must implement to support widgets
// like TableView.
type TableModel interface {
	// RowCount returns the number of rows in the model.
	RowCount() int

	// Value returns the value that should be displayed for the given cell.
	Value(row, col int) interface{}

	// RowsReset returns the event that the model should publish when the number
	// of its rows changes.
	RowsReset() *Event

	// RowChanged returns the event that the model should publish when a row was
	// changed.
	RowChanged() *IntEvent

	// RowsChanged returns the event that the model should publish when a
	// contiguous range of items was changed.
	RowsChanged() *IntRangeEvent

	// RowsInserted returns the event that the model should publish when a
	// contiguous range of items was inserted. If the model supports sorting, it
	// is assumed to be sorted before the model publishes the event.
	RowsInserted() *IntRangeEvent

	// RowsRemoved returns the event that the model should publish when a
	// contiguous range of items was removed.
	RowsRemoved() *IntRangeEvent
}

// TableModelBase implements the RowsReset and RowChanged methods of the
// TableModel interface.
type TableModelBase struct {
	rowsResetPublisher    EventPublisher
	rowChangedPublisher   IntEventPublisher
	rowsChangedPublisher  IntRangeEventPublisher
	rowsInsertedPublisher IntRangeEventPublisher
	rowsRemovedPublisher  IntRangeEventPublisher
}

func (tmb *TableModelBase) RowsReset() *Event {
	return tmb.rowsResetPublisher.Event()
}

func (tmb *TableModelBase) RowChanged() *IntEvent {
	return tmb.rowChangedPublisher.Event()
}

func (tmb *TableModelBase) RowsChanged() *IntRangeEvent {
	return tmb.rowsChangedPublisher.Event()
}

func (tmb *TableModelBase) RowsInserted() *IntRangeEvent {
	return tmb.rowsInsertedPublisher.Event()
}

func (tmb *TableModelBase) RowsRemoved() *IntRangeEvent {
	return tmb.rowsRemovedPublisher.Event()
}

func (tmb *TableModelBase) PublishRowsReset() {
	tmb.rowsResetPublisher.Publish()
}

func (tmb *TableModelBase) PublishRowChanged(row int) {
	tmb.rowChangedPublisher.Publish(row)
}

func (tmb *TableModelBase) PublishRowsChanged(from, to int) {
	tmb.rowsChangedPublisher.Publish(from, to)
}

func (tmb *TableModelBase) PublishRowsInserted(from, to int) {
	tmb.rowsInsertedPublisher.Publish(from, to)
}

func (tmb *TableModelBase) PublishRowsRemoved(from, to int) {
	tmb.rowsRemovedPublisher.Publish(from, to)
}

// ReflectTableModel provides an alternative to the TableModel interface. It
// uses reflection to obtain data.
type ReflectTableModel interface {
	// Items returns the model data, which must be a slice of pointer to struct.
	Items() interface{}

	// RowsReset returns the event that the model should publish when the
	// number of its items changes.
	RowsReset() *Event

	// RowChanged returns the event that the model should publish when an item
	// was changed.
	RowChanged() *IntEvent

	// RowsChanged returns the event that the model should publish when a
	// contiguous range of items was changed.
	RowsChanged() *IntRangeEvent

	// RowsInserted returns the event that the model should publish when a
	// contiguous range of items was inserted. If the model supports sorting, it
	// is assumed to be sorted before the model publishes the event.
	RowsInserted() *IntRangeEvent

	// RowsRemoved returns the event that the model should publish when a
	// contiguous range of items was removed.
	RowsRemoved() *IntRangeEvent

	setValueFunc(value func(row, col int) interface{})
}

// ReflectTableModelBase implements the ItemsReset and ItemChanged methods of
// the ReflectTableModel interface.
type ReflectTableModelBase struct {
	TableModelBase
	value func(row, col int) interface{}
}

func (rtmb *ReflectTableModelBase) setValueFunc(value func(row, col int) interface{}) {
	rtmb.value = value
}

func (rtmb *ReflectTableModelBase) Value(row, col int) interface{} {
	return rtmb.value(row, col)
}

type interceptedSorter interface {
	sorterBase() *SorterBase
	setSortFunc(sort func(col int, order SortOrder) error)
}

// SortedReflectTableModelBase implements the RowsReset and RowChanged methods
// of the ReflectTableModel interface as well as the Sorter interface for
// pre-implemented in-memory sorting.
type SortedReflectTableModelBase struct {
	ReflectTableModelBase
	SorterBase
	sort func(col int, order SortOrder) error
}

func (srtmb *SortedReflectTableModelBase) setSortFunc(sort func(col int, order SortOrder) error) {
	srtmb.sort = sort
}

func (srtmb *SortedReflectTableModelBase) sorterBase() *SorterBase {
	return &srtmb.SorterBase
}

func (srtmb *SortedReflectTableModelBase) Sort(col int, order SortOrder) error {
	if srtmb.sort != nil {
		return srtmb.sort(col, order)
	}

	return srtmb.SorterBase.Sort(col, order)
}

// Populator is an interface that can be implemented by Reflect*Models and slice
// types to populate themselves on demand.
//
// Widgets like TableView, ListBox and ComboBox support lazy population of a
// Reflect*Model or slice, if it implements this interface.
type Populator interface {
	// Populate initializes the slot specified by index.
	//
	// For best performance it is probably a good idea to populate more than a
	// single slot of the slice at once.
	Populate(index int) error
}

// ImageProvider is the interface that a model must implement to support
// displaying an item image.
type ImageProvider interface {
	// Image returns the image to display for the item at index index.
	//
	// Supported types are *walk.Bitmap, *walk.Icon and string. A string will be
	// interpreted as a file path and the icon associated with the file will be
	// used. It is not supported to use strings together with the other options
	// in the same model instance.
	Image(index int) interface{}
}

// CellStyler is the interface that must be implemented to provide a tabular
// widget like TableView with cell display style information.
type CellStyler interface {
	// StyleCell is called for each cell to pick up cell style information.
	StyleCell(style *CellStyle)
}

// CellStyle carries information about the display style of a cell in a tabular widget
// like TableView.
type CellStyle struct {
	row             int
	col             int
	bounds          Rectangle // in native pixels
	hdc             win.HDC
	dpi             int
	canvas          *Canvas
	BackgroundColor Color
	TextColor       Color
	Font            *Font

	// Image is the image to display in the cell.
	//
	// Supported types are *walk.Bitmap, *walk.Icon and string. A string will be
	// interpreted as a file path and the icon associated with the file will be
	// used. It is not supported to use strings together with the other options
	// in the same model instance.
	Image interface{}
}

func (cs *CellStyle) Row() int {
	return cs.row
}

func (cs *CellStyle) Col() int {
	return cs.col
}

func (cs *CellStyle) Bounds() Rectangle {
	return RectangleTo96DPI(cs.bounds, cs.dpi)
}

func (cs *CellStyle) BoundsPixels() Rectangle {
	return cs.bounds
}

func (cs *CellStyle) Canvas() *Canvas {
	if cs.canvas != nil {
		cs.canvas.dpi = cs.dpi
		return cs.canvas
	}

	if cs.hdc != 0 {
		cs.canvas, _ = newCanvasFromHDC(cs.hdc)
		cs.canvas.dpi = cs.dpi
	}

	return cs.canvas
}

// IDProvider is the interface that must be implemented by models to enable
// widgets like TableView to attempt keeping the current item when the model
// publishes a reset event.
type IDProvider interface {
	ID(index int) interface{}
}

// ListItemStyler is the interface that must be implemented to provide a list
// widget like ListBox with item display style information.
type ListItemStyler interface {
	// ItemHeightDependsOnWidth returns whether item height depends on width.
	ItemHeightDependsOnWidth() bool

	// DefaultItemHeight returns the initial height in native pixels for any item.
	DefaultItemHeight() int

	// ItemHeight is called for each item to retrieve the height of the item. width parameter and
	// return value are specified in native pixels.
	ItemHeight(index int, width int) int

	// StyleItem is called for each item to pick up item style information.
	StyleItem(style *ListItemStyle)
}

// ListItemStyle carries information about the display style of an item in a list widget
// like ListBox.
type ListItemStyle struct {
	BackgroundColor    Color
	TextColor          Color
	defaultTextColor   Color
	LineColor          Color
	Font               *Font
	index              int
	hoverIndex         int
	rc                 win.RECT
	bounds             Rectangle // in native pixels
	state              uint32
	hTheme             win.HTHEME
	hwnd               win.HWND
	hdc                win.HDC
	dpi                int
	canvas             *Canvas
	highContrastActive bool
}

func (lis *ListItemStyle) Index() int {
	return lis.index
}

func (lis *ListItemStyle) Bounds() Rectangle {
	return RectangleTo96DPI(lis.bounds, lis.dpi)
}

func (lis *ListItemStyle) BoundsPixels() Rectangle {
	return lis.bounds
}

func (lis *ListItemStyle) Canvas() *Canvas {
	if lis.canvas != nil {
		lis.canvas.dpi = lis.dpi
		return lis.canvas
	}

	if lis.hdc != 0 {
		lis.canvas, _ = newCanvasFromHDC(lis.hdc)
		lis.canvas.dpi = lis.dpi
	}

	return lis.canvas
}

func (lis *ListItemStyle) DrawBackground() error {
	canvas := lis.Canvas()
	if canvas == nil {
		return nil
	}

	stateID := lis.stateID()

	if lis.hTheme != 0 && stateID != win.LISS_NORMAL {
		if win.FAILED(win.DrawThemeBackground(lis.hTheme, lis.hdc, win.LVP_LISTITEM, stateID, &lis.rc, nil)) {
			return newError("DrawThemeBackground failed")
		}
	} else {
		brush, err := NewSolidColorBrush(lis.BackgroundColor)
		if err != nil {
			return err
		}
		defer brush.Dispose()

		if err := canvas.FillRectanglePixels(brush, lis.bounds); err != nil {
			return err
		}

		if lis.highContrastActive && (lis.index == lis.hoverIndex || stateID != win.LISS_NORMAL) {
			pen, err := NewCosmeticPen(PenSolid, lis.LineColor)
			if err != nil {
				return err
			}
			defer pen.Dispose()

			if err := canvas.DrawRectanglePixels(pen, lis.bounds); err != nil {
				return err
			}
		}
	}

	return nil
}

// DrawText draws text inside given bounds specified in native pixels.
func (lis *ListItemStyle) DrawText(text string, bounds Rectangle, format DrawTextFormat) error {
	if lis.hTheme != 0 && lis.TextColor == lis.defaultTextColor {
		if lis.Font != nil {
			hFontOld := win.SelectObject(lis.hdc, win.HGDIOBJ(lis.Font.handleForDPI(lis.dpi)))
			defer win.SelectObject(lis.hdc, hFontOld)
		}
		rc := bounds.toRECT()

		if win.FAILED(win.DrawThemeTextEx(lis.hTheme, lis.hdc, win.LVP_LISTITEM, lis.stateID(), syscall.StringToUTF16Ptr(text), int32(len(([]rune)(text))), uint32(format), &rc, nil)) {
			return newError("DrawThemeTextEx failed")
		}
	} else {
		if canvas := lis.Canvas(); canvas != nil {
			if err := canvas.DrawTextPixels(text, lis.Font, lis.TextColor, bounds, format); err != nil {
				return err
			}
		}
	}

	return nil
}

func (lis *ListItemStyle) stateID() int32 {
	if lis.state&win.ODS_CHECKED != 0 {
		if win.GetFocus() == lis.hwnd {
			if lis.index == lis.hoverIndex {
				return win.LISS_HOTSELECTED
			} else {
				return win.LISS_SELECTED
			}
		} else {
			return win.LISS_SELECTEDNOTFOCUS
		}
	} else if lis.index == lis.hoverIndex {
		return win.LISS_HOT
	}

	return win.LISS_NORMAL
}

// ItemChecker is the interface that a model must implement to support check
// boxes in a widget like TableView.
type ItemChecker interface {
	// Checked returns if the specified item is checked.
	Checked(index int) bool

	// SetChecked sets if the specified item is checked.
	SetChecked(index int, checked bool) error
}

// SortOrder specifies the order by which items are sorted.
type SortOrder int

const (
	// SortAscending specifies ascending sort order.
	SortAscending SortOrder = iota

	// SortDescending specifies descending sort order.
	SortDescending
)

// Sorter is the interface that a model must implement to support sorting with a
// widget like TableView.
type Sorter interface {
	// ColumnSortable returns whether column col is sortable.
	ColumnSortable(col int) bool

	// Sort sorts column col in order order.
	//
	// If col is -1 then no column is to be sorted. Sort must publish the event
	// returned from SortChanged() after sorting.
	Sort(col int, order SortOrder) error

	// SortChanged returns an event that is published after sorting.
	SortChanged() *Event

	// SortedColumn returns the index of the currently sorted column, or -1 if
	// no column is currently sorted.
	SortedColumn() int

	// SortOrder returns the current sort order.
	SortOrder() SortOrder
}

// SorterBase implements the Sorter interface.
//
// You still need to provide your own implementation of at least the Sort method
// to actually sort and reset the model. Your Sort method should call the
// SorterBase implementation so the SortChanged event, that e.g. a TableView
// widget depends on, is published.
type SorterBase struct {
	changedPublisher EventPublisher
	col              int
	order            SortOrder
}

func (sb *SorterBase) ColumnSortable(col int) bool {
	return true
}

func (sb *SorterBase) Sort(col int, order SortOrder) error {
	sb.col, sb.order = col, order

	sb.changedPublisher.Publish()

	return nil
}

func (sb *SorterBase) SortChanged() *Event {
	return sb.changedPublisher.Event()
}

func (sb *SorterBase) SortedColumn() int {
	return sb.col
}

func (sb *SorterBase) SortOrder() SortOrder {
	return sb.order
}

// Imager provides access to an image of objects like tree items.
type Imager interface {
	// Image returns the image to display for an item.
	//
	// Supported types are *walk.Bitmap, *walk.Icon and string. A string will be
	// interpreted as a file path and the icon associated with the file will be
	// used. It is not supported to use strings together with the other options
	// in the same model instance.
	Image() interface{}
}

// TreeItem represents an item in a TreeView widget.
type TreeItem interface {
	// Text returns the text of the item.
	Text() string

	// Parent returns the parent of the item.
	Parent() TreeItem

	// ChildCount returns the number of children of the item.
	ChildCount() int

	// ChildAt returns the child at the specified index.
	ChildAt(index int) TreeItem
}

// HasChilder enables widgets like TreeView to determine if an item has any
// child, without enforcing to fully count all children.
type HasChilder interface {
	HasChild() bool
}

// TreeModel provides widgets like TreeView with item data.
type TreeModel interface {
	// LazyPopulation returns if the model prefers on-demand population.
	//
	// This is useful for models that potentially contain huge amounts of items,
	// e.g. a model that represents a file system.
	LazyPopulation() bool

	// RootCount returns the number of root items.
	RootCount() int

	// RootAt returns the root item at the specified index.
	RootAt(index int) TreeItem

	// ItemsReset returns the event that the model should publish when the
	// descendants of the specified item, or all items if no item is specified,
	// are reset.
	ItemsReset() *TreeItemEvent

	// ItemChanged returns the event that the model should publish when an item
	// was changed.
	ItemChanged() *TreeItemEvent

	// ItemInserted returns the event that the model should publish when an item
	// was inserted into the model.
	ItemInserted() *TreeItemEvent

	// ItemRemoved returns the event that the model should publish when an item
	// was removed from the model.
	ItemRemoved() *TreeItemEvent
}

// TreeModelBase partially implements the TreeModel interface.
//
// You still need to provide your own implementation of at least the
// RootCount and RootAt methods. If your model needs lazy population,
// you will also have to implement LazyPopulation.
type TreeModelBase struct {
	itemsResetPublisher   TreeItemEventPublisher
	itemChangedPublisher  TreeItemEventPublisher
	itemInsertedPublisher TreeItemEventPublisher
	itemRemovedPublisher  TreeItemEventPublisher
}

func (tmb *TreeModelBase) LazyPopulation() bool {
	return false
}

func (tmb *TreeModelBase) ItemsReset() *TreeItemEvent {
	return tmb.itemsResetPublisher.Event()
}

func (tmb *TreeModelBase) ItemChanged() *TreeItemEvent {
	return tmb.itemChangedPublisher.Event()
}

func (tmb *TreeModelBase) ItemInserted() *TreeItemEvent {
	return tmb.itemInsertedPublisher.Event()
}

func (tmb *TreeModelBase) ItemRemoved() *TreeItemEvent {
	return tmb.itemRemovedPublisher.Event()
}

func (tmb *TreeModelBase) PublishItemsReset(parent TreeItem) {
	tmb.itemsResetPublisher.Publish(parent)
}

func (tmb *TreeModelBase) PublishItemChanged(item TreeItem) {
	tmb.itemChangedPublisher.Publish(item)
}

func (tmb *TreeModelBase) PublishItemInserted(item TreeItem) {
	tmb.itemInsertedPublisher.Publish(item)
}

func (tmb *TreeModelBase) PublishItemRemoved(item TreeItem) {
	tmb.itemRemovedPublisher.Publish(item)
}
