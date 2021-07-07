// Copyright 2021 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/lxn/win"
)

const openGLWindowClass = `\o/ Walk_OpenGL_Class \o/`

func init() {
	AppendToWalkInit(func() {
		MustRegisterWindowClass(openGLWindowClass)
	})
}

// OpenGLContext is passed to GL lifecycle callbacks
type OpenGLContext struct {
	widget *OpenGL
	hDc    win.HDC
}

// DC returns the device context handle
func (c *OpenGLContext) DC() win.HDC { return c.hDc }

// Handle returns the OpenGL rendering context handle
func (c *OpenGLContext) Handle() win.HGLRC { return c.widget.hGlrc }

// Widget returns the OpenGL widget
func (c *OpenGLContext) Widget() *OpenGL { return c.widget }

// GLFunc sets up, paints, or tears down OpenGL content
type GLFunc func(*OpenGLContext) error

// GLTickFunc returns whether a window repaint is required
type GLTickFunc func(*OpenGLContext) bool

type OpenGL struct {
	WidgetBase
	hGlrc win.HGLRC

	setup, paint, teardown GLFunc

	pixFmt       int32
	ppfd         *win.PIXELFORMATDESCRIPTOR
	contextAttrs []int32

	tickMu sync.Mutex
	tickFn GLTickFunc
	tickTk *time.Ticker
	tickCh chan struct{}
}

// NewOpenGL creates and initializes an OpenGL widget.
//
// pixFmtAttrs (can be null) is a list of attributes passed to
// wglCreateContextAttribsARB. contextAttrs (can be null) is a list of
// attributes passed to wglCreateContextAttribsARB.
//
// If wglCreateContextAttribsARB is unavailable, an equivalent pixel format
// descriptor will be used (attributes without a corresponding property will be
// ignored). If wglCreateContextAttribsARB is unavailable, wglCreateContext will
// be used.
func NewOpenGL(parent Container, style uint32, setup, paint, teardown GLFunc, pixFmtAttrs, contextAttrs []int32) (*OpenGL, error) {
	gl := &OpenGL{setup: setup, paint: paint, teardown: teardown}
	err := gl.init(parent, style)
	if err != nil {
		return nil, err
	}

	if pixFmtAttrs == nil {
		dc := win.GetDC(gl.hWnd)
		defer win.ReleaseDC(gl.hWnd, dc)

		pixFmtAttrs = []int32{
			win.WGL_SUPPORT_OPENGL_ARB, 1,
			win.WGL_DRAW_TO_WINDOW_ARB, 1,
			win.WGL_PIXEL_TYPE_ARB, win.WGL_TYPE_RGBA_ARB,
			win.WGL_COLOR_BITS_ARB, win.GetDeviceCaps(dc, win.BITSPIXEL),
			0,
		}
	} else {
		pixFmtAttrs = openGLValidateAttribs(pixFmtAttrs)
	}

	if contextAttrs == nil {
		gl.contextAttrs = []int32{
			win.WGL_CONTEXT_PROFILE_MASK_ARB, win.WGL_CONTEXT_CORE_PROFILE_BIT_ARB,
			0,
		}
	} else {
		gl.contextAttrs = openGLValidateAttribs(contextAttrs)
	}

	gl.pixFmt, gl.ppfd, err = openGLChoosePixelFormat(pixFmtAttrs)
	if err != nil {
		return nil, err
	}

	return gl, nil
}

func (gl *OpenGL) init(parent Container, style uint32) error {
	if err := InitWidget(
		gl,
		parent,
		customWidgetWindowClass,
		win.WS_VISIBLE|uint32(style),
		0); err != nil {
		return err
	}

	return nil
}

func (*OpenGL) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	return NewGreedyLayoutItem()
}

func (gl *OpenGL) context(dc win.HDC) *OpenGLContext {
	return &OpenGLContext{
		widget: gl,
		hDc:    dc,
	}
}

