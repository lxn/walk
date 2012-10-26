// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

func InitWidget(d Widget, w walk.Widget, customInit func() error) error {
	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	// Widget
	name, disabled, hidden, font, toolTipText, minSize, maxSize, stretchFactor, row, rowSpan, column, columnSpan, contextMenuActions, onKeyDown, onMouseDown, onMouseMove, onMouseUp, onSizeChanged := d.WidgetInfo()

	w.SetName(name)

	if err := w.SetToolTipText(toolTipText); err != nil {
		return err
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

	// Validatable
	if dva, ok := d.(Validatable); ok {
		wva := w.(walk.Validatable)

		if dvr := dva.ValidatableInfo(); dvr != nil {
			wvr, err := dvr.Create()
			if err != nil {
				return err
			}

			wva.SetValidator(wvr)
		}
	}

	// Container
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

			for _, child := range children {
				if err := child.Create(wc); err != nil {
					return err
				}
			}

			if db, err := dataBinder.create(); err != nil {
				return err
			} else if db != nil {
				wc.SetDataBinder(db)
			}
		}
	}

	// Custom
	if customInit != nil {
		if err := customInit(); err != nil {
			return err
		}
	}

	// Call Reset on DataBinder after customInit, so a Dialog gets a chance to first
	// wire up its DefaultButton to the CanSubmitChanged event of a DataBinder.
	if _, ok := d.(Container); ok {
		if wc, ok := w.(walk.Container); ok {
			db := wc.DataBinder()
			if db != nil {
				if err := db.Reset(); err != nil {
					return err
				}
			}
		}
	}

	// Widget continued
	w.SetEnabled(!disabled)
	w.SetVisible(!hidden)

	if font != nil {
		if f, err := font.Create(); err != nil {
			return err
		} else if f != nil {
			w.SetFont(f)
		}
	}

	succeeded = true

	return nil
}
