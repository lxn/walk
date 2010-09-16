// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printing

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/comdlg32"
	. "walk/winapi/kernel32"
	. "walk/winapi/gdi32"
	. "walk/winapi/winspool"
)

type optBool int

const (
	obUnspecified optBool = iota
	obFalse
	obTrue
)

func toOptBool(b bool) optBool {
	if b {
		return obTrue
	}

	return obFalse
}

type devModeField int

const (
	dmOrientation devModeField = iota
	dmPaperSize
	dmPaperLength
	dmPaperWidth
	dmCopies
	dmDefaultSource
	dmPrintQuality
	dmColor
	dmDuplex
	dmResolution
	dmTTOption
	dmCollate
)

type Duplex int16

const (
	duplexUnspecified Duplex = iota
	DuplexSimplex
	DuplexVertical
	DuplexHorizontal
)

type PrinterInfo struct {
	driverName  string
	printerName string
	pageInfo    *PageInfo
	outputPort  string
	collate     optBool
	copies      int16
	duplex      Duplex
	fromPage    int
	toPage      int
	maxPage     int
	minPage     int
}

func NewPrinterInfo() *PrinterInfo {
	p := &PrinterInfo{copies: -1, duplex: duplexUnspecified, maxPage: 9999, pageInfo: NewPageInfo()}
	p.pageInfo.SetPrinterInfo(p)

	return p
}

func (p *PrinterInfo) printerNameFallbackPtr() *uint16 {
	if p.printerName != "" {
		return syscall.StringToUTF16Ptr(p.printerName)
	}

	return defaultPrinterNamePtr()
}

func (p *PrinterInfo) printerNamePtr() *uint16 {
	return p.printerNameFallbackPtr()
}

func (p *PrinterInfo) printerNameFallback() string {
	if p.printerName != "" {
		return p.printerName
	}

	return defaultPrinterName()
}

func (p *PrinterInfo) createDCFromHDevMode(hDevMode HGLOBAL) HDC {
	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	hdc := CreateDC(syscall.StringToUTF16Ptr(p.driverName), p.printerNameFallbackPtr(), nil, devMode)
	if hdc == 0 {
		panic("CreateDC failed")
	}

	return hdc
}

func (p *PrinterInfo) createDC() HDC {
	hDevMode := p.hDevMode()
	defer GlobalUnlock(hDevMode)

	return p.createDCFromHDevMode(hDevMode)
}

func (p *PrinterInfo) createICFromHDevMode(hDevMode HGLOBAL) HDC {
	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	hic := CreateIC(syscall.StringToUTF16Ptr(p.driverName), p.printerNameFallbackPtr(), nil, devMode)
	if hic == 0 {
		panic("CreateIC failed")
	}

	return hic
}

func (p *PrinterInfo) createIC() HDC {
	hDevMode := p.hDevMode()
	defer GlobalUnlock(hDevMode)

	return p.createICFromHDevMode(hDevMode)
}

func (p *PrinterInfo) devCapsPrinter(capability uint16, pOutput *uint16, defaultValue uint, printerName *uint16) uint {
	val := DeviceCapabilities(printerName, syscall.StringToUTF16Ptr(p.outputPort), capability, pOutput, nil)

	if val == ^uint(0) { // ^0 == -1
		return defaultValue
	}

	return val
}

func (p *PrinterInfo) devCaps(capability uint16, pOutput *uint16, defaultValue uint) uint {
	return p.devCapsPrinter(capability, pOutput, defaultValue, p.printerNamePtr())
}

func (p *PrinterInfo) getDeviceCaps(capability int) int {
	hic := p.createIC()
	defer DeleteDC(hic)

	return GetDeviceCaps(hic, capability)
}

func (p *PrinterInfo) hDevMode() HGLOBAL {
	printerNamePtr := p.printerNamePtr()

	bufSize := DocumentProperties(0, 0, printerNamePtr, nil, nil, 0)
	if bufSize <= 1 {
		panic("DocumentProperties failed")
	}

	succeeded := false

	hDevMode := GlobalAlloc(GMEM_MOVEABLE, uintptr(bufSize))
	defer func() {
		if !succeeded {
			GlobalFree(hDevMode)
		}
	}()

	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	ret := DocumentProperties(0, 0, printerNamePtr, devMode, nil, DM_OUT_BUFFER)
	if ret < 0 {
		panic("DocumentProperties failed")
	}

	if p.collate != obUnspecified {
		devMode.DmCollate = int16(p.collate)
	}
	if p.copies > -1 {
		devMode.DmCopies = p.copies
	}
	if p.duplex != duplexUnspecified {
		devMode.DmDuplex = int16(p.duplex)
	}

	succeeded = true

	return hDevMode
}

