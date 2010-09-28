// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
	"walk/drawing"
	. "walk/winapi/gdi32"
	. "walk/winapi/user32"
)

type ImageView struct {
	CustomWidget
	image drawing.Image
}

func NewImageView(parent IContainer) (*ImageView, os.Error) {
	iv := &ImageView{}

	cw, err := NewCustomWidget(parent, 0, func(surface *drawing.Surface, updateBounds drawing.Rectangle) os.Error {
		return iv.drawImage(surface, updateBounds)
	})
	if err != nil {
		return nil, err
	}

	iv.CustomWidget = *cw

	widgetsByHWnd[iv.hWnd] = iv
	customWidgetsByHWND[iv.hWnd] = &iv.CustomWidget

	return iv, nil
}

func (iv *ImageView) Image() drawing.Image {
	return iv.image
}

func (iv *ImageView) SetImage(value drawing.Image) os.Error {
	iv.image = value

	cb, err := iv.ClientBounds()
	if err != nil {
		return err
	}

	r := &RECT{cb.X, cb.Y, cb.X + cb.Width - 1, cb.Y + cb.Height - 1}

	if !InvalidateRect(iv.hWnd, r, true) {
		return newError("InvalidateRect failed")
	}

	return nil
}

func (iv *ImageView) drawImage(surface *drawing.Surface, updateBounds drawing.Rectangle) os.Error {
	if iv.image == nil {
		return nil
	}

	bounds, err := iv.ClientBounds()
	if err != nil {
		return err
	}

	return surface.DrawImageStretched(iv.image, bounds)
}
