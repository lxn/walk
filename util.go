// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"bytes"
	"math"
	"math/big"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/xlplbo/win"
)

var (
	decimalSepB      byte
	decimalSepUint16 uint16
	decimalSepS      string
	groupSepB        byte
	groupSepUint16   uint16
	groupSepS        string
)

func init() {
	AppendToWalkInit(func() {
		var buf [4]uint16

		win.GetLocaleInfo(win.LOCALE_USER_DEFAULT, win.LOCALE_SDECIMAL, &buf[0], int32(len(buf)))
		decimalSepB = byte(buf[0])
		decimalSepS = syscall.UTF16ToString(buf[0:1])
		decimalSepUint16 = buf[0]

		win.GetLocaleInfo(win.LOCALE_USER_DEFAULT, win.LOCALE_STHOUSAND, &buf[0], int32(len(buf)))
		groupSepB = byte(buf[0])
		groupSepS = syscall.UTF16ToString(buf[0:1])
		groupSepUint16 = buf[0]
	})
}

func maxi(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func mini(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

func uint16IndexUint16(s []uint16, v uint16) int {
	for i, u := range s {
		if u == v {
			return i
		}
	}

	return -1
}

func uint16ContainsUint16(s []uint16, v uint16) bool {
	return uint16IndexUint16(s, v) != -1
}

func uint16CountUint16(s []uint16, v uint16) int {
	var count int

	for _, u := range s {
		if u == v {
			count++
		}
	}

	return count
}

func uint16RemoveUint16(s []uint16, v uint16) []uint16 {
	count := uint16CountUint16(s, v)
	if count == 0 {
		return s
	}

	ret := make([]uint16, 0, len(s)-count)

	for _, u := range s {
		if u != v {
			ret = append(ret, u)
		}
	}

	return ret
}

func assertFloat64Or(value interface{}, defaultValue float64) float64 {
	if f, ok := value.(float64); ok {
		return f
	}

	return defaultValue
}

func assertIntOr(value interface{}, defaultValue int) int {
	if n, ok := value.(int); ok {
		return n
	}

	return defaultValue
}

func assertStringOr(value interface{}, defaultValue string) string {
	if s, ok := value.(string); ok {
		return s
	}

	return defaultValue
}

func assertTimeOr(value interface{}, defaultValue time.Time) time.Time {
	if t, ok := value.(time.Time); ok {
		return t
	}

	return defaultValue
}

func ParseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)

	t := FormatFloatGrouped(1000, 2)

	replaceSep := func(new string, index func(string, func(rune) bool) int) {
		i := index(t, func(r rune) bool {
			return r < '0' || r > '9'
		})

		var sep string
		if i > -1 {
			sep = string(t[i])
		}
		if sep != "" {
			s = strings.Replace(s, string(sep), new, -1)
		}
	}

	replaceSep("", strings.IndexFunc)
	replaceSep(".", strings.LastIndexFunc)

	return strconv.ParseFloat(s, 64)
}

func FormatFloat(f float64, prec int) string {
	return formatFloatString(strconv.FormatFloat(f, 'f', prec, 64), prec, false)
}

func FormatFloatGrouped(f float64, prec int) string {
	return formatFloatString(strconv.FormatFloat(f, 'f', maxi(1, prec), 64), prec, true)
}

func formatBigRat(r *big.Rat, prec int) string {
	return formatFloatString(r.FloatString(prec), prec, false)
}

func formatBigRatGrouped(r *big.Rat, prec int) string {
	return formatFloatString(r.FloatString(prec), prec, true)
}

func formatFloatString(s string, prec int, grouped bool) string {
	switch s {
	case "NaN", "-Inf", "+Inf":
		return s
	}

	s = strings.Replace(s, ".", decimalSepS, 1)
	if !grouped {
		return s
	}

	b := new(bytes.Buffer)

	var firstDigit int
	if len(s) > 0 && s[0] == '-' {
		firstDigit = 1
		b.WriteByte('-')
		s = s[1:]
	}

	intLen := len(s) - maxi(1, prec) - 1

	n := intLen % 3
	if n != 0 {
		b.WriteString(s[:n])
	}
	for i := n; i < intLen; i += 3 {
		if b.Len() > firstDigit {
			b.WriteByte(groupSepB)
		}
		b.WriteString(s[i : i+3])
	}

	b.WriteString(s[intLen:])

	s = b.String()

	if prec == 0 {
		s = s[:len(s)-2]
	}

	return s
}

