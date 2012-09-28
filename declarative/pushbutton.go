// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type PushButton struct {
	Widget     **walk.PushButton
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Font       Font
	Text       string
	OnClicked  walk.EventHandler
}

func (pb PushButton) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewPushButton(parent)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	w.SetName(pb.Name)

	f, err := pb.Font.Create()
	if err != nil {
		return nil, err
	}

	if f != nil {
		w.SetFont(f)
	}

	if err := w.SetText(pb.Text); err != nil {
		return nil, err
	}

	if pb.OnClicked != nil {
		w.Clicked().Attach(pb.OnClicked)
	}

	if pb.Widget != nil {
		*pb.Widget = w
	}

	succeeded = true

	return w, nil
}
