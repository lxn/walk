// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
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

type TreeView struct {
	Widget
	items                  *TreeViewItemList
	itemCollapsedHandlers  vector.Vector
	itemCollapsingHandlers vector.Vector
	itemExpandedHandlers   vector.Vector
	itemExpandingHandlers  vector.Vector
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
	tv.itemCollapsedHandlers.Push(handler)
}

func (tv *TreeView) RemoveItemCollapsedHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemCollapsedHandlers {
		if h.(TreeViewItemEventHandler) == handler {
			tv.itemCollapsedHandlers.Delete(i)
			break
		}
	}
}

func (tv *TreeView) raiseItemCollapsed(item *TreeViewItem) {
	for _, handlerIface := range tv.itemCollapsedHandlers {
		handler := handlerIface.(TreeViewItemEventHandler)
		handler(&treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item})
	}
}

func (tv *TreeView) AddItemCollapsingHandler(handler TreeViewItemEventHandler) {
	tv.itemCollapsingHandlers.Push(handler)
}

func (tv *TreeView) RemoveItemCollapsingHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemCollapsingHandlers {
		if h.(TreeViewItemEventHandler) == handler {
			tv.itemCollapsingHandlers.Delete(i)
			break
		}
	}
}

func (tv *TreeView) raiseItemCollapsing(item *TreeViewItem) {
	for _, handlerIface := range tv.itemCollapsingHandlers {
		handler := handlerIface.(TreeViewItemEventHandler)
		handler(&treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item})
	}
}

func (tv *TreeView) AddItemExpandedHandler(handler TreeViewItemEventHandler) {
	tv.itemExpandedHandlers.Push(handler)
}

func (tv *TreeView) RemoveItemExpandedHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemExpandedHandlers {
		if h.(TreeViewItemEventHandler) == handler {
			tv.itemExpandedHandlers.Delete(i)
			break
		}
	}
}

func (tv *TreeView) raiseItemExpanded(item *TreeViewItem) {
	for _, handlerIface := range tv.itemExpandedHandlers {
		handler := handlerIface.(TreeViewItemEventHandler)
		handler(&treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item})
	}
}

func (tv *TreeView) AddItemExpandingHandler(handler TreeViewItemEventHandler) {
	tv.itemExpandingHandlers.Push(handler)
}

func (tv *TreeView) RemoveItemExpandingHandler(handler TreeViewItemEventHandler) {
	for i, h := range tv.itemExpandingHandlers {
		if h.(TreeViewItemEventHandler) == handler {
			tv.itemExpandingHandlers.Delete(i)
			break
		}
	}
}

func (tv *TreeView) raiseItemExpanding(item *TreeViewItem) {
	for _, handlerIface := range tv.itemExpandingHandlers {
		handler := handlerIface.(TreeViewItemEventHandler)
		handler(&treeViewItemEventArgs{eventArgs: eventArgs{widgetsByHWnd[tv.hWnd]}, item: item})
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

	return
}

func (tv *TreeView) onRemovingTreeViewItem(index int, item *TreeViewItem) (err os.Error) {
	panic("not implemented")
}

func (tv *TreeView) onClearingTreeViewItems() (err os.Error) {
	panic("not implemented")
}