func (p *PrinterInfo) hDevNames() HGLOBAL {
	driverName := syscall.StringToUTF16(p.driverName)
	printerName := syscall.StringToUTF16(UTF16PtrToString(p.printerNamePtr()))
	outputPort := syscall.StringToUTF16(p.outputPort)

	driverLen := len(driverName) + 1
	printerLen := len(printerName) + 1
	portLen := len(outputPort) + 1

	var dn DEVNAMES
	hDevNames := GlobalAlloc(GHND, uintptr(unsafe.Sizeof(dn)+(driverLen+printerLen+portLen)*2))
	devNames := (*DEVNAMES)(GlobalLock(hDevNames))
	defer GlobalUnlock(hDevNames)

	devNames.WDriverOffset = uint16(unsafe.Sizeof(*devNames) / 2)
	devNames.WDeviceOffset = devNames.WDriverOffset + uint16(driverLen)
	devNames.WOutputOffset = devNames.WDeviceOffset + uint16(printerLen)

	MoveMemory(unsafe.Pointer(uintptr(unsafe.Pointer(devNames))+uintptr(devNames.WDriverOffset)), unsafe.Pointer(&driverName[0]), uintptr(driverLen*2))
	MoveMemory(unsafe.Pointer(uintptr(unsafe.Pointer(devNames))+uintptr(devNames.WDeviceOffset)), unsafe.Pointer(&printerName[0]), uintptr(printerLen*2))
	MoveMemory(unsafe.Pointer(uintptr(unsafe.Pointer(devNames))+uintptr(devNames.WOutputOffset)), unsafe.Pointer(&outputPort[0]), uintptr(portLen*2))

	return hDevNames
}

func (p *PrinterInfo) valueFromHDevMode(field devModeField, defaultValue int, hDevMode HGLOBAL) int {
	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	switch field {
	case dmOrientation:
		return int(devMode.DmOrientation)

	case dmPaperSize:
		return int(devMode.DmPaperSize)

	case dmPaperLength:
		return int(devMode.DmPaperLength)

	case dmPaperWidth:
		return int(devMode.DmPaperWidth)

	case dmCopies:
		return int(devMode.DmCopies)

	case dmDefaultSource:
		return int(devMode.DmDefaultSource)

	case dmPrintQuality:
		return int(devMode.DmPrintQuality)

	case dmColor:
		return int(devMode.DmColor)

	case dmDuplex:
		return int(devMode.DmDuplex)

	case dmResolution:
		return int(devMode.DmYResolution)

	case dmTTOption:
		return int(devMode.DmTTOption)

	case dmCollate:
		return int(devMode.DmCollate)
	}

	return defaultValue
}

func (p *PrinterInfo) valueFromDevMode(field devModeField, defaultValue int) int {
	hDevMode := p.hDevMode()
	defer GlobalFree(hDevMode)

	return p.valueFromHDevMode(field, defaultValue, hDevMode)
}

func (p *PrinterInfo) devNamePtr(devNamesPtr unsafe.Pointer, index int) *uint16 {
	if devNamesPtr == nil || index < 0 || index > 2 {
		panic("invalid argument")
	}

	var offset uint16
	var offsetPtr unsafe.Pointer

	offsetPtr = unsafe.Pointer(uintptr(devNamesPtr) + uintptr(index*unsafe.Sizeof(offset)))

	MoveMemory(unsafe.Pointer(&offset), offsetPtr, uintptr(unsafe.Sizeof(offset)))

	return (*uint16)(unsafe.Pointer(uintptr(devNamesPtr) + uintptr(offset)))
}

func (p *PrinterInfo) devName(devNamesPtr unsafe.Pointer, index int) string {
	return UTF16PtrToString(p.devNamePtr(devNamesPtr, index))
}

