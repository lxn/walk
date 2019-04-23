// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"image"
	"path/filepath"
	"syscall"
)

import (
	"github.com/lxn/win"
)

var defaultIconSize Size

func init() {
	defaultIconSize = Size{int(win.GetSystemMetrics(win.SM_CXICON)), int(win.GetSystemMetrics(win.SM_CYICON))}
}

// Icon is a bitmap that supports transparency and combining multiple
// variants of an image in different resolutions.
type Icon struct {
	hIcon   win.HICON
	size    Size
	isStock bool
}

func IconApplication() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION)), defaultIconSize, true}
}

func IconError() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_ERROR)), defaultIconSize, true}
}

func IconQuestion() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_QUESTION)), defaultIconSize, true}
}

func IconWarning() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_WARNING)), defaultIconSize, true}
}

func IconInformation() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_INFORMATION)), defaultIconSize, true}
}

func IconWinLogo() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_WINLOGO)), defaultIconSize, true}
}

func IconShield() *Icon {
	return &Icon{win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_SHIELD)), defaultIconSize, true}
}

// NewIconFromFile returns a new Icon, using the specified icon image file and default size.
func NewIconFromFile(filePath string) (*Icon, error) {
	return NewIconFromFileWithSize(filePath, Size{})
}

// NewIconFromFileWithSize returns a new Icon, using the specified icon image file and size.
func NewIconFromFileWithSize(filePath string, size Size) (*Icon, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, wrapError(err)
	}

	flags := win.LR_LOADFROMFILE
	if size.Width == 0 || size.Height == 0 {
		flags |= win.LR_DEFAULTSIZE
	}

	hIcon := win.HICON(win.LoadImage(
		0,
		syscall.StringToUTF16Ptr(absFilePath),
		win.IMAGE_ICON,
		int32(size.Width),
		int32(size.Height),
		uint32(flags)))
	if hIcon == 0 {
		return nil, lastError("LoadImage")
	}

	if size.Width == 0 || size.Height == 0 {
		size = defaultIconSize
	}

	return &Icon{hIcon: hIcon, size: size}, nil
}

// NewIconFromResource returns a new Icon of default size, using the specified icon resource.
func NewIconFromResource(name string) (*Icon, error) {
	return NewIconFromResourceWithSize(name, Size{})
}

// NewIconFromResourceWithSize returns a new Icon of size size, using the specified icon resource.
func NewIconFromResourceWithSize(name string, size Size) (*Icon, error) {
	return newIconFromResource(syscall.StringToUTF16Ptr(name), size)
}

// NewIconFromResourceId returns a new Icon of default size, using the specified icon resource.
func NewIconFromResourceId(id int) (*Icon, error) {
	return NewIconFromResourceIdWithSize(id, Size{})
}

// NewIconFromResourceIdWithSize returns a new Icon of size size, using the specified icon resource.
func NewIconFromResourceIdWithSize(id int, size Size) (*Icon, error) {
	return newIconFromResource(win.MAKEINTRESOURCE(uintptr(id)), size)
}

func newIconFromResource(res *uint16, size Size) (ic *Icon, err error) {
	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		err = lastError("GetModuleHandle")
		return
	}

	var flags uint32
	if size.Width == 0 || size.Height == 0 {
		flags |= win.LR_DEFAULTSIZE
	}

	hIcon := win.HICON(win.LoadImage(
		hInst,
		res,
		win.IMAGE_ICON,
		int32(size.Width),
		int32(size.Height),
		flags))
	if hIcon == 0 {
		return nil, lastError("LoadImage")
	}

	if size.Width == 0 || size.Height == 0 {
		size = defaultIconSize
	}

	ic = &Icon{hIcon: hIcon, size: size}

	return
}

// NewIconFromImage returns a new Icon, using the specified image.Image as source.
func NewIconFromImage(im image.Image) (ic *Icon, err error) {
	hIcon, err := createAlphaCursorOrIconFromImage(im, image.Pt(0, 0), true)
	if err != nil {
		return nil, err
	}
	return &Icon{hIcon: hIcon}, nil
}

// NewIconFromBitmap returns a new Icon, using the specified Bitmap as source.
func NewIconFromBitmap(bmp *Bitmap) (ic *Icon, err error) {
	hIcon, err := createAlphaCursorOrIconFromBitmap(bmp, Point{}, true)
	if err != nil {
		return nil, err
	}
	return &Icon{hIcon: hIcon}, nil
}

// NewIconFromHICON returns a new Icon, using the specified win.HICON as source.
func NewIconFromHICON(hIcon win.HICON) (ic *Icon, err error) {
	return &Icon{hIcon: hIcon}, nil
}

// Dispose releases the operating system resources associated with the Icon.
func (i *Icon) Dispose() {
	if i.isStock || i.hIcon == 0 {
		return
	}

	win.DestroyIcon(i.hIcon)
	i.hIcon = 0
}

func (i *Icon) draw(hdc win.HDC, location Point) error {
	s := i.Size()

	return i.drawStretched(hdc, Rectangle{location.X, location.Y, s.Width, s.Height})
}

func (i *Icon) drawStretched(hdc win.HDC, bounds Rectangle) error {
	if !win.DrawIconEx(hdc, int32(bounds.X), int32(bounds.Y), i.hIcon, int32(bounds.Width), int32(bounds.Height), 0, 0, win.DI_NORMAL) {
		return lastError("DrawIconEx")
	}

	return nil
}

// Size returns the size of the Icon.
func (i *Icon) Size() Size {
	return i.size
}

// create an Alpha Icon or Cursor from an Image
// http://support.microsoft.com/kb/318876
func createAlphaCursorOrIconFromImage(im image.Image, hotspot image.Point, fIcon bool) (win.HICON, error) {
	bmp, err := NewBitmapFromImage(im)
	if err != nil {
		return 0, err
	}
	defer bmp.Dispose()

	return createAlphaCursorOrIconFromBitmap(bmp, Point{hotspot.X, hotspot.Y}, fIcon)
}

func createAlphaCursorOrIconFromBitmap(bmp *Bitmap, hotspot Point, fIcon bool) (win.HICON, error) {
	// Create an empty mask bitmap.
	size := bmp.Size()
	hMonoBitmap := win.CreateBitmap(int32(size.Width), int32(size.Height), 1, 1, nil)
	if hMonoBitmap == 0 {
		return 0, newError("CreateBitmap failed")
	}
	defer win.DeleteObject(win.HGDIOBJ(hMonoBitmap))

	var ii win.ICONINFO
	if fIcon {
		ii.FIcon = win.TRUE
	}
	ii.XHotspot = uint32(hotspot.X)
	ii.YHotspot = uint32(hotspot.Y)
	ii.HbmMask = hMonoBitmap
	ii.HbmColor = bmp.hBmp

	// Create the alpha cursor with the alpha DIB section.
	hIconOrCursor := win.CreateIconIndirect(&ii)

	return hIconOrCursor, nil
}
