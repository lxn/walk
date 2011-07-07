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

import . "walk/winapi"

// App-specific message ids for internal use in Walk.
// TODO: Document reserved range somewhere (when we have an idea how many we need).
const (
	notifyIconMessageId = WM_APP + iota
)

// LayoutFlags specify how a Widget wants to be treated when used with a Layout.
// 
// These flags are interpreted in respect to Widget.SizeHint.
type LayoutFlags byte

const (
	// ShrinkableHorz allows a Widget to be shrunk horizontally.
	ShrinkableHorz LayoutFlags = 1 << iota

	// ShrinkableVert allows a Widget to be shrunk vertically.
	ShrinkableVert

	// GrowableHorz allows a Widget to be enlarged horizontally.
	GrowableHorz

	// GrowableVert allows a Widget to be enlarged vertically.
	GrowableVert

	// GreedyHorz specifies that the widget prefers to take up as much space as
	// possible, horizontally.
	GreedyHorz

	// GreedyVert specifies that the widget prefers to take up as much space as
	// possible, vertically.
	GreedyVert
)

// Widget is an interface that provides operations common to all widgets.
type Widget interface {
	// Background returns the background Brush of the Widget.
	//
	// By default this is nil.
	Background() Brush

	// BaseWidget returns a *WidgetBase, a pointer to an instance of the struct 
	// that implements most operations common to all widgets.
	BaseWidget() *WidgetBase

	// Bounds returns the outer bounding box Rectangle of the Widget, including
	// decorations.
	//
	// For a RootWidget, like *MainWindow or *Dialog, the Rectangle is in screen
	// coordinates, for a child Widget the coordinates are relative to its 
	// parent.
	Bounds() Rectangle

	// BringToTop moves the Widget to the top of the keyboard focus order.
	BringToTop() os.Error

	// ClientBounds returns the inner bounding box Rectangle of the Widget,
	// excluding decorations.
	ClientBounds() Rectangle

	// ContextMenu returns the context menu of the Widget.
	//
	// By default this is nil.
	ContextMenu() *Menu

	// CreateCanvas creates and returns a *Canvas that can be used to draw
	// inside the ClientBounds of the Widget.
	//
	// Remember to call the Dispose method on the canvas to release resources, 
	// when you no longer need it. 
	CreateCanvas() (*Canvas, os.Error)

	// Cursor returns the Cursor of the Widget.
	//
	// By default this is nil.
	Cursor() Cursor

	// Dispose releases the operating system resources, associated with the 
	// Widget.
	//
	// If a user closes a *MainWindow or *Dialog, it is automatically released.
	// Also, if a Container is disposed of, all its descendants will be released
	// as well.
	Dispose()

	// Enabled returns if the Widget is enabled for user interaction.
	Enabled() bool

	// Font returns the *Font of the Widget.
	//
	// By default this is a MS Shell Dlg 2, 8 point font.
	Font() *Font

	// Height returns the outer height of the Widget, including decorations.
	Height() int

	// Invalidate schedules a full repaint of the Widget.
	Invalidate() os.Error

	// IsDisposed returns if the Widget has been disposed of.
	IsDisposed() bool

	// KeyDown returns a *KeyEvent that you can attach to for handling key down
	// events for the Widget.
	KeyDown() *KeyEvent

	// LayoutFlags returns a combination of LayoutFlags that specify how the
	// Widget wants to be treated by Layout implementations.
	LayoutFlags() LayoutFlags

	// MaxSize returns the maximum allowed outer Size for the Widget, including
	// decorations.
	//
	// For child widgets, this is only relevant when the parent of the Widget 
	// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
	MaxSize() Size

	// MinSize returns the minimum allowed outer Size for the Widget, including
	// decorations.
	//
	// For child widgets, this is only relevant when the parent of the Widget 
	// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
	MinSize() Size

	// MinSizeHint returns the minimum outer Size, including decorations, that 
	// makes sense for the respective type of Widget.
	MinSizeHint() Size

	// MouseDown returns a *MouseEvent that you can attach to for handling 
	// mouse down events for the Widget.
	MouseDown() *MouseEvent

	// MouseMove returns a *MouseEvent that you can attach to for handling 
	// mouse move events for the Widget.
	MouseMove() *MouseEvent

	// MouseUp returns a *MouseEvent that you can attach to for handling 
	// mouse up events for the Widget.
	MouseUp() *MouseEvent

	// Name returns the name of the Widget.
	Name() string

	// Parent returns the Container of the Widget.
	//
	// For RootWidgets, like *MainWindow and *Dialog, this is always nil.
	Parent() Container

	// RootWidget returns the root of the UI hierarchy of the Widget, which is
	// usually a *MainWindow or *Dialog.
	RootWidget() RootWidget

	// SetBackground sets the background Brush of the Widget.
	SetBackground(value Brush)

	// SetBounds sets the outer bounding box Rectangle of the Widget, including
	// decorations.
	//
	// For a RootWidget, like *MainWindow or *Dialog, the Rectangle is in screen
	// coordinates, for a child Widget the coordinates are relative to its 
	// parent.
	SetBounds(value Rectangle) os.Error

	// SetClientSize sets the Size of the inner bounding box of the Widget,
	// excluding decorations.
	SetClientSize(value Size) os.Error

	// SetContextMenu sets the context menu of the Widget.
	SetContextMenu(value *Menu)

	// SetCursor sets the Cursor of the Widget.
	SetCursor(value Cursor)

	// SetEnabled sets if the Widget is enabled for user interaction.
	SetEnabled(value bool)

	// SetFocus sets the keyboard input focus to the Widget.
	SetFocus() os.Error

	// SetFont sets the *Font of the Widget.
	SetFont(value *Font)

	// SetHeight sets the outer height of the Widget, including decorations.
	SetHeight(value int) os.Error

	// SetMinMaxSize sets the minimum and maximum outer Size of the Widget,
	// including decorations.
	//
	// Use walk.Size{} to make the respective limit be ignored.
	SetMinMaxSize(min, max Size) os.Error

	// SetName sets the name of the Widget.
	//
	// This is important if you want to make use of the built-in UI persistence.
	// Some widgets support automatic state persistence. See Settings for 
	// details.
	SetName(name string)

	// SetParent sets the parent of the Widget and adds the Widget to the 
	// Children list of the Container.
	SetParent(value Container) os.Error

	// SetSize sets the outer Size of the Widget, including decorations.
	SetSize(value Size) os.Error

	// SetSuspended sets if the Widget is suspended for layout and repainting 
	// purposes.
	//
	// You should call SetSuspended(true), before doing a batch of modifications
	// that would cause multiple layout or drawing updates. Remember to call
	// SetSuspended(false) afterwards, which will update the Widget accordingly.
	SetSuspended(suspend bool)

	// SetVisible sets if the Widget is visible.
	SetVisible(value bool)

	// SetWidth sets the outer width of the Widget, including decorations.
	SetWidth(value int) os.Error

	// SetX sets the x coordinate of the Widget, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
	// child Widgets.
	SetX(value int) os.Error

	// SetY sets the y coordinate of the Widget, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
	// child Widgets.
	SetY(value int) os.Error

	// Size returns the outer Size of the Widget, including decorations.
	Size() Size

	// SizeChanged returns an *Event that you can attach to for handling size
	// changed events for the Widget.
	SizeChanged() *Event

	// SizeHint returns the preferred Size for the respective type of Widget.
	SizeHint() Size

	// Suspended returns if the Widget is suspended for layout and repainting 
	// purposes.
	Suspended() bool

	// Visible returns if the Widget is visible.
	Visible() bool

	// Width returns the outer width of the Widget, including decorations.
	Width() int

	// X returns the x coordinate of the Widget, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
	// child Widgets.
	X() int

	// Y returns the y coordinate of the Widget, relative to the screen for
	// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
	// child Widgets.
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

// WidgetBase implements many operations common to all Widgets.
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

func (wb *WidgetBase) hasStyleBits(bits uint) bool {
	style := uint(GetWindowLong(wb.hWnd, GWL_STYLE))

	return style&bits == bits
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

// Name returns the name of the *WidgetBase.
func (wb *WidgetBase) Name() string {
	return wb.name
}

// SetName sets the name of the *WidgetBase.
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

// BaseWidget simply returns the receiver.
func (wb *WidgetBase) BaseWidget() *WidgetBase {
	return wb
}

// Dispose releases the operating system resources, associated with the 
// *WidgetBase.
//
// If a user closes a *MainWindow or *Dialog, it is automatically released.
// Also, if a Container is disposed of, all its descendants will be released
// as well.
func (wb *WidgetBase) Dispose() {
	if wb.hWnd != 0 {
		DestroyWindow(wb.hWnd)
		wb.hWnd = 0
	}
}

// IsDisposed returns if the *WidgetBase has been disposed of.
func (wb *WidgetBase) IsDisposed() bool {
	return wb.hWnd == 0
}

// RootWidget returns the root of the UI hierarchy of the *WidgetBase, which is
// usually a *MainWindow or *Dialog.
func (wb *WidgetBase) RootWidget() RootWidget {
	return rootWidget(wb)
}

// ContextMenu returns the context menu of the *WidgetBase.
//
// By default this is nil.
func (wb *WidgetBase) ContextMenu() *Menu {
	return wb.contextMenu
}

// SetContextMenu sets the context menu of the *WidgetBase.
func (wb *WidgetBase) SetContextMenu(value *Menu) {
	wb.contextMenu = value
}

// Background returns the background Brush of the *WidgetBase.
//
// By default this is nil.
func (wb *WidgetBase) Background() Brush {
	return wb.background
}

// SetBackground sets the background Brush of the *WidgetBase.
func (wb *WidgetBase) SetBackground(value Brush) {
	wb.background = value
}

// Cursor returns the Cursor of the *WidgetBase.
//
// By default this is nil.
func (wb *WidgetBase) Cursor() Cursor {
	return wb.cursor
}

// SetCursor sets the Cursor of the *WidgetBase.
func (wb *WidgetBase) SetCursor(value Cursor) {
	wb.cursor = value
}

// Enabled returns if the *WidgetBase is enabled for user interaction.
func (wb *WidgetBase) Enabled() bool {
	return IsWindowEnabled(wb.hWnd)
}

// SetEnabled sets if the *WidgetBase is enabled for user interaction.
func (wb *WidgetBase) SetEnabled(value bool) {
	EnableWindow(wb.hWnd, value)
}

// Font returns the *Font of the *WidgetBase.
//
// By default this is a MS Shell Dlg 2, 8 point font.
func (wb *WidgetBase) Font() *Font {
	return wb.font
}

func setWidgetFont(hwnd HWND, font *Font) {
	SendMessage(hwnd, WM_SETFONT, uintptr(font.handleForDPI(0)), 1)
}

// SetFont sets the *Font of the *WidgetBase.
func (wb *WidgetBase) SetFont(value *Font) {
	if value != wb.font {
		setWidgetFont(wb.hWnd, value)

		wb.font = value
	}
}

// Suspended returns if the *WidgetBase is suspended for layout and repainting 
// purposes.
func (wb *WidgetBase) Suspended() bool {
	return wb.suspended
}

// SetSuspended sets if the *WidgetBase is suspended for layout and repainting 
// purposes.
//
// You should call SetSuspended(true), before doing a batch of modifications
// that would cause multiple layout or drawing updates. Remember to call
// SetSuspended(false) afterwards, which will update the *WidgetBase 
// accordingly.
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

// Invalidate schedules a full repaint of the *WidgetBase.
func (wb *WidgetBase) Invalidate() os.Error {
	if !InvalidateRect(wb.hWnd, nil, true) {
		return newError("InvalidateRect failed")
	}

	return nil
}

// Parent returns the Container of the *WidgetBase.
//
// For RootWidgets, like *MainWindow and *Dialog, this is always nil.
func (wb *WidgetBase) Parent() Container {
	return wb.parent
}

// SetParent sets the parent of the *WidgetBase and adds the *WidgetBase to the 
// Children list of the Container.
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

// Visible returns if the *WidgetBase is visible.
func (wb *WidgetBase) Visible() bool {
	return IsWindowVisible(wb.hWnd)
}

// SetVisible sets if the *WidgetBase is visible.
func (wb *WidgetBase) SetVisible(visible bool) {
	var cmd int
	if visible {
		cmd = SW_SHOW
	} else {
		cmd = SW_HIDE
	}
	ShowWindow(wb.hWnd, cmd)
}

// BringToTop moves the *WidgetBase to the top of the keyboard focus order.
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

// Bounds returns the outer bounding box Rectangle of the *WidgetBase, including
// decorations.
//
// For a RootWidget, like *MainWindow or *Dialog, the Rectangle is in screen
// coordinates, for a child Widget the coordinates are relative to its parent.
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

// SetBounds returns the outer bounding box Rectangle of the *WidgetBase, 
// including decorations.
//
// For a RootWidget, like *MainWindow or *Dialog, the Rectangle is in screen
// coordinates, for a child Widget the coordinates are relative to its parent.
func (wb *WidgetBase) SetBounds(bounds Rectangle) os.Error {
	if !MoveWindow(wb.hWnd, bounds.X, bounds.Y, bounds.Width, bounds.Height, true) {
		return lastError("MoveWindow")
	}

	return nil
}

// MinSize returns the minimum allowed outer Size for the *WidgetBase, including
// decorations.
//
// For child widgets, this is only relevant when the parent of the *WidgetBase 
// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
func (wb *WidgetBase) MinSize() Size {
	return wb.minSize
}

// MaxSize returns the maximum allowed outer Size for the *WidgetBase, including
// decorations.
//
// For child widgets, this is only relevant when the parent of the *WidgetBase 
// has a Layout. RootWidgets, like *MainWindow and *Dialog, also honor this.
func (wb *WidgetBase) MaxSize() Size {
	return wb.maxSize
}

// SetMinMaxSize sets the minimum and maximum outer Size of the *WidgetBase,
// including decorations.
//
// Use walk.Size{} to make the respective limit be ignored.
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

// LayoutFlags returns a combination of LayoutFlags that specify how the
// *WidgetBase wants to be treated by Layout implementations.
func (wb *WidgetBase) LayoutFlags() LayoutFlags {
	return 0
}

// MinSizeHint returns the minimum outer Size, including decorations, that 
// makes sense for the respective type of Widget.
func (wb *WidgetBase) MinSizeHint() Size {
	return wb.widget.SizeHint()
}

// SizeHint returns a default Size that should be "overidden" by a concrete
// Widget type.
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

// Size returns the outer Size of the *WidgetBase, including decorations.
func (wb *WidgetBase) Size() Size {
	return wb.Bounds().Size()
}

// SetSize sets the outer Size of the *WidgetBase, including decorations.
func (wb *WidgetBase) SetSize(size Size) os.Error {
	bounds := wb.Bounds()

	return wb.SetBounds(bounds.SetSize(size))
}

// X returns the x coordinate of the *WidgetBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
// child Widgets.
func (wb *WidgetBase) X() int {
	return wb.Bounds().X
}

// SetX sets the x coordinate of the *WidgetBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
// child Widgets.
func (wb *WidgetBase) SetX(value int) os.Error {
	bounds := wb.Bounds()
	bounds.X = value

	return wb.SetBounds(bounds)
}

// Y returns the y coordinate of the *WidgetBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
// child Widgets.
func (wb *WidgetBase) Y() int {
	return wb.Bounds().Y
}

// SetY sets the y coordinate of the *WidgetBase, relative to the screen for
// RootWidgets like *MainWindow or *Dialog and relative to the parent for 
// child Widgets.
func (wb *WidgetBase) SetY(value int) os.Error {
	bounds := wb.Bounds()
	bounds.Y = value

	return wb.SetBounds(bounds)
}

// Width returns the outer width of the *WidgetBase, including decorations.
func (wb *WidgetBase) Width() int {
	return wb.Bounds().Width
}

// SetWidth sets the outer width of the *WidgetBase, including decorations.
func (wb *WidgetBase) SetWidth(value int) os.Error {
	bounds := wb.Bounds()
	bounds.Width = value

	return wb.SetBounds(bounds)
}

// Height returns the outer height of the *WidgetBase, including decorations.
func (wb *WidgetBase) Height() int {
	return wb.Bounds().Height
}

// SetHeight sets the outer height of the *WidgetBase, including decorations.
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

// ClientBounds returns the inner bounding box Rectangle of the *WidgetBase,
// excluding decorations.
func (wb *WidgetBase) ClientBounds() Rectangle {
	return widgetClientBounds(wb.hWnd)
}

func (wb *WidgetBase) sizeFromClientSize(clientSize Size) Size {
	s := wb.Size()
	cs := wb.ClientBounds().Size()
	ncs := Size{s.Width - cs.Width, s.Height - cs.Height}

	return Size{clientSize.Width + ncs.Width, clientSize.Height + ncs.Height}
}

// SetClientSize sets the Size of the inner bounding box of the *WidgetBase,
// excluding decorations.
func (wb *WidgetBase) SetClientSize(value Size) os.Error {
	return wb.SetSize(wb.sizeFromClientSize(value))
}

// SetFocus sets the keyboard input focus to the *WidgetBase.
func (wb *WidgetBase) SetFocus() os.Error {
	if SetFocus(wb.hWnd) == 0 {
		return lastError("SetFocus")
	}

	return nil
}

// CreateCanvas creates and returns a *Canvas that can be used to draw
// inside the ClientBounds of the *WidgetBase.
//
// Remember to call the Dispose method on the canvas to release resources, 
// when you no longer need it. 
func (wb *WidgetBase) CreateCanvas() (*Canvas, os.Error) {
	return newCanvasFromHWND(wb.hWnd)
}

func (wb *WidgetBase) setTheme(appName string) os.Error {
	if hr := SetWindowTheme(wb.hWnd, syscall.StringToUTF16Ptr(appName), nil); FAILED(hr) {
		return errorFromHRESULT("SetWindowTheme", hr)
	}

	return nil
}

// KeyDown returns a *KeyEvent that you can attach to for handling key down
// events for the *WidgetBase.
func (wb *WidgetBase) KeyDown() *KeyEvent {
	return wb.keyDownPublisher.Event()
}

// MouseDown returns a *MouseEvent that you can attach to for handling 
// mouse down events for the *WidgetBase.
func (wb *WidgetBase) MouseDown() *MouseEvent {
	return wb.mouseDownPublisher.Event()
}

// MouseMove returns a *MouseEvent that you can attach to for handling 
// mouse move events for the *WidgetBase.
func (wb *WidgetBase) MouseMove() *MouseEvent {
	return wb.mouseMovePublisher.Event()
}

// MouseUp returns a *MouseEvent that you can attach to for handling 
// mouse up events for the *WidgetBase.
func (wb *WidgetBase) MouseUp() *MouseEvent {
	return wb.mouseUpPublisher.Event()
}

func (wb *WidgetBase) publishMouseEvent(publisher *MouseEventPublisher, wParam, lParam uintptr) {
	x := int(GET_X_LPARAM(lParam))
	y := int(GET_Y_LPARAM(lParam))

	publisher.Publish(x, y, 0)
}

// SizeChanged returns an *Event that you can attach to for handling size
// changed events for the *WidgetBase.
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
		if len(appSingleton.panickingPublisher.event.handlers) > 0 {
			var err os.Error
			if x := recover(); x != nil {
				if e, ok := x.(os.Error); ok {
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
