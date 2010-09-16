// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

type TreeView struct {
	Widget
	items *TreeViewItemList
}

func NewTreeView(parent IContainer) (*TreeView, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("SysTreeView32"), nil,
		TVS_HASBUTTONS|TVS_HASLINES|TVS_LINESATROOT|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	tv := &TreeView{Widget: Widget{hWnd: hWnd, parent: parent}}

	tv.items = newTreeViewItemList(tv)

	tv.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = tv

	parent.Children().Add(tv)

	return tv, nil
}

func (*TreeView) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz | ShrinkVert | GrowVert
}

func (tv *TreeView) PreferredSize() drawing.Size {
	return tv.dialogBaseUnitsToPixels(drawing.Size{100, 100})
}

func (tv *TreeView) Items() *TreeViewItemList {
	return tv.items
}

func (tv *TreeView) onInsertingTreeViewItem(parent *TreeViewItem, index int, item *TreeViewItem) (err os.Error) {
	var tvi TVITEM
	var tvins TVINSERTSTRUCT

	tvi.Mask = TVIF_TEXT
	tvi.PszText = syscall.StringToUTF16Ptr(item.text)

	tvins.Item = tvi

	if parent == nil {
		tvins.HParent = TVI_ROOT
	} else {
		tvins.HParent = parent.handle
	}

	if index == 0 {
		tvins.HInsertAfter = TVI_LAST
	} else {
		var items *TreeViewItemList
		if parent == nil {
			items = tv.items
		} else {
			items = parent.children
		}
		tvins.HInsertAfter = items.At(index - 1).handle
	}

	item.handle = HTREEITEM(SendMessage(tv.hWnd, TVM_INSERTITEM, 0, uintptr(unsafe.Pointer(&tvins))))
	if item.handle == 0 {
		err = newError("TVM_INSERTITEM failed")
	} else {
		item.children.observer = tv
	}

	return
}

func (tv *TreeView) onRemovingTreeViewItem(index int, item *TreeViewItem) (err os.Error) {
	panic("not implemented")
}

func (tv *TreeView) onClearingTreeViewItems() (err os.Error) {
	panic("not implemented")
}
