// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

const progressIndicatorClass = `\o/ Walk_ProgressIndicator_Class \0/`

type ProgressIndicator struct {
}

func NewProgressIndicator()(*ProgressIndicator, error) {

	// Check that the Windows version is at least 6.1 (yes, Win 7 is version 6.1).
	var classFactoryPtr unsafe.Pointer
	if hr := CoGetClassObject(&CLSID_TaskbarList, CLSCTX_ALL, nil, &IID_IClassFactory, &classFactoryPtr); FAILED(hr) {
		return nil, errorFromHRESULT("CoGetClassObject", hr)
	}

	var taskbarList3ObjectPtr unsafe.Pointer

	if hr := classFactory.CreateInstance(nil, &IID_ITaskbarList3, &taskbarList3ObjectPtr); FAILED(hr) {
		return nil, errorFromHRESULT("IClassFactory.CreateInstance", hr)
	}

	taskbarList3Object := (*ITaskbarList3Object)(taskbarList3ObjectPtr)
	

}
