// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

import (
	"github.com/lxn/win"
)

// App-specific message ids for internal use in Walk.
// TODO: Document reserved range somewhere (when we have an idea how many we need).
const (
	notifyIconMessageId = win.WM_APP + iota
)

// Window is an interface that provides operations common to all windows.
type Window interface {
	// AsWindowBase returns a *WindowBase, a pointer to an instance of the
	// struct that implements most operations common to all windows.
	AsWindowBase() *WindowBase

	// Background returns the background Brush of the Window.
	//
	// By default this is nil.
	Background() Brush

	// Bounds returns the outer bounding box Rectangle of the Window, including
	// decorations.
	//
	// For a Form, like *MainWindow or *Dialog, the Rectangle is in screen
	// coordinates, for a child Window the coordinates are relative to its
	// parent.
	Bounds() Rectangle

	// BringToTop moves the Window to the top of the keyboard focus order.
	BringToTop() error

	// ClientBounds returns the inner bounding box Rectangle of the Window,
	// excluding decorations.
	ClientBounds() Rectangle

	// ContextMenu returns the context menu of the Window.
	//
	// By default this is nil.
	ContextMenu() *Menu

	// CreateCanvas creates and returns a *Canvas that can be used to draw
	// inside the ClientBounds of the Window.
	//
	// Remember to call the Dispose method on the canvas to release resources,
	// when you no longer need it.
	CreateCanvas() (*Canvas, error)

	// Cursor returns the Cursor of the Window.
	//
	// By default this is nil.
	Cursor() Cursor

	// Dispose releases the operating system resources, associated with the
	// Window.
	//
	// If a user closes a *MainWindow or *Dialog, it is automatically released.
	// Also, if a Container is disposed of, all its descendants will be released
	// as well.
	Dispose()

	// Enabled returns if the Window is enabled for user interaction.
	Enabled() bool

	// Font returns the *Font of the Window.
	//
	// By default this is a MS Shell Dlg 2, 8 point font.
	Font() *Font

	// Handle returns the window handle of the Window.
	Handle() win.HWND

	// Height returns the outer height of the Window, including decorations.
	Height() int

	// Invalidate schedules a full repaint of the Window.
	Invalidate() error

	// IsDisposed returns if the Window has been disposed of.
	IsDisposed() bool

	// KeyDown returns a *KeyEvent that you can attach to for handling key down
	// events for the Window.
	KeyDown() *KeyEvent

	// KeyPress returns a *KeyEvent that you can attach to for handling key
	// press events for the Window.
	KeyPress() *KeyEvent

	// KeyUp returns a *KeyEvent that you can attach to for handling key up
	// events for the Window.
	KeyUp() *KeyEvent

	// MaxSize returns the maximum allowed outer Size for the Window, including
	// decorations.
	//
	// For child windows, this is only relevant when the parent of the Window
	// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
	MaxSize() Size

	// MinSize returns the minimum allowed outer Size for the Window, including
	// decorations.
	//
	// For child windows, this is only relevant when the parent of the Window
	// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
	MinSize() Size

	// MouseDown returns a *MouseEvent that you can attach to for handling
	// mouse down events for the Window.
	MouseDown() *MouseEvent

	// MouseMove returns a *MouseEvent that you can attach to for handling
	// mouse move events for the Window.
	MouseMove() *MouseEvent

	// MouseUp returns a *MouseEvent that you can attach to for handling
	// mouse up events for the Window.
	MouseUp() *MouseEvent

	// Name returns the name of the Window.
	Name() string

	// SendMessage sends a message to the window and returns the result.
	SendMessage(msg uint32, wParam, lParam uintptr) uintptr

	// SetBackground sets the background Brush of the Window.
	SetBackground(value Brush)

	// SetBounds sets the outer bounding box Rectangle of the Window, including
	// decorations.
	//
	// For a Form, like *MainWindow or *Dialog, the Rectangle is in screen
	// coordinates, for a child Window the coordinates are relative to its
	// parent.
	SetBounds(value Rectangle) error

	// SetClientSize sets the Size of the inner bounding box of the Window,
	// excluding decorations.
	SetClientSize(value Size) error

	// SetContextMenu sets the context menu of the Window.
	SetContextMenu(value *Menu)

	// SetCursor sets the Cursor of the Window.
	SetCursor(value Cursor)

	// SetEnabled sets if the Window is enabled for user interaction.
	SetEnabled(value bool)

	// SetFocus sets the keyboard input focus to the Window.
	SetFocus() error

	// SetFont sets the *Font of the Window.
	SetFont(value *Font)

	// SetHeight sets the outer height of the Window, including decorations.
	SetHeight(value int) error

	// SetMinMaxSize sets the minimum and maximum outer Size of the Window,
	// including decorations.
	//
	// Use walk.Size{} to make the respective limit be ignored.
	SetMinMaxSize(min, max Size) error

	// SetName sets the name of the Window.
	//
	// This is important if you want to make use of the built-in UI persistence.
	// Some windows support automatic state persistence. See Settings for
	// details.
	SetName(name string)

	// SetSize sets the outer Size of the Window, including decorations.
	SetSize(value Size) error

	// SetSuspended sets if the Window is suspended for layout and repainting
	// purposes.
	//
	// You should call SetSuspended(true), before doing a batch of modifications
	// that would cause multiple layout or drawing updates. Remember to call
	// SetSuspended(false) afterwards, which will update the Window accordingly.
	SetSuspended(suspend bool)

	// SetVisible sets if the Window is visible.
	SetVisible(value bool)

	// SetWidth sets the outer width of the Window, including decorations.
	SetWidth(value int) error

	// SetX sets the x coordinate of the Window, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for
	// child Windows.
	SetX(value int) error

	// SetY sets the y coordinate of the Window, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for
	// child Windows.
	SetY(value int) error

	// Size returns the outer Size of the Window, including decorations.
	Size() Size

	// SizeChanged returns an *Event that you can attach to for handling size
	// changed events for the Window.
	SizeChanged() *Event

	// Suspended returns if the Window is suspended for layout and repainting
	// purposes.
	Suspended() bool

	// Synchronize enqueues func f to be called some time later by the main
	// goroutine from inside a message loop.
	Synchronize(f func())

	// Visible returns if the Window is visible.
	Visible() bool

	// Width returns the outer width of the Window, including decorations.
	Width() int

	// WndProc is the window procedure of the window.
	//
	// When implementing your own WndProc to add or modify behavior, call the
	// WndProc of the embedded window for messages you don't handle yourself.
	WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr

	// X returns the x coordinate of the Window, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for
	// child Windows.
	X() int

	// Y returns the y coordinate of the Window, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for
	// child Windows.
	Y() int
}

