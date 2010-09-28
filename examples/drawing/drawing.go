// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"runtime"
)

import (
	"walk/drawing"
	"walk/gui"
)

type MainWindow struct {
	*gui.MainWindow
	paintWidget *gui.CustomWidget
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func createBitmap() *drawing.Bitmap {
	bounds := drawing.Rectangle{Width: 200, Height: 200}

	bmp, err := drawing.NewBitmap(bounds.Size())
	panicIfErr(err)

	surface, err := drawing.NewSurfaceFromBitmap(bmp)
	panicIfErr(err)
	defer surface.Dispose()

	brushBmp, err := drawing.NewBitmapFromFile("img/plus.png")
	panicIfErr(err)
	defer brushBmp.Dispose()

	brush, err := drawing.NewBitmapBrush(brushBmp)
	panicIfErr(err)
	defer brush.Dispose()

	panicIfErr(surface.FillRectangle(brush, bounds))

	font, err := drawing.NewFont("Times New Roman", 48, drawing.FontBold|drawing.FontItalic)
	panicIfErr(err)
	defer font.Dispose()

	panicIfErr(surface.DrawText("Runtime Created Bitmap", font, drawing.RGB(0, 0, 0), bounds, drawing.TextWordbreak))

	return bmp
}

func (mw *MainWindow) drawStuff(surface *drawing.Surface, updateBounds drawing.Rectangle) os.Error {
	bmp := createBitmap()
	defer bmp.Dispose()

	bounds, err := mw.paintWidget.ClientBounds()
	panicIfErr(err)

	rectPen, err := drawing.NewCosmeticPen(drawing.PenSolid, drawing.RGB(255, 0, 0))
	panicIfErr(err)
	defer rectPen.Dispose()

	panicIfErr(surface.DrawRectangle(rectPen, bounds))

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

	bmpSize := bmp.Size()
	panicIfErr(surface.DrawImage(bmp, drawing.Point{bounds.Width/2 - bmpSize.Width/2, bounds.Height/2 - bmpSize.Height/2}))

	return nil
}

func runMainWindow() {
	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)
	defer mainWnd.Dispose()

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetText("Walk Drawing Example"))

	mw.ClientArea().SetLayout(gui.NewVBoxLayout())

	mw.paintWidget, err = gui.NewCustomWidget(mw.ClientArea(), 0, func(surface *drawing.Surface, updateBounds drawing.Rectangle) os.Error {
		return mw.drawStuff(surface, updateBounds)
	})
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