func applyEnabledToDescendants(window Window, enabled bool) {
	wb := window.AsWindowBase()
	wb.applyEnabled(enabled)

	walkDescendants(window, func(w Window) bool {
		if w.Handle() == wb.hWnd {
			return true
		}

		if enabled && !w.AsWindowBase().enabled {
			return false
		}

		w.(applyEnableder).applyEnabled(enabled)

		return true
	})
}

var seenInApplyFontToDescendantsDuringDPIChange map[*WindowBase]bool

func applyFontToDescendants(window Window, font *Font) {
	wb := window.AsWindowBase()
	wb.applyFont(font)

	walkDescendants(window, func(w Window) bool {
		if w.Handle() == wb.hWnd {
			return true
		}

		if w.AsWindowBase().font != nil {
			return false
		}

		if seenInApplyFontToDescendantsDuringDPIChange != nil {
			wb := w.AsWindowBase()
			if seenInApplyFontToDescendantsDuringDPIChange[wb] {
				return true
			}
			seenInApplyFontToDescendantsDuringDPIChange[wb] = true
		}

		w.(applyFonter).applyFont(font)

		return true
	})
}

func applySysColorsToDescendants(window Window) {
	wb := window.AsWindowBase()
	wb.ApplySysColors()

	walkDescendants(window, func(w Window) bool {
		if w.Handle() == wb.hWnd {
			return true
		}

		w.(ApplySysColorser).ApplySysColors()

		return true
	})
}

var seenInApplyDPIToDescendantsDuringDPIChange map[*WindowBase]bool

func applyDPIToDescendants(window Window, dpi int) {
	wb := window.AsWindowBase()
	wb.ApplyDPI(dpi)

	walkDescendants(window, func(w Window) bool {
		if w.Handle() == wb.hWnd {
			return true
		}

		if seenInApplyDPIToDescendantsDuringDPIChange != nil {
			wb := w.AsWindowBase()
			if seenInApplyDPIToDescendantsDuringDPIChange[wb] {
				return true
			}
			seenInApplyDPIToDescendantsDuringDPIChange[wb] = true
		}

		w.(ApplyDPIer).ApplyDPI(dpi)

		return true
	})
}

func walkDescendants(window Window, f func(w Window) bool) {
	window = window.AsWindowBase().window

	if window == nil || !f(window) {
		return
	}

	var children []*WidgetBase

	switch w := window.(type) {
	case *NumberEdit:
		if w.edit != nil {
			children = append(children, w.edit.AsWidgetBase())
		}

	case *TabWidget:
		for _, p := range w.Pages().items {
			children = append(children, p.AsWidgetBase())
		}

	case Container:
		if c := w.Children(); c != nil {
			children = c.items
		} else {
			children = nil
		}
	}

	for _, wb := range children {
		walkDescendants(wb.window.(Widget), f)
	}
}

func less(a, b interface{}, order SortOrder) bool {
	if _, ok := a.(error); ok {
		_, bIsErr := b.(error)

		return order == SortAscending == !bIsErr
	}
	if _, ok := b.(error); ok {
		return order == SortDescending
	}

	if a == nil {
		return order == SortAscending == (b != nil)
	}
	if b == nil {
		return order == SortDescending
	}

	switch av := a.(type) {
	case string:
		if bv, ok := b.(string); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case int:
		if bv, ok := b.(int); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case float64:
		if bv, ok := b.(float64); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case float32:
		if bv, ok := b.(float32); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case int64:
		if bv, ok := b.(int64); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case int32:
		if bv, ok := b.(int32); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case int16:
		if bv, ok := b.(int16); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case int8:
		if bv, ok := b.(int8); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case uint:
		if bv, ok := b.(uint); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case uint64:
		if bv, ok := b.(uint64); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case uint32:
		if bv, ok := b.(uint32); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case uint16:
		if bv, ok := b.(uint16); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case uint8:
		if bv, ok := b.(uint8); ok {
			if order == SortAscending {
				return av < bv
			} else {
				return bv < av
			}
		}

	case time.Time:
		if bv, ok := b.(time.Time); ok {
			if order == SortAscending {
				return av.Before(bv)
			} else {
				return bv.Before(av)
			}
		}

	case bool:
		if bv, ok := b.(bool); ok {
			if order == SortAscending {
				return !av && bv
			} else {
				return !bv && av
			}
		}
	}

	return false
}

