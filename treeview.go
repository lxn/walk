// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

type treeViewItemInfo struct {
	handle       HTREEITEM
	child2Handle map[TreeItem]HTREEITEM
}

type TreeView struct {
	WidgetBase
	model                         TreeModel
	lazyPopulation                bool
	itemsResetEventHandlerHandle  int
	itemChangedEventHandlerHandle int
	item2Info                     map[TreeItem]*treeViewItemInfo
	handle2Item                   map[HTREEITEM]TreeItem
	currItem                      TreeItem
	hIml                          HIMAGELIST
	usingSysIml                   bool
	imageUintptr2Index            map[uintptr]int32
	filePath2IconIndex            map[string]int32
	itemCollapsedPublisher        TreeItemEventPublisher
	itemExpandedPublisher         TreeItemEventPublisher
	currentItemChangedPublisher   EventPublisher
}

func NewTreeView(parent Container) (*TreeView, error) {
	tv := new(TreeView)

	if err := InitChildWidget(
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

	succeeded = true

	return tv, nil
}

func (*TreeView) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (tv *TreeView) SizeHint() Size {
	return tv.dialogBaseUnitsToPixels(Size{100, 100})
}

func (tv *TreeView) Dispose() {
	tv.WidgetBase.Dispose()

	tv.disposeImageListAndCaches()
}

func (tv *TreeView) Model() TreeModel {
	return tv.model
}

func (tv *TreeView) SetModel(model TreeModel) error {
	if tv.model != nil {
		tv.model.ItemsReset().Detach(tv.itemsResetEventHandlerHandle)
		tv.model.ItemChanged().Detach(tv.itemChangedEventHandlerHandle)

		tv.disposeImageListAndCaches()
	}

	tv.model = model

	if model != nil {
		tv.lazyPopulation = model.LazyPopulation()

		tv.itemsResetEventHandlerHandle = model.ItemsReset().Attach(func(parent TreeItem) {
			if parent == nil {
				tv.resetItems()
			} else if tv.item2Info[parent] != nil {
				if err := tv.removeDescendants(parent); err != nil {
					return
				}

				if err := tv.insertChildren(parent); err != nil {
					return
				}
			}
		})

		tv.itemChangedEventHandlerHandle = model.ItemChanged().Attach(func(item TreeItem) {
			if item == nil || tv.item2Info[item] == nil {
				return
			}

			if err := tv.updateItem(item); err != nil {
				return
			}
		})
	}

	return tv.resetItems()
}

func (tv *TreeView) CurrentItem() TreeItem {
	return tv.currItem
}

func (tv *TreeView) SetCurrentItem(item TreeItem) error {
	if item == tv.currItem {
		return nil
	}

	var handle HTREEITEM
	if item != nil {
		if info := tv.item2Info[item]; info == nil {
			return newError("invalid item")
		} else {
			handle = info.handle
		}
	}

	if 0 == tv.SendMessage(TVM_SELECTITEM, TVGN_CARET, uintptr(handle)) {
		return newError("SendMessage(TVM_SELECTITEM) failed")
	}

	tv.currItem = item

	return nil
}

func (tv *TreeView) ItemAt(x, y int) TreeItem {
	hti := TVHITTESTINFO{Pt: POINT{int32(x), int32(y)}}

	tv.SendMessage(TVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))

	if item, ok := tv.handle2Item[hti.HItem]; ok {
		return item
	}

	return nil
}

func (tv *TreeView) resetItems() error {
	tv.SetSuspended(true)
	defer tv.SetSuspended(false)

	if err := tv.clearItems(); err != nil {
		return err
	}

	if tv.model == nil {
		return nil
	}

	if err := tv.insertRoots(); err != nil {
		return err
	}

	return nil
}

func (tv *TreeView) clearItems() error {
	if 0 == tv.SendMessage(TVM_DELETEITEM, 0, 0) {
		return newError("SendMessage(TVM_DELETEITEM) failed")
	}

	tv.item2Info = make(map[TreeItem]*treeViewItemInfo)
	tv.handle2Item = make(map[HTREEITEM]TreeItem)

	return nil
}

func (tv *TreeView) insertRoots() error {
	count := tv.model.RootCount()

	for i := 0; i < count; i++ {
		if _, err := tv.insertItem(i, tv.model.RootAt(i)); err != nil {
			return err
		}
	}

	return nil
}

func (tv *TreeView) applyImageListForImage(image interface{}) {
	tv.hIml, tv.usingSysIml, _ = imageListForImage(image)

	tv.SendMessage(TVM_SETIMAGELIST, 0, uintptr(tv.hIml))

	tv.imageUintptr2Index = make(map[uintptr]int32)
	tv.filePath2IconIndex = make(map[string]int32)
}

func (tv *TreeView) disposeImageListAndCaches() {
	if tv.hIml != 0 && !tv.usingSysIml {
		ImageList_Destroy(tv.hIml)
	}
	tv.hIml = 0

	tv.imageUintptr2Index = nil
	tv.filePath2IconIndex = nil
}

func (tv *TreeView) setTVITEMImageInfo(tvi *TVITEM, item TreeItem) {
	if imager, ok := item.(Imager); ok {
		if tv.hIml == 0 {
			tv.applyImageListForImage(imager.Image())
		}

		// FIXME: If not setting TVIF_SELECTEDIMAGE and tvi.ISelectedImage, 
		// some default icon will show up, even though we have not asked for it.

		tvi.Mask |= TVIF_IMAGE | TVIF_SELECTEDIMAGE
		tvi.IImage = imageIndexMaybeAdd(
			imager.Image(),
			tv.hIml,
			tv.usingSysIml,
			tv.imageUintptr2Index,
			tv.filePath2IconIndex)

		tvi.ISelectedImage = tvi.IImage
	}
}

