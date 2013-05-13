// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	mw := new(MyMainWindow)

	var outTE *walk.TextEdit

	foo := new(Foo)

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Walk Data Binding Example",
		MinSize:  Size{300, 200},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text: "Edit Foo",
				OnClicked: func() {
					res, err := RunFooDialog(mw, foo)
					if err != nil {
						log.Print(err)
					} else if res == walk.DlgCmdOK {
						outTE.SetText(fmt.Sprintf("%+v", foo))
					}
				},
			},
			Label{
				Text: "foo:",
			},
			TextEdit{
				AssignTo: &outTE,
				ReadOnly: true,
				Text:     fmt.Sprintf("%+v", foo),
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

type Foo struct {
	Name     string
	AnimalId int
	Weight   float64
	Fruit    string
	Eaten    bool
	Date     time.Time
	Memo     string
}

type Animal struct {
	Id   int
	Name string
}

func Animals() []*Animal {
	return []*Animal{
		{1, "Dog"},
		{2, "Cat"},
		{3, "Bird"},
		{4, "Fish"},
		{5, "Elephant"},
	}
}

type MyMainWindow struct {
	*walk.MainWindow
}

func RunFooDialog(owner walk.RootWidget, foo *Foo) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var ep walk.ErrorPresenter
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         "Foo Details",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			DataSource:     foo,
			ErrorPresenter: ErrorPresenterRef{&ep},
		},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{},
				Children: []Widget{
					Label{
						Row:    0,
						Column: 0,
						Text:   "Name:",
					},
					LineEdit{
						Row:    0,
						Column: 1,
						Text:   Bind("Name"),
					},
					Label{
						Row:    1,
						Column: 0,
						Text:   "Animal:",
					},
					ComboBox{
						Row:           1,
						Column:        1,
						Value:         Bind("AnimalId", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         Animals(),
					},
					Label{
						Row:    2,
						Column: 0,
						Text:   "Weight:",
					},
					NumberEdit{
						Row:      2,
						Column:   1,
						Value:    Bind("Weight", Range{0.01, 9999.99}),
						Suffix:   " kg",
						Decimals: 2,
					},
					Label{
						Row:    3,
						Column: 0,
						Text:   "Fruit:",
					},
					ComboBox{
						Row:      3,
						Column:   1,
						Editable: true,
						Value:    Bind("Fruit"),
						Model:    []string{"Banana", "Orange", "Cherry"},
					},
					Label{
						Row:    4,
						Column: 0,
						Text:   "Eaten:",
					},
					CheckBox{
						Row:     4,
						Column:  1,
						Checked: Bind("Eaten"),
					},
					Label{
						Row:    5,
						Column: 0,
						Text:   "Date:",
					},
					DateEdit{
						Row:    5,
						Column: 1,
						Date:   Bind("Date"),
					},
					VSpacer{
						Row:    6,
						Column: 0,
						Size:   8,
					},
					Label{
						Row:    7,
						Column: 0,
						Text:   "Memo:",
					},
					TextEdit{
						Row:        8,
						Column:     0,
						ColumnSpan: 2,
						MinSize:    Size{100, 50},
						Text:       Bind("Memo"),
					},
					LineErrorPresenter{
						AssignTo:   &ep,
						Row:        9,
						Column:     0,
						ColumnSpan: 2,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
