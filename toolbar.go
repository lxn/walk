// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

type ToolBar struct {
	WidgetBase
	imageList          *ImageList
	actions            *ActionList
	defaultButtonWidth int
	maxTextRows        int
}

func newToolBar(parent Container, style uint32) (*ToolBar, error) {
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

func NewToolBar(parent Container) (*ToolBar, error) {
	return newToolBar(parent, TBSTYLE_WRAPABLE)
}

func NewVerticalToolBar(parent Container) (*ToolBar, error) {
	tb, err := newToolBar(parent, CCS_VERT|CCS_NORESIZE)
	if err != nil {
		return nil, err
	}

	tb.defaultButtonWidth = 100

	return tb, nil
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

func (tb *ToolBar) MinSizeHint() Size {
	return tb.SizeHint()
}

func (tb *ToolBar) SizeHint() Size {
	if tb.actions.Len() == 0 {
		return Size{}
	}

	size := uint32(SendMessage(tb.hWnd, TB_GETBUTTONSIZE, 0, 0))

	width := tb.defaultButtonWidth
	if width == 0 {
		width = int(LOWORD(size))
	}

	height := int(HIWORD(size))

	return Size{width, height}
}

func (tb *ToolBar) applyDefaultButtonWidth() error {
	if tb.defaultButtonWidth == 0 {
		return nil
	}

	lParam := uintptr(
		MAKELONG(uint16(tb.defaultButtonWidth), uint16(tb.defaultButtonWidth)))
	if 0 == SendMessage(tb.hWnd, TB_SETBUTTONWIDTH, 0, lParam) {
		return newError("SendMessage(TB_SETBUTTONWIDTH)")
	}

	size := uint32(SendMessage(tb.hWnd, TB_GETBUTTONSIZE, 0, 0))
	height := HIWORD(size)

	lParam = uintptr(MAKELONG(uint16(tb.defaultButtonWidth), height))
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
func (tb *ToolBar) SetDefaultButtonWidth(width int) error {
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

func (tb *ToolBar) MaxTextRows() int {
	return tb.maxTextRows
}

func (tb *ToolBar) SetMaxTextRows(maxTextRows int) error {
	if 0 == SendMessage(tb.hWnd, TB_SETMAXTEXTROWS, uintptr(maxTextRows), 0) {
		return newError("SendMessage(TB_SETMAXTEXTROWS)")
	}

	tb.maxTextRows = maxTextRows

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

func (tb *ToolBar) imageIndex(image *Bitmap) (imageIndex int32, err error) {
	imageIndex = -1
	if image != nil {
		// FIXME: Protect against duplicate insertion
		if imageIndex, err = tb.imageList.AddMasked(image); err != nil {
			return
		}
	}

	return
}

func (tb *ToolBar) wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_NOTIFY:
		nmm := (*NMMOUSE)(unsafe.Pointer(lParam))

		switch int32(nmm.Hdr.Code) {
		case NM_CLICK:
			actionId := uint16(nmm.DwItemSpec)
			if action := actionsById[actionId]; action != nil {
				action.raiseTriggered()
			}
		}
	}

	return tb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}

func (tb *ToolBar) initButtonForAction(action *Action, state, style *byte, image *int32, text *uintptr) (err error) {
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

func (tb *ToolBar) onActionChanged(action *Action) error {
	tbbi := TBBUTTONINFO{
		DwMask: TBIF_IMAGE | TBIF_STATE | TBIF_STYLE | TBIF_TEXT,
	}

	tbbi.CbSize = uint32(unsafe.Sizeof(tbbi))

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

func (tb *ToolBar) onInsertingAction(index int, action *Action) error {
	tbb := TBBUTTON{
		IdCommand: int32(action.id),
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

func (tb *ToolBar) removeAt(index int) error {
	action := tb.actions.At(index)
	action.removeChangedHandler(tb)

	if 0 == SendMessage(tb.hWnd, TB_DELETEBUTTON, uintptr(index), 0) {
		return newError("SendMessage(TB_DELETEBUTTON) failed")
	}

	return nil
}

func (tb *ToolBar) onRemovingAction(index int, action *Action) error {
	return tb.removeAt(index)
}

func (tb *ToolBar) onClearingActions() error {
	for i := tb.actions.Len() - 1; i >= 0; i-- {
		if err := tb.removeAt(i); err != nil {
			return err
		}
	}

	return nil
}
