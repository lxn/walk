// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	if _, err := (MainWindow{
		Title:   "Walk LinkLabel Example",
		MinSize: Size{300, 200},
		Layout:  VBox{},
		Children: []Widget{
			LinkLabel{
				MaxSize: Size{100, 0},
				Text:    `I can contain multiple links like <a id="this" href="https://golang.org">this</a> or <a id="that" href="https://github.com/lxn/walk">that one</a>.`,
				OnLinkActivated: func(link *walk.LinkLabelLink) {
					log.Printf("id: '%s', url: '%s'\n", link.Id(), link.URL())
				},
			},
		},
	}).Run(); err != nil {
		log.Fatal(err)
	}
}
