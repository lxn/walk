// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi/comdlg32"
)

type FileDialog struct {
	Title          string
	FilePath       string
	InitialDirPath string
	Filter         string
	FilterIndex    int
}

func (dlg *FileDialog) show(owner RootWidget, fun func(ofn *OPENFILENAME) bool) (accepted bool, err os.Error) {
	ofn := &OPENFILENAME{}

	ofn.LStructSize = uint(unsafe.Sizeof(*ofn))
	ofn.HwndOwner = owner.Handle()

	filter := make([]uint16, len(dlg.Filter)+1)
	copy(filter, syscall.StringToUTF16(dlg.Filter))
	// Replace '|' with the expected '\0'.
	for i, c := range filter {
		if byte(c) == '|' {
			filter[i] = uint16(0)
		}
	}
	ofn.LpstrFilter = &filter[0]
	ofn.NFilterIndex = uint(dlg.FilterIndex)

	filePath := make([]uint16, 1024)
	copy(filePath, syscall.StringToUTF16(dlg.FilePath))
	ofn.LpstrFile = &filePath[0]
	ofn.NMaxFile = uint(len(filePath))

	ofn.LpstrInitialDir = syscall.StringToUTF16Ptr(dlg.InitialDirPath)
	ofn.LpstrTitle = syscall.StringToUTF16Ptr(dlg.Title)
	ofn.Flags = OFN_FILEMUSTEXIST

	if !fun(ofn) {
		errno := CommDlgExtendedError()
		if errno != 0 {
			err = newError(fmt.Sprintf("Error %d", errno))
		}
		return
	}

	dlg.FilePath = syscall.UTF16ToString(filePath)

	accepted = true

	return
}

func (dlg *FileDialog) ShowOpen(owner RootWidget) (accepted bool, err os.Error) {
	return dlg.show(owner, GetOpenFileName)
}

func (dlg *FileDialog) ShowSave(owner RootWidget) (accepted bool, err os.Error) {
	return dlg.show(owner, GetSaveFileName)
}
