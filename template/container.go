package template

import (
	"walk"
)

type box struct {
	orientation walk.Orientation
	lf          LayoutFlag
	elems       []Widget
}

func HBox(lf LayoutFlag, elems ...Widget) Widget {
	return &box{walk.Horizontal, lf, elems}
}
func VBox(lf LayoutFlag, elems ...Widget) Widget {
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

func setupContainer(composite walk.Container,
	defaults LayoutFlag,
	options LayoutFlag,
	children []Widget) (err error) {
	if defaults != nil {
		defaults.SetupLayout(composite.Layout(), true)
	}
	if options != nil {
		options.SetupLayout(composite.Layout(), false)
	}
	for _, child := range children {
		if err = child.CreateElement(composite, defaults); err != nil {
			return
		}
	}
	return
}

func CreateBox(
	parent walk.Container,
	orientation walk.Orientation,
	defaults LayoutFlag,
	options LayoutFlag,
	children []Widget) error {

	composite, err := walk.NewComposite(parent)
	if err != nil {
		return err
	}
	if orientation == walk.Horizontal {
		composite.SetLayout(walk.NewHBoxLayout())
	} else {
		composite.SetLayout(walk.NewVBoxLayout())
	}
	if err = setupContainer(composite, options, defaults, children); err != nil {
		composite.Dispose()
		return err
	}
	return nil
}

type HBoxComposite struct {
	Options  LayoutFlag
	Children []Widget
}

func (b HBoxComposite) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return CreateBox(parent, walk.Horizontal, defaults, b.Options, b.Children)
}

type VBoxComposite struct {
	Options  LayoutFlag
	Children []Widget
	//FurtherDefaults LayoutFlag
}

func (b VBoxComposite) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	//if b.FurtherDefaults != nil {
	//	defaults = LayoutFlags(defaults, b.FurtherDefaults)
	//}
	return CreateBox(parent, walk.Vertical, defaults, b.Options, b.Children)
}

type HBoxCompositeChildren []Widget
func (b HBoxCompositeChildren) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return CreateBox(parent, walk.Horizontal, defaults, nil, []Widget(b))
}

type VBoxCompositeChildren []Widget
func (b VBoxCompositeChildren) CreateElement(parent walk.Container, defaults LayoutFlag) error {
	return CreateBox(parent, walk.Vertical, defaults, nil, []Widget(b))
}

