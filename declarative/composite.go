// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Composite struct {
	Widget     **walk.Composite
	Name       string
	HStretch   int
	VStretch   int
	Row        int
	RowSpan    int
	Column     int
	ColumnSpan int
	Layout     Layout
	Children   []Widget
}

func (c Composite) Create(parent walk.Container) (walk.Widget, error) {
	w, err := walk.NewComposite(parent)
	if err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	w.SetName(c.Name)

	if c.Layout != nil {
		l, err := c.Layout.Create()
		if err != nil {
			return nil, err
		}

		if err := w.SetLayout(l); err != nil {
			return nil, err
		}
	}

	for _, child := range c.Children {
		if _, err := child.Create(w); err != nil {
			return nil, err
		}
	}

	if c.Widget != nil {
		*c.Widget = w
	}

	succeeded = true

	return w, nil
}
