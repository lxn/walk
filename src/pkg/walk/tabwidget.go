// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

import . "walk/winapi"

const tabWidgetWindowClass = `\o/ Walk_TabWidget_Class \o/`

var tabWidgetWindowClassRegistered bool

type TabWidget struct {
	WidgetBase
	hWndTab                      HWND
	pages                        *TabPageList
	currentIndex                 int
	currentIndexChangedPublisher EventPublisher
	persistent                   bool
}

func NewTabWidget(parent Container) (*TabWidget, os.Error) {
	ensureRegisteredWindowClass(tabWidgetWindowClass, &tabWidgetWindowClassRegistered)

	tw := &TabWidget{currentIndex: -1}
	tw.pages = newTabPageList(tw)

	if err := initChildWidget(
		tw,
		parent,
		tabWidgetWindowClass,
		WS_VISIBLE,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			tw.Dispose()
		}
	}()

	tw.SetPersistent(true)

	tw.hWndTab = CreateWindowEx(
		0, syscall.StringToUTF16Ptr("SysTabControl32"), nil,
		WS_CHILD|WS_CLIPSIBLINGS|WS_TABSTOP|WS_VISIBLE,
		0, 0, 0, 0, tw.hWnd, 0, 0, nil)
	if tw.hWndTab == 0 {
		return nil, lastError("CreateWindowEx")
	}
	SendMessage(tw.hWndTab, WM_SETFONT, uintptr(defaultFont.handleForDPI(0)), 1)

	succeeded = true

	return tw, nil
}

func (tw *TabWidget) LayoutFlags() LayoutFlags {
	if tw.pages.Len() == 0 {
		return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
	}

	var flags LayoutFlags

	for i := tw.pages.Len() - 1; i >= 0; i-- {
		flags |= tw.pages.At(i).LayoutFlags()
	}

	return flags
}

func (tw *TabWidget) MinSizeHint() Size {
	if tw.pages.Len() == 0 {
		return tw.SizeHint()
	}

	var min Size

	for i := tw.pages.Len() - 1; i >= 0; i-- {
		s := tw.pages.At(i).MinSizeHint()

		min.Width = maxi(min.Width, s.Width)
		min.Height = maxi(min.Height, s.Height)
	}

	b := tw.Bounds()
	pb := tw.pages.At(0).Bounds()

	size := Size{b.Width - pb.Width + min.Width, b.Height - pb.Height + min.Height}

	return size

}

func (tw *TabWidget) SizeHint() Size {
	return Size{100, 100}
}

func (tw *TabWidget) CurrentIndex() int {
	return tw.currentIndex
}

func (tw *TabWidget) SetCurrentIndex(index int) os.Error {
	if index == tw.currentIndex {
		return nil
	}

	if index < 0 || index >= tw.pages.Len() {
		return newError("invalid index")
	}

	ret := int(SendMessage(tw.hWndTab, TCM_SETCURSEL, uintptr(index), 0))
	if ret == -1 {
		return newError("SendMessage(TCM_SETCURSEL) failed")
	}

	// FIXME: The SendMessage(TCM_SETCURSEL) call above doesn't cause a
	// TCN_SELCHANGE notification, so we use this workaround.
	tw.onSelChange()

	return nil
}

func (tw *TabWidget) CurrentIndexChanged() *Event {
	return tw.currentIndexChangedPublisher.Event()
}

func (tw *TabWidget) Pages() *TabPageList {
	return tw.pages
}

func (tw *TabWidget) Persistent() bool {
	return tw.persistent
}

func (tw *TabWidget) SetPersistent(value bool) {
	tw.persistent = value
}

func (tw *TabWidget) SaveState() os.Error {
	tw.putState(strconv.Itoa(tw.CurrentIndex()))

	for _, page := range tw.pages.items {
		if err := page.SaveState(); err != nil {
			return err
		}
	}

	return nil
}

func (tw *TabWidget) RestoreState() os.Error {
	state, err := tw.getState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	index, err := strconv.Atoi(state)
	if err != nil {
		return err
	}
	if index >= 0 && index < tw.pages.Len() {
		if err := tw.SetCurrentIndex(index); err != nil {
			return err
		}
	}

	for _, page := range tw.pages.items {
		if err := page.RestoreState(); err != nil {
			return err
		}
	}

	return nil
}

func (tw *TabWidget) resizePages() {
	var r RECT
	if !GetWindowRect(tw.hWndTab, &r) {
		lastError("GetWindowRect")
		return
	}

	p := POINT{r.Left, r.Top}
	if !ScreenToClient(tw.hWnd, &p) {
		newError("ScreenToClient failed")
		return
	}

	r = RECT{p.X, p.Y, r.Right - r.Left + p.X, r.Bottom - r.Top + p.Y}

	SendMessage(tw.hWndTab, TCM_ADJUSTRECT, 0, uintptr(unsafe.Pointer(&r)))

	for _, page := range tw.pages.items {
		if err := page.SetBounds(Rectangle{r.Left - 2, r.Top, r.Right - r.Left + 2, r.Bottom - r.Top}); err != nil {
			return
		}
	}
}

