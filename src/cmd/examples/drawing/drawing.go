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

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func createBitmap() *walk.Bitmap {
	bounds := walk.Rectangle{Width: 200, Height: 200}

	bmp, err := walk.NewBitmap(bounds.Size())
	panicIfErr(err)

	succeeded := false
	defer func() {
		if !succeeded {
			bmp.Dispose()
		}
	}()

	canvas, err := walk.NewCanvasFromImage(bmp)
	panicIfErr(err)
	defer canvas.Dispose()

	brushBmp, err := walk.NewBitmapFromFile("../img/plus.png")
	panicIfErr(err)
	defer brushBmp.Dispose()

	brush, err := walk.NewBitmapBrush(brushBmp)
	panicIfErr(err)
	defer brush.Dispose()

	panicIfErr(canvas.FillRectangle(brush, bounds))

	font, err := walk.NewFont("Times New Roman", 40, walk.FontBold|walk.FontItalic)
	panicIfErr(err)
	defer font.Dispose()

	panicIfErr(canvas.DrawText("Walk Drawing Example", font, walk.RGB(0, 0, 0), bounds, walk.TextWordbreak))

	succeeded = true
	return bmp
}

func (mw *MainWindow) drawStuff(canvas *walk.Canvas, updateBounds walk.Rectangle) os.Error {
	bmp := createBitmap()
	defer bmp.Dispose()

	bounds := mw.paintWidget.ClientBounds()

	rectPen, err := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(255, 0, 0))
	panicIfErr(err)
	defer rectPen.Dispose()

	panicIfErr(canvas.DrawRectangle(rectPen, bounds))

	ellipseBrush, err := walk.NewHatchBrush(walk.RGB(0, 255, 0), walk.HatchCross)
	panicIfErr(err)
	defer ellipseBrush.Dispose()

	panicIfErr(canvas.FillEllipse(ellipseBrush, bounds))

	linesBrush, err := walk.NewSolidColorBrush(walk.RGB(0, 0, 255))
	panicIfErr(err)
	defer linesBrush.Dispose()

	linesPen, err := walk.NewGeometricPen(walk.PenDash, 8, linesBrush)
	panicIfErr(err)
	defer linesPen.Dispose()

	panicIfErr(canvas.DrawLine(linesPen, walk.Point{bounds.X, bounds.Y}, walk.Point{bounds.Width, bounds.Height}))
	panicIfErr(canvas.DrawLine(linesPen, walk.Point{bounds.X, bounds.Height}, walk.Point{bounds.Width, bounds.Y}))

	bmpSize := bmp.Size()
	panicIfErr(canvas.DrawImage(bmp, walk.Point{(bounds.Width - bmpSize.Width) / 2, (bounds.Height - bmpSize.Height) / 2}))

	return nil
}

func main() {
	runtime.LockOSThread()

	mainWnd, err := walk.NewMainWindow()
	panicIfErr(err)

	mw := &MainWindow{MainWindow: mainWnd}
	panicIfErr(mw.SetTitle("Walk Drawing Example"))

	panicIfErr(mw.ClientArea().SetLayout(walk.NewVBoxLayout()))

	mw.paintWidget, err = walk.NewCustomWidget(mw.ClientArea(), 0, func(canvas *walk.Canvas, updateBounds walk.Rectangle) os.Error {
		return mw.drawStuff(canvas, updateBounds)
	})
	panicIfErr(err)
	mw.paintWidget.SetClearsBackground(true)
	mw.paintWidget.SetInvalidatesOnResize(true)

	panicIfErr(mw.SetMinMaxSize(walk.Size{320, 240}, walk.Size{}))
	panicIfErr(mw.SetSize(walk.Size{800, 600}))
	mw.Show()

	os.Exit(mw.Run())
}
