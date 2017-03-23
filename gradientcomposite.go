// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
	"unsafe"
)

type GradientComposite struct {
	*Composite
	vertical                 bool
	color1                   Color
	color2                   Color
	verticalChangedPublisher EventPublisher
	color1ChangedPublisher   EventPublisher
	color2ChangedPublisher   EventPublisher
	brush                    *BitmapBrush
}

func NewGradientComposite(parent Container) (*GradientComposite, error) {
	composite, err := NewComposite(parent)
	if err != nil {
		return nil, err
	}

	gc := &GradientComposite{Composite: composite}

	succeeded := false
	defer func() {
		if !succeeded {
			gc.Dispose()
		}
	}()

	if err := InitWrapperWindow(gc); err != nil {
		return nil, err
	}

	gc.MustRegisterProperty("Vertical", NewBoolProperty(
		func() bool {
			return gc.Vertical()
		},
		func(b bool) error {
			gc.SetVertical(b)
			return nil
		},
		gc.verticalChangedPublisher.Event()))

	gc.MustRegisterProperty("Color1", NewProperty(
		func() interface{} {
			return float64(uint32(gc.Color1()))
		},
		func(v interface{}) error {
			var c Color

			switch v := v.(type) {
			case Color:
				c = v

			case uint32:
				c = Color(v)

			case float64:
				c = Color(uint32(v))

			default:
				return ErrInvalidType
			}

			return gc.SetColor1(c)
		},
		gc.color1ChangedPublisher.Event()))

	gc.MustRegisterProperty("Color2", NewProperty(
		func() interface{} {
			return float64(uint32(gc.Color2()))
		},
		func(v interface{}) error {
			var c Color

			switch v := v.(type) {
			case Color:
				c = v

			case uint32:
				c = Color(v)

			case float64:
				c = Color(uint32(v))

			default:
				return ErrInvalidType
			}

			return gc.SetColor2(c)
		},
		gc.color2ChangedPublisher.Event()))

	succeeded = true

	return gc, nil
}

func (gc *GradientComposite) Vertical() bool {
	return gc.vertical
}

func (gc *GradientComposite) SetVertical(vertical bool) (err error) {
	if vertical == gc.vertical {
		return nil
	}

	old := gc.vertical

	defer func() {
		if err != nil {
			gc.vertical = old
		}
	}()

	gc.vertical = vertical

	if err = gc.updateBackground(); err != nil {
		return
	}

	gc.verticalChangedPublisher.Publish()

	return
}

func (gc *GradientComposite) Color1() Color {
	return gc.color1
}

func (gc *GradientComposite) SetColor1(c Color) (err error) {
	if c == gc.color1 {
		return nil
	}

	old := gc.color1

	defer func() {
		if err != nil {
			gc.color1 = old
		}
	}()

	gc.color1 = c

	if err = gc.updateBackground(); err != nil {
		return
	}

	gc.color1ChangedPublisher.Publish()

	return
}

func (gc *GradientComposite) Color2() Color {
	return gc.color2
}

func (gc *GradientComposite) SetColor2(c Color) (err error) {
	if c == gc.color2 {
		return nil
	}

	old := gc.color2

	defer func() {
		if err != nil {
			gc.color2 = old
		}
	}()

	gc.color2 = c

	if err = gc.updateBackground(); err != nil {
		return
	}

	gc.color2ChangedPublisher.Publish()

	return
}

func (gc *GradientComposite) updateBackground() error {
	size := gc.ClientBounds().Size()
	if size.Width < 1 || size.Height < 1 {
		return nil
	}

	if gc.brush != nil {
		gc.brush.Dispose()
		gc.brush.Bitmap().Dispose()
		gc.brush = nil
	}

	bmp, err := NewBitmap(size)
	if err != nil {
		return err
	}
	defer func() {
		if gc.brush == nil {
			bmp.Dispose()
		}
	}()

	canvas, err := NewCanvasFromImage(bmp)
	if err != nil {
		return err
	}
	defer canvas.Dispose()

	vertices := [2]win.TRIVERTEX{
		{
			X:     0,
			Y:     0,
			Red:   uint16(gc.color1.R()) * 256,
			Green: uint16(gc.color1.G()) * 256,
			Blue:  uint16(gc.color1.B()) * 256,
			Alpha: 0,
		}, {
			X:     int32(size.Width),
			Y:     int32(size.Height),
			Red:   uint16(gc.color2.R()) * 256,
			Green: uint16(gc.color2.G()) * 256,
			Blue:  uint16(gc.color2.B()) * 256,
			Alpha: 0,
		},
	}

	indices := win.GRADIENT_RECT{
		UpperLeft:  0,
		LowerRight: 1,
	}

	var orientation uint32
	if gc.vertical {
		orientation = win.GRADIENT_FILL_RECT_V
	} else {
		orientation = win.GRADIENT_FILL_RECT_H
	}

	if !win.GradientFill(canvas.HDC(), &vertices[0], 2, unsafe.Pointer(&indices), 1, orientation) {
		return newErr("GradientFill failed")
	}

	gc.brush, err = NewBitmapBrush(bmp)
	if err != nil {
		return err
	}

	gc.SetBackground(gc.brush)

	return nil
}

func (gc *GradientComposite) Dispose() {
	if gc.brush != nil {
		gc.SetBackground(nil)
		gc.brush.Dispose()
		gc.brush.Bitmap().Dispose()
		gc.brush = nil
	}

	gc.Composite.Dispose()
}

func (gc *GradientComposite) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_SIZE, win.WM_SIZING:
		size := gc.ClientBounds().Size()
		if gc.brush != nil && gc.brush.bitmap.size == size {
			break
		}

		gc.updateBackground()
	}

	return gc.Composite.WndProc(hwnd, msg, wParam, lParam)
}
