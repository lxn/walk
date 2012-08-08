// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

import (
	"github.com/lxn/walk"
)

func main() {
	walk.Initialize(walk.InitParams{})
	defer walk.Shutdown()

	if _, err := runMainWindow(); err != nil {
		log.Fatal(err)
	}
}
