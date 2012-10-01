// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type Layout interface {
	Create() (walk.Layout, error)
}

type Widget interface {
	Create(parent walk.Container) error
	WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action)
}

type Container interface {
	ContainerInfo() (Layout, []Widget)
}

type MenuItem interface {
	createMenuAction(menu *walk.Menu) (*walk.Action, error)
}

type ToolBarItem interface {
	createToolBarAction() (*walk.Action, error)
}
