// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"unsafe"
)

import . "github.com/lxn/go-winapi"


type ProgressIndicator struct {
	hwnd HWND
	taskbarList3 *ITaskbarList3
	length uint64
}

type PIState int
const (
    PINoProgress	PIState = TBPF_NOPROGRESS
	PIIndeterminate PIState = TBPF_INDETERMINATE
	PINormal        PIState = TBPF_NORMAL
	PIError         PIState = TBPF_ERROR
	PIPaused        PIState = TBPF_PAUSED
)


func newTaskbarList3(hwnd HWND)(*ProgressIndicator, error) {

	// Check that the Windows version is at least 6.1 (yes, Win 7 is version 6.1).
	var classFactoryPtr unsafe.Pointer
	if hr := CoGetClassObject(&CLSID_TaskbarList, CLSCTX_ALL, nil, &IID_IClassFactory, &classFactoryPtr); FAILED(hr) {
		return nil, errorFromHRESULT("CoGetClassObject", hr)
	}

	var taskbarList3ObjectPtr unsafe.Pointer
	classFactory := (*IClassFactory)(classFactoryPtr)
	defer classFactory.Release()

	if hr := classFactory.CreateInstance(nil, &IID_ITaskbarList3, &taskbarList3ObjectPtr); FAILED(hr) {
		return nil, errorFromHRESULT("IClassFactory.CreateInstance", hr)
	}

	return &ProgressIndicator{taskbarList3:(*ITaskbarList3)(taskbarList3ObjectPtr), hwnd:hwnd},  nil
}

func(pi *ProgressIndicator) SetState(state PIState) error {
	if hr := pi.taskbarList3.SetProgressState(pi.hwnd, (int)(state)); FAILED(hr){
		return errorFromHRESULT("ITaskbarList3.setprogressState", hr)
	}
	return nil
}

func (pi* ProgressIndicator) SetLength(length uint64) {
	pi.length = length
}

func(pi *ProgressIndicator) SetValue(pos uint64) error {
	if hr := pi.taskbarList3.SetProgressValue(pi.hwnd, pos, pi.length); FAILED(hr){
		return errorFromHRESULT("ITaskbarList3.SetProgressValue", hr)
	}
	return nil
}