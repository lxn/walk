// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

var labelOrigWndProcPtr uintptr
var _ subclassedWidget = &Label{}

type Label struct {
	WidgetBase
}

func NewLabel(parent Container) (*Label, os.Error) {
	l := &Label{}

	if err := initChildWidget(
		l,
		parent,
		"STATIC",
		WS_VISIBLE|SS_CENTERIMAGE,
		0); err != nil {
		return nil, err
	}

	return l, nil
}

func (*Label) origWndProcPtr() uintptr {
	return labelOrigWndProcPtr
}

func (*Label) setOrigWndProcPtr(ptr uintptr) {
	labelOrigWndProcPtr = ptr
}

func (*Label) LayoutFlags() LayoutFlags {
	return GrowableVert
}

func (l *Label) MinSizeHint() Size {
	return l.calculateTextSize()
}

func (l *Label) SizeHint() Size {
	return l.calculateTextSize()
}
