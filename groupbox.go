// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const groupBoxWindowClass = `\o/ Walk_GroupBox_Class \o/`

func init() {
	AppendToWalkInit(func() {
		MustRegisterWindowClass(groupBoxWindowClass)
	})
}

type GroupBox struct {
	WidgetBase
	hWndGroupBox          win.HWND
	checkBox              *CheckBox
	composite             *Composite
	headerHeight          int // in native pixels
	titleChangedPublisher EventPublisher
}

func NewGroupBox(parent Container) (*GroupBox, error) {
	gb := new(GroupBox)

	if err := InitWidget(
		gb,
		parent,
		groupBoxWindowClass,
		win.WS_VISIBLE,
		win.WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			gb.Dispose()
		}
	}()

	gb.hWndGroupBox = win.CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		win.WS_CHILD|win.WS_VISIBLE|win.BS_GROUPBOX,
		0, 0, 80, 24, gb.hWnd, 0, 0, nil)
	if gb.hWndGroupBox == 0 {
		return nil, lastError("CreateWindowEx(BUTTON)")
	}
	win.SetWindowLong(gb.hWndGroupBox, win.GWL_ID, 1)

	gb.applyFont(gb.Font())
	gb.updateHeaderHeight()

	var err error

	gb.checkBox, err = NewCheckBox(gb)
	if err != nil {
		return nil, err
	}
	win.SetWindowLong(gb.checkBox.hWnd, win.GWL_ID, 2)

	gb.SetCheckable(false)
	gb.checkBox.SetChecked(true)

	gb.checkBox.CheckedChanged().Attach(func() {
		gb.applyEnabledFromCheckBox(gb.checkBox.Checked())
	})

	setWindowVisible(gb.checkBox.hWnd, false)

	gb.composite, err = NewComposite(gb)
	if err != nil {
		return nil, err
	}
	win.SetWindowLong(gb.composite.hWnd, win.GWL_ID, 3)
	gb.composite.name = "composite"

	win.SetWindowPos(gb.checkBox.hWnd, win.HWND_TOP, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)

	gb.SetBackground(NullBrush())

	gb.MustRegisterProperty("Title", NewProperty(
		func() interface{} {
			return gb.Title()
		},
		func(v interface{}) error {
			return gb.SetTitle(assertStringOr(v, ""))
		},
		gb.titleChangedPublisher.Event()))

	gb.MustRegisterProperty("Checked", NewBoolProperty(
		func() bool {
			return gb.Checked()
		},
		func(v bool) error {
			gb.SetChecked(v)
			return nil
		},
		gb.CheckedChanged()))

	succeeded = true

	return gb, nil
}

func (gb *GroupBox) AsContainerBase() *ContainerBase {
	if gb.composite == nil {
		return nil
	}

	return gb.composite.AsContainerBase()
}

func (gb *GroupBox) ClientBoundsPixels() Rectangle {
	cb := windowClientBounds(gb.hWndGroupBox)

	if gb.Layout() == nil {
		return cb
	}

	if gb.Checkable() {
		s := createLayoutItemForWidget(gb.checkBox).(MinSizer).MinSize()

		cb.Y += s.Height
		cb.Height -= s.Height
	}

	padding := gb.IntFrom96DPI(1)
	return Rectangle{cb.X + padding, cb.Y + gb.headerHeight, cb.Width - 2*padding, cb.Height - gb.headerHeight - 2*padding}
}

func (gb *GroupBox) updateHeaderHeight() {
	gb.headerHeight = gb.calculateTextSizeImpl("gM").Height
}

func (gb *GroupBox) Persistent() bool {
	return gb.composite.Persistent()
}

func (gb *GroupBox) SetPersistent(value bool) {
	gb.composite.SetPersistent(value)
}

func (gb *GroupBox) SaveState() error {
	return gb.composite.SaveState()
}

func (gb *GroupBox) RestoreState() error {
	return gb.composite.RestoreState()
}

