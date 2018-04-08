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

	AssignTo             **walk.WebView
	OnURLChanged         walk.EventHandler
	URL                  Property
	OnShortcutsEnabled   walk.EventHandler
	ShortcutsEnabled     Property
	OnContextMenuEnabled walk.EventHandler
	ContextMenuEnabled   Property
	BeforeNavigate2      walk.WvBeforeNavigate2EventHandler
	NavigateComplete2    walk.WvNavigateComplete2EventHandler
	DownloadBegin        walk.WvDownloadBeginEventHandler
	DownloadComplete     walk.WvDownloadCompleteEventHandler
	DocumentComplete     walk.WvDocumentCompleteEventHandler
	NavigateError        walk.WvNavigateErrorEventHandler
	NewWindow3           walk.WvNewWindow3EventHandler
	OnQuit               walk.WvOnQuitEventHandler
	WindowClosing        walk.WvWindowClosingEventHandler
	OnStatusBar          walk.WvOnStatusBarEventHandler
	OnTheaterMode        walk.WvOnTheaterModeEventHandler
	OnToolBar            walk.WvOnToolBarEventHandler
	OnVisible            walk.WvOnVisibleEventHandler
	CommandStateChange   walk.WvCommandStateChangeEventHandler
	ProgressChange       walk.WvProgressChangeEventHandler
	StatusTextChange     walk.WvStatusTextChangeEventHandler
	TitleChange          walk.WvTitleChangeEventHandler
}

func (wv WebView) Create(builder *Builder) error {
	w, err := walk.NewWebView(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(wv, w, func() error {
		if wv.OnURLChanged != nil {
			w.URLChanged().Attach(wv.OnURLChanged)
		}
		if wv.OnShortcutsEnabled != nil {
			w.ShortcutsEnabledChanged().Attach(wv.OnShortcutsEnabled)
		}
		if wv.OnContextMenuEnabled != nil {
			w.ContextMenuEnabledChanged().Attach(wv.OnContextMenuEnabled)
		}
		if wv.BeforeNavigate2 != nil {
			w.BeforeNavigate2().Attach(wv.BeforeNavigate2)
		}
		if wv.NavigateComplete2 != nil {
			w.NavigateComplete2().Attach(wv.NavigateComplete2)
		}
		if wv.DownloadBegin != nil {
			w.DownloadBegin().Attach(wv.DownloadBegin)
		}
		if wv.DownloadComplete != nil {
			w.DownloadComplete().Attach(wv.DownloadComplete)
		}
		if wv.DocumentComplete != nil {
			w.DocumentComplete().Attach(wv.DocumentComplete)
		}
		if wv.NavigateError != nil {
			w.NavigateError().Attach(wv.NavigateError)
		}
		if wv.NewWindow3 != nil {
			w.NewWindow3().Attach(wv.NewWindow3)
		}
		if wv.OnQuit != nil {
			w.OnQuit().Attach(wv.OnQuit)
		}
		if wv.WindowClosing != nil {
			w.WindowClosing().Attach(wv.WindowClosing)
		}
		if wv.OnStatusBar != nil {
			w.OnStatusBar().Attach(wv.OnStatusBar)
		}
		if wv.OnTheaterMode != nil {
			w.OnTheaterMode().Attach(wv.OnTheaterMode)
		}
		if wv.OnToolBar != nil {
			w.OnToolBar().Attach(wv.OnToolBar)
		}
		if wv.OnVisible != nil {
			w.OnVisible().Attach(wv.OnVisible)
		}
		if wv.CommandStateChange != nil {
			w.CommandStateChange().Attach(wv.CommandStateChange)
		}
		if wv.ProgressChange != nil {
			w.ProgressChange().Attach(wv.ProgressChange)
		}
		if wv.StatusTextChange != nil {
			w.StatusTextChange().Attach(wv.StatusTextChange)
		}
		if wv.TitleChange != nil {
			w.TitleChange().Attach(wv.TitleChange)
		}

		if wv.AssignTo != nil {
			*wv.AssignTo = w
		}

		return nil
	})
}
