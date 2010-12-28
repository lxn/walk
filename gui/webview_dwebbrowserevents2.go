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

type webViewDWebBrowserEvents2Callbacks struct {
	QueryInterface             *syscall.Callback
	AddRef                     *syscall.Callback
	Release                    *syscall.Callback
	GetTypeInfoCount           *syscall.Callback
	GetTypeInfo                *syscall.Callback
	GetIDsOfNames              *syscall.Callback
	Invoke                     *syscall.Callback
	BeforeNavigate2            *syscall.Callback
	ClientToHostWindow         *syscall.Callback
	CommandStateChange         *syscall.Callback
	DocumentComplete           *syscall.Callback
	DownloadBegin              *syscall.Callback
	DownloadComplete           *syscall.Callback
	FileDownload               *syscall.Callback
	NavigateComplete2          *syscall.Callback
	NavigateError              *syscall.Callback
	NewProcess                 *syscall.Callback
	NewWindow2                 *syscall.Callback
	NewWindow3                 *syscall.Callback
	OnFullScreen               *syscall.Callback
	OnMenuBar                  *syscall.Callback
	OnQuit                     *syscall.Callback
	OnStatusBar                *syscall.Callback
	OnTheaterMode              *syscall.Callback
	OnToolBar                  *syscall.Callback
	OnVisible                  *syscall.Callback
	PrintTemplateInstantiation *syscall.Callback
	PrintTemplateTeardown      *syscall.Callback
	PrivacyImpactedStateChange *syscall.Callback
	ProgressChange             *syscall.Callback
	PropertyChange             *syscall.Callback
	RedirectXDomainBlocked     *syscall.Callback
	SetPhishingFilterStatus    *syscall.Callback
	SetSecureLockIcon          *syscall.Callback
	StatusTextChange           *syscall.Callback
	ThirdPartyUrlBlocked       *syscall.Callback
	TitleChange                *syscall.Callback
	UpdatePageStatus           *syscall.Callback
	WindowClosing              *syscall.Callback
	WindowSetHeight            *syscall.Callback
	WindowSetLeft              *syscall.Callback
	WindowSetResizable         *syscall.Callback
	WindowSetTop               *syscall.Callback
	WindowSetWidth             *syscall.Callback
	WindowStateChanged         *syscall.Callback
}

var webViewDWebBrowserEvents2Cbs = &webViewDWebBrowserEvents2Callbacks{
	syscall.NewCallback(webView_DWebBrowserEvents2_QueryInterface, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_AddRef, 4),
	syscall.NewCallback(webView_DWebBrowserEvents2_Release, 4),
	syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfoCount, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfo, 4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_GetIDsOfNames, 4+4+4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_Invoke, 4+4+4+4+2+4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_BeforeNavigate2, 4+4+4+4+4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_ClientToHostWindow, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_CommandStateChange, 4+4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_DocumentComplete, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_DownloadBegin, 4),
	syscall.NewCallback(webView_DWebBrowserEvents2_DownloadComplete, 4),
	syscall.NewCallback(webView_DWebBrowserEvents2_FileDownload, 4+2+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_NavigateComplete2, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_NavigateError, 4+4+4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_NewProcess, 4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_NewWindow2, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_NewWindow3, 4+4+4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnFullScreen, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnMenuBar, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnQuit, 4),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnStatusBar, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnTheaterMode, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnToolBar, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_OnVisible, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_PrintTemplateInstantiation, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_PrintTemplateTeardown, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_PrivacyImpactedStateChange, 4+16),
	syscall.NewCallback(webView_DWebBrowserEvents2_ProgressChange, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_PropertyChange, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_RedirectXDomainBlocked, 4+4+4+4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_SetPhishingFilterStatus, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_SetSecureLockIcon, 4+16),
	syscall.NewCallback(webView_DWebBrowserEvents2_StatusTextChange, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_ThirdPartyUrlBlocked, 4+4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_TitleChange, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_UpdatePageStatus, 4),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowClosing, 4+2+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetHeight, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetLeft, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetResizable, 4+2),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetTop, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowSetWidth, 4+4),
	syscall.NewCallback(webView_DWebBrowserEvents2_WindowStateChanged, 4+4+4),
}

