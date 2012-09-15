// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"
)

type Foo struct {
	FilePath string
	Bar      *Bar
	Number   float64
	Date     time.Time
	Memo     string
}

type Bar struct {
	Name string
	Baz  string
}
