// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
)

import . "github.com/lxn/go-winapi"

var webViewIOleInPlaceFrameVtbl *IOleInPlaceFrameVtbl

func init() {
	webViewIOleInPlaceFrameVtbl = &IOleInPlaceFrameVtbl{
		syscall.NewCallback(webView_IOleInPlaceFrame_QueryInterface),
		syscall.NewCallback(webView_IOleInPlaceFrame_AddRef),
		syscall.NewCallback(webView_IOleInPlaceFrame_Release),
		syscall.NewCallback(webView_IOleInPlaceFrame_GetWindow),
		syscall.NewCallback(webView_IOleInPlaceFrame_ContextSensitiveHelp),
		syscall.NewCallback(webView_IOleInPlaceFrame_GetBorder),
		syscall.NewCallback(webView_IOleInPlaceFrame_RequestBorderSpace),
		syscall.NewCallback(webView_IOleInPlaceFrame_SetBorderSpace),
		syscall.NewCallback(webView_IOleInPlaceFrame_SetActiveObject),
		syscall.NewCallback(webView_IOleInPlaceFrame_InsertMenus),
		syscall.NewCallback(webView_IOleInPlaceFrame_SetMenu),
		syscall.NewCallback(webView_IOleInPlaceFrame_RemoveMenus),
		syscall.NewCallback(webView_IOleInPlaceFrame_SetStatusText),
		syscall.NewCallback(webView_IOleInPlaceFrame_EnableModeless),
		syscall.NewCallback(webView_IOleInPlaceFrame_TranslateAccelerator),
	}
}

type webViewIOleInPlaceFrame struct {
	IOleInPlaceFrame
	webView *WebView
}

func webView_IOleInPlaceFrame_QueryInterface(inPlaceFrame *webViewIOleInPlaceFrame, riid REFIID, ppvObj *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_AddRef(inPlaceFrame *webViewIOleInPlaceFrame) uintptr {
	return 1
}

func webView_IOleInPlaceFrame_Release(inPlaceFrame *webViewIOleInPlaceFrame) uintptr {
	return 1
}

func webView_IOleInPlaceFrame_GetWindow(inPlaceFrame *webViewIOleInPlaceFrame, lphwnd *HWND) uintptr {
	*lphwnd = inPlaceFrame.webView.hWnd

	return S_OK
}

func webView_IOleInPlaceFrame_ContextSensitiveHelp(inPlaceFrame *webViewIOleInPlaceFrame, fEnterMode BOOL) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_GetBorder(inPlaceFrame *webViewIOleInPlaceFrame, lprectBorder *RECT) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_RequestBorderSpace(inPlaceFrame *webViewIOleInPlaceFrame, pborderwidths uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetBorderSpace(inPlaceFrame *webViewIOleInPlaceFrame, pborderwidths uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetActiveObject(inPlaceFrame *webViewIOleInPlaceFrame, pActiveObject uintptr, pszObjName *uint16) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_InsertMenus(inPlaceFrame *webViewIOleInPlaceFrame, hmenuShared HMENU, lpMenuWidths uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetMenu(inPlaceFrame *webViewIOleInPlaceFrame, hmenuShared HMENU, holemenu HMENU, hwndActiveObject HWND) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_RemoveMenus(inPlaceFrame *webViewIOleInPlaceFrame, hmenuShared HMENU) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceFrame_SetStatusText(inPlaceFrame *webViewIOleInPlaceFrame, pszStatusText *uint16) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_EnableModeless(inPlaceFrame *webViewIOleInPlaceFrame, fEnable BOOL) uintptr {
	return S_OK
}

func webView_IOleInPlaceFrame_TranslateAccelerator(inPlaceFrame *webViewIOleInPlaceFrame, lpmsg *MSG, wID uint32) uintptr {
	return E_NOTIMPL
}
