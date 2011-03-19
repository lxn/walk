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
	"walk"
)

type MainWindow struct {
	*walk.MainWindow
	treeView   *walk.TreeView
	selTvwItem *walk.TreeViewItem
	listView   *walk.ListView
	preview    *walk.WebView
}

func (mw *MainWindow) showError(err os.Error) {
	walk.MsgBox(mw, "Error", err.String(), walk.MsgBoxOK|walk.MsgBoxIconError)
}

func (mw *MainWindow) populateTreeViewItem(parent *walk.TreeViewItem) {
	mw.treeView.SetSuspended(true)
	defer mw.treeView.SetSuspended(false)

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
	mw.listView.SetSuspended(true)
	defer mw.listView.SetSuspended(false)

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

			item := walk.NewListViewItem()
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

func pathForTreeViewItem(item *walk.TreeViewItem) string {
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

func newTreeViewItem(text string) *walk.TreeViewItem {
	item := walk.NewTreeViewItem()
	item.SetText(text)

	// For now, we add a dummy child to make the item expandable.
	panicIfErr(item.Children().Add(walk.NewTreeViewItem()))

	return item
}

func main() {
	runtime.LockOSThread()

	mainWnd, err := walk.NewMainWindow()
	panicIfErr(err)

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetTitle("Walk File Browser Example"))
	panicIfErr(mw.ClientArea().SetLayout(walk.NewHBoxLayout()))

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
		walk.MsgBox(mw, "About", "Walk File Browser Example", walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	panicIfErr(helpMenu.Actions().Add(aboutAction))

	splitter, err := walk.NewSplitter(mw.ClientArea())
	panicIfErr(err)

	mw.treeView, err = walk.NewTreeView(splitter)
	panicIfErr(err)
	panicIfErr(mw.treeView.SetMinMaxSize(walk.Size{}, walk.Size{200, 0}))

	mw.treeView.ItemExpanded().Attach(func(item *walk.TreeViewItem) {
		children := item.Children()
		if children.Len() == 1 && children.At(0).Text() == "" {
			mw.populateTreeViewItem(item)
		}
	})

	mw.treeView.SelectionChanged().Attach(func(old, new *walk.TreeViewItem) {
		mw.selTvwItem = new
		mw.populateListView(pathForTreeViewItem(new))
	})

	drives, err := walk.DriveNames()
	panicIfErr(err)

	mw.treeView.SetSuspended(true)
	for _, drive := range drives {
		driveItem := newTreeViewItem(drive[:2])
		panicIfErr(mw.treeView.Items().Add(driveItem))
	}
	mw.treeView.SetSuspended(false)

	mw.listView, err = walk.NewListView(splitter)
	panicIfErr(err)
	panicIfErr(mw.listView.SetSingleItemSelection(true))
	panicIfErr(mw.listView.SetMinMaxSize(walk.Size{}, walk.Size{422, 0}))

	mw.listView.CurrentIndexChanged().Attach(func() {
		index := mw.listView.CurrentIndex()
		var url string
		if index > -1 {
			item := mw.listView.Items().At(index)
			panicIfErr(err)

			url = path.Join(pathForTreeViewItem(mw.selTvwItem), item.Texts()[0])
		}

		err := mw.preview.SetURL(url)
		panicIfErr(err)
	})

	nameCol := walk.NewListViewColumn()
	nameCol.SetTitle("Name")
	nameCol.SetWidth(200)
	panicIfErr(mw.listView.Columns().Add(nameCol))

	sizeCol := walk.NewListViewColumn()
	sizeCol.SetTitle("Size")
	sizeCol.SetWidth(80)
	sizeCol.SetAlignment(walk.AlignFar)
	panicIfErr(mw.listView.Columns().Add(sizeCol))

	modCol := walk.NewListViewColumn()
	modCol.SetTitle("Modified")
	modCol.SetWidth(120)
	panicIfErr(mw.listView.Columns().Add(modCol))

	mw.preview, err = walk.NewWebView(splitter)
	panicIfErr(err)

	panicIfErr(mw.SetMinMaxSize(walk.Size{600, 400}, walk.Size{}))
	panicIfErr(mw.SetSize(walk.Size{800, 600}))
	mw.Show()

	os.Exit(mw.Run())
}
