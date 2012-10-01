// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
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

func (mw *MyMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Nothing to see here...", walk.MsgBoxIconInformation|walk.MsgBoxOK)
}

func (mw *MyMainWindow) showDialogAction_Triggered() {
	dlg := new(MyDialog)

	var acceptPB *walk.PushButton
	var cancelPB *walk.PushButton

	if err := (Dialog{
		AssignTo:      &dlg.Dialog,
		Title:         "My Dialog",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{400, 300},
		Size:          Size{400, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{},
				Children: []Widget{
					Label{Row: 0, Column: 0, Text: "A LineEdit:"},
					LineEdit{Row: 0, Column: 1},
					ToolButton{Row: 0, Column: 2, Text: "..."},
					Label{Row: 1, Column: 0, Text: "Another LineEdit:"},
					LineEdit{Row: 1, Column: 1},
					Label{Row: 2, Column: 0, Text: "A ComboBox:"},
					ComboBox{Row: 2, Column: 1},
					VSpacer{Row: 3, Column: 0, Size: 10},
					Label{Row: 4, Column: 0, ColumnSpan: 2, Text: "A TextEdit:"},
					TextEdit{Row: 5, Column: 0, ColumnSpan: 2},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{AssignTo: &acceptPB, Text: "OK", OnClicked: func() { dlg.Accept() }},
					PushButton{AssignTo: &cancelPB, Text: "Cancel", OnClicked: func() { dlg.Cancel() }},
				},
			},
		},
	}.Create(mw)); err != nil {
		log.Fatal(err)
	}

	dlg.Run()
}

func main() {
	walk.Initialize(walk.InitParams{})
	defer walk.Shutdown()

	mw := new(MyMainWindow)

	openImage, err := walk.NewBitmapFromFile("../img/open.png")
	if err != nil {
		log.Fatal(err)
	}

	var openAction *walk.Action
	var recentMenu *walk.Menu

	menuActions, err := CreateMenuActions(
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Action{
					AssignTo:    &openAction,
					Text:        "&Open",
					Image:       openImage,
					OnTriggered: func() { mw.openAction_Triggered() },
				},
				Menu{
					AssignTo: &recentMenu,
					Text:     "Recent",
				},
				Action{},
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

	toolBarActions, err := CreateToolBarActions(
		ActionRef{openAction},
		Action{Text: "Show Dialog", OnTriggered: func() { mw.showDialogAction_Triggered() }})
	if err != nil {
		log.Fatal(err)
	}

	if err := (MainWindow{
		AssignTo:       &mw.MainWindow,
		Title:          "FTPS cycle finder",
		MenuActions:    menuActions,
		ToolBarActions: toolBarActions,
		MinSize:        Size{600, 400},
		Size:           Size{800, 600},
		Layout:         HBox{Margins: Margins{6, 6, 6, 6}},
		Children: []Widget{
			ToolBar{Orientation: Vertical, Actions: toolBarActions},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							Label{Text: "File"},
							LineEdit{ContextMenuActions: []*walk.Action{openAction}},
							ToolButton{Text: "..."},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							PushButton{Text: "Check"},
							PushButton{Text: "Check and Fix"},
							PushButton{Text: "Clear"},
							HSpacer{},
							Label{Text: "Parameter"},
							LineEdit{MaxLength: 10},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							LineEdit{Text: "Ready.", ReadOnly: true},
							ProgressBar{StretchFactor: 10},
						},
					},
					TextEdit{ReadOnly: true},
				},
			},
		},
	}.Create(nil)); err != nil {
		log.Fatal(err)
	}

	mw.Show()
	mw.Run()
}
