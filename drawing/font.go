// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drawing

import (
	"container/vector"
	"os"
	"syscall"
)

import (
	. "walk/winapi/gdi32"
	. "walk/winapi/user32"
)

// Font flags
const (
	font_bold          = 0x01
	font_italic        = 0x02
	font_underline     = 0x04
	font_strikeOut     = 0x08
	font_suspendUpdate = 0x40
	font_dirty         = 0x80
)

var (
	screenDPIX int
	screenDPIY int
)

func init() {
	// Retrieve screen DPI
	hDC := GetDC(0)
	screenDPIX = GetDeviceCaps(hDC, LOGPIXELSX)
	screenDPIY = GetDeviceCaps(hDC, LOGPIXELSY)
	ReleaseDC(0, hDC)
}

type FontChangeHandler func(font *Font)

type Font struct {
	hFont            HFONT
	family           string
	pointSize        float
	flags            byte
	changingHandlers vector.Vector
	changedHandlers  vector.Vector
}

func NewFont() *Font {
	return &Font{}
}

func (f *Font) create() HFONT {
	var lf LOGFONT

	lf.LfHeight = int(f.pointSize * float(screenDPIY) / float(72))
	if f.Bold() {
		lf.LfWeight = FW_BOLD
	} else {
		lf.LfWeight = FW_NORMAL
	}
	if f.Italic() {
		lf.LfItalic = 1
	}
	if f.Underline() {
		lf.LfUnderline = 1
	}
	if f.StrikeOut() {
		lf.LfStrikeOut = 1
	}
	lf.LfCharSet = DEFAULT_CHARSET
	lf.LfOutPrecision = OUT_TT_PRECIS
	lf.LfClipPrecision = CLIP_DEFAULT_PRECIS
	lf.LfQuality = CLEARTYPE_QUALITY
	lf.LfPitchAndFamily = VARIABLE_PITCH | FF_SWISS

	src := syscall.StringToUTF16(f.family)
	dest := lf.LfFaceName[0:]
	copy(dest, src)

	return CreateFontIndirect(&lf)
}

func (f *Font) update() (err os.Error) {
	f.raiseChanging()

	hFont := f.create()
	if hFont == 0 {
		err = newError("failed to create font")
		return
	}

	f.Dispose()

	f.hFont = hFont

	f.setFlag(font_dirty, false)

	f.raiseChanged()

	return
}

func (f *Font) Handle() HFONT {
	return f.hFont
}

func (f *Font) flag(flag byte) bool {
	return (f.flags & flag) != 0
}

func (f *Font) setFlag(flag byte, value bool) {
	if value {
		f.flags |= flag
	} else {
		f.flags &^= flag
	}
}

func (f *Font) setFlagValue(flag byte, value bool) (err os.Error) {
	old := f.flag(flag)
	if value != old {
		wasDirty := f.isDirty()

		f.setFlag(flag, value)

		err = f.setDirty()
		if err != nil {
			f.setFlag(flag, old)
			if wasDirty {
				f.setDirty()
			}
		}
	}

	return
}

func (f *Font) isDirty() bool {
	return f.flag(font_dirty)
}

func (f *Font) setDirty() (err os.Error) {
	f.setFlag(font_dirty, true)

	if !f.flag(font_suspendUpdate) {
		err = f.update()
	}

	return
}

func (f *Font) Dispose() {
	if f.hFont != 0 {
		DeleteObject(HGDIOBJ(f.hFont))
		f.hFont = 0
	}
}

func (f *Font) IsDisposed() bool {
	return f.hFont == 0
}

func (f *Font) BeginEdit() {
	f.setFlag(font_suspendUpdate, true)
}

func (f *Font) EndEdit() (err os.Error) {
	f.setFlag(font_suspendUpdate, false)

	if f.isDirty() {
		err = f.update()
	}

	return
}

func (f *Font) Family() string {
	return f.family
}

func (f *Font) SetFamily(value string) (err os.Error) {
	old := f.family
	if value != old {
		wasDirty := f.isDirty()

		f.family = value

		err = f.setDirty()
		if err != nil {
			f.family = old
			if wasDirty {
				f.setDirty()
			}
		}
	}

	return
}

func (f *Font) Bold() bool {
	return f.flag(font_bold)
}

func (f *Font) SetBold(value bool) os.Error {
	return f.setFlagValue(font_bold, value)
}

func (f *Font) Italic() bool {
	return f.flag(font_italic)
}

func (f *Font) SetItalic(value bool) os.Error {
	return f.setFlagValue(font_italic, value)
}

func (f *Font) StrikeOut() bool {
	return f.flag(font_strikeOut)
}

func (f *Font) SetStrikeOut(value bool) os.Error {
	return f.setFlagValue(font_strikeOut, value)
}

func (f *Font) Underline() bool {
	return f.flag(font_underline)
}

func (f *Font) SetUnderline(value bool) os.Error {
	return f.setFlagValue(font_underline, value)
}

func (f *Font) PointSize() float {
	return f.pointSize
}

func (f *Font) SetPointSize(value float) (err os.Error) {
	old := f.pointSize
	if value != old {
		wasDirty := f.isDirty()

		f.pointSize = value

		err = f.setDirty()
		if err != nil {
			f.pointSize = old
			if wasDirty {
				f.setDirty()
			}
		}
	}

	return
}

func (f *Font) AddChangingHandler(handler FontChangeHandler) {
	f.changingHandlers.Push(handler)
}

func (f *Font) RemoveChangingHandler(handler FontChangeHandler) {
	for i, h := range f.changingHandlers {
		if h.(FontChangeHandler) == handler {
			f.changingHandlers.Delete(i)
			break
		}
	}
}

func (f *Font) AddChangedHandler(handler FontChangeHandler) {
	f.changedHandlers.Push(handler)
}

func (f *Font) RemoveChangedHandler(handler FontChangeHandler) {
	for i, h := range f.changedHandlers {
		if h.(FontChangeHandler) == handler {
			f.changedHandlers.Delete(i)
			break
		}
	}
}

func (f *Font) raiseChanging() {
	for _, handlerIface := range f.changingHandlers {
		handler := handlerIface.(FontChangeHandler)
		handler(f)
	}
}

func (f *Font) raiseChanged() {
	for _, handlerIface := range f.changedHandlers {
		handler := handlerIface.(FontChangeHandler)
		handler(f)
	}
}
