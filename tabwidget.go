// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"strconv"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const tabWidgetWindowClass = `\o/ Walk_TabWidget_Class \o/`

func init() {
	AppendToWalkInit(func() {
		MustRegisterWindowClass(tabWidgetWindowClass)
		tabWidgetTabWndProcPtr = syscall.NewCallback(tabWidgetTabWndProc)
	})
}

type TabWidget struct {
	WidgetBase
	hWndTab                      win.HWND
	tabOrigWndProcPtr            uintptr
	imageList                    *ImageList
	pages                        *TabPageList
	currentIndex                 int
	currentIndexChangedPublisher EventPublisher
	persistent                   bool
}

func NewTabWidget(parent Container) (*TabWidget, error) {
	tw := &TabWidget{currentIndex: -1}
	tw.pages = newTabPageList(tw)

	if err := InitWidget(
		tw,
		parent,
		tabWidgetWindowClass,
		win.WS_VISIBLE,
		win.WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			tw.Dispose()
		}
	}()

	tw.SetPersistent(true)

	tw.hWndTab = win.CreateWindowEx(
		0, syscall.StringToUTF16Ptr("SysTabControl32"), nil,
		win.WS_CHILD|win.WS_CLIPSIBLINGS|win.WS_TABSTOP|win.WS_VISIBLE,
		0, 0, 0, 0, tw.hWnd, 0, 0, nil)
	if tw.hWndTab == 0 {
		return nil, lastError("CreateWindowEx")
	}

	win.SetWindowLongPtr(tw.hWndTab, win.GWLP_USERDATA, uintptr(unsafe.Pointer(tw)))
	tw.tabOrigWndProcPtr = win.SetWindowLongPtr(tw.hWndTab, win.GWLP_WNDPROC, tabWidgetTabWndProcPtr)

	dpi := int(win.GetDpiForWindow(tw.hWndTab))
	win.SendMessage(tw.hWndTab, win.WM_SETFONT, uintptr(defaultFont.handleForDPI(dpi)), 1)

	tw.applyFont(tw.Font())

	tw.MustRegisterProperty("HasCurrentPage", NewReadOnlyBoolProperty(
		func() bool {
			return tw.CurrentIndex() != -1
		},
		tw.CurrentIndexChanged()))

	tw.MustRegisterProperty("CurrentIndex", NewProperty(
		func() interface{} {
			return tw.CurrentIndex()
		},
		func(v interface{}) error {
			return tw.SetCurrentIndex(assertIntOr(v, -1))
		},
		tw.CurrentIndexChanged()))

	succeeded = true

	return tw, nil
}

func (tw *TabWidget) Dispose() {
	tw.WidgetBase.Dispose()

	if tw.imageList != nil {
		tw.imageList.Dispose()
		tw.imageList = nil
	}
}

func (tw *TabWidget) applyEnabled(enabled bool) {
	tw.WidgetBase.applyEnabled(enabled)

	setWindowEnabled(tw.hWndTab, enabled)

	applyEnabledToDescendants(tw, enabled)
}

func (tw *TabWidget) applyFont(font *Font) {
	tw.WidgetBase.applyFont(font)

	SetWindowFont(tw.hWndTab, font)

	// FIXME: won't work with ApplyDPI
	// applyFontToDescendants(tw, font)
}

func (tw *TabWidget) ApplyDPI(dpi int) {
	tw.WidgetBase.ApplyDPI(dpi)

	var maskColor Color
	var size Size
	if tw.imageList != nil {
		maskColor = tw.imageList.maskColor
		size = SizeFrom96DPI(tw.imageList.imageSize96dpi, dpi)
	} else {
		size = SizeFrom96DPI(Size{16, 16}, dpi)
	}

	iml, err := NewImageListForDPI(size, maskColor, dpi)
	if err != nil {
		return
	}

	win.SendMessage(tw.hWndTab, win.TCM_SETIMAGELIST, 0, uintptr(iml.hIml))

	if tw.imageList != nil {
		tw.imageList.Dispose()
	}

	tw.imageList = iml

	for _, page := range tw.pages.items {
		tw.onPageChanged(page)
	}
}

