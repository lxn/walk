// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

var webViewDWebBrowserEvents2Vtbl *DWebBrowserEvents2Vtbl

func init() {
	webViewDWebBrowserEvents2Vtbl = &DWebBrowserEvents2Vtbl{
		syscall.NewCallback(webView_DWebBrowserEvents2_QueryInterface),
		syscall.NewCallback(webView_DWebBrowserEvents2_AddRef),
		syscall.NewCallback(webView_DWebBrowserEvents2_Release),
		syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfoCount),
		syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfo),
		syscall.NewCallback(webView_DWebBrowserEvents2_GetIDsOfNames),
		syscall.NewCallback(webView_DWebBrowserEvents2_Invoke),
		syscall.NewCallback(webView_DWebBrowserEvents2_StatusTextChange),
		syscall.NewCallback(webView_DWebBrowserEvents2_ProgressChange),
		syscall.NewCallback(webView_DWebBrowserEvents2_CommandStateChange),
		syscall.NewCallback(webView_DWebBrowserEvents2_DownloadBegin),
		syscall.NewCallback(webView_DWebBrowserEvents2_DownloadComplete),
		syscall.NewCallback(webView_DWebBrowserEvents2_TitleChange),
		syscall.NewCallback(webView_DWebBrowserEvents2_PropertyChange),
		syscall.NewCallback(webView_DWebBrowserEvents2_BeforeNavigate2),
		syscall.NewCallback(webView_DWebBrowserEvents2_NewWindow2),
		syscall.NewCallback(webView_DWebBrowserEvents2_NavigateComplete2),
		syscall.NewCallback(webView_DWebBrowserEvents2_DocumentComplete),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnQuit),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnVisible),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnToolBar),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnMenuBar),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnStatusBar),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnFullScreen),
		syscall.NewCallback(webView_DWebBrowserEvents2_OnTheaterMode),
		syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetResizable),
		syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetLeft),
		syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetTop),
		syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetWidth),
		syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetHeight),
		syscall.NewCallback(webView_DWebBrowserEvents2_WindowClosing),
		syscall.NewCallback(webView_DWebBrowserEvents2_ClientToHostWindow),
		syscall.NewCallback(webView_DWebBrowserEvents2_SetSecureLockIcon),
		syscall.NewCallback(webView_DWebBrowserEvents2_FileDownload),
		syscall.NewCallback(webView_DWebBrowserEvents2_NavigateError),
		syscall.NewCallback(webView_DWebBrowserEvents2_PrintTemplateInstantiation),
		syscall.NewCallback(webView_DWebBrowserEvents2_PrintTemplateTeardown),
		syscall.NewCallback(webView_DWebBrowserEvents2_UpdatePageStatus),
		syscall.NewCallback(webView_DWebBrowserEvents2_PrivacyImpactedStateChange),
		syscall.NewCallback(webView_DWebBrowserEvents2_NewWindow3),
	}
}

type webViewDWebBrowserEvents2 struct {
	DWebBrowserEvents2
}

func webView_DWebBrowserEvents2_QueryInterface(wbe2 *webViewDWebBrowserEvents2, riid REFIID, ppvObject *unsafe.Pointer) HRESULT {
	// Just reuse the QueryInterface implementation we have for IOleClientSite.
	// We need to adjust object, which initially points at our
	// webViewDWebBrowserEvents2, so it refers to the containing
	// webViewIOleClientSite for the call.
	var clientSite IOleClientSite
	var webViewInPlaceSite webViewIOleInPlaceSite
	var docHostUIHandler webViewIDocHostUIHandler

	ptr := uintptr(unsafe.Pointer(wbe2)) - uintptr(unsafe.Sizeof(clientSite)) -
		uintptr(unsafe.Sizeof(webViewInPlaceSite)) - uintptr(unsafe.Sizeof(docHostUIHandler))

	return webView_IOleClientSite_QueryInterface((*webViewIOleClientSite)(unsafe.Pointer(ptr)), riid, ppvObject)
}

