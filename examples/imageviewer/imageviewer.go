// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"runtime"
)

import (
	"walk/drawing"
	"walk/gui"
)

type MainWindow struct {
	*gui.MainWindow
	imageView    *gui.ImageView
	prevFilePath string
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func (mw *MainWindow) updateTitle(filePath string) {
	title := "Walk Image Viewer Example"

	if filePath != "" {
		title += fmt.Sprintf(" [%s]", filePath)
	}

	panicIfErr(mw.SetText(title))
}

func (mw *MainWindow) openBitmap() {
	dlg := &gui.FileDialog{}

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Bitmap Files (*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select a Bitmap"

	ok, err := dlg.ShowOpen(mw)
	panicIfErr(err)
	if !ok {
		return
	}

	mw.prevFilePath = dlg.FilePath

	bmp, err := drawing.NewBitmapFromFile(dlg.FilePath)
	panicIfErr(err)

	prevImage := mw.imageView.Image()
	if prevImage != nil {
		prevImage.Dispose()
	}

	panicIfErr(mw.imageView.SetImage(bmp))

	mw.updateTitle(dlg.FilePath)
}

func runMainWindow() {
	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)
	defer mainWnd.Dispose()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.ClientArea().SetLayout(gui.NewVBoxLayout())
	mw.updateTitle("")

	mw.imageView, err = gui.NewImageView(mw.ClientArea())
	panicIfErr(err)

	imageList, err := gui.NewImageList(drawing.Size{16, 16}, 0)
	panicIfErr(err)
	mw.ToolBar().SetImageList(imageList)

	fileMenu, err := gui.NewMenu()
	panicIfErr(err)
	_, fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	fileMenuAction.SetText("File")

	openBmp, err := drawing.NewBitmapFromFile("img/open.png")
	panicIfErr(err)

	openAction := gui.NewAction()
	openAction.SetImage(openBmp)
	openAction.SetText("Open")
	openAction.AddTriggeredHandler(func(args gui.EventArgs) { mw.openBitmap() })
	fileMenu.Actions().Add(openAction)
	mw.ToolBar().Actions().Add(openAction)

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
		gui.MsgBox(mw, "About", "Walk Image Viewer Example", gui.MsgBoxOK|gui.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	panicIfErr(mw.RunMessageLoop())
}

func main() {
	runtime.LockOSThread()

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Error:", x)
		}
	}()

	runMainWindow()
}
