// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

var defaultTVRowBGColor Color = Color(GetSysColor(COLOR_WINDOW))

const (
	tableViewCurrentIndexChangedTimerId = 1 + iota
	tableViewSelectedIndexesChangedTimerId
)

// TableView is a model based widget for record centric, tabular data.
//
// TableView is implemented as a virtual mode list view to support quite large 
// amounts of data.
type TableView struct {
	WidgetBase
	model                           TableModel
	itemChecker                     ItemChecker
	imageProvider                   ImageProvider
	hasAppliedImageList             bool
	imageList                       *ImageList
	imageUintptr2Index              map[uintptr]int32
	filePath2IconIndex              map[string]int32
	rowsResetHandlerHandle          int
	rowChangedHandlerHandle         int
	sortChangedHandlerHandle        int
	columns                         []TableColumn
	currentIndex                    int
	currentIndexChangedPublisher    EventPublisher
	selectedIndexes                 *IndexList
	selectedIndexesChangedPublisher EventPublisher
	itemActivatedPublisher          EventPublisher
	columnClickedPublisher          IntEventPublisher
	lastColumnStretched             bool
	inEraseBkgnd                    bool
	persistent                      bool
	itemStateChangedEventDelay      int
	alternatingRowBGColor           Color
}

// NewTableView creates and returns a *TableView as child of the specified
// Container.
func NewTableView(parent Container) (*TableView, error) {
	tv := &TableView{
		alternatingRowBGColor: defaultTVRowBGColor,
		imageUintptr2Index:    make(map[uintptr]int32),
		filePath2IconIndex:    make(map[string]int32),
		selectedIndexes:       NewIndexList(nil),
	}

	if err := initChildWidget(
		tv,
		parent,
		"SysListView32",
		WS_TABSTOP|WS_VISIBLE|LVS_OWNERDATA|LVS_SHOWSELALWAYS|LVS_REPORT,
		WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			tv.Dispose()
		}
	}()

	tv.SetPersistent(true)

	exStyle := SendMessage(tv.hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	exStyle |= LVS_EX_DOUBLEBUFFER | LVS_EX_FULLROWSELECT
	SendMessage(tv.hWnd, LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)

	if err := tv.setTheme("Explorer"); err != nil {
		return nil, err
	}

	tv.currentIndex = -1

	succeeded = true

	return tv, nil
}

// Dispose releases the operating system resources, associated with the 
// *TableView.
func (tv *TableView) Dispose() {
	tv.detachModel()

	if tv.hWnd != 0 {
		if !KillTimer(tv.hWnd, tableViewCurrentIndexChangedTimerId) {
			lastError("KillTimer")
		}
		if !KillTimer(tv.hWnd, tableViewSelectedIndexesChangedTimerId) {
			lastError("KillTimer")
		}

		tv.WidgetBase.Dispose()
	}
}

// LayoutFlags returns a combination of LayoutFlags that specify how the
// *TableView wants to be treated by Layout implementations.
func (*TableView) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

// MinSizeHint returns the minimum outer Size, including decorations, that 
// makes sense for the *TableView.
func (tv *TableView) MinSizeHint() Size {
	return Size{10, 10}
}

// SizeHint returns a sensible Size for a *TableView.
func (tv *TableView) SizeHint() Size {
	return Size{100, 100}
}

