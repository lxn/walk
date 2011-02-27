// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

var toolBarOrigWndProcPtr uintptr
var _ subclassedWidget = &ToolBar{}

type ToolBar struct {
	WidgetBase
	imageList      *ImageList
	actions        *ActionList
	minButtonWidth uint16
	maxButtonWidth uint16
}

func newToolBar(parent Container, style uint) (*ToolBar, os.Error) {
	tb := &ToolBar{}
	tb.actions = newActionList(tb)

	if err := initChildWidget(
		tb,
		parent,
		"ToolbarWindow32",
		CCS_NODIVIDER|style,
		0); err != nil {
		return nil, err
	}

	return tb, nil
}

func NewToolBar(parent Container) (*ToolBar, os.Error) {
	return newToolBar(parent, TBSTYLE_WRAPABLE)
}

func NewVerticalToolBar(parent Container) (*ToolBar, os.Error) {
	return newToolBar(parent, CCS_VERT|CCS_NORESIZE)
}

func (*ToolBar) origWndProcPtr() uintptr {
	return toolBarOrigWndProcPtr
}

func (*ToolBar) setOrigWndProcPtr(ptr uintptr) {
	toolBarOrigWndProcPtr = ptr
}

func (tb *ToolBar) LayoutFlags() LayoutFlags {
	style := GetWindowLong(tb.hWnd, GWL_STYLE)

	if style&CCS_VERT > 0 {
		return VShrink | VGrow
	}

	// FIXME: Since reimplementation of BoxLayout we must return 0 here,
	// otherwise the ToolBar contained in MainWindow will eat half the space.  
	return 0 //HShrink | HGrow
}

func (tb *ToolBar) PreferredSize() Size {
	if tb.actions.Len() == 0 {
		return Size{}
	}

	style := GetWindowLong(tb.hWnd, GWL_STYLE)

	if style&CCS_VERT > 0 && tb.minButtonWidth > 0 {
		return Size{int(tb.minButtonWidth), 44}
	}

	// FIXME: Figure out how to do this.
	return Size{44, 44}
}

func (tb *ToolBar) ButtonWidthLimits() (min, max uint16) {
	return tb.minButtonWidth, tb.maxButtonWidth
}

func (tb *ToolBar) SetButtonWidthLimits(min, max uint16) os.Error {
	if SendMessage(tb.hWnd, TB_SETBUTTONWIDTH, 0, uintptr(MAKELONG(min, max))) == 0 {
		return newError("TB_SETBUTTONWIDTH failed")
	}

	tb.minButtonWidth = min
	tb.maxButtonWidth = max

	return nil
}

func (tb *ToolBar) Actions() *ActionList {
	return tb.actions
}

func (tb *ToolBar) ImageList() *ImageList {
	return tb.imageList
}

func (tb *ToolBar) SetImageList(value *ImageList) {
	var hIml HIMAGELIST

	if value != nil {
		hIml = value.hIml
	}

	SendMessage(tb.hWnd, TB_SETIMAGELIST, 0, uintptr(hIml))

	tb.imageList = value
}

func (tb *ToolBar) imageIndex(image *Bitmap) (imageIndex int, err os.Error) {
	imageIndex = -1
	if image != nil {
		// FIXME: Protect against duplicate insertion
		imageIndex, err = tb.imageList.AddMasked(image)
		if err != nil {
			return
		}
	}

	return
}

func (tb *ToolBar) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_NOTIFY:
		nmm := (*NMMOUSE)(unsafe.Pointer(lParam))

		switch nmm.Hdr.Code {
		case NM_CLICK:
			actionId := uint16(nmm.DwItemSpec)
			if action := actionsById[actionId]; action != nil {
				action.raiseTriggered()
			}
		}
	}

	return tb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}

func (tb *ToolBar) onActionChanged(action *Action) (err os.Error) {
	imageIndex, err := tb.imageIndex(action.image)
	if err != nil {
		return
	}

	tbbi := TBBUTTONINFO{
		DwMask:  TBIF_BYINDEX | TBIF_IMAGE | TBIF_STATE | TBIF_STYLE | TBIF_TEXT,
		IImage:  imageIndex,
		FsState: TBSTATE_WRAP,
		FsStyle: BTNS_BUTTON,
		PszText: syscall.StringToUTF16Ptr(action.Text()),
	}
	tbbi.CbSize = uint(unsafe.Sizeof(tbbi))
	if action.checked {
		tbbi.FsState |= TBSTATE_CHECKED
	}
	if action.enabled {
		tbbi.FsState |= TBSTATE_ENABLED
	}
	if action.checkable {
		tbbi.FsStyle |= BTNS_CHECK
	}
	if action.exclusive {
		tbbi.FsStyle |= BTNS_GROUP
	}

	if 0 == SendMessage(tb.hWnd, TB_SETBUTTONINFO, uintptr(tb.actions.Index(action)), uintptr(unsafe.Pointer(&tbbi))) {
		err = newError("SendMessage(TB_SETBUTTONINFO) failed")
	}

	return
}

func (tb *ToolBar) onInsertingAction(index int, action *Action) (err os.Error) {
	imageIndex, err := tb.imageIndex(action.image)
	if err != nil {
		return
	}

	tbb := TBBUTTON{
		IBitmap:   imageIndex,
		IdCommand: int(action.id),
		FsState:   TBSTATE_WRAP,
		FsStyle:/*BTNS_AUTOSIZE |*/ BTNS_BUTTON,
		IString: uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(action.Text()))),
	}
	if action.checked {
		tbb.FsState |= TBSTATE_CHECKED
	}
	if action.enabled {
		tbb.FsState |= TBSTATE_ENABLED
	}
	if action.checkable {
		tbb.FsStyle |= BTNS_CHECK
	}
	if action.exclusive {
		tbb.FsStyle |= BTNS_GROUP
	}

	tb.SetVisible(true)

	SendMessage(tb.hWnd, TB_BUTTONSTRUCTSIZE, uintptr(unsafe.Sizeof(tbb)), 0)
	SendMessage(tb.hWnd, TB_ADDBUTTONS, 1, uintptr(unsafe.Pointer(&tbb)))
	SendMessage(tb.hWnd, TB_AUTOSIZE, 0, 0)

	return
}

func (tb *ToolBar) removeAt(index int) (err os.Error) {
	if 0 == SendMessage(tb.hWnd, TB_DELETEBUTTON, uintptr(index), 0) {
		err = newError("SendMessage(TB_DELETEBUTTON) failed")
	}

	return
}

func (tb *ToolBar) onRemovingAction(index int, action *Action) os.Error {
	return tb.removeAt(index)
}

func (tb *ToolBar) onClearingActions() (err os.Error) {
	for i := tb.actions.Len() - 1; i >= 0; i-- {
		err = tb.removeAt(i)
		if err != nil {
			return
		}
	}

	return
}
