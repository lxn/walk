package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	var (
		e error
		winMain *walk.MainWindow
		label   *walk.Label
	)

	winMain, e = walk.NewMainWindow()
	if e != nil { panic(e) }

	// Define hotkey handler
	winMain.Hotkey().Attach(func(hkid int) {
		label.SetText("")
		switch hkid {
		case 1:
			label.SetText("Global hotkey 1 pressed: Ctrl+Alt+X")
		case 2:
			label.SetText("Global hotkey 2 pressed: Alt+Shift+D")
		}
	})

	// Register hotkeys globally
	walk.RegisterGlobalHotKey(winMain, 1, walk.Shortcut{Modifiers: walk.ModControl | walk.ModAlt, Key: walk.KeyX})
	walk.RegisterGlobalHotKey(winMain, 2, walk.Shortcut{Modifiers: walk.ModShift   | walk.ModAlt, Key: walk.KeyD})

	MainWindow {
		AssignTo: &winMain,
		Size: Size{400, 120},
		Layout: VBox{},
		Children: []Widget {
			Label {
				Text: "Focus on another window and press Ctrl+Alt+X or Alt+Shift+D",
			},
			Label {
				AssignTo: &label,
				Text: "-",
				Font: Font{PointSize: 12},
			},
		},
	}.Run()

}
