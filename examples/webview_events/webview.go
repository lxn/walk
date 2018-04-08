// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/lxn/win"
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
				AssignTo:           &mainWin.wv,
				Name:               "wv",
				URL:                "https://github.com/lxn/walk",
				ShortcutsEnabled:   true,
				ContextMenuEnabled: true,
				BeforeNavigate2:    mainWin.webView_BeforeNavigate2,
				NavigateComplete2:  mainWin.webView_NavigateComplete2,
				DownloadBegin:      mainWin.webView_DownloadBegin,
				DocumentComplete:   mainWin.webView_DocumentComplete,
				NavigateError:      mainWin.webView_NavigateError,
				NewWindow3:         mainWin.webView_NewWindow3,
				OnQuit:             mainWin.webView_OnQuit,
				WindowClosing:      mainWin.webView_WindowClosing,
				OnStatusBar:        mainWin.webView_OnStatusBar,
				OnTheaterMode:      mainWin.webView_OnTheaterMode,
				OnToolBar:          mainWin.webView_OnToolBar,
				OnVisible:          mainWin.webView_OnVisible,
				CommandStateChange: mainWin.webView_CommandStateChange,
				ProgressChange:     mainWin.webView_ProgressChange,
				StatusTextChange:   mainWin.webView_StatusTextChange,
				TitleChange:        mainWin.webView_TitleChange,
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

func (mainWin *MainWin) webView_BeforeNavigate2(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	flags *win.VARIANT,
	targetFrameName *win.VARIANT,
	postData *win.VARIANT,
	headers *win.VARIANT,
	cancel *win.VARIANT_BOOL) {

	fmt.Printf("webView_BeforeNavigate2\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
	fmt.Printf("flags = %+v\r\n", flags)
	if flags != nil {
		fmt.Printf("    flags = %+v\r\n", flags.LVal())
	}
	fmt.Printf("targetFrameName = %+v\r\n", targetFrameName)
	if targetFrameName != nil && targetFrameName.BstrVal() != nil {
		fmt.Printf("  targetFrameName = %+v\r\n", win.BSTRToString(targetFrameName.BstrVal()))
	}
	fmt.Printf("postData = %+v\r\n", postData)
	if postData != nil {
		fmt.Printf("    postData = %+v\r\n", postData.PVarVal())
	}
	fmt.Printf("headers = %+v\r\n", headers)
	if headers != nil && headers.BstrVal() != nil {
		fmt.Printf("  headers = %+v\r\n", win.BSTRToString(headers.BstrVal()))
	}
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("  *cancel = %+v\r\n", *cancel)
	}

}