// WindowBase implements many operations common to all Windows.
type WindowBase struct {
	window                  Window
	hWnd                    win.HWND
	origWndProcPtr          uintptr
	name                    string
	font                    *Font
	contextMenu             *Menu
	keyDownPublisher        KeyEventPublisher
	keyPressPublisher       KeyEventPublisher
	keyUpPublisher          KeyEventPublisher
	mouseDownPublisher      MouseEventPublisher
	mouseUpPublisher        MouseEventPublisher
	mouseMovePublisher      MouseEventPublisher
	sizeChangedPublisher    EventPublisher
	maxSize                 Size
	minSize                 Size
	background              Brush
	cursor                  Cursor
	suspended               bool
	visible                 bool
	enabled                 bool
	name2Property           map[string]Property
	enabledProperty         Property
	enabledChangedPublisher EventPublisher
	visibleProperty         Property
	visibleChangedPublisher EventPublisher
}

var (
	registeredWindowClasses = make(map[string]bool)
	defaultWndProcPtr       = syscall.NewCallback(defaultWndProc)
	hwnd2WindowBase         = make(map[win.HWND]*WindowBase)
)

// MustRegisterWindowClass registers the specified window class.
//
// MustRegisterWindowClass must be called once for every window type that is not
// based on any system provided control, before calling InitChildWidget or
// InitWidget. Calling MustRegisterWindowClass twice with the same className
// results in a panic.
func MustRegisterWindowClass(className string) {
	if registeredWindowClasses[className] {
		panic("window class already registered")
	}

	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		panic("GetModuleHandle")
	}

	hIcon := win.LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(win.IDI_APPLICATION))))
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = defaultWndProcPtr
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}

	registeredWindowClasses[className] = true
}

