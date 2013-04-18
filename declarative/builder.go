// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	"github.com/lxn/walk"
)

type declWidget struct {
	d Widget
	w walk.Widget
}

type Builder struct {
	level         int
	parent        walk.Container
	declWidgets   []declWidget
	name2Widget   map[string]walk.Widget
	deferredFuncs []func() error
}

func NewBuilder(parent walk.Container) *Builder {
	return &Builder{
		parent:      parent,
		name2Widget: make(map[string]walk.Widget),
	}
}

func (b *Builder) Parent() walk.Container {
	return b.parent
}

func (b *Builder) Defer(f func() error) {
	b.deferredFuncs = append(b.deferredFuncs, f)
}

func (b *Builder) InitWidget(d Widget, w walk.Widget, customInit func() error) error {
	b.level++
	defer func() {
		b.level--
	}()

	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	b.declWidgets = append(b.declWidgets, declWidget{d, w})

	// Widget
	name, _, _, font, toolTipText, minSize, maxSize, stretchFactor, row, rowSpan, column, columnSpan, contextMenuActions, onKeyDown, onMouseDown, onMouseMove, onMouseUp, onSizeChanged := d.WidgetInfo()

	w.SetName(name)

	if name != "" {
		b.name2Widget[name] = w
	}

	if toolTipText != "" {
		if err := w.SetToolTipText(toolTipText); err != nil {
			return err
		}
	}

	if err := w.SetMinMaxSize(minSize.toW(), maxSize.toW()); err != nil {
		return err
	}

	if len(contextMenuActions) > 0 {
		cm, err := walk.NewMenu()
		if err != nil {
			return err
		}
		if err := addToActionList(cm.Actions(), contextMenuActions); err != nil {
			return err
		}
		w.SetContextMenu(cm)
	}

	if onKeyDown != nil {
		w.KeyDown().Attach(onKeyDown)
	}

	if onMouseDown != nil {
		w.MouseDown().Attach(onMouseDown)
	}

	if onMouseMove != nil {
		w.MouseMove().Attach(onMouseMove)
	}

	if onMouseUp != nil {
		w.MouseUp().Attach(onMouseUp)
	}

	if onSizeChanged != nil {
		w.SizeChanged().Attach(onSizeChanged)
	}

	if p := w.Parent(); p != nil {
		switch l := p.Layout().(type) {
		case *walk.BoxLayout:
			if stretchFactor < 1 {
				stretchFactor = 1
			}
			if err := l.SetStretchFactor(w, stretchFactor); err != nil {
				return err
			}

		case *walk.GridLayout:
			cs := columnSpan
			if cs < 1 {
				cs = 1
			}
			rs := rowSpan
			if rs < 1 {
				rs = 1
			}
			r := walk.Rectangle{column, row, cs, rs}

			if err := l.SetRange(w, r); err != nil {
				return err
			}
		}
	}

	oldParent := b.parent

	// Container
	var db *walk.DataBinder
	if dc, ok := d.(Container); ok {
		if wc, ok := w.(walk.Container); ok {
			dataBinder, layout, children := dc.ContainerInfo()

			if layout != nil {
				l, err := layout.Create()
				if err != nil {
					return err
				}

				if err := wc.SetLayout(l); err != nil {
					return err
				}
			}

			b.parent = wc
			defer func() {
				b.parent = oldParent
			}()

			for _, child := range children {
				if err := child.Create(b); err != nil {
					return err
				}
			}

			var err error
			if db, err = dataBinder.create(); err != nil {
				return err
			}
		}
	}

	// Custom
	if customInit != nil {
		if err := customInit(); err != nil {
			return err
		}
	}

	b.parent = oldParent

	// Widget continued
	if font != nil {
		if f, err := font.Create(); err != nil {
			return err
		} else if f != nil {
			w.SetFont(f)
		}
	}

	if b.level == 1 {
		if err := b.initProperties(); err != nil {
			return err
		}
	}

	// Call Reset on DataBinder after customInit, so a Dialog gets a chance to first
	// wire up its DefaultButton to the CanSubmitChanged event of a DataBinder.
	if db != nil {
		if _, ok := d.(Container); ok {
			if wc, ok := w.(walk.Container); ok {
				// FIXME: Currently SetDataBinder must be called after initProperties.
				wc.SetDataBinder(db)

				if err := db.Reset(); err != nil {
					return err
				}
			}
		}
	}

	if b.level == 1 {
		for _, f := range b.deferredFuncs {
			if err := f(); err != nil {
				return err
			}
		}
	}

	succeeded = true

	return nil
}

func (b *Builder) initProperties() error {
	for _, dw := range b.declWidgets {
		d, w := dw.d, dw.w

		sv := reflect.ValueOf(d)
		st := sv.Type()
		if st.Kind() != reflect.Struct {
			panic("d must be a struct value")
		}

		wb := w.BaseWidget()

		fieldCount := st.NumField()
		for i := 0; i < fieldCount; i++ {
			sf := st.Field(i)

			prop := wb.Property(sf.Name)

			switch val := sv.Field(i).Interface().(type) {
			case nil:
				// nop

			case Bind:
				if prop == nil {
					panic(sf.Name + " is not a property")
				}

				prop.SetSource(val.To)
				if val.Validator != nil {
					validator, err := val.Validator.Create()
					if err != nil {
						return err
					}
					prop.SetValidator(validator)
				}

			case BindTo:
				if prop == nil {
					panic(sf.Name + " is not a property")
				}

				prop.SetSource(val.Name)

			case BindProperty:
				if prop == nil {
					panic(sf.Name + " is not a registered property")
				}

				parts := strings.Split(val.Name, ".")
				if len(parts) != 2 {
					panic("invalid BindProperty syntax: " + val.Name)
				}

				var srcProp *walk.Property
				if sw, ok := b.name2Widget[parts[0]]; ok {
					sbw := sw.BaseWidget()
					srcProp = sbw.Property(parts[1])
					if srcProp == nil {
						panic("unknown source property: " + parts[1])
					}
					prop.SetSource(srcProp)
				} else {
					panic("unknown widget: " + parts[0])
				}

			default:
				if prop == nil {
					continue
				}

				v := prop.Get()
				valt, vt := reflect.TypeOf(val), reflect.TypeOf(v)

				if valt != vt {
					panic(fmt.Sprintf("cannot assign value %v of type %T to property %s of type %T", val, val, prop.Name(), v))
				}
				if err := prop.Set(val); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
