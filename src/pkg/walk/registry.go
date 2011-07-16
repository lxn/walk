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

type RegistryKey struct {
	hKey HKEY
}

func ClassesRootKey() *RegistryKey {
	return &RegistryKey{HKEY_CLASSES_ROOT}
}

func CurrentUserKey() *RegistryKey {
	return &RegistryKey{HKEY_CURRENT_USER}
}

func LocalMachineKey() *RegistryKey {
	return &RegistryKey{HKEY_LOCAL_MACHINE}
}

func RegistryKeyString(rootKey *RegistryKey, subKeyPath, valueName string) (value string, err os.Error) {
	var hKey HKEY
	if RegOpenKeyEx(
		rootKey.hKey,
		syscall.StringToUTF16Ptr(subKeyPath),
		0,
		KEY_READ,
		&hKey) != ERROR_SUCCESS {

		return "", newError("RegistryKeyString: Failed to open subkey.")
	}
	defer RegCloseKey(hKey)

	var typ uint32
	var data []uint16
	var bufSize uint32

	if ERROR_SUCCESS != RegQueryValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		&typ,
		nil,
		&bufSize) {

		return "", newError("RegQueryValueEx #1")
	}

	data = make([]uint16, bufSize/2+1)

	if ERROR_SUCCESS != RegQueryValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&data[0])),
		&bufSize) {

		return "", newError("RegQueryValueEx #2")
	}

	return syscall.UTF16ToString(data), nil
}