// InitWindow initializes a window.
//
// Widgets should be initialized using InitWidget instead.
func InitWindow(window, parent Window, className string, style, exStyle uint32) error {
	wb := window.AsWindowBase()
	wb.window = window
	wb.enabled = true
	wb.visible = true

	wb.name2Property = make(map[string]Property)

	var hwndParent win.HWND
	if parent != nil {
		hwndParent = parent.Handle()

		if widget, ok := window.(Widget); ok {
			if container, ok := parent.(Container); ok {
				widget.AsWidgetBase().parent = container
			}
		}
	}

	wb.hWnd = win.CreateWindowEx(
		exStyle,
		syscall.StringToUTF16Ptr(className),
		nil,
		style|win.WS_CLIPSIBLINGS,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
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

	hwnd2WindowBase[wb.hWnd] = wb

	if !registeredWindowClasses[className] {
		// We subclass all windows of system classes.
		wb.origWndProcPtr = win.SetWindowLongPtr(wb.hWnd, win.GWLP_WNDPROC, defaultWndProcPtr)
		if wb.origWndProcPtr == 0 {
			return lastError("SetWindowLongPtr")
		}
	}

	setWindowFont(wb.hWnd, defaultFont)

	if form, ok := window.(Form); ok {
		if fb := form.AsFormBase(); fb != nil {
			if err := fb.init(form); err != nil {
				return err
			}
		}
	}

	if container, ok := window.(Container); ok {
		if cb := container.AsContainerBase(); cb != nil {
			if err := cb.init(container); err != nil {
				return err
			}
		}
	}

	if widget, ok := window.(Widget); ok {
		if wb := widget.AsWidgetBase(); wb != nil {
			if err := wb.init(widget); err != nil {
				return err
			}
		}
	}

	wb.enabledProperty = NewBoolProperty(
		func() bool {
			return wb.window.Enabled()
		},
		func(b bool) error {
			wb.window.SetEnabled(b)
			return nil
		},
		wb.enabledChangedPublisher.Event())

	wb.visibleProperty = NewBoolProperty(
		func() bool {
			return window.Visible()
		},
		func(b bool) error {
			wb.window.SetVisible(b)
			return nil
		},
		wb.visibleChangedPublisher.Event())

	wb.MustRegisterProperty("Enabled", wb.enabledProperty)
	wb.MustRegisterProperty("Visible", wb.visibleProperty)

	succeeded = true

	return nil
}

// InitWrapperWindow initializes a window that wraps (embeds) another window.
//
// Calling this method is necessary, if you want to be able to override the
// WndProc method of the embedded window. The embedded window should only be
// used as inseparable part of the wrapper window to avoid undefined behavior.
func InitWrapperWindow(window Window) error {
	wb := window.AsWindowBase()

	wb.window = window

	if widget, ok := window.(Widget); ok {
		widgetBase := widget.AsWidgetBase()

		if widgetBase.parent != nil {
			children := widgetBase.parent.Children().items

			for i, w := range children {
				if w.AsWidgetBase() == widgetBase {
					children[i] = widget
					break
				}
			}
		}
	}

	return nil
}

func (wb *WindowBase) MustRegisterProperty(name string, property Property) {
	if property == nil {
		panic("property must not be nil")
	}
	if wb.name2Property[name] != nil {
		panic("property already registered")
	}

	wb.name2Property[name] = property
}

func (wb *WindowBase) Property(name string) Property {
	return wb.name2Property[name]
}

func (wb *WindowBase) hasStyleBits(bits uint) bool {
	style := uint(win.GetWindowLong(wb.hWnd, win.GWL_STYLE))

	return style&bits == bits
}

func (wb *WindowBase) setAndClearStyleBits(set, clear uint32) error {
	style := uint32(win.GetWindowLong(wb.hWnd, win.GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	var newStyle uint32
	newStyle = (style | set) &^ clear

	if newStyle != style {
		win.SetLastError(0)
		if win.SetWindowLong(wb.hWnd, win.GWL_STYLE, int32(newStyle)) == 0 {
			return lastError("SetWindowLong")
		}
	}

	return nil
}

func (wb *WindowBase) ensureStyleBits(bits uint32, set bool) error {
	var setBits uint32
	var clearBits uint32
	if set {
		setBits = bits
	} else {
		clearBits = bits
	}
	return wb.setAndClearStyleBits(setBits, clearBits)
}

// Handle returns the window handle of the Window.
func (wb *WindowBase) Handle() win.HWND {
	return wb.hWnd
}

// SendMessage sends a message to the window and returns the result.
func (wb *WindowBase) SendMessage(msg uint32, wParam, lParam uintptr) uintptr {
	return win.SendMessage(wb.hWnd, msg, wParam, lParam)
}

// Name returns the name of the *WindowBase.
func (wb *WindowBase) Name() string {
	return wb.name
}

// SetName sets the name of the *WindowBase.
func (wb *WindowBase) SetName(name string) {
	wb.name = name
}

func (wb *WindowBase) writePath(buf *bytes.Buffer) {
	hWndParent := win.GetAncestor(wb.hWnd, win.GA_PARENT)
	if pwi := windowFromHandle(hWndParent); pwi != nil {
		pwi.AsWindowBase().writePath(buf)
		buf.WriteByte('/')
	}

	buf.WriteString(wb.name)
}

func (wb *WindowBase) path() string {
	buf := bytes.NewBuffer(nil)

	wb.writePath(buf)

	return buf.String()
}

// WindowBase simply returns the receiver.
func (wb *WindowBase) AsWindowBase() *WindowBase {
	return wb
}

// Dispose releases the operating system resources, associated with the
// *WindowBase.
//
// If a user closes a *MainWindow or *Dialog, it is automatically released.
// Also, if a Container is disposed of, all its descendants will be released
// as well.
func (wb *WindowBase) Dispose() {
	hWnd := wb.hWnd
	if hWnd != 0 {
		switch w := wb.window.(type) {
		case *ToolTip:
		case Widget:
			globalToolTip.RemoveTool(w)
		}

		wb.hWnd = 0
		win.DestroyWindow(hWnd)
	}

	if cm := wb.contextMenu; cm != nil {
		cm.actions.Clear()
	}

	for _, p := range wb.name2Property {
		p.SetSource(nil)
	}
}

// IsDisposed returns if the *WindowBase has been disposed of.
func (wb *WindowBase) IsDisposed() bool {
	return wb.hWnd == 0
}

// ContextMenu returns the context menu of the *WindowBase.
//
// By default this is nil.
func (wb *WindowBase) ContextMenu() *Menu {
	return wb.contextMenu
}

// SetContextMenu sets the context menu of the *WindowBase.
func (wb *WindowBase) SetContextMenu(value *Menu) {
	wb.contextMenu = value
}

// Background returns the background Brush of the *WindowBase.
//
// By default this is nil.
func (wb *WindowBase) Background() Brush {
	return wb.background
}

// SetBackground sets the background Brush of the *WindowBase.
func (wb *WindowBase) SetBackground(value Brush) {
	wb.background = value

	wb.Invalidate()
}

// Cursor returns the Cursor of the *WindowBase.
//
// By default this is nil.
func (wb *WindowBase) Cursor() Cursor {
	return wb.cursor
}

// SetCursor sets the Cursor of the *WindowBase.
func (wb *WindowBase) SetCursor(value Cursor) {
	wb.cursor = value
}

// Enabled returns if the *WindowBase is enabled for user interaction.
func (wb *WindowBase) Enabled() bool {
	return wb.enabled
}

// SetEnabled sets if the *WindowBase is enabled for user interaction.
func (wb *WindowBase) SetEnabled(value bool) {
	wb.enabled = value

	win.EnableWindow(wb.hWnd, wb.window.Enabled())

	wb.enabledChangedPublisher.Publish()
}

// Font returns the *Font of the *WindowBase.
//
// By default this is a MS Shell Dlg 2, 8 point font.
func (wb *WindowBase) Font() *Font {
	if wb.font != nil {
		return wb.font
	}

	return defaultFont
}

func setWindowFont(hwnd win.HWND, font *Font) {
	win.SendMessage(hwnd, win.WM_SETFONT, uintptr(font.handleForDPI(0)), 1)
}

// SetFont sets the *Font of the *WindowBase.
func (wb *WindowBase) SetFont(value *Font) {
	if value != wb.font {
		setWindowFont(wb.hWnd, value)

		wb.font = value
	}
}

// Suspended returns if the *WindowBase is suspended for layout and repainting
// purposes.
func (wb *WindowBase) Suspended() bool {
	return wb.suspended
}

// SetSuspended sets if the *WindowBase is suspended for layout and repainting
// purposes.
//
// You should call SetSuspended(true), before doing a batch of modifications
// that would cause multiple layout or drawing updates. Remember to call
// SetSuspended(false) afterwards, which will update the *WindowBase
// accordingly.
func (wb *WindowBase) SetSuspended(suspend bool) {
	if suspend == wb.suspended {
		return
	}

	var wParam int
	if suspend {
		wParam = 0
	} else {
		wParam = 1
	}

	wb.SendMessage(win.WM_SETREDRAW, uintptr(wParam), 0)

	wb.suspended = suspend
}

// Invalidate schedules a full repaint of the *WindowBase.
func (wb *WindowBase) Invalidate() error {
	if !win.InvalidateRect(wb.hWnd, nil, true) {
		return newError("InvalidateRect failed")
	}

	return nil
}

func windowText(hwnd win.HWND) string {
	textLength := win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0)
	buf := make([]uint16, textLength+1)
	win.SendMessage(hwnd, win.WM_GETTEXT, uintptr(textLength+1), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func setWindowText(hwnd win.HWND, text string) error {
	if win.TRUE != win.SendMessage(hwnd, win.WM_SETTEXT, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text)))) {
		return newError("WM_SETTEXT failed")
	}

	return nil
}

// Visible returns if the *WindowBase is visible.
func (wb *WindowBase) Visible() bool {
	return win.IsWindowVisible(wb.hWnd)
}

// SetVisible sets if the *WindowBase is visible.
func (wb *WindowBase) SetVisible(visible bool) {
	var cmd int32
	if visible {
		cmd = win.SW_SHOW
	} else {
		cmd = win.SW_HIDE
	}
	win.ShowWindow(wb.hWnd, cmd)

	wb.visible = visible

	if widget, ok := wb.window.(Widget); ok {
		widget.AsWidgetBase().updateParentLayout()
	}

	wb.visibleChangedPublisher.Publish()
}

// BringToTop moves the *WindowBase to the top of the keyboard focus order.
func (wb *WindowBase) BringToTop() error {
	if !win.SetWindowPos(wb.hWnd, win.HWND_TOP, 0, 0, 0, 0, win.SWP_NOACTIVATE|win.SWP_NOMOVE|win.SWP_NOSIZE) {
		return lastError("SetWindowPos")
	}

	return nil
}

// Bounds returns the outer bounding box Rectangle of the *WindowBase, including
// decorations.
//
// The coordinates are relative to the screen.
func (wb *WindowBase) Bounds() Rectangle {
	var r win.RECT

	if !win.GetWindowRect(wb.hWnd, &r) {
		lastError("GetWindowRect")
		return Rectangle{}
	}

	return Rectangle{
		int(r.Left),
		int(r.Top),
		int(r.Right - r.Left),
		int(r.Bottom - r.Top),
	}
}

// SetBounds returns the outer bounding box Rectangle of the *WindowBase,
// including decorations.
//
// For a Form, like *MainWindow or *Dialog, the Rectangle is in screen
// coordinates, for a child Window the coordinates are relative to its parent.
func (wb *WindowBase) SetBounds(bounds Rectangle) error {
	if !win.MoveWindow(
		wb.hWnd,
		int32(bounds.X),
		int32(bounds.Y),
		int32(bounds.Width),
		int32(bounds.Height),
		true) {

		return lastError("MoveWindow")
	}

	return nil
}

// MinSize returns the minimum allowed outer Size for the *WindowBase, including
// decorations.
//
// For child windows, this is only relevant when the parent of the *WindowBase
// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
func (wb *WindowBase) MinSize() Size {
	return wb.minSize
}

// MaxSize returns the maximum allowed outer Size for the *WindowBase, including
// decorations.
//
// For child windows, this is only relevant when the parent of the *WindowBase
// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
func (wb *WindowBase) MaxSize() Size {
	return wb.maxSize
}

// SetMinMaxSize sets the minimum and maximum outer Size of the *WindowBase,
// including decorations.
//
// Use walk.Size{} to make the respective limit be ignored.
func (wb *WindowBase) SetMinMaxSize(min, max Size) error {
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

var dialogBaseUnitsUTF16StringPtr = syscall.StringToUTF16Ptr("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func (wb *WindowBase) dialogBaseUnits() Size {
	// The window may use a font different from that in WindowBase,
	// like e.g. NumberEdit does, so we try to use the right one.
	window := windowFromHandle(wb.hWnd)

	hdc := win.GetDC(wb.hWnd)
	defer win.ReleaseDC(wb.hWnd, hdc)

	hFont := window.Font().handleForDPI(0)
	hFontOld := win.SelectObject(hdc, win.HGDIOBJ(hFont))
	defer win.SelectObject(hdc, win.HGDIOBJ(hFontOld))

	var tm win.TEXTMETRIC
	if !win.GetTextMetrics(hdc, &tm) {
		newError("GetTextMetrics failed")
	}

	var size win.SIZE
	if !win.GetTextExtentPoint32(
		hdc,
		dialogBaseUnitsUTF16StringPtr,
		52,
		&size) {
		newError("GetTextExtentPoint32 failed")
	}

	return Size{int((size.CX/26 + 1) / 2), int(tm.TmHeight)}
}

func (wb *WindowBase) dialogBaseUnitsToPixels(dlus Size) (pixels Size) {
	// FIXME: Cache dialog base units on font change.
	base := wb.dialogBaseUnits()

	return Size{
		int(win.MulDiv(int32(dlus.Width), int32(base.Width), 4)),
		int(win.MulDiv(int32(dlus.Height), int32(base.Height), 8)),
	}
}

func (wb *WindowBase) calculateTextSizeImpl(text string) Size {
	hdc := win.GetDC(wb.hWnd)
	if hdc == 0 {
		newError("GetDC failed")
		return Size{}
	}
	defer win.ReleaseDC(wb.hWnd, hdc)

	hFontOld := win.SelectObject(hdc, win.HGDIOBJ(wb.Font().handleForDPI(0)))
	defer win.SelectObject(hdc, hFontOld)

	var size Size
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		var s win.SIZE
		str := syscall.StringToUTF16(strings.TrimRight(line, "\r "))

		if !win.GetTextExtentPoint32(hdc, &str[0], int32(len(str)-1), &s) {
			newError("GetTextExtentPoint32 failed")
			return Size{}
		}

		size.Width = maxi(size.Width, int(s.CX))
		size.Height += int(s.CY)
	}

	return size
}

func (wb *WindowBase) calculateTextSize() Size {
	return wb.calculateTextSizeImpl(windowText(wb.hWnd))
}

// Size returns the outer Size of the *WindowBase, including decorations.
func (wb *WindowBase) Size() Size {
	return wb.window.Bounds().Size()
}

// SetSize sets the outer Size of the *WindowBase, including decorations.
func (wb *WindowBase) SetSize(size Size) error {
	bounds := wb.window.Bounds()

	return wb.SetBounds(bounds.SetSize(size))
}

// X returns the x coordinate of the *WindowBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for
// child Windows.
func (wb *WindowBase) X() int {
	return wb.window.Bounds().X
}

// SetX sets the x coordinate of the *WindowBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for
// child Windows.
func (wb *WindowBase) SetX(value int) error {
	bounds := wb.window.Bounds()
	bounds.X = value

	return wb.SetBounds(bounds)
}

// Y returns the y coordinate of the *WindowBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for
// child Windows.
func (wb *WindowBase) Y() int {
	return wb.window.Bounds().Y
}

// SetY sets the y coordinate of the *WindowBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for
// child Windows.
func (wb *WindowBase) SetY(value int) error {
	bounds := wb.window.Bounds()
	bounds.Y = value

	return wb.SetBounds(bounds)
}

// Width returns the outer width of the *WindowBase, including decorations.
func (wb *WindowBase) Width() int {
	return wb.window.Bounds().Width
}

// SetWidth sets the outer width of the *WindowBase, including decorations.
func (wb *WindowBase) SetWidth(value int) error {
	bounds := wb.window.Bounds()
	bounds.Width = value

	return wb.SetBounds(bounds)
}

// Height returns the outer height of the *WindowBase, including decorations.
func (wb *WindowBase) Height() int {
	return wb.window.Bounds().Height
}

// SetHeight sets the outer height of the *WindowBase, including decorations.
func (wb *WindowBase) SetHeight(value int) error {
	bounds := wb.window.Bounds()
	bounds.Height = value

	return wb.SetBounds(bounds)
}

func windowClientBounds(hwnd win.HWND) Rectangle {
	var r win.RECT

	if !win.GetClientRect(hwnd, &r) {
		lastError("GetClientRect")
		return Rectangle{}
	}

	return Rectangle{
		int(r.Left),
		int(r.Top),
		int(r.Right - r.Left),
		int(r.Bottom - r.Top),
	}
}

// ClientBounds returns the inner bounding box Rectangle of the *WindowBase,
// excluding decorations.
func (wb *WindowBase) ClientBounds() Rectangle {
	return windowClientBounds(wb.hWnd)
}

func (wb *WindowBase) sizeFromClientSize(clientSize Size) Size {
	s := wb.Size()
	cs := wb.ClientBounds().Size()
	ncs := Size{s.Width - cs.Width, s.Height - cs.Height}

	return Size{clientSize.Width + ncs.Width, clientSize.Height + ncs.Height}
}

// SetClientSize sets the Size of the inner bounding box of the *WindowBase,
// excluding decorations.
func (wb *WindowBase) SetClientSize(value Size) error {
	return wb.SetSize(wb.sizeFromClientSize(value))
}

// SetFocus sets the keyboard input focus to the *WindowBase.
func (wb *WindowBase) SetFocus() error {
	if win.SetFocus(wb.hWnd) == 0 {
		return lastError("SetFocus")
	}

	return nil
}

// CreateCanvas creates and returns a *Canvas that can be used to draw
// inside the ClientBounds of the *WindowBase.
//
// Remember to call the Dispose method on the canvas to release resources,
// when you no longer need it.
func (wb *WindowBase) CreateCanvas() (*Canvas, error) {
	return newCanvasFromHWND(wb.hWnd)
}

func (wb *WindowBase) setTheme(appName string) error {
	if hr := win.SetWindowTheme(wb.hWnd, syscall.StringToUTF16Ptr(appName), nil); win.FAILED(hr) {
		return errorFromHRESULT("SetWindowTheme", hr)
	}

	return nil
}

// KeyDown returns a *KeyEvent that you can attach to for handling key down
// events for the *WindowBase.
func (wb *WindowBase) KeyDown() *KeyEvent {
	return wb.keyDownPublisher.Event()
}

// KeyPress returns a *KeyEvent that you can attach to for handling key press
// events for the *WindowBase.
func (wb *WindowBase) KeyPress() *KeyEvent {
	return wb.keyPressPublisher.Event()
}

// KeyUp returns a *KeyEvent that you can attach to for handling key up
// events for the *WindowBase.
func (wb *WindowBase) KeyUp() *KeyEvent {
	return wb.keyUpPublisher.Event()
}

// MouseDown returns a *MouseEvent that you can attach to for handling
// mouse down events for the *WindowBase.
func (wb *WindowBase) MouseDown() *MouseEvent {
	return wb.mouseDownPublisher.Event()
}

// MouseMove returns a *MouseEvent that you can attach to for handling
// mouse move events for the *WindowBase.
func (wb *WindowBase) MouseMove() *MouseEvent {
	return wb.mouseMovePublisher.Event()
}

// MouseUp returns a *MouseEvent that you can attach to for handling
// mouse up events for the *WindowBase.
func (wb *WindowBase) MouseUp() *MouseEvent {
	return wb.mouseUpPublisher.Event()
}

func (wb *WindowBase) publishMouseEvent(publisher *MouseEventPublisher, wParam, lParam uintptr) {
	x := int(win.GET_X_LPARAM(lParam))
	y := int(win.GET_Y_LPARAM(lParam))

	var button MouseButton
	switch true {
	case wParam&win.MK_LBUTTON > 0:
		button = LeftButton

	case wParam&win.MK_MBUTTON > 0:
		button = MiddleButton

	case wParam&win.MK_RBUTTON > 0:
		button = RightButton
	}

	publisher.Publish(x, y, button)
}

// SizeChanged returns an *Event that you can attach to for handling size
// changed events for the *WindowBase.
func (wb *WindowBase) SizeChanged() *Event {
	return wb.sizeChangedPublisher.Event()
}

// Synchronize enqueues func f to be called some time later by the main
// goroutine from inside a message loop.
func (wb *WindowBase) Synchronize(f func()) {
	synchronize(f)

	win.PostMessage(wb.hWnd, syncMsgId, 0, 0)
}

func (wb *WindowBase) getState() (string, error) {
	settings := appSingleton.settings
	if settings == nil {
		return "", newError("App().Settings() must not be nil")
	}

	state, _ := settings.Get(wb.path())
	return state, nil
}

func (wb *WindowBase) putState(state string) error {
	settings := appSingleton.settings
	if settings == nil {
		return newError("App().Settings() must not be nil")
	}

	p := wb.path()
	if strings.HasPrefix(p, "/") ||
		strings.HasSuffix(p, "/") ||
		strings.Contains(p, "//") {

		return nil
	}

	return settings.PutExpiring(p, state)
}

func windowFromHandle(hwnd win.HWND) Window {
	if wb := hwnd2WindowBase[hwnd]; wb != nil {
		return wb.window
	}

	return nil
}

func defaultWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	defer func() {
		if len(appSingleton.panickingPublisher.event.handlers) > 0 {
			var err error
			if x := recover(); x != nil {
				if e, ok := x.(error); ok {
					err = wrapErrorNoPanic(e)
				} else {
					err = newErrorNoPanic(fmt.Sprint(x))
				}
			}
			if err != nil {
				appSingleton.panickingPublisher.Publish(err)
			}
		}
	}()

	if msg == notifyIconMessageId {
		return notifyIconWndProc(hwnd, msg, wParam, lParam)
	}

	wi := windowFromHandle(hwnd)
	if wi == nil {
		return win.DefWindowProc(hwnd, msg, wParam, lParam)
	}

	result = wi.WndProc(hwnd, msg, wParam, lParam)

	return
}

