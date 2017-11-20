// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"math"
)

var (
	borderGlowAlpha = []float64{0.4, 0.2, 0.15, 0.1, 0.15}

	DefaultDropShadowEffect, _ = NewDropShadowEffect(RGB(63, 63, 63), false)
)

type WidgetGraphicsEffect interface {
	Draw(widget Widget, canvas *Canvas) error
	Enabled() bool
	SetEnabled(enabled bool)
	EnabledChanged() *Event
}

type widgetGraphicsEffectBase struct {
	color                   Color
	bitmap                  *Bitmap
	enabledChangedPublisher EventPublisher
	enabled                 bool
}

func (wgeb *widgetGraphicsEffectBase) create(color Color, enabled bool) error {
	bitmap, err := NewBitmapWithTransparentPixels(Size{12, 12})
	if err != nil {
		return err
	}

	canvas, err := NewCanvasFromImage(bitmap)
	if err != nil {
		return err
	}
	defer canvas.Dispose()

	var succeeded bool
	defer func() {
		if !succeeded {
			bitmap.Dispose()
		}
	}()

	for i := 1; i <= 5; i++ {
		bmp, err := NewBitmapWithTransparentPixels(Size{i*2 + 2, i*2 + 2})
		if err != nil {
			return err
		}
		defer bmp.Dispose()

		bmpCanvas, err := NewCanvasFromImage(bmp)
		if err != nil {
			return err
		}
		defer bmpCanvas.Dispose()

		color := RGB(
			byte(math.Min(1.0, float64(color.R())/255.0-0.1+0.1*float64(i))*255.0),
			byte(math.Min(1.0, float64(color.G())/255.0-0.1+0.1*float64(i))*255.0),
			byte(math.Min(1.0, float64(color.B())/255.0-0.1+0.1*float64(i))*255.0),
		)

		brush, err := NewSolidColorBrush(color)
		if err != nil {
			return err
		}
		defer brush.Dispose()

		if err := bmpCanvas.FillRoundedRectangle(brush, Rectangle{0, 0, i*2 + 2, i*2 + 2}, Size{i * 2, i * 2}); err != nil {
			return err
		}

		bmpCanvas.Dispose()

		opacity := byte(borderGlowAlpha[i-1] * 255.0)

		canvas.DrawBitmapWithOpacity(bmp, Rectangle{5 - i, 5 - i, i*2 + 2, i*2 + 2}, opacity)
	}

	succeeded = true

	wgeb.color = color
	wgeb.bitmap = bitmap
	wgeb.enabled = enabled

	return nil
}

func (wgeb *widgetGraphicsEffectBase) Enabled() bool {
	return wgeb.enabled
}

func (wgeb *widgetGraphicsEffectBase) SetEnabled(enabled bool) {
	if enabled == wgeb.enabled {
		return
	}

	wgeb.enabled = true

	wgeb.enabledChangedPublisher.Publish()
}

func (wgeb *widgetGraphicsEffectBase) EnabledChanged() *Event {
	return wgeb.enabledChangedPublisher.Event()
}

func (wgeb *widgetGraphicsEffectBase) Dispose() {
	if wgeb.bitmap != nil {
		wgeb.bitmap.Dispose()
		wgeb.bitmap = nil
	}
}

type BorderGlowEffect struct {
	widgetGraphicsEffectBase
}

func NewBorderGlowEffect(color Color, enabled bool) (*BorderGlowEffect, error) {
	bge := new(BorderGlowEffect)

	if err := bge.create(color, enabled); err != nil {
		return nil, err
	}

	return bge, nil
}

func (bge *BorderGlowEffect) Draw(widget Widget, canvas *Canvas) error {
	b := widget.Bounds()

	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X - 5, b.Y - 5, 5, 5}, Rectangle{0, 0, 5, 5})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X, b.Y - 5, b.Width, 5}, Rectangle{5 + 1, 0, 1, 5})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X + b.Width, b.Y - 5, 5, 5}, Rectangle{5 + 2, 0, 5, 5})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X + b.Width, b.Y, 5, b.Height}, Rectangle{5 + 2, 5 + 1, 5, 1})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X + b.Width, b.Y + b.Height, 5, 5}, Rectangle{5 + 2, 5 + 2, 5, 5})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X, b.Y + b.Height, b.Width, 5}, Rectangle{5 + 1, 5 + 2, 1, 5})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X - 5, b.Y + b.Height, 5, 5}, Rectangle{0, 5 + 2, 5, 5})
	canvas.DrawBitmapPart(bge.bitmap, Rectangle{b.X - 5, b.Y, 5, b.Height}, Rectangle{0, 5 + 1, 5, 1})

	return nil
}

