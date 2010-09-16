// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
	. "walk/winapi/comctl32"
)

type TreeViewItem struct {
	handle   HTREEITEM
	children *TreeViewItemList
	parent   *TreeViewItem
	text     string
}

func NewTreeViewItem() *TreeViewItem {
	tvi := &TreeViewItem{}

	tvi.children = newTreeViewItemList(nil)
	tvi.children.parent = tvi

	return tvi
}

func (tvi *TreeViewItem) Children() *TreeViewItemList {
	return tvi.children
}

func (tvi *TreeViewItem) Parent() *TreeViewItem {
	return tvi.parent
}

func (tvi *TreeViewItem) Text() string {
	return tvi.text
}

func (tvi *TreeViewItem) SetText(value string) os.Error {
	tvi.text = value

	return nil
}
