// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"
	//"os"
	//"strings"
)

import "github.com/lxn/walk"

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	myWindow.SetLayout(walk.NewVBoxLayout())

	splitter, _ := walk.NewSplitter(myWindow)
	splitter.SetOrientation(walk.Vertical)

	btn, _ := walk.NewPushButton(splitter)
	btn.SetText("init")
	btn.Clicked().Attach(func(){
		fmt.Println("Hi")
		myWindow.ProgressIndicator().SetState(walk.PINormal)
		myWindow.ProgressIndicator().SetLength(100)
	})

	btn2, _ := walk.NewPushButton(splitter)
	btn2.SetText("COMEON")
	btn2.Clicked().Attach(func(){
		go func(){
			var i uint
			for i = 0; i < 100; i ++ {
				fmt.Println("Hello")
				time.Sleep(100 * time.Millisecond)
				myWindow.ProgressIndicator().SetValue(i)
			}
		}()
	})
	
	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()

}