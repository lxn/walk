// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path"
	"runtime"
	"strings"
)

import (
	"walk"
)

type MainWindow struct {
	*walk.MainWindow
	tabWidget    *walk.TabWidget
	prevFilePath string
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func (mw *MainWindow) openImage() {
	dlg := &walk.FileDialog{}

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select an Image"

	ok, err := dlg.ShowOpen(mw)
	panicIfErr(err)
	if !ok {
		return
	}

	mw.prevFilePath = dlg.FilePath

	img, err := walk.NewImageFromFile(dlg.FilePath)
	panicIfErr(err)

	var succeeded bool
	defer func() {
		if !succeeded {
			img.Dispose()
		}
	}()

	page, err := walk.NewTabPage()
	panicIfErr(err)
	panicIfErr(page.SetTitle(path.Base(strings.Replace(dlg.FilePath, "\\", "/", -1))))
	panicIfErr(page.SetLayout(walk.NewHBoxLayout()))

	defer func() {
		if !succeeded {
			page.Dispose()
		}
	}()

	imageView, err := walk.NewImageView(page)
	panicIfErr(err)

	defer func() {
		if !succeeded {
			imageView.Dispose()
		}
	}()

	panicIfErr(imageView.SetImage(img))
	panicIfErr(mw.tabWidget.Pages().Add(page))
	panicIfErr(mw.tabWidget.SetCurrentIndex(mw.tabWidget.Pages().Len() - 1))

	succeeded = true
}

func main() {
	runtime.LockOSThread()

	mainWnd, err := walk.NewMainWindow()
	panicIfErr(err)

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.ClientArea().SetLayout(walk.NewVBoxLayout()))
	panicIfErr(mw.SetTitle("Walk Image Viewer Example"))

	mw.tabWidget, err = walk.NewTabWidget(mw.ClientArea())
	panicIfErr(err)

	imageList, err := walk.NewImageList(walk.Size{16, 16}, 0)
	panicIfErr(err)
	mw.ToolBar().SetImageList(imageList)

	fileMenu, err := walk.NewMenu()
	panicIfErr(err)
	fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	panicIfErr(fileMenuAction.SetText("File"))

	openBmp, err := walk.NewBitmapFromFile("../img/open.png")
	panicIfErr(err)

	openAction := walk.NewAction()
	openAction.SetImage(openBmp)
	panicIfErr(openAction.SetText("Open"))
	openAction.Triggered().Attach(func() { mw.openImage() })
	panicIfErr(fileMenu.Actions().Add(openAction))
	panicIfErr(mw.ToolBar().Actions().Add(openAction))

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
		walk.MsgBox(mw, "About", "Walk Image Viewer Example", walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	panicIfErr(helpMenu.Actions().Add(aboutAction))

	panicIfErr(mw.SetMinMaxSize(walk.Size{320, 240}, walk.Size{}))
	panicIfErr(mw.SetSize(walk.Size{800, 600}))
	mw.Show()

	os.Exit(mw.Run())
}
