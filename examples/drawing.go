// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

import (
	"walk/drawing"
	"walk/gui"
)

type MainWindow struct {
	*gui.MainWindow
	treeView *gui.TreeView
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func (mw *MainWindow) drawStuff() {
	bounds, err := mw.treeView.ClientBounds()
	panicIfErr(err)

	surface, err := mw.treeView.GetDrawingSurface()
	panicIfErr(err)
	defer surface.Dispose()

	rectPen, err := drawing.NewCosmeticPen(drawing.PenSolid, drawing.RGB(255, 0, 0))
	panicIfErr(err)
	defer rectPen.Dispose()

	panicIfErr(surface.DrawRectangle(rectPen, bounds))

	font := drawing.NewFont()
	font.BeginEdit()
	font.SetFamily("Tahoma")
	font.SetPointSize(36)
	font.SetBold(true)
	panicIfErr(font.EndEdit())
	defer font.Dispose()

	text := strings.Repeat("Hello! ", 10)
	panicIfErr(surface.DrawText(text, font, drawing.RGB(255, 192, 128), bounds, drawing.TextWordbreak))

	ellipseBrush, err := drawing.NewHatchBrush(drawing.RGB(0, 255, 0), drawing.HatchCross)
	panicIfErr(err)
	defer ellipseBrush.Dispose()

	panicIfErr(surface.FillEllipse(ellipseBrush, bounds))

	linesBrush, err := drawing.NewSolidColorBrush(drawing.RGB(0, 0, 255))
	panicIfErr(err)
	defer linesBrush.Dispose()

	linesPen, err := drawing.NewGeometricPen(drawing.PenDash, 8, linesBrush)
	panicIfErr(err)
	defer linesPen.Dispose()

	panicIfErr(surface.DrawLine(linesPen, drawing.Point{bounds.X, bounds.Y}, drawing.Point{bounds.Width, bounds.Height}))
	panicIfErr(surface.DrawLine(linesPen, drawing.Point{bounds.X, bounds.Height}, drawing.Point{bounds.Width, bounds.Y}))
}

func runMainWindow() {
	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)
	defer mainWnd.Dispose()

	mw := &MainWindow{MainWindow: mainWnd}

	mw.ClientArea().SetLayout(gui.NewVBoxLayout())

	drawButton, err := gui.NewPushButton(mw.ClientArea())
	panicIfErr(err)
	panicIfErr(drawButton.SetText("Draw Stuff!"))
	drawButton.AddClickedHandler(func(args gui.EventArgs) { mw.drawStuff() })

	mw.treeView, err = gui.NewTreeView(mw.ClientArea())
	panicIfErr(err)

	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	panicIfErr(mw.RunMessageLoop())
}

func main() {
	runtime.LockOSThread()

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Error:", x)
		}
	}()

	runMainWindow()
}
