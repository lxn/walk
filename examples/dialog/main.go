// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

func main() {
	if _, err := runMainWindow(); err != nil {
		log.Fatal(err)
	}
}
