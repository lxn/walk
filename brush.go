// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

type HatchStyle int

const (
	HatchHorizontal       HatchStyle = win.HS_HORIZONTAL
	HatchVertical         HatchStyle = win.HS_VERTICAL
	HatchForwardDiagonal  HatchStyle = win.HS_FDIAGONAL
	HatchBackwardDiagonal HatchStyle = win.HS_BDIAGONAL
	HatchCross            HatchStyle = win.HS_CROSS
	HatchDiagonalCross    HatchStyle = win.HS_DIAGCROSS
)

type SystemColor int

const (
	SysColor3DDkShadow              SystemColor = win.COLOR_3DDKSHADOW
	SysColor3DFace                  SystemColor = win.COLOR_3DFACE
	SysColor3DHighlight             SystemColor = win.COLOR_3DHIGHLIGHT
	SysColor3DLight                 SystemColor = win.COLOR_3DLIGHT
	SysColor3DShadow                SystemColor = win.COLOR_3DSHADOW
	SysColorActiveBorder            SystemColor = win.COLOR_ACTIVEBORDER
	SysColorActiveCaption           SystemColor = win.COLOR_ACTIVECAPTION
	SysColorAppWorkspace            SystemColor = win.COLOR_APPWORKSPACE
	SysColorBackground              SystemColor = win.COLOR_BACKGROUND
	SysColorDesktop                 SystemColor = win.COLOR_DESKTOP
	SysColorBtnFace                 SystemColor = win.COLOR_BTNFACE
	SysColorBtnHighlight            SystemColor = win.COLOR_BTNHIGHLIGHT
	SysColorBtnShadow               SystemColor = win.COLOR_BTNSHADOW
	SysColorBtnText                 SystemColor = win.COLOR_BTNTEXT
	SysColorCaptionText             SystemColor = win.COLOR_CAPTIONTEXT
	SysColorGrayText                SystemColor = win.COLOR_GRAYTEXT
	SysColorHighlight               SystemColor = win.COLOR_HIGHLIGHT
	SysColorHighlightText           SystemColor = win.COLOR_HIGHLIGHTTEXT
	SysColorInactiveBorder          SystemColor = win.COLOR_INACTIVEBORDER
	SysColorInactiveCaption         SystemColor = win.COLOR_INACTIVECAPTION
	SysColorInactiveCaptionText     SystemColor = win.COLOR_INACTIVECAPTIONTEXT
	SysColorInfoBk                  SystemColor = win.COLOR_INFOBK
	SysColorInfoText                SystemColor = win.COLOR_INFOTEXT
	SysColorMenu                    SystemColor = win.COLOR_MENU
	SysColorMenuText                SystemColor = win.COLOR_MENUTEXT
	SysColorScrollBar               SystemColor = win.COLOR_SCROLLBAR
	SysColorWindow                  SystemColor = win.COLOR_WINDOW
	SysColorWindowFrame             SystemColor = win.COLOR_WINDOWFRAME
	SysColorWindowText              SystemColor = win.COLOR_WINDOWTEXT
	SysColorHotLight                SystemColor = win.COLOR_HOTLIGHT
	SysColorGradientActiveCaption   SystemColor = win.COLOR_GRADIENTACTIVECAPTION
	SysColorGradientInactiveCaption SystemColor = win.COLOR_GRADIENTINACTIVECAPTION
)

type Brush interface {
	Dispose()
	handle() win.HBRUSH
	logbrush() *win.LOGBRUSH
	attachWindow(wb *WindowBase)
	detachWindow(wb *WindowBase)
}

type brushBase struct {
	hBrush win.HBRUSH
	wb2int map[*WindowBase]int
}

func (bb *brushBase) Dispose() {
	if bb.hBrush != 0 {
		win.DeleteObject(win.HGDIOBJ(bb.hBrush))

		bb.hBrush = 0
	}
}

func (bb *brushBase) handle() win.HBRUSH {
	return bb.hBrush
}

func (bb *brushBase) attachWindow(wb *WindowBase) {
	if wb == nil {
		return
	}

	if bb.wb2int == nil {
		bb.wb2int = make(map[*WindowBase]int)
	}

	bb.wb2int[wb] = -1
}

func (bb *brushBase) detachWindow(wb *WindowBase) {
	if bb.wb2int == nil || wb == nil {
		return
	}

	delete(bb.wb2int, wb)

	if len(bb.wb2int) == 0 {
		bb.Dispose()
	}
}

