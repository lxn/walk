package template

import "walk"

type WidgetFlag interface {
	SetupWidget(widget walk.Widget) error
}

type WidgetFlagFunc func(widget walk.Widget) error
func (f WidgetFlagFunc) SetupWidget(widget walk.Widget) error {
	return f(widget)
}

////////////////////////////////////////////////////////////////////////////////

// WidgetFlags allows the use of more than one widget flags
func WidgetFlags(ww ...WidgetFlag) WidgetFlag {
	return WidgetFlagFunc(func(widget walk.Widget) error {
		for _, w := range ww {
			if err := w.SetupWidget(widget); err != nil {
				return err
			}
		}
		return nil
	})
}

type ReadOnlySetter interface {
	SetReadOnly(readOnly bool) error
}

var ReadOnly WidgetFlagFunc = func(widget walk.Widget) error {
	return widget.(ReadOnlySetter).SetReadOnly(true)
}

type SetMaxLengther interface {
	SetMaxLength(value int)
}

func MaxLength(l int) WidgetFlag {
	return WidgetFlagFunc(func(widget walk.Widget) error {
		widget.(SetMaxLengther).SetMaxLength(l)
		return nil
	})
}

func StretchFaktor(f int) WidgetFlag {
	return WidgetFlagFunc(func(widget walk.Widget) error {
		return widget.Parent().Layout().(*walk.BoxLayout).SetStretchFactor(widget, f)
	})
}

type SetFonter interface {
	SetFont(value *walk.Font)
}

func Font(f *walk.Font) WidgetFlag {
	return WidgetFlagFunc(func(widget walk.Widget) error {
		widget.(SetFonter).SetFont(f)
		return nil
	})
}