func dpiForHDC(hdc win.HDC) int {
	if hwnd := win.WindowFromDC(hdc); hwnd != 0 {
		return int(win.GetDpiForWindow(hwnd))
	}

	return int(win.GetDeviceCaps(hdc, win.LOGPIXELSX))
}

// IntFrom96DPI converts from 1/96" units to native pixels.
func IntFrom96DPI(value, dpi int) int {
	return scaleInt(value, float64(dpi)/96.0)
}

// IntTo96DPI converts from native pixels to 1/96" units.
func IntTo96DPI(value, dpi int) int {
	return scaleInt(value, 96.0/float64(dpi))
}

func scaleInt(value int, scale float64) int {
	return int(math.Round(float64(value) * scale))
}

// MarginsFrom96DPI converts from 1/96" units to native pixels.
func MarginsFrom96DPI(value Margins, dpi int) Margins {
	return scaleMargins(value, float64(dpi)/96.0)
}

// MarginsTo96DPI converts from native pixels to 1/96" units.
func MarginsTo96DPI(value Margins, dpi int) Margins {
	return scaleMargins(value, 96.0/float64(dpi))
}

func scaleMargins(value Margins, scale float64) Margins {
	return Margins{
		HNear: scaleInt(value.HNear, scale),
		VNear: scaleInt(value.VNear, scale),
		HFar:  scaleInt(value.HFar, scale),
		VFar:  scaleInt(value.VFar, scale),
	}
}

// PointFrom96DPI converts from 1/96" units to native pixels.
func PointFrom96DPI(value Point, dpi int) Point {
	return scalePoint(value, float64(dpi)/96.0)
}

// PointTo96DPI converts from native pixels to 1/96" units.
func PointTo96DPI(value Point, dpi int) Point {
	return scalePoint(value, 96.0/float64(dpi))
}

func scalePoint(value Point, scale float64) Point {
	return Point{
		X: scaleInt(value.X, scale),
		Y: scaleInt(value.Y, scale),
	}
}

// RectangleFrom96DPI converts from 1/96" units to native pixels.
func RectangleFrom96DPI(value Rectangle, dpi int) Rectangle {
	return scaleRectangle(value, float64(dpi)/96.0)
}

// RectangleTo96DPI converts from native pixels to 1/96" units.
func RectangleTo96DPI(value Rectangle, dpi int) Rectangle {
	return scaleRectangle(value, 96.0/float64(dpi))
}

func scaleRectangle(value Rectangle, scale float64) Rectangle {
	return Rectangle{
		X:      scaleInt(value.X, scale),
		Y:      scaleInt(value.Y, scale),
		Width:  scaleInt(value.Width, scale),
		Height: scaleInt(value.Height, scale),
	}
}

// SizeFrom96DPI converts from 1/96" units to native pixels.
func SizeFrom96DPI(value Size, dpi int) Size {
	return scaleSize(value, float64(dpi)/96.0)
}

// SizeTo96DPI converts from native pixels to 1/96" units.
func SizeTo96DPI(value Size, dpi int) Size {
	return scaleSize(value, 96.0/float64(dpi))
}

func scaleSize(value Size, scale float64) Size {
	return Size{
		Width:  scaleInt(value.Width, scale),
		Height: scaleInt(value.Height, scale),
	}
}
