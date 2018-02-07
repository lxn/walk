// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	walk.Resources.SetRootDirPath("../img")

	mw := new(AppMainWindow)

	cfg := &MultiPageMainWindowConfig{
		Name:    "mainWindow",
		MinSize: Size{600, 400},
		MenuItems: []MenuItem{
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "About",
						OnTriggered: func() { mw.aboutAction_Triggered() },
					},
				},
			},
		},
		OnCurrentPageChanged: func() {
			mw.updateTitle(mw.CurrentPageTitle())
		},
		PageCfgs: []PageConfig{
			{"Foo", "document-new.png", newFooPage},
			{"Bar", "document-properties.png", newBarPage},
			{"Baz", "system-shutdown.png", newBazPage},
		},
	}

	mpmw, err := NewMultiPageMainWindow(cfg)
	if err != nil {
		panic(err)
	}

	mw.MultiPageMainWindow = mpmw

	mw.updateTitle(mw.CurrentPageTitle())

	mw.Run()
}

type AppMainWindow struct {
	*MultiPageMainWindow
}

func (mw *AppMainWindow) updateTitle(prefix string) {
	var buf bytes.Buffer

	if prefix != "" {
		buf.WriteString(prefix)
		buf.WriteString(" - ")
	}

	buf.WriteString("Walk Multiple Pages Example")

	mw.SetTitle(buf.String())
}

func (mw *AppMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw,
		"About Walk Multiple Pages Example",
		"An example that demonstrates a main window that supports multiple pages.",
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

type FooPage struct {
	*walk.Composite
}

func newFooPage(parent walk.Container) (Page, error) {
	p := new(FooPage)

	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "fooPage",
		Layout:   HBox{},
		Children: []Widget{
			HSpacer{},
			Label{Text: "I'm the Foo page"},
			HSpacer{},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}

	return p, nil
}

type BarPage struct {
	*walk.Composite
}

func newBarPage(parent walk.Container) (Page, error) {
	p := new(BarPage)

	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "barPage",
		Layout:   HBox{},
		Children: []Widget{
			HSpacer{},
			Label{Text: "I'm the Bar page"},
			HSpacer{},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}

	return p, nil
}

type BazPage struct {
	*walk.Composite
}

func newBazPage(parent walk.Container) (Page, error) {
	p := new(BazPage)

	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "bazPage",
		Layout:   HBox{},
		Children: []Widget{
			HSpacer{},
			Label{Text: "I'm the Baz page"},
			HSpacer{},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}

	return p, nil
}
