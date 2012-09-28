// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TextEdit struct {
	Widget     **walk.TextEdit
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Font       Font
	Text       string
	ReadOnly   bool
}

func (te TextEdit) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewTextEdit(parent)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	w.SetName(te.Name)

	f, err := te.Font.Create()
	if err != nil {
		return nil, err
	}

	if f != nil {
		w.SetFont(f)
	}

	if err := w.SetText(te.Text); err != nil {
		return nil, err
	}

	w.SetReadOnly(te.ReadOnly)

	if te.Widget != nil {
		*te.Widget = w
	}

	succeeded = true

	return w, nil
}
