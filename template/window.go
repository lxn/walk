package template

import (
	"walk"
)

func MainWindow(title string, defaults LayoutFlag, tmpl GuiTemplate) (w *walk.MainWindow, err error) {
	w, err = walk.NewMainWindow()
	if err == nil {
		err = w.SetTitle(title)
		w.SetLayout(walk.NewVBoxLayout())
		if err == nil {
			err = tmpl.CreateElement(w, defaults)
		}
	}
	if err != nil {
		w.Dispose()
		w = nil
	}
	return
}


