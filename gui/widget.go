// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/crutches"
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/gdi32"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

type LayoutFlags byte

const (
	ShrinkHorz LayoutFlags = 1 << iota
	GrowHorz
	ShrinkVert
	GrowVert
)

type IWidget interface {
	Handle() HWND
	Bounds() (drawing.Rectangle, os.Error)
	SetBounds(value drawing.Rectangle) os.Error
	ClientBounds() (drawing.Rectangle, os.Error)
	ContextMenu() *Menu
	SetContextMenu(value *Menu)
	Dispose()
	IsDisposed() bool
	Enabled() (bool, os.Error)
	SetEnabled(value bool) os.Error
	Font() *drawing.Font
	SetFont(value *drawing.Font)
	GroupStart() (bool, os.Error)
	SetGroupStart(value bool) os.Error
	Height() (int, os.Error)
	SetHeight(value int) os.Error
	LayoutFlags() LayoutFlags
	MaxSize() (drawing.Size, os.Error)
	SetMaxSize(value drawing.Size) os.Error
	MinSize() (drawing.Size, os.Error)
	SetMinSize(value drawing.Size) os.Error
	Parent() IContainer
	SetParent(value IContainer) os.Error
	PreferredSize() drawing.Size
	Size() (drawing.Size, os.Error)
	SetSize(value drawing.Size) os.Error
	Text() string
	SetText(value string) os.Error
	Visible() (bool, os.Error)
	SetVisible(value bool) os.Error
	Width() (int, os.Error)
	SetWidth(value int) os.Error
	X() (int, os.Error)
	SetX(value int) os.Error
	Y() (int, os.Error)
	SetY(value int) os.Error
	SetFocus() os.Error
	AddSizeChangedHandler(handler EventHandler)
	RemoveSizeChangedHandler(handler EventHandler)
	RootWidget() RootWidget
	GetDrawingSurface() (*drawing.Surface, os.Error)
}

type widgetInternal interface {
	IWidget
	raiseEvent(msg *MSG) os.Error
}

type Widget struct {
	hWnd                HWND
	parent              IContainer
	font                *drawing.Font
	contextMenu         *Menu
	keyDownHandlers     vector.Vector
	mouseDownHandlers   vector.Vector
	sizeChangedHandlers vector.Vector
	maxSize             drawing.Size
	minSize             drawing.Size
}

var (
	trackedMouseDownHWnd HWND
	widgetsByHWnd        map[HWND]widgetInternal = make(map[HWND]widgetInternal)
)

func rootWidget(w IWidget) RootWidget {
	if w == nil {
		return nil
	}

	for w.Parent() != nil {
		w = w.Parent()
	}

	return widgetsByHWnd[w.Handle()].(RootWidget)
}

func (w *Widget) Handle() HWND {
	return w.hWnd
}

func (w *Widget) Dispose() {
	if w.hWnd != 0 {
		DestroyWindow(w.hWnd)
		w.hWnd = 0
	}
}

func (w *Widget) IsDisposed() bool {
	return w.hWnd == 0
}

func (w *Widget) RootWidget() RootWidget {
	return rootWidget(w)
}

func (w *Widget) ContextMenu() *Menu {
	return w.contextMenu
}

func (w *Widget) SetContextMenu(value *Menu) {
	w.contextMenu = value
}

func (w *Widget) Enabled() (bool, os.Error) {
	ret := GetWindowLong(w.hWnd, GWL_STYLE)
	if ret == 0 {
		return false, lastError("GetWindowLong")
	}

	return (ret & WS_DISABLED) == 0, nil
}

func (w *Widget) SetEnabled(value bool) os.Error {
	style := GetWindowLong(w.hWnd, GWL_STYLE)
	if style == 0 {
		return lastError("GetWindowLong")
	}
	if value {
		style &^= WS_DISABLED
	} else {
		style |= WS_DISABLED
	}

	SetLastError(0)
	ret := SetWindowLong(w.hWnd, GWL_STYLE, style)
	if ret == 0 {
		return lastError("SetWindowLong")
	}

	SendMessage(w.hWnd, WM_ENABLE, uintptr(BoolToBOOL(value)), 0)

	return nil
}

func (w *Widget) Font() *drawing.Font {
	return w.font
}

func (w *Widget) SetFont(value *drawing.Font) {
	if value != w.font {
		SendMessage(w.hWnd, WM_SETFONT, uintptr(value.Handle()), 1)

		w.font = value
	}
}

