// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
)

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	mw, _ := walk.NewMainWindow()
	mw.SetTitle("Walk WebView Example")
	mw.SetLayout(walk.NewVBoxLayout())
	mw.SetMinMaxSize(walk.Size{600, 400}, walk.Size{})
	mw.SetSize(walk.Size{800, 600})

	var webView *walk.WebView

	urlLineEdit, _ := walk.NewLineEdit(mw)
	urlLineEdit.ReturnPressed().Attach(func() {
		webView.SetURL(urlLineEdit.Text())
	})

	webView, _ = walk.NewWebView(mw)
	webView.SetURL("http://golang.org")

	mw.Show()
	mw.Run()
}
