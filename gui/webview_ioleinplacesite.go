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
	syscall.NewCallback(webView_IOleInPlaceSite_QueryInterface, 4+4+4),
	syscall.NewCallback(webView_IOleInPlaceSite_AddRef, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_Release, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_GetWindow, 4+4),
	syscall.NewCallback(webView_IOleInPlaceSite_ContextSensitiveHelp, 4+4),
	syscall.NewCallback(webView_IOleInPlaceSite_CanInPlaceActivate, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_OnInPlaceActivate, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_OnUIActivate, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_GetWindowContext, 4+4+4+4+4+4),
	syscall.NewCallback(webView_IOleInPlaceSite_Scroll, 4+8),
	syscall.NewCallback(webView_IOleInPlaceSite_OnUIDeactivate, 4+4),
	syscall.NewCallback(webView_IOleInPlaceSite_OnInPlaceDeactivate, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_DiscardUndoState, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_DeactivateAndUndo, 4),
	syscall.NewCallback(webView_IOleInPlaceSite_OnPosRectChange, 4+4),
}

var webViewIOleInPlaceSiteVtbl *IOleInPlaceSiteVtbl

func init() {
	webViewIOleInPlaceSiteVtbl = &IOleInPlaceSiteVtbl{
		uintptr(webViewIOleInPlaceSiteCbs.QueryInterface.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.AddRef.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.Release.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.GetWindow.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.ContextSensitiveHelp.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.CanInPlaceActivate.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.OnInPlaceActivate.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.OnUIActivate.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.GetWindowContext.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.Scroll.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.OnUIDeactivate.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.OnInPlaceDeactivate.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.DiscardUndoState.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.DeactivateAndUndo.ExtFnEntry()),
		uintptr(webViewIOleInPlaceSiteCbs.OnPosRectChange.ExtFnEntry()),
	}
}

type webViewIOleInPlaceSite struct {
	IOleInPlaceSite
	inPlaceFrame webViewIOleInPlaceFrame
}

func webView_IOleInPlaceSite_QueryInterface(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_QueryInterface")

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
	log.Println("webView_IOleInPlaceSite_AddRef")

	return 1
}

func webView_IOleInPlaceSite_Release(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_Release")

	return 1
}

func webView_IOleInPlaceSite_GetWindow(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_GetWindow")

	p := (*struct {
		inPlaceSite *webViewIOleInPlaceSite
		lphwnd      *HWND
	})(unsafe.Pointer(args))

	*p.lphwnd = p.inPlaceSite.inPlaceFrame.webView.hWnd

	return S_OK
}

func webView_IOleInPlaceSite_ContextSensitiveHelp(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_ContextSensitiveHelp")

	return E_NOTIMPL
}

func webView_IOleInPlaceSite_CanInPlaceActivate(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_CanInPlaceActivate")

	return S_OK
}

func webView_IOleInPlaceSite_OnInPlaceActivate(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_OnInPlaceActivate")

	return S_OK
}

func webView_IOleInPlaceSite_OnUIActivate(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_OnUIActivate")

	return S_OK
}

func webView_IOleInPlaceSite_GetWindowContext(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_GetWindowContext")

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
	log.Println("webView_IOleInPlaceSite_Scroll")

	return E_NOTIMPL
}

func webView_IOleInPlaceSite_OnUIDeactivate(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_OnUIDeactivate")

	return S_OK
}

func webView_IOleInPlaceSite_OnInPlaceDeactivate(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_OnInPlaceDeactivate")

	return S_OK
}

func webView_IOleInPlaceSite_DiscardUndoState(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_DiscardUndoState")

	return E_NOTIMPL
}

func webView_IOleInPlaceSite_DeactivateAndUndo(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_DeactivateAndUndo")

	return E_NOTIMPL
}

func webView_IOleInPlaceSite_OnPosRectChange(args *uintptr) uintptr {
	log.Println("webView_IOleInPlaceSite_OnPosRectChange")

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
