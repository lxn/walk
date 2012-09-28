package template

import (
	"errors"
	"walk"
)

type LayoutFlag interface {
	SetupLayout(layout walk.Layout, isdefault bool) error
}

type LayoutFlagFunc func(layout walk.Layout, isdefault bool) error
func (f LayoutFlagFunc) SetupLayout(layout walk.Layout, isdefault bool) error {
	return f(layout, isdefault)
}

////////////////////////////////////////////////////////////////////////////////

func LayoutFlags(ll ...LayoutFlag) LayoutFlag {
	return LayoutFlagFunc(func(layout walk.Layout, isdefault bool) error {
		for _, l := range ll {
			if err := l.SetupLayout(layout, isdefault); err != nil {
				return err
			}
		}
		return nil
	})
}

func Spacing(s int) LayoutFlag {
	return LayoutFlagFunc(func(layout walk.Layout, isdefault bool) error {
		if b, ok := layout.(*walk.BoxLayout); ok {
			b.SetSpacing(s)
		}
		if !isdefault {
			return errors.New("Setting spacing on non-box layout")
		}
		return nil
	})
}

func Margins(hn, vn, hf, vf int) LayoutFlag {
	return LayoutFlagFunc(func(layout walk.Layout, isdefault bool) error {
		return layout.SetMargins(walk.Margins{hn, vn, hf, vf})
	})
}

