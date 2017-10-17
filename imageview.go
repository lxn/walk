// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
	"math"
	"strconv"
)

var imageViewBackgroundBrush, _ = NewSystemColorBrush(win.COLOR_APPWORKSPACE)

type ImageView struct {
	*CustomWidget
	image                  Image
	imageChangedPublisher  EventPublisher
	margin                 int
	marginChangedPublisher EventPublisher
}

func NewImageView(parent Container) (*ImageView, error) {
	iv := new(ImageView)

	cw, err := NewCustomWidget(parent, 0, func(canvas *Canvas, updateBounds Rectangle) error {
		return iv.drawImage(canvas, updateBounds)
	})
	if err != nil {
		return nil, err
	}

	iv.CustomWidget = cw

	iv.window = iv

	iv.SetInvalidatesOnResize(true)
	iv.SetPaintMode(PaintNoErase)

	iv.MustRegisterProperty("Image", NewProperty(
		func() interface{} {
			return iv.Image()
		},
		func(v interface{}) error {
			var img Image

			switch val := v.(type) {
			case Image:
				img = val

			case int:
				var err error
				if img, err = Resources.Image(strconv.Itoa(val)); err != nil {
					return err
				}

			case string:
				var err error
				if img, err = Resources.Image(val); err != nil {
					return err
				}

			default:
				return ErrInvalidType
			}

			return iv.SetImage(img)
		},
		iv.imageChangedPublisher.Event()))

	iv.MustRegisterProperty("Margin", NewProperty(
		func() interface{} {
			return iv.Margin()
		},
		func(v interface{}) error {
			return iv.SetMargin(v.(int))
		},
		iv.MarginChanged()))

	return iv, nil
}

func (iv *ImageView) Image() Image {
	return iv.image
}

func (iv *ImageView) SetImage(value Image) error {
	if value == iv.image {
		return nil
	}

	iv.image = value

	_, isMetafile := value.(*Metafile)
	iv.SetClearsBackground(isMetafile)

	err := iv.Invalidate()

	iv.imageChangedPublisher.Publish()

	return err
}

func (iv *ImageView) ImageChanged() *Event {
	return iv.imageChangedPublisher.Event()
}

func (iv *ImageView) Margin() int {
	return iv.margin
}

func (iv *ImageView) SetMargin(margin int) error {
	if margin == iv.margin {
		return nil
	}

	iv.margin = margin

	err := iv.Invalidate()

	iv.marginChangedPublisher.Publish()

	return err
}

func (iv *ImageView) MarginChanged() *Event {
	return iv.marginChangedPublisher.Event()
}

func (iv *ImageView) drawImage(canvas *Canvas, updateBounds Rectangle) error {
	if iv.image == nil {
		return nil
	}

	cb := iv.ClientBounds()

	canvas.FillRectangle(imageViewBackgroundBrush, cb)

	cb.Width -= iv.margin * 2
	cb.Height -= iv.margin * 2

	s := iv.image.Size()

	var scale float64
	if s.Width > cb.Width || s.Height > cb.Height {
		sx := float64(cb.Width) / float64(s.Width)
		sy := float64(cb.Height) / float64(s.Height)

		scale = math.Min(sx, sy)
	} else {
		scale = 1.0
	}

	w := int(float64(s.Width) * scale)
	h := int(float64(s.Height) * scale)
	x := iv.margin + (cb.Width-w)/2
	y := iv.margin + (cb.Height-h)/2

	return canvas.DrawImageStretched(iv.image, Rectangle{X: x, Y: y, Width: w, Height: h})
}
