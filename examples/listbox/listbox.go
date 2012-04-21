package main

import (
	"./walk"
	"os"
	"fmt"
	"strings"
)

func main() {
	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	myWindow.SetLayout(walk.NewVBoxLayout())
	myWindow.SetTitle("Go GUI example")

	// usrGroup, _ := walk.NewGroupBox(myWindow)
	// usrGroup.SetTitle("User Enviroment")
	// sysGroup, _ := walk.NewGroupBox(myWindow)
	// sysGroup.SetTitle("System Enviroment")

	// usrGroup.SetLayout(walk.NewVBoxLayout())
	// myButton1, _ := walk.NewPushButton(usrGroup)
	// myButton1.SetText("XXXX")

	// sysGroup.SetLayout(walk.NewVBoxLayout())
	// myButton2, _ := walk.NewPushButton(sysGroup)
	// myButton2.SetText("XXXXYY")	

	myButton1, _ := walk.NewPushButton(myWindow)
	myButton1.SetText("XXXX")
	// myButton2, _ := walk.NewLabel(myWindow)
	// myButton2.SetText("XXXXYYYYYYY")

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
		lb.AddItem(key)
	}

	lb.SelectedIndexChanged().Attach(func() {
		myButton1.SetText(lb.SelectedItem())
		fmt.Println("SelectedIndex:", lb.SelectedIndex())
		fmt.Println("SelectedItem:",lb.SelectedItem())
	})
	lb.DBClicked().Attach(func() { 
		value, _ := envMap[lb.SelectedItem()]
		walk.MsgBox(myWindow, "About", value, walk.MsgBoxOK|walk.MsgBoxIconInformation)
	})
	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()
}
