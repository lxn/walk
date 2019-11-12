// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/win"
)

type Container interface {
	Window
	AsContainerBase() *ContainerBase
	Children() *WidgetList
	Layout() Layout
	SetLayout(value Layout) error
	DataBinder() *DataBinder
	SetDataBinder(dbm *DataBinder)
}

type ContainerBase struct {
	WidgetBase
	layout      Layout
	children    *WidgetList
	dataBinder  *DataBinder
	nextChildID int32
	persistent  bool
}

func (cb *ContainerBase) AsWidgetBase() *WidgetBase {
	return &cb.WidgetBase
}

func (cb *ContainerBase) AsContainerBase() *ContainerBase {
	return cb
}

func (cb *ContainerBase) NextChildID() int32 {
	cb.nextChildID++
	return cb.nextChildID
}

func (cb *ContainerBase) applyEnabled(enabled bool) {
	cb.WidgetBase.applyEnabled(enabled)

	applyEnabledToDescendants(cb.window.(Widget), enabled)
}

func (cb *ContainerBase) applyFont(font *Font) {
	cb.WidgetBase.applyFont(font)

	applyFontToDescendants(cb.window.(Widget), font)
}

func (cb *ContainerBase) ApplySysColors() {
	cb.WidgetBase.ApplySysColors()

	applySysColorsToDescendants(cb.window.(Widget))
}

func (cb *ContainerBase) ApplyDPI(dpi int) {
	cb.WidgetBase.ApplyDPI(dpi)

	applyDPIToDescendants(cb.window.(Widget), dpi)

	if cb.layout != nil {
		if ums, ok := cb.layout.(interface {
			updateMargins()
			updateSpacing()
		}); ok {
			ums.updateMargins()
			ums.updateSpacing()
		}

		cb.RequestLayout()
	}
}

func (cb *ContainerBase) Children() *WidgetList {
	return cb.children
}

func (cb *ContainerBase) Layout() Layout {
	return cb.layout
}

func (cb *ContainerBase) SetLayout(value Layout) error {
	if cb.layout != value {
		if cb.layout != nil {
			cb.layout.SetContainer(nil)
		}

		cb.layout = value

		if value != nil && value.Container() != Container(cb) {
			value.SetContainer(cb)
		}
	}

	return nil
}

func (cb *ContainerBase) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	return cb.layout.CreateLayoutItem(ctx)
}

func (cb *ContainerBase) DataBinder() *DataBinder {
	return cb.dataBinder
}

func (cb *ContainerBase) SetDataBinder(db *DataBinder) {
	if db == cb.dataBinder {
		return
	}

	if cb.dataBinder != nil {
		cb.dataBinder.SetBoundWidgets(nil)
	}

	cb.dataBinder = db

	if db != nil {
		var boundWidgets []Widget

		walkDescendants(cb.window, func(w Window) bool {
			if w.Handle() == cb.hWnd {
				return true
			}

			if c, ok := w.(Container); ok && c.DataBinder() != nil {
				return false
			}

			for _, prop := range w.AsWindowBase().name2Property {
				if _, ok := prop.Source().(string); ok {
					boundWidgets = append(boundWidgets, w.(Widget))
					break
				}
			}

			return true
		})

		db.SetBoundWidgets(boundWidgets)
	}
}