/*func (p *PrinterInfo) hDevMode() (val HGLOBAL, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.hDevMode(), nil
}

func (p *PrinterInfo) hDevMode2(pageInfo *PageInfo) (hDevMode HGLOBAL, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)

            if hDevMode != 0 {
                GlobalFree(hDevMode)
                hDevMode = 0
            }
		}
	}()

	hDevMode = p.hDevMode()
	pageInfo.copyToHDevMode(hDevMode)

	return
}*/

func (p *PrinterInfo) initFromHDevMode(hDevMode HGLOBAL) (err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if hDevMode == 0 {
		panic("invalid hDevMode")
	}

	devMode := (*DEVMODE)(GlobalLock(hDevMode))
	defer GlobalUnlock(hDevMode)

	p.collate = toOptBool(devMode.DmCollate == 1)
	p.copies = devMode.DmCopies
	p.duplex = Duplex(devMode.DmDuplex)

	return nil
}

func (p *PrinterInfo) initFromHDevNames(hDevNames HGLOBAL) (err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if hDevNames == 0 {
		panic("invalid hDevNames")
	}

	devNamesPtr := GlobalLock(hDevNames)
	defer GlobalUnlock(hDevNames)

	p.driverName = p.devName(devNamesPtr, 0)
	p.printerName = p.devName(devNamesPtr, 1)
	p.outputPort = p.devName(devNamesPtr, 2)

	return nil
}

func (p *PrinterInfo) Collate() (val bool, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.collate != obUnspecified {
		return p.collate == obTrue, nil
	}

	return p.valueFromDevMode(dmCollate, 0) == 1, nil
}

func (p *PrinterInfo) SetCollate(value bool) {
	p.collate = toOptBool(value)
}

func (p *PrinterInfo) Copies() (val int16, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.copies != -1 {
		return p.copies, nil
	}

	return int16(p.valueFromDevMode(dmCopies, 1)), nil
}

func (p *PrinterInfo) SetCopies(value int16) os.Error {
	if value < 0 {
		return newError("invalid value for copies")
	}

	p.copies = value

	return nil
}

func (p *PrinterInfo) DriverName() string {
	return p.driverName
}

func (p *PrinterInfo) SetDriverName(value string) {
	p.driverName = value
}

func (p *PrinterInfo) Duplex() (val Duplex, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	if p.duplex != duplexUnspecified {
		return p.duplex, nil
	}

	return Duplex(p.valueFromDevMode(dmDuplex, int(DuplexSimplex))), nil
}

func (p *PrinterInfo) SetDuplex(value Duplex) {
	p.duplex = value
}

func (p *PrinterInfo) FromPage() int {
	return p.fromPage
}

func (p *PrinterInfo) SetFromPage(value int) os.Error {
	if value < 0 {
		return newError("invalid value for FromPage")
	}

	p.fromPage = value

	return nil
}

func (p *PrinterInfo) IsPlotter() bool {
	return p.getDeviceCaps(DT_RASPRINTER) == 0
}

func (p *PrinterInfo) IsValid() bool {
	return p.devCaps(DC_COPIES, nil, ^uint(0)) != ^uint(0)
}

func (p *PrinterInfo) MaxCopies() int16 {
	return int16(p.devCaps(DC_COPIES, nil, 1))
}

func (p *PrinterInfo) MaxPage() int {
	return p.maxPage
}

func (p *PrinterInfo) SetMaxPage(value int) os.Error {
	if value < 0 {
		return newError("invalid value for MaxPage")
	}

	p.maxPage = value

	return nil
}

func (p *PrinterInfo) MinPage() int {
	return p.minPage
}

func (p *PrinterInfo) SetMinPage(value int) os.Error {
	if value < 0 {
		return newError("invalid value for MinPage")
	}

	p.minPage = value

	return nil
}

func (p *PrinterInfo) OutputPort() string {
	return p.outputPort
}

func (p *PrinterInfo) PageInfo() *PageInfo {
	return p.pageInfo
}

func (p *PrinterInfo) SetPageInfo(value *PageInfo) {
	p.pageInfo = value
}