var webViewDWebBrowserEvents2Vtbl *DWebBrowserEvents2Vtbl

func init() {
	webViewDWebBrowserEvents2Vtbl = &DWebBrowserEvents2Vtbl{
		uintptr(webViewDWebBrowserEvents2Cbs.QueryInterface.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.AddRef.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.Release.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.GetTypeInfoCount.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.GetTypeInfo.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.GetIDsOfNames.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.Invoke.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.BeforeNavigate2.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.ClientToHostWindow.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.CommandStateChange.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.DocumentComplete.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.DownloadBegin.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.DownloadComplete.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.FileDownload.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.NavigateComplete2.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.NavigateError.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.NewProcess.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.NewWindow2.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.NewWindow3.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnFullScreen.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnMenuBar.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnQuit.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnStatusBar.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnTheaterMode.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnToolBar.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.OnVisible.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.PrintTemplateInstantiation.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.PrintTemplateTeardown.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.PrivacyImpactedStateChange.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.ProgressChange.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.PropertyChange.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.RedirectXDomainBlocked.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.SetPhishingFilterStatus.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.SetSecureLockIcon.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.StatusTextChange.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.ThirdPartyUrlBlocked.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.TitleChange.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.UpdatePageStatus.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowClosing.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowSetHeight.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowSetLeft.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowSetResizable.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowSetTop.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowSetWidth.ExtFnEntry()),
		uintptr(webViewDWebBrowserEvents2Cbs.WindowStateChanged.ExtFnEntry()),
	}
}

type webViewDWebBrowserEvents2 struct {
	DWebBrowserEvents2
}

func webView_DWebBrowserEvents2_QueryInterface(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_QueryInterface")

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
	log.Println("webView_DWebBrowserEvents2_AddRef")

	return 1
}

func webView_DWebBrowserEvents2_Release(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_Release")

	return 1
}

func webView_DWebBrowserEvents2_GetTypeInfoCount(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_GetTypeInfoCount")

	/*	p := (*struct {
			wbe2    *webViewDWebBrowserEvents2
			pctinfo *uint
		})(unsafe.Pointer(args))

		*p.pctinfo = 0

		return S_OK*/

	return E_NOTIMPL
}

func webView_DWebBrowserEvents2_GetTypeInfo(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_GetTypeInfo")

	/*	p := (*struct {
				wbe2         *webViewDWebBrowserEvents2
			})(unsafe.Pointer(args))

		    unsigned int  iTInfo,         
		    LCID  lcid,                   
		    ITypeInfo FAR* FAR*  ppTInfo*/

	return E_NOTIMPL
}

func webView_DWebBrowserEvents2_GetIDsOfNames(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_GetIDsOfNames")

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
	log.Println("webView_DWebBrowserEvents2_Invoke")

	p := (*struct {
		wbe2         *webViewDWebBrowserEvents2
		dispIdMember DISPID
		riid         REFIID
		lcid         uint /* LCID */
		wFlags       uint16
		pDispParams  *DISPPARAMS
		pVarResult   *VARIANT
		pExcepInfo   unsafe.Pointer /* *EXCEPINFO */
		puArgErr     *uint
	})(unsafe.Pointer(args))

	log.Printf("p: %+v\n", p)

	return DISP_E_MEMBERNOTFOUND
}

func webView_DWebBrowserEvents2_BeforeNavigate2(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_BeforeNavigate2")

	return 0
}

func webView_DWebBrowserEvents2_ClientToHostWindow(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_ClientToHostWindow")

	return 0
}

func webView_DWebBrowserEvents2_CommandStateChange(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_CommandStateChange")

	return 0
}

