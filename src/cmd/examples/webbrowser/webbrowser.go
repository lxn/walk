// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"runtime"
)

import (
	"walk"
)

type MainWindow struct {
	*walk.MainWindow
	urlLineEdit *walk.LineEdit
	webView     *walk.WebView
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.LockOSThread()

	mainWnd, err := walk.NewMainWindow()
	panicIfErr(err)

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetText("Walk Web Browser Example"))
	panicIfErr(mw.ClientArea().SetLayout(walk.NewVBoxLayout()))

	fileMenu, err := walk.NewMenu()
	panicIfErr(err)
	fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	panicIfErr(fileMenuAction.SetText("File"))

	exitAction := walk.NewAction()
	panicIfErr(exitAction.SetText("Exit"))
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	panicIfErr(fileMenu.Actions().Add(exitAction))

	helpMenu, err := walk.NewMenu()
	panicIfErr(err)
	helpMenuAction, err := mw.Menu().Actions().AddMenu(helpMenu)
	panicIfErr(err)
	panicIfErr(helpMenuAction.SetText("Help"))

	aboutAction := walk.NewAction()
	panicIfErr(aboutAction.SetText("About"))
	aboutAction.Triggered().Attach(func() {
		walk.MsgBox(mw, "About", "Walk Web Browser Example", walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	panicIfErr(helpMenu.Actions().Add(aboutAction))

	mw.urlLineEdit, err = walk.NewLineEdit(mw.ClientArea())
	panicIfErr(err)
	mw.urlLineEdit.ReturnPressed().Attach(func() {
		panicIfErr(mw.webView.SetURL(mw.urlLineEdit.Text()))
	})

	mw.webView, err = walk.NewWebView(mw.ClientArea())
	panicIfErr(err)

	panicIfErr(mw.webView.SetURL("http://golang.org"))

	panicIfErr(mw.SetMinMaxSize(walk.Size{600, 400}, walk.Size{}))
	panicIfErr(mw.SetSize(walk.Size{800, 600}))
	mw.Show()

	os.Exit(mw.Run())
}
