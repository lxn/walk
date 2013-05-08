// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

import (
	. "github.com/lxn/walk/declarative"
)

type Foo struct {
	Bar string
	Baz int
}

func main() {
	foo := &Foo{"B", 0}

	MainWindow{
		Title:   "Walk RadioButton Example",
		MinSize: Size{320, 240},
		Layout:  VBox{},
		DataBinder: DataBinder{
			DataSource: foo,
			AutoSubmit: true,
			OnSubmitted: func() {
				fmt.Println(foo)
			},
		},
		Children: []Widget{
			RadioButton{
				Name:                 "aRB",
				Text:                 "A",
				CheckedValue:         Bind("Bar"),
				CheckedDiscriminator: "A",
			},
			RadioButton{
				Name:                 "bRB",
				Text:                 "B",
				CheckedValue:         Bind("Bar"),
				CheckedDiscriminator: "B",
			},
			RadioButton{
				Name:                 "cRB",
				Text:                 "C",
				CheckedValue:         Bind("Bar"),
				CheckedDiscriminator: "C",
			},
			Label{
				Text:    "A",
				Enabled: Bind("aRB.Checked"),
			},
			Label{
				Text:    "B",
				Enabled: Bind("bRB.Checked"),
			},
			Label{
				Text:    "C",
				Enabled: Bind("cRB.Checked"),
			},
			// These will become their own group, because they are separated
			// from the other radio buttons by the labels.
			RadioButton{
				Name:                 "oneRB",
				Text:                 "1",
				CheckedValue:         Bind("Baz"),
				CheckedDiscriminator: 1,
			},
			RadioButton{
				Name:                 "twoRB",
				Text:                 "2",
				CheckedValue:         Bind("Baz"),
				CheckedDiscriminator: 2,
			},
			RadioButton{
				Name:                 "threeRB",
				Text:                 "3",
				CheckedValue:         Bind("Baz"),
				CheckedDiscriminator: 3,
			},
			Label{
				Text:    "1",
				Enabled: Bind("oneRB.Checked"),
			},
			Label{
				Text:    "2",
				Enabled: Bind("twoRB.Checked"),
			},
			Label{
				Text:    "3",
				Enabled: Bind("threeRB.Checked"),
			},
		},
	}.Run()
}
