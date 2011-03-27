// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"fmt"
	"os"
	"strings"
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
	ShrinkableHorz LayoutFlags = 1 << iota
	ShrinkableVert
	GrowableHorz
	GrowableVert
	GreedyHorz
	GreedyVert
)

type Widget interface {
	Background() Brush
	BaseWidget() *WidgetBase
	Bounds() Rectangle
	BringToTop() os.Error
	ClientBounds() Rectangle
	ContextMenu() *Menu
	CreateCanvas() (*Canvas, os.Error)
	Cursor() Cursor
	Dispose()
	Enabled() bool
	Font() *Font
	Height() int
	Invalidate() os.Error
	IsDisposed() bool
	KeyDown() *KeyEvent
	LayoutFlags() LayoutFlags
	MaxSize() Size
	MinSize() Size
	MinSizeHint() Size
	MouseDown() *MouseEvent
	MouseMove() *MouseEvent
	MouseUp() *MouseEvent
	Name() string
	Parent() Container
	RootWidget() RootWidget
	SetBackground(value Brush)
	SetBounds(value Rectangle) os.Error
	SetClientSize(value Size) os.Error
	SetContextMenu(value *Menu)
	SetCursor(value Cursor)
	SetEnabled(value bool)
	SetFocus() os.Error
	SetFont(value *Font)
	SetHeight(value int) os.Error
	SetMinMaxSize(min, max Size) os.Error
	SetName(name string)
	SetParent(value Container) os.Error
	SetSize(value Size) os.Error
	SetSuspended(suspend bool)
	SetVisible(value bool)
	SetWidth(value int) os.Error
	SetX(value int) os.Error
	SetY(value int) os.Error
	Size() Size
	SizeChanged() *Event
	SizeHint() Size
	Suspended() bool
	Visible() bool
	Width() int
	X() int
	Y() int
}

type widgetInternal interface {
	Widget
	path() string
	wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr
	writePath(buf *bytes.Buffer)
}

type subclassedWidget interface {
	widgetInternal
	origWndProcPtr() uintptr
	setOrigWndProcPtr(ptr uintptr)
}

type WidgetBase struct {
	widget               widgetInternal
	hWnd                 HWND
	name                 string
	parent               Container
	font                 *Font
	contextMenu          *Menu
	keyDownPublisher     KeyEventPublisher
	mouseDownPublisher   MouseEventPublisher
	mouseUpPublisher     MouseEventPublisher
	mouseMovePublisher   MouseEventPublisher
	sizeChangedPublisher EventPublisher
	maxSize              Size
	minSize              Size
	background           Brush
	cursor               Cursor
	layoutFlags          LayoutFlags
	suspended            bool
}

var widgetWndProcPtr uintptr = syscall.NewCallback(widgetWndProc)

func ensureRegisteredWindowClass(className string, registered *bool) {
	if registered == nil {
		panic("registered cannot be nil")
	}

	if *registered {
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

	var wc WNDCLASSEX
	wc.CbSize = uint(unsafe.Sizeof(wc))
	wc.LpfnWndProc = widgetWndProcPtr
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}

	*registered = true
}

func initWidget(widget widgetInternal, parent Widget, className string, style, exStyle uint) os.Error {
	wb := widget.BaseWidget()
	wb.widget = widget

	var hwndParent HWND
	if parent != nil {
		hwndParent = parent.BaseWidget().hWnd

		if container, ok := parent.(Container); ok {
			wb.parent = container
		}
	}

	wb.hWnd = CreateWindowEx(
		exStyle,
		syscall.StringToUTF16Ptr(className),
		nil,
		style|WS_CLIPSIBLINGS,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		hwndParent,
		0,
		0,
		nil)
	if wb.hWnd == 0 {
		return lastError("CreateWindowEx")
	}

	succeeded := false
	defer func() {
		if !succeeded {
			wb.Dispose()
		}
	}()

	SetWindowLongPtr(wb.hWnd, GWLP_USERDATA, uintptr(unsafe.Pointer(wb)))

	if subclassed, ok := widget.(subclassedWidget); ok {
		origWndProcPtr := SetWindowLongPtr(wb.hWnd, GWLP_WNDPROC, widgetWndProcPtr)
		if origWndProcPtr == 0 {
			return lastError("SetWindowLongPtr")
		}

		if subclassed.origWndProcPtr() == 0 {
			subclassed.setOrigWndProcPtr(origWndProcPtr)
		}
	}

	wb.SetFont(defaultFont)

	succeeded = true

	return nil
}

