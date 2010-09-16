// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printing

import (
	"os"
)

import (
	"walk/drawing"
	. "walk/winapi/gdi32"
	. "walk/winapi/kernel32"
)

type ColorMode int16

const (
	colorUnspecified ColorMode = 0
	ColorMonochrome  ColorMode = DMCOLOR_MONOCHROME
	ColorColor       ColorMode = DMCOLOR_COLOR
)

type Orientation int16

const (
	orientUnspecified Orientation = 0
	OrientPortrait    Orientation = DMORIENT_PORTRAIT
	OrientLandscape   Orientation = DMORIENT_LANDSCAPE
)

type Margins struct {
	Left, Top, Right, Bottom int
}

type PageInfo struct {
	printerInfo *PrinterInfo
	paperSize   *PaperSize
	paperSource *PaperSource
	resolution  *Resolution
	colorMode   ColorMode
	orientation Orientation
	margins     *Margins
}

func NewPageInfo() *PageInfo {
	return &PageInfo{margins: &Margins{}}
}

func (p *PageInfo) Bounds() (val *drawing.Rectangle, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	hDevMode := p.printerInfo.hDevMode()
	defer GlobalFree(hDevMode)

	return p.boundsFromHDevMode(hDevMode), nil
}

func (p *PageInfo) ColorMode() (val ColorMode, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.colorMode != colorUnspecified {
		return p.colorMode, nil
	}

	return ColorMode(p.printerInfo.valueFromDevMode(dmColor, 1)), nil
}

func (p *PageInfo) SetColorMode(value ColorMode) {
	p.colorMode = value
}

func (p *PageInfo) Orientation() (val Orientation, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.orientation == orientUnspecified {
		return p.orientation, nil
	}

	return Orientation(p.printerInfo.valueFromDevMode(dmOrientation, DMORIENT_PORTRAIT)), nil
}

func (p *PageInfo) SetOrientation(value Orientation) (err os.Error) {
	orientation, err := p.Orientation()
	if err != nil {
		return
	}

	if value != orientation {
		m := p.margins

		l := m.Left
		t := m.Top
		r := m.Right
		b := m.Bottom

		if value == OrientLandscape {
			m.Left = t
			m.Top = r
			m.Right = b
			m.Bottom = l
		} else {
			m.Left = b
			m.Top = l
			m.Right = t
			m.Bottom = r
		}
	}

	p.orientation = value

	return
}

func (p *PageInfo) Margins() *Margins {
	return p.margins
}

func (p *PageInfo) SetMargins(value *Margins) {
	p.margins = value
}

func (p *PageInfo) PaperSize() (val *PaperSize, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.paperSizeFallbackHDevMode(0), nil
}

func (p *PageInfo) SetPaperSize(value *PaperSize) {
	p.paperSize = value
}

func (p *PageInfo) PaperSource() (val *PaperSource, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.paperSource != nil {
		return p.paperSource, nil
	}

	hDevMode := p.printerInfo.hDevMode()
	defer GlobalFree(hDevMode)

	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	return p.paperSourceFromDevMode(devMode), nil
}

func (p *PageInfo) SetPaperSource(value *PaperSource) {
	p.paperSource = value
}

func (p *PageInfo) Resolution() (val *Resolution, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.resolution != nil {
		return p.resolution, nil
	}

	hDevMode := p.printerInfo.hDevMode()
	defer GlobalFree(hDevMode)

	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	return p.resolutionFromDevMode(devMode), nil
}

func (p *PageInfo) SetResolution(value *Resolution) {
	p.resolution = value
}

func (p *PageInfo) PrinterInfo() *PrinterInfo {
	return p.printerInfo
}

func (p *PageInfo) SetPrinterInfo(value *PrinterInfo) {
	if value == nil {
		value = NewPrinterInfo()
	}

	p.printerInfo = value
}

func (p *PageInfo) SetCustomPaperSize(name string, width, height int) {
	p.paperSize = &PaperSize{typ: PaperCustom, name: name, width: width, height: height}
}

func (p *PageInfo) SetPaperSizeFromType(typ PaperSizeType) os.Error {
	if typ == PaperCustom {
		return newError("call SetCustomPaperSize to set a custom paper size")
	}

	sizes, err := p.printerInfo.PaperSizes()
	if err != nil {
		return err
	}

	for _, size := range sizes {
		if size.typ == typ {
			p.paperSize = size
			return nil
		}
	}

	return newError("invalid paper size")
}

func (p *PageInfo) SetCustomResolution(x, y int) {
	p.resolution = &Resolution{typ: ResCustom, x: x, y: y}
}

