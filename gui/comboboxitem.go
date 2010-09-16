// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

type ComboBoxItem struct {
	text string
}

func NewComboBoxItem() *ComboBoxItem {
	return &ComboBoxItem{}
}

func (cbi *ComboBoxItem) Text() string {
	return cbi.text
}

func (cbi *ComboBoxItem) SetText(value string) os.Error {
	panic("not implemented")
}
