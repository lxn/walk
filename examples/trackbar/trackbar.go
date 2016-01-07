// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var tbv, tbh *walk.TrackBar
	var maxEdit, minEdit, valueEdit *walk.NumberEdit

	data := struct{ Min, Max, Value int }{0, 100, 30}

	MainWindow{
		Title:   "Walk TrackBar Example",
		MinSize: Size{320, 240},
		Layout:  HBox{},
		Children: []Widget{
			TrackBar{
				AssignTo:    &tbv,
				MinValue:    data.Min,
				MaxValue:    data.Max,
				Value:       data.Value,
				Orientation: Vertical,
				OnValueChanged: func() {
					data.Value = tbv.Value()
					valueEdit.SetValue(float64(data.Value))

				},
			},
			Composite{
				Layout:        Grid{Columns: 3},
				StretchFactor: 4,
				Children: []Widget{
					Label{Text: "Min value"},
					Label{Text: "Value"},
					Label{Text: "Max value"},
					NumberEdit{
						AssignTo: &minEdit,
						Value:    float64(data.Min),
						OnValueChanged: func() {
							data.Min = int(minEdit.Value())
							tbh.SetRange(data.Min, data.Max)
							tbv.SetRange(data.Min, data.Max)
						},
					},
					NumberEdit{
						AssignTo: &valueEdit,
						Value:    float64(data.Value),
						OnValueChanged: func() {
							data.Value = int(valueEdit.Value())
							tbh.SetValue(data.Value)
							tbv.SetValue(data.Value)
						},
					},
					NumberEdit{
						AssignTo: &maxEdit,
						Value:    float64(data.Max),
						OnValueChanged: func() {
							data.Max = int(maxEdit.Value())
							tbh.SetRange(data.Min, data.Max)
							tbv.SetRange(data.Min, data.Max)
						},
					},
					TrackBar{
						ColumnSpan: 3,
						AssignTo:   &tbh,
						MinValue:   data.Min,
						MaxValue:   data.Max,
						Value:      data.Value,
						OnValueChanged: func() {
							data.Value = tbh.Value()
							valueEdit.SetValue(float64(data.Value))
						},
					},
					PushButton{
						ColumnSpan: 3,
						Text:       "Print state",
						OnClicked: func() {
							log.Printf("H: < %d | %d | %d >\n", tbh.MinValue(), tbh.Value(), tbh.MaxValue())
							log.Printf("V: < %d | %d | %d >\n", tbv.MinValue(), tbv.Value(), tbv.MaxValue())
						},
					},
				},
			},
		},
	}.Run()
}
