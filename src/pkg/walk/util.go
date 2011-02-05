// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

import (
	. "walk/winapi"
	. "walk/winapi/kernel32"
)

func callStack() string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("=======================================================\n")

	i := 0
	for {
		pc, file, line, ok := runtime.Caller(i + 1)
		if !ok {
			break
		}
		if i > 0 {
			buf.WriteString("-------------------------------------------------------\n")
		}

		fun := runtime.FuncForPC(pc)
		name := fun.Name()

		buf.WriteString(fmt.Sprintf("%s (%s, Line %d)\n", name, file, line))

		i++
	}

	buf.WriteString("=======================================================\n")

	return buf.String()
}

func printCallStack() {
	fmt.Print(callStack())
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func toError(x interface{}) os.Error {
	switch x := x.(type) {
	case os.Error:
		return x

	case string:
		return newError(x)
	}

	return newError(fmt.Sprintf("Error: %v", x))
}

func newError(message string) os.Error {
	return os.NewError(fmt.Sprintf("%s\nCall Stack:\n%s", message, callStack()))
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

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

/*func boolToBOOL(value bool) BOOL {
	if value {
		return 1
	}

	return 0
}*/
