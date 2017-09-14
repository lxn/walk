// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	walk.Resources.SetRootDirPath("../img")

	MainWindow{
		Title:   "Walk ImageView Example",
		MinSize: Size{300, 200},
		Layout:  HBox{},
		Children: []Widget{
			ImageView{
				Image: "check.ico",
			},
			ImageView{
				Image: "plus.png",
			},
		},
	}.Run()
}
