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
	. "walk/winapi/gdi32"
	. "walk/winapi/ole32"
	. "walk/winapi/user32"
)

type webViewIOleInPlaceSiteCallbacks struct {
	QueryInterface       *syscall.Callback
	AddRef               *syscall.Callback
	Release              *syscall.Callback
	GetWindow            *syscall.Callback
	ContextSensitiveHelp *syscall.Callback
	CanInPlaceActivate   *syscall.Callback
	OnInPlaceActivate    *syscall.Callback
	OnUIActivate         *syscall.Callback
	GetWindowContext     *syscall.Callback
	Scroll               *syscall.Callback
	OnUIDeactivate       *syscall.Callback
	OnInPlaceDeactivate  *syscall.Callback
	DiscardUndoState     *syscall.Callback
	DeactivateAndUndo    *syscall.Callback
	OnPosRectChange      *syscall.Callback
}

var webViewIOleInPlaceSiteCbs = &webViewIOleInPlaceSiteCallbacks{
	syscall.NewCallback(webView_IOleInPlaceSite_QueryInterface, 1+2),
	syscall.NewCallback(webView_IOleInPlaceSite_AddRef, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_Release, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_GetWindow, 1+1),
	syscall.NewCallback(webView_IOleInPlaceSite_ContextSensitiveHelp, 1+1),
	syscall.NewCallback(webView_IOleInPlaceSite_CanInPlaceActivate, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_OnInPlaceActivate, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_OnUIActivate, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_GetWindowContext, 1+5),
	syscall.NewCallback(webView_IOleInPlaceSite_Scroll, 1+1*2),
	syscall.NewCallback(webView_IOleInPlaceSite_OnUIDeactivate, 1+1),
	syscall.NewCallback(webView_IOleInPlaceSite_OnInPlaceDeactivate, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_DiscardUndoState, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_DeactivateAndUndo, 1+0),
	syscall.NewCallback(webView_IOleInPlaceSite_OnPosRectChange, 1+1),
}

var webViewIOleInPlaceSiteVtbl *IOleInPlaceSiteVtbl

func init() {
	webViewIOleInPlaceSiteVtbl = &IOleInPlaceSiteVtbl{
		webViewIOleInPlaceSiteCbs.QueryInterface.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.AddRef.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.Release.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.GetWindow.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.ContextSensitiveHelp.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.CanInPlaceActivate.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.OnInPlaceActivate.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.OnUIActivate.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.GetWindowContext.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.Scroll.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.OnUIDeactivate.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.OnInPlaceDeactivate.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.DiscardUndoState.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.DeactivateAndUndo.ExtFnEntry(),
		webViewIOleInPlaceSiteCbs.OnPosRectChange.ExtFnEntry(),
	}
}

type webViewIOleInPlaceSite struct {
	IOleInPlaceSite
	inPlaceFrame webViewIOleInPlaceFrame
}

func webView_IOleInPlaceSite_QueryInterface(args *uintptr) uintptr {
	p := (*struct {
		object    uintptr
		riid      REFIID
		ppvObject *unsafe.Pointer
	})(unsafe.Pointer(args))

	// Just reuse the QueryInterface implementation we have for IOleClientSite.
	// We need to adjust object from the webViewIDocHostUIHandler to the
	// containing webViewIOleInPlaceSite.
	var clientSite IOleClientSite

	p.object -= uintptr(unsafe.Sizeof(clientSite))

	return webView_IOleClientSite_QueryInterface(args)
}

func webView_IOleInPlaceSite_AddRef(args *uintptr) uintptr {
	return 1
}

func webView_IOleInPlaceSite_Release(args *uintptr) uintptr {
	return 1
}

func webView_IOleInPlaceSite_GetWindow(args *uintptr) uintptr {
	p := (*struct {
		inPlaceSite *webViewIOleInPlaceSite
		lphwnd      *HWND
	})(unsafe.Pointer(args))

	*p.lphwnd = p.inPlaceSite.inPlaceFrame.webView.hWnd

	return S_OK
}

func webView_IOleInPlaceSite_ContextSensitiveHelp(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceSite_CanInPlaceActivate(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceSite_OnInPlaceActivate(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceSite_OnUIActivate(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceSite_GetWindowContext(args *uintptr) uintptr {
	p := (*struct {
		inPlaceSite  *webViewIOleInPlaceSite
		lplpFrame    **webViewIOleInPlaceFrame
		lplpDoc      *unsafe.Pointer
		lprcPosRect  *RECT
		lprcClipRect *RECT
		lpFrameInfo  *OLEINPLACEFRAMEINFO
	})(unsafe.Pointer(args))

	*p.lplpFrame = &p.inPlaceSite.inPlaceFrame
	*p.lplpDoc = nil

	p.lpFrameInfo.FMDIApp = FALSE
	p.lpFrameInfo.HwndFrame = p.inPlaceSite.inPlaceFrame.webView.hWnd
	p.lpFrameInfo.Haccel = 0
	p.lpFrameInfo.CAccelEntries = 0

	return S_OK
}

func webView_IOleInPlaceSite_Scroll(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceSite_OnUIDeactivate(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceSite_OnInPlaceDeactivate(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleInPlaceSite_DiscardUndoState(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceSite_DeactivateAndUndo(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleInPlaceSite_OnPosRectChange(args *uintptr) uintptr {
	p := (*struct {
		inPlaceSite *webViewIOleInPlaceSite
		lprcPosRect *RECT
	})(unsafe.Pointer(args))

	browserObject := p.inPlaceSite.inPlaceFrame.webView.browserObject
	var inPlaceObjectPtr unsafe.Pointer
	if hr := browserObject.QueryInterface(&IID_IOleInPlaceObject, &inPlaceObjectPtr); FAILED(hr) {
		return uintptr(hr)
	}
	inPlaceObject := (*IOleInPlaceObject)(inPlaceObjectPtr)
	defer inPlaceObject.Release()

	return uintptr(inPlaceObject.SetObjectRects(p.lprcPosRect, p.lprcPosRect))
}
