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

	bg := SolidColorBrush{Color: walk.RGB(255, 191, 0)}

	MainWindow{
		Title:  "Walk ImageView Example",
		Size:   Size{400, 600},
		Layout: Grid{Columns: 2},
		Children: []Widget{
			Label{
				Text: "ImageViewModeIdeal",
			},
			ImageView{
				Background: bg,
				Image:      "open.png",
				Margin:     10,
				Mode:       ImageViewModeIdeal,
			},
			Label{
				Text: "ImageViewModeCorner",
			},
			ImageView{
				Background: bg,
				Image:      "open.png",
				Margin:     10,
				Mode:       ImageViewModeCorner,
			},
			Label{
				Text: "ImageViewModeCenter",
			},
			ImageView{
				Background: bg,
				Image:      "open.png",
				Margin:     10,
				Mode:       ImageViewModeCenter,
			},
			Label{
				Text: "ImageViewModeShrink",
			},
			ImageView{
				Background: bg,
				Image:      "open.png",
				Margin:     10,
				Mode:       ImageViewModeShrink,
			},
			Label{
				Text: "ImageViewModeZoom",
			},
			ImageView{
				Background: bg,
				Image:      "open.png",
				Margin:     10,
				Mode:       ImageViewModeZoom,
			},
			Label{
				Text: "ImageViewModeStretch",
			},
			ImageView{
				Background: bg,
				Image:      "open.png",
				Margin:     10,
				Mode:       ImageViewModeStretch,
			},
		},
	}.Run()
}
