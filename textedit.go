// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type TextEdit struct {
	WidgetBase
	readOnlyChangedPublisher EventPublisher
	textChangedPublisher     EventPublisher
	textColor                Color
	compactHeight            bool
	margins                  Size // in native pixels
	lastHeight               int
	origWordbreakProcPtr     uintptr
}

func NewTextEdit(parent Container) (*TextEdit, error) {
	return NewTextEditWithStyle(parent, 0)
}

func NewTextEditWithStyle(parent Container, style uint32) (*TextEdit, error) {
	te := new(TextEdit)

	if err := InitWidget(
		te,
		parent,
		"EDIT",
		win.WS_TABSTOP|win.WS_VISIBLE|win.ES_MULTILINE|win.ES_WANTRETURN|style,
		win.WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	te.origWordbreakProcPtr = te.SendMessage(win.EM_GETWORDBREAKPROC, 0, 0)

	te.GraphicsEffects().Add(InteractionEffect)
	te.GraphicsEffects().Add(FocusEffect)

	te.MustRegisterProperty("ReadOnly", NewProperty(
		func() interface{} {
			return te.ReadOnly()
		},
		func(v interface{}) error {
			return te.SetReadOnly(v.(bool))
		},
		te.readOnlyChangedPublisher.Event()))

	te.MustRegisterProperty("Text", NewProperty(
		func() interface{} {
			return te.Text()
		},
		func(v interface{}) error {
			return te.SetText(assertStringOr(v, ""))
		},
		te.textChangedPublisher.Event()))

	return te, nil
}

func (te *TextEdit) applyFont(font *Font) {
	te.WidgetBase.applyFont(font)

	te.updateMargins()
}

func (te *TextEdit) updateMargins() {
	// 56 works at least from 96 to 192 DPI, so until a better solution comes up, this is it.
	defaultSize := te.dialogBaseUnitsToPixels(Size{56, 12})

	var rc win.RECT
	te.SendMessage(win.EM_GETRECT, 0, uintptr(unsafe.Pointer(&rc)))

	if te.hasExtendedStyleBits(win.WS_EX_CLIENTEDGE) {
		width := te.WidthPixels()
		if width == 0 {
			width = defaultSize.Width
		}
		te.margins.Width = width - int(rc.Right-rc.Left)
	} else {
		te.margins.Width = int(rc.Left) * 2
	}

	lineHeight := te.calculateTextSizeImpl("gM").Height
	te.margins.Height = defaultSize.Height - lineHeight
}

var drawTextCompatibleEditWordbreakProcPtr uintptr

func init() {
	AppendToWalkInit(func() {
		drawTextCompatibleEditWordbreakProcPtr = syscall.NewCallback(drawTextCompatibleEditWordbreakProc)
	})
}

func drawTextCompatibleEditWordbreakProc(lpch *uint16, ichCurrent, cch, code uintptr) uintptr {
	switch code {
	case win.WB_LEFT:
		for i := int(ichCurrent); i >= 0; i-- {
			if *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(lpch)) + uintptr(i)*2)) == 32 {
				return uintptr(i)
			}
		}

	case win.WB_RIGHT:
		for i := int(ichCurrent); i < int(cch); i++ {
			if *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(lpch)) + uintptr(i)*2)) == 32 {
				return uintptr(i)
			}
		}

	case win.WB_ISDELIMITER:
		if *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(lpch)) + ichCurrent*2)) == 32 {
			return 1
		}
	}

	return 0
}

func (te *TextEdit) Text() string {
	return te.text()
}

func (te *TextEdit) TextLength() int {
	return int(te.SendMessage(win.WM_GETTEXTLENGTH, 0, 0))
}

func (te *TextEdit) SetText(text string) (err error) {
	if text == te.Text() {
		return nil
	}

	var oldLineCount int
	if te.compactHeight {
		oldLineCount = int(te.SendMessage(win.EM_GETLINECOUNT, 0, 0))
	}
	err = te.setText(text)
	if te.compactHeight {
		if newLineCount := int(te.SendMessage(win.EM_GETLINECOUNT, 0, 0)); newLineCount != oldLineCount {
			te.RequestLayout()
		}
	}
	te.textChangedPublisher.Publish()
	return
}

func (te *TextEdit) CompactHeight() bool {
	return te.compactHeight
}

func (te *TextEdit) SetCompactHeight(enabled bool) {
	if enabled == te.compactHeight {
		return
	}

	te.compactHeight = enabled

	var ptr uintptr
	if enabled {
		te.updateMargins()
		ptr = drawTextCompatibleEditWordbreakProcPtr
	} else {
		ptr = te.origWordbreakProcPtr
	}
	te.SendMessage(win.EM_SETWORDBREAKPROC, 0, ptr)

	te.RequestLayout()
}

func (te *TextEdit) TextAlignment() Alignment1D {
	switch win.GetWindowLong(te.hWnd, win.GWL_STYLE) & (win.ES_LEFT | win.ES_CENTER | win.ES_RIGHT) {
	case win.ES_CENTER:
		return AlignCenter

	case win.ES_RIGHT:
		return AlignFar
	}

	return AlignNear
}

func (te *TextEdit) SetTextAlignment(alignment Alignment1D) error {
	if alignment == AlignDefault {
		alignment = AlignNear
	}

	var bit uint32

	switch alignment {
	case AlignCenter:
		bit = win.ES_CENTER

	case AlignFar:
		bit = win.ES_RIGHT

	default:
		bit = win.ES_LEFT
	}

	return te.setAndClearStyleBits(bit, win.ES_LEFT|win.ES_CENTER|win.ES_RIGHT)
}