func (p *PrinterInfo) paperSizes() []*PaperSize {
	printerNamePtr := p.printerNamePtr()

	count := p.devCapsPrinter(DC_PAPERS, nil, ^uint(0), printerNamePtr)
	if count <= 0 {
		return nil
	}

	names := make([]uint16, 64*count)
	p.devCapsPrinter(DC_PAPERNAMES, &names[0], ^uint(0), printerNamePtr)

	types := make([]uint16, count)
	p.devCapsPrinter(DC_PAPERS, &types[0], ^uint(0), printerNamePtr)

	sizes := make([]SIZE, count)
	p.devCapsPrinter(DC_PAPERSIZE, (*uint16)(unsafe.Pointer(&sizes[0])), ^uint(0), printerNamePtr)

	papers := make([]*PaperSize, count)
	for i := range papers {
		name := syscall.UTF16ToString(names[i*64 : i*64+64])
		typ := PaperSizeType(types[i])
		size := sizes[i]

		papers[i] = &PaperSize{name: name, typ: typ, width: size.CX, height: size.CY}
	}

	return papers
}

func (p *PrinterInfo) PaperSizes() (val []*PaperSize, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.paperSizes(), nil
}

func (p *PrinterInfo) paperSources() []*PaperSource {
	printerNamePtr := p.printerNamePtr()

	count := p.devCapsPrinter(DC_BINS, nil, ^uint(0), printerNamePtr)
	if count <= 0 {
		return nil
	}

	names := make([]uint16, 24*count)
	p.devCapsPrinter(DC_BINNAMES, &names[0], ^uint(0), printerNamePtr)

	types := make([]uint16, count)
	p.devCapsPrinter(DC_BINS, &types[0], ^uint(0), printerNamePtr)

	sources := make([]*PaperSource, count)
	for i := range sources {
		name := syscall.UTF16ToString(names[i*24 : i*24+24])
		typ := PaperSourceType(types[i])

		sources[i] = &PaperSource{name: name, typ: typ}
	}

	return sources
}

func (p *PrinterInfo) PaperSources() (val []*PaperSource, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.paperSources(), nil
}

func (p *PrinterInfo) PrinterName() (val string, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.printerNameFallback(), nil
}

func (p *PrinterInfo) SetPrinterName(value string) {
	p.printerName = value
}

func (p *PrinterInfo) resolutions() []*Resolution {
	printerNamePtr := p.printerNamePtr()

	count := p.devCapsPrinter(DC_ENUMRESOLUTIONS, nil, ^uint(0), printerNamePtr)
	if count <= 0 {
		return nil
	}

	sizes := make([]SIZE, count)
	p.devCapsPrinter(DC_ENUMRESOLUTIONS, (*uint16)(unsafe.Pointer(&sizes[0])), ^uint(0), printerNamePtr)

	resolutions := make([]*Resolution, count+4)

	resolutions[0] = &Resolution{typ: ResHigh, x: -1, y: -1}
	resolutions[1] = &Resolution{typ: ResMedium, x: -1, y: -1}
	resolutions[2] = &Resolution{typ: ResLow, x: -1, y: -1}
	resolutions[3] = &Resolution{typ: ResDraft, x: -1, y: -1}

	for i := 0; i < len(sizes); i++ {
		size := sizes[i]

		resolutions[i+4] = &Resolution{typ: ResCustom, x: size.CX, y: size.CY}
	}

	return resolutions
}

func (p *PrinterInfo) Resolutions() (val []*Resolution, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.resolutions(), nil
}

func (p *PrinterInfo) SupportsColor() (val bool, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.getDeviceCaps(BITSPIXEL) > 1, nil
}

func (p *PrinterInfo) SupportsDuplex() (val bool, err os.Error) {
	defer func() {
		if x := recover(); x != nil {
			err = toError(x)
		}
	}()

	return p.devCaps(DC_DUPLEX, nil, 0) == 1, nil
}

func (p *PrinterInfo) ToPage() int {
	return p.toPage
}

func (p *PrinterInfo) SetToPage(value int) os.Error {
	if value < 0 {
		return newError("invalid value for ToPage")
	}

	p.toPage = value

	return nil
}

