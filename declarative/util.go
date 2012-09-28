// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

func initWidget(d Widget, w walk.Widget) error {
	// Widget
	if p := w.Parent(); p != nil {
		hStretch, vStretch, row, rowSpan, col, colSpan := d.LayoutParams()

		switch l := p.Layout().(type) {
		case *walk.BoxLayout:
			var sf int
			if l.Orientation() == walk.Horizontal {
				sf = hStretch
			} else {
				sf = vStretch
			}
			if sf < 1 {
				sf = 1
			}
			if err := l.SetStretchFactor(w, sf); err != nil {
				return err
			}

		case *walk.GridLayout:
			cs := colSpan
			if cs < 1 {
				cs = 1
			}
			rs := rowSpan
			if rs < 1 {
				rs = 1
			}
			r := walk.Rectangle{col, row, cs, rs}

			if err := l.SetRange(w, r); err != nil {
				return err
			}
		}
	}

	// Fonter
	if fonter, ok := d.(Fonter); ok {
		if f, err := fonter.Font_().Create(); err != nil {
			return err
		} else if f != nil {
			w.SetFont(f)
		}
	}

	// Container
	if dc, ok := d.(Container); ok {
		if wc, ok := w.(walk.Container); ok {
			if dl := dc.Layout_(); dl != nil {
				l, err := dl.Create()
				if err != nil {
					return err
				}

				if err := wc.SetLayout(l); err != nil {
					return err
				}
			}

			for _, child := range dc.Children_() {
				if _, err := child.Create(wc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