func (gl *OpenGL) createContext() bool {
	dc := win.GetDC(gl.hWnd)
	defer win.ReleaseDC(gl.hWnd, dc)

	if !win.SetPixelFormat(dc, gl.pixFmt, gl.ppfd) {
		processError(lastError("SetPixelFormat"))
		return false
	}

	if win.HasWglCreateContextAttribsARB() {
		gl.hGlrc = win.WglCreateContextAttribsARB(dc, 0, &gl.contextAttrs[0])
		if gl.hGlrc == 0 {
			processError(lastError("WglCreateContextAttribsARB"))
			return false
		}
	} else {
		gl.hGlrc = win.WglCreateContext(dc)
		if gl.hGlrc == 0 {
			processError(lastError("WglCreateContext"))
			return false
		}
	}

	win.WglMakeCurrent(dc, gl.hGlrc)
	defer win.WglMakeCurrent(0, 0)

	if gl.setup != nil {
		err := gl.setup(gl.context(dc))
		if err != nil {
			newError(fmt.Sprintf("setup func failed: %v", err))
			return false
		}
	}

	return true
}

func (gl *OpenGL) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_PAINT:
		if gl.hGlrc == 0 {
			if !gl.createContext() {
				break
			}
		}

		if gl.paint == nil {
			newError("paint func is nil")
			break
		}

		var ps win.PAINTSTRUCT
		var dc win.HDC
		if wParam == 0 {
			dc = win.BeginPaint(gl.hWnd, &ps)
			if dc == 0 {
				processError(lastError("BeginPaint"))
				break
			}
			defer win.EndPaint(gl.hWnd, &ps)
		} else {
			dc = win.HDC(wParam)
		}

		win.WglMakeCurrent(dc, gl.hGlrc)
		defer win.WglMakeCurrent(0, 0)

		err := gl.paint(gl.context(dc))
		if err != nil {
			newError(fmt.Sprintf("paint func failed: %v", err))
		}

		return 0

	case win.WM_DESTROY:
		if gl.teardown != nil {
			dc := win.GetDC(gl.hWnd)
			win.WglMakeCurrent(dc, gl.hGlrc)

			err := gl.teardown(gl.context(dc))
			if err != nil {
				newError(fmt.Sprintf("teardown func failed: %v", err))
			}

			win.WglMakeCurrent(0, 0)
			win.ReleaseDC(gl.hWnd, dc)
		}

		win.WglDeleteContext(gl.hGlrc)
		gl.hGlrc = 0
		break // continue with base teardown
	}

	return gl.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

func openGLValidateAttribs(attribs []int32) []int32 {
	if len(attribs) == 0 {
		// length zero, no attributes
		return []int32{0}
	}

	if len(attribs)%2 == 0 {
		// length divisible by 2, add null terminator
		return append(attribs, 0)
	}

	if attribs[len(attribs)-1] != 0 {
		// length is not divisible by 2, but last element is not null
		panic("WGL attributes list must be null terminated or contain an even number of elements")
	}

	return attribs
}

