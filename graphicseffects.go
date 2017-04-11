// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
	"log"
	"math"
	"unsafe"
)

var (
	borderGlowAlpha = []float64{0.4, 0.2, 0.15, 0.1, 0.15}

	defaultDropShadowEffect = NewDropShadowEffect(RGB(63, 63, 63))
)

type WidgetGraphicsEffect interface {
	Draw(widget Widget, renderTarget *win.ID2D1RenderTarget) error
}

type BorderGlowEffect struct {
	color Color
}

func NewBorderGlowEffect(color Color) *BorderGlowEffect {
	return &BorderGlowEffect{color: color}
}

var (
	id2d1Factory *win.ID2D1Factory
)

func init() {
	var factory unsafe.Pointer

	if !win.SUCCEEDED(win.D2D1CreateFactory(
		win.D2D1_FACTORY_TYPE_SINGLE_THREADED,
		&win.IID_ID2D1Factory,
		&win.D2D1_FACTORY_OPTIONS{DebugLevel: win.D2D1_DEBUG_LEVEL_NONE},
		&factory)) {
		log.Println("D2D1CreateFactory failed")
	}

	id2d1Factory = (*win.ID2D1Factory)(factory)
}

func (bge *BorderGlowEffect) Draw(widget Widget, renderTarget *win.ID2D1RenderTarget) error {
	bounds := widget.Bounds()

	for i := 1; i <= 5; i++ {
		width := float32(i)

		color := win.D2D1_COLOR_F{
			R: float32(math.Min(1.0, float64(bge.color.R())/255.0-0.1+0.1*float64(width))),
			G: float32(math.Min(1.0, float64(bge.color.G())/255.0-0.1+0.1*float64(width))),
			B: float32(math.Min(1.0, float64(bge.color.B())/255.0-0.1+0.1*float64(width))),
			A: float32(borderGlowAlpha[i-1]),
		}

		var scBrush *win.ID2D1SolidColorBrush
		if hr := renderTarget.CreateSolidColorBrush(&color, nil, &scBrush); !win.SUCCEEDED(hr) {
			return errorFromHRESULT("ID2D1RenderTarget.CreateSolidColorBrush", hr)
		}
		defer scBrush.Release()

		rr := win.D2D1_ROUNDED_RECT{
			Rect: win.D2D1_RECT_F{
				Left:   float32(bounds.X) - width,
				Top:    float32(bounds.Y) - width,
				Right:  float32(bounds.X+bounds.Width) + width,
				Bottom: float32(bounds.Y+bounds.Height) + width,
			},
			RadiusX: width,
			RadiusY: width,
		}

		brush := (*win.ID2D1Brush)(unsafe.Pointer(scBrush))

		// DrawRoundedRectangle does not work, because syscall does not support float args,
		// so we have to fill the whole thing...
		renderTarget.FillRoundedRectangle(&rr, brush)
	}

	return nil
}

type DropShadowEffect struct {
	color Color
}

func NewDropShadowEffect(color Color) *DropShadowEffect {
	return &DropShadowEffect{color: color}
}

func (dse *DropShadowEffect) Draw(widget Widget, renderTarget *win.ID2D1RenderTarget) error {
	bounds := widget.Bounds()

	for i := 1; i <= 5; i++ {
		width := float32(i)

		color := win.D2D1_COLOR_F{
			R: float32(math.Min(1.0, float64(dse.color.R())/255.0-0.1+0.1*float64(width))),
			G: float32(math.Min(1.0, float64(dse.color.G())/255.0-0.1+0.1*float64(width))),
			B: float32(math.Min(1.0, float64(dse.color.B())/255.0-0.1+0.1*float64(width))),
			A: float32(borderGlowAlpha[i-1]),
		}

		var scBrush *win.ID2D1SolidColorBrush
		if hr := renderTarget.CreateSolidColorBrush(&color, nil, &scBrush); !win.SUCCEEDED(hr) {
			return errorFromHRESULT("ID2D1RenderTarget.CreateSolidColorBrush", hr)
		}
		defer scBrush.Release()

		rr := win.D2D1_ROUNDED_RECT{
			Rect: win.D2D1_RECT_F{
				Left:   float32(bounds.X+10) - width,
				Top:    float32(bounds.Y+10) - width,
				Right:  float32(bounds.X+bounds.Width) + width,
				Bottom: float32(bounds.Y+bounds.Height) + width,
			},
			RadiusX: width,
			RadiusY: width,
		}

		brush := (*win.ID2D1Brush)(unsafe.Pointer(scBrush))

		// DrawRoundedRectangle does not work, because syscall does not support float args,
		// so we have to fill the whole thing...
		renderTarget.FillRoundedRectangle(&rr, brush)
	}

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

func (l *WidgetGraphicsEffectList) Add(item WidgetGraphicsEffect) error {
	return l.Insert(len(l.items), item)
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

func (l *WidgetGraphicsEffectList) Index(item WidgetGraphicsEffect) int {
	for i, widget := range l.items {
		if widget == item {
			return i
		}
	}

	return -1
}

func (l *WidgetGraphicsEffectList) Contains(item WidgetGraphicsEffect) bool {
	return l.Index(item) > -1
}

func (l *WidgetGraphicsEffectList) insertIntoSlice(index int, item WidgetGraphicsEffect) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
}

func (l *WidgetGraphicsEffectList) Insert(index int, item WidgetGraphicsEffect) error {
	observer := l.observer

	l.insertIntoSlice(index, item)

	if observer != nil {
		if err := observer.onInsertedGraphicsEffect(index, item); err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return err
		}
	}

	return nil
}

func (l *WidgetGraphicsEffectList) Len() int {
	return len(l.items)
}

func (l *WidgetGraphicsEffectList) Remove(item WidgetGraphicsEffect) error {
	index := l.Index(item)
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
