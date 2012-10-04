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
	createAction(menu *walk.Menu) (*walk.Action, error)
}

type topLevelWindowInfo struct {
	Name               string
	Disabled           bool
	Hidden             bool
	Font               Font
	MinSize            Size
	MaxSize            Size
	ContextMenuActions []*walk.Action
	Layout             Layout
	Children           []Widget
}

func (topLevelWindowInfo) Create(parent walk.Container) error {
	return nil
}

func (i topLevelWindowInfo) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return i.Name, i.Disabled, i.Hidden, &i.Font, i.MinSize, i.MaxSize, 0, 0, 0, 0, 0, i.ContextMenuActions
}

func (i topLevelWindowInfo) ContainerInfo() (Layout, []Widget) {
	return i.Layout, i.Children
}
