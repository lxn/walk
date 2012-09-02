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
type MyDialog struct {
    *walk.Dialog
    ui myDialogUI
}

func RunMyDialog(owner walk.RootWidget) (int, error) {
    dlg := new(MyDialog)
    if err := dlg.init(owner); err != nil {
        return 0, err
    }

    //TODO: Do further required setup, e.g. for event handling, here.

	dlg.ui.noProgressBtn.Clicked().Attach(func(){
		fmt.Println("Hi")
		dlg.ProgressIndicator().SetState(walk.PINoProgress)
	})

	dlg.ui.normalBtn.Clicked().Attach(func(){
		fmt.Println("Hi")
		dlg.ProgressIndicator().SetState(walk.PINormal)
	})

	dlg.ui.errBtn.Clicked().Attach(func(){
		fmt.Println("Hi")
		dlg.ProgressIndicator().SetState(walk.PIError)
	})

	dlg.ui.pausedBtn.Clicked().Attach(func(){
		fmt.Println("Hi")
		dlg.ProgressIndicator().SetState(walk.PIPaused)
	})

	dlg.ui.startBtn.Clicked().Attach(func(){
		go func(){
		dlg.ProgressIndicator().SetLength(100)
			var i uint32
			for i = 0; i < 100; i ++ {
				fmt.Println("Hello", i)
				time.Sleep(100 * time.Millisecond)
				if err := dlg.ProgressIndicator().SetValue(i); err != nil {
					fmt.Println(err)
				}
			}
		}()
	})

    return dlg.Run(), nil
}
func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	RunMyDialog(nil)
}