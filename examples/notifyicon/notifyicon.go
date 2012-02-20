// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"walk"
)

func main() {
	// Initialize walk and specify that we want errors to be panics.
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, _ := walk.NewMainWindow()

	// We load our icon from a file.
	icon, _ := walk.NewIconFromFile("../img/x.ico")

	// Create the notify icon and make sure we clean it up on exit.
	notifyIcon, _ := walk.NewNotifyIcon()
	defer notifyIcon.Dispose()

	// Set the icon and a tool tip text.
	notifyIcon.SetIcon(icon)
	notifyIcon.SetToolTip("Click for info or use the context menu to exit.")

	// When the left mouse button is pressed, bring up our about box.
	notifyIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		walk.MsgBox(mw, "About", "Walk NotifyIcon Example", walk.MsgBoxIconInformation)
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	exitAction.SetText("E&xit")
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	notifyIcon.ContextMenu().Actions().Add(exitAction)

	// The notify icon is hidden initially, so we have to make it visible.
	notifyIcon.SetVisible(true)

	// Run the message loop.
	mw.Run()
}
