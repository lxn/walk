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

var toolBarOrigWndProcPtr uintptr
var _ subclassedWidget = &ToolBar{}

type ToolBar struct {
	WidgetBase
	imageList          *ImageList
	actions            *ActionList
	defaultButtonWidth int
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
	tb, err := newToolBar(parent, CCS_VERT|CCS_NORESIZE)
	if err != nil {
		return nil, err
	}

	tb.defaultButtonWidth = 100

	return tb, nil
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
		return ShrinkableVert | GrowableVert | GreedyVert
	}

	// FIXME: Since reimplementation of BoxLayout we must return 0 here,
	// otherwise the ToolBar contained in MainWindow will eat half the space.  
	return 0 //ShrinkableHorz | GrowableHorz
}

func (tb *ToolBar) SizeHint() Size {
	if tb.actions.Len() == 0 {
		return Size{}
	}

	size := uint(SendMessage(tb.hWnd, TB_GETBUTTONSIZE, 0, 0))

	width := tb.defaultButtonWidth
	if width == 0 {
		width = int(LOWORD(size))
	}

	height := int(HIWORD(size))

	return Size{width, height}
}

func (tb *ToolBar) applyDefaultButtonWidth() os.Error {
	if tb.defaultButtonWidth == 0 {
		return nil
	}

	size := uint(SendMessage(tb.hWnd, TB_GETBUTTONSIZE, 0, 0))
	height := HIWORD(size)

	lParam := uintptr(MAKELONG(uint16(tb.defaultButtonWidth), height))
	if FALSE == SendMessage(tb.hWnd, TB_SETBUTTONSIZE, 0, lParam) {
		return newError("SendMessage(TB_SETBUTTONSIZE)")
	}

	return nil
}

// DefaultButtonWidth returns the default button width of the ToolBar.
//
// The default value for a horizontal ToolBar is 0, resulting in automatic
// sizing behavior. For a vertical ToolBar, the default is 100 pixels. 
func (tb *ToolBar) DefaultButtonWidth() int {
	return tb.defaultButtonWidth
}

// SetDefaultButtonWidth sets the default button width of the ToolBar.
//
// Calling this method affects all buttons in the ToolBar, no matter if they are 
// added before or after the call. A width of 0 results in automatic sizing 
// behavior. Negative values are not allowed.
func (tb *ToolBar) SetDefaultButtonWidth(width int) os.Error {
	if width == tb.defaultButtonWidth {
		return nil
	}

	if width < 0 {
		return newError("width must be >= 0")
	}

	old := tb.defaultButtonWidth

	tb.defaultButtonWidth = width

	for _, action := range tb.actions.actions {
		if err := tb.onActionChanged(action); err != nil {
			tb.defaultButtonWidth = old

			return err
		}
	}

	return tb.applyDefaultButtonWidth()
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
		if imageIndex, err = tb.imageList.AddMasked(image); err != nil {
			return
		}
	}

	return
}

func (tb *ToolBar) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_NOTIFY:
		nmm := (*NMMOUSE)(unsafe.Pointer(lParam))

		switch int(nmm.Hdr.Code) {
		case NM_CLICK:
			actionId := uint16(nmm.DwItemSpec)
			if action := actionsById[actionId]; action != nil {
				action.raiseTriggered()
			}
		}
	}

	return tb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}

func (tb *ToolBar) initButtonForAction(action *Action, state, style *byte, image *int, text *uintptr) (err os.Error) {
	if tb.hasStyleBits(CCS_VERT) {
		*state |= TBSTATE_WRAP
	} else if tb.defaultButtonWidth == 0 {
		*style |= BTNS_AUTOSIZE
	}

	if action.checked {
		*state |= TBSTATE_CHECKED
	}

	if action.enabled {
		*state |= TBSTATE_ENABLED
	}

	if action.checkable {
		*style |= BTNS_CHECK
	}

	if action.exclusive {
		*style |= BTNS_GROUP
	}

	if *image, err = tb.imageIndex(action.image); err != nil {
		return
	}

	*text = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(action.Text())))

	return
}

func (tb *ToolBar) onActionChanged(action *Action) os.Error {
	tbbi := TBBUTTONINFO{
		DwMask: TBIF_IMAGE | TBIF_STATE | TBIF_STYLE | TBIF_TEXT,
	}

	tbbi.CbSize = uint(unsafe.Sizeof(tbbi))

	if err := tb.initButtonForAction(
		action,
		&tbbi.FsState,
		&tbbi.FsStyle,
		&tbbi.IImage,
		&tbbi.PszText); err != nil {

		return err
	}

	if 0 == SendMessage(
		tb.hWnd,
		TB_SETBUTTONINFO,
		uintptr(action.id),
		uintptr(unsafe.Pointer(&tbbi))) {

		return newError("SendMessage(TB_SETBUTTONINFO) failed")
	}

	return nil
}

func (tb *ToolBar) onInsertingAction(index int, action *Action) os.Error {
	tbb := TBBUTTON{
		IdCommand: int(action.id),
	}

	if err := tb.initButtonForAction(
		action,
		&tbb.FsState,
		&tbb.FsStyle,
		&tbb.IBitmap,
		&tbb.IString); err != nil {

		return err
	}

	tb.SetVisible(true)

	SendMessage(tb.hWnd, TB_BUTTONSTRUCTSIZE, uintptr(unsafe.Sizeof(tbb)), 0)

	if FALSE == SendMessage(tb.hWnd, TB_ADDBUTTONS, 1, uintptr(unsafe.Pointer(&tbb))) {
		return newError("SendMessage(TB_ADDBUTTONS)")
	}

	if err := tb.applyDefaultButtonWidth(); err != nil {
		return err
	}

	SendMessage(tb.hWnd, TB_AUTOSIZE, 0, 0)

	action.addChangedHandler(tb)

	return nil
}

func (tb *ToolBar) removeAt(index int) os.Error {
	action := tb.actions.At(index)
	action.removeChangedHandler(tb)

	if 0 == SendMessage(tb.hWnd, TB_DELETEBUTTON, uintptr(index), 0) {
		return newError("SendMessage(TB_DELETEBUTTON) failed")
	}

	return nil
}

func (tb *ToolBar) onRemovingAction(index int, action *Action) os.Error {
	return tb.removeAt(index)
}

func (tb *ToolBar) onClearingActions() os.Error {
	for i := tb.actions.Len() - 1; i >= 0; i-- {
		if err := tb.removeAt(i); err != nil {
			return err
		}
	}

	return nil
}