func (wb *WindowBase) handleKeyDown(wParam, lParam uintptr) {
	key := Key(wParam)

	if uint32(lParam)>>30 == 0 {
		wb.keyDownPublisher.Publish(key)

		// Using TranslateAccelerators refused to work, so we handle them
		// ourselves, at least for now.
		shortcut := Shortcut{ModifiersDown(), key}
		if action, ok := shortcut2Action[shortcut]; ok {
			if action.Visible() && action.Enabled() {
				action.raiseTriggered()
			}
		}
	}

	switch key {
	case KeyAlt, KeyControl, KeyShift:
		// nop

	default:
		wb.keyPressPublisher.Publish(key)
	}
}

func (wb *WindowBase) handleKeyUp(wParam, lParam uintptr) {
	wb.keyUpPublisher.Publish(Key(wParam))
}

// WndProc is the window procedure of the window.
//
// When implementing your own WndProc to add or modify behavior, call the
// WndProc of the embedded window for messages you don't handle yourself.
func (wb *WindowBase) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_ERASEBKGND:
		if wb.background == nil {
			break
		}

		canvas, err := newCanvasFromHDC(win.HDC(wParam))
		if err != nil {
			break
		}
		defer canvas.Dispose()

		if err := canvas.FillRectangle(wb.background, wb.ClientBounds()); err != nil {
			break
		}

		return 1

	case win.WM_LBUTTONDOWN, win.WM_MBUTTONDOWN, win.WM_RBUTTONDOWN:
		if msg == win.WM_LBUTTONDOWN && wb.origWndProcPtr == 0 {
			// Only call SetCapture if this is no subclassed control.
			// (Otherwise e.g. WM_COMMAND(BN_CLICKED) would no longer
			// be generated for PushButton.)
			win.SetCapture(wb.hWnd)
		}
		wb.publishMouseEvent(&wb.mouseDownPublisher, wParam, lParam)

	case win.WM_LBUTTONUP, win.WM_MBUTTONUP, win.WM_RBUTTONUP:
		if msg == win.WM_LBUTTONUP && wb.origWndProcPtr == 0 {
			// See WM_LBUTTONDOWN for why we require origWndProcPtr == 0 here.
			if !win.ReleaseCapture() {
				lastError("ReleaseCapture")
			}
		}
		wb.publishMouseEvent(&wb.mouseUpPublisher, wParam, lParam)

	case win.WM_MOUSEMOVE:
		wb.publishMouseEvent(&wb.mouseMovePublisher, wParam, lParam)

	case win.WM_SETCURSOR:
		if wb.cursor != nil {
			win.SetCursor(wb.cursor.handle())
			return 0
		}

	case win.WM_CONTEXTMENU:
		sourceWindow := windowFromHandle(win.HWND(wParam))
		if sourceWindow == nil {
			break
		}

		x := win.GET_X_LPARAM(lParam)
		y := win.GET_Y_LPARAM(lParam)

		contextMenu := sourceWindow.ContextMenu()

		var handle win.HWND
		if widget, ok := sourceWindow.(Widget); ok {
			handle = ancestor(widget).Handle()
		} else {
			handle = sourceWindow.Handle()
		}

		if contextMenu != nil {
			win.TrackPopupMenuEx(
				contextMenu.hMenu,
				win.TPM_NOANIMATION,
				x,
				y,
				handle,
				nil)
			return 0
		}

	case win.WM_KEYDOWN:
		wb.handleKeyDown(wParam, lParam)

	case win.WM_KEYUP:
		wb.handleKeyUp(wParam, lParam)

	case win.WM_SIZE, win.WM_SIZING:
		wb.sizeChangedPublisher.Publish()

	case win.WM_DESTROY:
		switch w := wb.window.(type) {
		case *ToolTip:
		case Widget:
			globalToolTip.RemoveTool(w)
		}

		delete(hwnd2WindowBase, hwnd)

		wb.hWnd = 0
		wb.window.Dispose()
	}

	if window := windowFromHandle(hwnd); window != nil {
		origWndProcPtr := window.AsWindowBase().origWndProcPtr
		if origWndProcPtr != 0 {
			return win.CallWindowProc(origWndProcPtr, hwnd, msg, wParam, lParam)
		}
	}

	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}
