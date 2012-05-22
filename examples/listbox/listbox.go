// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

import "github.com/lxn/walk"

type EnvItem struct {
	varName string
	value   string
}

type EnvModel struct{
	walk.ListModelBase
	envItems []EnvItem
	itemsResetPublisher  walk.EventPublisher
	itemChangedPublisher walk.IntEventPublisher	
}

func NewEnvModel() *EnvModel{
	em := &EnvModel{}
	em.envItems = make([]EnvItem, 0)
	return em
}

func (em *EnvModel) ItemCount()int{
	return len(em.envItems)
}

func (em *EnvModel) Value( index int) interface{} {
	return em.envItems[index].varName
}



func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	myWindow.SetLayout(walk.NewVBoxLayout())
	myWindow.SetTitle("Go GUI example")

	
	splitter, _ := walk.NewSplitter(myWindow)
	splitter.SetOrientation(walk.Vertical)
	//splitter.SetHeight(270)

	lb, _ := walk.NewListBox(splitter)
	//lb.SetHeight(100)
	//lb.SetMinMaxSize(walk.Size{Width:100, Height:100}, walk.Size{})
	//lb.SetSize(walk.Size{Width:100, Height:100})


	valueEdit, _ := walk.NewTextEdit(splitter)
	valueEdit.SetReadOnly(true)
	//valueEdit.SetHeight(300)
	

	buttonCompositor, _ := walk.NewComposite(myWindow)
	hbox := walk.NewHBoxLayout()
	buttonCompositor.SetLayout(hbox)
	buttonCompositor.SetHeight(30)
	
	myButton1, _ := walk.NewPushButton(buttonCompositor)
	myButton1.SetText("New")

	myButton2, _ := walk.NewPushButton(buttonCompositor)
	myButton2.SetText("Edit")

	myButton3, _ := walk.NewPushButton(buttonCompositor)
	myButton3.SetText("Delete")


	//env model
	em := NewEnvModel()

	for _, env := range os.Environ() {
		i := strings.Index(env, "=")
		if i == 0 {
			continue
		}
		varName := env[0:i]
		value := env[i+1:]
		envItem := EnvItem{varName, value}
		
		em.envItems = append(em.envItems, envItem)
	}

	fmt.Println("The len of Model", em.ItemCount())
	lb.SetModel(em)
	lb.CurrentIndexChanged().Attach(func() {
		if curVar, ok := em.Value(lb.CurrentIndex()).(string); ok {
			value := em.envItems[lb.CurrentIndex()].value
			value = strings.Replace(value, ";", "\r\n", -1)
			valueEdit.SetText(value)
			fmt.Println("CurrentIndex:", lb.CurrentIndex())
			fmt.Println("CurrentEnvVarName:",curVar)
		}
	})
	lb.DblClicked().Attach(func() { 
		value := em.envItems[lb.CurrentIndex()].value
		value = strings.Replace(value, ";", "\r\n", -1)
		valueEdit.SetText(value)
		walk.MsgBox(myWindow, "About", value, walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()
}