// ReorderColumnsEnabled returns if the user can reorder columns by dragging and
// dropping column headers.
func (tv *TableView) ReorderColumnsEnabled() bool {
	exStyle := SendMessage(tv.hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	return exStyle&LVS_EX_HEADERDRAGDROP > 0
}

// SetReorderColumnsEnabled sets if the user can reorder columns by dragging and
// dropping column headers.
func (tv *TableView) SetReorderColumnsEnabled(enabled bool) {
	exStyle := SendMessage(tv.hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	if enabled {
		exStyle |= LVS_EX_HEADERDRAGDROP
	} else {
		exStyle &^= LVS_EX_HEADERDRAGDROP
	}
	SendMessage(tv.hWnd, LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)
}

// AlternatingRowBGColor returns the alternating row background color.
func (tv *TableView) AlternatingRowBGColor() Color {
	return tv.alternatingRowBGColor
}

// SetAlternatingRowBGColor sets the alternating row background color.
func (tv *TableView) SetAlternatingRowBGColor(c Color) {
	tv.alternatingRowBGColor = c

	tv.Invalidate()
}

func (tv *TableView) attachModel() {
	tv.rowsResetHandlerHandle = tv.model.RowsReset().Attach(func() {
		tv.setItemCount()

		tv.SetCurrentIndex(-1)
	})

	tv.rowChangedHandlerHandle = tv.model.RowChanged().Attach(func(row int) {
		if FALSE == SendMessage(tv.hWnd, LVM_UPDATE, uintptr(row), 0) {
			newError("SendMessage(LVM_UPDATE)")
		}
	})

	if sorter, ok := tv.model.(Sorter); ok {
		tv.sortChangedHandlerHandle = sorter.SortChanged().Attach(func() {
			col := sorter.SortedColumn()
			tv.setSelectedColumnIndex(col)
			tv.setSortIcon(col, sorter.SortOrder())
			tv.Invalidate()
		})
	}
}

func (tv *TableView) detachModel() {
	tv.model.RowsReset().Detach(tv.rowsResetHandlerHandle)
	tv.model.RowChanged().Detach(tv.rowChangedHandlerHandle)
	if sorter, ok := tv.model.(Sorter); ok {
		sorter.SortChanged().Detach(tv.sortChangedHandlerHandle)
	}
}

// Model returns the TableModel that provides data to the *TableView.
func (tv *TableView) Model() TableModel {
	return tv.model
}

// SetModel sets the TableModel that provides data to the *TableView.
func (tv *TableView) SetModel(model TableModel) error {
	tv.SetSuspended(true)
	defer tv.SetSuspended(false)

	if tv.model != nil {
		for _ = range tv.columns {
			if FALSE == SendMessage(tv.hWnd, LVM_DELETECOLUMN, 0, 0) {
				return newError("SendMessage(LVM_DELETECOLUMN)")
			}
		}

		tv.detachModel()
	}

	tv.model = model

	tv.itemChecker, _ = model.(ItemChecker)
	tv.imageProvider, _ = model.(ImageProvider)

	if tv.imageList != nil {
		SendMessage(tv.hWnd, LVM_SETIMAGELIST, LVSIL_SMALL, 0)
		tv.imageList.Dispose()
		tv.imageList = nil
	}
	tv.hasAppliedImageList = false

	if model != nil {
		tv.attachModel()

		tv.columns = model.Columns()

		for i, column := range tv.columns {
			if column.Format == "" {
				tv.columns[i].Format = "%v"
			}

			var lvc LVCOLUMN

			lvc.Mask = LVCF_FMT | LVCF_WIDTH | LVCF_TEXT | LVCF_SUBITEM
			lvc.ISubItem = int32(i)
			lvc.PszText = syscall.StringToUTF16Ptr(column.Title)
			if column.Width > 0 {
				lvc.Cx = int32(column.Width)
			} else {
				lvc.Cx = 100
			}

			switch column.Alignment {
			case AlignCenter:
				lvc.Fmt = 2

			case AlignFar:
				lvc.Fmt = 1
			}

			j := SendMessage(tv.hWnd, LVM_INSERTCOLUMN, uintptr(i), uintptr(unsafe.Pointer(&lvc)))
			if int(j) == -1 {
				return newError("TableView.SetModel: Failed to insert column.")
			}
		}

		if sorter, ok := tv.model.(Sorter); ok {
			col := sorter.SortedColumn()
			tv.setSelectedColumnIndex(col)
			tv.setSortIcon(col, sorter.SortOrder())
		}

		return tv.setItemCount()
	}

	return nil
}

func (tv *TableView) setItemCount() error {
	var count int

	if tv.model != nil {
		count = tv.model.RowCount()
	}

	if 0 == SendMessage(tv.hWnd, LVM_SETITEMCOUNT, uintptr(count), 0) {
		return newError("SendMessage(LVM_SETITEMCOUNT)")
	}

	return nil
}

// CheckBoxes returns if the *TableView has check boxes.
func (tv *TableView) CheckBoxes() bool {
	return SendMessage(tv.hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)&LVS_EX_CHECKBOXES > 0
}

// SetCheckBoxes sets if the *TableView has check boxes.
func (tv *TableView) SetCheckBoxes(value bool) {
	exStyle := SendMessage(tv.hWnd, LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	oldStyle := exStyle
	if value {
		exStyle |= LVS_EX_CHECKBOXES
	} else {
		exStyle &^= LVS_EX_CHECKBOXES
	}
	if exStyle != oldStyle {
		SendMessage(tv.hWnd, LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)
	}

	mask := SendMessage(tv.hWnd, LVM_GETCALLBACKMASK, 0, 0)

	if value {
		mask |= LVIS_STATEIMAGEMASK
	} else {
		mask &^= LVIS_STATEIMAGEMASK
	}

	if FALSE == SendMessage(tv.hWnd, LVM_SETCALLBACKMASK, mask, 0) {
		newError("SendMessage(LVM_SETCALLBACKMASK)")
	}
}

func (tv *TableView) selectedColumnIndex() int {
	return int(SendMessage(tv.hWnd, LVM_GETSELECTEDCOLUMN, 0, 0))
}

func (tv *TableView) setSelectedColumnIndex(value int) {
	SendMessage(tv.hWnd, LVM_SETSELECTEDCOLUMN, uintptr(value), 0)
}

func (tv *TableView) setSortIcon(index int, order SortOrder) error {
	headerHwnd := HWND(SendMessage(tv.hWnd, LVM_GETHEADER, 0, 0))

	count := len(tv.model.Columns())

	for i := 0; i < count; i++ {
		item := HDITEM{
			Mask: HDI_FORMAT,
		}

		iPtr := uintptr(i)
		itemPtr := uintptr(unsafe.Pointer(&item))

		if SendMessage(headerHwnd, HDM_GETITEM, iPtr, itemPtr) == 0 {
			return newError("SendMessage(HDM_GETITEM)")
		}

		if i == index {
			switch order {
			case SortAscending:
				item.Fmt &^= HDF_SORTDOWN
				item.Fmt |= HDF_SORTUP

			case SortDescending:
				item.Fmt &^= HDF_SORTUP
				item.Fmt |= HDF_SORTDOWN
			}
		} else {
			item.Fmt &^= HDF_SORTDOWN | HDF_SORTUP
		}

		if SendMessage(headerHwnd, HDM_SETITEM, iPtr, itemPtr) == 0 {
			return newError("SendMessage(HDM_SETITEM)")
		}
	}

	return nil
}

// ColumnClicked returns the event that is published after a column header was
// clicked.
func (tv *TableView) ColumnClicked() *IntEvent {
	return tv.columnClickedPublisher.Event()
}

// ItemActivated returns the event that is published after an item was 
// activated.
//
// An item is activated when it is double clicked or the enter key is pressed
// when the item is selected. 
func (tv *TableView) ItemActivated() *Event {
	return tv.itemActivatedPublisher.Event()
}

// CurrentIndex returns the index of the current item, or -1 if there is no
// current item.
func (tv *TableView) CurrentIndex() int {
	return tv.currentIndex
}

// SetCurrentIndex sets the index of the current item.
//
// Call this with a value of -1 to have no current item.
func (tv *TableView) SetCurrentIndex(value int) error {
	var lvi LVITEM

	lvi.StateMask = LVIS_FOCUSED | LVIS_SELECTED
	if value > -1 {
		lvi.State = LVIS_FOCUSED | LVIS_SELECTED
	}

	if FALSE == SendMessage(tv.hWnd, LVM_SETITEMSTATE, uintptr(value), uintptr(unsafe.Pointer(&lvi))) {
		return newError("SendMessage(LVM_SETITEMSTATE)")
	}

	if value != -1 {
		if FALSE == SendMessage(tv.hWnd, LVM_ENSUREVISIBLE, uintptr(value), uintptr(0)) {
			return newError("SendMessage(LVM_ENSUREVISIBLE)")
		}
	}

	tv.currentIndex = value

	if value == -1 {
		tv.currentIndexChangedPublisher.Publish()
	}

	return nil
}

// CurrentIndexChanged is the event that is published after CurrentIndex has
// changed.
func (tv *TableView) CurrentIndexChanged() *Event {
	return tv.currentIndexChangedPublisher.Event()
}

// SingleItemSelection returns if only a single item can be selected at once.
//
// By default multiple items can be selected at once.
func (tv *TableView) SingleItemSelection() bool {
	style := uint(GetWindowLong(tv.hWnd, GWL_STYLE))
	if style == 0 {
		lastError("GetWindowLong")
		return false
	}

	return style&LVS_SINGLESEL > 0
}

// SetSingleItemSelection sets if only a single item can be selected at once.
func (tv *TableView) SetSingleItemSelection(value bool) error {
	return tv.ensureStyleBits(LVS_SINGLESEL, value)
}

// SelectedIndexes returns a list of the currently selected item indexes.
func (tv *TableView) SelectedIndexes() *IndexList {
	return tv.selectedIndexes
}

// ItemStateChangedEventDelay returns the delay in milliseconds, between the
// moment the state of an item in the *TableView changes and the moment the
// associated event is published.
//
// By default there is no delay.
func (tv *TableView) ItemStateChangedEventDelay() int {
	return tv.itemStateChangedEventDelay
}

// SetItemStateChangedEventDelay sets the delay in milliseconds, between the
// moment the state of an item in the *TableView changes and the moment the
// associated event is published.
//
// An example where this may be useful is a master-details scenario. If the 
// master TableView is configured to delay the event, you can avoid pointless
// updates of the details TableView, if the user uses arrow keys to rapidly
// navigate the master view.
func (tv *TableView) SetItemStateChangedEventDelay(delay int) {
	tv.itemStateChangedEventDelay = delay
}

func (tv *TableView) updateSelectedIndexes() {
	count := int(SendMessage(tv.hWnd, LVM_GETSELECTEDCOUNT, 0, 0))
	indexes := make([]int, count)

	j := -1
	for i := 0; i < count; i++ {
		j = int(SendMessage(tv.hWnd, LVM_GETNEXTITEM, uintptr(j), LVNI_SELECTED))
		indexes[i] = j
	}

	changed := len(indexes) != len(tv.selectedIndexes.items)
	if !changed {
		for i := 0; i < len(indexes); i++ {
			if indexes[i] != tv.selectedIndexes.items[i] {
				changed = true
				break
			}
		}
	}

	if changed {
		tv.selectedIndexes.items = indexes
		if tv.itemStateChangedEventDelay > 0 {
			if 0 == SetTimer(
				tv.hWnd,
				tableViewSelectedIndexesChangedTimerId,
				uint32(tv.itemStateChangedEventDelay),
				0) {

				lastError("SetTimer")
			}
		} else {
			tv.selectedIndexesChangedPublisher.Publish()
		}
	}
}

// SelectedIndexesChanged returns the event that is published when the list of
// selected item indexes changed.
func (tv *TableView) SelectedIndexesChanged() *Event {
	return tv.selectedIndexesChangedPublisher.Event()
}

// LastColumnStretched returns if the last column should take up all remaining
// horizontal space of the *TableView.
func (tv *TableView) LastColumnStretched() bool {
	return tv.lastColumnStretched
}

// SetLastColumnStretched sets if the last column should take up all remaining
// horizontal space of the *TableView.
//
// The effect of setting this is persistent.
func (tv *TableView) SetLastColumnStretched(value bool) error {
	if value {
		if err := tv.StretchLastColumn(); err != nil {
			return err
		}
	}

	tv.lastColumnStretched = value

	return nil
}

// StretchLastColumn makes the last column take up all remaining horizontal 
// space of the *TableView.
//
// The effect of this is not persistent.
func (tv *TableView) StretchLastColumn() error {
	colCount := len(tv.columns)
	if colCount == 0 {
		return nil
	}

	if 0 == SendMessage(tv.hWnd, LVM_SETCOLUMNWIDTH, uintptr(colCount-1), LVSCW_AUTOSIZE_USEHEADER) {
		return newError("LVM_SETCOLUMNWIDTH failed")
	}

	return nil
}

// Persistent returns if the *TableView should persist its UI state, like column
// widths. See *App.Settings for details.
func (tv *TableView) Persistent() bool {
	return tv.persistent
}

// SetPersistent sets if the *TableView should persist its UI state, like column
// widths. See *App.Settings for details.
func (tv *TableView) SetPersistent(value bool) {
	tv.persistent = value
}

// SaveState writes the UI state of the *TableView to the settings.
func (tv *TableView) SaveState() error {
	buf := bytes.NewBuffer(nil)

	count := len(tv.columns)
	for i := 0; i < count; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}

		width := SendMessage(tv.hWnd, LVM_GETCOLUMNWIDTH, uintptr(i), 0)
		if width == 0 {
			width = 100
		}

		buf.WriteString(strconv.Itoa(int(width)))
	}

	buf.WriteString(";")

	indices := make([]int32, count)
	lParam := uintptr(unsafe.Pointer(&indices[0]))

	SendMessage(tv.hWnd, LVM_GETCOLUMNORDERARRAY, uintptr(count), lParam)

	for i, idx := range indices {
		if i > 0 {
			buf.WriteString(" ")
		}

		buf.WriteString(strconv.Itoa(int(idx)))
	}

	return tv.putState(buf.String())
}

// RestoreState restores the UI state of the *TableView from the settings.
func (tv *TableView) RestoreState() error {
	state, err := tv.getState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	parts := strings.Split(state, ";")

	widthStrs := strings.Split(parts[0], " ")

	// FIXME: Solve this in a better way.
	if len(widthStrs) != len(tv.columns) {
		log.Print("*TableView.RestoreState: failed due to unexpected column count (FIXME!)")
		return nil
	}

	tv.SetSuspended(true)
	defer tv.SetSuspended(false)

	for i, str := range widthStrs {
		width, err := strconv.Atoi(str)
		if err != nil {
			return err
		}

		if FALSE == SendMessage(tv.hWnd, LVM_SETCOLUMNWIDTH, uintptr(i), uintptr(width)) {
			return newError("LVM_SETCOLUMNWIDTH failed")
		}
	}

	if len(parts) > 1 {
		indexStrs := strings.Split(parts[1], " ")

		indices := make([]int32, len(indexStrs))

		var failed bool
		for i, s := range indexStrs {
			idx, err := strconv.Atoi(s)
			if err != nil {
				failed = true
				break
			}
			indices[i] = int32(idx)
		}

		if !failed {
			wParam := uintptr(len(indices))
			lParam := uintptr(unsafe.Pointer(&indices[0]))
			SendMessage(tv.hWnd, LVM_SETCOLUMNORDERARRAY, wParam, lParam)
		}
	}

	return nil
}

func (tv *TableView) applyImageList(image interface{}) {
	var himl HIMAGELIST

	if filePath, ok := image.(string); ok {
		_, himl = tv.iconIndexAndHIml(filePath)
	} else {
		imgSize := Size{
			int(GetSystemMetrics(SM_CXSMICON)),
			int(GetSystemMetrics(SM_CYSMICON)),
		}

		var err error
		if tv.imageList, err = NewImageList(imgSize, 0); err != nil {
			return
		}

		himl = tv.imageList.hIml
	}

	if himl != 0 {
		SendMessage(tv.hWnd, LVM_SETIMAGELIST, LVSIL_SMALL, uintptr(himl))

		tv.hasAppliedImageList = true
	}
}

func (tv *TableView) iconIndexAndHIml(filePath string) (int32, HIMAGELIST) {
	var shfi SHFILEINFO

	if hIml := HIMAGELIST(SHGetFileInfo(
		syscall.StringToUTF16Ptr(filePath),
		0,
		&shfi,
		uint32(unsafe.Sizeof(shfi)),
		SHGFI_SYSICONINDEX|SHGFI_SMALLICON)); hIml != 0 {

		return shfi.IIcon, hIml
	}

	return -1, 0
}

func (tv *TableView) imageIndex(image interface{}) int32 {
	imageIndex := int32(-1)

	if image != nil {
		var ptr uintptr
		switch img := image.(type) {
		case *Bitmap:
			ptr = uintptr(unsafe.Pointer(img))

		case *Icon:
			ptr = uintptr(unsafe.Pointer(img))
		}

		if ptr == 0 {
			return -1
		}

		if imageIndex, ok := tv.imageUintptr2Index[ptr]; ok {
			return imageIndex
		}

		switch img := image.(type) {
		case *Bitmap:
			imageIndex = ImageList_AddMasked(tv.imageList.hIml, img.hBmp, 0)

		case *Icon:
			imageIndex = ImageList_ReplaceIcon(tv.imageList.hIml, -1, img.hIcon)
		}

		if imageIndex > -1 {
			tv.imageUintptr2Index[ptr] = imageIndex
		}
	}

	return imageIndex
}

func (tv *TableView) toggleItemChecked(index int) error {
	checked := tv.itemChecker.Checked(index)

	if err := tv.itemChecker.SetChecked(index, !checked); err != nil {
		return wrapError(err)
	}

	if FALSE == SendMessage(tv.hWnd, LVM_UPDATE, uintptr(index), 0) {
		return newError("SendMessage(LVM_UPDATE)")
	}

	return nil
}

func (tv *TableView) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_ERASEBKGND:
		if tv.lastColumnStretched && !tv.inEraseBkgnd {
			tv.inEraseBkgnd = true
			defer func() {
				tv.inEraseBkgnd = false
			}()
			tv.StretchLastColumn()
		}
		return 1

	case WM_GETDLGCODE:
		if wParam == VK_RETURN {
			return DLGC_WANTALLKEYS
		}

	case WM_LBUTTONDOWN, WM_RBUTTONDOWN, WM_LBUTTONDBLCLK, WM_RBUTTONDBLCLK:
		var hti LVHITTESTINFO
		hti.Pt = POINT{GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam)}
		SendMessage(tv.hWnd, LVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))

		if hti.Flags == LVHT_NOWHERE && tv.SingleItemSelection() {
			// We keep the current item, if in single item selection mode. 
			tv.SetFocus()
			return 0
		}

		switch msg {
		case WM_LBUTTONDOWN, WM_RBUTTONDOWN:
			if hti.Flags == LVHT_ONITEMSTATEICON &&
				tv.itemChecker != nil &&
				tv.CheckBoxes() {

				tv.toggleItemChecked(int(hti.IItem))
			}
		}

	case WM_KEYDOWN:
		if wParam == VK_SPACE &&
			tv.currentIndex > -1 &&
			tv.itemChecker != nil &&
			tv.CheckBoxes() {

			tv.toggleItemChecked(tv.currentIndex)
		}

	case WM_NOTIFY:
		switch int(((*NMHDR)(unsafe.Pointer(lParam))).Code) {
		case LVN_GETDISPINFO:
			di := (*NMLVDISPINFO)(unsafe.Pointer(lParam))

			row := int(di.Item.IItem)
			col := int(di.Item.ISubItem)

			if di.Item.Mask&LVIF_TEXT > 0 {
				var text string
				switch val := tv.model.Value(row, col).(type) {
				case string:
					text = val

				case float32:
					prec := tv.columns[col].Precision
					if prec == 0 {
						prec = 2
					}
					text, _ = formatFloat(float64(val), prec)

				case float64:
					prec := tv.columns[col].Precision
					if prec == 0 {
						prec = 2
					}
					text, _ = formatFloat(val, prec)

				case time.Time:
					text = val.Format(tv.columns[col].Format)

				case *big.Rat:
					prec := tv.columns[col].Precision
					if prec == 0 {
						prec = 2
					}
					text, _ = formatRat(val, prec)

				default:
					text = fmt.Sprintf(tv.columns[col].Format, val)
				}

				utf16 := syscall.StringToUTF16(text)
				buf := (*[256]uint16)(unsafe.Pointer(di.Item.PszText))
				max := mini(len(utf16), int(di.Item.CchTextMax))
				copy((*buf)[:], utf16[:max])
			}

			if tv.imageProvider != nil && di.Item.Mask&LVIF_IMAGE > 0 {
				image := tv.imageProvider.Image(row)

				if !tv.hasAppliedImageList {
					tv.applyImageList(image)
				}

				if tv.hasAppliedImageList {
					if tv.imageList != nil {
						di.Item.IImage = int32(tv.imageIndex(image))
					} else if filePath, ok := image.(string); ok {
						if iIcon, ok := tv.filePath2IconIndex[filePath]; ok {
							di.Item.IImage = iIcon
							break
						}

						if iIcon, _ := tv.iconIndexAndHIml(filePath); iIcon != -1 {
							tv.filePath2IconIndex[filePath] = iIcon

							di.Item.IImage = iIcon
						}
					}
				}
			}

			if di.Item.StateMask&LVIS_STATEIMAGEMASK > 0 &&
				tv.itemChecker != nil {
				checked := tv.itemChecker.Checked(row)

				if checked {
					di.Item.State = 0x2000
				} else {
					di.Item.State = 0x1000
				}
			}

		case NM_CUSTOMDRAW:
			if tv.alternatingRowBGColor != defaultTVRowBGColor {
				nmlvcd := (*NMLVCUSTOMDRAW)(unsafe.Pointer(lParam))

				switch nmlvcd.Nmcd.DwDrawStage {
				case CDDS_PREPAINT:
					return CDRF_NOTIFYITEMDRAW

				case CDDS_ITEMPREPAINT:
					if nmlvcd.Nmcd.DwItemSpec%2 == 1 {
						nmlvcd.ClrTextBk = COLORREF(tv.alternatingRowBGColor)
					}
				}
			}

			return CDRF_DODEFAULT

		case LVN_COLUMNCLICK:
			nmlv := (*NMLISTVIEW)(unsafe.Pointer(lParam))

			col := int(nmlv.ISubItem)
			tv.columnClickedPublisher.Publish(col)

			if sorter, ok := tv.model.(Sorter); ok && sorter.ColumnSortable(col) {
				prevCol := sorter.SortedColumn()
				var order SortOrder
				if col != prevCol || sorter.SortOrder() == SortDescending {
					order = SortAscending
				} else {
					order = SortDescending
				}
				sorter.Sort(col, order)
			}

		case LVN_ITEMCHANGED:
			nmlv := (*NMLISTVIEW)(unsafe.Pointer(lParam))
			selectedNow := nmlv.UNewState&LVIS_SELECTED > 0
			selectedBefore := nmlv.UOldState&LVIS_SELECTED > 0
			if selectedNow && !selectedBefore {
				tv.currentIndex = int(nmlv.IItem)
				if tv.itemStateChangedEventDelay > 0 {
					if 0 == SetTimer(
						tv.hWnd,
						tableViewCurrentIndexChangedTimerId,
						uint32(tv.itemStateChangedEventDelay),
						0) {

						lastError("SetTimer")
					}
				} else {
					tv.currentIndexChangedPublisher.Publish()
				}
			}
			if !tv.SingleItemSelection() {
				tv.updateSelectedIndexes()
			}

		case LVN_ITEMACTIVATE:
			tv.itemActivatedPublisher.Publish()
		}

	case WM_TIMER:
		switch wParam {
		case tableViewCurrentIndexChangedTimerId:
			tv.currentIndexChangedPublisher.Publish()

		case tableViewSelectedIndexesChangedTimerId:
			tv.selectedIndexesChangedPublisher.Publish()
		}
	}

	return tv.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
