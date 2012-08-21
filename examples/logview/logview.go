// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	//"fmt"
	//"os"
	//"strings"
	"log"
	"time"
)

import "github.com/ZhuBiCen/walk"

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	myWindow.SetLayout(walk.NewVBoxLayout())
	myWindow.SetTitle("LogView example")

	
	logView, _ := walk.NewLogView(myWindow)
	logView.PostAppendText("XXX")
	log.SetOutput(logView)
	
	go func(){
		for i := 0; i < 10000; i++{
		time.Sleep(10 * time.Millisecond)	
			log.Println("Text" + "\r\n")
		}
	}()
	
	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()
}
