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
	. "walk/winapi/oleaut32"
	. "walk/winapi/shdocvw"
)

type webViewIOleClientSiteCallbacks struct {
	QueryInterface         *syscall.Callback
	AddRef                 *syscall.Callback
	Release                *syscall.Callback
	SaveObject             *syscall.Callback
	GetMoniker             *syscall.Callback
	GetContainer           *syscall.Callback
	ShowObject             *syscall.Callback
	OnShowWindow           *syscall.Callback
	RequestNewObjectLayout *syscall.Callback
}

var webViewIOleClientSiteCbs = &webViewIOleClientSiteCallbacks{
	syscall.NewCallback(webView_IOleClientSite_QueryInterface, 1+2),
	syscall.NewCallback(webView_IOleClientSite_AddRef, 1+0),
	syscall.NewCallback(webView_IOleClientSite_Release, 1+0),
	syscall.NewCallback(webView_IOleClientSite_SaveObject, 1+0),
	syscall.NewCallback(webView_IOleClientSite_GetMoniker, 1+3),
	syscall.NewCallback(webView_IOleClientSite_GetContainer, 1+1),
	syscall.NewCallback(webView_IOleClientSite_ShowObject, 1+0),
	syscall.NewCallback(webView_IOleClientSite_OnShowWindow, 1+1),
	syscall.NewCallback(webView_IOleClientSite_RequestNewObjectLayout, 1+0),
}

var webViewIOleClientSiteVtbl *IOleClientSiteVtbl

func init() {
	webViewIOleClientSiteVtbl = &IOleClientSiteVtbl{
		webViewIOleClientSiteCbs.QueryInterface.ExtFnEntry(),
		webViewIOleClientSiteCbs.AddRef.ExtFnEntry(),
		webViewIOleClientSiteCbs.Release.ExtFnEntry(),
		webViewIOleClientSiteCbs.SaveObject.ExtFnEntry(),
		webViewIOleClientSiteCbs.GetMoniker.ExtFnEntry(),
		webViewIOleClientSiteCbs.GetContainer.ExtFnEntry(),
		webViewIOleClientSiteCbs.ShowObject.ExtFnEntry(),
		webViewIOleClientSiteCbs.OnShowWindow.ExtFnEntry(),
		webViewIOleClientSiteCbs.RequestNewObjectLayout.ExtFnEntry(),
	}
}

type webViewIOleClientSite struct {
	IOleClientSite
	inPlaceSite       webViewIOleInPlaceSite
	docHostUIHandler  webViewIDocHostUIHandler
	webBrowserEvents2 webViewDWebBrowserEvents2
}

func webView_IOleClientSite_QueryInterface(args *uintptr) uintptr {
	p := (*struct {
		clientSite *webViewIOleClientSite
		riid       REFIID
		ppvObject  *unsafe.Pointer
	})(unsafe.Pointer(args))

	if EqualREFIID(p.riid, &IID_IUnknown) {
		*p.ppvObject = unsafe.Pointer(p.clientSite)
	} else if EqualREFIID(p.riid, &IID_IOleClientSite) {
		*p.ppvObject = unsafe.Pointer(p.clientSite)
	} else if EqualREFIID(p.riid, &IID_IOleInPlaceSite) {
		*p.ppvObject = unsafe.Pointer(&p.clientSite.inPlaceSite)
	} else if EqualREFIID(p.riid, &IID_IDocHostUIHandler) {
		*p.ppvObject = unsafe.Pointer(&p.clientSite.docHostUIHandler)
	} else if EqualREFIID(p.riid, &IID_IDispatch) {
		*p.ppvObject = unsafe.Pointer(&p.clientSite.webBrowserEvents2)
	} else {
		*p.ppvObject = nil
		return E_NOINTERFACE
	}

	return S_OK
}

func webView_IOleClientSite_AddRef(args *uintptr) uintptr {
	return 1
}

func webView_IOleClientSite_Release(args *uintptr) uintptr {
	return 1
}

func webView_IOleClientSite_SaveObject(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleClientSite_GetMoniker(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleClientSite_GetContainer(args *uintptr) uintptr {
	p := (*struct {
		clientSite  *webViewIOleClientSite
		ppContainer *unsafe.Pointer
	})(unsafe.Pointer(args))

	*p.ppContainer = nil

	return E_NOINTERFACE
}

func webView_IOleClientSite_ShowObject(args *uintptr) uintptr {
	return S_OK
}

func webView_IOleClientSite_OnShowWindow(args *uintptr) uintptr {
	return E_NOTIMPL
}

func webView_IOleClientSite_RequestNewObjectLayout(args *uintptr) uintptr {
	return E_NOTIMPL
}
