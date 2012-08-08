// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

import (
	"github.com/lxn/walk"
)

type MainWindow struct {
	*walk.MainWindow
	ui mainWindowUI
}

func runMainWindow() (int, error) {
	mw := new(MainWindow)
	if err := mw.init(); err != nil {
		return 0, err
	}
	defer mw.Dispose()

	foo := new(Foo)

	mw.ui.pushButton.Clicked().Attach(func() {
		res, err := runFooDialog(mw, foo)
		if err != nil {
			return
		}

		var s string
		if res == walk.DlgCmdOK {
			s = "accepted"
		} else {
			s = "canceled"
		}

		mw.ui.textEdit.SetText(fmt.Sprintf("%s, foo: %+v", s, foo))
	})

	mw.Show()

	return mw.Run(), nil
}
