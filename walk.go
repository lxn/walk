// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

func LogErrors() bool {
	return logErrors
}

func SetLogErrors(v bool) {
	logErrors = v
}

func PanicOnError() bool {
	return panicOnError
}

func SetPanicOnError(v bool) {
	panicOnError = v
}

func TranslationFunc() TranslationFunction {
	return translation
}

func SetTranslationFunc(f TranslationFunction) {
	translation = f
}

type TranslationFunction func(source string, context ...string) string

var translation TranslationFunction

func tr(source string, context ...string) string {
	if translation == nil {
		return source
	}

	return translation(source, context...)
}
