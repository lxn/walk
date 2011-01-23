// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/gdi32"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
	. "walk/winapi/uxtheme"
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
	BeginUpdate()
	EndUpdate()
	Invalidate() os.Error
	Bounds() (drawing.Rectangle, os.Error)
	SetBounds(value drawing.Rectangle) os.Error
	ClientBounds() (drawing.Rectangle, os.Error)
	ContextMenu() *Menu
	Cursor() Cursor
	SetCursor(value Cursor)
	SetContextMenu(value *Menu)
	Dispose()
	IsDisposed() bool
	Enabled() bool
	SetEnabled(value bool)
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
	RootWidget() RootWidget
	GetDrawingSurface() (*drawing.Surface, os.Error)
}

type widgetInternal interface {
	IWidget
	wndProc(msg *MSG, origWndProcPtr uintptr) uintptr
}

type Widget struct {
	hWnd                 HWND
	parent               IContainer
	font                 *drawing.Font
	contextMenu          *Menu
	keyDownPublisher     KeyEventPublisher
	mouseDownPublisher   MouseEventPublisher
	mouseUpPublisher     MouseEventPublisher
	mouseMovePublisher   MouseEventPublisher
	sizeChangedPublisher EventPublisher
	maxSize              drawing.Size
	minSize              drawing.Size
	cursor               Cursor
}

var (
	widgetsByHWnd map[HWND]widgetInternal = make(map[HWND]widgetInternal)
)

func ensureRegisteredWindowClass(className string, windowProc syscall.CallbackFunc, callback **syscall.Callback) {
	if callback == nil {
		panic("callback cannot be nil")
	}

	if *callback != nil {
		return
	}

	hInst := GetModuleHandle(nil)
	if hInst == 0 {
		panic("GetModuleHandle failed")
	}

	hIcon := LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(IDI_APPLICATION))))
	if hIcon == 0 {
		panic("LoadIcon failed")
	}

	hCursor := LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor failed")
	}

	*callback = syscall.NewCallback(windowProc, 4)

	var wc WNDCLASSEX
	wc.CbSize = uint(unsafe.Sizeof(wc))
	wc.LpfnWndProc = (*callback).ExtFnEntry()
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

func msgFromCallbackArgs(args *uintptr) *MSG {
	p := (*[4]int32)(unsafe.Pointer(args))

	return &MSG{
		HWnd:    HWND(p[0]),
		Message: uint(p[1]),
		WParam:  uintptr(p[2]),
		LParam:  uintptr(p[3]),
	}
}

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

func (w *Widget) Cursor() Cursor {
	return w.cursor
}

func (w *Widget) SetCursor(value Cursor) {
	w.cursor = value
}

func (w *Widget) Enabled() bool {
	return IsWindowEnabled(w.hWnd)
}

func (w *Widget) SetEnabled(value bool) {
	EnableWindow(w.hWnd, value)
}

func (w *Widget) Font() *drawing.Font {
	return w.font
}

func (w *Widget) SetFont(value *drawing.Font) {
	if value != w.font {
		SendMessage(w.hWnd, WM_SETFONT, uintptr(value.HandleForDPI(0)), 1)

		w.font = value
	}
}

func (w *Widget) BeginUpdate() {
	SendMessage(w.hWnd, WM_SETREDRAW, 0, 0)
}

func (w *Widget) EndUpdate() {
	SendMessage(w.hWnd, WM_SETREDRAW, 1, 0)
}

