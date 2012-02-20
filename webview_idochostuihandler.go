// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "walk/winapi"

var webViewIDocHostUIHandlerVtbl *IDocHostUIHandlerVtbl

func init() {
	webViewIDocHostUIHandlerVtbl = &IDocHostUIHandlerVtbl{
		syscall.NewCallback(webView_IDocHostUIHandler_QueryInterface),
		syscall.NewCallback(webView_IDocHostUIHandler_AddRef),
		syscall.NewCallback(webView_IDocHostUIHandler_Release),
		syscall.NewCallback(webView_IDocHostUIHandler_ShowContextMenu),
		syscall.NewCallback(webView_IDocHostUIHandler_GetHostInfo),
		syscall.NewCallback(webView_IDocHostUIHandler_ShowUI),
		syscall.NewCallback(webView_IDocHostUIHandler_HideUI),
		syscall.NewCallback(webView_IDocHostUIHandler_UpdateUI),
		syscall.NewCallback(webView_IDocHostUIHandler_EnableModeless),
		syscall.NewCallback(webView_IDocHostUIHandler_OnDocWindowActivate),
		syscall.NewCallback(webView_IDocHostUIHandler_OnFrameWindowActivate),
		syscall.NewCallback(webView_IDocHostUIHandler_ResizeBorder),
		syscall.NewCallback(webView_IDocHostUIHandler_TranslateAccelerator),
		syscall.NewCallback(webView_IDocHostUIHandler_GetOptionKeyPath),
		syscall.NewCallback(webView_IDocHostUIHandler_GetDropTarget),
		syscall.NewCallback(webView_IDocHostUIHandler_GetExternal),
		syscall.NewCallback(webView_IDocHostUIHandler_TranslateUrl),
		syscall.NewCallback(webView_IDocHostUIHandler_FilterDataObject),
	}
}

type webViewIDocHostUIHandler struct {
	IDocHostUIHandler
}

func webView_IDocHostUIHandler_QueryInterface(docHostUIHandler *webViewIDocHostUIHandler, riid REFIID, ppvObject *unsafe.Pointer) HRESULT {
	// Just reuse the QueryInterface implementation we have for IOleClientSite.
	// We need to adjust object, which initially points at our
	// webViewIDocHostUIHandler, so it refers to the containing
	// webViewIOleClientSite for the call.
	var clientSite IOleClientSite
	var webViewInPlaceSite webViewIOleInPlaceSite

	ptr := uintptr(unsafe.Pointer(docHostUIHandler)) - uintptr(unsafe.Sizeof(clientSite)) -
		uintptr(unsafe.Sizeof(webViewInPlaceSite))

	return webView_IOleClientSite_QueryInterface((*webViewIOleClientSite)(unsafe.Pointer(ptr)), riid, ppvObject)
}

func webView_IDocHostUIHandler_AddRef(docHostUIHandler *webViewIDocHostUIHandler) HRESULT {
	return 1
}

func webView_IDocHostUIHandler_Release(docHostUIHandler *webViewIDocHostUIHandler) HRESULT {
	return 1
}

func webView_IDocHostUIHandler_ShowContextMenu(docHostUIHandler *webViewIDocHostUIHandler, dwID uint, ppt *POINT, pcmdtReserved *IUnknown, pdispReserved uintptr) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_GetHostInfo(docHostUIHandler *webViewIDocHostUIHandler, pInfo *DOCHOSTUIINFO) HRESULT {
	pInfo.CbSize = uint32(unsafe.Sizeof(*pInfo))
	pInfo.DwFlags = DOCHOSTUIFLAG_NO3DBORDER
	pInfo.DwDoubleClick = DOCHOSTUIDBLCLK_DEFAULT

	return S_OK
}

func webView_IDocHostUIHandler_ShowUI(docHostUIHandler *webViewIDocHostUIHandler, dwID uint, pActiveObject uintptr, pCommandTarget uintptr, pFrame *IOleInPlaceFrame, pDoc uintptr) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_HideUI(docHostUIHandler *webViewIDocHostUIHandler) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_UpdateUI(docHostUIHandler *webViewIDocHostUIHandler) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_EnableModeless(docHostUIHandler *webViewIDocHostUIHandler, fEnable BOOL) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_OnDocWindowActivate(docHostUIHandler *webViewIDocHostUIHandler, fActivate BOOL) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_OnFrameWindowActivate(docHostUIHandler *webViewIDocHostUIHandler, fActivate BOOL) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_ResizeBorder(docHostUIHandler *webViewIDocHostUIHandler, prcBorder *RECT, pUIWindow uintptr, fRameWindow BOOL) HRESULT {
	return S_OK
}

func webView_IDocHostUIHandler_TranslateAccelerator(docHostUIHandler *webViewIDocHostUIHandler, lpMsg *MSG, pguidCmdGroup *GUID, nCmdID uint) HRESULT {
	return S_FALSE
}

func webView_IDocHostUIHandler_GetOptionKeyPath(docHostUIHandler *webViewIDocHostUIHandler, pchKey *uint16, dw uint) HRESULT {
	return S_FALSE
}

func webView_IDocHostUIHandler_GetDropTarget(docHostUIHandler *webViewIDocHostUIHandler, pDropTarget uintptr, ppDropTarget *uintptr) HRESULT {
	return S_FALSE
}

func webView_IDocHostUIHandler_GetExternal(docHostUIHandler *webViewIDocHostUIHandler, ppDispatch *uintptr) HRESULT {
	*ppDispatch = 0

	return S_FALSE
}

func webView_IDocHostUIHandler_TranslateUrl(docHostUIHandler *webViewIDocHostUIHandler, dwTranslate uint, pchURLIn *uint16, ppchURLOut **uint16) HRESULT {
	*ppchURLOut = nil

	return S_FALSE
}

func webView_IDocHostUIHandler_FilterDataObject(docHostUIHandler *webViewIDocHostUIHandler, pDO uintptr, ppDORet *uintptr) HRESULT {
	*ppDORet = 0

	return S_FALSE
}
