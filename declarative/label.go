// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Label struct {
	Widget     **walk.Label
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Font       Font
	Text       string
}

func (l Label) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewLabel(parent)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	w.SetName(l.Name)

	f, err := l.Font.Create()
	if err != nil {
		return nil, err
	}

	if f != nil {
		w.SetFont(f)
	}

	if err := w.SetText(l.Text); err != nil {
		return nil, err
	}

	if l.Widget != nil {
		*l.Widget = w
	}

	succeeded = true

	return w, nil
}
