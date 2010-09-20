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

// Font flags
const (
	font_bold      = 0x01
	font_italic    = 0x02
	font_underline = 0x04
	font_strikeOut = 0x08
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

// FontInfo carries information about a font.
type FontInfo struct {
	Family    string
	PointSize float
	Bold      bool
	Italic    bool
	StrikeOut bool
	Underline bool
}

// Font represents a typographic typeface that is used for text drawing
// operations and on many GUI widgets.
type Font struct {
	hFont     HFONT
	family    string
	pointSize float
	flags     byte
}

// NewFont returns a new Font, initialized using the contents of the specified
// FontInfo.
//
// Family and PointSize are required.
func NewFont(info *FontInfo) (*Font, os.Error) {
	if info == nil {
		return nil, newError("info cannot be nil")
	}
	if info.Family == "" {
		return nil, newError("Family is required")
	}
	if info.PointSize < 0.000001 {
		return nil, newError("PointSize must be > 0")
	}

	font := &Font{family: info.Family, pointSize: info.PointSize}

	var lf LOGFONT

	lf.LfHeight = int(info.PointSize * float(screenDPIY) / float(72))
	if info.Bold {
		font.flags |= font_bold
		lf.LfWeight = FW_BOLD
	} else {
		lf.LfWeight = FW_NORMAL
	}
	if info.Italic {
		font.flags |= font_italic
		lf.LfItalic = 1
	}
	if info.Underline {
		font.flags |= font_underline
		lf.LfUnderline = 1
	}
	if info.StrikeOut {
		font.flags |= font_strikeOut
		lf.LfStrikeOut = 1
	}
	lf.LfCharSet = DEFAULT_CHARSET
	lf.LfOutPrecision = OUT_TT_PRECIS
	lf.LfClipPrecision = CLIP_DEFAULT_PRECIS
	lf.LfQuality = CLEARTYPE_QUALITY
	lf.LfPitchAndFamily = VARIABLE_PITCH | FF_SWISS

	src := syscall.StringToUTF16(info.Family)
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
	return f.flags&font_bold != 0
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
	return f.flags&font_italic != 0
}

// Handle returns the os resource handle of the font.
func (f *Font) Handle() HFONT {
	return f.hFont
}

// StrikeOut returns if text drawn using the Font appears striked out.
func (f *Font) StrikeOut() bool {
	return f.flags&font_strikeOut != 0
}

// Underline returns if text drawn using the font appears underlined.
func (f *Font) Underline() bool {
	return f.flags&font_underline != 0
}

// PointSize returns the size of the Font in point units.
func (f *Font) PointSize() float {
	return f.pointSize
}
