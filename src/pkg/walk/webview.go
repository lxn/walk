// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import . "walk/winapi"

const webViewWindowClass = `\o/ Walk_WebView_Class \o/`

var webViewWindowClassRegistered bool

type WebView struct {
	WidgetBase
	clientSite    webViewIOleClientSite
	browserObject *IOleObject
}

func NewWebView(parent Container) (*WebView, os.Error) {
	ensureRegisteredWindowClass(webViewWindowClass, &webViewWindowClassRegistered)

	wv := &WebView{
		clientSite: webViewIOleClientSite{
			IOleClientSite: IOleClientSite{
				LpVtbl: webViewIOleClientSiteVtbl,
			},
			inPlaceSite: webViewIOleInPlaceSite{
				IOleInPlaceSite: IOleInPlaceSite{
					LpVtbl: webViewIOleInPlaceSiteVtbl,
				},
				inPlaceFrame: webViewIOleInPlaceFrame{
					IOleInPlaceFrame: IOleInPlaceFrame{
						LpVtbl: webViewIOleInPlaceFrameVtbl,
					},
				},
			},
			docHostUIHandler: webViewIDocHostUIHandler{
				IDocHostUIHandler: IDocHostUIHandler{
					LpVtbl: webViewIDocHostUIHandlerVtbl,
				},
			},
			webBrowserEvents2: webViewDWebBrowserEvents2{
				DWebBrowserEvents2: DWebBrowserEvents2{
					LpVtbl: webViewDWebBrowserEvents2Vtbl,
				},
			},
		},
	}

	if err := initChildWidget(
		wv,
		parent,
		webViewWindowClass,
		WS_VISIBLE,
		0); err != nil {
		return nil, err
	}

	wv.clientSite.inPlaceSite.inPlaceFrame.webView = wv

	succeeded := false

	defer func() {
		if !succeeded {
			wv.Dispose()
		}
	}()

	var classFactoryPtr unsafe.Pointer
	if hr := CoGetClassObject(&CLSID_WebBrowser, CLSCTX_INPROC_HANDLER|CLSCTX_INPROC_SERVER, nil, &IID_IClassFactory, &classFactoryPtr); FAILED(hr) {
		return nil, errorFromHRESULT("CoGetClassObject", hr)
	}
	classFactory := (*IClassFactory)(classFactoryPtr)
	defer classFactory.Release()

	var browserObjectPtr unsafe.Pointer
	if hr := classFactory.CreateInstance(nil, &IID_IOleObject, &browserObjectPtr); FAILED(hr) {
		return nil, errorFromHRESULT("IClassFactory.CreateInstance", hr)
	}
	browserObject := (*IOleObject)(browserObjectPtr)

	wv.browserObject = browserObject

	if hr := browserObject.SetClientSite((*IOleClientSite)(unsafe.Pointer(&wv.clientSite))); FAILED(hr) {
		return nil, errorFromHRESULT("IOleObject.SetClientSite", hr)
	}

	if hr := browserObject.SetHostNames(syscall.StringToUTF16Ptr("Walk.WebView"), nil); FAILED(hr) {
		return nil, errorFromHRESULT("IOleObject.SetHostNames", hr)
	}

	if hr := OleSetContainedObject((*IUnknown)(unsafe.Pointer(browserObject)), true); FAILED(hr) {
		return nil, errorFromHRESULT("OleSetContainedObject", hr)
	}

	var rect RECT
	GetClientRect(wv.hWnd, &rect)

	if hr := browserObject.DoVerb(OLEIVERB_SHOW, nil, (*IOleClientSite)(unsafe.Pointer(&wv.clientSite)), -1, wv.hWnd, &rect); FAILED(hr) {
		return nil, errorFromHRESULT("IOleObject.DoVerb", hr)
	}

	// FIXME: Reactivate after fixing crash
	/*	var cpcPtr unsafe.Pointer
		if hr := browserObject.QueryInterface(&IID_IConnectionPointContainer, &cpcPtr); FAILED(hr) {
			return nil, errorFromHRESULT("IOleObject.QueryInterface(IID_IConnectionPointContainer)", hr)
		}
		cpc := (*IConnectionPointContainer)(cpcPtr)
		defer cpc.Release()

		var cp *IConnectionPoint
		if hr := cpc.FindConnectionPoint(&DIID_DWebBrowserEvents2, &cp); FAILED(hr) {
			return nil, errorFromHRESULT("IConnectionPointContainer.FindConnectionPoint(DIID_DWebBrowserEvents2)", hr)
		}
		defer cp.Release()

		var cookie uint
		if hr := cp.Advise(unsafe.Pointer(&wv.clientSite.webBrowserEvents2), &cookie); FAILED(hr) {
			return nil, errorFromHRESULT("IConnectionPoint.Advise", hr)
		}*/

	wv.onResize()

	succeeded = true

	return wv, nil
}

func (wv *WebView) Dispose() {
	if wv.browserObject != nil {
		wv.browserObject.Close(OLECLOSE_NOSAVE)
		wv.browserObject.Release()

		wv.browserObject = nil
	}

	wv.WidgetBase.Dispose()
}

func (*WebView) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (*WebView) SizeHint() Size {
	return Size{100, 100}
}

func (wv *WebView) URL() (url string, err os.Error) {
	err = wv.withWebBrowser2(func(webBrowser2 *IWebBrowser2) os.Error {
		var urlBstr *uint16 /*BSTR*/
		if hr := webBrowser2.Get_LocationURL(&urlBstr); FAILED(hr) {
			return errorFromHRESULT("IWebBrowser2.Get_LocationURL", hr)
		}
		defer SysFreeString(urlBstr)

		url = BSTRToString(urlBstr)

		return nil
	})

	return
}

func (wv *WebView) SetURL(url string) os.Error {
	return wv.withWebBrowser2(func(webBrowser2 *IWebBrowser2) os.Error {
		urlBstr := StringToVariantBSTR(url)
		flags := IntToVariantI4(0)
		targetFrameName := StringToVariantBSTR("_self")

		if hr := webBrowser2.Navigate2(urlBstr, flags, targetFrameName, nil, nil); FAILED(hr) {
			return errorFromHRESULT("IWebBrowser2.Navigate2", hr)
		}

		return nil
	})
}

func (wv *WebView) withWebBrowser2(f func(webBrowser2 *IWebBrowser2) os.Error) os.Error {
	var webBrowser2Ptr unsafe.Pointer
	if hr := wv.browserObject.QueryInterface(&IID_IWebBrowser2, &webBrowser2Ptr); FAILED(hr) {
		return errorFromHRESULT("IOleObject.QueryInterface", hr)
	}
	webBrowser2 := (*IWebBrowser2)(webBrowser2Ptr)
	defer webBrowser2.Release()

	return f(webBrowser2)
}

func (wv *WebView) onResize() {
	// FIXME: handle error?
	wv.withWebBrowser2(func(webBrowser2 *IWebBrowser2) os.Error {
		bounds := wv.ClientBounds()

		webBrowser2.Put_Left(0)
		webBrowser2.Put_Top(0)
		webBrowser2.Put_Width(bounds.Width)
		webBrowser2.Put_Height(bounds.Height)

		return nil
	})
}

func (wv *WebView) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_SIZE, WM_SIZING:
		if wv.clientSite.inPlaceSite.inPlaceFrame.webView == nil {
			break
		}

		wv.onResize()
	}

	return wv.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}
