// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"os"
	"syscall"
)

import (
	. "walk/winapi/kernel32"
	. "walk/winapi/shell32"
)

func knownFolderPath(id CSIDL) (string, os.Error) {
	var buf [MAX_PATH]uint16

	if !ShGetSpecialFolderPath(0, &buf[0], id, false) {
		return "", newError("ShGetSpecialFolderPath failed")
	}

	return syscall.UTF16ToString(buf[0:]), nil
}

func AppData() (string, os.Error) {
	return knownFolderPath(CSIDL_APPDATA)
}

func CommonAppData() (string, os.Error) {
	return knownFolderPath(CSIDL_COMMON_APPDATA)
}

func LocalAppData() (string, os.Error) {
	return knownFolderPath(CSIDL_LOCAL_APPDATA)
}
