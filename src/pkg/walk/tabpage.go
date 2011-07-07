// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
)

import . "walk/winapi"

const tabPageWindowClass = `\o/ Walk_TabPage_Class \o/`

var tabPageWindowClassRegistered bool

type TabPage struct {
	ContainerBase
	title     string
	tabWidget *TabWidget
}

func NewTabPage() (*TabPage, os.Error) {
	ensureRegisteredWindowClass(tabPageWindowClass, &tabPageWindowClassRegistered)

	tp := &TabPage{}

	if err := initWidget(
		tp,
		nil,
		tabPageWindowClass,
		WS_POPUP,
		WS_EX_CONTROLPARENT /*|WS_EX_TRANSPARENT*/ ); err != nil {
		return nil, err
	}

	tp.children = newWidgetList(tp)

	// FIXME: The next line, together with WS_EX_TRANSPARENT, would make the tab
	// page background transparent, but it doesn't work on XP :(
	//	tp.SetBackground(NullBrush())

	return tp, nil
}

func (tp *TabPage) Title() string {
	return tp.title
}

func (tp *TabPage) SetTitle(value string) os.Error {
	tp.title = value

	if tp.tabWidget == nil {
		return nil
	}

	return tp.tabWidget.onPageChanged(tp)
}

func (tp *TabPage) tcItem() *TCITEM {
	text := syscall.StringToUTF16(tp.Title())

	item := &TCITEM{
		Mask:       TCIF_TEXT,
		PszText:    &text[0],
		CchTextMax: len(text),
	}

	return item
}
