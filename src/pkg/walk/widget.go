// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"log"
	"os"
	"syscall"
	"unsafe"
)

import (
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
	Name() string
	SetName(name string)
	BeginUpdate()
	EndUpdate()
	Invalidate() os.Error
	Bounds() (Rectangle, os.Error)
	SetBounds(value Rectangle) os.Error
	ClientBounds() (Rectangle, os.Error)
	ContextMenu() *Menu
	Cursor() Cursor
	SetCursor(value Cursor)
	SetContextMenu(value *Menu)
	Dispose()
	IsDisposed() bool
	Enabled() bool
	SetEnabled(value bool)
	Font() *Font
	SetFont(value *Font)
	GroupStart() (bool, os.Error)
	SetGroupStart(value bool) os.Error
	Height() (int, os.Error)
	SetHeight(value int) os.Error
	LayoutFlags() LayoutFlags
	MaxSize() (Size, os.Error)
	SetMaxSize(value Size) os.Error
	MinSize() (Size, os.Error)
	SetMinSize(value Size) os.Error
	Parent() IContainer
	SetParent(value IContainer) os.Error
	PreferredSize() Size
	Size() (Size, os.Error)
	SetSize(value Size) os.Error
	Text() string
	SetText(value string) os.Error
	Visible() bool
	SetVisible(value bool)
	Width() (int, os.Error)
	SetWidth(value int) os.Error
	X() (int, os.Error)
	SetX(value int) os.Error
	Y() (int, os.Error)
	SetY(value int) os.Error
	SetFocus() os.Error
	RootWidget() RootWidget
	GetDrawingSurface() (*Surface, os.Error)
}

type widgetInternal interface {
	IWidget
	wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr
	path() string
	writePath(buf *bytes.Buffer)
}

type Widget struct {
	hWnd                 HWND
	name                 string
	parent               IContainer
	font                 *Font
	contextMenu          *Menu
	keyDownPublisher     KeyEventPublisher
	mouseDownPublisher   MouseEventPublisher
	mouseUpPublisher     MouseEventPublisher
	mouseMovePublisher   MouseEventPublisher
	sizeChangedPublisher EventPublisher
	maxSize              Size
	minSize              Size
	cursor               Cursor
}

var (
	widgetsByHWnd map[HWND]widgetInternal = make(map[HWND]widgetInternal)
)

