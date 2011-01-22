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
	syscall.NewCallback(webView_IOleClientSite_QueryInterface, 4+4+4),
	syscall.NewCallback(webView_IOleClientSite_AddRef, 4),
	syscall.NewCallback(webView_IOleClientSite_Release, 4),
	syscall.NewCallback(webView_IOleClientSite_SaveObject, 4),
	syscall.NewCallback(webView_IOleClientSite_GetMoniker, 4+4+4+4),
	syscall.NewCallback(webView_IOleClientSite_GetContainer, 4+4),
	syscall.NewCallback(webView_IOleClientSite_ShowObject, 4),
	syscall.NewCallback(webView_IOleClientSite_OnShowWindow, 4+4),
	syscall.NewCallback(webView_IOleClientSite_RequestNewObjectLayout, 4),
}

var webViewIOleClientSiteVtbl *IOleClientSiteVtbl

func init() {
	webViewIOleClientSiteVtbl = &IOleClientSiteVtbl{
		uintptr(webViewIOleClientSiteCbs.QueryInterface.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.AddRef.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.Release.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.SaveObject.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.GetMoniker.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.GetContainer.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.ShowObject.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.OnShowWindow.ExtFnEntry()),
		uintptr(webViewIOleClientSiteCbs.RequestNewObjectLayout.ExtFnEntry()),
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
		log.Println("webView_IOleClientSite_QueryInterface (IID_IUnknown)")
		*p.ppvObject = unsafe.Pointer(p.clientSite)
	} else if EqualREFIID(p.riid, &IID_IOleClientSite) {
		log.Println("webView_IOleClientSite_QueryInterface (IID_IOleClientSite)")
		*p.ppvObject = unsafe.Pointer(p.clientSite)
	} else if EqualREFIID(p.riid, &IID_IOleInPlaceSite) {
		log.Println("webView_IOleClientSite_QueryInterface (IID_IOleInPlaceSite)")
		*p.ppvObject = unsafe.Pointer(&p.clientSite.inPlaceSite)
	} else if EqualREFIID(p.riid, &IID_IDocHostUIHandler) {
		log.Println("webView_IOleClientSite_QueryInterface (IID_IDocHostUIHandler)")
		*p.ppvObject = unsafe.Pointer(&p.clientSite.docHostUIHandler)
	} else if EqualREFIID(p.riid, &IID_IDispatch) {
		log.Println("webView_IOleClientSite_QueryInterface (IID_IDispatch)")
		*p.ppvObject = unsafe.Pointer(&p.clientSite.webBrowserEvents2)
	} else {
		log.Println("webView_IOleClientSite_QueryInterface (?)")
		*p.ppvObject = nil
		return E_NOINTERFACE
	}

	return S_OK
}

func webView_IOleClientSite_AddRef(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_AddRef")

	return 1
}

func webView_IOleClientSite_Release(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_Release")

	return 1
}

func webView_IOleClientSite_SaveObject(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_SaveObject")

	return E_NOTIMPL
}

func webView_IOleClientSite_GetMoniker(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_GetMoniker")

	return E_NOTIMPL
}

func webView_IOleClientSite_GetContainer(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_GetContainer")

	p := (*struct {
		clientSite  *webViewIOleClientSite
		ppContainer *unsafe.Pointer
	})(unsafe.Pointer(args))

	*p.ppContainer = nil

	return E_NOINTERFACE
}

func webView_IOleClientSite_ShowObject(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_ShowObject")

	return S_OK
}

func webView_IOleClientSite_OnShowWindow(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_OnShowWindow")

	return E_NOTIMPL
}

func webView_IOleClientSite_RequestNewObjectLayout(args *uintptr) uintptr {
	log.Println("webView_IOleClientSite_RequestNewObjectLayout")

	return E_NOTIMPL
}
