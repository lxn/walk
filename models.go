// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

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
}

// ListModelBase implements the ItemsReset and ItemChanged methods of the
// ListModel interface.
type ListModelBase struct {
	itemsResetPublisher  EventPublisher
	itemChangedPublisher IntEventPublisher
}

func (lmb *ListModelBase) ItemsReset() *Event {
	return lmb.itemsResetPublisher.Event()
}

func (lmb *ListModelBase) ItemChanged() *IntEvent {
	return lmb.itemChangedPublisher.Event()
}

func (lmb *ListModelBase) PublishItemsReset() {
	lmb.itemsResetPublisher.Publish()
}

func (lmb *ListModelBase) PublishItemChanged(index int) {
	lmb.itemChangedPublisher.Publish(index)
}

// TableColumn provides column information for widgets like TableView.
type TableColumn struct {
	// Name is the optional name of the column.
	Name string

	// Title is the text to display in the column header.
	Title string

	// Format is the format string for converting a value into a string.
	Format string

	// Precision is the number of decimal places for formatting a big.Rat.
	Precision int

	// Width is the width of the column in pixels.
	Width int

	// Alignment is the alignment of the column (who would have thought).
	Alignment Alignment1D
}

// TableModel is the interface that a model must implement to support widgets
// like TableView.
type TableModel interface {
	// Columns returns information about the columns of the model.
	Columns() []TableColumn

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
}

// TableModelBase implements the RowsReset and RowChanged methods of the
// TableModel interface.
type TableModelBase struct {
	rowsResetPublisher  EventPublisher
	rowChangedPublisher IntEventPublisher
}

func (tmb *TableModelBase) RowsReset() *Event {
	return tmb.rowsResetPublisher.Event()
}

func (tmb *TableModelBase) RowChanged() *IntEvent {
	return tmb.rowChangedPublisher.Event()
}

func (tmb *TableModelBase) PublishRowsReset() {
	tmb.rowsResetPublisher.Publish()
}

func (tmb *TableModelBase) PublishRowChanged(row int) {
	tmb.rowChangedPublisher.Publish(row)
}

// ItemChecker is the interface that a model must implement to support check 
// boxes in a widget like TableView.
type ItemChecker interface {
	// Checked returns if the specified item is checked.
	Checked(index int) bool

	// SetChecked sets if the specified item is checked.
	SetChecked(index int, checked bool) error
}