func ensureRegisteredWindowClass(className string, windowProc interface{}, callback *uintptr) {
	if callback == nil {
		panic("callback cannot be nil")
	}

	if *callback != 0 {
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

	*callback = syscall.NewCallback(windowProc)

	var wc WNDCLASSEX
	wc.CbSize = uint(unsafe.Sizeof(wc))
	wc.LpfnWndProc = *callback
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

/*func msgFromCallbackArgs(args *uintptr) *MSG {
	p := (*[4]int32)(unsafe.Pointer(args))

	return &MSG{
		HWnd:    HWND(p[0]),
		Message: uint(p[1]),
		WParam:  uintptr(p[2]),
		LParam:  uintptr(p[3]),
	}
}*/

func rootWidget(w IWidget) RootWidget {
	if w == nil {
		return nil
	}

	hWndRoot := GetAncestor(w.Handle(), GA_ROOT)

	rw, _ := widgetsByHWnd[hWndRoot].(RootWidget)
	return rw
}

func (w *Widget) setAndClearStyleBits(set, clear uint) os.Error {
	style := uint(GetWindowLong(w.hWnd, GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	var newStyle uint
	newStyle = (style | set) &^ clear

	if newStyle != style {
		SetLastError(0)
		if SetWindowLong(w.hWnd, GWL_STYLE, int(newStyle)) == 0 {
			return lastError("SetWindowLong")
		}
	}

	return nil
}

func (w *Widget) ensureStyleBits(bits uint, set bool) os.Error {
	var setBits uint
	var clearBits uint
	if set {
		setBits = bits
	} else {
		clearBits = bits
	}
	return w.setAndClearStyleBits(setBits, clearBits)
}

func (w *Widget) Name() string {
	return w.name
}

func (w *Widget) SetName(name string) {
	w.name = name
}

func (w *Widget) writePath(buf *bytes.Buffer) {
	hWndParent := GetAncestor(w.hWnd, GA_PARENT)
	if pw, ok := widgetsByHWnd[hWndParent]; ok {
		if pwi, ok := pw.(widgetInternal); ok {
			pwi.writePath(buf)
			buf.WriteByte('/')
		}
	}

	buf.WriteString(w.name)
}

func (w *Widget) path() string {
	buf := bytes.NewBuffer(nil)

	w.writePath(buf)

	return buf.String()
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

func (w *Widget) Font() *Font {
	return w.font
}

func (w *Widget) SetFont(value *Font) {
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

func (w *Widget) Visible() bool {
	return IsWindowVisible(w.hWnd)
}

func (w *Widget) SetVisible(visible bool) {
	var cmd int
	if visible {
		cmd = SW_SHOW
	} else {
		cmd = SW_HIDE
	}
	ShowWindow(w.hWnd, cmd)
}

func (w *Widget) Bounds() (Rectangle, os.Error) {
	var r RECT

	if !GetWindowRect(w.hWnd, &r) {
		return Rectangle{}, lastError("GetWindowRect")
	}

	b := Rectangle{X: r.Left, Y: r.Top, Width: r.Right - r.Left, Height: r.Bottom - r.Top}

	if w.parent != nil {
		p := POINT{b.X, b.Y}
		if !ScreenToClient(w.parent.Handle(), &p) {
			return Rectangle{}, newError("ScreenToClient failed")
		}
		b.X = p.X
		b.Y = p.Y
	}

	return b, nil
}

func (w *Widget) SetBounds(bounds Rectangle) os.Error {
	if !MoveWindow(w.hWnd, bounds.X, bounds.Y, bounds.Width, bounds.Height, true) {
		return lastError("MoveWindow")
	}

	return nil
}

func (w *Widget) MaxSize() (Size, os.Error) {
	return w.maxSize, nil
}

func (w *Widget) SetMaxSize(value Size) os.Error {
	w.maxSize = value

	return nil
}

func (w *Widget) MinSize() (Size, os.Error) {
	return w.minSize, nil
}

func (w *Widget) SetMinSize(value Size) os.Error {
	w.minSize = value

	return nil
}

func (w *Widget) dialogBaseUnits() Size {
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

	return Size{(size.CX/26 + 1) / 2, int(tm.TmHeight)}
}

func (w *Widget) dialogBaseUnitsToPixels(dlus Size) (pixels Size) {
	// FIXME: Cache dialog base units on font change.
	base := w.dialogBaseUnits()

	return Size{MulDiv(dlus.Width, base.Width, 4), MulDiv(dlus.Height, base.Height, 8)}
}

func (w *Widget) LayoutFlags() LayoutFlags {
	// FIXME: Figure out how to do this, if at all.
	return 0
}

func (w *Widget) PreferredSize() Size {
	// FIXME: Figure out how to do this, if at all.
	return w.dialogBaseUnitsToPixels(Size{10, 10})
}

func (w *Widget) Size() (size Size, err os.Error) {
	bounds, err := w.Bounds()
	if err != nil {
		return
	}

	size = bounds.Size()
	return
}

func (w *Widget) SetSize(size Size) (err os.Error) {
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

func (w *Widget) ClientBounds() (Rectangle, os.Error) {
	var r RECT

	if !GetClientRect(w.hWnd, &r) {
		return Rectangle{}, lastError("GetClientRect")
	}

	return Rectangle{X: r.Left, Y: r.Top, Width: r.Right - r.Left, Height: r.Bottom - r.Top}, nil
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

func (w *Widget) GetDrawingSurface() (*Surface, os.Error) {
	return NewSurfaceFromHWND(w.hWnd)
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

func (w *Widget) publishMouseEvent(publisher *MouseEventPublisher, wParam, lParam uintptr) {
	x := int(GET_X_LPARAM(lParam))
	y := int(GET_Y_LPARAM(lParam))

	publisher.Publish(x, y, 0)
}

func (w *Widget) SizeChanged() *Event {
	return w.sizeChangedPublisher.Event()
}

func (w *Widget) persistState(restore bool) {
	settings := appSingleton.settings
	if settings != nil {
		if widget, ok := widgetsByHWnd[w.hWnd]; ok {
			if persistable, ok := widget.(Persistable); ok && persistable.Persistent() {
				if restore {
					if err := persistable.RestoreState(); err != nil {
						log.Println(err)
					}
				} else {
					if err := persistable.SaveState(); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func (w *Widget) getState() (string, os.Error) {
	settings := appSingleton.settings
	if settings == nil {
		return "", newError("App().Settings() must not be nil")
	}

	state, _ := settings.Get(w.path())
	return state, nil
}

func (w *Widget) putState(state string) os.Error {
	settings := appSingleton.settings
	if settings == nil {
		return newError("App().Settings() must not be nil")
	}

	return settings.Put(w.path(), state)
}

func (w *Widget) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr, origWndProcPtr uintptr) uintptr {
	switch msg {
	case WM_LBUTTONDOWN:
		SetCapture(w.hWnd)
		w.publishMouseEvent(&w.mouseDownPublisher, wParam, lParam)

	case WM_LBUTTONUP:
		if !ReleaseCapture() {
			log.Println(lastError("ReleaseCapture"))
		}
		w.publishMouseEvent(&w.mouseUpPublisher, wParam, lParam)

	case WM_MOUSEMOVE:
		w.publishMouseEvent(&w.mouseMovePublisher, wParam, lParam)

	case WM_SETCURSOR:
		if w.cursor != nil {
			SetCursor(w.cursor.handle())
			return 0
		}

	case WM_CONTEXTMENU:
		sourceWidget := widgetsByHWnd[HWND(wParam)]
		if sourceWidget == nil {
			break
		}

		x := int(GET_X_LPARAM(lParam))
		y := int(GET_Y_LPARAM(lParam))

		contextMenu := sourceWidget.ContextMenu()

		if contextMenu != nil {
			TrackPopupMenuEx(contextMenu.hMenu, TPM_NOANIMATION, x, y, rootWidget(sourceWidget).Handle(), nil)
		}
		return 0

	case WM_KEYDOWN:
		w.keyDownPublisher.Publish(int(wParam))

	case WM_SIZE, WM_SIZING:
		w.sizeChangedPublisher.Publish()

	case WM_GETMINMAXINFO:
		mmi := (*MINMAXINFO)(unsafe.Pointer(lParam))
		mmi.PtMinTrackSize = POINT{w.minSize.Width, w.minSize.Height}
		return 0

	case WM_SHOWWINDOW:
		w.persistState(wParam != 0)

	case WM_DESTROY:
		w.persistState(false)
	}

	if origWndProcPtr != 0 {
		return CallWindowProc(origWndProcPtr, hwnd, msg, wParam, lParam)
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

func (w *Widget) runMessageLoop() int {
	var msg MSG

	for w.hWnd != 0 {
		switch GetMessage(&msg, 0, 0, 0) {
		case 0:
			return int(msg.WParam)

		case -1:
			return -1
		}

		if !IsDialogMessage(w.hWnd, &msg) {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}

	return 0
}
