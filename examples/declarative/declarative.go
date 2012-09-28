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

func (mw *MyMainWindow) checkPushButton_Clicked() {
	walk.MsgBox(mw, "Check", "Doing some checking now...", walk.MsgBoxIconInformation|walk.MsgBoxOK)
}

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	mw := new(MyMainWindow)

	MainWindow{
		Widget: &mw.MainWindow,
		Title:  "FTPS cycle finder",
		Layout: VBox{
			Margins: &Margins{6, 6, 6, 6},
		},
		Children: []Widget{
			Composite{
				Layout: HBox{Margins: new(Margins)},
				Children: []Widget{
					Label{Text: "File"},
					LineEdit{},
					ToolButton{Text: "..."},
				},
			},
			Composite{
				Layout: HBox{Margins: new(Margins)},
				Children: []Widget{
					PushButton{Text: "Check", OnClicked: func() { mw.checkPushButton_Clicked() }},
					PushButton{Text: "Check and Fix"},
					PushButton{Text: "Clear"},
					HSpacer{},
					Label{Text: "Parameter"},
					LineEdit{MaxLength: 10},
				},
			},
			Composite{
				Layout: HBox{Margins: new(Margins)},
				Children: []Widget{
					LineEdit{Text: "Ready.", ReadOnly: true},
					ProgressBar{HStretch: 10},
				},
			},
			TextEdit{ReadOnly: true},
		},
	}.Create(nil)

	mw.Show()
	mw.Run()
}