func (tw *TabWidget) CurrentIndex() int {
	return tw.currentIndex
}

func (tw *TabWidget) SetCurrentIndex(index int) error {
	if index == tw.currentIndex {
		return nil
	}

	if index < 0 || index >= tw.pages.Len() {
		return newError("invalid index")
	}

	ret := int(win.SendMessage(tw.hWndTab, win.TCM_SETCURSEL, uintptr(index), 0))
	if ret == -1 {
		return newError("SendMessage(TCM_SETCURSEL) failed")
	}

	// FIXME: The SendMessage(TCM_SETCURSEL) call above doesn't cause a
	// TCN_SELCHANGE notification, so we use this workaround.
	tw.onSelChange()

	return nil
}

func (tw *TabWidget) CurrentIndexChanged() *Event {
	return tw.currentIndexChangedPublisher.Event()
}

func (tw *TabWidget) Pages() *TabPageList {
	return tw.pages
}

func (tw *TabWidget) Persistent() bool {
	return tw.persistent
}

func (tw *TabWidget) SetPersistent(value bool) {
	tw.persistent = value
}

func (tw *TabWidget) SaveState() error {
	tw.WriteState(strconv.Itoa(tw.CurrentIndex()))

	for _, page := range tw.pages.items {
		if err := page.SaveState(); err != nil {
			return err
		}
	}

	return nil
}

func (tw *TabWidget) RestoreState() error {
	state, err := tw.ReadState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	index, err := strconv.Atoi(state)
	if err != nil {
		return err
	}
	if index >= 0 && index < tw.pages.Len() {
		if err := tw.SetCurrentIndex(index); err != nil {
			return err
		}
	}

	for _, page := range tw.pages.items {
		if err := page.RestoreState(); err != nil {
			return err
		}
	}

	return nil
}

func (tw *TabWidget) resizePages() {
	bounds := tw.pageBounds()

	for _, page := range tw.pages.items {
		page.SetBoundsPixels(bounds)
	}
}

// pageBounds returns page bounds in native pixels.
func (tw *TabWidget) pageBounds() Rectangle {
	var r win.RECT
	if !win.GetWindowRect(tw.hWndTab, &r) {
		lastError("GetWindowRect")
		return Rectangle{}
	}

	p := win.POINT{
		r.Left,
		r.Top,
	}
	if !win.ScreenToClient(tw.hWnd, &p) {
		newError("ScreenToClient failed")
		return Rectangle{}
	}

	r = win.RECT{
		p.X,
		p.Y,
		r.Right - r.Left + p.X,
		r.Bottom - r.Top + p.Y,
	}
	win.SendMessage(tw.hWndTab, win.TCM_ADJUSTRECT, 0, uintptr(unsafe.Pointer(&r)))

	adjustment := 2 * int32(tw.IntFrom96DPI(1))
	return Rectangle{
		int(r.Left - adjustment),
		int(r.Top),
		int(r.Right - r.Left + adjustment),
		int(r.Bottom - r.Top),
	}
}

func (tw *TabWidget) onResize(width, height int32) {
	if !win.MoveWindow(tw.hWndTab, 0, 0, width, height, true) {
		lastError("MoveWindow")
		return
	}

	tw.resizePages()
}

func (tw *TabWidget) onSelChange() {
	pageCount := tw.pages.Len()

	if tw.currentIndex > -1 && tw.currentIndex < pageCount {
		page := tw.pages.At(tw.currentIndex)
		page.SetVisible(false)
	}

	tw.currentIndex = int(int32(win.SendMessage(tw.hWndTab, win.TCM_GETCURSEL, 0, 0)))

	if tw.currentIndex > -1 && tw.currentIndex < pageCount {
		page := tw.pages.At(tw.currentIndex)
		page.SetVisible(true)
		tw.RequestLayout()
		page.Invalidate()
		tw.pages.At(tw.currentIndex).focusFirstCandidateDescendant()
	}

	tw.Invalidate()

	tw.currentIndexChangedPublisher.Publish()
}