func (cb *ContainerBase) forEachPersistableChild(f func(p Persistable) error) error {
	if cb.children == nil {
		return nil
	}

	for _, wb := range cb.children.items {
		if persistable, ok := wb.window.(Persistable); ok && persistable.Persistent() {
			if err := f(persistable); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cb *ContainerBase) Persistent() bool {
	return cb.persistent
}

func (cb *ContainerBase) SetPersistent(value bool) {
	cb.persistent = value
}

func (cb *ContainerBase) SaveState() error {
	return cb.forEachPersistableChild(func(p Persistable) error {
		return p.SaveState()
	})
}

func (cb *ContainerBase) RestoreState() error {
	return cb.forEachPersistableChild(func(p Persistable) error {
		return p.RestoreState()
	})
}

func (cb *ContainerBase) doPaint() error {
	var ps win.PAINTSTRUCT

	hdc := win.BeginPaint(cb.hWnd, &ps)
	defer win.EndPaint(cb.hWnd, &ps)

	canvas, err := newCanvasFromHDC(hdc)
	if err != nil {
		return err
	}
	defer canvas.Dispose()

	for _, wb := range cb.children.items {
		widget := wb.window.(Widget)

		for _, effect := range widget.GraphicsEffects().items {
			switch effect {
			case InteractionEffect:
				type ReadOnlyer interface {
					ReadOnly() bool
				}
				if ro, ok := widget.(ReadOnlyer); ok {
					if ro.ReadOnly() {
						continue
					}
				}

				if hwnd := widget.Handle(); !win.IsWindowEnabled(hwnd) || !win.IsWindowVisible(hwnd) {
					continue
				}

			case FocusEffect:
				continue
			}

			b := widget.BoundsPixels().toRECT()
			win.ExcludeClipRect(hdc, b.Left, b.Top, b.Right, b.Bottom)

			if err := effect.Draw(widget, canvas); err != nil {
				return err
			}
		}
	}

	if FocusEffect != nil {
		hwndFocused := win.GetFocus()
		var widget Widget
		if wnd := windowFromHandle(hwndFocused); wnd != nil {
			widget, _ = wnd.(Widget)
		}
		for hwndFocused != 0 && (widget == nil || widget.Parent() == nil) {
			hwndFocused = win.GetParent(hwndFocused)
			if wnd := windowFromHandle(hwndFocused); wnd != nil {
				widget, _ = wnd.(Widget)
			}
		}

		if widget != nil && widget.Parent() != nil && widget.Parent().Handle() == cb.hWnd {
			for _, effect := range widget.GraphicsEffects().items {
				if effect == FocusEffect {
					b := widget.BoundsPixels().toRECT()
					win.ExcludeClipRect(hdc, b.Left, b.Top, b.Right, b.Bottom)

					if err := FocusEffect.Draw(widget, canvas); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (cb *ContainerBase) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_CTLCOLOREDIT, win.WM_CTLCOLORSTATIC:
		if hBrush := cb.handleWMCTLCOLOR(wParam, lParam); hBrush != 0 {
			return hBrush
		}

	case win.WM_PAINT:
		if FocusEffect == nil && InteractionEffect == nil && ValidationErrorEffect == nil {
			break
		}

		// If it fails, what can we do about it? Panic? That's extreme. So just ignore it.
		_ = cb.doPaint()

		return 0

	case win.WM_COMMAND:
		if lParam == 0 {
			switch win.HIWORD(uint32(wParam)) {
			case 0:
				cmdId := win.LOWORD(uint32(wParam))
				switch cmdId {
				case win.IDOK, win.IDCANCEL:
					form := ancestor(cb)
					if form == nil {
						break
					}

					dlg, ok := form.(dialogish)
					if !ok {
						break
					}

					var button *PushButton
					if cmdId == win.IDOK {
						button = dlg.DefaultButton()
					} else {
						button = dlg.CancelButton()
					}

					if button != nil && button.Visible() && button.Enabled() {
						button.raiseClicked()
					}

					break
				}

				// Menu
				actionId := uint16(win.LOWORD(uint32(wParam)))
				if action, ok := actionsById[actionId]; ok {
					action.raiseTriggered()
					return 0
				}

			case 1:
				// Accelerator
			}
		} else {
			// The window that sent the notification shall handle it itself.
			hwndSrc := win.GetDlgItem(cb.hWnd, int32(win.LOWORD(uint32(wParam))))

			var toolBarOnly bool
			if hwndSrc == 0 {
				toolBarOnly = true
				hwndSrc = win.HWND(lParam)
			}

			if window := windowFromHandle(hwndSrc); window != nil {
				if _, ok := window.(*ToolBar); toolBarOnly && !ok {
					break
				}

				window.WndProc(hwnd, msg, wParam, lParam)
				return 0
			}
		}

	case win.WM_MEASUREITEM:
		mis := (*win.MEASUREITEMSTRUCT)(unsafe.Pointer(lParam))
		if window := windowFromHandle(win.GetDlgItem(hwnd, int32(mis.CtlID))); window != nil {
			// The window that sent the notification shall handle it itself.
			return window.WndProc(hwnd, msg, wParam, lParam)
		}

	case win.WM_DRAWITEM:
		dis := (*win.DRAWITEMSTRUCT)(unsafe.Pointer(lParam))
		if window := windowFromHandle(dis.HwndItem); window != nil {
			// The window that sent the notification shall handle it itself.
			return window.WndProc(hwnd, msg, wParam, lParam)
		}

	case win.WM_NOTIFY:
		nmh := (*win.NMHDR)(unsafe.Pointer(lParam))
		if window := windowFromHandle(nmh.HwndFrom); window != nil {
			// The window that sent the notification shall handle it itself.
			return window.WndProc(hwnd, msg, wParam, lParam)
		}

	case win.WM_HSCROLL, win.WM_VSCROLL:
		if window := windowFromHandle(win.HWND(lParam)); window != nil {
			// The window that sent the notification shall handle it itself.
			return window.WndProc(hwnd, msg, wParam, lParam)
		}

	case win.WM_WINDOWPOSCHANGED:
		wp := (*win.WINDOWPOS)(unsafe.Pointer(lParam))

		if wp.Flags&win.SWP_NOSIZE != 0 || cb.Layout() == nil {
			break
		}

		if cb.background == nullBrushSingleton {
			cb.Invalidate()
		}
	}

	return cb.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

func (cb *ContainerBase) onInsertingWidget(index int, widget Widget) (err error) {
	return nil
}

func (cb *ContainerBase) onInsertedWidget(index int, widget Widget) (err error) {
	if parent := widget.Parent(); parent == nil || parent.Handle() != cb.hWnd {
		if err = widget.SetParent(cb.window.(Container)); err != nil {
			return
		}
	}

	cb.RequestLayout()

	widget.(applyFonter).applyFont(cb.Font())

	return
}

func (cb *ContainerBase) onRemovingWidget(index int, widget Widget) (err error) {
	if widget.Parent() == nil {
		return
	}

	if widget.Parent().Handle() == cb.hWnd {
		err = widget.SetParent(nil)
	}

	return
}

func (cb *ContainerBase) onRemovedWidget(index int, widget Widget) (err error) {
	cb.RequestLayout()

	return
}

func (cb *ContainerBase) onClearingWidgets() (err error) {
	for i := cb.children.Len() - 1; i >= 0; i-- {
		widget := cb.children.At(i)

		if parent := widget.Parent(); parent != nil && parent.Handle() == cb.hWnd {
			if err = widget.SetParent(nil); err != nil {
				return
			}
		}
	}

	return
}

func (cb *ContainerBase) onClearedWidgets() (err error) {
	cb.RequestLayout()

	return
}

func (cb *ContainerBase) focusFirstCandidateDescendant() {
	window := firstFocusableDescendant(cb)
	if window == nil {
		return
	}

	if err := window.SetFocus(); err != nil {
		return
	}

	if textSel, ok := window.(textSelectable); ok {
		time.AfterFunc(time.Millisecond, func() {
			window.Synchronize(func() {
				if window.Focused() {
					textSel.SetTextSelection(0, -1)
				}
			})
		})
	}
}

func firstFocusableDescendantCallback(hwnd win.HWND, lParam uintptr) uintptr {
	if !win.IsWindowVisible(hwnd) || !win.IsWindowEnabled(hwnd) {
		return 1
	}

	if win.GetWindowLong(hwnd, win.GWL_STYLE)&win.WS_TABSTOP > 0 {
		if rb, ok := windowFromHandle(hwnd).(radioButtonish); ok {
			if !rb.radioButton().Checked() {
				return 1
			}
		}

		hwndPtr := (*win.HWND)(unsafe.Pointer(lParam))
		*hwndPtr = hwnd
		return 0
	}

	return 1
}

var firstFocusableDescendantCallbackPtr uintptr

func init() {
	AppendToWalkInit(func() {
		firstFocusableDescendantCallbackPtr = syscall.NewCallback(firstFocusableDescendantCallback)
	})
}

func firstFocusableDescendant(container Container) Window {
	var hwnd win.HWND

	win.EnumChildWindows(container.Handle(), firstFocusableDescendantCallbackPtr, uintptr(unsafe.Pointer(&hwnd)))

	window := windowFromHandle(hwnd)

	for hwnd != 0 && window == nil {
		hwnd = win.GetParent(hwnd)
		window = windowFromHandle(hwnd)
	}

	return window
}

type textSelectable interface {
	SetTextSelection(start, end int)
}

func DescendantByName(container Container, name string) Widget {
	var widget Widget

	walkDescendants(container.AsContainerBase(), func(w Window) bool {
		if w.Name() == name {
			widget = w.(Widget)
			return false
		}

		return true
	})

	if widget == nil {
		return nil
	}

	return widget
}
