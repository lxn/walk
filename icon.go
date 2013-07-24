// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"image"
	"path/filepath"
	"syscall"
)

import (
	. "github.com/lxn/go-winapi"
)

// Icon is a bitmap that supports transparency and combining multiple
// variants of an image in different resolutions.
type Icon struct {
	hIcon   HICON
	isStock bool
}

func IconApplication() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_APPLICATION)), true}
}

func IconError() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_ERROR)), true}
}

func IconQuestion() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_QUESTION)), true}
}

func IconWarning() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_WARNING)), true}
}

func IconInformation() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_INFORMATION)), true}
}

func IconWinLogo() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_WINLOGO)), true}
}

func IconShield() *Icon {
	return &Icon{LoadIcon(0, MAKEINTRESOURCE(IDI_SHIELD)), true}
}

// NewIconFromFile returns a new Icon, using the specified icon image file.
func NewIconFromFile(filePath string) (*Icon, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, wrapError(err)
	}

	hIcon := HICON(LoadImage(
		0,
		syscall.StringToUTF16Ptr(absFilePath),
		IMAGE_ICON,
		0,
		0,
		LR_DEFAULTSIZE|LR_LOADFROMFILE))
	if hIcon == 0 {
		return nil, lastError("LoadImage")
	}

	return &Icon{hIcon: hIcon}, nil
}

// NewIconFromResource returns a new Icon, using the specified icon resource.
func NewIconFromResource(resName string) (ic *Icon, err error) {
	hInst := GetModuleHandle(nil)
	if hInst == 0 {
		err = lastError("GetModuleHandle")
		return
	}
	if hIcon := LoadIcon(hInst, syscall.StringToUTF16Ptr(resName)); hIcon == 0 {
		err = lastError("LoadIcon")
	} else {
		ic = &Icon{hIcon: hIcon}
	}
	return
}

func NewIconFromImage(im image.Image) (ic *Icon, err error) {
	hIcon, err := createAlphaCursorOrIconFromImage(im, image.Pt(0, 0), true)
	if err != nil {
		return nil, err
	}
	return &Icon{hIcon: hIcon}, nil
}

// Dispose releases the operating system resources associated with the Icon.
func (i *Icon) Dispose() error {
	if i.isStock || i.hIcon == 0 {
		return nil
	}

	if !DestroyIcon(i.hIcon) {
		return lastError("DestroyIcon")
	}

	i.hIcon = 0

	return nil
}

// create an Alpha Icon or Cursor from an Image
// http://support.microsoft.com/kb/318876
func createAlphaCursorOrIconFromImage(im image.Image, hotspot image.Point, fIcon bool) (HICON, error) {

	hBitmap, err := hBitmapFromImage(im)
	if err != nil {
		return 0, err
	}
	defer DeleteObject(HGDIOBJ(hBitmap))

	// Create an empty mask bitmap.
	hMonoBitmap := CreateBitmap(int32(im.Bounds().Dx()), int32(im.Bounds().Dy()), 1, 1, nil)
	if hMonoBitmap == 0 {
		return 0, newError("CreateBitmap failed")
	}
	defer DeleteObject(HGDIOBJ(hMonoBitmap))

	var ii ICONINFO
	if fIcon {
		ii.FIcon = TRUE
	}
	ii.XHotspot = uint32(hotspot.X)
	ii.YHotspot = uint32(hotspot.Y)
	ii.HbmMask = hMonoBitmap
	ii.HbmColor = hBitmap

	// Create the alpha cursor with the alpha DIB section.
	hIconOrCursor := CreateIconIndirect(&ii)

	return hIconOrCursor, nil
}
