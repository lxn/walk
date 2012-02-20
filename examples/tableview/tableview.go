// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"rand"
	"strings"
	"time"
)

import "github.com/lxn/walk"

type Foo struct {
	Bar     string
	Baz     float64
	Quux    int64
	checked bool
}

type FooModel struct {
	items               []*Foo
	rowsResetPublisher  walk.EventPublisher
	rowChangedPublisher walk.IntEventPublisher
}

// Make sure we implement all required interfaces.
var _ walk.TableModel = &FooModel{}
var _ walk.ItemChecker = &FooModel{}

// Called by the TableView from SetModel to retrieve column information. 
func (m *FooModel) Columns() []walk.TableColumn {
	return []walk.TableColumn{
		{Title: "#"},
		{Title: "Bar"},
		{Title: "Baz", Format: "%.2f", Alignment: walk.AlignFar},
		{Title: "Quux", Format: "2006-01-02 15:04:05", Width: 150},
	}
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return row

	case 1:
		return item.Bar

	case 2:
		return item.Baz

	case 3:
		return time.SecondsToLocalTime(item.Quux)
	}

	panic("unexpected col")
}

// The TableView attaches to this event to synchronize its internal item count.
func (m *FooModel) RowsReset() *walk.Event {
	return m.rowsResetPublisher.Event()
}

// The TableView attaches to this event to get notified when a row changed and
// needs to be repainted.
func (m *FooModel) RowChanged() *walk.IntEvent {
	return m.rowChangedPublisher.Event()
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) os.Error {
	m.items[row].checked = checked

	return nil
}

func (m *FooModel) ResetRows() {
	// Create some random data.
	m.items = make([]*Foo, rand.Intn(50000))

	now := time.Seconds()

	for i := range m.items {
		m.items[i] = &Foo{
			Bar:  strings.Repeat("*", rand.Intn(5)+1),
			Baz:  rand.Float64() * 1000,
			Quux: rand.Int63n(now),
		}
	}

	// Notify TableView and other interested parties about the reset.
	m.rowsResetPublisher.Publish()
}

type MainWindow struct {
	*walk.MainWindow
	model *FooModel
}

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	rand.Seed(time.Seconds())

	mainWnd, _ := walk.NewMainWindow()

	mw := &MainWindow{
		MainWindow: mainWnd,
		model:      &FooModel{},
	}

	// We want the model to be populated right from the beginning.
	mw.model.ResetRows()

	mw.SetLayout(walk.NewVBoxLayout())
	mw.SetTitle("Walk TableView Example")

	fileMenu, _ := walk.NewMenu()
	fileMenuAction, _ := mw.Menu().Actions().AddMenu(fileMenu)
	fileMenuAction.SetText("&File")

	exitAction := walk.NewAction()
	exitAction.SetText("E&xit")
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	fileMenu.Actions().Add(exitAction)

	helpMenu, _ := walk.NewMenu()
	helpMenuAction, _ := mw.Menu().Actions().AddMenu(helpMenu)
	helpMenuAction.SetText("&Help")

	aboutAction := walk.NewAction()
	aboutAction.SetText("&About")
	aboutAction.Triggered().Attach(func() {
		walk.MsgBox(mw, "About", "Walk TableView Example", walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	resetRowsButton, _ := walk.NewPushButton(mw)
	resetRowsButton.SetText("Reset Rows")

	resetRowsButton.Clicked().Attach(func() {
		// Get some fresh data.
		mw.model.ResetRows()
	})

	tableView, _ := walk.NewTableView(mw)

	// Everybody loves check boxes.
	tableView.SetCheckBoxes(true)

	// Don't forget to set the model.
	tableView.SetModel(mw.model)

	mw.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	mw.SetSize(walk.Size{800, 600})
	mw.Show()

	mw.Run()
}
