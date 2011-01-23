// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/ole32"
	. "walk/winapi/user32"
)

type webViewIOleInPlaceFrameCallbacks struct {
	QueryInterface       *syscall.Callback
	AddRef               *syscall.Callback
	Release              *syscall.Callback
	GetWindow            *syscall.Callback
	ContextSensitiveHelp *syscall.Callback
	GetBorder            *syscall.Callback
	RequestBorderSpace   *syscall.Callback
	SetBorderSpace       *syscall.Callback
	SetActiveObject      *syscall.Callback
	InsertMenus          *syscall.Callback
	SetMenu              *syscall.Callback
	RemoveMenus          *syscall.Callback
	SetStatusText        *syscall.Callback
	EnableModeless       *syscall.Callback
	TranslateAccelerator *syscall.Callback
}

var webViewIOleInPlaceFrameCbs = &webViewIOleInPlaceFrameCallbacks{
	syscall.NewCallback(webView_IOleInPlaceFrame_QueryInterface, 1+2),
	syscall.NewCallback(webView_IOleInPlaceFrame_AddRef, 1+0),
	syscall.NewCallback(webView_IOleInPlaceFrame_Release, 1+0),
	syscall.NewCallback(webView_IOleInPlaceFrame_GetWindow, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_ContextSensitiveHelp, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_GetBorder, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_RequestBorderSpace, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetBorderSpace, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetActiveObject, 1+2),
	syscall.NewCallback(webView_IOleInPlaceFrame_InsertMenus, 1+2),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetMenu, 1+3),
	syscall.NewCallback(webView_IOleInPlaceFrame_RemoveMenus, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetStatusText, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_EnableModeless, 1+1),
	syscall.NewCallback(webView_IOleInPlaceFrame_TranslateAccelerator, 1+2),
}

var webViewIOleInPlaceFrameVtbl *IOleInPlaceFrameVtbl

func init() {
	webViewIOleInPlaceFrameVtbl = &IOleInPlaceFrameVtbl{
		webViewIOleInPlaceFrameCbs.QueryInterface.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.AddRef.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.Release.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.GetWindow.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.ContextSensitiveHelp.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.GetBorder.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.RequestBorderSpace.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.SetBorderSpace.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.SetActiveObject.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.InsertMenus.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.SetMenu.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.RemoveMenus.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.SetStatusText.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.EnableModeless.ExtFnEntry(),
		webViewIOleInPlaceFrameCbs.TranslateAccelerator.ExtFnEntry(),
	}
}

type webViewIOleInPlaceFrame struct {
	IOleInPlaceFrame
	webView *WebView
}

func webView_IOleInPlaceFrame_QueryInterface(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_AddRef(args *uintptr) uintptr {
	return 1
}

func webView_IOleInPlaceFrame_Release(args *uintptr) uintptr {
	return 1
}

func webView_IOleInPlaceFrame_GetWindow(args *uintptr) uintptr {
	p := (*struct {
		inPlaceFrame *webViewIOleInPlaceFrame
		hwnd         *HWND
	})(unsafe.Pointer(args))

	*p.hwnd = p.inPlaceFrame.webView.hWnd

	return S_OK
}

func webView_IOleInPlaceFrame_ContextSensitiveHelp(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_GetBorder(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_RequestBorderSpace(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetBorderSpace(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetActiveObject(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_InsertMenus(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetMenu(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_RemoveMenus(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetStatusText(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_EnableModeless(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_TranslateAccelerator(args *uintptr) uintptr {
	return E_NOTIMPL
}