type DropShadowEffect struct {
	widgetGraphicsEffectBase
}

func NewDropShadowEffect(color Color, enabled bool) (*DropShadowEffect, error) {
	dse := new(DropShadowEffect)

	if err := dse.create(color, enabled); err != nil {
		return nil, err
	}

	return dse, nil
}

func (dse *DropShadowEffect) Draw(widget Widget, canvas *Canvas) error {
	b := widget.Bounds()

	canvas.DrawBitmapPart(dse.bitmap, Rectangle{b.X + b.Width, b.Y + 10 - 5, 5, 5}, Rectangle{5 + 2, 0, 5, 5})
	canvas.DrawBitmapPart(dse.bitmap, Rectangle{b.X + b.Width, b.Y + 10, 5, b.Height - 10}, Rectangle{5 + 2, 5 + 1, 5, 1})
	canvas.DrawBitmapPart(dse.bitmap, Rectangle{b.X + b.Width, b.Y + b.Height, 5, 5}, Rectangle{5 + 2, 5 + 2, 5, 5})
	canvas.DrawBitmapPart(dse.bitmap, Rectangle{b.X + 10, b.Y + b.Height, b.Width - 10, 5}, Rectangle{5 + 1, 5 + 2, 1, 5})
	canvas.DrawBitmapPart(dse.bitmap, Rectangle{b.X + 10 - 5, b.Y + b.Height, 5, 5}, Rectangle{0, 5 + 2, 5, 5})

	return nil
}

type widgetGraphicsEffectListObserver interface {
	onInsertedGraphicsEffect(index int, effect WidgetGraphicsEffect) error
	onRemovedGraphicsEffect(index int, effect WidgetGraphicsEffect) error
	onClearedGraphicsEffects() error
	onGraphicsEffectEnabledChanged(effect WidgetGraphicsEffect)
}

type WidgetGraphicsEffectList struct {
	items    []widgetGraphicsEffectListItem
	observer widgetGraphicsEffectListObserver
}

type widgetGraphicsEffectListItem struct {
	effect               WidgetGraphicsEffect
	enabledChangedHandle int
}

func newWidgetGraphicsEffectList(observer widgetGraphicsEffectListObserver) *WidgetGraphicsEffectList {
	return &WidgetGraphicsEffectList{observer: observer}
}

func (l *WidgetGraphicsEffectList) Add(item WidgetGraphicsEffect) error {
	return l.Insert(len(l.items), item)
}

func (l *WidgetGraphicsEffectList) At(index int) WidgetGraphicsEffect {
	return l.items[index].effect
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
		if item.effect == effect {
			return i
		}
	}

	return -1
}

func (l *WidgetGraphicsEffectList) Contains(effect WidgetGraphicsEffect) bool {
	return l.Index(effect) > -1
}

func (l *WidgetGraphicsEffectList) insertIntoSlice(index int, effect WidgetGraphicsEffect) {
	l.items = append(l.items, widgetGraphicsEffectListItem{nil, 0})
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = widgetGraphicsEffectListItem{
		effect: effect,
		enabledChangedHandle: effect.EnabledChanged().Attach(func() {
			l.observer.onGraphicsEffectEnabledChanged(effect)
		}),
	}
}

func (l *WidgetGraphicsEffectList) Insert(index int, effect WidgetGraphicsEffect) error {
	observer := l.observer

	l.insertIntoSlice(index, effect)

	if observer != nil {
		if err := observer.onInsertedGraphicsEffect(index, effect); err != nil {
			item := l.items[index]
			item.effect.EnabledChanged().Detach(item.enabledChangedHandle)

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

	item.effect.EnabledChanged().Detach(item.enabledChangedHandle)

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		if err := observer.onRemovedGraphicsEffect(index, item.effect); err != nil {
			l.insertIntoSlice(index, item.effect)
			return err
		}
	}

	return nil
}
