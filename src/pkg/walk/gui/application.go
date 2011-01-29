// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	. "walk/winapi/user32"
)

func Exit(exitCode int) {
	PostQuitMessage(exitCode)
}
