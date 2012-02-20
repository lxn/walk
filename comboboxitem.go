// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

type ComboBoxItem struct {
	text     string
	userData interface{}
}

func NewComboBoxItem() *ComboBoxItem {
	return &ComboBoxItem{}
}

func (cbi *ComboBoxItem) Text() string {
	return cbi.text
}

func (cbi *ComboBoxItem) SetText(value string) os.Error {
	cbi.text = value

	// FIXME: Update ComboBox
	return nil
}

func (cbi *ComboBoxItem) UserData() interface{} {
	return cbi.userData
}

func (cbi *ComboBoxItem) SetUserData(value interface{}) {
	cbi.userData = value
}