func initChildWidget(widget widgetInternal, parent Widget, className string, style, exStyle uint) os.Error {
	if parent == nil {
		return newError("parent cannot be nil")
	}

	if err := initWidget(widget, parent, className, style|WS_CHILD, exStyle); err != nil {
		return err
	}

	if container, ok := parent.(Container); ok {
		if container.Children() == nil {
			// Required by parents like MainWindow and GroupBox.
			if SetParent(widget.BaseWidget().hWnd, container.BaseWidget().hWnd) == 0 {
				return lastError("SetParent")
			}
		} else {
			if err := container.Children().Add(widget); err != nil {
				return err
			}
		}
	}

	return nil
}

func rootWidget(w Widget) RootWidget {
	if w == nil {
		return nil
	}

	hWndRoot := GetAncestor(w.BaseWidget().hWnd, GA_ROOT)

	rw, _ := widgetFromHWND(hWndRoot).(RootWidget)
	return rw
}

func (wb *WidgetBase) setAndClearStyleBits(set, clear uint) os.Error {
	style := uint(GetWindowLong(wb.hWnd, GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	var newStyle uint
	newStyle = (style | set) &^ clear

	if newStyle != style {
		SetLastError(0)
		if SetWindowLong(wb.hWnd, GWL_STYLE, int(newStyle)) == 0 {
			return lastError("SetWindowLong")
		}
	}

	return nil
}

func (wb *WidgetBase) ensureStyleBits(bits uint, set bool) os.Error {
	var setBits uint
	var clearBits uint
	if set {
		setBits = bits
	} else {
		clearBits = bits
	}
	return wb.setAndClearStyleBits(setBits, clearBits)
}

func (wb *WidgetBase) Name() string {
	return wb.name
}

func (wb *WidgetBase) SetName(name string) {
	wb.name = name
}

func (wb *WidgetBase) writePath(buf *bytes.Buffer) {
	hWndParent := GetAncestor(wb.hWnd, GA_PARENT)
	if pwi := widgetFromHWND(hWndParent); pwi != nil {
		pwi.writePath(buf)
		buf.WriteByte('/')
	}

	buf.WriteString(wb.name)
}

func (wb *WidgetBase) path() string {
	buf := bytes.NewBuffer(nil)

	wb.writePath(buf)

	return buf.String()
}

func (wb *WidgetBase) BaseWidget() *WidgetBase {
	return wb
}

func (wb *WidgetBase) Dispose() {
	if wb.hWnd != 0 {
		DestroyWindow(wb.hWnd)
		wb.hWnd = 0
	}
}

func (wb *WidgetBase) IsDisposed() bool {
	return wb.hWnd == 0
}

func (wb *WidgetBase) RootWidget() RootWidget {
	return rootWidget(wb)
}

func (wb *WidgetBase) ContextMenu() *Menu {
	return wb.contextMenu
}

func (wb *WidgetBase) SetContextMenu(value *Menu) {
	wb.contextMenu = value
}

func (wb *WidgetBase) Background() Brush {
	return wb.background
}

func (wb *WidgetBase) SetBackground(value Brush) {
	wb.background = value
}

func (wb *WidgetBase) Cursor() Cursor {
	return wb.cursor
}

func (wb *WidgetBase) SetCursor(value Cursor) {
	wb.cursor = value
}

func (wb *WidgetBase) Enabled() bool {
	return IsWindowEnabled(wb.hWnd)
}

func (wb *WidgetBase) SetEnabled(value bool) {
	EnableWindow(wb.hWnd, value)
}

func (wb *WidgetBase) Font() *Font {
	return wb.font
}

func setWidgetFont(hwnd HWND, font *Font) {
	SendMessage(hwnd, WM_SETFONT, uintptr(font.handleForDPI(0)), 1)
}

func (wb *WidgetBase) SetFont(value *Font) {
	if value != wb.font {
		setWidgetFont(wb.hWnd, value)

		wb.font = value
	}
}

func (wb *WidgetBase) Suspended() bool {
	return wb.suspended
}

func (wb *WidgetBase) SetSuspended(suspend bool) {
	if suspend == wb.suspended {
		return
	}

	var wParam int
	if suspend {
		wParam = 0
	} else {
		wParam = 1
	}

	SendMessage(wb.hWnd, WM_SETREDRAW, uintptr(wParam), 0)

	wb.suspended = suspend
}

func (wb *WidgetBase) Invalidate() os.Error {
	if !InvalidateRect(wb.hWnd, nil, true) {
		return newError("InvalidateRect failed")
	}

	return nil
}

func (wb *WidgetBase) Parent() Container {
	return wb.parent
}

func (wb *WidgetBase) SetParent(value Container) (err os.Error) {
	if value == wb.parent {
		return nil
	}

	style := uint(GetWindowLong(wb.hWnd, GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	if value == nil {
		style &^= WS_CHILD
		style |= WS_POPUP

		if SetParent(wb.hWnd, 0) == 0 {
			return lastError("SetParent")
		}
		SetLastError(0)
		if SetWindowLong(wb.hWnd, GWL_STYLE, int(style)) == 0 {
			return lastError("SetWindowLong")
		}
	} else {
		style |= WS_CHILD
		style &^= WS_POPUP

		SetLastError(0)
		if SetWindowLong(wb.hWnd, GWL_STYLE, int(style)) == 0 {
			return lastError("SetWindowLong")
		}
		if SetParent(wb.hWnd, value.BaseWidget().hWnd) == 0 {
			return lastError("SetParent")
		}
	}

	b := wb.Bounds()

	if !SetWindowPos(wb.hWnd, HWND_BOTTOM, b.X, b.Y, b.Width, b.Height, SWP_FRAMECHANGED) {
		return lastError("SetWindowPos")
	}

	oldParent := wb.parent

	wb.parent = value

	if oldParent != nil {
		oldParent.Children().Remove(wb)
	}

	if value != nil && !value.Children().containsHandle(wb.hWnd) {
		value.Children().Add(wb)
	}

	return nil
}

func widgetText(hwnd HWND) string {
	textLength := SendMessage(hwnd, WM_GETTEXTLENGTH, 0, 0)
	buf := make([]uint16, textLength+1)
	SendMessage(hwnd, WM_GETTEXT, uintptr(textLength+1), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func setWidgetText(hwnd HWND, text string) os.Error {
	if TRUE != SendMessage(hwnd, WM_SETTEXT, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text)))) {
		return newError("WM_SETTEXT failed")
	}

	return nil
}

func (wb *WidgetBase) Visible() bool {
	return IsWindowVisible(wb.hWnd)
}

func (wb *WidgetBase) SetVisible(visible bool) {
	var cmd int
	if visible {
		cmd = SW_SHOW
	} else {
		cmd = SW_HIDE
	}
	ShowWindow(wb.hWnd, cmd)
}

func (wb *WidgetBase) BringToTop() os.Error {
	if wb.parent != nil {
		if err := wb.parent.BringToTop(); err != nil {
			return err
		}
	}

	if !SetWindowPos(wb.hWnd, HWND_TOP, 0, 0, 0, 0, SWP_NOACTIVATE|SWP_NOMOVE|SWP_NOSIZE) {
		return lastError("SetWindowPos")
	}

	return nil
}

func (wb *WidgetBase) Bounds() Rectangle {
	var r RECT

	if !GetWindowRect(wb.hWnd, &r) {
		lastError("GetWindowRect")
		return Rectangle{}
	}

	b := Rectangle{X: r.Left, Y: r.Top, Width: r.Right - r.Left, Height: r.Bottom - r.Top}

	if wb.parent != nil {
		p := POINT{b.X, b.Y}
		if !ScreenToClient(wb.parent.BaseWidget().hWnd, &p) {
			newError("ScreenToClient failed")
			return Rectangle{}
		}
		b.X = p.X
		b.Y = p.Y
	}

	return b
}

func (wb *WidgetBase) SetBounds(bounds Rectangle) os.Error {
	if !MoveWindow(wb.hWnd, bounds.X, bounds.Y, bounds.Width, bounds.Height, true) {
		return lastError("MoveWindow")
	}

	return nil
}

func (wb *WidgetBase) MinSize() Size {
	return wb.minSize
}

func (wb *WidgetBase) MaxSize() Size {
	return wb.maxSize
}

func (wb *WidgetBase) SetMinMaxSize(min, max Size) os.Error {
	if min.Width < 0 || min.Height < 0 {
		return newError("min must be positive")
	}
	if max.Width > 0 && max.Width < min.Width ||
		max.Height > 0 && max.Height < min.Height {
		return newError("max must be greater as or equal to min")
	}

	wb.minSize = min
	wb.maxSize = max

	return nil
}

func (wb *WidgetBase) dialogBaseUnits() Size {
	// The widget may use a font different from that in WidgetBase,
	// like e.g. NumberEdit does, so we try to use the right one.
	widget := widgetFromHWND(wb.hWnd)

	hdc := GetDC(wb.hWnd)
	defer ReleaseDC(wb.hWnd, hdc)

	hFont := widget.Font().handleForDPI(0) //HFONT(SendMessage(wb.hWnd, WM_GETFONT, 0, 0))
	hFontOld := SelectObject(hdc, HGDIOBJ(hFont))
	defer SelectObject(hdc, HGDIOBJ(hFontOld))

	var tm TEXTMETRIC
	if !GetTextMetrics(hdc, &tm) {
		newError("GetTextMetrics failed")
	}

	var size SIZE
	if !GetTextExtentPoint32(
		hdc,
		syscall.StringToUTF16Ptr("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
		52,
		&size) {
		newError("GetTextExtentPoint32 failed")
	}

	return Size{(size.CX/26 + 1) / 2, int(tm.TmHeight)}
}

func (wb *WidgetBase) dialogBaseUnitsToPixels(dlus Size) (pixels Size) {
	// FIXME: Cache dialog base units on font change.
	base := wb.dialogBaseUnits()

	return Size{MulDiv(dlus.Width, base.Width, 4), MulDiv(dlus.Height, base.Height, 8)}
}

func (wb *WidgetBase) LayoutFlags() LayoutFlags {
	return 0
}

func (wb *WidgetBase) MinSizeHint() Size {
	return wb.widget.SizeHint()
}

func (wb *WidgetBase) SizeHint() Size {
	return Size{10, 10}
}

func (wb *WidgetBase) calculateTextSize() Size {
	hdc := GetDC(wb.hWnd)
	if hdc == 0 {
		newError("GetDC failed")
		return Size{}
	}
	defer ReleaseDC(wb.hWnd, hdc)

	hFontOld := SelectObject(hdc, HGDIOBJ(wb.Font().handleForDPI(0)))
	defer SelectObject(hdc, hFontOld)

	var size Size
	lines := strings.Split(widgetText(wb.hWnd), "\n", -1)

	for _, line := range lines {
		var s SIZE
		str := syscall.StringToUTF16(strings.TrimRight(line, "\r "))

		if !GetTextExtentPoint32(hdc, &str[0], len(str)-1, &s) {
			newError("GetTextExtentPoint32 failed")
			return Size{}
		}

		size.Width = maxi(size.Width, s.CX)
		size.Height += s.CY
	}

	return size
}

func (wb *WidgetBase) updateParentLayout() os.Error {
	if wb.parent == nil || wb.parent.Layout() == nil {
		return nil
	}

	return wb.parent.Layout().Update(false)
}

func (wb *WidgetBase) Size() Size {
	return wb.Bounds().Size()
}

func (wb *WidgetBase) SetSize(size Size) os.Error {
	bounds := wb.Bounds()

	return wb.SetBounds(bounds.SetSize(size))
}

func (wb *WidgetBase) X() int {
	return wb.Bounds().X
}

func (wb *WidgetBase) SetX(value int) os.Error {
	bounds := wb.Bounds()
	bounds.X = value

	return wb.SetBounds(bounds)
}

func (wb *WidgetBase) Y() int {
	return wb.Bounds().Y
}

func (wb *WidgetBase) SetY(value int) os.Error {
	bounds := wb.Bounds()
	bounds.Y = value

	return wb.SetBounds(bounds)
}

func (wb *WidgetBase) Width() int {
	return wb.Bounds().Width
}

func (wb *WidgetBase) SetWidth(value int) os.Error {
	bounds := wb.Bounds()
	bounds.Width = value

	return wb.SetBounds(bounds)
}

func (wb *WidgetBase) Height() int {
	return wb.Bounds().Height
}

func (wb *WidgetBase) SetHeight(value int) os.Error {
	bounds := wb.Bounds()
	bounds.Height = value

	return wb.SetBounds(bounds)
}

func widgetClientBounds(hwnd HWND) Rectangle {
	var r RECT

	if !GetClientRect(hwnd, &r) {
		lastError("GetClientRect")
		return Rectangle{}
	}

	return Rectangle{X: r.Left, Y: r.Top, Width: r.Right - r.Left, Height: r.Bottom - r.Top}
}

func (wb *WidgetBase) ClientBounds() Rectangle {
	return widgetClientBounds(wb.hWnd)
}

func (wb *WidgetBase) sizeFromClientSize(clientSize Size) Size {
	s := wb.Size()
	cs := wb.ClientBounds().Size()
	ncs := Size{s.Width - cs.Width, s.Height - cs.Height}

	return Size{clientSize.Width + ncs.Width, clientSize.Height + ncs.Height}
}

func (wb *WidgetBase) SetClientSize(value Size) os.Error {
	return wb.SetSize(wb.sizeFromClientSize(value))
}

func (wb *WidgetBase) SetFocus() os.Error {
	if SetFocus(wb.hWnd) == 0 {
		return lastError("SetFocus")
	}

	return nil
}

func (wb *WidgetBase) CreateCanvas() (*Canvas, os.Error) {
	return newCanvasFromHWND(wb.hWnd)
}

func (wb *WidgetBase) setTheme(appName string) os.Error {
	if hr := SetWindowTheme(wb.hWnd, syscall.StringToUTF16Ptr(appName), nil); FAILED(hr) {
		return errorFromHRESULT("SetWindowTheme", hr)
	}

	return nil
}

func (wb *WidgetBase) KeyDown() *KeyEvent {
	return wb.keyDownPublisher.Event()
}

func (wb *WidgetBase) MouseDown() *MouseEvent {
	return wb.mouseDownPublisher.Event()
}

func (wb *WidgetBase) MouseMove() *MouseEvent {
	return wb.mouseMovePublisher.Event()
}

func (wb *WidgetBase) MouseUp() *MouseEvent {
	return wb.mouseUpPublisher.Event()
}

func (wb *WidgetBase) publishMouseEvent(publisher *MouseEventPublisher, wParam, lParam uintptr) {
	x := int(GET_X_LPARAM(lParam))
	y := int(GET_Y_LPARAM(lParam))

	publisher.Publish(x, y, 0)
}

func (wb *WidgetBase) SizeChanged() *Event {
	return wb.sizeChangedPublisher.Event()
}

func (wb *WidgetBase) persistState(restore bool) {
	settings := appSingleton.settings
	if settings != nil {
		widget := widgetFromHWND(wb.hWnd)
		if persistable, ok := widget.(Persistable); ok && persistable.Persistent() {
			if restore {
				persistable.RestoreState()
			} else {
				persistable.SaveState()
			}
		}
	}
}

func (wb *WidgetBase) getState() (string, os.Error) {
	settings := appSingleton.settings
	if settings == nil {
		return "", newError("App().Settings() must not be nil")
	}

	state, _ := settings.Get(wb.path())
	return state, nil
}

func (wb *WidgetBase) putState(state string) os.Error {
	settings := appSingleton.settings
	if settings == nil {
		return newError("App().Settings() must not be nil")
	}

	return settings.Put(wb.path(), state)
}

func widgetFromHWND(hwnd HWND) widgetInternal {
	ptr := GetWindowLongPtr(hwnd, GWLP_USERDATA)
	if ptr == 0 {
		return nil
	}

	wb := (*WidgetBase)(unsafe.Pointer(ptr))

	return wb.widget
}

func widgetWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) (result uintptr) {
	defer func() {
		var err os.Error
		if x := recover(); x != nil {
			if e, ok := x.(os.Error); ok {
				err = e
			} else {
				err = newError(fmt.Sprint(x))
			}
		}
		if err != nil {
			appSingleton.panickingPublisher.Publish(err)
		}
	}()

	wi := widgetFromHWND(hwnd)
	if wi == nil {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	result = wi.wndProc(hwnd, msg, wParam, lParam)

	return
}

func (wb *WidgetBase) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_ERASEBKGND:
		if wb.background == nil {
			break
		}

		canvas, err := newCanvasFromHDC(HDC(wParam))
		if err != nil {
			break
		}
		defer canvas.Dispose()

		if err := canvas.FillRectangle(wb.background, wb.ClientBounds()); err != nil {
			break
		}

		return 1

	case WM_LBUTTONDOWN:
		if _, isSubclassed := widgetFromHWND(wb.hWnd).(subclassedWidget); !isSubclassed {
			// Only call SetCapture if this is no subclassed control.
			// (Otherwise e.g. WM_COMMAND(BN_CLICKED) would no longer
			// be generated for PushButton.)
			SetCapture(wb.hWnd)
		}
		wb.publishMouseEvent(&wb.mouseDownPublisher, wParam, lParam)

	case WM_LBUTTONUP:
		if _, isSubclassed := widgetFromHWND(wb.hWnd).(subclassedWidget); !isSubclassed {
			// See WM_LBUTTONDOWN for why we require origWndProcPtr == 0 here.
			if !ReleaseCapture() {
				lastError("ReleaseCapture")
			}
		}
		wb.publishMouseEvent(&wb.mouseUpPublisher, wParam, lParam)

	case WM_MOUSEMOVE:
		wb.publishMouseEvent(&wb.mouseMovePublisher, wParam, lParam)

	case WM_SETCURSOR:
		if wb.cursor != nil {
			SetCursor(wb.cursor.handle())
			return 0
		}

	case WM_CONTEXTMENU:
		sourceWidget := widgetFromHWND(HWND(wParam))
		if sourceWidget == nil {
			break
		}

		x := int(GET_X_LPARAM(lParam))
		y := int(GET_Y_LPARAM(lParam))

		contextMenu := sourceWidget.ContextMenu()

		if contextMenu != nil {
			TrackPopupMenuEx(contextMenu.hMenu, TPM_NOANIMATION, x, y, rootWidget(sourceWidget).BaseWidget().hWnd, nil)
			return 0
		}

	case WM_KEYDOWN:
		wb.keyDownPublisher.Publish(int(wParam))

	case WM_SIZE, WM_SIZING:
		wb.sizeChangedPublisher.Publish()

	case WM_SHOWWINDOW:
		wb.persistState(wParam != 0)

	case WM_DESTROY:
		wb.persistState(false)
	}

	if widget := widgetFromHWND(hwnd); widget != nil {
		if subclassed, ok := widget.(subclassedWidget); ok {
			return CallWindowProc(subclassed.origWndProcPtr(), hwnd, msg, wParam, lParam)
		}
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

func (w *WidgetBase) runMessageLoop() int {
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