type nullBrush struct {
	brushBase
}

func newNullBrush() *nullBrush {
	lb := &win.LOGBRUSH{LbStyle: win.BS_NULL}

	hBrush := win.CreateBrushIndirect(lb)
	if hBrush == 0 {
		panic("failed to create null brush")
	}

	return &nullBrush{brushBase: brushBase{hBrush: hBrush}}
}

func (b *nullBrush) Dispose() {
	if b == nullBrushSingleton {
		return
	}

	b.brushBase.Dispose()
}

func (b *nullBrush) logbrush() *win.LOGBRUSH {
	return &win.LOGBRUSH{LbStyle: win.BS_NULL}
}

var nullBrushSingleton Brush = newNullBrush()

func NullBrush() Brush {
	return nullBrushSingleton
}

type SystemColorBrush struct {
	brushBase
	sysColor SystemColor
}

var sysColorBtnFaceBrush, _ = NewSystemColorBrush(SysColorBtnFace)

func NewSystemColorBrush(sysColor SystemColor) (*SystemColorBrush, error) {
	hBrush := win.GetSysColorBrush(int(sysColor))
	if hBrush == 0 {
		return nil, newError("GetSysColorBrush failed")
	}

	return &SystemColorBrush{brushBase: brushBase{hBrush: hBrush}, sysColor: sysColor}, nil
}

func (b *SystemColorBrush) Color() Color {
	return Color(win.GetSysColor(int(b.sysColor)))
}

func (b *SystemColorBrush) SystemColor() SystemColor {
	return b.sysColor
}

func (b *SystemColorBrush) Dispose() {
	// nop
}

func (b *SystemColorBrush) logbrush() *win.LOGBRUSH {
	return &win.LOGBRUSH{
		LbStyle: win.BS_SOLID,
		LbColor: win.COLORREF(win.GetSysColor(int(b.sysColor))),
	}
}

type SolidColorBrush struct {
	brushBase
	color Color
}

func NewSolidColorBrush(color Color) (*SolidColorBrush, error) {
	lb := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: win.COLORREF(color)}

	hBrush := win.CreateBrushIndirect(lb)
	if hBrush == 0 {
		return nil, newError("CreateBrushIndirect failed")
	}

	return &SolidColorBrush{brushBase: brushBase{hBrush: hBrush}, color: color}, nil
}

func (b *SolidColorBrush) Color() Color {
	return b.color
}

func (b *SolidColorBrush) logbrush() *win.LOGBRUSH {
	return &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: win.COLORREF(b.color)}
}

type HatchBrush struct {
	brushBase
	color Color
	style HatchStyle
}

func NewHatchBrush(color Color, style HatchStyle) (*HatchBrush, error) {
	lb := &win.LOGBRUSH{LbStyle: win.BS_HATCHED, LbColor: win.COLORREF(color), LbHatch: uintptr(style)}

	hBrush := win.CreateBrushIndirect(lb)
	if hBrush == 0 {
		return nil, newError("CreateBrushIndirect failed")
	}

	return &HatchBrush{brushBase: brushBase{hBrush: hBrush}, color: color, style: style}, nil
}

func (b *HatchBrush) Color() Color {
	return b.color
}

func (b *HatchBrush) logbrush() *win.LOGBRUSH {
	return &win.LOGBRUSH{LbStyle: win.BS_HATCHED, LbColor: win.COLORREF(b.color), LbHatch: uintptr(b.style)}
}

func (b *HatchBrush) Style() HatchStyle {
	return b.style
}

type BitmapBrush struct {
	brushBase
	bitmap *Bitmap
}

func NewBitmapBrush(bitmap *Bitmap) (*BitmapBrush, error) {
	if bitmap == nil {
		return nil, newError("bitmap cannot be nil")
	}

	hBrush := win.CreatePatternBrush(bitmap.hBmp)
	if hBrush == 0 {
		return nil, newError("CreatePatternBrush failed")
	}

	return &BitmapBrush{brushBase: brushBase{hBrush: hBrush}, bitmap: bitmap}, nil
}

func (b *BitmapBrush) logbrush() *win.LOGBRUSH {
	return &win.LOGBRUSH{LbStyle: win.BS_DIBPATTERN, LbColor: win.DIB_RGB_COLORS, LbHatch: uintptr(b.bitmap.hPackedDIB)}
}

func (b *BitmapBrush) Bitmap() *Bitmap {
	return b.bitmap
}
