// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"path/filepath"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

// Icon is a bitmap that supports transparency and combining multiple 
// variants of an image in different resolutions.
type Icon struct {
	hIcon HICON
}

// NewIconFromFile returns a new Icon, using the specified icon image file.
func NewIconFromFile(filePath string) (*Icon, os.Error) {
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

// Dispose releases the operating system resources associated with the Icon.
func (i *Icon) Dispose() os.Error {
	if i.hIcon == 0 {
		return nil
	}

	if !DestroyIcon(i.hIcon) {
		return lastError("DestroyIcon")
	}

	i.hIcon = 0

	return nil
}