func (w *Widget) Parent() IContainer {
	return w.parent
}

func (w *Widget) SetParent(value IContainer) (err os.Error) {
	if value == w.parent {
		return nil
	}

	style := uint(GetWindowLong(w.hWnd, GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	if value == nil {
		style &^= WS_CHILD
		style |= WS_POPUP

		if SetParent(w.hWnd, 0) == 0 {
			return lastError("SetParent")
		}
		SetLastError(0)
		if SetWindowLong(w.hWnd, GWL_STYLE, int(style)) == 0 {
			return lastError("SetWindowLong")
		}
	} else {
		style |= WS_CHILD
		style &^= WS_POPUP

		SetLastError(0)
		if SetWindowLong(w.hWnd, GWL_STYLE, int(style)) == 0 {
			return lastError("SetWindowLong")
		}
		if SetParent(w.hWnd, value.Handle()) == 0 {
			return lastError("SetParent")
		}
	}

	b, err := w.Bounds()
	if err != nil {
		return err
	}

	if !SetWindowPos(w.hWnd, HWND_BOTTOM, b.X, b.Y, b.Width, b.Height, SWP_FRAMECHANGED) {
		return lastError("SetWindowPos")
	}

	oldParent := w.parent

	w.parent = value

	if oldParent != nil {
		oldParent.Children().Remove(w)
	}

	if value != nil && !value.Children().ContainsHandle(w.hWnd) {
		value.Children().Add(w)
	}

	return nil
}

func (w *Widget) Text() string {
	textLength := SendMessage(w.hWnd, WM_GETTEXTLENGTH, 0, 0)
	buf := make([]uint16, textLength+1)
	SendMessage(w.hWnd, WM_GETTEXT, uintptr(textLength+1), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func (w *Widget) SetText(value string) os.Error {
	if TRUE != SendMessage(w.hWnd, WM_SETTEXT, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value)))) {
		return newError("WM_SETTEXT failed")
	}

	return nil
}

func (w *Widget) Visible() (bool, os.Error) {
	style := GetWindowLong(w.hWnd, GWL_STYLE)
	if style == 0 {
		return false, lastError("GetWindowLong")
	}

	return (style & WS_VISIBLE) != 0, nil
}

func (w *Widget) SetVisible(value bool) os.Error {
	style := GetWindowLong(w.hWnd, GWL_STYLE)
	if style == 0 {
		return lastError("GetWindowLong")
	}

	if value {
		style |= WS_VISIBLE
	} else {
		style &^= WS_VISIBLE
	}

	SetLastError(0)
	if SetWindowLong(w.hWnd, GWL_STYLE, style) == 0 {
		return lastError("SetWindowLong")
	}

	SendMessage(w.hWnd, WM_SHOWWINDOW, uintptr(BoolToBOOL(value)), 0)

	return nil
}

func (w *Widget) Bounds() (drawing.Rectangle, os.Error) {
	var r RECT

	if !GetWindowRect(w.hWnd, &r) {
		return drawing.Rectangle{}, lastError("GetWindowRect")
	}

	b := drawing.Rectangle{X: r.Left, Y: r.Top, Width: r.Right - r.Left, Height: r.Bottom - r.Top}

	if w.parent != nil {
		p := POINT{b.X, b.Y}
		if !ScreenToClient(w.hWnd, &p) {
			return drawing.Rectangle{}, newError("ScreenToClient failed")
		}
		b.X = p.X
		b.Y = p.Y
	}

	return b, nil
}

func (w *Widget) SetBounds(bounds drawing.Rectangle) os.Error {
	if !MoveWindow(w.hWnd, bounds.X, bounds.Y, bounds.Width, bounds.Height, true) {
		return lastError("MoveWindow")
	}

	return nil
}

func (w *Widget) MaxSize() (drawing.Size, os.Error) {
	return w.maxSize, nil
}

func (w *Widget) SetMaxSize(value drawing.Size) os.Error {
	w.maxSize = value

	return nil
}

func (w *Widget) MinSize() (drawing.Size, os.Error) {
	return w.minSize, nil
}

func (w *Widget) SetMinSize(value drawing.Size) os.Error {
	w.minSize = value

	return nil
}

func (w *Widget) dialogBaseUnits() drawing.Size {
	// FIXME: Error handling
	hFont := HFONT(SendMessage(w.hWnd, WM_GETFONT, 0, 0))
	hdc := GetDC(w.hWnd)
	hFontOld := SelectObject(hdc, HGDIOBJ(hFont))

	var tm TEXTMETRIC
	GetTextMetrics(hdc, &tm)

	var size SIZE
	GetTextExtentPoint32(
		hdc,
		syscall.StringToUTF16Ptr("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
		52,
		&size)

	SelectObject(hdc, HGDIOBJ(hFontOld))
	ReleaseDC(w.hWnd, hdc)

	return drawing.Size{(size.CX/26 + 1) / 2, int(tm.TmHeight)}
}

func (w *Widget) dialogBaseUnitsToPixels(dlus drawing.Size) (pixels drawing.Size) {
	// FIXME: Cache dialog base units on font change.
	base := w.dialogBaseUnits()

	return drawing.Size{MulDiv(dlus.Width, base.Width, 4), MulDiv(dlus.Height, base.Height, 8)}
}

func (w *Widget) LayoutFlags() LayoutFlags {
	// FIXME: Figure out how to do this, if at all.
	return 0
}

func (w *Widget) PreferredSize() drawing.Size {
	// FIXME: Figure out how to do this, if at all.
	return w.dialogBaseUnitsToPixels(drawing.Size{10, 10})
}

func (w *Widget) Size() (size drawing.Size, err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	size = bounds.Size()
	return
}

func (w *Widget) SetSize(size drawing.Size) (err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	return w.SetBounds(bounds.SetSize(size))
}

func (w *Widget) X() (x int, err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	x = bounds.X
	return
}

func (w *Widget) SetX(value int) (err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	bounds.X = value

	return w.SetBounds(bounds)
}

func (w *Widget) Y() (y int, err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	y = bounds.Y
	return
}

func (w *Widget) SetY(value int) (err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	bounds.Y = value

	return w.SetBounds(bounds)
}

func (w *Widget) Width() (width int, err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	width = bounds.Width
	return
}

func (w *Widget) SetWidth(value int) (err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	bounds.Width = value

	return w.SetBounds(bounds)
}

func (w *Widget) Height() (height int, err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	height = bounds.Height
	return
}

func (w *Widget) SetHeight(value int) (err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	bounds.Height = value

	return w.SetBounds(bounds)
}

func (w *Widget) ClientBounds() (drawing.Rectangle, os.Error) {
	var r RECT

	if !GetClientRect(w.hWnd, &r) {
		return drawing.Rectangle{}, lastError("GetClientRect")
	}

	return drawing.Rectangle{X: r.Left, Y: r.Top, Width: r.Right - r.Left, Height: r.Bottom - r.Top}, nil
}

func (w *Widget) SetFocus() os.Error {
	if SetFocus(w.hWnd) == 0 {
		return lastError("SetFocus")
	}

	return nil
}

func (w *Widget) GroupStart() (bool, os.Error) {
	style := GetWindowLong(w.hWnd, GWL_STYLE)
	if style == 0 {
		return false, lastError("GetWindowLong")
	}

	return (style & WS_GROUP) != 0, nil
}

func (w *Widget) SetGroupStart(value bool) os.Error {
	style := GetWindowLong(w.hWnd, GWL_STYLE)
	if style == 0 {
		return lastError("GetWindowLong")
	}

	if value {
		style |= WS_GROUP
	} else {
		style &^= WS_GROUP
	}

	SetLastError(0)
	if SetWindowLong(w.hWnd, GWL_STYLE, style) == 0 {
		return lastError("SetWindowLong")
	}

	return nil
}

func (w *Widget) GetDrawingSurface() (*drawing.Surface, os.Error) {
	return drawing.NewSurfaceFromWidget(w.hWnd)
}

func (w *Widget) AddKeyDownHandler(handler KeyEventHandler) {
	w.keyDownHandlers.Push(handler)
}

func (w *Widget) RemoveKeyDownHandler(handler KeyEventHandler) {
	for i, h := range w.keyDownHandlers {
		if h.(KeyEventHandler) == handler {
			w.keyDownHandlers.Delete(i)
			break
		}
	}
}

func (w *Widget) raiseKeyDown(args KeyEventArgs) {
	for _, handlerIface := range w.keyDownHandlers {
		handler := handlerIface.(KeyEventHandler)
		handler(args)
	}
}

func (w *Widget) AddMouseDownHandler(handler MouseEventHandler) {
	w.mouseDownHandlers.Push(handler)
}

func (w *Widget) AddSizeChangedHandler(handler EventHandler) {
	w.sizeChangedHandlers.Push(handler)
}

func (w *Widget) RemoveSizeChangedHandler(handler EventHandler) {
	for i, h := range w.sizeChangedHandlers {
		if h.(EventHandler) == handler {
			w.sizeChangedHandlers.Delete(i)
			break
		}
	}
}

func (w *Widget) raiseSizeChanged() {
	for _, handlerIface := range w.sizeChangedHandlers {
		handler := handlerIface.(EventHandler)
		handler(&eventArgs{widgetsByHWnd[w.hWnd]})
	}
}

func (w *Widget) raiseEvent(msg *MSG) os.Error {
	switch msg.Message {
	case WM_LBUTTONDOWN:
		for _, handlerIface := range w.mouseDownHandlers {
			handler := handlerIface.(MouseEventHandler)
			handler(&mouseEventArgs{eventArgs: eventArgs{sender: widgetsByHWnd[w.hWnd]}})
		}

	case crutches.ContextMenuMsgId():
		hWnd := HWND(msg.WParam)
		sourceWidget := widgetsByHWnd[hWnd]
		x := int(GET_X_LPARAM(msg.LParam))
		y := int(GET_Y_LPARAM(msg.LParam))

		contextMenu := sourceWidget.ContextMenu()

		if contextMenu != nil {
			TrackPopupMenuEx(contextMenu.hMenu, TPM_NOANIMATION, x, y, rootWidget(sourceWidget).Handle(), nil)
		}
		//		MessageBox(rootWidget(w), "WM_CONTEXTMENU", "The context menu would go here.", MB_OK|MB_ICONINFORMATION)

	case WM_KEYDOWN:
		w.raiseKeyDown(&keyEventArgs{eventArgs: eventArgs{widgetsByHWnd[w.hWnd]}, key: int(msg.WParam)})

	case crutches.ResizeMsgId():
		w.raiseSizeChanged()

	case /*WM_COMMAND,*/ crutches.CommandMsgId():
		switch HIWORD(uint(msg.WParam)) {
		case 0:
			// menu
			actionId := uint16(LOWORD(uint(msg.WParam)))
			if action, ok := actionsById[actionId]; ok {
				action.raiseTriggered()
			}
		}
	}

	return nil
}

func (w *Widget) runMessageLoop() os.Error {
	var msg MSG
	var cm crutches.Message

	// Copy handle for later use when closing the window.
	hWnd := w.hWnd

	for {
		ret := GetMessage(&msg, 0, 0, 0)

		switch msg.Message {
		case WM_MOUSEMOVE, WM_PAINT, WM_TIMER:

			//        case WM_COMMAND, commandMsgId:
			//			fmt.Printf("Widget.runMessageLoop: WM_COMMAND, commandMsgId: msg: %v\n", msg)

			//        default:
			//            fmt.Printf("Widget.runMessageLoop: msg: %v\n", msg)
		}

		switch ret {
		case 0:
			fmt.Println("Widget.runMessageLoop: GetMessage returned 0, exiting")
			return nil

		case -1:
			return newError("GetMessage returned -1")
		}

		rootHWnd := GetAncestor(msg.HWnd, GA_ROOT)
		if rootHWnd == 0 {
			rootHWnd = msg.HWnd
		}

		if !IsDialogMessage(rootHWnd, &msg) {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}

		if widget, ok := widgetsByHWnd[msg.HWnd]; ok {
			widget.raiseEvent(&msg)
		}

		for {
			if crutches.GetCustomMessage(&cm) < 1 {
				break
			}

			//			fmt.Printf("Widget.runMessageLoop: cm: %v\n", cm)

			var cmMsg MSG
			cmMsg.HWnd = cm.Hwnd
			cmMsg.Message = cm.Msg
			cmMsg.WParam = cm.WParam
			cmMsg.LParam = cm.LParam

			if widget, ok := widgetsByHWnd[cmMsg.HWnd]; ok {
				widget.raiseEvent(&cmMsg)
			}
		}

		switch msg.Message {
		case crutches.CloseMsgId():
			// We compare to local var, because Dispose has zeroed the original.
			if msg.HWnd == hWnd {
				return nil
			}

		case WM_LBUTTONDOWN:
			trackedMouseDownHWnd = msg.HWnd

		case WM_LBUTTONUP:
			trackedMouseDownHWnd = 0
		}
	}

	return nil
}