func (tw *TabWidget) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if tw.hWndTab != 0 {
		switch msg {
		case win.WM_ERASEBKGND:
			return 1

		case win.WM_WINDOWPOSCHANGED:
			wp := (*win.WINDOWPOS)(unsafe.Pointer(lParam))

			if wp.Flags&win.SWP_NOSIZE != 0 {
				break
			}

			tw.onResize(wp.Cx, wp.Cy)

		case win.WM_NOTIFY:
			nmhdr := (*win.NMHDR)(unsafe.Pointer(lParam))

			switch int32(nmhdr.Code) {
			case win.TCN_SELCHANGE:
				tw.onSelChange()
			}
		}
	}

	return tw.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

var tabWidgetTabWndProcPtr uintptr

func tabWidgetTabWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	tw := (*TabWidget)(unsafe.Pointer(win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)))

	switch msg {
	case win.WM_MOUSEMOVE:
		win.InvalidateRect(hwnd, nil, true)

	case win.WM_ERASEBKGND:
		return 1

	case win.WM_PAINT:
		var ps win.PAINTSTRUCT

		hdc := win.BeginPaint(hwnd, &ps)
		defer win.EndPaint(hwnd, &ps)

		cb := tw.ClientBoundsPixels()

		dpi := tw.DPI()
		bitmap, err := NewBitmapForDPI(cb.Size(), dpi)
		if err != nil {
			break
		}
		defer bitmap.Dispose()

		canvas, err := NewCanvasFromImage(bitmap)
		if err != nil {
			break
		}
		defer canvas.Dispose()

		themed := win.IsAppThemed()

		if !themed {
			if err := canvas.FillRectanglePixels(sysColorBtnFaceBrush, cb); err != nil {
				break
			}
		}

		win.SendMessage(hwnd, win.WM_PRINTCLIENT, uintptr(canvas.hdc), uintptr(win.PRF_CLIENT|win.PRF_CHILDREN|win.PRF_ERASEBKGND))

		parent := tw.Parent()
		if parent == nil {
			break
		}

		// Draw background of free area not occupied by tab items.
		if bg, wnd := parent.AsWindowBase().backgroundEffective(); bg != nil {
			tw.prepareDCForBackground(canvas.hdc, hwnd, wnd)

			hRgn := win.CreateRectRgn(0, 0, 0, 0)
			defer win.DeleteObject(win.HGDIOBJ(hRgn))

			var rc win.RECT

			adjustment := SizeFrom96DPI(Size{1, 1}, dpi).toSIZE()
			count := tw.pages.Len()
			for i := 0; i < count; i++ {
				if 0 == win.SendMessage(hwnd, win.TCM_GETITEMRECT, uintptr(i), uintptr(unsafe.Pointer(&rc))) {
					break
				}

				if i == tw.currentIndex {
					rc.Left -= 2 * adjustment.CX
					rc.Top -= 2 * adjustment.CY
					rc.Right += 2 * adjustment.CX
				} else {
					if i == count-1 && themed {
						rc.Right -= 2 * adjustment.CX
					}
				}

				hRgnTab := win.CreateRectRgn(rc.Left, rc.Top, rc.Right, rc.Bottom)
				win.CombineRgn(hRgn, hRgn, hRgnTab, win.RGN_OR)
				win.DeleteObject(win.HGDIOBJ(hRgnTab))
			}

			hRgnRC := win.CreateRectRgn(0, 0, int32(cb.Width), rc.Bottom)
			win.CombineRgn(hRgn, hRgnRC, hRgn, win.RGN_DIFF)
			win.DeleteObject(win.HGDIOBJ(hRgnRC))

			if !win.FillRgn(canvas.hdc, hRgn, bg.handle()) {
				break
			}
		}

		// Draw current tab item.
		if tw.currentIndex != -1 {
			page := tw.pages.At(tw.CurrentIndex())

			if bg, wnd := page.AsWindowBase().backgroundEffective(); bg != nil &&
				bg != tabPageBackgroundBrush &&
				(page.layout == nil || !page.layout.Margins().isZero()) {

				tw.prepareDCForBackground(canvas.hdc, hwnd, wnd)

				var rc win.RECT
				if 0 == win.SendMessage(hwnd, win.TCM_GETITEMRECT, uintptr(tw.currentIndex), uintptr(unsafe.Pointer(&rc))) {
					break
				}

				adjustment := SizeFrom96DPI(Size{6, 1}, dpi).toSIZE()
				hRgn := win.CreateRectRgn(rc.Left, rc.Top, rc.Right, rc.Bottom+2*adjustment.CY)
				defer win.DeleteObject(win.HGDIOBJ(hRgn))
				if !win.FillRgn(canvas.hdc, hRgn, bg.handle()) {
					break
				}

				if page.image != nil {
					x := rc.Left + adjustment.CX
					y := rc.Top
					s := int32(IntFrom96DPI(16, dpi))

					bmp, err := iconCache.Bitmap(page.image, dpi)
					if err == nil {
						if imageCanvas, err := NewCanvasFromImage(bmp); err == nil {
							defer imageCanvas.Dispose()

							if !win.TransparentBlt(
								canvas.hdc, x, y, s, s,
								imageCanvas.hdc, 0, 0, int32(bmp.size.Width), int32(bmp.size.Height),
								0) {
								break
							}
						}

						rc.Left += s + adjustment.CX
					}
				}

				rc.Left += adjustment.CX
				rc.Top += adjustment.CY

				title := syscall.StringToUTF16(page.title)

				if themed {
					hTheme := win.OpenThemeData(hwnd, syscall.StringToUTF16Ptr("tab"))
					defer win.CloseThemeData(hTheme)

					options := win.DTTOPTS{DwFlags: win.DTT_GLOWSIZE, IGlowSize: int32(IntFrom96DPI(3, dpi))}
					options.DwSize = uint32(unsafe.Sizeof(options))
					if hr := win.DrawThemeTextEx(hTheme, canvas.hdc, 0, win.TIS_SELECTED, &title[0], int32(len(title)), 0, &rc, &options); !win.SUCCEEDED(hr) {
						break
					}
				} else {
					if 0 == win.DrawTextEx(canvas.hdc, &title[0], int32(len(title)), &rc, 0, nil) {
						break
					}
				}
			}
		}

		if !win.BitBlt(hdc, 0, 0, int32(cb.Width), int32(cb.Height), canvas.hdc, 0, 0, win.SRCCOPY) {
			break
		}

		return 0

	case win.WM_LBUTTONDOWN:
		x := win.GET_X_LPARAM(lParam)
		y := win.GET_Y_LPARAM(lParam)

		hti := win.TCHITTESTINFO{
			Pt: win.POINT{x, y},
		}

		i := int(win.SendMessage(hwnd, win.TCM_HITTEST, 0, uintptr(unsafe.Pointer(&hti))))

		if i == -1 {
			break
		}

		ret := win.CallWindowProc(tw.tabOrigWndProcPtr, hwnd, msg, wParam, lParam)

		tw.pages.At(i).focusFirstCandidateDescendant()

		return ret
	}

	return win.CallWindowProc(tw.tabOrigWndProcPtr, hwnd, msg, wParam, lParam)
}

