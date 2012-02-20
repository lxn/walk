// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import . "github.com/lxn/go-winapi"

type PaperSizeType int16

const (
	PaperLetter      PaperSizeType = DMPAPER_LETTER
	PaperLetterSmall PaperSizeType = DMPAPER_LETTERSMALL
	PaperTabloid     PaperSizeType = DMPAPER_TABLOID
	PaperLedger      PaperSizeType = DMPAPER_LEDGER
	PaperLegal       PaperSizeType = DMPAPER_LEGAL
	PaperStatement   PaperSizeType = DMPAPER_STATEMENT
	PaperExecutive   PaperSizeType = DMPAPER_EXECUTIVE
	PaperA3          PaperSizeType = DMPAPER_A3
	PaperA4          PaperSizeType = DMPAPER_A4
	PaperA4Small     PaperSizeType = DMPAPER_A4SMALL
	PaperA5          PaperSizeType = DMPAPER_A5
	PaperB4          PaperSizeType = DMPAPER_B4
	PaperB5          PaperSizeType = DMPAPER_B5
	PaperCustom      PaperSizeType = DMPAPER_USER
)

type PaperSize struct {
	name   string
	typ    PaperSizeType
	width  int
	height int
}

func (ps *PaperSize) Name() string {
	return ps.name
}

func (ps *PaperSize) Type() PaperSizeType {
	return ps.typ
}

func (ps *PaperSize) Width() int {
	return ps.width
}

func (ps *PaperSize) Height() int {
	return ps.height
}