func (tv *TreeView) insertItem(index int, item TreeItem) (HTREEITEM, error) {
	var tvins TVINSERTSTRUCT
	tvi := &tvins.Item

	tvi.Mask = TVIF_CHILDREN | TVIF_TEXT
	tvi.PszText = syscall.StringToUTF16Ptr(item.Text())
	tvi.CChildren = I_CHILDRENCALLBACK

	tv.setTVITEMImageInfo(tvi, item)

	parent := item.Parent()

	if parent == nil {
		tvins.HParent = TVI_ROOT
	} else {
		info := tv.item2Info[parent]
		if info == nil {
			return 0, newError("invalid parent")
		}
		tvins.HParent = info.handle
	}

	if index == 0 {
		tvins.HInsertAfter = TVI_LAST
	} else {
		var prevItem TreeItem
		if parent == nil {
			prevItem = tv.model.RootAt(index - 1)
		} else {
			prevItem = parent.ChildAt(index - 1)
		}
		info := tv.item2Info[prevItem]
		if info == nil {
			return 0, newError("invalid prev item")
		}
		tvins.HInsertAfter = info.handle
	}

	hItem := HTREEITEM(tv.SendMessage(TVM_INSERTITEM, 0, uintptr(unsafe.Pointer(&tvins))))
	if hItem == 0 {
		return 0, newError("TVM_INSERTITEM failed")
	}
	tv.item2Info[item] = &treeViewItemInfo{hItem, make(map[TreeItem]HTREEITEM)}
	tv.handle2Item[hItem] = item

	if !tv.lazyPopulation {
		if err := tv.insertChildren(item); err != nil {
			return 0, err
		}
	}

	return hItem, nil
}

func (tv *TreeView) insertChildren(parent TreeItem) error {
	info := tv.item2Info[parent]

	count := parent.ChildCount()
	for i := 0; i < count; i++ {
		child := parent.ChildAt(i)

		if handle, err := tv.insertItem(i, child); err != nil {
			return err
		} else {
			info.child2Handle[child] = handle
		}
	}

	return nil
}

func (tv *TreeView) updateItem(item TreeItem) error {
	tvi := &TVITEM{
		Mask:    TVIF_TEXT,
		HItem:   tv.item2Info[item].handle,
		PszText: syscall.StringToUTF16Ptr(item.Text()),
	}

	tv.setTVITEMImageInfo(tvi, item)

	if 0 == tv.SendMessage(TVM_SETITEM, 0, uintptr(unsafe.Pointer(tvi))) {
		return newError("SendMessage(TVM_SETITEM) failed")
	}

	return nil
}

func (tv *TreeView) removeItem(item TreeItem) error {
	if err := tv.removeDescendants(item); err != nil {
		return err
	}

	info := tv.item2Info[item]
	if info == nil {
		return newError("invalid item")
	}

	if 0 == tv.SendMessage(TVM_DELETEITEM, 0, uintptr(info.handle)) {
		return newError("SendMessage(TVM_DELETEITEM) failed")
	}

	if parentInfo := tv.item2Info[item.Parent()]; parentInfo != nil {
		delete(parentInfo.child2Handle, item)
	}
	delete(tv.item2Info, item)
	delete(tv.handle2Item, info.handle)

	return nil
}

func (tv *TreeView) removeDescendants(parent TreeItem) error {
	for item, _ := range tv.item2Info[parent].child2Handle {
		if err := tv.removeItem(item); err != nil {
			return err
		}
	}

	return nil
}

func (tv *TreeView) ItemCollapsed() *TreeItemEvent {
	return tv.itemCollapsedPublisher.Event()
}

func (tv *TreeView) ItemExpanded() *TreeItemEvent {
	return tv.itemExpandedPublisher.Event()
}

func (tv *TreeView) CurrentItemChanged() *Event {
	return tv.currentItemChangedPublisher.Event()
}

func (tv *TreeView) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_NOTIFY:
		nmhdr := (*NMHDR)(unsafe.Pointer(lParam))

		switch nmhdr.Code {
		case TVN_GETDISPINFO:
			nmtvdi := (*NMTVDISPINFO)(unsafe.Pointer(lParam))
			item := tv.handle2Item[nmtvdi.Item.HItem]

			if nmtvdi.Item.Mask&TVIF_CHILDREN != 0 {
				nmtvdi.Item.CChildren = int32(item.ChildCount())
			}

		case TVN_ITEMEXPANDING:
			nmtv := (*NMTREEVIEW)(unsafe.Pointer(lParam))
			item := tv.handle2Item[nmtv.ItemNew.HItem]

			if nmtv.Action == TVE_EXPAND && tv.lazyPopulation {
				info := tv.item2Info[item]
				if len(info.child2Handle) == 0 {
					tv.insertChildren(item)
				}
			}

		case TVN_ITEMEXPANDED:
			nmtv := (*NMTREEVIEW)(unsafe.Pointer(lParam))
			item := tv.handle2Item[nmtv.ItemNew.HItem]

			switch nmtv.Action {
			case TVE_COLLAPSE:
				tv.itemCollapsedPublisher.Publish(item)

			case TVE_COLLAPSERESET:

			case TVE_EXPAND:
				tv.itemExpandedPublisher.Publish(item)

			case TVE_EXPANDPARTIAL:

			case TVE_TOGGLE:
			}

		case TVN_SELCHANGED:
			nmtv := (*NMTREEVIEW)(unsafe.Pointer(lParam))

			tv.currItem = tv.handle2Item[nmtv.ItemNew.HItem]

			tv.currentItemChangedPublisher.Publish()
		}
	}

	return tv.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
