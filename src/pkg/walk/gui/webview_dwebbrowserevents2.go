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

type webViewDWebBrowserEvents2Callbacks struct {
	QueryInterface             *syscall.Callback
	AddRef                     *syscall.Callback
	Release                    *syscall.Callback
	GetTypeInfoCount           *syscall.Callback
	GetTypeInfo                *syscall.Callback
	GetIDsOfNames              *syscall.Callback
	Invoke                     *syscall.Callback
	StatusTextChange           *syscall.Callback
	ProgressChange             *syscall.Callback
	CommandStateChange         *syscall.Callback
	DownloadBegin              *syscall.Callback
	DownloadComplete           *syscall.Callback
	TitleChange                *syscall.Callback
	PropertyChange             *syscall.Callback
	BeforeNavigate2            *syscall.Callback
	NewWindow2                 *syscall.Callback
	NavigateComplete2          *syscall.Callback
	DocumentComplete           *syscall.Callback
	OnQuit                     *syscall.Callback
	OnVisible                  *syscall.Callback
	OnToolBar                  *syscall.Callback
	OnMenuBar                  *syscall.Callback
	OnStatusBar                *syscall.Callback
	OnFullScreen               *syscall.Callback
	OnTheaterMode              *syscall.Callback
	WindowSetResizable         *syscall.Callback
	WindowSetLeft              *syscall.Callback
	WindowSetTop               *syscall.Callback
	WindowSetWidth             *syscall.Callback
	WindowSetHeight            *syscall.Callback
	WindowClosing              *syscall.Callback
	ClientToHostWindow         *syscall.Callback
	SetSecureLockIcon          *syscall.Callback
	FileDownload               *syscall.Callback
	NavigateError              *syscall.Callback
	PrintTemplateInstantiation *syscall.Callback
	PrintTemplateTeardown      *syscall.Callback
	UpdatePageStatus           *syscall.Callback
	PrivacyImpactedStateChange *syscall.Callback
	NewWindow3                 *syscall.Callback
}

var webViewDWebBrowserEvents2Cbs = &webViewDWebBrowserEvents2Callbacks{
	syscall.NewCallback(webView_DWebBrowserEvents2_QueryInterface, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_AddRef, 1+0),
	syscall.NewCallback(webView_DWebBrowserEvents2_Release, 1+0),
	syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfoCount, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfo, 1+3),
	syscall.NewCallback(webView_DWebBrowserEvents2_GetIDsOfNames, 1+5),
	syscall.NewCallback(webView_DWebBrowserEvents2_Invoke, 1+8),

	syscall.NewCallback(webView_DWebBrowserEvents2_StatusTextChange, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_ProgressChange, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_CommandStateChange, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_DownloadBegin, 1+0),
	syscall.NewCallback(webView_DWebBrowserEvents2_DownloadComplete, 1+0),
	syscall.NewCallback(webView_DWebBrowserEvents2_TitleChange, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_PropertyChange, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_BeforeNavigate2, 1+7),
	syscall.NewCallback(webView_DWebBrowserEvents2_NewWindow2, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_NavigateComplete2, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_DocumentComplete, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnQuit, 1+0),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnVisible, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnToolBar, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnMenuBar, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnStatusBar, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnFullScreen, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnTheaterMode, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetResizable, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetLeft, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetTop, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetWidth, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetHeight, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowClosing, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_ClientToHostWindow, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_SetSecureLockIcon, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_FileDownload, 1+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_NavigateError, 1+5),
	syscall.NewCallback(webView_DWebBrowserEvents2_PrintTemplateInstantiation, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_PrintTemplateTeardown, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_UpdatePageStatus, 1+0),
	syscall.NewCallback(webView_DWebBrowserEvents2_PrivacyImpactedStateChange, 1+1),
	syscall.NewCallback(webView_DWebBrowserEvents2_NewWindow3, 1+5),
}

var webViewDWebBrowserEvents2Vtbl *DWebBrowserEvents2Vtbl

func init() {
	webViewDWebBrowserEvents2Vtbl = &DWebBrowserEvents2Vtbl{
		webViewDWebBrowserEvents2Cbs.QueryInterface.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.AddRef.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.Release.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.GetTypeInfoCount.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.GetTypeInfo.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.GetIDsOfNames.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.Invoke.ExtFnEntry(),

		webViewDWebBrowserEvents2Cbs.StatusTextChange.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.ProgressChange.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.CommandStateChange.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.DownloadBegin.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.DownloadComplete.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.TitleChange.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.PropertyChange.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.BeforeNavigate2.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.NewWindow2.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.NavigateComplete2.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.DocumentComplete.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnQuit.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnVisible.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnToolBar.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnMenuBar.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnStatusBar.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnFullScreen.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.OnTheaterMode.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.WindowSetResizable.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.WindowSetLeft.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.WindowSetTop.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.WindowSetWidth.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.WindowSetHeight.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.WindowClosing.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.ClientToHostWindow.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.SetSecureLockIcon.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.FileDownload.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.NavigateError.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.PrintTemplateInstantiation.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.PrintTemplateTeardown.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.UpdatePageStatus.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.PrivacyImpactedStateChange.ExtFnEntry(),
		webViewDWebBrowserEvents2Cbs.NewWindow3.ExtFnEntry(),
	}
}

type webViewDWebBrowserEvents2 struct {
	DWebBrowserEvents2
}

func webView_DWebBrowserEvents2_QueryInterface(args *uintptr) uintptr {
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
}

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
