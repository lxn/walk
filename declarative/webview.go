// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type WebView struct {
	// Window

	Background         Brush
	ContextMenuItems   []MenuItem
	Enabled            Property
	Font               Font
	MaxSize            Size
	MinSize            Size
	Name               string
	OnBoundsChanged    walk.EventHandler
	OnKeyDown          walk.KeyEventHandler
	OnKeyPress         walk.KeyEventHandler
	OnKeyUp            walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	Persistent         bool
	RightToLeftReading bool
	ToolTipText        Property
	Visible            Property

	// Widget

	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	GraphicsEffects    []walk.WidgetGraphicsEffect
	Row                int
	RowSpan            int
	StretchFactor      int

	// WebView

	AssignTo                   **walk.WebView
	NativeContextMenuEnabled   Property
	OnBrowserVisibleChanged    walk.WebViewBrowserVisibleChangedEventHandler
	OnCommandStateChanged      walk.WebViewCommandStateChangedEventHandler
	OnDocumentCompleted        walk.WebViewDocumentCompletedEventHandler
	OnDownloaded               walk.EventHandler
	OnDownloading              walk.EventHandler
	OnNativeContextMenuEnabled walk.EventHandler
	OnNavigated                walk.WebViewNavigatedEventHandler
	OnNavigatedError           walk.WebViewNavigatedErrorEventHandler
	OnNavigating               walk.WebViewNavigatingEventHandler
	OnNewWindow                walk.WebViewNewWindowEventHandler
	OnProgressChanged          walk.WebViewProgressChangedEventHandler
	OnShortcutsEnabled         walk.EventHandler
	OnStatusBarVisibleChanged  walk.WebViewStatusBarVisibleChangedEventHandler
	OnStatusTextChanged        walk.WebViewStatusTextChangedEventHandler
	OnTheaterModeChanged       walk.WebViewTheaterModeChangedEventHandler
	OnTitleChanged             walk.WebViewTitleChangedEventHandler
	OnToolBarVisibleChanged    walk.WebViewToolBarVisibleChangedEventHandler
	OnURLChanged               walk.EventHandler
	OnQuitting                 walk.EventHandler
	OnWindowClosing            walk.WebViewWindowClosingEventHandler
	URL                        Property
	ShortcutsEnabled           Property
}

func (wv WebView) Create(builder *Builder) error {
	w, err := walk.NewWebView(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(wv, w, func() error {
		if wv.OnBrowserVisibleChanged != nil {
			w.BrowserVisibleChanged().Attach(wv.OnBrowserVisibleChanged)
		}
		if wv.OnCommandStateChanged != nil {
			w.CommandStateChanged().Attach(wv.OnCommandStateChanged)
		}
		if wv.OnDocumentCompleted != nil {
			w.DocumentCompleted().Attach(wv.OnDocumentCompleted)
		}
		if wv.OnDownloaded != nil {
			w.Downloaded().Attach(wv.OnDownloaded)
		}
		if wv.OnDownloading != nil {
			w.Downloading().Attach(wv.OnDownloading)
		}
		if wv.OnNativeContextMenuEnabled != nil {
			w.NativeContextMenuEnabledChanged().Attach(wv.OnNativeContextMenuEnabled)
		}
		if wv.OnNavigated != nil {
			w.Navigated().Attach(wv.OnNavigated)
		}
		if wv.OnNavigatedError != nil {
			w.NavigatedError().Attach(wv.OnNavigatedError)
		}
		if wv.OnNavigating != nil {
			w.Navigating().Attach(wv.OnNavigating)
		}
		if wv.OnNewWindow != nil {
			w.NewWindow().Attach(wv.OnNewWindow)
		}
		if wv.OnProgressChanged != nil {
			w.ProgressChanged().Attach(wv.OnProgressChanged)
		}
		if wv.OnURLChanged != nil {
			w.URLChanged().Attach(wv.OnURLChanged)
		}
		if wv.OnShortcutsEnabled != nil {
			w.ShortcutsEnabledChanged().Attach(wv.OnShortcutsEnabled)
		}
		if wv.OnStatusBarVisibleChanged != nil {
			w.StatusBarVisibleChanged().Attach(wv.OnStatusBarVisibleChanged)
		}
		if wv.OnStatusTextChanged != nil {
			w.StatusTextChanged().Attach(wv.OnStatusTextChanged)
		}
		if wv.OnTitleChanged != nil {
			w.TitleChanged().Attach(wv.OnTitleChanged)
		}
		if wv.OnTheaterModeChanged != nil {
			w.TheaterModeChanged().Attach(wv.OnTheaterModeChanged)
		}
		if wv.OnToolBarVisibleChanged != nil {
			w.ToolBarVisibleChanged().Attach(wv.OnToolBarVisibleChanged)
		}
		if wv.OnQuitting != nil {
			w.Quitting().Attach(wv.OnQuitting)
		}
		if wv.OnWindowClosing != nil {
			w.WindowClosing().Attach(wv.OnWindowClosing)
		}

		if wv.AssignTo != nil {
			*wv.AssignTo = w
		}

		return nil
	})
}
