// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"log"
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
	syscall.NewCallback(webView_IOleInPlaceFrame_QueryInterface, 4+4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_AddRef, 4),
	syscall.NewCallback(webView_IOleInPlaceFrame_Release, 4),
	syscall.NewCallback(webView_IOleInPlaceFrame_GetWindow, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_ContextSensitiveHelp, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_GetBorder, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_RequestBorderSpace, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetBorderSpace, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetActiveObject, 4+4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_InsertMenus, 4+4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetMenu, 4+4+4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_RemoveMenus, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_SetStatusText, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_EnableModeless, 4+4),
	syscall.NewCallback(webView_IOleInPlaceFrame_TranslateAccelerator, 4+4+2),
}

var webViewIOleInPlaceFrameVtbl *IOleInPlaceFrameVtbl

func init() {
	webViewIOleInPlaceFrameVtbl = &IOleInPlaceFrameVtbl{
		uintptr(webViewIOleInPlaceFrameCbs.QueryInterface.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.AddRef.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.Release.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.GetWindow.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.ContextSensitiveHelp.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.GetBorder.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.RequestBorderSpace.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.SetBorderSpace.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.SetActiveObject.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.InsertMenus.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.SetMenu.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.RemoveMenus.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.SetStatusText.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.EnableModeless.ExtFnEntry()),
		uintptr(webViewIOleInPlaceFrameCbs.TranslateAccelerator.ExtFnEntry()),
	}
}

type webViewIOleInPlaceFrame struct {
	IOleInPlaceFrame
	webView *WebView
}

func webView_IOleInPlaceFrame_QueryInterface(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_QueryInterface")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_AddRef(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_AddRef")

	return 1
}

func webView_IOleInPlaceFrame_Release(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_Release")

	return 1
}

func webView_IOleInPlaceFrame_GetWindow(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_GetWindow")

	p := (*struct {
		inPlaceFrame *webViewIOleInPlaceFrame
		hwnd         *HWND
	})(unsafe.Pointer(args))

	*p.hwnd = p.inPlaceFrame.webView.hWnd

	return S_OK
}

func webView_IOleInPlaceFrame_ContextSensitiveHelp(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_ContextSensitiveHelp")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_GetBorder(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_GetBorder")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_RequestBorderSpace(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_RequestBorderSpace")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetBorderSpace(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_SetBorderSpace")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetActiveObject(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_SetActiveObject")

	return S_OK
}

func webView_IOleInPlaceFrame_InsertMenus(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_InsertMenus")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetMenu(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_SetMenu")

	return S_OK
}

func webView_IOleInPlaceFrame_RemoveMenus(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_RemoveMenus")

	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetStatusText(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_SetStatusText")

	return S_OK
}

func webView_IOleInPlaceFrame_EnableModeless(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_EnableModeless")

	return S_OK
}

func webView_IOleInPlaceFrame_TranslateAccelerator(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceFrame_TranslateAccelerator")

	return E_NOTIMPL
}
