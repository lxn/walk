// Copyright 2013 The Walk Authors. All rights reserved.
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

func main() {
	mw := new(MyMainWindow)

	var openAction, showAboutBoxAction *walk.Action
	var recentMenu *walk.Menu

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Walk Actions Example",
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open",
						Image:       "../img/open.png",
						Enabled:     Bind("openEnabledCB.Checked"),
						Visible:     Bind("openVisibleCB.Checked"),
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
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "About",
						OnTriggered: func() { mw.showAboutBoxAction_Triggered() },
					},
				},
			},
		},
		ToolBarItems: []MenuItem{
			ActionRef{&openAction},
			Separator{},
			ActionRef{&showAboutBoxAction},
		},
		ContextMenuItems: []MenuItem{
			ActionRef{&showAboutBoxAction},
		},
		MinSize: Size{300, 200},
		Layout:  VBox{},
		Children: []Widget{
			CheckBox{
				Name:    "openEnabledCB",
				Text:    "Open Enabled",
				Checked: true,
			},
			CheckBox{
				Name:    "openVisibleCB",
				Text:    "Open Visible",
				Checked: true,
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

func (mw *MyMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "About", "Walk Actions Example", walk.MsgBoxIconInformation)
}
