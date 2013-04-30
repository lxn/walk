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
	title                 string
	tabWidget             *TabWidget
	titleProperty         Property
	titleChangedPublisher EventPublisher
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

	tp.titleProperty = NewProperty(
		func() interface{} {
			return tp.Title()
		},
		func(v interface{}) error {
			return tp.SetTitle(v.(string))
		},
		tp.titleChangedPublisher.Event())

	tp.MustRegisterProperty("Title", tp.titleProperty)

	return tp, nil
}

func (tp *TabPage) Enabled() bool {
	if tp.tabWidget != nil {
		return tp.tabWidget.Enabled() && tp.enabled
	}

	return tp.enabled
}

func (tp *TabPage) Font() *Font {
	if tp.font != nil {
		return tp.font
	} else if tp.tabWidget != nil {
		return tp.tabWidget.Font()
	}

	return defaultFont
}

func (tp *TabPage) Title() string {
	return tp.title
}

func (tp *TabPage) SetTitle(value string) error {
	tp.title = value

	tp.titleChangedPublisher.Publish()

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
