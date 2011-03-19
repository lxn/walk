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
	title string
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

func (tp *TabPage) Title() string {
	return tp.title
}

func (tp *TabPage) SetTitle(value string) os.Error {
	tp.title = value

	return nil
}
