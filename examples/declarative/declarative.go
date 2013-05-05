// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"
)

import (
	"github.com/lxn/polyglot"
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

func Foos() []*Foo {
	return []*Foo{
		{1, "One"},
		{2, "Two"},
		{3, "Three"},
	}
}

type Bar struct {
	Key  string
	Text string
}

func Bars() []*Bar {
	return []*Bar{
		{"one", "1"},
		{"two", "2"},
		{"three", "3"},
	}
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
	var errorPresenter walk.ErrorPresenter
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
		DataBinder: DataBinder{
			AssignTo:       &dataBinder,
			DataSource:     db.DataSource,
			ErrorPresenter: ErrorPresenterRef{&errorPresenter},
		},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout:   Grid{},
				Children: db.Widgets,
			},
			Composite{
				Layout: VBox{Margins: Margins{9, 0, 9, 0}},
				Children: []Widget{
					LineErrorPresenter{AssignTo: &errorPresenter},
				},
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
		Label{
			Row: 0, Column: 0,
			Text: "Name:",
		},
		LineEdit{
			Row: 0, Column: 1,
			Name:    "nameLE",
			Enabled: Bind("enabledCB.Checked"),
			Text:    Bind("Name", Regexp{`^[A-Z][a-z]*$`}),
		},

		Label{
			Row: 1, Column: 0,
			Text: "No.:",
		},
		LineEdit{
			Row: 1, Column: 1,
			Name:      "noLE",
			Enabled:   Bind("enabledCB.Checked"),
			Text:      Bind("No", Regexp{`^[\d]{3}[ ]{1}[\d]{3}$`}),
			MaxLength: 7,
			CueBanner: "### ###",
		},

		Label{
			Row: 2, Column: 0,
			Text: "Foo (int BindingValue):",
		},
		ComboBox{
			Row: 2, Column: 1,
			Name:          "fooIdCB",
			BindingMember: "Id",
			DisplayMember: "Text",
			Model:         Foos(),
			Value:         Bind("FooId", SelRequired{}),
		},

		Label{
			Row: 3, Column: 0,
			Text: "Bar (string BindingValue):",
		},
		ComboBox{
			Row: 3, Column: 1,
			BindingMember: "Key",
			DisplayMember: "Text",
			Model:         Bars(),
			Value:         Bind("BarKey"),
		},

		Label{
			Row: 4, Column: 0,
			Text: "String:",
		},
		ComboBox{
			Row: 4, Column: 1,
			Editable: true,
			Model:    []string{"One", "Two", "Three"},
			Value:    Bind("String"),
		},

		Label{
			Row: 5, Column: 0,
			Text: "Float64:",
		},
		NumberEdit{
			Row: 5, Column: 1,
			Value:    Bind("Float64", Range{0.01, 999.99}),
			Decimals: 2,
		},

		Label{
			Row: 6, Column: 0,
			Text: "Int:",
		},
		NumberEdit{
			Row: 6, Column: 1,
			Value: Bind("Int"),
		},

		Label{
			Row: 7, Column: 0,
			Text: "Date:",
		},
		DateEdit{
			Row: 7, Column: 1,
			Date: Bind("Date"),
		},

		Label{
			Row: 8, Column: 0,
			Text: "Enabled:",
		},
		CheckBox{
			Row: 8, Column: 1,
			Name:    "enabledCB",
			Checked: Bind("Enabled"),
		},

		VSpacer{
			Row: 9, Column: 0,
			Size: 10,
		},

		Label{
			Row: 10, Column: 0, ColumnSpan: 2,
			Text: "Memo:",
		},
		TextEdit{
			Row: 11, Column: 0, ColumnSpan: 2,
			Text: Bind("Memo"),
		},
	}

	type Item struct {
		Name    string
		No      string
		FooId   int
		BarKey  string
		String  string
		Float64 float64
		Int     int
		Date    time.Time
		Enabled bool
		Memo    string
	}

	item := &Item{
		Name:    "Name",
		Int:     67890,
		Date:    time.Now(),
		Enabled: true,
		Memo:    "Memo",
	}

	db := &DialogBuilder{
		Title:      "My Dialog",
		Owner:      mw,
		Dialog:     &dlg.Dialog,
		Widgets:    widgets,
		MinSize:    Size{0, 480},
		DataSource: item,
	}

	if err := db.Build(); err != nil {
		log.Fatal(err)
	}

	if dlg.Run() == walk.DlgCmdOK {
		log.Printf("item: %+v", item)
	}
}

var trDict *polyglot.Dict

func tr(source string, context ...string) string {
	return trDict.Translation(source, context...)
}

func main() {
	walk.SetTranslationFunc(tr)

	var err error
	if trDict, err = polyglot.NewDict("../../l10n", "en"); err != nil {
		log.Fatal(err)
	}

	mw := new(MyMainWindow)

	var openAction *walk.Action
	var recentMenu *walk.Menu

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Walk Declarative Example",
		MenuItems: []MenuItem{
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
			},
		},
		ToolBarItems: []MenuItem{
			ActionRef{&openAction},
			Separator{},
			Action{
				Text:        "Show Dialog",
				OnTriggered: func() { mw.showDialogAction_Triggered() },
			},
		},
		MinSize: Size{600, 400},
		Size:    Size{1024, 768},
		Layout:  HBox{MarginsZero: true},
		Children: []Widget{
			TabWidget{
				ContentMarginsZero: true,
				Pages: []TabPage{
					//					TabPage{Title: "golang.org/doc/", Content: WebView{URL: "http://golang.org/doc/"}},
					//					TabPage{Title: "golang.org/ref/", Content: WebView{URL: "http://golang.org/ref/"}},
					//					TabPage{Title: "golang.org/pkg/", Content: WebView{URL: "http://golang.org/pkg/"}},
					TabPage{
						Title:  "Composite Stuff",
						Layout: Grid{},
						Children: []Widget{
							TextEdit{Row: 0, Column: 0, RowSpan: 4},
							PushButton{Row: 0, Column: 1, Text: "Foo"},
							PushButton{Row: 1, Column: 1, Text: "Bar"},
							PushButton{Row: 2, Column: 1, Text: "Baz"},
							VSpacer{Row: 3, Column: 1},
						},
					},
				},
			},
		},
	}.Create()); err != nil {
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

	mw.Run()
}
