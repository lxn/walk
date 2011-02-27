// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

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
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

var listViewOrigWndProcPtr uintptr
var _ subclassedWidget = &ListView{}

type SelectedIndexList struct {
	items []int
}

func NewSelectedIndexList(items []int) *SelectedIndexList {
	return &SelectedIndexList{items}
}

func (l *SelectedIndexList) At(index int) int {
	return l.items[index]
}

func (l *SelectedIndexList) Len() int {
	return len(l.items)
}

const (
	selectedIndexChangedTimerId = 1 + iota
	selectedIndexesChangedTimerId
)

type ListView struct {
	WidgetBase
	columns                         *ListViewColumnList
	items                           *ListViewItemList
	selectedIndex                   int
	selectedIndexChangedPublisher   EventPublisher
	selectedIndexes                 *SelectedIndexList
	selectedIndexesChangedPublisher EventPublisher
	itemActivatedPublisher          EventPublisher
	columnClickedPublisher          IntEventPublisher
	lastColumnStretched             bool
	persistent                      bool
	selectionChangedDelay           int
}

func NewListView(parent Container) (*ListView, os.Error) {
	lv := &ListView{selectedIndexes: NewSelectedIndexList(nil)}

	if err := initChildWidget(
		lv,
		parent,
		"SysListView32",
		WS_TABSTOP|WS_VISIBLE|LVS_SHOWSELALWAYS|LVS_REPORT,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			lv.Dispose()
		}
	}()

	lv.SetPersistent(true)

	exStyle := SendMessage(lv.hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	exStyle |= LVS_EX_DOUBLEBUFFER | LVS_EX_FULLROWSELECT
	SendMessage(lv.hWnd, LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)

	if err := lv.setTheme("Explorer"); err != nil {
		return nil, err
	}

	lv.columns = newListViewColumnList(lv)
	lv.items = newListViewItemList(lv)
	lv.selectedIndex = -1

	succeeded = true

	return lv, nil
}

func (*ListView) origWndProcPtr() uintptr {
	return listViewOrigWndProcPtr
}

func (*ListView) setOrigWndProcPtr(ptr uintptr) {
	listViewOrigWndProcPtr = ptr
}

func (lv *ListView) Dispose() {
	if lv.hWnd != 0 {
		if !KillTimer(lv.hWnd, selectedIndexChangedTimerId) {
			log.Print(lastError("KillTimer"))
		}
		if !KillTimer(lv.hWnd, selectedIndexesChangedTimerId) {
			log.Print(lastError("KillTimer"))
		}

		lv.WidgetBase.Dispose()
	}
}

func (*ListView) LayoutFlags() LayoutFlags {
	return HShrink | HGrow | VShrink | VGrow
}

func (lv *ListView) PreferredSize() Size {
	return lv.dialogBaseUnitsToPixels(Size{100, 100})
}

func (lv *ListView) Columns() *ListViewColumnList {
	return lv.columns
}

func (lv *ListView) SelectedColumnIndex() int {
	return int(SendMessage(lv.hWnd, LVM_GETSELECTEDCOLUMN, 0, 0))
}

func (lv *ListView) SetSelectedColumnIndex(value int) {
	SendMessage(lv.hWnd, LVM_SETSELECTEDCOLUMN, uintptr(value), 0)
}

func (lv *ListView) ColumnClicked() *IntEvent {
	return lv.columnClickedPublisher.Event()
}

/* FIXME: Our items list and that of the control get out of sync here.
func listViewCompare(lParam1, lParam2, lParamSort uintptr) int {
	compare := (*func(i, j int) int)(unsafe.Pointer(lParamSort))
	return (*compare)(int(lParam1), int(lParam2))
}

var listViewCompareCallback = uintptr(unsafe.Pointer(syscall.NewCallback(listViewCompare)))

func (lv *ListView) Sort(compare func(i, j int) int) os.Error {
	if compare == nil {
		return newError("compare must not be nil")
	}

	if FALSE == SendMessage(lv.hWnd, LVM_SORTITEMSEX, uintptr(unsafe.Pointer(&compare)), listViewCompareCallback) {
		return newError("LVM_SORTITEMSEX failed")
	}

	return nil
}*/

