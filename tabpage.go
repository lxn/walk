// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import "syscall"

import . "github.com/lxn/go-winapi"

const tabPageWindowClass = `\o/ Walk_TabPage_Class \o/`

var tabPageBackgroundBrush Brush

func init() {
	MustRegisterWindowClass(tabPageWindowClass)

	tabPageBackgroundBrush, _ = NewSystemColorBrush(COLOR_WINDOW)
}

type TabPage struct {
	ContainerBase
	title     string
	tabWidget *TabWidget
}

func NewTabPage() (*TabPage, error) {
	tp := &TabPage{}

	if err := InitWidget(
		tp,
		nil,
		tabPageWindowClass,
		WS_POPUP,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	tp.children = newWidgetList(tp)

	tp.SetBackground(tabPageBackgroundBrush)

	return tp, nil
}

func (tp *TabPage) Dispose() {
	b := tp.Background()
	if b != nil {
		b.Dispose()
		tp.SetBackground(nil)
	}

	tp.ContainerBase.Dispose()
}

func (tp *TabPage) Title() string {
	return tp.title
}

func (tp *TabPage) SetTitle(value string) error {
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
		CchTextMax: int32(len(text)),
	}

	return item
}
