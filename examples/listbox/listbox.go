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

type EnvModel struct{
	names  []string
	values []string
	itemsResetPublisher  walk.EventPublisher
	itemChangedPublisher walk.IntEventPublisher	
}

func NewEnvModel() *EnvModel{
	em := &EnvModel{}
	em.names = make([]string, 0)
	em.values = make([]string, 0)
	return em
}

func (em *EnvModel) ItemCount()int{
	return len(em.names)
}

func (em *EnvModel) Value( index int) interface{} {
	return em.names[index]
}

func (em *EnvModel) ItemsReset() *walk.Event {
	return em.itemsResetPublisher.Event()
}

func (em *EnvModel) ItemChanged() *walk.IntEvent {
	return em.itemChangedPublisher.Event()
}


func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	myWindow.SetLayout(walk.NewVBoxLayout())
	myWindow.SetTitle("Go GUI example")

	myButton1, _ := walk.NewPushButton(myWindow)
	myButton1.SetText("XXXX")

	lb, _ := walk.NewListBox(myWindow)

	em := NewEnvModel()

	for _, env := range os.Environ() {
		i := strings.Index(env, "=")
		if i == 0 {
			continue
		}
		key := env[0:i]
		value := env[i+1:]
		em.names = append(em.names, key)
		em.values = append(em.values, value)
	}

	lb.SetModel(em)
	lb.CurrentIndexChanged().Attach(func() {
		myButton1.SetText(lb.CurrentString())
		fmt.Println("CurrentIndex:", lb.CurrentIndex())
		fmt.Println("CurrentString",lb.CurrentString())
	})
	lb.DblClicked().Attach(func() { 
		value := em.values[lb.CurrentIndex()]
		walk.MsgBox(myWindow, "About", value, walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()
}
