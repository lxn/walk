// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

var webViewIOleClientSiteVtbl *IOleClientSiteVtbl

func init() {
	webViewIOleClientSiteVtbl = &IOleClientSiteVtbl{
		syscall.NewCallback(webView_IOleClientSite_QueryInterface),
		syscall.NewCallback(webView_IOleClientSite_AddRef),
		syscall.NewCallback(webView_IOleClientSite_Release),
		syscall.NewCallback(webView_IOleClientSite_SaveObject),
		syscall.NewCallback(webView_IOleClientSite_GetMoniker),
		syscall.NewCallback(webView_IOleClientSite_GetContainer),
		syscall.NewCallback(webView_IOleClientSite_ShowObject),
		syscall.NewCallback(webView_IOleClientSite_OnShowWindow),
		syscall.NewCallback(webView_IOleClientSite_RequestNewObjectLayout),
	}
}

type webViewIOleClientSite struct {
	IOleClientSite
	inPlaceSite       webViewIOleInPlaceSite
	docHostUIHandler  webViewIDocHostUIHandler
	webBrowserEvents2 webViewDWebBrowserEvents2
}

func webView_IOleClientSite_QueryInterface(clientSite *webViewIOleClientSite, riid REFIID, ppvObject *unsafe.Pointer) uintptr {
	if EqualREFIID(riid, &IID_IUnknown) {
		*ppvObject = unsafe.Pointer(clientSite)
	} else if EqualREFIID(riid, &IID_IOleClientSite) {
		*ppvObject = unsafe.Pointer(clientSite)
	} else if EqualREFIID(riid, &IID_IOleInPlaceSite) {
		*ppvObject = unsafe.Pointer(&clientSite.inPlaceSite)
	} else if EqualREFIID(riid, &IID_IDocHostUIHandler) {
		*ppvObject = unsafe.Pointer(&clientSite.docHostUIHandler)
		// FIXME: Reactivate after fixing crash
		//	} else if EqualREFIID(riid, &IID_IDispatch) {
		//		*ppvObject = unsafe.Pointer(&clientSite.webBrowserEvents2)
	} else {
		*ppvObject = nil
		return E_NOINTERFACE
	}

	return S_OK
}

func webView_IOleClientSite_AddRef(clientSite *webViewIOleClientSite) uintptr {
	return 1
}

func webView_IOleClientSite_Release(clientSite *webViewIOleClientSite) uintptr {
	return 1
}

func webView_IOleClientSite_SaveObject(clientSite *webViewIOleClientSite) uintptr {
	return E_NOTIMPL
}

func webView_IOleClientSite_GetMoniker(clientSite *webViewIOleClientSite, dwAssign, dwWhichMoniker uint32, ppmk *unsafe.Pointer) uintptr {
	return E_NOTIMPL
}

func webView_IOleClientSite_GetContainer(clientSite *webViewIOleClientSite, ppContainer *unsafe.Pointer) uintptr {
	*ppContainer = nil

	return E_NOINTERFACE
}

func webView_IOleClientSite_ShowObject(clientSite *webViewIOleClientSite) uintptr {
	return S_OK
}

func webView_IOleClientSite_OnShowWindow(clientSite *webViewIOleClientSite, fShow BOOL) uintptr {
	return E_NOTIMPL
}

func webView_IOleClientSite_RequestNewObjectLayout(clientSite *webViewIOleClientSite) uintptr {
	return E_NOTIMPL
}