func openGLChoosePixelFormat(attribs []int32) (int32, *win.PIXELFORMATDESCRIPTOR, error) {
	// create dummy window, dispose after
	w, err := NewMainWindow()
	if err != nil {
		return 0, nil, err
	}
	defer w.Dispose()

	// get the device context, release after
	hDC := win.GetDC(w.hWnd)
	if hDC == 0 {
		return 0, nil, lastError("GetDC")
	}
	defer win.ReleaseDC(w.hWnd, hDC)

	// legacy pixel format
	pfd := win.PIXELFORMATDESCRIPTOR{
		NSize:       uint16(unsafe.Sizeof(win.PIXELFORMATDESCRIPTOR{})),
		NVersion:    1,
		DwLayerMask: win.PFD_MAIN_PLANE,
	}
	wglConvertAttributes(&pfd, attribs)

	pixelFormat := win.ChoosePixelFormat(hDC, &pfd)
	if pixelFormat == 0 {
		return 0, nil, lastError("ChoosePixelFormat")
	}

	if !win.SetPixelFormat(hDC, pixelFormat, &pfd) {
		return 0, nil, lastError("SetPixelFormat")
	}

	// create dummy context, dispose after
	hRC := win.WglCreateContext(hDC)
	if hRC == 0 {
		return 0, nil, lastError("WglCreateContext")
	}
	defer win.WglDeleteContext(hRC)

	if !win.WglMakeCurrent(hDC, hRC) {
		return 0, nil, lastError("WglMakeCurrent")
	}
	defer win.WglMakeCurrent(0, 0)

	// get WGL extension functions
	win.InitWglExt()

	// use the legacy format if the extension is not available
	if !win.HasWglChoosePixelFormatARB() {
		return pixelFormat, &pfd, nil
	}

	var formatARB int32
	var numFormats uint32
	if !win.WglChoosePixelFormatARB(hDC, &attribs[0], nil, 1, &formatARB, &numFormats) {
		return 0, nil, lastError("WglChoosePixelFormatARB")
	}

	// use the legacy format if no acceptable formats are found
	if numFormats == 0 {
		return pixelFormat, &pfd, nil
	}

	// update the PFD
	if !win.DescribePixelFormat(hDC, formatARB, uint32(unsafe.Sizeof(pfd)), &pfd) {
		return 0, nil, lastError("DescribePixelFormat")
	}

	return formatARB, &pfd, nil
}

func wglConvertAttributes(ppfd *win.PIXELFORMATDESCRIPTOR, attribs []int32) {
	amap := map[int32]int32{}
	for i, n := 0, len(attribs)/2; i < n; i++ {
		amap[attribs[2*i]] = attribs[2*i+1]
	}

	ppfd.CColorBits = byte(amap[win.WGL_COLOR_BITS_ARB])
	ppfd.CRedBits = byte(amap[win.WGL_RED_BITS_ARB])
	ppfd.CRedShift = byte(amap[win.WGL_RED_SHIFT_ARB])
	ppfd.CGreenBits = byte(amap[win.WGL_GREEN_BITS_ARB])
	ppfd.CGreenShift = byte(amap[win.WGL_GREEN_SHIFT_ARB])
	ppfd.CBlueBits = byte(amap[win.WGL_BLUE_BITS_ARB])
	ppfd.CBlueShift = byte(amap[win.WGL_BLUE_SHIFT_ARB])
	ppfd.CAlphaBits = byte(amap[win.WGL_ALPHA_BITS_ARB])
	ppfd.CAlphaShift = byte(amap[win.WGL_ALPHA_SHIFT_ARB])
	ppfd.CAccumBits = byte(amap[win.WGL_ACCUM_BITS_ARB])
	ppfd.CAccumRedBits = byte(amap[win.WGL_ACCUM_RED_BITS_ARB])
	ppfd.CAccumGreenBits = byte(amap[win.WGL_ACCUM_GREEN_BITS_ARB])
	ppfd.CAccumBlueBits = byte(amap[win.WGL_ACCUM_BLUE_BITS_ARB])
	ppfd.CAccumAlphaBits = byte(amap[win.WGL_ACCUM_ALPHA_BITS_ARB])
	ppfd.CDepthBits = byte(amap[win.WGL_DEPTH_BITS_ARB])
	ppfd.CStencilBits = byte(amap[win.WGL_STENCIL_BITS_ARB])
	ppfd.CAuxBuffers = byte(amap[win.WGL_AUX_BUFFERS_ARB])

	if amap[win.WGL_SUPPORT_OPENGL_ARB] == 1 {
		ppfd.DwFlags |= win.PFD_SUPPORT_OPENGL
	}

	if amap[win.WGL_DRAW_TO_WINDOW_ARB] == 1 {
		ppfd.DwFlags |= win.PFD_DRAW_TO_WINDOW
	}

	switch amap[win.WGL_PIXEL_TYPE_ARB] {
	case win.WGL_TYPE_RGBA_ARB, 0:
		ppfd.IPixelType = win.PFD_TYPE_RGBA
	case win.WGL_TYPE_COLORINDEX_ARB:
		ppfd.IPixelType = win.PFD_TYPE_COLORINDEX
	}
}