func (tw *TabWidget) onPageChanged(page *TabPage) (err error) {
	index := tw.pages.Index(page)
	item := tw.tcitemFromPage(page)

	if 0 == win.SendMessage(tw.hWndTab, win.TCM_SETITEM, uintptr(index), uintptr(unsafe.Pointer(item))) {
		return newError("SendMessage(TCM_SETITEM) failed")
	}

	return nil
}

func (tw *TabWidget) onInsertingPage(index int, page *TabPage) (err error) {
	return nil
}

func (tw *TabWidget) onInsertedPage(index int, page *TabPage) (err error) {
	item := tw.tcitemFromPage(page)

	if idx := int(win.SendMessage(tw.hWndTab, win.TCM_INSERTITEM, uintptr(index), uintptr(unsafe.Pointer(item)))); idx == -1 {
		return newError("SendMessage(TCM_INSERTITEM) failed")
	}

	page.SetVisible(false)

	style := uint32(win.GetWindowLong(page.hWnd, win.GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	style |= win.WS_CHILD
	style &^= win.WS_POPUP

	win.SetLastError(0)
	if win.SetWindowLong(page.hWnd, win.GWL_STYLE, int32(style)) == 0 {
		return lastError("SetWindowLong")
	}

	if win.SetParent(page.hWnd, tw.hWnd) == 0 {
		return lastError("SetParent")
	}

	if tw.pages.Len() == 1 {
		page.SetVisible(true)
		tw.SetCurrentIndex(0)
	}

	tw.resizePages()

	page.tabWidget = tw

	page.applyFont(tw.Font())

	tw.Invalidate()

	return
}

func (tw *TabWidget) removePage(page *TabPage) (err error) {
	page.SetVisible(false)

	style := uint32(win.GetWindowLong(page.hWnd, win.GWL_STYLE))
	if style == 0 {
		return lastError("GetWindowLong")
	}

	style &^= win.WS_CHILD
	style |= win.WS_POPUP

	win.SetLastError(0)
	if win.SetWindowLong(page.hWnd, win.GWL_STYLE, int32(style)) == 0 {
		return lastError("SetWindowLong")
	}

	page.tabWidget = nil

	return page.SetParent(nil)
}

func (tw *TabWidget) onRemovingPage(index int, page *TabPage) (err error) {
	return nil
}

func (tw *TabWidget) onRemovedPage(index int, page *TabPage) (err error) {
	err = tw.removePage(page)
	if err != nil {
		return
	}

	win.SendMessage(tw.hWndTab, win.TCM_DELETEITEM, uintptr(index), 0)

	if tw.pages.Len() > 0 {
		tw.currentIndex = 0
		win.SendMessage(tw.hWndTab, win.TCM_SETCURSEL, uintptr(tw.currentIndex), 0)
	} else {
		tw.currentIndex = -1
	}
	tw.onSelChange()

	return

	// FIXME: Either make use of this unreachable code or remove it.
	if index == tw.currentIndex {
		// removal of current visible tabpage...
		tw.currentIndex = -1

		// select new tabpage if any :
		if tw.pages.Len() > 0 {
			// are we removing the rightmost page ?
			if index == tw.pages.Len()-1 {
				// If so, select the page on the left
				index -= 1
			}
		}
	}

	tw.SetCurrentIndex(index)

	tw.Invalidate()

	return
}

func (tw *TabWidget) onClearingPages(pages []*TabPage) (err error) {
	return nil
}

func (tw *TabWidget) onClearedPages(pages []*TabPage) (err error) {
	win.SendMessage(tw.hWndTab, win.TCM_DELETEALLITEMS, 0, 0)
	for _, page := range pages {
		tw.removePage(page)
	}
	tw.currentIndex = -1

	tw.Invalidate()

	return nil
}

func (tw *TabWidget) tcitemFromPage(page *TabPage) *win.TCITEM {
	var imageIndex int32 = -1
	if page.image != nil {
		if bmp, err := iconCache.Bitmap(page.image, tw.DPI()); err == nil {
			imageIndex, _ = tw.imageIndex(bmp)
		}
	}

	text := syscall.StringToUTF16(page.title)

	item := &win.TCITEM{
		Mask:       win.TCIF_IMAGE | win.TCIF_TEXT,
		IImage:     imageIndex,
		PszText:    &text[0],
		CchTextMax: int32(len(text)),
	}

	return item
}

func (tw *TabWidget) imageIndex(image *Bitmap) (index int32, err error) {
	index = -1
	if image != nil {
		if tw.imageList == nil {
			dpi := tw.DPI()
			if tw.imageList, err = NewImageListForDPI(SizeFrom96DPI(Size{16, 16}, dpi), 0, dpi); err != nil {
				return
			}

			win.SendMessage(tw.hWndTab, win.TCM_SETIMAGELIST, 0, uintptr(tw.imageList.hIml))
		}

		if index, err = tw.imageList.AddMasked(image); err != nil {
			return
		}
	}

	return
}

func (tw *TabWidget) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	pages := make([]LayoutItem, tw.pages.Len())

	bounds := tw.pageBounds()

	li := &tabWidgetLayoutItem{
		pagePos:      bounds.Location(),
		currentIndex: tw.CurrentIndex(),
	}

	for i := tw.pages.Len() - 1; i >= 0; i-- {
		var page LayoutItem
		if p := tw.pages.At(i); p.Layout() != nil {
			page = CreateLayoutItemsForContainerWithContext(p, ctx)
		} else {
			page = NewGreedyLayoutItem()
		}

		lib := page.AsLayoutItemBase()
		lib.ctx = ctx
		lib.parent = li
		pages[i] = page
	}

	li.children = pages

	return li
}

type tabWidgetLayoutItem struct {
	ContainerLayoutItemBase
	pagePos      Point // in native pixels
	currentIndex int
}

func (li *tabWidgetLayoutItem) LayoutFlags() LayoutFlags {
	if len(li.children) == 0 {
		return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
	}

	var flags LayoutFlags

	for _, page := range li.children {
		flags |= page.LayoutFlags()
	}

	return flags
}

func (li *tabWidgetLayoutItem) MinSize() Size {
	if len(li.children) == 0 {
		return li.IdealSize()
	}

	var min Size

	for _, page := range li.children {
		if ms, ok := page.(MinSizer); ok {
			s := ms.MinSize()

			min.Width = maxi(min.Width, s.Width)
			min.Height = maxi(min.Height, s.Height)
		}
	}

	s := li.geometry.Size
	ps := li.children[0].Geometry().Size

	size := Size{s.Width - ps.Width + min.Width, s.Height - ps.Height + min.Height}

	return size
}

func (li *tabWidgetLayoutItem) MinSizeForSize(size Size) Size {
	return li.MinSize()
}

func (li *tabWidgetLayoutItem) HasHeightForWidth() bool {
	if len(li.children) == 0 {
		return false
	}

	for _, page := range li.children {
		if hfw, ok := page.(HeightForWidther); ok && hfw.HasHeightForWidth() {
			return true
		}
	}

	return false
}

func (li *tabWidgetLayoutItem) HeightForWidth(width int) int {
	if len(li.children) == 0 {
		return 0
	}

	var height int
	margin := li.geometry.Size
	pageSize := li.children[0].Geometry().Size

	margin.Width -= pageSize.Width
	margin.Height -= pageSize.Height

	for _, page := range li.children {
		if hfw, ok := page.(HeightForWidther); ok && hfw.HasHeightForWidth() {
			h := hfw.HeightForWidth(width + margin.Width)

			height = maxi(height, h)
		}
	}

	return height + margin.Height
}

func (li *tabWidgetLayoutItem) IdealSize() Size {
	return li.MinSize()
}

func (li *tabWidgetLayoutItem) PerformLayout() []LayoutResultItem {
	if li.currentIndex > -1 {
		page := li.children[li.currentIndex]

		adjustment := IntFrom96DPI(1, li.ctx.dpi)
		return []LayoutResultItem{
			{
				Item: page,
				Bounds: Rectangle{
					X:      li.pagePos.X,
					Y:      li.pagePos.Y,
					Width:  li.geometry.Size.Width - li.pagePos.X*2 - adjustment,
					Height: li.geometry.Size.Height - li.pagePos.Y - 2*adjustment,
				},
			},
		}
	}

	return nil
}
