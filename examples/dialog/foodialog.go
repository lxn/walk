// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
)

type FooDialog struct {
	*walk.Dialog
	ui fooDialogUI
}

func runFooDialog(owner walk.RootWidget, foo *Foo) (int, error) {
	dlg := new(FooDialog)
	if err := dlg.init(owner); err != nil {
		return 0, err
	}

	if err := dlg.SetDefaultButton(dlg.ui.acceptPushButton); err != nil {
		return 0, err
	}

	dlg.ui.acceptPushButton.Clicked().Attach(func() {
		// Populate foo from widgets.
		foo.FilePath = dlg.ui.filePathLineEdit.Text()
		foo.Bar = barModel.items[dlg.ui.barComboBox.CurrentIndex()]
		foo.Number = dlg.ui.numberNumberEdit.Value()
		foo.Date = dlg.ui.dateDateEdit.Value()
		foo.Memo = dlg.ui.memoTextEdit.Text()

		dlg.Accept()
	})

	if err := dlg.SetCancelButton(dlg.ui.cancelPushButton); err != nil {
		return 0, err
	}

	dlg.ui.cancelPushButton.Clicked().Attach(func() {
		dlg.Cancel()
	})

	dlg.ui.filePathToolButton.Clicked().Attach(func() {
		d := &walk.FileDialog{}

		d.FilePath = dlg.ui.filePathLineEdit.Text()
		d.Filter = "Text Files (*.txt)|*.txt"
		d.Title = "Choose a text file."

		if ok, _ := d.ShowOpen(dlg); !ok {
			return
		}

		dlg.ui.filePathLineEdit.SetText(d.FilePath)
	})

	dlg.ui.barComboBox.CurrentIndexChanged().Attach(func() {
		// Enable accept button only if a Bar was selected.
		dlg.ui.acceptPushButton.SetEnabled(dlg.ui.barComboBox.CurrentIndex() > -1)
	})

	dlg.ui.barComboBox.SetModel(barModel)

	// Populate widgets from foo.
	dlg.ui.filePathLineEdit.SetText(foo.FilePath)
	dlg.ui.barComboBox.SetCurrentIndex(barModel.Index(foo.Bar))
	dlg.ui.numberNumberEdit.SetValue(foo.Number)
	dlg.ui.dateDateEdit.SetValue(foo.Date)
	dlg.ui.memoTextEdit.SetText(foo.Memo)

	return dlg.Run(), nil
}
