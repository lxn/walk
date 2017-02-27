// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var le *walk.LineEdit
	var wv *walk.WebView

	MainWindow{
		Icon:    Bind("'../img/' + icon(wv.URL) + '.ico'"),
		Title:   "Walk WebView Example'",
		MinSize: Size{800, 600},
		Layout:  VBox{MarginsZero: true},
		Children: []Widget{
			LineEdit{
				AssignTo: &le,
				Text:     Bind("wv.URL"),
				OnKeyDown: func(key walk.Key) {
					if key == walk.KeyReturn {
						wv.SetURL(le.Text())
					}
				},
			},
			WebView{
				AssignTo: &wv,
				Name:     "wv",
				URL:      "https://github.com/lxn/walk",
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"icon": func(args ...interface{}) (interface{}, error) {
				if strings.HasPrefix(args[0].(string), "https") {
					return "check", nil
				}

				return "stop", nil
			},
		},
	}.Run()
}
