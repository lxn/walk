package template

import (
	"walk"
)

type GuiTemplate interface {
	CreateElement(parent walk.Container, defaults LayoutFlag) error
}

type GuiTemplateFunc func(parent walk.Container, defaults LayoutFlag) error
func (f GuiTemplateFunc) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return f(parent, defaults)
}

type SimpleGuiTemplateFunc func(parent walk.Container) error
func (f SimpleGuiTemplateFunc) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return f(parent)
}

type SimpleGuiTemplate struct {
	wf WidgetFlag
	fn func(parent walk.Container) (walk.Widget, error)
}
func (t *SimpleGuiTemplate) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	w, err := t.fn(parent)
	if err == nil && t.wf != nil {
		t.wf.SetupWidget(w)
	}
	return err
}


