// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"syscall"
)

import (
	. "walk/winapi"
	. "walk/winapi/kernel32"
)

var (
	logErrors    bool
	panicOnError bool
)

type Error struct {
	inner   os.Error
	message string
	stack   []byte
}

func (err *Error) Inner() os.Error {
	return err.inner
}

func (err *Error) Message() string {
	if err.message != "" {
		return err.message
	}

	if err.inner != nil {
		if walkErr, ok := err.inner.(*Error); ok {
			return walkErr.Message()
		} else {
			return err.inner.String()
		}
	}

	return ""
}

func (err *Error) Stack() []byte {
	return err.stack
}

func (err *Error) String() string {
	return fmt.Sprintf("%s\n\nStack:\n%s", err.Message(), err.stack)
}

func processErrorNoPanic(err os.Error) os.Error {
	if logErrors {
		if walkErr, ok := err.(*Error); ok {
			log.Print(walkErr.String())
		} else {
			log.Printf("%s\n\nStack:\n%s", err, debug.Stack())
		}
	}

	return err
}

func processError(err os.Error) os.Error {
	processErrorNoPanic(err)

	if panicOnError {
		panic(err)
	}

	return err
}

func newErr(message string) os.Error {
	return &Error{message: message, stack: debug.Stack()}
}

func newError(message string) os.Error {
	return processError(newErr(message))
}

func newErrorNoPanic(message string) os.Error {
	return processErrorNoPanic(newErr(message))
}

func lastError(win32FuncName string) os.Error {
	if errno := GetLastError(); errno != ERROR_SUCCESS {
		return newError(fmt.Sprintf("%s: %s", win32FuncName, syscall.Errstr(int(errno))))
	}

	return newError(win32FuncName)
}

func errorFromHRESULT(funcName string, hr HRESULT) os.Error {
	return newError(fmt.Sprintf("%s: %s", funcName, syscall.Errstr(int(hr))))
}

func wrapErr(err os.Error) os.Error {
	if _, ok := err.(*Error); ok {
		return err
	}

	return &Error{inner: err, stack: debug.Stack()}
}

func wrapErrorNoPanic(err os.Error) os.Error {
	return processErrorNoPanic(wrapErr(err))
}

func wrapError(err os.Error) os.Error {
	return processError(wrapErr(err))
}

func toErrorNoPanic(x interface{}) os.Error {
	switch x := x.(type) {
	case *Error:
		return x

	case os.Error:
		return wrapErrorNoPanic(x)

	case string:
		return newErrorNoPanic(x)
	}

	return newErrorNoPanic(fmt.Sprintf("Error: %v", x))
}

func toError(x interface{}) os.Error {
	err := toErrorNoPanic(x)

	if panicOnError {
		panic(err)
	}

	return err
}
