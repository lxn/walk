// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

import "walk"

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
	parent.Children().Clear()

	dirPath := pathForTreeViewItem(parent)

	dir, err := os.Open(dirPath)
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

			parent.Children().Add(child)
		}
	}
}

func (mw *MainWindow) populateListView(dirPath string) {
	mw.listView.SetSuspended(true)
	defer mw.listView.SetSuspended(false)

	mw.listView.Items().Clear()

	dir, err := os.Open(dirPath)
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

			mw.listView.Items().Add(item)
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
	item.Children().Add(walk.NewTreeViewItem())

	return item
}

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	mainWnd, _ := walk.NewMainWindow()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.SetTitle("Walk File Browser Example")
	mw.SetLayout(walk.NewHBoxLayout())

	fileMenu, _ := walk.NewMenu()
	fileMenuAction, _ := mw.Menu().Actions().AddMenu(fileMenu)
	fileMenuAction.SetText("&File")

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
		walk.MsgBox(mw, "About", "Walk File Browser Example", walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	splitter, _ := walk.NewSplitter(mw)

	mw.treeView, _ = walk.NewTreeView(splitter)

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

	drives, _ := walk.DriveNames()

	mw.treeView.SetSuspended(true)
	for _, drive := range drives {
		driveItem := newTreeViewItem(drive[:2])
		mw.treeView.Items().Add(driveItem)
	}
	mw.treeView.SetSuspended(false)

	mw.listView, _ = walk.NewListView(splitter)
	mw.listView.SetSingleItemSelection(true)

	mw.listView.CurrentIndexChanged().Attach(func() {
		index := mw.listView.CurrentIndex()
		var url string
		if index > -1 {
			item := mw.listView.Items().At(index)

			url = path.Join(pathForTreeViewItem(mw.selTvwItem), item.Texts()[0])
		}

		mw.preview.SetURL(url)
	})

	nameCol := walk.NewListViewColumn()
	nameCol.SetTitle("Name")
	nameCol.SetWidth(200)
	mw.listView.Columns().Add(nameCol)

	sizeCol := walk.NewListViewColumn()
	sizeCol.SetTitle("Size")
	sizeCol.SetWidth(80)
	sizeCol.SetAlignment(walk.AlignFar)
	mw.listView.Columns().Add(sizeCol)

	modCol := walk.NewListViewColumn()
	modCol.SetTitle("Modified")
	modCol.SetWidth(120)
	mw.listView.Columns().Add(modCol)

	mw.preview, _ = walk.NewWebView(splitter)

	mw.SetMinMaxSize(walk.Size{600, 400}, walk.Size{})
	mw.SetSize(walk.Size{800, 600})
	mw.Show()

	mw.Run()
}
