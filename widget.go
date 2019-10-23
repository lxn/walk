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

	// Alignment returns the alignment of the Widget.
	Alignment() Alignment2D

	// AlwaysConsumeSpace returns if the Widget should consume space even if it
	// is not visible.
	AlwaysConsumeSpace() bool

	// AsWidgetBase returns a *WidgetBase that implements Widget.
	AsWidgetBase() *WidgetBase

	// CreateLayoutItem creates and returns a new LayoutItem specific to the
	// concrete Widget type, that carries all data and logic required to layout
	// the Widget.
	CreateLayoutItem(ctx *LayoutContext) LayoutItem

	// GraphicsEffects returns a list of WidgetGraphicsEffects that are applied to the Widget.
	GraphicsEffects() *WidgetGraphicsEffectList

	// LayoutFlags returns a combination of LayoutFlags that specify how the
	// Widget wants to be treated by Layout implementations.
	LayoutFlags() LayoutFlags

	// MinSizeHint returns the minimum outer size in native pixels, including decorations, that
	// makes sense for the respective type of Widget.
	MinSizeHint() Size

	// Parent returns the Container of the Widget.
	Parent() Container

	// SetAlignment sets the alignment of the widget.
	SetAlignment(alignment Alignment2D) error

	// SetAlwaysConsumeSpace sets if the Widget should consume space even if it
	// is not visible.
	SetAlwaysConsumeSpace(b bool) error

	// SetParent sets the parent of the Widget and adds the Widget to the
	// Children list of the Container.
	SetParent(value Container) error

	// SetToolTipText sets the tool tip text of the Widget.
	SetToolTipText(s string) error

	// SizeHint returns the preferred size in native pixels for the respective type of Widget.
	SizeHint() Size

	// ToolTipText returns the tool tip text of the Widget.
	ToolTipText() string
}

type WidgetBase struct {
	WindowBase
	geometry                    Geometry
	parent                      Container
	toolTipTextProperty         Property
	toolTipTextChangedPublisher EventPublisher
	graphicsEffects             *WidgetGraphicsEffectList
	alignment                   Alignment2D
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

	tt, err := wb.group.CreateToolTip()
	if err != nil {
		return err
	}
	if err := tt.AddTool(wb.window.(Widget)); err != nil {
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

func (wb *WidgetBase) Dispose() {
	if wb.hWnd == 0 {
		return
	}

	if wb.parent != nil && win.GetParent(wb.hWnd) == wb.parent.Handle() {
		wb.SetParent(nil)
	}

	if tt := wb.group.ToolTip(); tt != nil {
		tt.RemoveTool(wb.window.(Widget))
	}

	wb.WindowBase.Dispose()
}

// AsWidgetBase just returns the receiver.
func (wb *WidgetBase) AsWidgetBase() *WidgetBase {
	return wb
}

// Bounds returns the outer bounding box rectangle of the WidgetBase, including
// decorations.
//
// The coordinates are relative to the parent of the Widget.
func (wb *WidgetBase) Bounds() Rectangle {
	return wb.RectangleTo96DPI(wb.BoundsPixels())
}

// BoundsPixels returns the outer bounding box rectangle of the WidgetBase, including
// decorations.
//
// The coordinates are relative to the parent of the Widget.
func (wb *WidgetBase) BoundsPixels() Rectangle {
	b := wb.WindowBase.BoundsPixels()

	if wb.parent != nil {
		p := b.Location().toPOINT()
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

	wb.RequestLayout()
}

// Alignment return the alignment ot the *WidgetBase.
func (wb *WidgetBase) Alignment() Alignment2D {
	return wb.alignment
}

// SetAlignment sets the alignment of the *WidgetBase.
func (wb *WidgetBase) SetAlignment(alignment Alignment2D) error {
	if alignment != wb.alignment {
		if alignment < AlignHVDefault || alignment > AlignHFarVFar {
			return newError("invalid Alignment value")
		}

		wb.alignment = alignment

		wb.RequestLayout()
	}

	return nil
}

// SetMinMaxSize sets the minimum and maximum outer size of the *WidgetBase,
// including decorations.
//
// Use walk.Size{} to make the respective limit be ignored.
func (wb *WidgetBase) SetMinMaxSize(min, max Size) (err error) {
	err = wb.WindowBase.SetMinMaxSize(min, max)

	wb.RequestLayout()

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

	wb.RequestLayout()

	return nil
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
		wb.SetVisible(false)

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

		if cb := parent.AsContainerBase(); cb != nil {
			win.SetWindowLong(wb.hWnd, win.GWL_ID, cb.NextChildID())
		}
	}

	b := wb.BoundsPixels()

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

// ToolTipText returns the tool tip text of the WidgetBase.
func (wb *WidgetBase) ToolTipText() string {
	if tt := wb.group.ToolTip(); tt != nil {
		return tt.Text(wb.window.(Widget))
	}
	return ""
}

// SetToolTipText sets the tool tip text of the WidgetBase.
func (wb *WidgetBase) SetToolTipText(s string) error {
	if tt := wb.group.ToolTip(); tt != nil {
		if err := tt.SetText(wb.window.(Widget), s); err != nil {
			return err
		}
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
	if !wb.hasActiveGraphicsEffects() {
		return
	}

	if wb.parent != nil && wb.parent.Layout() != nil {
		b := wb.BoundsPixels().toRECT()
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

func (wb *WidgetBase) hasActiveGraphicsEffects() bool {
	if wb.graphicsEffects == nil {
		return false
	}

	count := wb.graphicsEffects.Len()

	for _, gfx := range [...]WidgetGraphicsEffect{FocusEffect, InteractionEffect, ValidationErrorEffect} {
		if wb.graphicsEffects.Contains(gfx) {
			if gfx == nil {
				count--
			}
		}
	}

	return count > 0
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

func ancestor(w Widget) Form {
	if w == nil {
		return nil
	}

	hWndRoot := win.GetAncestor(w.Handle(), win.GA_ROOT)

	rw, _ := windowFromHandle(hWndRoot).(Form)
	return rw
}

func (wb *WidgetBase) LayoutFlags() LayoutFlags {
	return createLayoutItemForWidget(wb.window.(Widget)).LayoutFlags()
}

func (wb *WidgetBase) SizeHint() Size {
	if is, ok := createLayoutItemForWidget(wb.window.(Widget)).(IdealSizer); ok {
		return is.IdealSize()
	}

	return Size{}
}

func (wb *WidgetBase) MinSizeHint() Size {
	if ms, ok := createLayoutItemForWidget(wb.window.(Widget)).(MinSizer); ok {
		return ms.MinSize()
	}

	return Size{}
}