func (mainWin *MainWin) webView_NavigateComplete2(pDisp *win.IDispatch, url *win.VARIANT) {
	fmt.Printf("webView_NavigateComplete2\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
}

func (mainWin *MainWin) webView_DownloadBegin() {
	fmt.Printf("webView_DownloadBegin\r\n")
}

func (mainWin *MainWin) webView_DownloadComplete() {
	fmt.Printf("webView_DownloadComplete\r\n")
}

func (mainWin *MainWin) webView_DocumentComplete(pDisp *win.IDispatch, url *win.VARIANT) {
	fmt.Printf("webView_DocumentComplete\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
}

func (mainWin *MainWin) webView_NavigateError(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	targetFrameName *win.VARIANT,
	statusCode *win.VARIANT,
	cancel *win.VARIANT_BOOL) {

	fmt.Printf("webView_NavigateError\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
	fmt.Printf("targetFrameName = %+v\r\n", targetFrameName)
	if targetFrameName != nil && targetFrameName.BstrVal() != nil {
		fmt.Printf("  targetFrameName = %+v\r\n", win.BSTRToString(targetFrameName.BstrVal()))
	}
	fmt.Printf("statusCode = %+v\r\n", statusCode)
	if statusCode != nil {
		fmt.Printf("    statusCode = %+v\r\n", statusCode.LVal())
	}
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("  *cancel = %+v\r\n", *cancel)
	}
}

func (mainWin *MainWin) webView_NewWindow3(
	ppDisp **win.IDispatch,
	cancel *win.VARIANT_BOOL,
	dwFlags uint32,
	bstrUrlContext *uint16,
	bstrUrl *uint16) {

	fmt.Printf("webView_NewWindow3\r\n")
	fmt.Printf("ppDisp = %+v\r\n", ppDisp)
	if ppDisp != nil {
		fmt.Printf("    *ppDisp = %+v\r\n", *ppDisp)
	}
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("  *cancel = %+v\r\n", *cancel)
	}
	fmt.Printf("dwFlags = %+v\r\n", dwFlags)
	fmt.Printf("bstrUrlContext = %+v\r\n", bstrUrlContext)
	if bstrUrlContext != nil {
		fmt.Printf("  bstrUrlContext = %+v\r\n", win.BSTRToString(bstrUrlContext))
	}
	fmt.Printf("bstrUrl = %+v\r\n", bstrUrl)
	if bstrUrl != nil {
		fmt.Printf("  bstrUrl = %+v\r\n", win.BSTRToString(bstrUrl))
	}
}

func (mainWin *MainWin) webView_OnQuit() {
	fmt.Printf("webView_OnQuit\r\n")
}

func (mainWin *MainWin) webView_WindowClosing(bIsChildWindow win.VARIANT_BOOL, cancel *win.VARIANT_BOOL) {
	fmt.Printf("webView_WindowClosing\r\n")
	fmt.Printf("bIsChildWindow = %+v\r\n", bIsChildWindow)
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("*cancel = %+v\r\n", *cancel)
	}
}

func (mainWin *MainWin) webView_OnStatusBar(statusBar win.VARIANT_BOOL) {
	fmt.Printf("webView_OnStatusBar\r\n")
	fmt.Printf("statusBar = %+v\r\n", statusBar)
}

func (mainWin *MainWin) webView_OnTheaterMode(theaterMode win.VARIANT_BOOL) {
	fmt.Printf("webView_OnTheaterMode\r\n")
	fmt.Printf("theaterMode = %+v\r\n", theaterMode)
}

func (mainWin *MainWin) webView_OnToolBar(toolBar win.VARIANT_BOOL) {
	fmt.Printf("webView_OnToolBar\r\n")
	fmt.Printf("toolBar = %+v\r\n", toolBar)
}

func (mainWin *MainWin) webView_OnVisible(vVisible win.VARIANT_BOOL) {
	fmt.Printf("webView_OnVisible\r\n")
	fmt.Printf("vVisible = %+v\r\n", vVisible)
}

func (mainWin *MainWin) webView_CommandStateChange(command int32, enable win.VARIANT_BOOL) {
	fmt.Printf("webView_CommandStateChange\r\n")
	fmt.Printf("command = %+v\r\n", command)
	fmt.Printf("enable = %+v\r\n", enable)
}

func (mainWin *MainWin) webView_ProgressChange(nProgress int32, nProgressMax int32) {
	fmt.Printf("webView_ProgressChange\r\n")
	fmt.Printf("nProgress = %+v\r\n", nProgress)
	fmt.Printf("nProgressMax = %+v\r\n", nProgressMax)
}

func (mainWin *MainWin) webView_StatusTextChange(sText *uint16) {
	fmt.Printf("webView_StatusTextChange\r\n")
	fmt.Printf("sText = %+v\r\n", sText)
	if sText != nil {
		fmt.Printf("  sText = %+v\r\n", win.BSTRToString(sText))
	}
}

func (mainWin *MainWin) webView_TitleChange(sText *uint16) {
	fmt.Printf("webView_TitleChange\r\n")
	fmt.Printf("sText = %+v\r\n", sText)
	if sText != nil {
		fmt.Printf("  sText = %+v\r\n", win.BSTRToString(sText))
	}
}
