// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"syscall"
	"unsafe"
)

import (
	"github.com/lxn/win"
	"time"
)

var webViewDWebBrowserEvents2Vtbl *win.DWebBrowserEvents2Vtbl

func init() {
	webViewDWebBrowserEvents2Vtbl = &win.DWebBrowserEvents2Vtbl{
		syscall.NewCallback(webView_DWebBrowserEvents2_QueryInterface),
		syscall.NewCallback(webView_DWebBrowserEvents2_AddRef),
		syscall.NewCallback(webView_DWebBrowserEvents2_Release),
		syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfoCount),
		syscall.NewCallback(webView_DWebBrowserEvents2_GetTypeInfo),
		syscall.NewCallback(webView_DWebBrowserEvents2_GetIDsOfNames),
		syscall.NewCallback(webView_DWebBrowserEvents2_Invoke),
	}
}

type webViewDWebBrowserEvents2 struct {
	win.DWebBrowserEvents2
}

func webView_DWebBrowserEvents2_QueryInterface(wbe2 *webViewDWebBrowserEvents2, riid win.REFIID, ppvObject *unsafe.Pointer) uintptr {
	// Just reuse the QueryInterface implementation we have for IOleClientSite.
	// We need to adjust object, which initially points at our
	// webViewDWebBrowserEvents2, so it refers to the containing
	// webViewIOleClientSite for the call.
	var clientSite win.IOleClientSite
	var webViewInPlaceSite webViewIOleInPlaceSite
	var docHostUIHandler webViewIDocHostUIHandler

	ptr := uintptr(unsafe.Pointer(wbe2)) -
		uintptr(unsafe.Sizeof(clientSite)) -
		uintptr(unsafe.Sizeof(webViewInPlaceSite)) -
		uintptr(unsafe.Sizeof(docHostUIHandler))

	return webView_IOleClientSite_QueryInterface((*webViewIOleClientSite)(unsafe.Pointer(ptr)), riid, ppvObject)
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

	return win.E_NOTIMPL
}

func webView_DWebBrowserEvents2_GetTypeInfo(args *uintptr) uintptr {
	/*	p := (*struct {
				wbe2         *webViewDWebBrowserEvents2
			})(unsafe.Pointer(args))

		    unsigned int  iTInfo,
		    LCID  lcid,
		    ITypeInfo FAR* FAR*  ppTInfo*/

	return win.E_NOTIMPL
}

func webView_DWebBrowserEvents2_GetIDsOfNames(args *uintptr) uintptr {
	/*	p := (*struct {
		wbe2      *webViewDWebBrowserEvents2
		riid      REFIID
		rgszNames **uint16
		cNames    uint32
		lcid      LCID
		rgDispId  *DISPID
	})(unsafe.Pointer(args))*/

	return win.E_NOTIMPL
}

