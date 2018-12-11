// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
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

type Widget interface {
	Window

	// AlwaysConsumeSpace returns if the Widget should consume space even if it
	// is not visible.
	AlwaysConsumeSpace() bool

	// AsWidgetBase returns a *WidgetBase that implements Widget.
	AsWidgetBase() *WidgetBase

	// Form returns the root ancestor Form of the Widget.
	Form() Form

	// GraphicsEffects returns a list of WidgetGraphicsEffects that are applied to the Widget.
	GraphicsEffects() *WidgetGraphicsEffectList

	// LayoutFlags returns a combination of LayoutFlags that specify how the
	// Widget wants to be treated by Layout implementations.
	LayoutFlags() LayoutFlags

	// MinSizeHint returns the minimum outer Size, including decorations, that
	// makes sense for the respective type of Widget.
	MinSizeHint() Size

	// Parent returns the Container of the Widget.
	Parent() Container

	// SetAlwaysConsumeSpace sets if the Widget should consume space even if it
	// is not visible.
	SetAlwaysConsumeSpace(b bool) error

	// SetParent sets the parent of the Widget and adds the Widget to the
	// Children list of the Container.
	SetParent(value Container) error

	// SetToolTipText sets the tool tip text of the Widget.
	SetToolTipText(s string) error

	// SizeHint returns the preferred Size for the respective type of Widget.
	SizeHint() Size

	// ToolTipText returns the tool tip text of the Widget.
	ToolTipText() string
}

type WidgetBase struct {
	WindowBase
	parent                      Container
	toolTipTextProperty         Property
	toolTipTextChangedPublisher EventPublisher
	graphicsEffects             *WidgetGraphicsEffectList
	alwaysConsumeSpace          bool
}