func (tw *TabWidget) onResize(lParam uintptr) {
	r := RECT{0, 0, GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam)}
	if !MoveWindow(tw.hWndTab, r.Left, r.Top, r.Right-r.Left, r.Bottom-r.Top, true) {
		lastError("MoveWindow")
		return
	}

	tw.resizePages()
}

func (tw *TabWidget) onSelChange() {
	if tw.currentIndex != -1 {
		page := tw.pages.At(tw.currentIndex)
		page.SetVisible(false)
	}

	tw.currentIndex = int(SendMessage(tw.hWndTab, TCM_GETCURSEL, 0, 0))

	if tw.currentIndex != -1 {
		page := tw.pages.At(tw.currentIndex)
		page.SetVisible(true)
		page.Invalidate()
	}

	tw.currentIndexChangedPublisher.Publish()
}

func (tw *TabWidget) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	if tw.hWndTab != 0 {
		switch msg {
		case WM_SIZE, WM_SIZING:
			tw.onResize(lParam)

		case WM_NOTIFY:
			nmhdr := (*NMHDR)(unsafe.Pointer(lParam))

			switch int(nmhdr.Code) {
			case TCN_SELCHANGE:
				tw.onSelChange()
			}
		}
	}

	return tw.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}

func (tw *TabWidget) onPageChanged(page *TabPage) (err os.Error) {
	index := tw.pages.Index(page)
	item := page.tcItem()

	if 0 == SendMessage(tw.hWndTab, TCM_SETITEM, uintptr(index), uintptr(unsafe.Pointer(item))) {
		return newError("SendMessage(TCM_SETITEM) failed")
	}

	return nil
}

func (tw *TabWidget) onInsertingPage(index int, page *TabPage) (err os.Error) {
	return nil
}

func (tw *TabWidget) onInsertedPage(index int, page *TabPage) (err os.Error) {
	page.SetVisible(false)

	style := uint(GetWindowLong(page.hWnd, GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	style |= WS_CHILD
	style &^= WS_POPUP

	SetLastError(0)
	if SetWindowLong(page.hWnd, GWL_STYLE, int(style)) == 0 {
		return lastError("SetWindowLong")
	}

	if SetParent(page.hWnd, tw.hWnd) == 0 {
		return lastError("SetParent")
	}

	if tw.pages.Len() == 1 {
		page.SetVisible(true)
		tw.currentIndex = 0
	}

	item := page.tcItem()

	if idx := int(SendMessage(tw.hWndTab, TCM_INSERTITEM, uintptr(index), uintptr(unsafe.Pointer(item)))); idx == -1 {
		return newError("SendMessage(TCM_INSERTITEM) failed")
	}

	tw.resizePages()

	page.tabWidget = tw

	return
}

func (tw *TabWidget) removePage(page *TabPage) (err os.Error) {
	page.SetVisible(false)

	style := uint(GetWindowLong(page.hWnd, GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	style &^= WS_CHILD
	style |= WS_POPUP

	SetLastError(0)
	if SetWindowLong(page.hWnd, GWL_STYLE, int(style)) == 0 {
		return lastError("SetWindowLong")
	}

	page.tabWidget = nil

	return page.SetParent(nil)
}

func (tw *TabWidget) onRemovingPage(index int, page *TabPage) (err os.Error) {
	return nil
}

func (tw *TabWidget) onRemovedPage(index int, page *TabPage) (err os.Error) {
	err = tw.removePage(page)
	if err != nil {
		return
	}

	SendMessage(tw.hWndTab, TCM_DELETEITEM, uintptr(index), 0)

	tw.currentIndex = -1
	tw.onSelChange()

	return

	// FIXME: Either make use of this unreachable code or remove it.
	if index == tw.currentIndex {
		// removal of current visible tabpage...
		tw.currentIndex = -1

		// select new tabpage if any :
		if tw.pages.Len() > 0 {
			// are we removing the rightmost page ? 
			if index == tw.pages.Len()-1 {
				// If so, select the page on the left 
				index -= 1
			}
		}
	}

	tw.SetCurrentIndex(index)
	//tw.Invalidate()

	return
}

func (tw *TabWidget) onClearingPages(pages []*TabPage) (err os.Error) {
	return nil
}

func (tw *TabWidget) onClearedPages(pages []*TabPage) (err os.Error) {
	SendMessage(tw.hWndTab, TCM_DELETEALLITEMS, 0, 0)
	for _, page := range pages {
		tw.removePage(page)
	}
	tw.currentIndex = -1
	return nil
}