/*
func webView_DWebBrowserEvents2_Invoke(
	wbe2 *webViewDWebBrowserEvents2,
	dispIdMember win.DISPID,
	riid win.REFIID,
	lcid uint32, // LCID
	wFlags uint16,
	pDispParams *win.DISPPARAMS,
	pVarResult *win.VARIANT,
	pExcepInfo unsafe.Pointer, // *EXCEPINFO
	puArgErr *uint32) uintptr {
*/
func webView_DWebBrowserEvents2_Invoke(
	arg0 uintptr,
	arg1 uintptr,
	arg2 uintptr,
	arg3 uintptr,
	arg4 uintptr,
	arg5 uintptr,
	arg6 uintptr,
	arg7 uintptr,
	arg8 uintptr) uintptr {

	wbe2 := (*webViewDWebBrowserEvents2)(unsafe.Pointer(arg0))
	dispIdMember := *(*win.DISPID)(unsafe.Pointer(&arg1))
	//riid := *(*win.REFIID)(unsafe.Pointer(&arg2))
	//lcid := *(*uint32)(unsafe.Pointer(&arg3))
	//wFlags := *(*uint16)(unsafe.Pointer(&arg4))
	pDispParams := (*win.DISPPARAMS)(unsafe.Pointer(arg5))
	//pVarResult := (*win.VARIANT)(unsafe.Pointer(arg6))
	//pExcepInfo := unsafe.Pointer(arg7)
	//puArgErr := (*uint32)(unsafe.Pointer(arg8))

	var wb WidgetBase
	var wvcs webViewIOleClientSite

	wv := (*WebView)(unsafe.Pointer(uintptr(unsafe.Pointer(wbe2)) +
		uintptr(unsafe.Sizeof(*wbe2)) -
		uintptr(unsafe.Sizeof(wvcs)) -
		uintptr(unsafe.Sizeof(wb))))

	switch dispIdMember {
	case win.DISPID_BEFORENAVIGATE2:
		rgvargPtr := (*[7]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		arg := &WebViewNavigatingArg{
			pDisp:           (*rgvargPtr)[6].MustPDispatch(),
			url:             (*rgvargPtr)[5].MustPVariant(),
			flags:           (*rgvargPtr)[4].MustPVariant(),
			targetFrameName: (*rgvargPtr)[3].MustPVariant(),
			postData:        (*rgvargPtr)[2].MustPVariant(),
			headers:         (*rgvargPtr)[1].MustPVariant(),
			cancel:          (*rgvargPtr)[0].MustPBool(),
		}
		wv.navigatingEventPublisher.Publish(arg)

	case win.DISPID_NAVIGATECOMPLETE2:
		rgvargPtr := (*[2]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		url := (*rgvargPtr)[0].MustPVariant()
		urlStr := ""
		if url != nil && url.MustBSTR() != nil {
			urlStr = win.BSTRToString(url.MustBSTR())
		}
		wv.navigatedEventPublisher.Publish(urlStr)

		wv.urlChangedPublisher.Publish()

	case win.DISPID_DOWNLOADBEGIN:
		wv.downloadingEventPublisher.Publish()

	case win.DISPID_DOWNLOADCOMPLETE:
		wv.downloadedEventPublisher.Publish()

	case win.DISPID_DOCUMENTCOMPLETE:
		rgvargPtr := (*[2]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		url := (*rgvargPtr)[0].MustPVariant()
		urlStr := ""
		if url != nil && url.MustBSTR() != nil {
			urlStr = win.BSTRToString(url.MustBSTR())
		}

		// FIXME: Horrible hack to avoid glitch where the document is not displayed.
		time.AfterFunc(time.Millisecond*100, func() {
			wv.Synchronize(func() {
				b := wv.Bounds()
				b.Width++
				wv.SetBounds(b)
				b.Width--
				wv.SetBounds(b)
			})
		})

		wv.documentCompletedEventPublisher.Publish(urlStr)

	case win.DISPID_NAVIGATEERROR:
		rgvargPtr := (*[5]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		arg := &WebViewNavigatedErrorEventArg{
			pDisp:           (*rgvargPtr)[4].MustPDispatch(),
			url:             (*rgvargPtr)[3].MustPVariant(),
			targetFrameName: (*rgvargPtr)[2].MustPVariant(),
			statusCode:      (*rgvargPtr)[1].MustPVariant(),
			cancel:          (*rgvargPtr)[0].MustPBool(),
		}
		wv.navigatedErrorEventPublisher.Publish(arg)

	case win.DISPID_NEWWINDOW3:
		rgvargPtr := (*[5]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		arg := &WebViewNewWindowEventArg{
			ppDisp:         (*rgvargPtr)[4].MustPPDispatch(),
			cancel:         (*rgvargPtr)[3].MustPBool(),
			dwFlags:        (*rgvargPtr)[2].MustULong(),
			bstrUrlContext: (*rgvargPtr)[1].MustBSTR(),
			bstrUrl:        (*rgvargPtr)[0].MustBSTR(),
		}
		wv.newWindowEventPublisher.Publish(arg)

	case win.DISPID_ONQUIT:
		wv.quittingEventPublisher.Publish()

	case win.DISPID_WINDOWCLOSING:
		rgvargPtr := (*[2]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		arg := &WebViewWindowClosingEventArg{
			bIsChildWindow: (*rgvargPtr)[1].MustBool(),
			cancel:         (*rgvargPtr)[0].MustPBool(),
		}
		wv.windowClosingEventPublisher.Publish(arg)

	case win.DISPID_ONSTATUSBAR:
		rgvargPtr := (*[1]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		statusBar := (*rgvargPtr)[0].MustBool()
		if statusBar != win.VARIANT_FALSE {
			wv.statusBarVisible = true
		} else {
			wv.statusBarVisible = false
		}
		wv.statusBarVisibleChangedEventPublisher.Publish()

	case win.DISPID_ONTHEATERMODE:
		rgvargPtr := (*[1]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		theaterMode := (*rgvargPtr)[0].MustBool()
		if theaterMode != win.VARIANT_FALSE {
			wv.isTheaterMode = true
		} else {
			wv.isTheaterMode = false
		}
		wv.theaterModeChangedEventPublisher.Publish()

	case win.DISPID_ONTOOLBAR:
		rgvargPtr := (*[1]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		toolBar := (*rgvargPtr)[0].MustBool()
		if toolBar != win.VARIANT_FALSE {
			wv.toolBarVisible = true
		} else {
			wv.toolBarVisible = false
		}
		wv.toolBarVisibleChangedEventPublisher.Publish()

	case win.DISPID_ONVISIBLE:
		rgvargPtr := (*[1]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		vVisible := (*rgvargPtr)[0].MustBool()
		if vVisible != win.VARIANT_FALSE {
			wv.browserVisible = true
		} else {
			wv.browserVisible = false
		}
		wv.browserVisibleChangedEventPublisher.Publish()

	case win.DISPID_COMMANDSTATECHANGE:
		rgvargPtr := (*[2]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		arg := &WebViewCommandStateChangedEventArg{
			command: (*rgvargPtr)[1].MustLong(),
			enable:  (*rgvargPtr)[0].MustBool(),
		}
		wv.commandStateChangedEventPublisher.Publish(arg)

	case win.DISPID_PROGRESSCHANGE:
		rgvargPtr := (*[2]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		wv.progressValue = (*rgvargPtr)[1].MustLong()
		wv.progressMax = (*rgvargPtr)[0].MustLong()
		wv.progressChangedEventPublisher.Publish()

	case win.DISPID_STATUSTEXTCHANGE:
		rgvargPtr := (*[1]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		sText := (*rgvargPtr)[0].MustBSTR()
		if sText != nil {
			wv.statusText = win.BSTRToString(sText)
		} else {
			wv.statusText = ""
		}
		wv.statusTextChangedEventPublisher.Publish()

	case win.DISPID_TITLECHANGE:
		rgvargPtr := (*[1]win.VARIANTARG)(unsafe.Pointer(pDispParams.Rgvarg))
		sText := (*rgvargPtr)[0].MustBSTR()
		if sText != nil {
			wv.documentTitle = win.BSTRToString(sText)
		} else {
			wv.documentTitle = ""
		}
		wv.documentTitleChangedEventPublisher.Publish()
	}

	return win.DISP_E_MEMBERNOTFOUND
}
