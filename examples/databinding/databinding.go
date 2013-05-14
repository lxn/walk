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

	animal := new(Animal)

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Walk Data Binding Example",
		MinSize:  Size{300, 200},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text: "Edit Animal",
				OnClicked: func() {
					res, err := RunAnimalDialog(mw, animal)
					if err != nil {
						log.Print(err)
					} else if res == walk.DlgCmdOK {
						outTE.SetText(fmt.Sprintf("%+v", animal))
					}
				},
			},
			Label{
				Text: "animal:",
			},
			TextEdit{
				AssignTo: &outTE,
				ReadOnly: true,
				Text:     fmt.Sprintf("%+v", animal),
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

type Animal struct {
	Name          string
	ArrivalDate   time.Time
	SpeciesId     int
	Sex           Sex
	Weight        float64
	PreferredFood string
	Domesticated  bool
	Remarks       string
}

type Species struct {
	Id   int
	Name string
}

func KnownSpecies() []*Species {
	return []*Species{
		{1, "Dog"},
		{2, "Cat"},
		{3, "Bird"},
		{4, "Fish"},
		{5, "Elephant"},
	}
}

type Sex byte

const (
	SexMale Sex = 1 + iota
	SexFemale
	SexHermaphrodite
)

type MyMainWindow struct {
	*walk.MainWindow
}

func RunAnimalDialog(owner walk.RootWidget, animal *Animal) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var ep walk.ErrorPresenter
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         "Animal Details",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			DataSource:     animal,
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
						Text:   "Arrival Date:",
					},
					DateEdit{
						Row:    1,
						Column: 1,
						Date:   Bind("ArrivalDate"),
					},
					Label{
						Row:    2,
						Column: 0,
						Text:   "Species:",
					},
					ComboBox{
						Row:           2,
						Column:        1,
						Value:         Bind("SpeciesId", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownSpecies(),
					},
					RadioButtonGroupBox{
						Row:        3,
						Column:     0,
						ColumnSpan: 2,
						Title:      "Sex",
						Layout:     HBox{},
						DataMember: "Sex",
						Buttons: []RadioButton{
							RadioButton{
								Text:  "Male",
								Value: SexMale,
							},
							RadioButton{
								Text:  "Female",
								Value: SexFemale,
							},
							RadioButton{
								Text:  "Hermaphrodite",
								Value: SexHermaphrodite,
							},
						},
					},
					Label{
						Row:    4,
						Column: 0,
						Text:   "Weight:",
					},
					NumberEdit{
						Row:      4,
						Column:   1,
						Value:    Bind("Weight", Range{0.01, 9999.99}),
						Suffix:   " kg",
						Decimals: 2,
					},
					Label{
						Row:    5,
						Column: 0,
						Text:   "Preferred Food:",
					},
					ComboBox{
						Row:      5,
						Column:   1,
						Editable: true,
						Value:    Bind("PreferredFood"),
						Model:    []string{"Fruits", "Gras", "Fish", "Meat"},
					},
					Label{
						Row:    6,
						Column: 0,
						Text:   "Domesticated:",
					},
					CheckBox{
						Row:     6,
						Column:  1,
						Checked: Bind("Domesticated"),
					},
					VSpacer{
						Row:    7,
						Column: 0,
						Size:   8,
					},
					Label{
						Row:    8,
						Column: 0,
						Text:   "Remarks:",
					},
					TextEdit{
						Row:        9,
						Column:     0,
						ColumnSpan: 2,
						MinSize:    Size{100, 50},
						Text:       Bind("Remarks"),
					},
					LineErrorPresenter{
						AssignTo:   &ep,
						Row:        10,
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
