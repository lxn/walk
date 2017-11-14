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

	type Mode struct {
		Name  string
		Value ImageViewMode
	}

	modes := []Mode{
		{"ImageViewModeIdeal", ImageViewModeIdeal},
		{"ImageViewModeCorner", ImageViewModeCorner},
		{"ImageViewModeCenter", ImageViewModeCenter},
		{"ImageViewModeShrink", ImageViewModeShrink},
		{"ImageViewModeZoom", ImageViewModeZoom},
		{"ImageViewModeStretch", ImageViewModeStretch},
	}

	var widgets []Widget

	for _, mode := range modes {
		widgets = append(widgets,
			Label{
				Text: mode.Name,
			},
			ImageView{
				Background: SolidColorBrush{Color: walk.RGB(255, 191, 0)},
				Image:      "open.png",
				Margin:     10,
				Mode:       mode.Value,
			},
		)
	}

	MainWindow{
		Title:    "Walk ImageView Example",
		Size:     Size{400, 600},
		Layout:   Grid{Columns: 2},
		Children: widgets,
	}.Run()
}
