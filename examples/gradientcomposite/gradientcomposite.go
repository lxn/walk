// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	MainWindow{
		Title:   "Walk GradientComposite Example",
		MinSize: Size{360, 0},
		Layout:  HBox{MarginsZero: true},
		Children: []Widget{
			GradientComposite{
				Vertical: Bind("verticalCB.Checked"),
				Color1:   Bind("rgb(c1RedSld.Value, c1GreenSld.Value, c1BlueSld.Value)"),
				Color2:   Bind("rgb(c2RedSld.Value, c2GreenSld.Value, c2BlueSld.Value)"),
				Layout:   HBox{},
				Children: []Widget{
					GroupBox{
						Title:  "Gradient Parameters",
						Layout: VBox{},
						Children: []Widget{
							CheckBox{Name: "verticalCB", Text: "Vertical", Checked: true},
							GroupBox{
								Title:  "Color1",
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{Text: "Red:"},
									Slider{Name: "c1RedSld", Tracking: true, MaxValue: 255, Value: 95},
									Label{Text: "Green:"},
									Slider{Name: "c1GreenSld", Tracking: true, MaxValue: 255, Value: 191},
									Label{Text: "Blue:"},
									Slider{Name: "c1BlueSld", Tracking: true, MaxValue: 255, Value: 255},
								},
							},
							GroupBox{
								Title:  "Color2",
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{Text: "Red:"},
									Slider{Name: "c2RedSld", Tracking: true, MaxValue: 255, Value: 239},
									Label{Text: "Green:"},
									Slider{Name: "c2GreenSld", Tracking: true, MaxValue: 255, Value: 63},
									Label{Text: "Blue:"},
									Slider{Name: "c2BlueSld", Tracking: true, MaxValue: 255, Value: 0},
								},
							},
						},
					},
				},
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"rgb": func(args ...interface{}) (interface{}, error) {
				return walk.RGB(byte(args[0].(float64)), byte(args[1].(float64)), byte(args[2].(float64))), nil
			},
		},
	}.Run()
}
