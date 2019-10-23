// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import "math"

var (
	borderGlowAlpha = []float64{0.2, 0.1, 0.075, 0.05, 0.075}

	InteractionEffect WidgetGraphicsEffect
	FocusEffect       WidgetGraphicsEffect
)

type WidgetGraphicsEffect interface {
	Draw(widget Widget, canvas *Canvas) error
}

type widgetGraphicsEffectBase struct {
	color      Color
	dpi2Bitmap map[int]*Bitmap
}

func (wgeb *widgetGraphicsEffectBase) create(color Color) error {
	wgeb.color = color
	return nil
}

func (wgeb *widgetGraphicsEffectBase) Dispose() {
	if len(wgeb.dpi2Bitmap) == 0 {
		return
	}

	for dpi, bitmap := range wgeb.dpi2Bitmap {
		bitmap.Dispose()
		delete(wgeb.dpi2Bitmap, dpi)
	}
}

func (wgeb *widgetGraphicsEffectBase) bitmapForDPI(dpi int) (*Bitmap, error) {
	if wgeb.dpi2Bitmap == nil {
		wgeb.dpi2Bitmap = make(map[int]*Bitmap)
	} else if bitmap, ok := wgeb.dpi2Bitmap[dpi]; ok {
		return bitmap, nil
	}

	var disposables Disposables
	defer disposables.Treat()

	bitmap, err := NewBitmapWithTransparentPixelsForDPI(SizeFrom96DPI(Size{12, 12}, dpi), dpi)
	if err != nil {
		return nil, err
	}
	disposables.Add(bitmap)

	canvas, err := NewCanvasFromImage(bitmap)
	if err != nil {
		return nil, err
	}
	defer canvas.Dispose()

	for i := 1; i <= 5; i++ {
		size := SizeFrom96DPI(Size{i*2 + 2, i*2 + 2}, dpi)

		bmp, err := NewBitmapWithTransparentPixelsForDPI(size, dpi)
		if err != nil {
			return nil, err
		}
		defer bmp.Dispose()

		bmpCanvas, err := NewCanvasFromImage(bmp)
		if err != nil {
			return nil, err
		}
		defer bmpCanvas.Dispose()

		color := RGB(
			byte(math.Min(1.0, float64(wgeb.color.R())/255.0-0.1+0.1*float64(i))*255.0),
			byte(math.Min(1.0, float64(wgeb.color.G())/255.0-0.1+0.1*float64(i))*255.0),
			byte(math.Min(1.0, float64(wgeb.color.B())/255.0-0.1+0.1*float64(i))*255.0),
		)

		brush, err := NewSolidColorBrush(color)
		if err != nil {
			return nil, err
		}
		defer brush.Dispose()

		ellipseSize := SizeFrom96DPI(Size{i * 2, i * 2}, dpi)
		if err := bmpCanvas.FillRoundedRectanglePixels(brush, Rectangle{0, 0, size.Width, size.Height}, ellipseSize); err != nil {
			return nil, err
		}

		bmpCanvas.Dispose()

		opacity := byte(borderGlowAlpha[i-1] * 255.0)

		offset := PointFrom96DPI(Point{5 - i, 5 - i}, dpi)
		canvas.DrawBitmapWithOpacityPixels(bmp, Rectangle{offset.X, offset.Y, size.Width, size.Height}, opacity)
	}

	disposables.Spare()

	wgeb.dpi2Bitmap[dpi] = bitmap

	return bitmap, nil
}

type BorderGlowEffect struct {
	widgetGraphicsEffectBase
}

func NewBorderGlowEffect(color Color) (*BorderGlowEffect, error) {
	bge := new(BorderGlowEffect)

	if err := bge.create(color); err != nil {
		return nil, err
	}

	return bge, nil
}

func (bge *BorderGlowEffect) Draw(widget Widget, canvas *Canvas) error {
	b := widget.BoundsPixels()

	dpi := canvas.DPI()
	bitmap, err := bge.bitmapForDPI(dpi)
	if err != nil {
		return err
	}

	off1 := IntFrom96DPI(1, dpi)
	off2 := IntFrom96DPI(2, dpi)
	off5 := IntFrom96DPI(5, dpi)

	canvas.DrawBitmapPart(bitmap, Rectangle{b.X - off5, b.Y - off5, off5, off5}, Rectangle{0, 0, off5, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X, b.Y - off5, b.Width, off5}, Rectangle{off5 + off1, 0, off1, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + b.Width, b.Y - off5, off5, off5}, Rectangle{off5 + off2, 0, off5, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + b.Width, b.Y, off5, b.Height}, Rectangle{off5 + off2, off5 + off1, off5, off1})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + b.Width, b.Y + b.Height, off5, off5}, Rectangle{off5 + off2, off5 + off2, off5, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X, b.Y + b.Height, b.Width, off5}, Rectangle{off5 + off1, off5 + off2, off1, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X - off5, b.Y + b.Height, off5, off5}, Rectangle{0, off5 + off2, off5, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X - off5, b.Y, off5, b.Height}, Rectangle{0, off5 + off1, off5, off1})

	return nil
}

