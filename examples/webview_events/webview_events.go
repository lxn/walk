// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MainWin struct {
	*walk.MainWindow
	le *walk.LineEdit
	wv *walk.WebView
}

func main() {
	mainWin, err := NewMainWin()
	if err != nil {
		log.Fatal(err)
	}

	mainWin.Run()
}

func NewMainWin() (*MainWin, error) {
	mainWin := new(MainWin)

	err := MainWindow{
		AssignTo: &mainWin.MainWindow,
		Icon:     Bind("'../img/' + icon(mainWin.wv.URL) + '.ico'"),
		Title:    "Walk WebView Example (With Events Printing)",
		MinSize:  Size{800, 600},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			LineEdit{
				AssignTo: &mainWin.le,
				Text:     Bind("wv.URL"),
				OnKeyDown: func(key walk.Key) {
					if key == walk.KeyReturn {
						mainWin.wv.SetURL(mainWin.le.Text())
					}
				},
			},
			WebView{
				AssignTo:                  &mainWin.wv,
				Name:                      "wv",
				URL:                       "https://github.com/lxn/walk",
				ShortcutsEnabled:          true,
				NativeContextMenuEnabled:  true,
				OnNavigating:              mainWin.webView_OnNavigating,
				OnNavigated:               mainWin.webView_OnNavigated,
				OnDownloading:             mainWin.webView_OnDownloading,
				OnDocumentCompleted:       mainWin.webView_OnDocumentCompleted,
				OnNavigatedError:          mainWin.webView_OnNavigatedError,
				OnNewWindow:               mainWin.webView_OnNewWindow,
				OnQuitting:                mainWin.webView_OnQuitting,
				OnWindowClosing:           mainWin.webView_OnWindowClosing,
				OnStatusBarVisibleChanged: mainWin.webView_OnStatusBarVisibleChanged,
				OnTheaterModeChanged:      mainWin.webView_OnTheaterModeChanged,
				OnToolBarVisibleChanged:   mainWin.webView_OnToolBarVisibleChanged,
				OnBrowserVisibleChanged:   mainWin.webView_OnBrowserVisibleChanged,
				OnCommandStateChanged:     mainWin.webView_OnCommandStateChanged,
				OnProgressChanged:         mainWin.webView_OnProgressChanged,
				OnStatusTextChanged:       mainWin.webView_OnStatusTextChanged,
				OnTitleChanged:            mainWin.webView_OnTitleChanged,
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"icon": func(args ...interface{}) (interface{}, error) {
				if strings.HasPrefix(args[0].(string), "https") {
					return "check", nil
				}

				return "stop", nil
			},
		},
	}.Create()

	return mainWin, err
}

func (mainWin *MainWin) webView_OnNavigating(arg *walk.WebViewNavigatingArg) {
	fmt.Printf("webView_OnNavigating\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
	fmt.Printf("Flags = %+v\r\n", arg.Flags())
	fmt.Printf("Headers = %+v\r\n", arg.Headers())
	fmt.Printf("TargetFrameName = %+v\r\n", arg.TargetFrameName())
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	// if you want to cancel
	//arg.SetCancel(true)
}

func (mainWin *MainWin) webView_OnNavigated(arg *walk.WebViewNavigatedEventArg) {
	fmt.Printf("webView_OnNavigated\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
}

func (mainWin *MainWin) webView_OnDownloading() {
	fmt.Printf("webView_OnDownloading\r\n")
}

func (mainWin *MainWin) webView_OnDownloaded() {
	fmt.Printf("webView_OnDownloaded\r\n")
}

func (mainWin *MainWin) webView_OnDocumentCompleted(arg *walk.WebViewDocumentCompletedEventArg) {
	fmt.Printf("webView_OnDocumentCompleted\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
}

func (mainWin *MainWin) webView_OnNavigatedError(arg *walk.WebViewNavigatedErrorEventArg) {
	fmt.Printf("webView_OnNavigatedError\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
	fmt.Printf("TargetFrameName = %+v\r\n", arg.TargetFrameName())
	fmt.Printf("StatusCode = %+v\r\n", arg.StatusCode())
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	// if you want to cancel
	//arg.SetCancel(true)
}

func (mainWin *MainWin) webView_OnNewWindow(arg *walk.WebViewNewWindowEventArg) {
	fmt.Printf("webView_OnNewWindow\r\n")
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	fmt.Printf("Flags = %+v\r\n", arg.Flags())
	fmt.Printf("UrlContext = %+v\r\n", arg.UrlContext())
	fmt.Printf("Url = %+v\r\n", arg.Url())
	// if you want to cancel
	//arg.SetCancel(true)
}

func (mainWin *MainWin) webView_OnQuitting() {
	fmt.Printf("webView_OnQuitting\r\n")
}

func (mainWin *MainWin) webView_OnWindowClosing(arg *walk.WebViewWindowClosingEventArg) {
	fmt.Printf("webView_OnWindowClosing\r\n")
	fmt.Printf("IsChildWindow = %+v\r\n", arg.IsChildWindow())
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	// if you want to cancel
	//arg.SetCancel(true)
}

func (mainWin *MainWin) webView_OnStatusBarVisibleChanged(arg *walk.WebViewStatusBarVisibleChangedEventArg) {
	fmt.Printf("webView_OnStatusBarVisibleChanged\r\n")
	fmt.Printf("Visible = %+v\r\n", arg.Visible())
}

func (mainWin *MainWin) webView_OnTheaterModeChanged(arg *walk.WebViewTheaterModeChangedEventArg) {
	fmt.Printf("webView_OnTheaterModeChanged\r\n")
	fmt.Printf("IsTheaterMode = %+v\r\n", arg.IsTheaterMode())
}

func (mainWin *MainWin) webView_OnToolBarVisibleChanged(arg *walk.WebViewToolBarVisibleChangedEventArg) {
	fmt.Printf("webView_OnToolBarVisibleChanged\r\n")
	fmt.Printf("Visible = %+v\r\n", arg.Visible())
}

func (mainWin *MainWin) webView_OnBrowserVisibleChanged(arg *walk.WebViewBrowserVisibleChangedEventArg) {
	fmt.Printf("webView_OnBrowserVisibleChanged\r\n")
	fmt.Printf("Visible = %+v\r\n", arg.Visible())
}

func (mainWin *MainWin) webView_OnCommandStateChanged(arg *walk.WebViewCommandStateChangedEventArg) {
	fmt.Printf("webView_OnCommandStateChanged\r\n")
	fmt.Printf("Command = %+v\r\n", arg.Command())
	fmt.Printf("Enabled = %+v\r\n", arg.Enabled())
}

func (mainWin *MainWin) webView_OnProgressChanged(arg *walk.WebViewProgressChangedEventArg) {
	fmt.Printf("webView_OnProgressChanged\r\n")
	fmt.Printf("Progress = %+v\r\n", arg.Progress())
	fmt.Printf("ProgressMax = %+v\r\n", arg.ProgressMax())
}

func (mainWin *MainWin) webView_OnStatusTextChanged(arg *walk.WebViewStatusTextChangedEventArg) {
	fmt.Printf("webView_OnStatusTextChanged\r\n")
	fmt.Printf("StatusText = %+v\r\n", arg.StatusText())
}

func (mainWin *MainWin) webView_OnTitleChanged(arg *walk.WebViewTitleChangedEventArg) {
	fmt.Printf("webView_OnTitleChanged\r\n")
	fmt.Printf("Title = %+v\r\n", arg.Title())
}
