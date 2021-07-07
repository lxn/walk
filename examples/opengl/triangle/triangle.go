// Copyright 2021 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	MainWindow{
		Title:   "Walk OpenGL Example",
		MinSize: Size{320, 240},
		Layout:  HBox{},
		Children: []Widget{
			OpenGL{
				Setup: func(*walk.OpenGLContext) error {
					return gl.Init()
				},
				Paint: func(glc *walk.OpenGLContext) error {
					sz := glc.Widget().Size()
					gl.Viewport(0, 0, int32(sz.Width), int32(sz.Height))
					gl.Clear(gl.COLOR_BUFFER_BIT)
					gl.Begin(gl.TRIANGLES)
					gl.Color3f(1.0, 0.0, 0.0)
					gl.Vertex2i(0, 1)
					gl.Color3f(0.0, 1.0, 0.0)
					gl.Vertex2i(-1, -1)
					gl.Color3f(0.0, 0.0, 1.0)
					gl.Vertex2i(1, -1)
					gl.End()
					gl.Flush()
					return nil
				},
			},
		},
	}.Run()
}
