// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
}

type MyDialog struct {
	*walk.Dialog
}

type Foo struct {
	Id   int
	Text string
}

type FooModel struct {
	walk.ListModelBase
	items []*Foo
}

func NewFooModel() *FooModel {
	return &FooModel{
		items: []*Foo{
			{1, "One"},
			{2, "Two"},
			{3, "Three"},
		},
	}
}

func (m *FooModel) ItemCount() int {
	return len(m.items)
}

func (m *FooModel) Value(index int) interface{} {
	return m.items[index].Text
}

func (m *FooModel) BindingValue(index int) interface{} {
	return m.items[index].Id
}

type Bar struct {
	Key  string
	Text string
}

type BarModel struct {
	walk.ListModelBase
	items []*Bar
}

func NewBarModel() *BarModel {
	return &BarModel{
		items: []*Bar{
			{"one", "1"},
			{"two", "2"},
			{"three", "3"},
		},
	}
}

func (m *BarModel) ItemCount() int {
	return len(m.items)
}

func (m *BarModel) Value(index int) interface{} {
	return m.items[index].Text
}

func (m *BarModel) BindingValue(index int) interface{} {
	return m.items[index].Key
}

type DialogBuilder struct {
	Owner      walk.RootWidget
	Dialog     **walk.Dialog
	Widgets    []Widget
	Title      string
	Size       Size
	MinSize    Size
	DataSource interface{}
}

func (db *DialogBuilder) Build() error {
	var dataBinder *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	onAcceptClicked := func() {
		if err := dataBinder.Submit(); err != nil {
			log.Fatal(err)
		}

		db.Dialog.Accept()
	}

	return Dialog{
		AssignTo:      db.Dialog,
		Title:         db.Title,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       db.MinSize,
		Size:          db.Size,
		DataBinder:    DataBinder{AssignTo: &dataBinder, DataSource: db.DataSource},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout:   Grid{},
				Children: db.Widgets,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{AssignTo: &acceptPB, Text: "OK", OnClicked: onAcceptClicked},
					PushButton{AssignTo: &cancelPB, Text: "Cancel", OnClicked: func() { db.Dialog.Cancel() }},
				},
			},
		},
	}.Create(db.Owner)
}

func (mw *MyMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Nothing to see here...", walk.MsgBoxIconInformation|walk.MsgBoxOK)
}

func (mw *MyMainWindow) showDialogAction_Triggered() {
	dlg := new(MyDialog)

	widgets := []Widget{
		Label{Row: 0, Column: 0, Text: "Name:"},
		Label{Row: 0, Column: 1, BindTo: "Name"},
		Label{Row: 1, Column: 0, Text: "Short Text:"},
		LineEdit{Row: 1, Column: 1, BindTo: "ShortText"},
		Label{Row: 2, Column: 0, Text: "Foo (int BindingValue):"},
		ComboBox{Row: 2, Column: 1, BindTo: "FooId", Model: NewFooModel()},
		Label{Row: 3, Column: 0, Text: "Bar (string BindingValue):"},
		ComboBox{Row: 3, Column: 1, BindTo: "BarKey", Model: NewBarModel()},
		Label{Row: 4, Column: 0, Text: "Float64:"},
		NumberEdit{Row: 4, Column: 1, BindTo: "Float64", Decimals: 2},
		Label{Row: 5, Column: 0, Text: "Int:"},
		NumberEdit{Row: 5, Column: 1, BindTo: "Int"},
		Label{Row: 6, Column: 0, Text: "Date:"},
		DateEdit{Row: 6, Column: 1, BindTo: "Date"},
		Label{Row: 7, Column: 0, Text: "Checked:"},
		CheckBox{Row: 7, Column: 1, BindTo: "Checked"},
		VSpacer{Row: 8, Column: 0, Size: 10},
		Label{Row: 9, Column: 0, ColumnSpan: 2, Text: "Memo:"},
		TextEdit{Row: 10, Column: 0, ColumnSpan: 2, BindTo: "Memo"},
	}

	type Item struct {
		Name      string
		ShortText string
		FooId     int
		BarKey    string
		Float64   float64
		Int       int
		Date      time.Time
		Checked   bool
		Memo      string
	}

	item := &Item{
		Name:      "Name",
		ShortText: "ShortText",
		FooId:     2,
		BarKey:    "three",
		Float64:   123.45,
		Int:       67890,
		Date:      time.Now(),
		Checked:   true,
		Memo:      "Memo",
	}

	db := &DialogBuilder{
		Title:      "My Dialog",
		Owner:      mw,
		Dialog:     &dlg.Dialog,
		Widgets:    widgets,
		MinSize:    Size{320, 480},
		DataSource: item,
	}

	if err := db.Build(); err != nil {
		log.Fatal(err)
	}

	if dlg.Run() == walk.DlgCmdOK {
		log.Printf("item: %+v", item)
	}
}

func main() {
	walk.Initialize(walk.InitParams{})
	defer walk.Shutdown()

	mw := new(MyMainWindow)

	var openAction *walk.Action
	var recentMenu *walk.Menu

	menuActions, err := CreateActions(
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Action{
					AssignTo:    &openAction,
					Text:        "&Open",
					Image:       "../img/open.png",
					OnTriggered: func() { mw.openAction_Triggered() },
				},
				Menu{
					AssignTo: &recentMenu,
					Text:     "Recent",
				},
				Separator{},
				Action{
					Text:        "E&xit",
					OnTriggered: func() { walk.App().Exit(0) },
				},
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	openRecent1Action := walk.NewAction()
	openRecent1Action.SetText("Blah")
	recentMenu.Actions().Add(openRecent1Action)

	openRecent2Action := walk.NewAction()
	openRecent2Action.SetText("Yadda")
	recentMenu.Actions().Add(openRecent2Action)

	openRecent3Action := walk.NewAction()
	openRecent3Action.SetText("Oink")
	recentMenu.Actions().Add(openRecent3Action)

	toolBarActions, err := CreateActions(
		ActionRef{openAction},
		Separator{},
		Action{Text: "Show Dialog", OnTriggered: func() { mw.showDialogAction_Triggered() }})
	if err != nil {
		log.Fatal(err)
	}

	if err := (MainWindow{
		AssignTo:       &mw.MainWindow,
		Title:          "Walk Declarative Example",
		MenuActions:    menuActions,
		ToolBarActions: toolBarActions,
		MinSize:        Size{600, 400},
		Size:           Size{1024, 768},
		Layout:         HBox{MarginsZero: true},
		Children: []Widget{
			TabWidget{
				MarginsZero: true,
				PageTitles:  []string{"golang.org/doc/", "golang.org/ref/", "golang.org/pkg/"},
				Pages: []Widget{
					WebView{URL: "http://golang.org/doc/"},
					WebView{URL: "http://golang.org/ref/"},
					WebView{URL: "http://golang.org/pkg/"},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	mw.Show()
	mw.Run()
}
