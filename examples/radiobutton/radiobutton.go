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
}

func main() {
	foo := &Foo{"B"}

	MainWindow{
		Title:   "Walk RadioButton Example",
		MinSize: Size{320, 240},
		Layout:  VBox{},
		DataBinder: DataBinder{
			DataSource: foo,
			AutoSubmit: true,
			OnSubmitted: func() {
				fmt.Println(foo.Bar)
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
		},
	}.Run()
}
