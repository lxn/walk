// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path"
	"strings"
)

import (
	"github.com/lxn/walk"
)

type MainWindow struct {
	*walk.MainWindow
	tabWidget    *walk.TabWidget
	prevFilePath string
}

func (mw *MainWindow) openImage() {
	dlg := &walk.FileDialog{}

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select an Image"

	if ok, _ := dlg.ShowOpen(mw); !ok {
		return
	}

	mw.prevFilePath = dlg.FilePath

	img, _ := walk.NewImageFromFile(dlg.FilePath)

	var succeeded bool
	defer func() {
		if !succeeded {
			img.Dispose()
		}
	}()

	page, _ := walk.NewTabPage()
	page.SetTitle(path.Base(strings.Replace(dlg.FilePath, "\\", "/", -1)))
	page.SetLayout(walk.NewHBoxLayout())

	defer func() {
		if !succeeded {
			page.Dispose()
		}
	}()

	imageView, _ := walk.NewImageView(page)

	defer func() {
		if !succeeded {
			imageView.Dispose()
		}
	}()

	imageView.SetImage(img)
	mw.tabWidget.Pages().Add(page)
	mw.tabWidget.SetCurrentIndex(mw.tabWidget.Pages().Len() - 1)

	succeeded = true
}

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	mainWnd, _ := walk.NewMainWindow()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.SetLayout(walk.NewVBoxLayout())
	mw.SetTitle("Walk Image Viewer Example")

	mw.tabWidget, _ = walk.NewTabWidget(mw)

	imageList, _ := walk.NewImageList(walk.Size{16, 16}, 0)
	mw.ToolBar().SetImageList(imageList)

	fileMenu, _ := walk.NewMenu()
	fileMenuAction, _ := mw.Menu().Actions().AddMenu(fileMenu)
	fileMenuAction.SetText("&File")

	openBmp, _ := walk.NewBitmapFromFile("../img/open.png")

	openAction := walk.NewAction()
	openAction.SetImage(openBmp)
	openAction.SetText("&Open")
	openAction.Triggered().Attach(func() { mw.openImage() })
	fileMenu.Actions().Add(openAction)
	mw.ToolBar().Actions().Add(openAction)

	exitAction := walk.NewAction()
	exitAction.SetText("E&xit")
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	fileMenu.Actions().Add(exitAction)

	helpMenu, _ := walk.NewMenu()
	helpMenuAction, _ := mw.Menu().Actions().AddMenu(helpMenu)
	helpMenuAction.SetText("&Help")

	aboutAction := walk.NewAction()
	aboutAction.SetText("&About")
	aboutAction.Triggered().Attach(func() {
		walk.MsgBox(mw, "About", "Walk Image Viewer Example", walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	mw.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	mw.SetSize(walk.Size{800, 600})
	mw.Show()

	mw.Run()
}
