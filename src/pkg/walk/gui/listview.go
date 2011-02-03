// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

var listViewSubclassWndProcPtr uintptr
var listViewOrigWndProcPtr uintptr

func listViewSubclassWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	lv, ok := widgetsByHWnd[hwnd].(*ListView)
	if !ok {
		return CallWindowProc(listViewOrigWndProcPtr, hwnd, msg, wParam, lParam)
	}

	return lv.wndProc(hwnd, msg, wParam, lParam, listViewOrigWndProcPtr)
}

type ListView struct {
	Widget
	columns                       *ListViewColumnList
	items                         *ListViewItemList
	selectedIndex                 int
	selectedIndexChangedPublisher EventPublisher
	itemActivatedPublisher        EventPublisher
	lastColumnStretched           bool
	persistent                    bool
}

func NewListView(parent IContainer) (*ListView, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	if listViewSubclassWndProcPtr == 0 {
		listViewSubclassWndProcPtr = syscall.NewCallback(listViewSubclassWndProc)
	}

	hWnd := CreateWindowEx(
		WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("SysListView32"), nil,
		LVS_SINGLESEL|LVS_SHOWSELALWAYS|LVS_REPORT|WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	lv := &ListView{
		Widget: Widget{
			hWnd:   hWnd,
			parent: parent,
		},
	}

	succeeded := false
	defer func() {
		if !succeeded {
			lv.Dispose()
		}
	}()

	lv.SetPersistent(true)

	listViewOrigWndProcPtr = uintptr(SetWindowLong(hWnd, GWL_WNDPROC, int(listViewSubclassWndProcPtr)))
	if listViewOrigWndProcPtr == 0 {
		return nil, lastError("SetWindowLong")
	}

	exStyle := SendMessage(hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	exStyle |= LVS_EX_DOUBLEBUFFER | LVS_EX_FULLROWSELECT //| LVS_EX_GRIDLINES
	SendMessage(hWnd, LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)

	if err := lv.setTheme("Explorer"); err != nil {
		return nil, err
	}

	lv.columns = newListViewColumnList(lv)
	lv.items = newListViewItemList(lv)
	lv.selectedIndex = -1

	lv.SetFont(defaultFont)

	if err := parent.Children().Add(lv); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = lv

	succeeded = true

	return lv, nil
}

func (*ListView) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz | ShrinkVert | GrowVert
}

func (lv *ListView) PreferredSize() drawing.Size {
	return lv.dialogBaseUnitsToPixels(drawing.Size{100, 100})
}

func (lv *ListView) Columns() *ListViewColumnList {
	return lv.columns
}

func (lv *ListView) Items() *ListViewItemList {
	return lv.items
}

func (lv *ListView) SelectedIndex() int {
	return lv.selectedIndex
}

func (lv *ListView) SetSelectedIndex(value int) os.Error {
	var lvi LVITEM

	lvi.StateMask = LVIS_SELECTED
	if value > -1 {
		lvi.State = LVIS_SELECTED
	}

	if FALSE == SendMessage(lv.hWnd, LVM_SETITEMSTATE, uintptr(value), uintptr(unsafe.Pointer(&lvi))) {
		return newError("failed to set selected item")
	}

	lv.selectedIndex = value

	if value == -1 {
		lv.selectedIndexChangedPublisher.Publish(NewEventArgs(lv))
	}

	return nil
}

func (lv *ListView) LastColumnStretched() bool {
	return lv.lastColumnStretched
}

func (lv *ListView) SetLastColumnStretched(value bool) os.Error {
	if value {
		if err := lv.StretchLastColumn(); err != nil {
			return err
		}
	}

	lv.lastColumnStretched = value

	return nil
}

func (lv *ListView) StretchLastColumn() os.Error {
	colCount := lv.columns.Len()
	if colCount == 0 {
		return nil
	}

	if 0 == SendMessage(lv.hWnd, LVM_SETCOLUMNWIDTH, uintptr(colCount-1), LVSCW_AUTOSIZE_USEHEADER) {
		return newError("LVM_SETCOLUMNWIDTH failed")
	}

	return nil
}

func (lv *ListView) Persistent() bool {
	return lv.persistent
}

func (lv *ListView) SetPersistent(value bool) {
	lv.persistent = value
}

func (lv *ListView) SaveState() os.Error {
	buf := bytes.NewBuffer(nil)

	count := lv.columns.Len()
	for i := 0; i < count; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}

		width := SendMessage(lv.hWnd, LVM_GETCOLUMNWIDTH, uintptr(i), 0)
		if width == 0 {
			return newError("LVM_GETCOLUMNWIDTH failed")
		}

		buf.WriteString(strconv.Itoa(int(width)))
	}

	return lv.putState(buf.String())
}

func (lv *ListView) RestoreState() os.Error {
	state, err := lv.getState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	widthStrs := strings.Split(state, " ", -1)

	for i, str := range widthStrs {
		width, err := strconv.Atoi(str)
		if err != nil {
			return err
		}

		if FALSE == SendMessage(lv.hWnd, LVM_SETCOLUMNWIDTH, uintptr(i), uintptr(width)) {
			return newError("LVM_SETCOLUMNWIDTH failed")
		}
	}

	return nil
}

func (lv *ListView) ItemActivated() *Event {
	return lv.itemActivatedPublisher.Event()
}