func webView_DWebBrowserEvents2_DocumentComplete(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_DocumentComplete")

	return 0
}

func webView_DWebBrowserEvents2_DownloadBegin(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_DownloadBegin")

	return 0
}

func webView_DWebBrowserEvents2_DownloadComplete(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_DownloadComplete")

	return 0
}

func webView_DWebBrowserEvents2_FileDownload(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_FileDownload")

	return 0
}

func webView_DWebBrowserEvents2_NavigateComplete2(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_NavigateComplete2")

	return 0
}

func webView_DWebBrowserEvents2_NavigateError(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_NavigateError")

	return 0
}

func webView_DWebBrowserEvents2_NewProcess(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_NewProcess")

	return 0
}

func webView_DWebBrowserEvents2_NewWindow2(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_NewWindow2")

	return 0
}

func webView_DWebBrowserEvents2_NewWindow3(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_NewWindow3")

	return 0
}

func webView_DWebBrowserEvents2_OnFullScreen(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnFullScreen")

	return 0
}

func webView_DWebBrowserEvents2_OnMenuBar(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnMenuBar")

	return 0
}

func webView_DWebBrowserEvents2_OnQuit(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnQuit")

	return 0
}

func webView_DWebBrowserEvents2_OnStatusBar(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnStatusBar")

	return 0
}

func webView_DWebBrowserEvents2_OnTheaterMode(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnTheaterMode")

	return 0
}

func webView_DWebBrowserEvents2_OnToolBar(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnToolBar")

	return 0
}

func webView_DWebBrowserEvents2_OnVisible(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_OnVisible")

	return 0
}

func webView_DWebBrowserEvents2_PrintTemplateInstantiation(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_PrintTemplateInstantiation")

	return 0
}

func webView_DWebBrowserEvents2_PrintTemplateTeardown(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_PrintTemplateTeardown")

	return 0
}

func webView_DWebBrowserEvents2_PrivacyImpactedStateChange(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_PrivacyImpactedStateChange")

	return 0
}

func webView_DWebBrowserEvents2_ProgressChange(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_ProgressChange")

	return 0
}

func webView_DWebBrowserEvents2_PropertyChange(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_PropertyChange")

	return 0
}

func webView_DWebBrowserEvents2_RedirectXDomainBlocked(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_RedirectXDomainBlocked")

	return 0
}

func webView_DWebBrowserEvents2_SetPhishingFilterStatus(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_SetPhishingFilterStatus")

	return 0
}

func webView_DWebBrowserEvents2_SetSecureLockIcon(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_SetSecureLockIcon")

	return 0
}

func webView_DWebBrowserEvents2_StatusTextChange(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_StatusTextChange")

	return 0
}

func webView_DWebBrowserEvents2_ThirdPartyUrlBlocked(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_ThirdPartyUrlBlocked")

	return 0
}

func webView_DWebBrowserEvents2_TitleChange(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_TitleChange")

	return 0
}

func webView_DWebBrowserEvents2_UpdatePageStatus(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_UpdatePageStatus")

	return 0
}

func webView_DWebBrowserEvents2_WindowClosing(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowClosing")

	return 0
}

func webView_DWebBrowserEvents2_WindowSetHeight(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowSetHeight")

	return 0
}

func webView_DWebBrowserEvents2_WindowSetLeft(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowSetLeft")

	return 0
}

func webView_DWebBrowserEvents2_WindowSetResizable(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowSetResizable")

	return 0
}

func webView_DWebBrowserEvents2_WindowSetTop(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowSetTop")

	return 0
}

func webView_DWebBrowserEvents2_WindowSetWidth(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowSetWidth")

	return 0
}

func webView_DWebBrowserEvents2_WindowStateChanged(args *uintptr) uintptr {
	log.Println("webView_DWebBrowserEvents2_WindowStateChanged")

	return 0
}
