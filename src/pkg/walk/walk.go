// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"runtime"
)

import (
	"walk/winapi/gdiplus"
	"walk/winapi/ole32"
)

type InitParams struct {
	LogErrors    bool
	PanicOnError bool
}

func Initialize(params InitParams) {
	runtime.LockOSThread()

	logErrors = params.LogErrors
	panicOnError = params.PanicOnError

	// TODO: Should we setup winapi syscalls from here instead using init funcs?
}

func Shutdown() {
	gdiplus.GdiplusShutdown()
	ole32.OleUninitialize()
}
