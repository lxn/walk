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

func processError(err os.Error) os.Error {
	if logErrors {
		log.Print(err)
	}

	if panicOnError {
		panic(err)
	}

	return err
}

func newError(message string) os.Error {
	return processError(&Error{message: message, stack: debug.Stack()})
}

func lastError(win32FuncName string) os.Error {
	if errno := GetLastError(); errno != ERROR_SUCCESS {
		return newError(fmt.Sprintf("%s: %s", win32FuncName, syscall.Errstr(int(errno))))
	}

	return nil
}

func errorFromHRESULT(funcName string, hr HRESULT) os.Error {
	return newError(fmt.Sprintf("%s: %s", funcName, syscall.Errstr(int(hr))))
}

func wrapError(err os.Error) os.Error {
	return processError(&Error{inner: err, stack: debug.Stack()})
}

func toError(x interface{}) os.Error {
	switch x := x.(type) {
	case os.Error:
		return wrapError(x)

	case string:
		return newError(x)
	}

	return newError(fmt.Sprintf("Error: %v", x))
}
