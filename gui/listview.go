// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"bytes"
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

var listViewSubclassWndProcCallback *syscall.Callback
var listViewSubclassWndProcPtr uintptr
var listViewOrigWndProcPtr uintptr

func listViewSubclassWndProc(args *uintptr) uintptr {
	msg := msgFromCallbackArgs(args)

	lv, ok := widgetsByHWnd[msg.HWnd].(*ListView)
	if !ok {
		// Before CreateWindowEx returns, among others, WM_GETMINMAXINFO is sent.
		// FIXME: Find a way to properly handle this.
		return CallWindowProc(listViewOrigWndProcPtr, msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return lv.wndProc(msg, listViewOrigWndProcPtr)
}

type ListView struct {
	Widget
	columns                       *ListViewColumnList
	items                         *ListViewItemList
	prevSelIndex                  int
	selectedIndexChangedPublisher EventPublisher
	itemActivatedPublisher        EventPublisher
	lastColumnStretched           bool
}

func NewListView(parent IContainer) (*ListView, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	if listViewSubclassWndProcCallback == nil {
		listViewSubclassWndProcCallback = syscall.NewCallback(listViewSubclassWndProc, 4)
		listViewSubclassWndProcPtr = uintptr(listViewSubclassWndProcCallback.ExtFnEntry())
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
	lv.prevSelIndex = -1

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
	selCount := int(SendMessage(lv.hWnd, LVM_GETSELECTEDCOUNT, 0, 0))
	if selCount == 0 {
		return -1
	}

	return int(SendMessage(lv.hWnd, LVM_GETNEXTITEM, ^uintptr(0), LVNI_SELECTED)) // ^uintptr(0) == -1
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

func (lv *ListView) SaveState() (string, os.Error) {
	buf := bytes.NewBuffer(nil)

	count := lv.columns.Len()
	for i := 0; i < count; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}

		width := SendMessage(lv.hWnd, LVM_GETCOLUMNWIDTH, uintptr(i), 0)
		if width == 0 {
			return "", newError("LVM_GETCOLUMNWIDTH failed")
		}

		buf.WriteString(strconv.Itoa(int(width)))
	}

	return buf.String(), nil
}

func (lv *ListView) RestoreState(s string) os.Error {
	widthStrs := strings.Split(s, " ", -1)

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

func (lv *ListView) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	switch msg.Message {
	case WM_ERASEBKGND:
		if lv.lastColumnStretched {
			lv.StretchLastColumn()
		}
		return 1

	case WM_GETDLGCODE:
		if msg.WParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

	case WM_KEYDOWN:
		if msg.WParam == VK_RETURN && lv.SelectedIndex() > -1 {
			lv.itemActivatedPublisher.Publish(NewEventArgs(lv))
		}

	case WM_NOTIFY:
		switch int(((*NMHDR)(unsafe.Pointer(msg.LParam))).Code) {
		case LVN_ITEMCHANGED:
			if selIndex := lv.SelectedIndex(); selIndex != lv.prevSelIndex {
				lv.selectedIndexChangedPublisher.Publish(NewEventArgs(lv))
				lv.prevSelIndex = selIndex
			}

		case LVN_ITEMACTIVATE:
			lv.itemActivatedPublisher.Publish(NewEventArgs(lv))
		}
	}

	return lv.Widget.wndProc(msg, listViewOrigWndProcPtr)
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
	panic("not implemented")
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

	return
}

func (lv *ListView) onRemovingListViewItem(index int, item *ListViewItem) (err os.Error) {
	if 0 == SendMessage(lv.hWnd, LVM_DELETEITEM, uintptr(index), 0) {
		return newError("LVM_DELETEITEM failed")
	}

	return nil
}

func (lv *ListView) onClearingListViewItems() os.Error {
	if FALSE == SendMessage(lv.hWnd, LVM_DELETEALLITEMS, 0, 0) {
		return newError("LVM_DELETEALLITEMS failed")
	}

	return nil
}
