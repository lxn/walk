// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import . "walk/winapi"

var treeViewOrigWndProcPtr uintptr
var _ subclassedWidget = &TreeView{}

type TreeView struct {
	WidgetBase
	items                      *TreeViewItemList
	itemCollapsedPublisher     TreeViewItemEventPublisher
	itemCollapsingPublisher    TreeViewItemEventPublisher
	itemExpandedPublisher      TreeViewItemEventPublisher
	itemExpandingPublisher     TreeViewItemEventPublisher
	selectionChangedPublisher  TreeViewItemSelectionEventPublisher
	selectionChangingPublisher TreeViewItemSelectionEventPublisher
}

func NewTreeView(parent Container) (*TreeView, os.Error) {
	tv := &TreeView{}

	if err := initChildWidget(
		tv,
		parent,
		"SysTreeView32",
		WS_TABSTOP|WS_VISIBLE|TVS_HASBUTTONS|TVS_HASLINES|TVS_LINESATROOT|TVS_SHOWSELALWAYS,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			tv.Dispose()
		}
	}()

	if err := tv.setTheme("Explorer"); err != nil {
		return nil, err
	}

	tv.items = newTreeViewItemList(tv)

	succeeded = true

	return tv, nil
}

func (*TreeView) origWndProcPtr() uintptr {
	return treeViewOrigWndProcPtr
}

func (*TreeView) setOrigWndProcPtr(ptr uintptr) {
	treeViewOrigWndProcPtr = ptr
}

func (*TreeView) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (tv *TreeView) SizeHint() Size {
	return tv.dialogBaseUnitsToPixels(Size{100, 100})
}

func (tv *TreeView) Items() *TreeViewItemList {
	return tv.items
}

func (tv *TreeView) ItemCollapsed() *TreeViewItemEvent {
	return tv.itemCollapsedPublisher.Event()
}

func (tv *TreeView) ItemCollapsing() *TreeViewItemEvent {
	return tv.itemCollapsingPublisher.Event()
}

func (tv *TreeView) ItemExpanded() *TreeViewItemEvent {
	return tv.itemExpandedPublisher.Event()
}

func (tv *TreeView) ItemExpanding() *TreeViewItemEvent {
	return tv.itemExpandingPublisher.Event()
}

func (tv *TreeView) SelectionChanged() *TreeViewItemSelectionEvent {
	return tv.selectionChangedPublisher.Event()
}

func (tv *TreeView) SelectionChanging() *TreeViewItemSelectionEvent {
	return tv.selectionChangingPublisher.Event()
}

func (tv *TreeView) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_NOTIFY:
		nmtv := (*NMTREEVIEW)(unsafe.Pointer(lParam))

		switch nmtv.Hdr.Code {
		case TVN_ITEMEXPANDED:
			item := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))

			switch nmtv.Action {
			case TVE_COLLAPSE:
				tv.itemCollapsedPublisher.Publish(item)

			case TVE_COLLAPSERESET:

			case TVE_EXPAND:
				tv.itemExpandedPublisher.Publish(item)

			case TVE_EXPANDPARTIAL:

			case TVE_TOGGLE:
			}

		case TVN_ITEMEXPANDING:
			item := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))

			switch nmtv.Action {
			case TVE_COLLAPSE:
				tv.itemCollapsingPublisher.Publish(item)

			case TVE_COLLAPSERESET:

			case TVE_EXPAND:
				tv.itemExpandingPublisher.Publish(item)

			case TVE_EXPANDPARTIAL:

			case TVE_TOGGLE:
			}

		case TVN_SELCHANGED:
			old := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemOld.LParam))
			new := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))
			tv.selectionChangedPublisher.Publish(old, new)

		case TVN_SELCHANGING:
			old := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemOld.LParam))
			new := (*TreeViewItem)(unsafe.Pointer(nmtv.ItemNew.LParam))
			tv.selectionChangingPublisher.Publish(old, new)
		}
	}

	return tv.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
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
