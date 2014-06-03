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

var conditionsByName = make(map[string]walk.Condition)

func MustRegisterCondition(name string, condition walk.Condition) {
	if name == "" {
		panic(`name == ""`)
	}
	if condition == nil {
		panic("condition == nil")
	}
	if _, ok := conditionsByName[name]; ok {
		panic("name already registered")
	}

	conditionsByName[name] = condition
}

type declWidget struct {
	d Widget
	w walk.Window
}

type Builder struct {
	level                    int
	columns                  int
	row                      int
	col                      int
	parent                   walk.Container
	declWidgets              []declWidget
	name2Window              map[string]walk.Window
	deferredFuncs            []func() error
	knownCompositeConditions map[string]walk.Condition
}

func NewBuilder(parent walk.Container) *Builder {
	return &Builder{
		parent:                   parent,
		name2Window:              make(map[string]walk.Window),
		knownCompositeConditions: make(map[string]walk.Condition),
	}
}

func (b *Builder) Parent() walk.Container {
	return b.parent
}

func (b *Builder) Defer(f func() error) {
	b.deferredFuncs = append(b.deferredFuncs, f)
}

func (b *Builder) deferBuildMenuActions(menu *walk.Menu, items []MenuItem) {
	if len(items) > 0 {
		b.Defer(func() error {
			for _, item := range items {
				if _, err := item.createAction(b, menu); err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (b *Builder) deferBuildActions(actionList *walk.ActionList, items []MenuItem) {
	if len(items) > 0 {
		b.Defer(func() error {
			for _, item := range items {
				action, err := item.createAction(b, nil)
				if err != nil {
					return err
				}
				if err := actionList.Add(action); err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (b *Builder) InitWidget(d Widget, w walk.Window, customInit func() error) error {
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
	name, _, _, font, toolTipText, minSize, maxSize, stretchFactor, row, rowSpan, column, columnSpan, contextMenuItems, onKeyDown, onKeyPress, onKeyUp, onMouseDown, onMouseMove, onMouseUp, onSizeChanged := d.WidgetInfo()

	w.SetName(name)

	if name != "" {
		b.name2Window[name] = w
	}

	if toolTipText != "" {
		if widget, ok := w.(walk.Widget); ok {
			if err := widget.SetToolTipText(toolTipText); err != nil {
				return err
			}
		}
	}

	if err := w.SetMinMaxSize(minSize.toW(), maxSize.toW()); err != nil {
		return err
	}

	if len(contextMenuItems) > 0 {
		cm, err := walk.NewMenu()
		if err != nil {
			return err
		}

		b.deferBuildMenuActions(cm, contextMenuItems)

		w.SetContextMenu(cm)
	}

	if onKeyDown != nil {
		w.KeyDown().Attach(onKeyDown)
	}

	if onKeyPress != nil {
		w.KeyPress().Attach(onKeyPress)
	}

	if onKeyUp != nil {
		w.KeyUp().Attach(onKeyUp)
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

	if widget, ok := w.(walk.Widget); ok {
		if p := widget.Parent(); p != nil {
			switch l := p.Layout().(type) {
			case *walk.BoxLayout:
				if stretchFactor < 1 {
					stretchFactor = 1
				}
				if err := l.SetStretchFactor(widget, stretchFactor); err != nil {
					return err
				}

			case *walk.GridLayout:
				if rowSpan < 1 {
					rowSpan = 1
				}
				if columnSpan < 1 {
					columnSpan = 1
				}

				if b.columns > 0 && row == 0 && column == 0 {
					if b.col+columnSpan > b.columns {
						b.row++
						b.col = 0
					}

					row = b.row
					column = b.col

					b.col += columnSpan
				}

				r := walk.Rectangle{column, row, columnSpan, rowSpan}

				if err := l.SetRange(widget, r); err != nil {
					return err
				}
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

			if g, ok := layout.(Grid); ok {
				columns := b.columns
				defer func() {
					b.columns, b.row, b.col = columns, row, column+columnSpan
				}()

				b.columns = g.Columns
				b.row = 0
				b.col = 0
			}

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
				b.Defer(func() error {
					// FIXME: Currently SetDataBinder must be called after initProperties.
					wc.SetDataBinder(db)

					return db.Reset()
				})
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

		wb := w.AsWindowBase()

		fieldCount := st.NumField()
		for i := 0; i < fieldCount; i++ {
			sf := st.Field(i)

			prop := wb.Property(sf.Name)

			switch val := sv.Field(i).Interface().(type) {
			case nil:
				// nop

			case bindData:
				if prop == nil {
					panic(sf.Name + " is not a property")
				}

				src := b.conditionOrProperty(val)

				if src == nil {
					// No luck so far, so we assume the expression refers to
					// something in the data source.
					src = val.expression

					if val.validator != nil {
						validator, err := val.validator.Create()
						if err != nil {
							return err
						}
						if err := prop.SetValidator(validator); err != nil {
							return err
						}
					}
				}

				if err := prop.SetSource(src); err != nil {
					return err
				}

			case walk.Condition:
				if prop == nil {
					panic(sf.Name + " is not a property")
				}

				if err := prop.SetSource(val); err != nil {
					return err
				}

			default:
				if prop == nil {
					continue
				}

				v := prop.Get()
				valt, vt := reflect.TypeOf(val), reflect.TypeOf(v)

				if v != nil && valt != vt {
					panic(fmt.Sprintf("cannot assign value %v of type %T to property %s of type %T", val, val, sf.Name, v))
				}
				if err := prop.Set(val); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (b *Builder) conditionOrProperty(data Property) interface{} {
	switch val := data.(type) {
	case bindData:
		if c, ok := b.knownCompositeConditions[val.expression]; ok {
			return c
		} else if conds := strings.Split(val.expression, "&&"); len(conds) > 1 {
			// This looks like a composite condition.
			for i, s := range conds {
				conds[i] = strings.TrimSpace(s)
			}

			var conditions []walk.Condition

			for _, cond := range conds {
				if p := b.property(cond); p != nil {
					conditions = append(conditions, p.(walk.Condition))
				} else if c, ok := conditionsByName[cond]; ok {
					conditions = append(conditions, c)
				} else {
					panic("unknown condition or property name: " + cond)
				}
			}

			var condition walk.Condition
			if len(conditions) > 1 {
				condition = walk.NewAllCondition(conditions...)
				b.knownCompositeConditions[val.expression] = condition
			} else {
				condition = conditions[0]
			}

			return condition
		}

		if p := b.property(val.expression); p != nil {
			return p
		}

		return conditionsByName[val.expression]

	case walk.Condition:
		return val
	}

	return nil
}

func (b *Builder) property(expression string) walk.Property {
	if parts := strings.Split(expression, "."); len(parts) == 2 {
		if sw, ok := b.name2Window[parts[0]]; ok {
			return sw.AsWindowBase().Property(parts[1])
		}
	}

	return nil
}
