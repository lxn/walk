package template

import (
	"walk"
)

type Widget interface {
	CreateElement(parent walk.Container, defaults LayoutFlag) error
}

type WidgetFunc func(parent walk.Container, defaults LayoutFlag) error

func (f WidgetFunc) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return f(parent, defaults)
}

type SimpleWidgetFunc func(parent walk.Container) error

func (f SimpleWidgetFunc) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return f(parent)
}

type SimpleWidget struct {
	wf WidgetFlag
	fn func(parent walk.Container) (walk.Widget, error)
}

func (t *SimpleWidget) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	w, err := t.fn(parent)
	if err == nil && t.wf != nil {
		t.wf.SetupWidget(w)
	}
	return err
}