func (gb *GroupBox) applyEnabled(enabled bool) {
	gb.WidgetBase.applyEnabled(enabled)

	if gb.hWndGroupBox != 0 {
		setWindowEnabled(gb.hWndGroupBox, enabled)
	}

	if gb.checkBox != nil {
		gb.checkBox.applyEnabled(enabled)
	}

	if gb.composite != nil {
		gb.composite.applyEnabled(enabled)
	}
}

func (gb *GroupBox) applyEnabledFromCheckBox(enabled bool) {
	if gb.hWndGroupBox != 0 {
		setWindowEnabled(gb.hWndGroupBox, enabled)
	}

	if gb.composite != nil {
		gb.composite.applyEnabled(enabled)
	}
}

func (gb *GroupBox) applyFont(font *Font) {
	gb.WidgetBase.applyFont(font)

	if gb.checkBox != nil {
		gb.checkBox.applyFont(font)
	}

	if gb.hWndGroupBox != 0 {
		SetWindowFont(gb.hWndGroupBox, font)
	}

	if gb.composite != nil {
		gb.composite.applyFont(font)
	}

	gb.updateHeaderHeight()
}

func (gb *GroupBox) SetSuspended(suspend bool) {
	gb.composite.SetSuspended(suspend)
	gb.WidgetBase.SetSuspended(suspend)
	gb.Invalidate()
}

func (gb *GroupBox) DataBinder() *DataBinder {
	return gb.composite.dataBinder
}

func (gb *GroupBox) SetDataBinder(dataBinder *DataBinder) {
	gb.composite.SetDataBinder(dataBinder)
}

func (gb *GroupBox) Title() string {
	if gb.Checkable() {
		return gb.checkBox.Text()
	}

	return windowText(gb.hWndGroupBox)
}

func (gb *GroupBox) SetTitle(title string) error {
	if gb.Checkable() {
		if err := setWindowText(gb.hWndGroupBox, ""); err != nil {
			return err
		}

		return gb.checkBox.SetText(title)
	}

	return setWindowText(gb.hWndGroupBox, title)
}

func (gb *GroupBox) Checkable() bool {
	return gb.checkBox.visible
}

func (gb *GroupBox) SetCheckable(checkable bool) {
	title := gb.Title()

	gb.checkBox.SetVisible(checkable)

	gb.SetTitle(title)

	gb.RequestLayout()
}

func (gb *GroupBox) Checked() bool {
	return gb.checkBox.Checked()
}

func (gb *GroupBox) SetChecked(checked bool) {
	gb.checkBox.SetChecked(checked)
}

func (gb *GroupBox) CheckedChanged() *Event {
	return gb.checkBox.CheckedChanged()
}

func (gb *GroupBox) ApplyDPI(dpi int) {
	gb.WidgetBase.ApplyDPI(dpi)
	if gb.checkBox != nil {
		gb.checkBox.ApplyDPI(dpi)
	}
	if gb.composite != nil {
		gb.composite.ApplyDPI(dpi)
	}
}

func (gb *GroupBox) Children() *WidgetList {
	if gb.composite == nil {
		// Without this we would get into trouble in NewComposite.
		return nil
	}

	return gb.composite.Children()
}

func (gb *GroupBox) Layout() Layout {
	if gb.composite == nil {
		// Without this we would get into trouble through the call to
		// SetCheckable in NewGroupBox.
		return nil
	}

	return gb.composite.Layout()
}

func (gb *GroupBox) SetLayout(value Layout) error {
	return gb.composite.SetLayout(value)
}

func (gb *GroupBox) MouseDown() *MouseEvent {
	return gb.composite.MouseDown()
}

func (gb *GroupBox) MouseMove() *MouseEvent {
	return gb.composite.MouseMove()
}

func (gb *GroupBox) MouseUp() *MouseEvent {
	return gb.composite.MouseUp()
}