func (lv *ListView) Items() *ListViewItemList {
	return lv.items
}

func (lv *ListView) ItemActivated() *Event {
	return lv.itemActivatedPublisher.Event()
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
		lv.selectedIndexChangedPublisher.Publish()
	}

	return nil
}

func (lv *ListView) SelectedIndexChanged() *Event {
	return lv.selectedIndexChangedPublisher.Event()
}

func (lv *ListView) SingleItemSelection() bool {
	style := uint(GetWindowLong(lv.hWnd, GWL_STYLE))
	if style == 0 {
		log.Print(lastError("GetWindowLong"))
		return false
	}

	return style&LVS_SINGLESEL > 0
}

func (lv *ListView) SetSingleItemSelection(value bool) os.Error {
	return lv.ensureStyleBits(LVS_SINGLESEL, value)
}

func (lv *ListView) SelectedIndexes() *SelectedIndexList {
	return lv.selectedIndexes
}

func (lv *ListView) SelectionChangedDelay() int {
	return lv.selectionChangedDelay
}

func (lv *ListView) SetSelectionChangedDelay(delay int) {
	lv.selectionChangedDelay = delay
}

func (lv *ListView) updateSelectedIndexes() {
	count := int(SendMessage(lv.hWnd, LVM_GETSELECTEDCOUNT, 0, 0))
	indexes := make([]int, count)

	j := -1
	for i := 0; i < count; i++ {
		j = int(SendMessage(lv.hWnd, LVM_GETNEXTITEM, uintptr(j), LVNI_SELECTED))
		indexes[i] = j
	}

	changed := len(indexes) != len(lv.selectedIndexes.items)
	if !changed {
		for i := 0; i < len(indexes); i++ {
			if indexes[i] != lv.selectedIndexes.items[i] {
				changed = true
				break
			}
		}
	}

	if changed {
		lv.selectedIndexes.items = indexes
		if lv.selectionChangedDelay > 0 {
			if 0 == SetTimer(lv.hWnd, selectedIndexesChangedTimerId, uint(lv.selectionChangedDelay), 0) {
				log.Print(lastError("SetTimer"))
			}
		} else {
			lv.selectedIndexesChangedPublisher.Publish()
		}
	}
}

func (lv *ListView) SelectedIndexesChanged() *Event {
	return lv.selectedIndexesChangedPublisher.Event()
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
			width = 100
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

	lv.SetSuspended(true)
	defer lv.SetSuspended(false)

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

func (lv *ListView) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
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

	case WM_NOTIFY:
		switch int(((*NMHDR)(unsafe.Pointer(lParam))).Code) {
		case LVN_COLUMNCLICK:
			nmlv := (*NMLISTVIEW)(unsafe.Pointer(lParam))
			lv.columnClickedPublisher.Publish(nmlv.ISubItem)

		case LVN_ITEMCHANGED:
			nmlv := (*NMLISTVIEW)(unsafe.Pointer(lParam))
			selectedNow := nmlv.UNewState&LVIS_SELECTED > 0
			selectedBefore := nmlv.UOldState&LVIS_SELECTED > 0
			if selectedNow && !selectedBefore {
				lv.selectedIndex = nmlv.IItem
				if lv.selectionChangedDelay > 0 {
					if 0 == SetTimer(lv.hWnd, selectedIndexChangedTimerId, uint(lv.selectionChangedDelay), 0) {
						log.Print(lastError("SetTimer"))
					}
				} else {
					lv.selectedIndexChangedPublisher.Publish()
				}
			}
			if !lv.SingleItemSelection() {
				lv.updateSelectedIndexes()
			}

		case LVN_ITEMACTIVATE:
			lv.itemActivatedPublisher.Publish()
		}

	case WM_TIMER:
		switch wParam {
		case selectedIndexChangedTimerId:
			lv.selectedIndexChangedPublisher.Publish()

		case selectedIndexesChangedTimerId:
			lv.selectedIndexesChangedPublisher.Publish()
		}
	}

	return lv.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
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

	switch column.Alignment() {
	case AlignCenter:
		lvc.Fmt = 2

	case AlignFar:
		lvc.Fmt = 1
	}

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
