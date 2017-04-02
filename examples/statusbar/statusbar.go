// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates the status bar, including a size gripper
// attached to the bottom of the main window.
// The status bar has two items, one is dynamically updated and one includes an icon.
package main

import (
	"log"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var mw *walk.MainWindow
	var sbi1 *walk.StatusBarItem

	icon, err := walk.NewIconFromFile("../img/x.ico")
	if err != nil {
		log.Fatal(err)
	}

	m := MainWindow{
		AssignTo: &mw,
		Title:    "Walk Statusbar Example",
		MinSize:  Size{600, 200},
		Layout:   VBox{MarginsZero: true},
		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				AssignTo:    &sbi1,
				Text:        "item 1",
				ToolTipText: "tool tip text...", // This does not show up!
				Width:       300,
			},
			StatusBarItem{
				Icon: icon,
				Text: "item 2",
			},
		},
		Children: []Widget{},
	}
	if err := m.Create(); err != nil {
		log.Fatal(err)
	}

	c := time.Tick(1 * time.Second)
	go func(c <-chan time.Time) {
		for now := range c {
			if sbi1 != nil {
				sbi1.SetText(now.String())
			}
		}
	}(c)

	mw.Run()
}