func (gb *GroupBox) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if gb.composite != nil {
		switch msg {
		case win.WM_CTLCOLORSTATIC:
			if hBrush := gb.handleWMCTLCOLOR(wParam, lParam); hBrush != 0 {
				return hBrush
			}

		case win.WM_COMMAND:
			hwndSrc := win.GetDlgItem(gb.hWnd, int32(win.LOWORD(uint32(wParam))))

			if window := windowFromHandle(hwndSrc); window != nil {
				window.WndProc(hwnd, msg, wParam, lParam)
			}

		case win.WM_NOTIFY:
			gb.composite.WndProc(hwnd, msg, wParam, lParam)

		case win.WM_SETTEXT:
			gb.titleChangedPublisher.Publish()

		case win.WM_PAINT:
			win.UpdateWindow(gb.checkBox.hWnd)

		case win.WM_WINDOWPOSCHANGED:
			wp := (*win.WINDOWPOS)(unsafe.Pointer(lParam))

			if wp.Flags&win.SWP_NOSIZE != 0 {
				break
			}

			offset := gb.headerHeight / 4
			wbcb := gb.WidgetBase.ClientBoundsPixels()
			if !win.MoveWindow(
				gb.hWndGroupBox,
				int32(wbcb.X),
				int32(wbcb.Y-offset),
				int32(wbcb.Width),
				int32(wbcb.Height),
				true) {

				lastError("MoveWindow")
				break
			}

			if gb.Checkable() {
				s := createLayoutItemForWidget(gb.checkBox).(MinSizer).MinSize()
				var x int
				if l := gb.Layout(); l != nil {
					x = gb.IntFrom96DPI(l.Margins().HNear)
				} else {
					x = gb.headerHeight * 2 / 3
				}
				gb.checkBox.SetBoundsPixels(Rectangle{x, gb.headerHeight, s.Width, s.Height})
			}

			gbcb := gb.ClientBoundsPixels()
			gbcb.Y -= offset
			gb.composite.SetBoundsPixels(gbcb)
		}
	}

	return gb.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

func (gb *GroupBox) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	compositePos := Point{gb.IntFrom96DPI(1), gb.headerHeight}
	if gb.Checkable() {
		idealSize := gb.checkBox.idealSize()

		compositePos.Y += idealSize.Height
	}

	li := &groupBoxLayoutItem{
		compositePos: compositePos,
		title:        gb.Title(),
	}

	gbli := CreateLayoutItemsForContainerWithContext(gb.composite, ctx)
	gbli.AsLayoutItemBase().parent = li

	li.children = append(li.children, gbli)

	return li
}

type groupBoxLayoutItem struct {
	ContainerLayoutItemBase
	compositePos Point // in native pixels
	title        string
}

func (li *groupBoxLayoutItem) LayoutFlags() LayoutFlags {
	return li.children[0].LayoutFlags()
}

func (li *groupBoxLayoutItem) MinSize() Size {
	min := li.children[0].(MinSizer).MinSize()
	min.Width += li.compositePos.X * 2
	min.Height += li.compositePos.Y + IntFrom96DPI(5, li.ctx.dpi)

	return min
}

func (li *groupBoxLayoutItem) MinSizeForSize(size Size) Size {
	return li.MinSize()
}

func (li *groupBoxLayoutItem) HasHeightForWidth() bool {
	return li.children[0].(HeightForWidther).HasHeightForWidth()
}

func (li *groupBoxLayoutItem) HeightForWidth(width int) int {
	return li.children[0].(HeightForWidther).HeightForWidth(width-li.compositePos.X*2) + li.compositePos.Y + IntFrom96DPI(5, li.ctx.dpi)
}

func (li *groupBoxLayoutItem) IdealSize() Size {
	size := li.children[0].(IdealSizer).IdealSize()
	size.Height += li.compositePos.Y
	return size
}

func (li *groupBoxLayoutItem) PerformLayout() []LayoutResultItem {
	return []LayoutResultItem{
		{
			Item:   li.children[0],
			Bounds: Rectangle{X: li.compositePos.X, Y: li.compositePos.Y, Width: li.geometry.Size.Width - li.compositePos.X*2, Height: li.geometry.Size.Height - li.compositePos.Y - IntFrom96DPI(5, li.ctx.dpi)},
		},
	}
}
