// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var le *walk.LineEdit
	var wv *walk.WebView

	MainWindow{
		Title:   "Walk WebView Example",
		MinSize: Size{800, 600},
		Layout:  VBox{},
		Children: []Widget{
			LineEdit{AssignTo: &le, OnReturnPressed: func() { wv.SetURL(le.Text()) }},
			WebView{AssignTo: &wv, URL: "http://golang.org"},
		},
	}.Run()
}
