// Copyright 2014 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"unsafe"

	"github.com/lxn/win"
)

const scrollViewWindowClass = `\o/ Walk_ScrollView_Class \o/`

func init() {
	AppendToWalkInit(func() {
		MustRegisterWindowClass(scrollViewWindowClass)
	})
}

type ScrollView struct {
	WidgetBase
	composite  *Composite
	horizontal bool
	vertical   bool
}

func NewScrollView(parent Container) (*ScrollView, error) {
	sv := &ScrollView{horizontal: true, vertical: true}

	if err := InitWidget(
		sv,
		parent,
		scrollViewWindowClass,
		win.WS_CHILD|win.WS_HSCROLL|win.WS_VISIBLE|win.WS_VSCROLL,
		win.WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			sv.Dispose()
		}
	}()

	var err error
	if sv.composite, err = NewComposite(sv); err != nil {
		return nil, err
	}

	sv.composite.SizeChanged().Attach(func() {
		sv.updateScrollBars()
	})

	sv.SetBackground(NullBrush())

	succeeded = true

	return sv, nil
}

func (sv *ScrollView) AsContainerBase() *ContainerBase {
	if sv.composite == nil {
		return nil
	}

	return sv.composite.AsContainerBase()
}

func (sv *ScrollView) ApplyDPI(dpi int) {
	sv.WidgetBase.ApplyDPI(dpi)
	sv.composite.ApplyDPI(dpi)
}

func (sv *ScrollView) Scrollbars() (horizontal, vertical bool) {
	horizontal = sv.horizontal
	vertical = sv.vertical

	return
}

func (sv *ScrollView) SetScrollbars(horizontal, vertical bool) {
	sv.horizontal = horizontal
	sv.vertical = vertical

	sv.ensureStyleBits(win.WS_HSCROLL, horizontal)
	sv.ensureStyleBits(win.WS_VSCROLL, vertical)
}

func (sv *ScrollView) SetSuspended(suspend bool) {
	sv.composite.SetSuspended(suspend)
	sv.WidgetBase.SetSuspended(suspend)
	sv.Invalidate()
}

func (sv *ScrollView) DataBinder() *DataBinder {
	return sv.composite.dataBinder
}

func (sv *ScrollView) SetDataBinder(dataBinder *DataBinder) {
	sv.composite.SetDataBinder(dataBinder)
}

func (sv *ScrollView) Children() *WidgetList {
	if sv.composite == nil {
		// Without this we would get into trouble in NewComposite.
		return nil
	}

	return sv.composite.Children()
}

func (sv *ScrollView) Layout() Layout {
	if sv.composite == nil {
		return nil
	}

	return sv.composite.Layout()
}

func (sv *ScrollView) SetLayout(value Layout) error {
	return sv.composite.SetLayout(value)
}

func (sv *ScrollView) Name() string {
	if sv.composite == nil {
		return ""
	}

	return sv.composite.Name()
}

func (sv *ScrollView) SetName(name string) {
	sv.composite.SetName(name)
}

func (sv *ScrollView) Persistent() bool {
	return sv.composite.Persistent()
}

func (sv *ScrollView) SetPersistent(value bool) {
	sv.composite.SetPersistent(value)
}

func (sv *ScrollView) SaveState() error {
	return sv.composite.SaveState()
}

func (sv *ScrollView) RestoreState() error {
	return sv.composite.RestoreState()
}

func (sv *ScrollView) MouseDown() *MouseEvent {
	return sv.composite.MouseDown()
}

func (sv *ScrollView) MouseMove() *MouseEvent {
	return sv.composite.MouseMove()
}

func (sv *ScrollView) MouseUp() *MouseEvent {
	return sv.composite.MouseUp()
}

