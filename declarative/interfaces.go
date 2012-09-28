// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Widget interface {
	Create(parent walk.Container) (walk.Widget, error)
}

type Layout interface {
	Create() (walk.Layout, error)
}

type Container interface {
	Layout_() Layout
	Children_() []Widget
}
