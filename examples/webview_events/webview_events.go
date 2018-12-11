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
				OnDownloaded:              mainWin.webView_OnDownloaded,
				OnDocumentCompleted:       mainWin.webView_OnDocumentCompleted,
				OnNavigatedError:          mainWin.webView_OnNavigatedError,
				OnNewWindow:               mainWin.webView_OnNewWindow,
				OnQuitting:                mainWin.webView_OnQuitting,
				OnWindowClosing:           mainWin.webView_OnWindowClosing,
				OnStatusBarVisibleChanged: mainWin.webView_OnStatusBarVisibleChanged,
				OnTheaterModeChanged:      mainWin.webView_OnTheaterModeChanged,
				OnToolBarVisibleChanged:   mainWin.webView_OnToolBarVisibleChanged,
				OnBrowserVisibleChanged:   mainWin.webView_OnBrowserVisibleChanged,
				OnCanGoBackChanged:        mainWin.webView_OnCanGoBackChanged,
				OnCanGoForwardChanged:     mainWin.webView_OnCanGoForwardChanged,
				OnToolBarEnabledChanged:   mainWin.webView_OnToolBarEnabledChanged,
				OnProgressChanged:         mainWin.webView_OnProgressChanged,
				OnStatusTextChanged:       mainWin.webView_OnStatusTextChanged,
				OnDocumentTitleChanged:    mainWin.webView_OnDocumentTitleChanged,
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

func (mainWin *MainWin) webView_OnNavigating(eventData *walk.WebViewNavigatingEventData) {
	fmt.Printf("webView_OnNavigating\r\n")
	fmt.Printf("Url = %+v\r\n", eventData.Url())
	fmt.Printf("Flags = %+v\r\n", eventData.Flags())
	fmt.Printf("PostData = %+v\r\n", eventData.PostData())
	fmt.Printf("Headers = %+v\r\n", eventData.Headers())
	fmt.Printf("TargetFrameName = %+v\r\n", eventData.TargetFrameName())
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	// if you want to cancel
	//eventData.SetCanceled(true)
}

func (mainWin *MainWin) webView_OnNavigated(url string) {
	fmt.Printf("webView_OnNavigated\r\n")
	fmt.Printf("url = %+v\r\n", url)
}

func (mainWin *MainWin) webView_OnDownloading() {
	fmt.Printf("webView_OnDownloading\r\n")
}

func (mainWin *MainWin) webView_OnDownloaded() {
	fmt.Printf("webView_OnDownloaded\r\n")
}

func (mainWin *MainWin) webView_OnDocumentCompleted(url string) {
	fmt.Printf("webView_OnDocumentCompleted\r\n")
	fmt.Printf("url = %+v\r\n", url)
}

func (mainWin *MainWin) webView_OnNavigatedError(eventData *walk.WebViewNavigatedErrorEventData) {
	fmt.Printf("webView_OnNavigatedError\r\n")
	fmt.Printf("Url = %+v\r\n", eventData.Url())
	fmt.Printf("TargetFrameName = %+v\r\n", eventData.TargetFrameName())
	fmt.Printf("StatusCode = %+v\r\n", eventData.StatusCode())
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	// if you want to cancel
	//eventData.SetCanceled(true)
}

func (mainWin *MainWin) webView_OnNewWindow(eventData *walk.WebViewNewWindowEventData) {
	fmt.Printf("webView_OnNewWindow\r\n")
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	fmt.Printf("Flags = %+v\r\n", eventData.Flags())
	fmt.Printf("UrlContext = %+v\r\n", eventData.UrlContext())
	fmt.Printf("Url = %+v\r\n", eventData.Url())
	// if you want to cancel
	//eventData.SetCancel(true)
}

func (mainWin *MainWin) webView_OnQuitting() {
	fmt.Printf("webView_OnQuitting\r\n")
}

func (mainWin *MainWin) webView_OnWindowClosing(eventData *walk.WebViewWindowClosingEventData) {
	fmt.Printf("webView_OnWindowClosing\r\n")
	fmt.Printf("IsChildWindow = %+v\r\n", eventData.IsChildWindow())
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	// if you want to cancel
	//eventData.SetCancel(true)
}

func (mainWin *MainWin) webView_OnStatusBarVisibleChanged() {
	fmt.Printf("webView_OnStatusBarVisibleChanged\r\n")
	fmt.Printf("StatusBarVisible = %+v\r\n", mainWin.wv.StatusBarVisible())
}

func (mainWin *MainWin) webView_OnTheaterModeChanged() {
	fmt.Printf("webView_OnTheaterModeChanged\r\n")
	fmt.Printf("IsTheaterMode = %+v\r\n", mainWin.wv.IsTheaterMode())
}

func (mainWin *MainWin) webView_OnToolBarVisibleChanged() {
	fmt.Printf("webView_OnToolBarVisibleChanged\r\n")
	fmt.Printf("ToolBarVisible = %+v\r\n", mainWin.wv.ToolBarVisible())
}

func (mainWin *MainWin) webView_OnBrowserVisibleChanged() {
	fmt.Printf("webView_OnBrowserVisibleChanged\r\n")
	fmt.Printf("BrowserVisible = %+v\r\n", mainWin.wv.BrowserVisible())
}

func (mainWin *MainWin) webView_OnCanGoBackChanged() {
	fmt.Printf("webView_OnCanGoBackChanged\r\n")
	fmt.Printf("CanGoBack = %+v\r\n", mainWin.wv.CanGoBack())
}

func (mainWin *MainWin) webView_OnCanGoForwardChanged() {
	fmt.Printf("webView_OnCanGoForwardChanged\r\n")
	fmt.Printf("CanGoForward = %+v\r\n", mainWin.wv.CanGoForward())
}

func (mainWin *MainWin) webView_OnToolBarEnabledChanged() {
	fmt.Printf("webView_OnToolBarEnabledChanged\r\n")
	fmt.Printf("ToolBarEnabled = %+v\r\n", mainWin.wv.ToolBarEnabled())
}

func (mainWin *MainWin) webView_OnProgressChanged() {
	fmt.Printf("webView_OnProgressChanged\r\n")
	fmt.Printf("ProgressValue = %+v\r\n", mainWin.wv.ProgressValue())
	fmt.Printf("ProgressMax = %+v\r\n", mainWin.wv.ProgressMax())
}

func (mainWin *MainWin) webView_OnStatusTextChanged() {
	fmt.Printf("webView_OnStatusTextChanged\r\n")
	fmt.Printf("StatusText = %+v\r\n", mainWin.wv.StatusText())
}

func (mainWin *MainWin) webView_OnDocumentTitleChanged() {
	fmt.Printf("webView_OnDocumentTitleChanged\r\n")
	fmt.Printf("DocumentTitle = %+v\r\n", mainWin.wv.DocumentTitle())
}