func (lv *ListView) SelectedIndexChanged() *Event {
	return lv.selectedIndexChangedPublisher.Event()
}

func (lv *ListView) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_ERASEBKGND:
		if lv.lastColumnStretched {
			lv.StretchLastColumn()
		}
		return 1

	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

	case WM_KEYDOWN:
		if wParam == VK_RETURN && lv.SelectedIndex() > -1 {
			lv.itemActivatedPublisher.Publish(NewEventArgs(lv))
		}

	case WM_NOTIFY:
		switch int(((*NMHDR)(unsafe.Pointer(lParam))).Code) {
		case LVN_ITEMCHANGED:
			nmlv := (*NMLISTVIEW)(unsafe.Pointer(lParam))
			selectedNow := nmlv.UNewState&LVIS_SELECTED > 0
			selectedBefore := nmlv.UOldState&LVIS_SELECTED > 0
			if selectedNow && !selectedBefore {
				lv.selectedIndex = nmlv.IItem
				lv.selectedIndexChangedPublisher.Publish(NewEventArgs(lv))
			}

		case LVN_ITEMACTIVATE:
			lv.itemActivatedPublisher.Publish(NewEventArgs(lv))
		}
	}

	return lv.Widget.wndProc(hwnd, msg, wParam, lParam, listViewOrigWndProcPtr)
}

func (lv *ListView) onListViewColumnChanged(column *ListViewColumn) {
	panic("not implemented")
}

func (lv *ListView) onInsertingListViewColumn(index int, column *ListViewColumn) (err os.Error) {
	var lvc LVCOLUMN

	lvc.Mask = LVCF_FMT | LVCF_WIDTH | LVCF_TEXT | LVCF_SUBITEM
	lvc.ISubItem = index
	lvc.PszText = syscall.StringToUTF16Ptr(column.Title())
	lvc.Cx = column.width
	lvc.Fmt = int(column.alignment)

	i := SendMessage(lv.hWnd, LVM_INSERTCOLUMN, uintptr(index), uintptr(unsafe.Pointer(&lvc)))
	if int(i) == -1 {
		return newError("ListView.onInsertingListViewColumn: Failed to insert column.")
	}

	column.addChangedHandler(lv)

	return
}

func (lv *ListView) onRemovingListViewColumn(index int, column *ListViewColumn) (err os.Error) {
	panic("not implemented")
}

func (lv *ListView) onClearingListViewColumns() (err os.Error) {
	panic("not implemented")
}

func (lv *ListView) onListViewItemChanged(item *ListViewItem) {
	var lvi LVITEM

	lvi.Mask = LVIF_TEXT
	lvi.IItem = lv.Items().Index(item)

	texts := item.Texts()

	colCount := lv.columns.Len()

	for colIndex := 0; colIndex < colCount; colIndex++ {
		lvi.ISubItem = colIndex

		if colIndex < len(texts) {
			lvi.PszText = syscall.StringToUTF16Ptr(texts[colIndex])
		} else {
			lvi.PszText = nil
		}

		ret := SendMessage(lv.hWnd, LVM_SETITEM, 0, uintptr(unsafe.Pointer(&lvi)))
		if ret == 0 {
			log.Println(newError("ListView.onInsertingListViewItem: Failed to set sub item."))
		}
	}
}

func (lv *ListView) onInsertingListViewItem(index int, item *ListViewItem) (err os.Error) {
	var lvi LVITEM

	lvi.Mask = LVIF_TEXT
	lvi.IItem = index

	texts := item.Texts()
	if len(texts) > 0 {
		lvi.PszText = syscall.StringToUTF16Ptr(texts[0])
	}

	i := SendMessage(lv.hWnd, LVM_INSERTITEM, 0, uintptr(unsafe.Pointer(&lvi)))
	if int(i) == -1 {
		return newError("ListView.onInsertingListViewItem: Failed to insert item.")
	}

	colCount := lv.columns.Len()

	for colIndex := 1; colIndex < colCount; colIndex++ {
		lvi.ISubItem = colIndex

		if colIndex < len(texts) {
			lvi.PszText = syscall.StringToUTF16Ptr(texts[colIndex])
		} else {
			lvi.PszText = nil
		}

		ret := SendMessage(lv.hWnd, LVM_SETITEM, 0, uintptr(unsafe.Pointer(&lvi)))
		if ret == 0 {
			return newError("ListView.onInsertingListViewItem: Failed to set sub item.")
		}
	}

	item.addChangedHandler(lv)

	return
}

func (lv *ListView) onRemovingListViewItem(index int, item *ListViewItem) (err os.Error) {
	if 0 == SendMessage(lv.hWnd, LVM_DELETEITEM, uintptr(index), 0) {
		return newError("LVM_DELETEITEM failed")
	}

	item.removeChangedHandler(lv)

	if index == lv.selectedIndex {
		return lv.SetSelectedIndex(-1)
	}

	return nil
}

func (lv *ListView) onClearingListViewItems() os.Error {
	if FALSE == SendMessage(lv.hWnd, LVM_DELETEALLITEMS, 0, 0) {
		return newError("LVM_DELETEALLITEMS failed")
	}

	for _, item := range lv.items.items {
		item.removeChangedHandler(lv)
	}

	return lv.SetSelectedIndex(-1)
}
