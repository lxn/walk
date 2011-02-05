// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/gdi32"
	. "walk/winapi/kernel32"
	. "walk/winapi/user32"
)

type FontStyle byte

// Font style flags
const (
	FontBold      FontStyle = 0x01
	FontItalic    FontStyle = 0x02
	FontUnderline FontStyle = 0x04
	FontStrikeOut FontStyle = 0x08
)

var (
	screenDPIX  int
	screenDPIY  int
	defaultFont *Font
)

func init() {
	// Retrieve screen DPI
	hDC := GetDC(0)
	defer ReleaseDC(0, hDC)
	screenDPIX = GetDeviceCaps(hDC, LOGPIXELSX)
	screenDPIY = GetDeviceCaps(hDC, LOGPIXELSY)

	// Initialize default font
	var ncm NONCLIENTMETRICS
	ncm.CbSize = uint(unsafe.Sizeof(ncm))

	if !SystemParametersInfo(SPI_GETNONCLIENTMETRICS, ncm.CbSize, unsafe.Pointer(&ncm), 0) {
		panic("SystemParametersInfo failed")
	}

	var err os.Error
	defaultFont, err = newFontFromLOGFONT(&ncm.LfMenuFont, screenDPIY)
	if err != nil {
		panic("failed to create default font")
	}
}

// Font represents a typographic typeface that is used for text drawing
// operations and on many GUI widgets.
type Font struct {
	dpi2hFont map[int]HFONT
	family    string
	pointSize int
	style     FontStyle
}

// NewFont returns a new Font with the specified attributes.
func NewFont(family string, pointSize int, style FontStyle) (*Font, os.Error) {
	if style > FontBold|FontItalic|FontUnderline|FontStrikeOut {
		return nil, newError("invalid style")
	}

	font := &Font{
		family:    family,
		pointSize: pointSize,
		style:     style,
		dpi2hFont: make(map[int]HFONT),
	}

	hFont := font.createForDPI(screenDPIY)
	if hFont == 0 {
		return nil, newError("CreateFontIndirect failed")
	}

	font.dpi2hFont[screenDPIY] = hFont
	font.dpi2hFont[0] = hFont // Make HFONT for screen easier accessible.

	return font, nil
}

func newFontFromLOGFONT(lf *LOGFONT, dpi int) (*Font, os.Error) {
	if lf == nil {
		return nil, newError("lf cannot be nil")
	}

	family := UTF16PtrToString(&lf.LfFaceName[0])
	pointSize := MulDiv(lf.LfHeight, 72, dpi)
	if pointSize < 0 {
		pointSize = -pointSize
	}

	var style FontStyle
	if lf.LfWeight > FW_NORMAL {
		style |= FontBold
	}
	if lf.LfItalic == TRUE {
		style |= FontItalic
	}
	if lf.LfUnderline == TRUE {
		style |= FontUnderline
	}
	if lf.LfStrikeOut == TRUE {
		style |= FontStrikeOut
	}

	return NewFont(family, pointSize, style)
}

func (f *Font) createForDPI(dpi int) HFONT {
	var lf LOGFONT

	lf.LfHeight = -MulDiv(f.pointSize, dpi, 72)
	if f.style&FontBold > 0 {
		lf.LfWeight = FW_BOLD
	} else {
		lf.LfWeight = FW_NORMAL
	}
	if f.style&FontItalic > 0 {
		lf.LfItalic = 1
	}
	if f.style&FontUnderline > 0 {
		lf.LfUnderline = 1
	}
	if f.style&FontStrikeOut > 0 {
		lf.LfStrikeOut = 1
	}
	lf.LfCharSet = DEFAULT_CHARSET
	lf.LfOutPrecision = OUT_TT_PRECIS
	lf.LfClipPrecision = CLIP_DEFAULT_PRECIS
	lf.LfQuality = CLEARTYPE_QUALITY
	lf.LfPitchAndFamily = VARIABLE_PITCH | FF_SWISS

	src := syscall.StringToUTF16(f.family)
	dest := lf.LfFaceName[:]
	copy(dest, src)

	return CreateFontIndirect(&lf)
}

// Bold returns if text drawn using the Font appears with
// greater weight than normal.
func (f *Font) Bold() bool {
	return f.style&FontBold > 0
}

// Dispose releases the os resources that were allocated for the Font.
//
// The Font can no longer be used for drawing operations or with GUI widgets
// after calling this method. It is safe to call Dispose multiple times.
func (f *Font) Dispose() {
	for dpi, hFont := range f.dpi2hFont {
		if dpi != 0 {
			DeleteObject(HGDIOBJ(hFont))
		}

		f.dpi2hFont[dpi] = 0, false
	}
}

// Family returns the family name of the Font.
func (f *Font) Family() string {
	return f.family
}

// Italic returns if text drawn using the Font appears slanted.
func (f *Font) Italic() bool {
	return f.style&FontItalic > 0
}

// HandleForDPI returns the os resource handle of the font for the specified
// DPI value.
//
// A value of 0 returns a HFONT suitable for the screen.
func (f *Font) handleForDPI(dpi int) HFONT {
	hFont := f.dpi2hFont[dpi]
	if hFont == 0 {
		hFont = f.createForDPI(dpi)
		f.dpi2hFont[dpi] = hFont
	}

	return hFont
}

// StrikeOut returns if text drawn using the Font appears striked out.
func (f *Font) StrikeOut() bool {
	return f.style&FontStrikeOut > 0
}

// Style returns the combination of style flags of the Font.
func (f *Font) Style() FontStyle {
	return f.style
}

// Underline returns if text drawn using the font appears underlined.
func (f *Font) Underline() bool {
	return f.style&FontUnderline > 0
}

// PointSize returns the size of the Font in point units.
func (f *Font) PointSize() int {
	return f.pointSize
}
