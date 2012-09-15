// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
)

type MainWindow struct {
	*walk.MainWindow
	urlLineEdit *walk.LineEdit
	webView     *walk.WebView
}

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	mainWnd, _ := walk.NewMainWindow()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.SetTitle("Walk WebView Example")
	mw.SetLayout(walk.NewVBoxLayout())

	mw.urlLineEdit, _ = walk.NewLineEdit(mw)
	mw.urlLineEdit.ReturnPressed().Attach(func() {
		mw.webView.SetURL(mw.urlLineEdit.Text())
	})

	mw.webView, _ = walk.NewWebView(mw)

	mw.webView.SetURL("http://golang.org")

	mw.SetMinMaxSize(walk.Size{600, 400}, walk.Size{})
	mw.SetSize(walk.Size{800, 600})
	mw.Show()

	mw.Run()
}
