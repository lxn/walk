// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type InitParams struct {
	LogErrors    bool
	PanicOnError bool
	Translation  func(source string, context ...string) string
}

func Initialize(params InitParams) {
	logErrors = params.LogErrors
	panicOnError = params.PanicOnError
	translation = params.Translation
}

func Shutdown() {
}

var translation func(source string, context ...string) string

func tr(source string, context ...string) string {
	if translation == nil {
		return source
	}

	return translation(source, context...)
}
