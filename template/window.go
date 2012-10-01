package template

import (
	"walk"
)

type MainWindow struct {
	Title    string
	Defaults LayoutFlag
	Widget   Widget
}

func (wt MainWindow) Create(parent walk.Container) (w *walk.MainWindow, err error) {
	if w, err = walk.NewMainWindow(); err != nil {
		return
	}
	w.SetLayout(walk.NewVBoxLayout())

	if err = wt.Widget.CreateElement(w, wt.Defaults); err != nil {
		return nil, err
	}
	return
}
