// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
)

type listViewColumnChangedHandler interface {
	onListViewColumnChanged(column *ListViewColumn)
}

type ListViewColumn struct {
	alignment       HorizontalAlignment
	width           int
	title           string
	changedHandlers vector.Vector
}

func NewListViewColumn() *ListViewColumn {
	return &ListViewColumn{width: 100}
}

func (c *ListViewColumn) Alignment() HorizontalAlignment {
	return c.alignment
}

func (c *ListViewColumn) SetAlignment(value HorizontalAlignment) {
	if value != c.alignment {
		c.alignment = value

		c.raiseChanged()
	}
}

func (c *ListViewColumn) Title() string {
	return c.title
}

func (c *ListViewColumn) SetTitle(value string) {
	if value != c.title {
		c.title = value

		c.raiseChanged()
	}
}

func (c *ListViewColumn) Width() int {
	return c.width
}

func (c *ListViewColumn) SetWidth(value int) {
	if value != c.width {
		c.width = value

		c.raiseChanged()
	}
}

func (c *ListViewColumn) addChangedHandler(handler listViewColumnChangedHandler) {
	c.changedHandlers.Push(handler)
}

func (c *ListViewColumn) removeChangedHandler(handler listViewColumnChangedHandler) {
	for i, h := range c.changedHandlers {
		if h.(listViewColumnChangedHandler) == handler {
			c.changedHandlers.Delete(i)
			break
		}
	}
}

func (c *ListViewColumn) raiseChanged() {
	for _, handlerIface := range c.changedHandlers {
		handler := handlerIface.(listViewColumnChangedHandler)
		handler.onListViewColumnChanged(c)
	}
}
