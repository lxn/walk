// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drawing

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/gdi32"
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
	screenDPIX int
	screenDPIY int
)

func init() {
	// Retrieve screen DPI
	hDC := GetDC(0)
	defer ReleaseDC(0, hDC)
	screenDPIX = GetDeviceCaps(hDC, LOGPIXELSX)
	screenDPIY = GetDeviceCaps(hDC, LOGPIXELSY)
}

// Font represents a typographic typeface that is used for text drawing
// operations and on many GUI widgets.
type Font struct {
	hFont     HFONT
	family    string
	pointSize float
	style     FontStyle
}

// NewFont returns a new Font with the specified attributes.
func NewFont(family string, pointSize float, style FontStyle) (*Font, os.Error) {
	if style > FontBold|FontItalic|FontUnderline|FontStrikeOut {
		return nil, newError("invalid style")
	}

	font := &Font{family: family, pointSize: pointSize, style: style}

	var lf LOGFONT

	lf.LfHeight = int(pointSize * float(screenDPIY) / float(72))
	if style&FontBold > 0 {
		lf.LfWeight = FW_BOLD
	} else {
		lf.LfWeight = FW_NORMAL
	}
	if style&FontItalic > 0 {
		lf.LfItalic = 1
	}
	if style&FontUnderline > 0 {
		lf.LfUnderline = 1
	}
	if style&FontStrikeOut > 0 {
		lf.LfStrikeOut = 1
	}
	lf.LfCharSet = DEFAULT_CHARSET
	lf.LfOutPrecision = OUT_TT_PRECIS
	lf.LfClipPrecision = CLIP_DEFAULT_PRECIS
	lf.LfQuality = CLEARTYPE_QUALITY
	lf.LfPitchAndFamily = VARIABLE_PITCH | FF_SWISS

	src := syscall.StringToUTF16(family)
	dest := lf.LfFaceName[0:]
	copy(dest, src)

	font.hFont = CreateFontIndirect(&lf)
	if font.hFont == 0 {
		return nil, newError("CreateFontIndirect failed")
	}

	return font, nil
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
	if f.hFont != 0 {
		DeleteObject(HGDIOBJ(f.hFont))
		f.hFont = 0
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

// Handle returns the os resource handle of the font.
func (f *Font) Handle() HFONT {
	return f.hFont
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
func (f *Font) PointSize() float {
	return f.pointSize
}
