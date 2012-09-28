// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type MainWindow struct {
	Widget   **walk.MainWindow
	Name     string
	Font     Font
	Title    string
	Layout   Layout
	Children []Widget
}

func (mw MainWindow) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewMainWindow()
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	if err := initWidget(mw, w); err != nil {
		return nil, err
	}

	w.SetName(mw.Name)

	if err := w.SetTitle(mw.Title); err != nil {
		return nil, err
	}

	if mw.Widget != nil {
		*mw.Widget = w
	}

	succeeded = true

	return w, nil
}

func (mw MainWindow) LayoutParams() (hStretch, vStretch, row, rowSpan, col, colSpan int) {
	return 0, 0, 0, 0, 0, 0
}

func (mw MainWindow) Layout_() Layout {
	return mw.Layout
}

func (mw MainWindow) Children_() []Widget {
	return mw.Children
}