/*func (p *PrinterInfo) ShowPageSetupDialog(ownerHandle HWND, _
                                    Optional ByVal Flags As PageSetupFlags, _
                                    Optional ByVal MinMargins As Margins) As Boolean

    Dim usrPageDlg As PAGESETUPDLG_TYPE
    Dim lngResult As Long


    With usrPageDlg
        .lStructSize = Len(usrPageDlg)
        .hwndOwner = hwndOwner
        If Not (Flags And PageSetupFlags.CustomMinMargins) _
                = PageSetupFlags.CustomMinMargins Then Flags = Flags Or PSD_DEFAULTMINMARGINS
        .Flags = PSD_INHUNDREDTHSOFMILLIMETERS Or Flags

        Dim objMargins As Margins
        Set objMargins = DefaultPageSettings.Margins

        If (Flags And PageSetupFlags.CustomMargins) > 0 Then
            With .rtMargin
                .Left = objMargins.Left
                .Top = objMargins.Top
                .Right = objMargins.Right
                .Bottom = objMargins.Bottom
            End With '.rtMargin
        Else
            .Flags = .Flags Or PageSetupFlags.CustomMargins

            With DefaultPageSettings.Margins
                usrPageDlg.rtMargin.Left = .Left
                usrPageDlg.rtMargin.Top = .Top
                usrPageDlg.rtMargin.Right = .Right
                usrPageDlg.rtMargin.Bottom = .Bottom
            End With 'mobjPageSettings.Margins
        End If

        If Not MinMargins Is Nothing Then
            With .rtMinMargin
                .Left = MinMargins.Left
                .Top = MinMargins.Top
                .Right = MinMargins.Right
                .Bottom = MinMargins.Bottom
            End With '.rtMinMargin
        End If

        .HDevMode = hDevMode
        .hDevNames = GetHdevnames

        DefaultPageSettings.mergeHDevMode .HDevMode

        lngResult = PageSetupDlg(usrPageDlg)

        If lngResult <> 0 Then
            initFromHDevMode .HDevMode
            SetHdevnames .hDevNames

            DefaultPageSettings.initFromHDevMode .HDevMode

            Dim lngPtrDevMode As Long

            lngPtrDevMode = GlobalLock(.HDevMode)

            GlobalUnlock .HDevMode

            With DefaultPageSettings.Margins
                .Left = usrPageDlg.rtMargin.Left
                .Top = usrPageDlg.rtMargin.Top
                .Right = usrPageDlg.rtMargin.Right
                .Bottom = usrPageDlg.rtMargin.Bottom
            End With 'mobjPageSettings.Margins

            GlobalFree .HDevMode
            GlobalFree .hDevNames

'            With .rtMargin
'                marginLeftMM = .Left / 100
'                marginTopMM = .Top / 100
'                marginRightMM = .Right / 100
'                marginBottomMM = .Bottom / 100
'            End With '.rtMargin

            ShowPageSetupDialog = True
        Else
            GlobalFree .HDevMode
            GlobalFree .hDevNames

            lngResult = CommDlgExtendedError
        End If
    End With 'usrPageDlg

End Function

Public Function ShowPrintDialog(ByVal hwndOwner As OLE_HANDLE) As Boolean

    On Error Resume Next

    Dim usrPrintDlg As PRINTDLG_TYPE
    Dim lngResult As Long


    With usrPrintDlg
        .lStructSize = Len(usrPrintDlg)
        .hwndOwner = hwndOwner

        .Flags = PD_PRINTSETUP Or PD_USEDEVMODECOPIESANDCOLLATE

        .HDevMode = hDevMode
        .hDevNames = GetHdevnames

        DefaultPageSettings.mergeHDevMode .HDevMode

        lngResult = PrintDlg(usrPrintDlg)

        If lngResult <> 0 Then
            initFromHDevMode .HDevMode
            SetHdevnames .hDevNames

            DefaultPageSettings.initFromHDevMode .HDevMode

            Dim lngPtrDevMode As Long

            lngPtrDevMode = GlobalLock(.HDevMode)

            GlobalUnlock .HDevMode

            GlobalFree .HDevMode
            GlobalFree .hDevNames

            ShowPrintDialog = True
        Else
            GlobalFree .HDevMode
            GlobalFree .hDevNames

            lngResult = CommDlgExtendedError
        End If
    End With 'usrPrintDlg

End Function*/
