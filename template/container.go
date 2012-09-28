package template

import (
	"walk"
)

type box struct {
	orientation walk.Orientation
	lf          LayoutFlag
	elems       []GuiTemplate
}
func HBox(lf LayoutFlag, elems ...GuiTemplate) GuiTemplate {
	return &box{walk.Horizontal, lf, elems}
}
func VBox(lf LayoutFlag, elems ...GuiTemplate) GuiTemplate {
	return &box{walk.Vertical, lf, elems}
}
func (b *box) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	composite, err := walk.NewComposite(parent)
	if err != nil {
		return err
	}
	var l *walk.BoxLayout
	if b.orientation == walk.Horizontal {
		l = walk.NewHBoxLayout()
	} else {
		l = walk.NewVBoxLayout()
	}
	composite.SetLayout(l)
	defaults.SetupLayout(l, true)
	if b.lf != nil {
		b.lf.SetupLayout(l, false)
	}

	for _, elem := range b.elems {
		err = elem.CreateElement(composite, defaults)
		if err != nil {
			composite.Dispose()
			return err
		}
	}
	return err
}