type DropShadowEffect struct {
	widgetGraphicsEffectBase
}

func NewDropShadowEffect(color Color) (*DropShadowEffect, error) {
	dse := new(DropShadowEffect)

	if err := dse.create(color); err != nil {
		return nil, err
	}

	return dse, nil
}

func (dse *DropShadowEffect) Draw(widget Widget, canvas *Canvas) error {
	b := widget.BoundsPixels()

	dpi := canvas.DPI()
	bitmap, err := dse.bitmapForDPI(dpi)
	if err != nil {
		return err
	}

	off1 := IntFrom96DPI(1, dpi)
	off2 := IntFrom96DPI(2, dpi)
	off5 := IntFrom96DPI(5, dpi)
	off10 := IntFrom96DPI(10, dpi)

	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + b.Width, b.Y + off10 - off5, off5, off5}, Rectangle{off5 + off2, 0, off5, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + b.Width, b.Y + off10, off5, b.Height - off10}, Rectangle{off5 + off2, off5 + off1, off5, off1})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + b.Width, b.Y + b.Height, off5, off5}, Rectangle{off5 + off2, off5 + off2, off5, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + off10, b.Y + b.Height, b.Width - off10, off5}, Rectangle{off5 + off1, off5 + off2, off1, off5})
	canvas.DrawBitmapPart(bitmap, Rectangle{b.X + off10 - off5, b.Y + b.Height, off5, off5}, Rectangle{0, off5 + off2, off5, off5})

	return nil
}

type widgetGraphicsEffectListObserver interface {
	onInsertedGraphicsEffect(index int, effect WidgetGraphicsEffect) error
	onRemovedGraphicsEffect(index int, effect WidgetGraphicsEffect) error
	onClearedGraphicsEffects() error
}

type WidgetGraphicsEffectList struct {
	items    []WidgetGraphicsEffect
	observer widgetGraphicsEffectListObserver
}

func newWidgetGraphicsEffectList(observer widgetGraphicsEffectListObserver) *WidgetGraphicsEffectList {
	return &WidgetGraphicsEffectList{observer: observer}
}

func (l *WidgetGraphicsEffectList) Add(effect WidgetGraphicsEffect) error {
	if effect == nil {
		return newError("effect == nil")
	}

	return l.Insert(len(l.items), effect)
}

func (l *WidgetGraphicsEffectList) At(index int) WidgetGraphicsEffect {
	return l.items[index]
}

func (l *WidgetGraphicsEffectList) Clear() error {
	observer := l.observer
	oldItems := l.items
	l.items = l.items[:0]

	if observer != nil {
		if err := observer.onClearedGraphicsEffects(); err != nil {
			l.items = oldItems
			return err
		}
	}

	return nil
}

func (l *WidgetGraphicsEffectList) Index(effect WidgetGraphicsEffect) int {
	for i, item := range l.items {
		if item == effect {
			return i
		}
	}

	return -1
}

func (l *WidgetGraphicsEffectList) Contains(effect WidgetGraphicsEffect) bool {
	return l.Index(effect) > -1
}

func (l *WidgetGraphicsEffectList) insertIntoSlice(index int, effect WidgetGraphicsEffect) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = effect
}

func (l *WidgetGraphicsEffectList) Insert(index int, effect WidgetGraphicsEffect) error {
	observer := l.observer

	l.insertIntoSlice(index, effect)

	if observer != nil {
		if err := observer.onInsertedGraphicsEffect(index, effect); err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return err
		}
	}

	return nil
}

func (l *WidgetGraphicsEffectList) Len() int {
	return len(l.items)
}

func (l *WidgetGraphicsEffectList) Remove(effect WidgetGraphicsEffect) error {
	index := l.Index(effect)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *WidgetGraphicsEffectList) RemoveAt(index int) error {
	observer := l.observer
	item := l.items[index]

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		if err := observer.onRemovedGraphicsEffect(index, item); err != nil {
			l.insertIntoSlice(index, item)
			return err
		}
	}

	return nil
}
