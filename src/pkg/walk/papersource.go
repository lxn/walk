// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "walk/winapi/gdi32"
)

type PaperSourceType int16

const (
	SourceUpper         PaperSourceType = DMBIN_UPPER
	SourceOnlyOne       PaperSourceType = DMBIN_ONLYONE
	SourceLower         PaperSourceType = DMBIN_LOWER
	SourceMiddle        PaperSourceType = DMBIN_MIDDLE
	SourceManual        PaperSourceType = DMBIN_MANUAL
	SourceEnvelope      PaperSourceType = DMBIN_ENVELOPE
	SourceEnvManual     PaperSourceType = DMBIN_ENVMANUAL
	SourceAuto          PaperSourceType = DMBIN_AUTO
	SourceTractor       PaperSourceType = DMBIN_TRACTOR
	SourceSmallFmt      PaperSourceType = DMBIN_SMALLFMT
	SourceLargeFmt      PaperSourceType = DMBIN_LARGEFMT
	SourceLargeCapacity PaperSourceType = DMBIN_LARGECAPACITY
	SourceCassette      PaperSourceType = DMBIN_CASSETTE
	SourceForm          PaperSourceType = DMBIN_FORMSOURCE
	SourceCustom        PaperSourceType = DMBIN_USER
)

type PaperSource struct {
	name string
	typ  PaperSourceType
}

func (ps *PaperSource) Name() string {
	return ps.name
}

func (ps *PaperSource) Type() PaperSourceType {
	return ps.typ
}
