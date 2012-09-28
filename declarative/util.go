// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

func initWidget(decl interface{}, widget walk.Widget) error {
	// Container
	if dc, ok := decl.(Container); ok {
		if wc, ok := widget.(walk.Container); ok {
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
