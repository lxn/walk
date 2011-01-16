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

type TreeViewItemEventArgs interface {
	EventArgs
	Item() *TreeViewItem
}

type treeViewItemEventArgs struct {
	eventArgs
	item *TreeViewItem
}

func (a *treeViewItemEventArgs) Item() *TreeViewItem {
	return a.item
}

type TreeViewItemEventHandler func(args TreeViewItemEventArgs)

type TreeViewItemSelectionEventArgs interface {
	EventArgs
	Old() *TreeViewItem
	New() *TreeViewItem
}

type treeViewItemSelectionEventArgs struct {
	eventArgs
	old *TreeViewItem
	new *TreeViewItem
}

func (a *treeViewItemSelectionEventArgs) Old() *TreeViewItem {
	return a.old
}

func (a *treeViewItemSelectionEventArgs) New() *TreeViewItem {
	return a.new
}

type TreeViewItemSelectionEventHandler func(args TreeViewItemSelectionEventArgs)

type TreeView struct {
	Widget
	items                     *TreeViewItemList
	itemCollapsedHandlers     []TreeViewItemEventHandler
	itemCollapsingHandlers    []TreeViewItemEventHandler
	itemExpandedHandlers      []TreeViewItemEventHandler
	itemExpandingHandlers     []TreeViewItemEventHandler
	selectionChangedHandlers  []TreeViewItemSelectionEventHandler
	selectionChangingHandlers []TreeViewItemSelectionEventHandler
}

func NewTreeView(parent IContainer) (*TreeView, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("SysTreeView32"), nil,
		TVS_HASBUTTONS|TVS_HASLINES|TVS_LINESATROOT|TVS_SHOWSELALWAYS|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	tv := &TreeView{Widget: Widget{hWnd: hWnd, parent: parent}}

	if err := tv.setTheme("Explorer"); err != nil {
		return nil, err
	}

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

func (tv *TreeView) AddItemCollapsedHandler(handler TreeViewItemEventHandler) {
	tv.itemCollapsedHandlers = append(tv.itemCollapsedHandlers, handler)
}

func (tv *TreeView) RemoveItemCollapsedHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemCollapsedHandlers {
		if h == handler {
			tv.itemCollapsedHandlers = append(tv.itemCollapsedHandlers[:i], tv.itemCollapsedHandlers[i+1:]...)
			break
		}
	}
}

func (tv *TreeView) raiseItemCollapsed(item *TreeViewItem) {
	args := &treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item}
	for _, handler := range tv.itemCollapsedHandlers {
		handler(args)
	}
}

func (tv *TreeView) AddItemCollapsingHandler(handler TreeViewItemEventHandler) {
	tv.itemCollapsingHandlers = append(tv.itemCollapsingHandlers, handler)
}

func (tv *TreeView) RemoveItemCollapsingHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemCollapsingHandlers {
		if h == handler {
			tv.itemCollapsingHandlers = append(tv.itemCollapsingHandlers[:i], tv.itemCollapsingHandlers[i+1:]...)
			break
		}
	}
}

func (tv *TreeView) raiseItemCollapsing(item *TreeViewItem) {
	args := &treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item}
	for _, handler := range tv.itemCollapsingHandlers {
		handler(args)
	}
}

func (tv *TreeView) AddItemExpandedHandler(handler TreeViewItemEventHandler) {
	tv.itemExpandedHandlers = append(tv.itemExpandedHandlers, handler)
}

func (tv *TreeView) RemoveItemExpandedHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemExpandedHandlers {
		if h == handler {
			tv.itemExpandedHandlers = append(tv.itemExpandedHandlers[:i], tv.itemExpandedHandlers[i+1:]...)
			break
		}
	}
}

func (tv *TreeView) raiseItemExpanded(item *TreeViewItem) {
	args := &treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item}
	for _, handler := range tv.itemExpandedHandlers {
		handler(args)
	}
}

func (tv *TreeView) AddItemExpandingHandler(handler TreeViewItemEventHandler) {
	tv.itemExpandingHandlers = append(tv.itemExpandingHandlers, handler)
}

func (tv *TreeView) RemoveItemExpandingHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemExpandingHandlers {
		if h == handler {
			tv.itemExpandingHandlers = append(tv.itemExpandingHandlers[:i], tv.itemExpandingHandlers[i+1:]...)
			break
		}
	}
}

