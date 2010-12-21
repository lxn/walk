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
	. "walk/winapi/shdocvw"
)

type webViewIDocHostUIHandlerCallbacks struct {
	QueryInterface        *syscall.Callback
	AddRef                *syscall.Callback
	Release               *syscall.Callback
	ShowContextMenu       *syscall.Callback
	GetHostInfo           *syscall.Callback
	ShowUI                *syscall.Callback
	HideUI                *syscall.Callback
	UpdateUI              *syscall.Callback
	EnableModeless        *syscall.Callback
	OnDocWindowActivate   *syscall.Callback
	OnFrameWindowActivate *syscall.Callback
	ResizeBorder          *syscall.Callback
	TranslateAccelerator  *syscall.Callback
	GetOptionKeyPath      *syscall.Callback
	GetDropTarget         *syscall.Callback
	GetExternal           *syscall.Callback
	TranslateUrl          *syscall.Callback
	FilterDataObject      *syscall.Callback
}

var webViewIDocHostUIHandlerCbs = &webViewIDocHostUIHandlerCallbacks{
	syscall.NewCallback(webView_IDocHostUIHandler_QueryInterface, 4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_AddRef, 4),
	syscall.NewCallback(webView_IDocHostUIHandler_Release, 4),
	syscall.NewCallback(webView_IDocHostUIHandler_ShowContextMenu, 4+4+4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_GetHostInfo, 4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_ShowUI, 4+4+4+4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_HideUI, 4),
	syscall.NewCallback(webView_IDocHostUIHandler_UpdateUI, 4),
	syscall.NewCallback(webView_IDocHostUIHandler_EnableModeless, 4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_OnDocWindowActivate, 4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_OnFrameWindowActivate, 4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_ResizeBorder, 4+4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_TranslateAccelerator, 4+4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_GetOptionKeyPath, 4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_GetDropTarget, 4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_GetExternal, 4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_TranslateUrl, 4+4+4+4),
	syscall.NewCallback(webView_IDocHostUIHandler_FilterDataObject, 4+4+4),
}

var webViewIDocHostUIHandlerVtbl *IDocHostUIHandlerVtbl

func init() {
	webViewIDocHostUIHandlerVtbl = &IDocHostUIHandlerVtbl{
		webViewIDocHostUIHandlerCbs.QueryInterface.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.AddRef.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.Release.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.ShowContextMenu.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.GetHostInfo.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.ShowUI.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.HideUI.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.UpdateUI.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.EnableModeless.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.OnDocWindowActivate.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.OnFrameWindowActivate.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.ResizeBorder.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.TranslateAccelerator.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.GetOptionKeyPath.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.GetDropTarget.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.GetExternal.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.TranslateUrl.ExtFnEntry(),
		webViewIDocHostUIHandlerCbs.FilterDataObject.ExtFnEntry(),
	}
}

type webViewIDocHostUIHandler struct {
	IDocHostUIHandler
}

func webView_IDocHostUIHandler_QueryInterface(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_QueryInterface")

	p := (*struct {
		object    uintptr
		riid      REFIID
		ppvObject *unsafe.Pointer
	})(unsafe.Pointer(args))

	// Just reuse the QueryInterface implementation we have for IOleClientSite.
	// We need to adjust object, which initially points at our
	// webViewIDocHostUIHandler, so it refers to the containing
	// webViewIOleClientSite for the call.
	var clientSite IOleClientSite
	var webViewInPlaceSite webViewIOleInPlaceSite

	ptr := int(p.object) - unsafe.Sizeof(clientSite) - unsafe.Sizeof(webViewInPlaceSite)
	p.object = uintptr(ptr)

	return webView_IOleClientSite_QueryInterface(args)
}

func webView_IDocHostUIHandler_AddRef(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_AddRef")

	return 1
}

func webView_IDocHostUIHandler_Release(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_Release")

	return 1
}

func webView_IDocHostUIHandler_ShowContextMenu(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_ShowContextMenu")

	return S_OK
}

func webView_IDocHostUIHandler_GetHostInfo(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_GetHostInfo")

	p := (*struct {
		docHostUIHandler *webViewIDocHostUIHandler
		pInfo            *DOCHOSTUIINFO
	})(unsafe.Pointer(args))

	p.pInfo.CbSize = uint(unsafe.Sizeof(*p.pInfo))
	p.pInfo.DwFlags = DOCHOSTUIFLAG_NO3DBORDER
	p.pInfo.DwDoubleClick = DOCHOSTUIDBLCLK_DEFAULT

	return S_OK
}

func webView_IDocHostUIHandler_ShowUI(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_ShowUI")

	return S_OK
}

func webView_IDocHostUIHandler_HideUI(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_HideUI")

	return S_OK
}

func webView_IDocHostUIHandler_UpdateUI(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_UpdateUI")

	return S_OK
}

func webView_IDocHostUIHandler_EnableModeless(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_EnableModeless")

	return S_OK
}

func webView_IDocHostUIHandler_OnDocWindowActivate(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_OnDocWindowActivate")

	return S_OK
}

func webView_IDocHostUIHandler_OnFrameWindowActivate(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_OnFrameWindowActivate")

	return S_OK
}

func webView_IDocHostUIHandler_ResizeBorder(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_ResizeBorder")

	return S_OK
}

func webView_IDocHostUIHandler_TranslateAccelerator(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_TranslateAccelerator")

	return S_FALSE
}

func webView_IDocHostUIHandler_GetOptionKeyPath(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_GetOptionKeyPath")

	return S_FALSE
}

func webView_IDocHostUIHandler_GetDropTarget(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_GetDropTarget")

	return S_FALSE
}

func webView_IDocHostUIHandler_GetExternal(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_GetExternal")

	p := (*struct {
		docHostUIHandler *webViewIDocHostUIHandler
		ppDispatch       *unsafe.Pointer
	})(unsafe.Pointer(args))

	*p.ppDispatch = nil

	return S_FALSE
}

func webView_IDocHostUIHandler_TranslateUrl(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_TranslateUrl")

	p := (*struct {
		docHostUIHandler *webViewIDocHostUIHandler
		dwTranslate      uint
		pchURLIn         *uint16
		ppchURLOut       **uint16
	})(unsafe.Pointer(args))

	*p.ppchURLOut = nil

	return S_FALSE
}

func webView_IDocHostUIHandler_FilterDataObject(args *uintptr) uintptr {
	log.Println("webView_IDocHostUIHandler_FilterDataObject")

	p := (*struct {
		docHostUIHandler *webViewIDocHostUIHandler
		pDO              unsafe.Pointer
		ppDORet          *unsafe.Pointer
	})(unsafe.Pointer(args))

	*p.ppDORet = nil

	return S_FALSE
}
