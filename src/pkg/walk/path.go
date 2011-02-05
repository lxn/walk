// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

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

func DriveNames() ([]string, os.Error) {
	bufLen := GetLogicalDriveStrings(0, nil)
	if bufLen == 0 {
		return nil, lastError("GetLogicalDriveStrings")
	}
	buf := make([]uint16, bufLen+1)

	bufLen = GetLogicalDriveStrings(bufLen+1, &buf[0])
	if bufLen == 0 {
		return nil, lastError("GetLogicalDriveStrings")
	}

	var names []string

	for i := 0; i < len(buf)-2; {
		name := syscall.UTF16ToString(buf[i:])
		names = append(names, name)
		i += len(name) + 1
	}

	return names, nil
}
