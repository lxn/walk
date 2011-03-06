// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

const tabPageWindowClass = `\o/ Walk_TabPage_Class \o/`

var tabPageWindowClassRegistered bool

type TabPage struct {
	ContainerBase
}

func NewTabPage() (*TabPage, os.Error) {
	ensureRegisteredWindowClass(tabPageWindowClass, &tabPageWindowClassRegistered)

	tp := &TabPage{}

	if err := initWidget(
		tp,
		nil,
		tabPageWindowClass,
		WS_POPUP,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	tp.children = newWidgetList(tp)

	return tp, nil
}

func (*TabPage) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (tp *TabPage) PreferredSize() Size {
	return tp.dialogBaseUnitsToPixels(Size{100, 100})
}
