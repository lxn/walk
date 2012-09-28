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

	w.SetName(mw.Name)

	f, err := mw.Font.Create()
	if err != nil {
		return nil, err
	}

	if f != nil {
		w.SetFont(f)
	}

	if err := w.SetTitle(mw.Title); err != nil {
		return nil, err
	}

	if mw.Layout != nil {
		l, err := mw.Layout.Create()
		if err != nil {
			return nil, err
		}

		if err := w.SetLayout(l); err != nil {
			return nil, err
		}
	}

	for _, child := range mw.Children {
		if _, err := child.Create(w); err != nil {
			return nil, err
		}
	}

	if mw.Widget != nil {
		*mw.Widget = w
	}

	succeeded = true

	return w, nil
}
