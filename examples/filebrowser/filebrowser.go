// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

import (
	"walk/drawing"
	"walk/gui"
	wpath "walk/path"
)

type MainWindow struct {
	*gui.MainWindow
	treeView   *gui.TreeView
	selTvwItem *gui.TreeViewItem
	listView   *gui.ListView
	preview    *gui.WebView
}

func (mw *MainWindow) showError(err os.Error) {
	gui.MsgBox(mw, "Error", err.String(), gui.MsgBoxOK|gui.MsgBoxIconError)
}

func (mw *MainWindow) populateTreeViewItem(parent *gui.TreeViewItem) {
	mw.treeView.BeginUpdate()
	defer mw.treeView.EndUpdate()

	// Remove dummy child
	parent.Children().Clear()

	dirPath := pathForTreeViewItem(parent)

	dir, err := os.Open(dirPath, os.O_RDONLY, 0)
	if err != nil {
		mw.showError(err)
		return
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	panicIfErr(err)

	for _, name := range names {
		fi, err := os.Stat(path.Join(dirPath, name))
		panicIfErr(err)

		if !excludePath(name) && fi.IsDirectory() {
			child := newTreeViewItem(name)

			_, err = parent.Children().Add(child)
			panicIfErr(err)
		}
	}
}

func (mw *MainWindow) populateListView(dirPath string) {
	mw.listView.BeginUpdate()
	defer mw.listView.EndUpdate()

	mw.listView.Items().Clear()

	dir, err := os.Open(dirPath, os.O_RDONLY, 0)
	if err != nil {
		mw.showError(err)
		return
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	panicIfErr(err)

	for _, name := range names {
		if !excludePath(name) {
			fullPath := path.Join(dirPath, name)

			fi, err := os.Stat(fullPath)
			if err != nil {
				mw.showError(err)
				continue
			}

			var size string
			if !fi.IsDirectory() {
				size = fmt.Sprintf("%d", fi.Size)
			}
			lastMod := time.SecondsToLocalTime(fi.Mtime_ns / 10e8).Format("2006-01-02 15:04:05")

			item := gui.NewListViewItem()
			texts := []string{name, size, lastMod}
			item.SetTexts(texts)

			_, err = mw.listView.Items().Add(item)
			panicIfErr(err)
		}
	}
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func pathForTreeViewItem(item *gui.TreeViewItem) string {
	var parts []string
	for item != nil {
		parts = append([]string{item.Text()}, parts...)
		item = item.Parent()
	}

	return path.Join(parts...)
}

func excludePath(path string) bool {
	if path == "System Volume Information" {
		return true
	}

	return false
}

func newTreeViewItem(text string) *gui.TreeViewItem {
	item := gui.NewTreeViewItem()
	item.SetText(text)

	// For now, we add a dummy child to make the item expandable.
	_, err := item.Children().Add(gui.NewTreeViewItem())
	panicIfErr(err)

	return item
}

func runMainWindow() (int, os.Error) {
	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)
	defer mainWnd.Dispose()

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetText("Walk File Browser Example"))
	mw.ClientArea().SetLayout(gui.NewHBoxLayout())

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
		gui.MsgBox(mw, "About", "Walk File Browser Example", gui.MsgBoxOK|gui.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	splitter, err := gui.NewSplitter(mw.ClientArea())
	panicIfErr(err)

	mw.treeView, err = gui.NewTreeView(splitter)
	panicIfErr(err)
	//	panicIfErr(mw.treeView.SetMaxSize(drawing.Size{200, 0}))

	mw.treeView.AddItemExpandedHandler(func(args gui.TreeViewItemEventArgs) {
		item := args.Item()
		children := item.Children()
		if children.Len() == 1 && children.At(0).Text() == "" {
			mw.populateTreeViewItem(item)
		}
	})

	mw.treeView.AddSelectionChangedHandler(func(args gui.TreeViewItemSelectionEventArgs) {
		mw.selTvwItem = args.New()
		mw.populateListView(pathForTreeViewItem(mw.selTvwItem))
	})

	drives, err := wpath.DriveNames()
	panicIfErr(err)

	mw.treeView.BeginUpdate()
	for _, drive := range drives {
		driveItem := newTreeViewItem(drive[:2])
		_, err = mw.treeView.Items().Add(driveItem)
		panicIfErr(err)
	}
	mw.treeView.EndUpdate()

	mw.listView, err = gui.NewListView(splitter)
	panicIfErr(err)

	mw.listView.AddSelectedIndexChangedHandler(func(args gui.EventArgs) {
		index := mw.listView.SelectedIndex()
		var url string
		if index > -1 {
			item := mw.listView.Items().At(index)
			panicIfErr(err)

			url = path.Join(pathForTreeViewItem(mw.selTvwItem), item.Texts()[0])
		}

		err := mw.preview.SetURL(url)
		panicIfErr(err)
	})

	nameCol := gui.NewListViewColumn()
	nameCol.SetTitle("Name")
	nameCol.SetWidth(300)
	_, err = mw.listView.Columns().Add(nameCol)
	panicIfErr(err)

	sizeCol := gui.NewListViewColumn()
	sizeCol.SetTitle("Size")
	sizeCol.SetAlignment(gui.RightAlignment)
	_, err = mw.listView.Columns().Add(sizeCol)
	panicIfErr(err)

	lastModCol := gui.NewListViewColumn()
	lastModCol.SetTitle("Last Modified")
	lastModCol.SetWidth(120)
	_, err = mw.listView.Columns().Add(lastModCol)
	panicIfErr(err)

	mw.preview, err = gui.NewWebView(splitter)
	panicIfErr(err)

	panicIfErr(mw.SetMinSize(drawing.Size{600, 400}))
	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	return mw.RunMessageLoop()
}

func main() {
	runtime.LockOSThread()

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Error:", x)
		}
	}()

	exitCode, err := runMainWindow()
	panicIfErr(err)
	os.Exit(exitCode)
}