func (w *Widget) Invalidate() os.Error {
	cb, err := w.ClientBounds()
	if err != nil {
		return err
	}

	r := &RECT{cb.X, cb.Y, cb.X + cb.Width, cb.Y + cb.Height}

	if !InvalidateRect(w.hWnd, r, true) {
		return newError("InvalidateRect failed")
	}

	return nil
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
		if style&WS_VISIBLE > 0 {
			return nil
		}

		style |= WS_VISIBLE
	} else {
		if style&WS_VISIBLE == 0 {
			return nil
		}

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
		if !ScreenToClient(w.parent.Handle(), &p) {
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
	return drawing.NewSurfaceFromHWND(w.hWnd)
}

func (w *Widget) setTheme(appName string) os.Error {
	if hr := SetWindowTheme(w.hWnd, syscall.StringToUTF16Ptr(appName), nil); FAILED(hr) {
		return errorFromHRESULT("SetWindowTheme", hr)
	}

	return nil
}

func (w *Widget) KeyDown() *KeyEvent {
	return w.keyDownPublisher.Event()
}

func (w *Widget) MouseDown() *MouseEvent {
	return w.mouseDownPublisher.Event()
}

func (w *Widget) MouseMove() *MouseEvent {
	return w.mouseMovePublisher.Event()
}

func (w *Widget) MouseUp() *MouseEvent {
	return w.mouseUpPublisher.Event()
}

func (w *Widget) mouseEventArgsFromMSG(msg *MSG) *MouseEventArgs {
	x := int(GET_X_LPARAM(msg.LParam))
	y := int(GET_Y_LPARAM(msg.LParam))

	return NewMouseEventArgs(widgetsByHWnd[w.hWnd], x, y, 0)
}

func (w *Widget) wndProc(msg *MSG, origWndProcPtr uintptr) uintptr {
	//	widget := widgetsByHWnd[w.hWnd]
	//	fmt.Printf("*Widget.wndProc: type: %T, msg: %+v\n", widget, msg)

	switch msg.Message {
	case WM_LBUTTONDOWN:
		SetCapture(w.hWnd)
		w.mouseDownPublisher.Publish(w.mouseEventArgsFromMSG(msg))

	case WM_LBUTTONUP:
		if !ReleaseCapture() {
			log.Println(lastError("ReleaseCapture"))
		}
		w.mouseUpPublisher.Publish(w.mouseEventArgsFromMSG(msg))

	case WM_MOUSEMOVE:
		w.mouseMovePublisher.Publish(w.mouseEventArgsFromMSG(msg))

	case WM_SETCURSOR:
		if w.cursor != nil {
			SetCursor(w.cursor.handle())
			return 0
		}

	case WM_CONTEXTMENU:
		sourceWidget := widgetsByHWnd[HWND(msg.WParam)]
		if sourceWidget == nil {
			break
		}

		x := int(GET_X_LPARAM(msg.LParam))
		y := int(GET_Y_LPARAM(msg.LParam))

		contextMenu := sourceWidget.ContextMenu()

		if contextMenu != nil {
			TrackPopupMenuEx(contextMenu.hMenu, TPM_NOANIMATION, x, y, rootWidget(sourceWidget).Handle(), nil)
		}
		return 0

	case WM_KEYDOWN:
		w.keyDownPublisher.Publish(NewKeyEventArgs(widgetsByHWnd[w.hWnd], int(msg.WParam)))

	case WM_SIZE, WM_SIZING:
		w.sizeChangedPublisher.Publish(NewEventArgs(widgetsByHWnd[w.hWnd]))

	case WM_GETMINMAXINFO:
		mmi := (*MINMAXINFO)(unsafe.Pointer(msg.LParam))
		mmi.PtMinTrackSize = POINT{w.minSize.Width, w.minSize.Height}
		return 0
	}

	if origWndProcPtr != 0 {
		return CallWindowProc(origWndProcPtr, msg.HWnd, msg.Message, msg.WParam, msg.LParam)
	}

	return DefWindowProc(msg.HWnd, msg.Message, msg.WParam, msg.LParam)
}

func (w *Widget) runMessageLoop() (int, os.Error) {
	var msg MSG

	for w.hWnd != 0 {
		switch GetMessage(&msg, 0, 0, 0) {
		case 0:
			return int(msg.WParam), nil

		case -1:
			return -1, lastError("GetMessage")
		}

		if !IsDialogMessage(w.hWnd, &msg) {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}

	return 0, nil
}
