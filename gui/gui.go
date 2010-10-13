// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
	"walk/drawing"
)

var defaultFont *drawing.Font

func init() {
	// Initialize default font
	var err os.Error
	defaultFont, err = drawing.NewFont("MS Shell Dlg", 8, 0)
	if err != nil {
		panic("failed to create default font")
	}
}
