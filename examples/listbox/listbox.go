// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

import "github.com/lxn/walk"

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	myWindow.SetLayout(walk.NewVBoxLayout())
	myWindow.SetTitle("Go GUI example")

	myButton1, _ := walk.NewPushButton(myWindow)
	myButton1.SetText("XXXX")

	envMap := make(map[string]string)

	lb, _ := walk.NewListBox(myWindow)
	for _, env := range os.Environ() {
		i := strings.Index(env, "=")
		if i == 0 {
			continue
		}
		key := env[0:i]
		value := env[i+1:]
		envMap[key] = value
		lb.AddString(key)
	}

	lb.CurrentIndexChanged().Attach(func() {
		myButton1.SetText(lb.CurrentString())
		fmt.Println("CurrentIndex:", lb.CurrentIndex())
		fmt.Println("CurrentString",lb.CurrentString())
	})
	lb.DblClicked().Attach(func() { 
		value, _ := envMap[lb.CurrentString()]
		walk.MsgBox(myWindow, "About", value, walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()
}
