// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"unsafe"
)

import (
	. "github.com/lxn/go-winapi"
)

type ProgressIndicator struct {
	hwnd         HWND
	taskbarList3 *ITaskbarList3
	completed    uint32
	total        uint32
	state        PIState
}

type PIState int

const (
	PINoProgress    PIState = TBPF_NOPROGRESS
	PIIndeterminate PIState = TBPF_INDETERMINATE
	PINormal        PIState = TBPF_NORMAL
	PIError         PIState = TBPF_ERROR
	PIPaused        PIState = TBPF_PAUSED
)

//newTaskbarList3 precondition: Windows version is at least 6.1 (yes, Win 7 is version 6.1).
func newTaskbarList3(hwnd HWND) (*ProgressIndicator, error) {
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

	return &ProgressIndicator{taskbarList3: (*ITaskbarList3)(taskbarList3ObjectPtr), hwnd: hwnd}, nil
}

func (pi *ProgressIndicator) SetState(state PIState) error {
	if hr := pi.taskbarList3.SetProgressState(pi.hwnd, (int)(state)); FAILED(hr) {
		return errorFromHRESULT("ITaskbarList3.setprogressState", hr)
	}
	pi.state = state
	return nil
}

func (pi *ProgressIndicator) State() PIState {
	return pi.state
}

func (pi *ProgressIndicator) SetTotal(total uint32) {
	pi.total = total
}

func (pi *ProgressIndicator) Total() uint32 {
	return pi.total
}

func (pi *ProgressIndicator) SetCompleted(completed uint32) error {
	if hr := pi.taskbarList3.SetProgressValue(pi.hwnd, completed, pi.total); FAILED(hr) {
		return errorFromHRESULT("ITaskbarList3.SetProgressValue", hr)
	}
	pi.completed = completed
	return nil
}

func (pi *ProgressIndicator) Completed() uint32 {
	return pi.completed
}
