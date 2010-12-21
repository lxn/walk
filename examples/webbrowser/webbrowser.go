// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"runtime"
)

import (
	"walk/drawing"
	"walk/gui"
)

type MainWindow struct {
	*gui.MainWindow
	webView *gui.WebView
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func runMainWindow() (int, os.Error) {
	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)
	defer mainWnd.Dispose()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.ClientArea().SetLayout(gui.NewVBoxLayout())

	fileMenu, err := gui.NewMenu()
	panicIfErr(err)
	_, fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	fileMenuAction.SetText("File")

	exitAction := gui.NewAction()
	exitAction.SetText("Exit")
	exitAction.AddTriggeredHandler(func(args gui.EventArgs) { gui.Exit(0) })
	fileMenu.Actions().Add(exitAction)

	helpMenu, err := gui.NewMenu()
	panicIfErr(err)
	_, helpMenuAction, err := mw.Menu().Actions().AddMenu(helpMenu)
	panicIfErr(err)
	helpMenuAction.SetText("Help")

	aboutAction := gui.NewAction()
	aboutAction.SetText("About")
	aboutAction.AddTriggeredHandler(func(args gui.EventArgs) {
		gui.MsgBox(mw, "About", "Walk Web Browser Example", gui.MsgBoxOK|gui.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	mw.webView, err = gui.NewWebView(mw.ClientArea())
	panicIfErr(err)

	panicIfErr(mw.webView.SetURL("http://golang.org"))

	panicIfErr(mw.SetMinSize(drawing.Size{600, 400}))
	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	return mw.RunMessageLoop()
}

func main() {
	runtime.LockOSThread()

	defer func() {
		if x := recover(); x != nil {
			log.Println("Error:", x)
		}
	}()

	exitCode, err := runMainWindow()
	panicIfErr(err)
	os.Exit(exitCode)
}
