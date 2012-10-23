// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
)

func main() {
	// Specify that we want errors to be panics.
	walk.SetPanicOnError(true)

	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, _ := walk.NewMainWindow()

	// We load our icon from a file.
	icon, _ := walk.NewIconFromFile("../img/x.ico")

	// Create the notify icon and make sure we clean it up on exit.
	ni, _ := walk.NewNotifyIcon()
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	ni.SetIcon(icon)
	ni.SetToolTip("Click for info or use the context menu to exit.")

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		ni.ShowCustom(
			"Walk NotifyIcon Example",
			"There are multiple ShowX methods sporting different icons.")
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	exitAction.SetText("E&xit")
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	ni.ContextMenu().Actions().Add(exitAction)

	// The notify icon is hidden initially, so we have to make it visible.
	ni.SetVisible(true)

	// Now that the icon is visible, we can bring up an info balloon.
	ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again.")

	// Run the message loop.
	mw.Run()
}