func (sv *ScrollView) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if sv.composite != nil {
		avoidBGArtifacts := func() {
			if sv.hasComplexBackground() {
				sv.composite.Invalidate()
			}
		}

		switch msg {
		case win.WM_HSCROLL:
			sv.composite.SetXPixels(sv.scroll(win.SB_HORZ, win.LOWORD(uint32(wParam))))
			if wParam == win.SB_ENDSCROLL {
				avoidBGArtifacts()
			}

		case win.WM_VSCROLL:
			sv.composite.SetYPixels(sv.scroll(win.SB_VERT, win.LOWORD(uint32(wParam))))
			if wParam == win.SB_ENDSCROLL {
				avoidBGArtifacts()
			}

		case win.WM_MOUSEWHEEL:
			if win.GetWindowLong(sv.hWnd, win.GWL_STYLE)&win.WS_VSCROLL == 0 {
				break
			}

			var cmd uint16
			if delta := int16(win.HIWORD(uint32(wParam))); delta < 0 {
				cmd = win.SB_LINEDOWN
			} else {
				cmd = win.SB_LINEUP
			}

			sv.composite.SetYPixels(sv.scroll(win.SB_VERT, cmd))
			avoidBGArtifacts()

			return 0

		case win.WM_COMMAND, win.WM_NOTIFY:
			sv.composite.WndProc(hwnd, msg, wParam, lParam)

		case win.WM_WINDOWPOSCHANGED:
			wp := (*win.WINDOWPOS)(unsafe.Pointer(lParam))

			if wp.Flags&win.SWP_NOSIZE != 0 {
				break
			}

			sv.updateScrollBars()

			if h, v := sv.Scrollbars(); !h || !v {
				sv.RequestLayout()
			}
		}
	}

	return sv.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

func (sv *ScrollView) updateScrollBars() {
	size := sv.SizePixels()
	compositeSize := sv.composite.SizePixels()

	var si win.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_PAGE | win.SIF_RANGE

	newCompositeBounds := Rectangle{Width: compositeSize.Width, Height: compositeSize.Height}

	if size != compositeSize {
		dpi := uint32(sv.DPI())

		vsbw := int(win.GetSystemMetricsForDpi(win.SM_CXVSCROLL, dpi))
		hsbh := int(win.GetSystemMetricsForDpi(win.SM_CYHSCROLL, dpi))

		if size.Width < compositeSize.Width && size.Height < compositeSize.Height {
			size.Width -= vsbw
			size.Height -= hsbh
		}
	}

	si.NMax = int32(compositeSize.Width - 1)
	si.NPage = uint32(size.Width)
	win.SetScrollInfo(sv.hWnd, win.SB_HORZ, &si, false)
	newCompositeBounds.X = sv.scroll(win.SB_HORZ, win.SB_THUMBPOSITION)

	si.NMax = int32(compositeSize.Height - 1)
	si.NPage = uint32(size.Height)
	win.SetScrollInfo(sv.hWnd, win.SB_VERT, &si, false)
	newCompositeBounds.Y = sv.scroll(win.SB_VERT, win.SB_THUMBPOSITION)

	sv.composite.SetBoundsPixels(newCompositeBounds)
}

// scroll scrolls and returns new position in native pixels.
func (sv *ScrollView) scroll(sb int32, cmd uint16) int {
	var pos int32
	var si win.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_PAGE | win.SIF_POS | win.SIF_RANGE | win.SIF_TRACKPOS

	win.GetScrollInfo(sv.hWnd, sb, &si)

	pos = si.NPos

	switch cmd {
	case win.SB_LINELEFT: // == win.SB_LINEUP
		pos -= int32(sv.IntFrom96DPI(20))

	case win.SB_LINERIGHT: // == win.SB_LINEDOWN
		pos += int32(sv.IntFrom96DPI(20))

	case win.SB_PAGELEFT: // == win.SB_PAGEUP
		pos -= int32(si.NPage)

	case win.SB_PAGERIGHT: // == win.SB_PAGEDOWN
		pos += int32(si.NPage)

	case win.SB_THUMBTRACK:
		pos = si.NTrackPos
	}

	if pos < 0 {
		pos = 0
	}
	if pos > si.NMax+1-int32(si.NPage) {
		pos = si.NMax + 1 - int32(si.NPage)
	}

	si.FMask = win.SIF_POS
	si.NPos = pos
	win.SetScrollInfo(sv.hWnd, sb, &si, true)

	return -int(pos)
}

