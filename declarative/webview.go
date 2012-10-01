// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type WebView struct {
	AssignTo           **walk.WebView
	Name               string
	Disabled           bool
	Hidden             bool
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	URL                string
}

func (wv WebView) Create(parent walk.Container) error {
	w, err := walk.NewWebView(parent)
	if err != nil {
		return err
	}

	return InitWidget(wv, w, func() error {
		if err := w.SetURL(wv.URL); err != nil {
			return err
		}

		if wv.AssignTo != nil {
			*wv.AssignTo = w
		}

		return nil
	})
}

func (wv WebView) WidgetInfo() (name string, disabled, hidden bool, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return wv.Name, wv.Disabled, wv.Hidden, wv.MinSize, wv.MaxSize, wv.StretchFactor, wv.Row, wv.RowSpan, wv.Column, wv.ColumnSpan, wv.ContextMenuActions
}
