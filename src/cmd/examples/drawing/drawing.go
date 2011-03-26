// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"runtime"
)

import (
	"walk"
)

type MainWindow struct {
	*walk.MainWindow
	paintWidget *walk.CustomWidget
}

func createBitmap() *walk.Bitmap {
	bounds := walk.Rectangle{Width: 200, Height: 200}

	bmp, _ := walk.NewBitmap(bounds.Size())

	succeeded := false
	defer func() {
		if !succeeded {
			bmp.Dispose()
		}
	}()

	canvas, _ := walk.NewCanvasFromImage(bmp)
	defer canvas.Dispose()

	brushBmp, _ := walk.NewBitmapFromFile("../img/plus.png")
	defer brushBmp.Dispose()

	brush, _ := walk.NewBitmapBrush(brushBmp)
	defer brush.Dispose()

	canvas.FillRectangle(brush, bounds)

	font, _ := walk.NewFont("Times New Roman", 40, walk.FontBold|walk.FontItalic)
	defer font.Dispose()

	canvas.DrawText("Walk Drawing Example", font, walk.RGB(0, 0, 0), bounds, walk.TextWordbreak)

	succeeded = true
	return bmp
}

func (mw *MainWindow) drawStuff(canvas *walk.Canvas, updateBounds walk.Rectangle) os.Error {
	bmp := createBitmap()
	defer bmp.Dispose()

	bounds := mw.paintWidget.ClientBounds()

	rectPen, _ := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(255, 0, 0))
	defer rectPen.Dispose()

	canvas.DrawRectangle(rectPen, bounds)

	ellipseBrush, _ := walk.NewHatchBrush(walk.RGB(0, 255, 0), walk.HatchCross)
	defer ellipseBrush.Dispose()

	canvas.FillEllipse(ellipseBrush, bounds)

	linesBrush, _ := walk.NewSolidColorBrush(walk.RGB(0, 0, 255))
	defer linesBrush.Dispose()

	linesPen, _ := walk.NewGeometricPen(walk.PenDash, 8, linesBrush)
	defer linesPen.Dispose()

	canvas.DrawLine(linesPen, walk.Point{bounds.X, bounds.Y}, walk.Point{bounds.Width, bounds.Height})
	canvas.DrawLine(linesPen, walk.Point{bounds.X, bounds.Height}, walk.Point{bounds.Width, bounds.Y})

	bmpSize := bmp.Size()
	canvas.DrawImage(bmp, walk.Point{(bounds.Width - bmpSize.Width) / 2, (bounds.Height - bmpSize.Height) / 2})

	return nil
}

func main() {
	runtime.LockOSThread()

	walk.PanicOnError = true

	mainWnd, _ := walk.NewMainWindow()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.SetTitle("Walk Drawing Example")

	mw.ClientArea().SetLayout(walk.NewVBoxLayout())

	mw.paintWidget, _ = walk.NewCustomWidget(mw.ClientArea(), 0, func(canvas *walk.Canvas, updateBounds walk.Rectangle) os.Error {
		return mw.drawStuff(canvas, updateBounds)
	})
	mw.paintWidget.SetClearsBackground(true)
	mw.paintWidget.SetInvalidatesOnResize(true)

	mw.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	mw.SetSize(walk.Size{800, 600})
	mw.Show()

	os.Exit(mw.Run())
}
