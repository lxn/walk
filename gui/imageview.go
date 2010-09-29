// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
)

import (
	"walk/drawing"
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

	iv.SetInvalidatesOnResize(true)

	widgetsByHWnd[iv.hWnd] = iv
	customWidgetsByHWND[iv.hWnd] = &iv.CustomWidget

	return iv, nil
}

func (iv *ImageView) Image() drawing.Image {
	return iv.image
}

func (iv *ImageView) SetImage(value drawing.Image) os.Error {
	iv.image = value

	return iv.Invalidate()
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
