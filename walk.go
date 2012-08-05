// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
	"runtime"
)

import "github.com/lxn/go-winapi"

type InitParams struct {
	LogErrors    bool
	PanicOnError bool
}

func Initialize(params InitParams) {
	runtime.LockOSThread()

	logErrors = params.LogErrors
	panicOnError = params.PanicOnError

	if hr := winapi.OleInitialize(); winapi.FAILED(hr) {
		panic(fmt.Sprint("OleInitialize Error: ", hr))
	}

	// TODO: Should we setup winapi syscalls from here instead using init funcs?
}

func Shutdown() {
	winapi.GdiplusShutdown()
	winapi.OleUninitialize()
}