// InitWidget initializes a Widget.
func InitWidget(widget Widget, parent Window, className string, style, exStyle uint32) error {
	if parent == nil {
		return newError("parent cannot be nil")
	}

	if err := InitWindow(widget, parent, className, style|win.WS_CHILD, exStyle); err != nil {
		return err
	}

	if container, ok := parent.(Container); ok {
		if container.Children() == nil {
			// Required by parents like MainWindow and GroupBox.
			if win.SetParent(widget.Handle(), container.Handle()) == 0 {
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

func (wb *WidgetBase) init(widget Widget) error {
	wb.graphicsEffects = newWidgetGraphicsEffectList(wb)

	if err := globalToolTip.AddTool(wb); err != nil {
		return err
	}

	wb.toolTipTextProperty = NewProperty(
		func() interface{} {
			return wb.window.(Widget).ToolTipText()
		},
		func(v interface{}) error {
			wb.window.(Widget).SetToolTipText(assertStringOr(v, ""))
			return nil
		},
		wb.toolTipTextChangedPublisher.Event())

	wb.MustRegisterProperty("ToolTipText", wb.toolTipTextProperty)

	return nil
}

// AsWidgetBase just returns the receiver.
func (wb *WidgetBase) AsWidgetBase() *WidgetBase {
	return wb
}

// Bounds returns the outer bounding box Rectangle of the WidgetBase, including
// decorations.
//
// The coordinates are relative to the parent of the Widget.
func (wb *WidgetBase) Bounds() Rectangle {
	b := wb.WindowBase.Bounds()

	if wb.parent != nil {
		p := win.POINT{int32(b.X), int32(b.Y)}
		if !win.ScreenToClient(wb.parent.Handle(), &p) {
			newError("ScreenToClient failed")
			return Rectangle{}
		}
		b.X = int(p.X)
		b.Y = int(p.Y)
	}

	return b
}

// BringToTop moves the WidgetBase to the top of the keyboard focus order.
func (wb *WidgetBase) BringToTop() error {
	if wb.parent != nil {
		if err := wb.parent.BringToTop(); err != nil {
			return err
		}
	}

	return wb.WindowBase.BringToTop()
}

// Enabled returns if the WidgetBase is enabled for user interaction.
func (wb *WidgetBase) Enabled() bool {
	if wb.parent != nil {
		return wb.enabled && wb.parent.Enabled()
	}

	return wb.enabled
}

// Font returns the Font of the WidgetBase.
//
// By default this is a MS Shell Dlg 2, 8 point font.
func (wb *WidgetBase) Font() *Font {
	if wb.font != nil {
		return wb.font
	} else if wb.parent != nil {
		return wb.parent.Font()
	}

	return defaultFont
}

func (wb *WidgetBase) applyFont(font *Font) {
	wb.WindowBase.applyFont(font)

	wb.updateParentLayout()
}

// Form returns the root ancestor Form of the Widget.
func (wb *WidgetBase) Form() Form {
	return ancestor(wb)
}

// LayoutFlags returns a combination of LayoutFlags that specify how the
// WidgetBase wants to be treated by Layout implementations.
func (wb *WidgetBase) LayoutFlags() LayoutFlags {
	return 0
}

// SetMinMaxSize sets the minimum and maximum outer Size of the *WidgetBase,
// including decorations.
//
// Use walk.Size{} to make the respective limit be ignored.
func (wb *WidgetBase) SetMinMaxSize(min, max Size) (err error) {
	err = wb.WindowBase.SetMinMaxSize(min, max)

	wb.updateParentLayout()

	return
}

// AlwaysConsumeSpace returns if the Widget should consume space even if it is
// not visible.
func (wb *WidgetBase) AlwaysConsumeSpace() bool {
	return wb.alwaysConsumeSpace
}

// SetAlwaysConsumeSpace sets if the Widget should consume space even if it is
// not visible.
func (wb *WidgetBase) SetAlwaysConsumeSpace(b bool) error {
	wb.alwaysConsumeSpace = b

	return wb.updateParentLayout()
}

// MinSizeHint returns the minimum outer Size, including decorations, that
// makes sense for the respective type of Widget.
func (wb *WidgetBase) MinSizeHint() Size {
	return Size{10, 10}
}

// Parent returns the Container of the WidgetBase.
func (wb *WidgetBase) Parent() Container {
	return wb.parent
}

// SetParent sets the parent of the WidgetBase and adds the WidgetBase to the
// Children list of the Container.
func (wb *WidgetBase) SetParent(parent Container) (err error) {
	if parent == wb.parent {
		return nil
	}

	style := uint32(win.GetWindowLong(wb.hWnd, win.GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	if parent == nil {
		style &^= win.WS_CHILD
		style |= win.WS_POPUP

		if win.SetParent(wb.hWnd, 0) == 0 {
			return lastError("SetParent")
		}
		win.SetLastError(0)
		if win.SetWindowLong(wb.hWnd, win.GWL_STYLE, int32(style)) == 0 {
			return lastError("SetWindowLong")
		}
	} else {
		style |= win.WS_CHILD
		style &^= win.WS_POPUP

		win.SetLastError(0)
		if win.SetWindowLong(wb.hWnd, win.GWL_STYLE, int32(style)) == 0 {
			return lastError("SetWindowLong")
		}
		if win.SetParent(wb.hWnd, parent.Handle()) == 0 {
			return lastError("SetParent")
		}
	}

	b := wb.Bounds()

	if !win.SetWindowPos(
		wb.hWnd,
		win.HWND_BOTTOM,
		int32(b.X),
		int32(b.Y),
		int32(b.Width),
		int32(b.Height),
		win.SWP_FRAMECHANGED) {

		return lastError("SetWindowPos")
	}

	oldParent := wb.parent

	wb.parent = parent

	var oldChildren, newChildren *WidgetList
	if oldParent != nil {
		oldChildren = oldParent.Children()
	}
	if parent != nil {
		newChildren = parent.Children()
	}

	if newChildren == oldChildren {
		return nil
	}

	widget := wb.window.(Widget)

	if oldChildren != nil {
		oldChildren.Remove(widget)
	}

	if newChildren != nil && !newChildren.containsHandle(wb.hWnd) {
		newChildren.Add(widget)
	}

	return nil
}

func (wb *WidgetBase) ForEachAncestor(f func(window Window) bool) {
	hwnd := win.GetParent(wb.hWnd)

	for hwnd != 0 {
		if window := windowFromHandle(hwnd); window != nil {
			if !f(window) {
				return
			}
		}

		hwnd = win.GetParent(hwnd)
	}
}

// SizeHint returns a default Size that should be "overidden" by a concrete
// Widget type.
func (wb *WidgetBase) SizeHint() Size {
	return wb.window.(Widget).MinSizeHint()
}

// ToolTipText returns the tool tip text of the WidgetBase.
func (wb *WidgetBase) ToolTipText() string {
	return globalToolTip.Text(wb.window.(Widget))
}

// SetToolTipText sets the tool tip text of the WidgetBase.
func (wb *WidgetBase) SetToolTipText(s string) error {
	if err := globalToolTip.SetText(wb.window.(Widget), s); err != nil {
		return err
	}

	wb.toolTipTextChangedPublisher.Publish()

	return nil
}

// GraphicsEffects returns a list of WidgetGraphicsEffects that are applied to the WidgetBase.
func (wb *WidgetBase) GraphicsEffects() *WidgetGraphicsEffectList {
	return wb.graphicsEffects
}

func (wb *WidgetBase) onInsertedGraphicsEffect(index int, effect WidgetGraphicsEffect) error {
	wb.invalidateBorderInParent()

	return nil
}

func (wb *WidgetBase) onRemovedGraphicsEffect(index int, effect WidgetGraphicsEffect) error {
	wb.invalidateBorderInParent()

	return nil
}

func (wb *WidgetBase) onClearedGraphicsEffects() error {
	wb.invalidateBorderInParent()

	return nil
}

func (wb *WidgetBase) invalidateBorderInParent() {
	if wb.parent != nil && wb.parent.Layout() != nil {
		b := wb.Bounds().toRECT()
		s := int32(wb.parent.Layout().Spacing())

		hwnd := wb.parent.Handle()

		rc := win.RECT{Left: b.Left - s, Top: b.Top - s, Right: b.Left, Bottom: b.Bottom + s}
		win.InvalidateRect(hwnd, &rc, true)

		rc = win.RECT{Left: b.Right, Top: b.Top - s, Right: b.Right + s, Bottom: b.Bottom + s}
		win.InvalidateRect(hwnd, &rc, true)

		rc = win.RECT{Left: b.Left, Top: b.Top - s, Right: b.Right, Bottom: b.Top}
		win.InvalidateRect(hwnd, &rc, true)

		rc = win.RECT{Left: b.Left, Top: b.Bottom, Right: b.Right, Bottom: b.Bottom + s}
		win.InvalidateRect(hwnd, &rc, true)
	}
}

func (wb *WidgetBase) hasComplexBackground() bool {
	if bg := wb.window.Background(); bg != nil && !bg.simple() {
		return false
	}

	var complex bool
	wb.ForEachAncestor(func(window Window) bool {
		if bg := window.Background(); bg != nil && !bg.simple() {
			complex = true
			return false
		}

		return true
	})

	return complex
}

func (wb *WidgetBase) updateParentLayout() error {
	parent := wb.window.(Widget).Parent()

	if parent == nil || parent.Layout() == nil || parent.Suspended() || !parent.Visible() {
		return nil
	}

	layout := parent.Layout()

	if !formResizeScheduled || len(inProgressEventsByForm[appSingleton.activeForm]) == 0 {
		clientSize := parent.ClientBounds().Size()
		minSize := layout.MinSize()

		if clientSize.Width < minSize.Width || clientSize.Height < minSize.Height {
			switch wnd := parent.(type) {
			case *ScrollView:
				ifContainerIsScrollViewDoCoolSpecialLayoutStuff(layout)
				return nil

			case Widget:
				return wnd.AsWidgetBase().updateParentLayout()

			case Form:
				if len(inProgressEventsByForm[appSingleton.activeForm]) > 0 {
					formResizeScheduled = true
				} else {
					bounds := wnd.Bounds()

					if wnd.AsFormBase().fixedSize() {
						bounds.Width, bounds.Height = 0, 0
					}

					wnd.SetBounds(bounds)

					return nil
				}
			}
		}
	}

	layout.Update(false)

	if FocusEffect != nil {
		if focusedWnd := windowFromHandle(win.GetFocus()); focusedWnd != nil && win.GetParent(focusedWnd.Handle()) == parent.Handle() {
			focusedWnd.(Widget).AsWidgetBase().invalidateBorderInParent()
		}
	}

	return nil
}

func ancestor(w Widget) Form {
	if w == nil {
		return nil
	}

	hWndRoot := win.GetAncestor(w.Handle(), win.GA_ROOT)

	rw, _ := windowFromHandle(hWndRoot).(Form)
	return rw
}

func minSizeEffective(w Widget) Size {
	s := maxSize(w.MinSize(), w.MinSizeHint())

	max := w.MaxSize()
	if max.Width > 0 && s.Width > max.Width {
		s.Width = max.Width
	}
	if max.Height > 0 && s.Height > max.Height {
		s.Height = max.Height
	}

	return s
}
