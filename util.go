// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"math/big"
	"strconv"
	"strings"
	"syscall"
)

import (
	. "github.com/lxn/go-winapi"
)

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

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)

	t, _ := formatFloat(1000, 2)

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

func formatFloat(f float64, prec int) (string, error) {
	return formatFloatString(strconv.FormatFloat(f, 'f', prec, 64), prec)
}

func formatRat(r *big.Rat, prec int) (string, error) {
	return formatFloatString(r.FloatString(prec), prec)
}

func formatFloatString(s string, prec int) (string, error) {
	// FIXME: Currently precision is ignored, because passing a *NUMBERFMT
	// with only NumDigits initialized causes GetNumberFormat to fail.
	sPtr := syscall.StringToUTF16Ptr(s)

	bufSize := GetNumberFormat(
		LOCALE_USER_DEFAULT,
		0,
		sPtr,
		nil,
		nil,
		0)

	if bufSize == 0 {
		switch s {
		case "NaN", "-Inf", "+Inf":
			return s, nil
		}

		return "", lastError("GetNumberFormat")
	}

	buf := make([]uint16, bufSize)

	if 0 == GetNumberFormat(
		LOCALE_USER_DEFAULT,
		0,
		sPtr,
		nil,
		&buf[0],
		bufSize) {

		return "", lastError("GetNumberFormat")
	}

	return UTF16PtrToString(&buf[0]), nil
}
