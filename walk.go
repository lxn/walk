// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
)

import "github.com/lxn/go-winapi"

type InitParams struct {
	LogErrors    bool
	PanicOnError bool
	Translation  func(source string, context ...string) string
}

func Initialize(params InitParams) {
	logErrors = params.LogErrors
	panicOnError = params.PanicOnError
	translation = params.Translation

	if hr := winapi.OleInitialize(); winapi.FAILED(hr) {
		panic(fmt.Sprint("OleInitialize Error: ", hr))
	}

	// TODO: Should we setup winapi syscalls from here instead using init funcs?
}

func Shutdown() {
	winapi.OleUninitialize()
}

var translation func(source string, context ...string) string

func tr(source string, context ...string) string {
	if translation == nil {
		return source
	}

	return translation(source, context...)
}
