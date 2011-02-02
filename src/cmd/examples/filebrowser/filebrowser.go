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
	panicIfErr(parent.Children().Clear())

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

			panicIfErr(parent.Children().Add(child))
		}
	}
}

func (mw *MainWindow) populateListView(dirPath string) {
	mw.listView.BeginUpdate()
	defer mw.listView.EndUpdate()

	panicIfErr(mw.listView.Items().Clear())

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

			panicIfErr(mw.listView.Items().Add(item))
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
	panicIfErr(item.Children().Add(gui.NewTreeViewItem()))

	return item
}

func main() {
	runtime.LockOSThread()

	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetText("Walk File Browser Example"))
	panicIfErr(mw.ClientArea().SetLayout(gui.NewHBoxLayout()))

	fileMenu, err := gui.NewMenu()
	panicIfErr(err)
	fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	panicIfErr(fileMenuAction.SetText("File"))

	exitAction := gui.NewAction()
	panicIfErr(exitAction.SetText("Exit"))
	exitAction.Triggered().Subscribe(func(args *gui.EventArgs) { gui.App().Exit(0) })
	panicIfErr(fileMenu.Actions().Add(exitAction))

	helpMenu, err := gui.NewMenu()
	panicIfErr(err)
	helpMenuAction, err := mw.Menu().Actions().AddMenu(helpMenu)
	panicIfErr(err)
	panicIfErr(helpMenuAction.SetText("Help"))

	aboutAction := gui.NewAction()
	panicIfErr(aboutAction.SetText("About"))
	aboutAction.Triggered().Subscribe(func(args *gui.EventArgs) {
		gui.MsgBox(mw, "About", "Walk File Browser Example", gui.MsgBoxOK|gui.MsgBoxIconInformation)
	})
	panicIfErr(helpMenu.Actions().Add(aboutAction))

	splitter, err := gui.NewSplitter(mw.ClientArea())
	panicIfErr(err)

	mw.treeView, err = gui.NewTreeView(splitter)
	panicIfErr(err)
	panicIfErr(mw.treeView.SetMaxSize(drawing.Size{200, 0}))

	mw.treeView.ItemExpanded().Subscribe(func(args *gui.TreeViewItemEventArgs) {
		item := args.Item()
		children := item.Children()
		if children.Len() == 1 && children.At(0).Text() == "" {
			mw.populateTreeViewItem(item)
		}
	})

	mw.treeView.SelectionChanged().Subscribe(func(args *gui.TreeViewItemSelectionEventArgs) {
		mw.selTvwItem = args.New()
		mw.populateListView(pathForTreeViewItem(mw.selTvwItem))
	})

	drives, err := wpath.DriveNames()
	panicIfErr(err)

	mw.treeView.BeginUpdate()
	for _, drive := range drives {
		driveItem := newTreeViewItem(drive[:2])
		panicIfErr(mw.treeView.Items().Add(driveItem))
	}
	mw.treeView.EndUpdate()

	mw.listView, err = gui.NewListView(splitter)
	panicIfErr(err)
	panicIfErr(mw.listView.SetMaxSize(drawing.Size{422, 0}))

	mw.listView.SelectedIndexChanged().Subscribe(func(args *gui.EventArgs) {
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
	nameCol.SetWidth(200)
	panicIfErr(mw.listView.Columns().Add(nameCol))

	sizeCol := gui.NewListViewColumn()
	sizeCol.SetTitle("Size")
	sizeCol.SetWidth(80)
	sizeCol.SetAlignment(gui.RightAlignment)
	panicIfErr(mw.listView.Columns().Add(sizeCol))

	modCol := gui.NewListViewColumn()
	modCol.SetTitle("Modified")
	modCol.SetWidth(120)
	panicIfErr(mw.listView.Columns().Add(modCol))

	mw.preview, err = gui.NewWebView(splitter)
	panicIfErr(err)

	panicIfErr(mw.SetMinSize(drawing.Size{600, 400}))
	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	os.Exit(mw.Run())
}