func (p *PageInfo) SetResolutionFromType(typ ResolutionType) os.Error {
	if typ == ResCustom {
		return newError("call SetCustomResolution to set a custom resolution")
	}

	resolutions, err := p.printerInfo.Resolutions()
	if err != nil {
		return err
	}

	for _, res := range resolutions {
		if res.typ == typ {
			p.resolution = res
			return nil
		}
	}

	return newError("invalid resolution")
}

func (p *PageInfo) SetPaperSourceFromType(typ PaperSourceType) os.Error {
	sources, err := p.printerInfo.PaperSources()
	if err != nil {
		return err
	}

	for _, source := range sources {
		if source.typ == typ {
			p.paperSource = source
			return nil
		}
	}

	return newError("invalid paper source")
}

func (p *PageInfo) mergeHDevMode(hDevMode HGLOBAL) {
	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	if p.colorMode != colorUnspecified {
		devMode.DmColor = int16(p.colorMode)
	}

	if p.orientation != orientUnspecified {
		devMode.DmOrientation = int16(p.orientation)
	}

	if p.paperSize != nil {
		devMode.DmPaperSize = int16(p.paperSize.typ)
		devMode.DmPaperLength = int16(p.paperSize.height)
		devMode.DmPaperWidth = int16(p.paperSize.width)
	}

	if p.paperSource != nil {
		devMode.DmDefaultSource = int16(p.paperSource.typ)
	}

	if p.resolution != nil {
		if p.resolution.typ == ResCustom {
			devMode.DmPrintQuality = int16(p.resolution.x)
			devMode.DmYResolution = int16(p.resolution.y)
		} else {
			devMode.DmPrintQuality = int16(p.resolution.typ)
		}
	}
}

func (p *PageInfo) boundsFromHDevMode(hDevMode HGLOBAL) *drawing.Rectangle {
	size := p.paperSizeFallbackHDevMode(hDevMode)

	if p.orientationFallbackHDevMode(hDevMode) == OrientLandscape {
		return &drawing.Rectangle{Width: size.height, Height: size.width}
	}

	return &drawing.Rectangle{Width: size.width, Height: size.height}
}

func (p *PageInfo) orientationFallbackHDevMode(hDevMode HGLOBAL) Orientation {
	if p.orientation != orientUnspecified {
		return p.orientation
	}

	return Orientation(p.printerInfo.valueFromHDevMode(dmOrientation, DMORIENT_PORTRAIT, hDevMode))
}

func (p *PageInfo) paperSizeFromDevMode(devMode *DEVMODE) *PaperSize {
	sizes := p.printerInfo.paperSizes()

	for _, size := range sizes {
		if int16(size.typ) == devMode.DmPaperSize {
			return size
		}
	}

	return &PaperSize{
		typ:    PaperCustom,
		name:   "Custom",
		width:  int(devMode.DmPaperWidth),
		height: int(devMode.DmPaperLength),
	}
}

func (p *PageInfo) paperSizeFallbackHDevMode(hDevMode HGLOBAL) *PaperSize {
	if p.paperSize != nil {
		return p.paperSize
	}

	if hDevMode == 0 {
		hDevMode := p.printerInfo.hDevMode()
		defer GlobalFree(hDevMode)
	}

	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	return p.paperSizeFromDevMode(devMode)
}

func (p *PageInfo) paperSourceFromDevMode(devMode *DEVMODE) *PaperSource {
	sources := p.printerInfo.paperSources()

	for _, source := range sources {
		if int16(source.typ) == devMode.DmDefaultSource {
			return source
		}
	}

	return &PaperSource{typ: PaperSourceType(devMode.DmDefaultSource), name: "Unknown"}
}

func (p *PageInfo) resolutionFromDevMode(devMode *DEVMODE) *Resolution {
	resolutions := p.printerInfo.resolutions()

	for _, res := range resolutions {
		if devMode.DmPrintQuality >= 0 {
			if res.x == int(devMode.DmPrintQuality) &&
				res.y == int(devMode.DmYResolution) {
				return res
			}
		}

		if int16(res.typ) == devMode.DmPrintQuality {
			return res
		}
	}

	return &Resolution{
		typ: ResCustom,
		x:   int(devMode.DmPrintQuality),
		y:   int(devMode.DmYResolution),
	}
}

func (p *PageInfo) initFromHDevMode(hDevMode HGLOBAL) (err os.Error) {
	if hDevMode == 0 {
		panic("invalid hDevMode")
	}

	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	p.colorMode = ColorMode(devMode.DmColor)
	p.orientation = Orientation(devMode.DmOrientation)
	p.paperSize = p.paperSizeFromDevMode(devMode)
	p.paperSource = p.paperSourceFromDevMode(devMode)
	p.resolution = p.resolutionFromDevMode(devMode)

	return nil
}