func (te *TextEdit) MaxLength() int {
	return int(te.SendMessage(win.EM_GETLIMITTEXT, 0, 0))
}

func (te *TextEdit) SetMaxLength(value int) {
	te.SendMessage(win.EM_SETLIMITTEXT, uintptr(value), 0)
}

func (te *TextEdit) ScrollToCaret() {
	te.SendMessage(win.EM_SCROLLCARET, 0, 0)
}

func (te *TextEdit) TextSelection() (start, end int) {
	te.SendMessage(win.EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (te *TextEdit) SetTextSelection(start, end int) {
	te.SendMessage(win.EM_SETSEL, uintptr(start), uintptr(end))
}

func (te *TextEdit) ReplaceSelectedText(text string, canUndo bool) {
	te.SendMessage(win.EM_REPLACESEL,
		uintptr(win.BoolToBOOL(canUndo)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))))
}

func (te *TextEdit) AppendText(value string) {
	s, e := te.TextSelection()
	l := te.TextLength()
	te.SetTextSelection(l, l)
	te.ReplaceSelectedText(value, false)
	te.SetTextSelection(s, e)
}

func (te *TextEdit) ReadOnly() bool {
	return te.hasStyleBits(win.ES_READONLY)
}

func (te *TextEdit) SetReadOnly(readOnly bool) error {
	if 0 == te.SendMessage(win.EM_SETREADONLY, uintptr(win.BoolToBOOL(readOnly)), 0) {
		return newError("SendMessage(EM_SETREADONLY)")
	}

	te.readOnlyChangedPublisher.Publish()

	return nil
}

func (te *TextEdit) TextChanged() *Event {
	return te.textChangedPublisher.Event()
}

func (te *TextEdit) TextColor() Color {
	return te.textColor
}

func (te *TextEdit) SetTextColor(c Color) {
	te.textColor = c

	te.Invalidate()
}

// ContextMenuLocation returns carret position in screen coordinates in native pixels.
func (te *TextEdit) ContextMenuLocation() Point {
	idx := int(te.SendMessage(win.EM_GETCARETINDEX, 0, 0))
	if idx < 0 {
		start, end := te.TextSelection()
		idx = (start + end) / 2
	}
	res := uint32(te.SendMessage(win.EM_POSFROMCHAR, uintptr(idx), 0))
	pt := win.POINT{int32(win.LOWORD(res)), int32(win.HIWORD(res))}
	windowTrimToClientBounds(te.hWnd, &pt)
	return pointPixelsFromPOINT(pt)
}

func (*TextEdit) NeedsWmSize() bool {
	return true
}

func (te *TextEdit) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		switch win.HIWORD(uint32(wParam)) {
		case win.EN_CHANGE:
			if te.compactHeight {
				if createLayoutItemForWidget(te).(MinSizer).MinSize().Height != te.HeightPixels() {
					te.RequestLayout()
				}
			}
			te.textChangedPublisher.Publish()
		}

	case win.WM_GETDLGCODE:
		if wParam == win.VK_RETURN {
			return win.DLGC_WANTALLKEYS
		}

		return win.DLGC_HASSETSEL | win.DLGC_WANTARROWS | win.DLGC_WANTCHARS

	case win.WM_KEYDOWN:
		if Key(wParam) == KeyA && ControlDown() {
			te.SetTextSelection(0, -1)
		}
	}

	return te.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

func (te *TextEdit) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	if te.margins.Width <= 0 {
		te.updateMargins()
	}

	return &textEditLayoutItem{
		width2Height:            make(map[int]int),
		compactHeight:           te.compactHeight,
		margins:                 te.margins,
		text:                    te.Text(),
		font:                    te.Font(),
		minWidth:                te.calculateTextSizeImpl("W").Width,
		nonCompactHeightMinSize: te.dialogBaseUnitsToPixels(Size{20, 12}),
	}
}

type textEditLayoutItem struct {
	LayoutItemBase
	mutex                   sync.Mutex
	width2Height            map[int]int // in native pixels
	nonCompactHeightMinSize Size        // in native pixels
	margins                 Size        // in native pixels
	text                    string
	font                    *Font
	minWidth                int // in native pixels
	compactHeight           bool
}

func (li *textEditLayoutItem) LayoutFlags() LayoutFlags {
	flags := ShrinkableHorz | GrowableHorz | GreedyHorz
	if !li.compactHeight {
		flags |= GreedyVert | GrowableVert | ShrinkableVert
	}
	return flags
}

func (li *textEditLayoutItem) IdealSize() Size {
	if li.compactHeight {
		return li.MinSize()
	} else {
		return SizeFrom96DPI(Size{100, 100}, li.ctx.dpi)
	}
}

func (li *textEditLayoutItem) MinSize() Size {
	if li.compactHeight {
		width := IntFrom96DPI(100, li.ctx.dpi)
		return Size{width, li.HeightForWidth(width)}
	} else {
		return li.nonCompactHeightMinSize
	}
}

func (li *textEditLayoutItem) HasHeightForWidth() bool {
	return li.compactHeight
}

func (li *textEditLayoutItem) HeightForWidth(width int) int {
	li.mutex.Lock()
	defer li.mutex.Unlock()

	if height, ok := li.width2Height[width]; ok {
		return height
	}

	size := calculateTextSize(li.text, li.font, li.ctx.dpi, width-li.margins.Width, li.handle)
	size.Height += li.margins.Height
	size.Height = maxi(size.Height, li.nonCompactHeightMinSize.Height)

	li.width2Height[width] = size.Height

	return size.Height
}
