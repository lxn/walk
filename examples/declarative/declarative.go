// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
}

func (mw *MyMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Nothing to see here...", walk.MsgBoxIconInformation|walk.MsgBoxOK)
}

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	mw := new(MyMainWindow)

	openImage, _ := walk.NewBitmapFromFile("../img/open.png")

	var openAction *walk.Action

	menuActions, _ := CreateMenuActions(
		[]MenuItem{
			SubMenu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open",
						Image:       openImage,
						OnTriggered: func() { mw.openAction_Triggered() },
					},
					Action{},
					Action{
						Text:        "E&xit",
						OnTriggered: func() { walk.App().Exit(0) },
					},
				},
			},
		})

	toolBarActions, _ := CreateToolBarActions(
		[]ToolBarItem{
			ActionRef{openAction},
			Action{Text: "NOP"},
		},
	)

	marg0 := &Margins{}

	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "FTPS cycle finder",
		Menu:     Menu{Actions: menuActions},
		ToolBar:  ToolBar{Actions: toolBarActions},
		Layout:   HBox{Margins: &Margins{6, 6, 6, 6}},
		Children: []Widget{
			ToolBar{Orientation: walk.Vertical, Actions: toolBarActions},
			Composite{
				Layout: VBox{Margins: marg0},
				Children: []Widget{
					Composite{
						Layout: HBox{Margins: marg0},
						Children: []Widget{
							Label{Text: "File"},
							LineEdit{ContextMenu: Menu{Actions: []*walk.Action{openAction}}},
							ToolButton{Text: "..."},
						},
					},
					Composite{
						Layout: HBox{Margins: marg0},
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
						Layout: HBox{Margins: marg0},
						Children: []Widget{
							LineEdit{Text: "Ready.", ReadOnly: true},
							ProgressBar{StretchFactor: 10},
						},
					},
					TextEdit{ReadOnly: true},
				},
			},
		},
	}.Create(nil)

	mw.Show()
	mw.Run()
}
