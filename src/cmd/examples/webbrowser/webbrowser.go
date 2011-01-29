// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"runtime"
)

import (
	"walk/winapi/user32"
)

import (
	"walk/drawing"
	"walk/gui"
)

type MainWindow struct {
	*gui.MainWindow
	urlLineEdit *gui.LineEdit
	webView     *gui.WebView
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.LockOSThread()

	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetText("Walk Web Browser Example"))
	panicIfErr(mw.ClientArea().SetLayout(gui.NewVBoxLayout()))

	fileMenu, err := gui.NewMenu()
	panicIfErr(err)
	fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	panicIfErr(fileMenuAction.SetText("File"))

	exitAction := gui.NewAction()
	panicIfErr(exitAction.SetText("Exit"))
	exitAction.Triggered().Subscribe(func(args *gui.EventArgs) { gui.Exit(0) })
	panicIfErr(fileMenu.Actions().Add(exitAction))

	helpMenu, err := gui.NewMenu()
	panicIfErr(err)
	helpMenuAction, err := mw.Menu().Actions().AddMenu(helpMenu)
	panicIfErr(err)
	panicIfErr(helpMenuAction.SetText("Help"))

	aboutAction := gui.NewAction()
	panicIfErr(aboutAction.SetText("About"))
	aboutAction.Triggered().Subscribe(func(args *gui.EventArgs) {
		gui.MsgBox(mw, "About", "Walk Web Browser Example", gui.MsgBoxOK|gui.MsgBoxIconInformation)
	})
	panicIfErr(helpMenu.Actions().Add(aboutAction))

	mw.urlLineEdit, err = gui.NewLineEdit(mw.ClientArea())
	panicIfErr(err)
	mw.urlLineEdit.KeyDown().Subscribe(func(args *gui.KeyEventArgs) {
		if args.Key() == user32.VK_RETURN {
			panicIfErr(mw.webView.SetURL(mw.urlLineEdit.Text()))
		}
	})

	mw.webView, err = gui.NewWebView(mw.ClientArea())
	panicIfErr(err)

	panicIfErr(mw.webView.SetURL("http://golang.org"))

	panicIfErr(mw.SetMinSize(drawing.Size{600, 400}))
	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	os.Exit(mw.Run())
}