func (tv *TreeView) raiseItemExpanding(item *TreeViewItem) {
	args := &treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item}
	for _, handler := range tv.itemExpandingHandlers {
		handler(args)
	}
}

func (tv *TreeView) AddSelectionChangingHandler(handler TreeViewItemSelectionEventHandler) {
	tv.selectionChangingHandlers = append(tv.selectionChangingHandlers, handler)
}

func (tv *TreeView) RemoveSelectionChangingHandler(handler TreeViewItemSelectionEventHandler) {
	for i, h := range tv.selectionChangingHandlers {
		if h == handler {
			tv.selectionChangingHandlers = append(tv.selectionChangingHandlers[:i], tv.selectionChangingHandlers[i+1:]...)
			break
		}
	}
}

func (tv *TreeView) raiseSelectionChanging(old, new *TreeViewItem) {
	args := &treeViewItemSelectionEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, old: old, new: new}
	for _, handler := range tv.selectionChangingHandlers {
		handler(args)
	}
}

func (tv *TreeView) AddSelectionChangedHandler(handler TreeViewItemSelectionEventHandler) {
	tv.selectionChangedHandlers = append(tv.selectionChangedHandlers, handler)
}

func (tv *TreeView) RemoveSelectionChangedHandler(handler TreeViewItemSelectionEventHandler) {
	for i, h := range tv.selectionChangedHandlers {
		if h == handler {
			tv.selectionChangedHandlers = append(tv.selectionChangedHandlers[:i], tv.selectionChangedHandlers[i+1:]...)
			break
		}
	}
}

func (tv *TreeView) raiseSelectionChanged(old, new *TreeViewItem) {
	args := &treeViewItemSelectionEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, old: old, new: new}
	for _, handler := range tv.selectionChangedHandlers {
		handler(args)
	}
}

func (tv *TreeView) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_NOTIFY:
		nmtv := (*NMTREEVIEW)(unsafe.Pointer(msg.LParam))

		switch nmtv.Hdr.Code {
		case TVN_ITEMEXPANDED:
			item := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))

			switch nmtv.Action {
			case TVE_COLLAPSE:
				tv.raiseItemCollapsed(item)

			case TVE_COLLAPSERESET:

			case TVE_EXPAND:
				tv.raiseItemExpanded(item)

			case TVE_EXPANDPARTIAL:

			case TVE_TOGGLE:
			}

		case TVN_ITEMEXPANDING:
			item := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))

			switch nmtv.Action {
			case TVE_COLLAPSE:
				tv.raiseItemCollapsing(item)

			case TVE_COLLAPSERESET:

			case TVE_EXPAND:
				tv.raiseItemExpanding(item)

			case TVE_EXPANDPARTIAL:

			case TVE_TOGGLE:
			}

		case TVN_SELCHANGED:
			old := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemOld.LParam))
			new := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))
			tv.raiseSelectionChanged(old, new)

		case TVN_SELCHANGING:
			old := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemOld.LParam))
			new := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))
			tv.raiseSelectionChanging(old, new)
		}
	}

	return tv.Widget.wndProc(msg, origWndProcPtr)
}

func (tv *TreeView) onInsertingTreeViewItem(parent *TreeViewItem, index int, item *TreeViewItem) (err os.Error) {
	var tvi TVITEM
	var tvins TVINSERTSTRUCT

	tvi.LParam = uintptr(unsafe.Pointer(item))
	tvi.Mask = TVIF_TEXT | TVIF_PARAM
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

	if err == nil {
		for i, child := range item.children.items {
			err = tv.onInsertingTreeViewItem(item, i, child)
			if err != nil {
				return
			}
		}
	}

	return
}

func (tv *TreeView) onRemovingTreeViewItem(index int, item *TreeViewItem) (err os.Error) {
	if 0 == SendMessage(tv.hWnd, TVM_DELETEITEM, 0, uintptr(item.handle)) {
		err = newError("SendMessage(TVM_DELETEITEM) failed")
	}

	return
}

func (tv *TreeView) onClearingTreeViewItems(parent *TreeViewItem) (err os.Error) {
	for i, child := range parent.children.items {
		err = tv.onRemovingTreeViewItem(i, child)
		if err != nil {
			return
		}
	}

	return
}
