// Copyright 2021 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import (
	"github.com/lxn/walk"
)

type OpenGL struct {
	// Window

	Accessibility      Accessibility
	Background         Brush
	ContextMenuItems   []MenuItem
	DoubleBuffering    bool
	Enabled            Property
	Font               Font
	MaxSize            Size
	MinSize            Size
	Name               string
	OnBoundsChanged    walk.EventHandler
	OnKeyDown          walk.KeyEventHandler
	OnKeyPress         walk.KeyEventHandler
	OnKeyUp            walk.KeyEventHandler
	OnMouseDown        walk.MouseEventHandler
	OnMouseMove        walk.MouseEventHandler
	OnMouseUp          walk.MouseEventHandler
	OnSizeChanged      walk.EventHandler
	Persistent         bool
	RightToLeftReading bool
	ToolTipText        Property
	Visible            Property

	// Widget

	Alignment          Alignment2D
	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	GraphicsEffects    []walk.WidgetGraphicsEffect
	Row                int
	RowSpan            int
	StretchFactor      int

	// OpenGL

	AssignTo    **walk.OpenGL
	Setup       walk.GLFunc
	Paint       walk.GLFunc
	Teardown    walk.GLFunc
	Style       int
	PixelFormat []int32 // for wglCreateContextAttribsARB
	Context     []int32 // for wglCreateContextAttribsARB
}

func (gl OpenGL) Create(builder *Builder) error {
	w, err := walk.NewOpenGL(builder.Parent(), uint32(gl.Style), gl.Setup, gl.Paint, gl.Teardown, gl.PixelFormat, gl.Context)
	if err != nil {
		return err
	}

	if gl.AssignTo != nil {
		*gl.AssignTo = w
	}

	return builder.InitWidget(gl, w, func() error { return nil })
}
