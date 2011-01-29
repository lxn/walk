// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package registry

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi/advapi32"
	. "walk/winapi/kernel32"
)

type Key struct {
	hKey HKEY
}

func ClassesRootKey() *Key {
	return &Key{HKEY_CLASSES_ROOT}
}

func CurrentUserKey() *Key {
	return &Key{HKEY_CURRENT_USER}
}

func LocalMachineKey() *Key {
	return &Key{HKEY_LOCAL_MACHINE}
}

func KeyString(rootKey *Key, subKeyPath, valueName string) (value string, err os.Error) {
	var hKey HKEY
	if RegOpenKeyEx(rootKey.hKey, syscall.StringToUTF16Ptr(subKeyPath), 0, KEY_READ, &hKey) != ERROR_SUCCESS {
		return "", newError("KeyString: Failed to open subkey.")
	}
	defer RegCloseKey(hKey)

	var typ uint
	var data []uint16
	var bufSize uint

	if RegQueryValueEx(hKey, syscall.StringToUTF16Ptr(valueName), (*uint)(unsafe.Pointer(nil)), &typ, nil, &bufSize) != ERROR_SUCCESS {
		return "", newError("KeyString: Failed to retrieve required buffer size.")
	}

	data = make([]uint16, bufSize/2+1)

	if RegQueryValueEx(hKey, syscall.StringToUTF16Ptr(valueName), (*uint)(unsafe.Pointer(nil)), &typ, (*byte)(unsafe.Pointer((&data[0]))), &bufSize) != ERROR_SUCCESS {
		return "", newError("KeyString: Failed to retrieve registry key value.")
	}

	return syscall.UTF16ToString(data), nil
}