/*func webView_DWebBrowserEvents2_AddRef(wbe2 *webViewDWebBrowserEvents2) uint {
	return 1
}

func webView_DWebBrowserEvents2_Release(wbe2 *webViewDWebBrowserEvents2) uint {
	return 1
}

func webView_DWebBrowserEvents2_GetTypeInfoCount(wbe2 *webViewDWebBrowserEvents2, pctinfo *uint) HRESULT {
	return E_NOTIMPL
}*/

/*func webView_DWebBrowserEvents2_QueryInterface(args *uintptr) uintptr {
	p := (*struct {
		object    uintptr
		riid      REFIID
		ppvObject *unsafe.Pointer
	})(unsafe.Pointer(args))

	// Just reuse the QueryInterface implementation we have for IOleClientSite.
	// We need to adjust object, which initially points at our
	// webViewDWebBrowserEvents2, so it refers to the containing
	// webViewIOleClientSite for the call.
	var clientSite IOleClientSite
	var webViewInPlaceSite webViewIOleInPlaceSite
	var docHostUIHandler webViewIDocHostUIHandler

	ptr := int(p.object) - unsafe.Sizeof(clientSite) - unsafe.Sizeof(webViewInPlaceSite) - unsafe.Sizeof(docHostUIHandler)
	p.object = uintptr(ptr)

	return webView_IOleClientSite_QueryInterface(args)
}*/

func webView_DWebBrowserEvents2_AddRef(args *uintptr) uintptr {
	return 1
}

func webView_DWebBrowserEvents2_Release(args *uintptr) uintptr {
	return 1
}

func webView_DWebBrowserEvents2_GetTypeInfoCount(args *uintptr) uintptr {
	/*	p := (*struct {
			wbe2    *webViewDWebBrowserEvents2
			pctinfo *uint
		})(unsafe.Pointer(args))

		*p.pctinfo = 0

		return S_OK*/

	return E_NOTIMPL
}

func webView_DWebBrowserEvents2_GetTypeInfo(args *uintptr) uintptr {
	/*	p := (*struct {
				wbe2         *webViewDWebBrowserEvents2
			})(unsafe.Pointer(args))

		    unsigned int  iTInfo,         
		    LCID  lcid,                   
		    ITypeInfo FAR* FAR*  ppTInfo*/

	return E_NOTIMPL
}

func webView_DWebBrowserEvents2_GetIDsOfNames(args *uintptr) uintptr {
	/*	p := (*struct {
				wbe2         *webViewDWebBrowserEvents2
			})(unsafe.Pointer(args))

		    REFIID             riid,                  
		    OLECHAR FAR* FAR*  rgszNames,  
		    unsigned int       cNames,          
		    LCID               lcid,                   
		    DISPID       FAR*  rgDispId*/

	return E_NOTIMPL
}

func webView_DWebBrowserEvents2_Invoke(args *uintptr) uintptr {
	/*p := (*struct {
		wbe2         *webViewDWebBrowserEvents2
		dispIdMember DISPID
		riid         REFIID
		lcid         uint // LCID
		wFlags       uint16
		pDispParams  *DISPPARAMS
		pVarResult   *VARIANT
		pExcepInfo   unsafe.Pointer // *EXCEPINFO
		puArgErr     *uint
	})(unsafe.Pointer(args))*/

	return DISP_E_MEMBERNOTFOUND
}

func webView_DWebBrowserEvents2_BeforeNavigate2(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_ClientToHostWindow(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_CommandStateChange(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_DocumentComplete(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_DownloadBegin(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_DownloadComplete(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_FileDownload(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_NavigateComplete2(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_NavigateError(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_NewWindow2(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_NewWindow3(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnFullScreen(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnMenuBar(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnQuit(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnStatusBar(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnTheaterMode(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnToolBar(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_OnVisible(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_PrintTemplateInstantiation(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_PrintTemplateTeardown(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_PrivacyImpactedStateChange(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_ProgressChange(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_PropertyChange(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_SetSecureLockIcon(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_StatusTextChange(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_TitleChange(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_UpdatePageStatus(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_WindowClosing(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_WindowSetHeight(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_WindowSetLeft(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_WindowSetResizable(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_WindowSetTop(args *uintptr) uintptr {
	return 0
}

func webView_DWebBrowserEvents2_WindowSetWidth(args *uintptr) uintptr {
	return 0
}
