// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/crutches"
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

type ToolBar struct {
	Widget
	imageList      *ImageList
	actions        *ActionList
	minButtonWidth uint16
	maxButtonWidth uint16
}

func newToolBar(parent IContainer, style uint) (*ToolBar, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("ToolbarWindow32"), nil,
		WS_CHILD|WS_VISIBLE|style,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	tb := &ToolBar{Widget: Widget{hWnd: hWnd, parent: parent}}
	tb.actions = newActionList(tb)

	widgetsByHWnd[hWnd] = tb

	//	exStyle := SendMessage(hWnd, TB_GETEXTENDEDSTYLE, 0, 0)
	//	exStyle |= TBSTYLE_EX_DOUBLEBUFFER
	//	SendMessage(hWnd, TB_SETEXTENDEDSTYLE, 0, LPARAM(exStyle))

	parent.Children().Add(tb)

	return tb, nil
}

func NewToolBar(parent IContainer) (*ToolBar, os.Error) {
	return newToolBar(parent, TBSTYLE_WRAPABLE)
}

func NewVerticalToolBar(parent IContainer) (*ToolBar, os.Error) {
	return newToolBar(parent, CCS_VERT|CCS_NORESIZE)
}

func (tb *ToolBar) LayoutFlags() LayoutFlags {
	style := GetWindowLong(tb.hWnd, GWL_STYLE)

	if style&CCS_VERT > 0 {
		return ShrinkVert | GrowVert
	}

	return ShrinkHorz | GrowHorz
}

func (tb *ToolBar) PreferredSize() drawing.Size {
	style := GetWindowLong(tb.hWnd, GWL_STYLE)

	if style&CCS_VERT > 0 && tb.minButtonWidth > 0 {
		return drawing.Size{int(tb.minButtonWidth), 44}
	}

	// FIXME: Figure out how to do this.
	return drawing.Size{44, 44}
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

func (tb *ToolBar) raiseEvent(msg *MSG) (err os.Error) {
	err = tb.Widget.raiseEvent(msg)
	if err != nil {
		return
	}

	switch msg.Message {
	case WM_NOTIFY:
		fmt.Printf("ToolBar.raiseEvent: WM_NOTIFY, msg: %v\n", msg)

		switch ((*NMHDR)(unsafe.Pointer(msg.LParam))).Code {
		case NM_CLICK:
			nmm := (*NMMOUSE)(unsafe.Pointer(msg.LParam))
			fmt.Printf("ToolBar.raiseEvent: WM_NOTIFY, NM_CLICK, nmm: %v\n", nmm)
		}

	case WM_COMMAND, crutches.CommandMsgId():
		switch HIWORD(uint(msg.WParam)) {
		case 0:
			// menu
			actionId := uint16(LOWORD(uint(msg.WParam)))
			if action, ok := actionsById[actionId]; ok {
				action.raiseTriggered()
			}
		}
	}

	return
}

func (tb *ToolBar) onActionChanged(action *Action) (err os.Error) {
	panic("not implemented")
}

func (tb *ToolBar) onInsertingAction(index int, action *Action) (err os.Error) {
	image := action.image
	imageIndex := -1
	if image != nil {
		// FIXME: Protect against duplicate insertion
		imageIndex, err = tb.imageList.AddMasked(image)
		if err != nil {
			return
		}
	}

	tbb := TBBUTTON{
		IBitmap:   imageIndex,
		IdCommand: int(action.id),
		FsState:   TBSTATE_ENABLED | TBSTATE_WRAP,
		FsStyle:/*BTNS_AUTOSIZE |*/ BTNS_BUTTON,
		IString: uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(action.Text()))),
	}

	SendMessage(tb.hWnd, TB_BUTTONSTRUCTSIZE, uintptr(unsafe.Sizeof(tbb)), 0)
	SendMessage(tb.hWnd, TB_ADDBUTTONSW, uintptr(1), uintptr(unsafe.Pointer(&tbb)))
	SendMessage(tb.hWnd, TB_AUTOSIZE, 0, 0)

	return
}

func (tb *ToolBar) onRemovingAction(index int, action *Action) (err os.Error) {
	panic("not implemented")
}

func (tb *ToolBar) onClearingActions() (err os.Error) {
	panic("not implemented")
}