func (sv *ScrollView) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	svli := new(scrollViewLayoutItem)
	svli.ctx = ctx
	cli := CreateLayoutItemsForContainerWithContext(sv.composite, ctx)
	cli.AsLayoutItemBase().parent = svli
	svli.children = append(svli.children, cli)

	if box, ok := cli.(*boxLayoutItem); ok {
		if len(box.children) > 0 {
			if _, ok := box.children[len(box.children)-1].(*spacerLayoutItem); !ok {
				// To retain the previous behavior with box layouts, we add a fake spacer at the end.
				// Maybe this should just be an option.
				box.children = append(box.children, &spacerLayoutItem{
					LayoutItemBase: LayoutItemBase{ctx: ctx},
					layoutFlags:    ShrinkableHorz | ShrinkableVert | GrowableVert | GreedyVert,
				})
			}
		}
	}

	svli.idealSize = cli.MinSize()

	h, v := sv.Scrollbars()

	if h {
		svli.layoutFlags |= ShrinkableHorz | GrowableHorz | GreedyHorz

		if !v {
			maxSize := SizeFrom96DPI(sv.maxSize96dpi, ctx.dpi)
			if svli.idealSize.Width > sv.geometry.ClientSize.Width && sv.geometry.ClientSize.Width > 0 && maxSize.Width == 0 ||
				svli.idealSize.Width > maxSize.Width && maxSize.Width > 0 {
				svli.sbSize.Height = int(win.GetSystemMetricsForDpi(win.SM_CYHSCROLL, uint32(ctx.dpi)))
				svli.idealSize.Height += svli.sbSize.Height
			}

			svli.minSize.Height = svli.idealSize.Height
		}
	}

	if v {
		svli.layoutFlags |= GreedyVert | GrowableVert | ShrinkableVert

		if !h {
			maxSize := SizeFrom96DPI(sv.maxSize96dpi, ctx.dpi)
			if svli.idealSize.Height > sv.geometry.ClientSize.Height && sv.geometry.ClientSize.Height > 0 && maxSize.Height == 0 ||
				svli.idealSize.Height > maxSize.Height && maxSize.Height > 0 {
				svli.sbSize.Width = int(win.GetSystemMetricsForDpi(win.SM_CXVSCROLL, uint32(ctx.dpi)))
				svli.idealSize.Width += svli.sbSize.Width
			}

			svli.minSize.Width = svli.idealSize.Width
		}
	}

	var si win.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_POS | win.SIF_RANGE

	win.GetScrollInfo(sv.hWnd, win.SB_HORZ, &si)
	svli.scrollX = float64(si.NPos) / float64(si.NMax)

	win.GetScrollInfo(sv.hWnd, win.SB_VERT, &si)
	svli.scrollY = float64(si.NPos) / float64(si.NMax)

	return svli
}

type scrollViewLayoutItem struct {
	ContainerLayoutItemBase
	idealSize   Size // in native pixels
	minSize     Size // in native pixels
	sbSize      Size // in native pixels
	layoutFlags LayoutFlags
	scrollX     float64
	scrollY     float64
}

func (li *scrollViewLayoutItem) LayoutFlags() LayoutFlags {
	return li.layoutFlags
}

func (li *scrollViewLayoutItem) IdealSize() Size {
	return li.idealSize
}

func (li *scrollViewLayoutItem) MinSize() Size {
	return li.minSize
}

func (li *scrollViewLayoutItem) MinSizeForSize(size Size) Size {
	return li.MinSize()
}

func (li *scrollViewLayoutItem) HasHeightForWidth() bool {
	return false
}

func (li *scrollViewLayoutItem) HeightForWidth(width int) int {
	return 0
}

func (li *scrollViewLayoutItem) PerformLayout() []LayoutResultItem {
	composite := li.children[0]

	clientSize := li.geometry.Size
	clientSize.Width -= li.sbSize.Width
	clientSize.Height -= li.sbSize.Height

	minSize := composite.(MinSizeForSizer).MinSizeForSize(clientSize)
	if hfw, ok := composite.(HeightForWidther); ok && hfw.HasHeightForWidth() {
		if minSize.Height > clientSize.Height {
			if minSize.Width > clientSize.Width {
				clientSize.Width = minSize.Width
				minSize = composite.(MinSizeForSizer).MinSizeForSize(clientSize)
			} else {
				clientSize.Width -= int(win.GetSystemMetricsForDpi(win.SM_CXVSCROLL, uint32(li.ctx.dpi)))
				minSize = composite.(MinSizeForSizer).MinSizeForSize(clientSize)
				if minSize.Width > clientSize.Width {
					clientSize.Width = minSize.Width
					minSize = composite.(MinSizeForSizer).MinSizeForSize(clientSize)
				}
			}
		}
	}

	s := maxSize(minSize, clientSize)

	var x, y int
	if clientSize.Width < minSize.Width {
		x = -int(float64(minSize.Width) * li.scrollX)
	}
	if clientSize.Height < minSize.Height {
		y = -int(float64(minSize.Height) * li.scrollY)
	}

	return []LayoutResultItem{
		{
			Item:   composite,
			Bounds: Rectangle{x, y, s.Width, s.Height},
		},
	}
}
