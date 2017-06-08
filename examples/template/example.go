package main

import (
	"time"
	"walk"
	. "walk/template"
)

func console_font() *walk.Font {
	for _, face := range []string{"Consolas", "Andale Mono", "Lucida Console", "Courier New"} {
		f, err := walk.NewFont(face, 10, walk.FontStyle(0))
		if err == nil {
			return f
		}
	}
	return nil
}

func main() {
	walk.Initialize(walk.InitParams{LogErrors:true, PanicOnError:true})
	defer walk.Shutdown()

	var (
		editFile, editParameter, editStatus *walk.LineEdit
		buttonFile, buttonChk, buttonFix, buttonClr *walk.Button
		progressBar *walk.ProgressBar
		output *walk.TextEdit
	)
	frame, err := (MainWindow{"Test tool",
		Spacing(6),
		VBoxComposite{
			Options: Margins(6,6,6,6),
			Children: []Widget{
				HBoxCompositeChildren{
					Label("File"),
					LineEdit(&editFile, "", nil),
					ToolButton(&buttonFile, "...", nil, nil),
				},
				HBoxCompositeChildren{
					PushButton(&buttonChk, "Check", nil, nil),
					PushButton(&buttonFix, "Check and Fix", nil, nil),
					PushButton(&buttonClr, "Clear", nil, nil),
					HSpacer(StretchFaktor(10)),
					Label("Parameter"),
					LineEdit(&editParameter, "", MaxLength(10)),
				},
				HBoxCompositeChildren{
					LineEdit(&editStatus, "Ready.", ReadOnly),
					ProgressBar(&progressBar, 0, 100, StretchFaktor(10)),
				},
				TextEdit(&output, "", WidgetFlags(Font(console_font()), ReadOnly)),
			},
		},
	}).Create(nil)

	if err != nil {
		panic(err)
	}

	procfn := func(ops ...string) {
		buttonChk.SetEnabled(false)
		buttonFix.SetEnabled(false)
		buttonClr.SetEnabled(false)
		editStatus.SetText("")
		go func() {
			for _, op := range ops {
				editStatus.SetText(op + "...")
				for i := 0; i < 100; i++ {
					progressBar.SetValue(i)
					time.Sleep(time.Duration(20) * time.Millisecond)
				}
				output.AppendText(op + " done on " + time.Now().String() + "\r\n")
			}
			editStatus.SetText("Finished.")
			buttonChk.SetEnabled(true)
			buttonFix.SetEnabled(true)
			buttonClr.SetEnabled(true)
		}()
	}

	buttonChk.Clicked().Attach(func() {
		procfn("Checking")
	})

	buttonFix.Clicked().Attach(func() {
		procfn("Checking", "Fixing")
	})

	buttonClr.Clicked().Attach(func() {
		output.SetText("")
	})

	frame.Show()
	frame.Run()
}

