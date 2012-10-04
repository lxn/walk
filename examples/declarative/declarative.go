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

type DialogDecl struct {
	Owner    walk.RootWidget
	Dialog   **walk.Dialog
	AcceptPB **walk.PushButton
	Widgets  []Widget
	Title    string
	Size     Size
	MinSize  Size
}

func (dd *DialogDecl) Create() error {
	var cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      dd.Dialog,
		Title:         dd.Title,
		DefaultButton: dd.AcceptPB,
		CancelButton:  &cancelPB,
		MinSize:       dd.MinSize,
		Size:          dd.Size,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout:   Grid{},
				Children: dd.Widgets,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{AssignTo: dd.AcceptPB, Text: "OK"},
					PushButton{AssignTo: &cancelPB, Text: "Cancel", OnClicked: func() { dd.Dialog.Cancel() }},
				},
			},
		},
	}.Create(dd.Owner)
}

func (mw *MyMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Nothing to see here...", walk.MsgBoxIconInformation|walk.MsgBoxOK)
}

func (mw *MyMainWindow) showDialogAction_Triggered() {
	dlg := new(MyDialog)

	var acceptPB *walk.PushButton
	var le1, le2 *walk.LineEdit

	widgets := []Widget{
		Label{Row: 0, Column: 0, Text: "A LineEdit:"},
		LineEdit{Row: 0, Column: 1, AssignTo: &le1, OnTextChanged: func() { le2.SetText(le1.Text()) }},
		ToolButton{Row: 0, Column: 2, Text: "..."},
		Label{Row: 1, Column: 0, Text: "Another LineEdit:"},
		LineEdit{Row: 1, Column: 1, AssignTo: &le2},
		Label{Row: 2, Column: 0, Text: "A ComboBox:"},
		ComboBox{Row: 2, Column: 1},
		VSpacer{Row: 3, Column: 0, Size: 10},
		Label{Row: 4, Column: 0, ColumnSpan: 2, Text: "A TextEdit:"},
		TextEdit{Row: 5, Column: 0, ColumnSpan: 2},
	}

	dd := &DialogDecl{
		Title:    "My Dialog",
		Owner:    mw,
		Dialog:   &dlg.Dialog,
		AcceptPB: &acceptPB,
		Widgets:  widgets,
		MinSize:  Size{400, 300},
	}

	if err := dd.Create(); err != nil {
		log.Fatal(err)
	}

	acceptPB.Clicked().Attach(func() {
		dlg.Accept()
	})

	dlg.Run()
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
