// Copyright 2019 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func main() {
	var mw *walk.MainWindow
	var lb *walk.ListBox
	var items []logEntry
	for i := 10000; i > 0; i-- {
		items = append(items, logEntry{time.Now().Add(-time.Second * time.Duration(i)), fmt.Sprintf("Some stuff just happend. %s", strings.TrimSpace(strings.Repeat("blah ", i%100)))})
	}
	model := &logModel{items: items}
	styler := &Styler{
		lb:                  &lb,
		model:               model,
		dpi2StampSize:       make(map[int]walk.Size),
		widthDPI2WsPerLine:  make(map[widthDPI]int),
		textWidthDPI2Height: make(map[textWidthDPI]int),
	}

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Walk ListBox Owner Drawing Example",
		MinSize:  Size{200, 200},
		Size:     Size{800, 600},
		Font:     Font{Family: "Segoe UI", PointSize: 9},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				DoubleBuffering: true,
				Layout:          VBox{},
				Children: []Widget{
					ListBox{
						AssignTo:       &lb,
						MultiSelection: true,
						Model:          model,
						ItemStyler:     styler,
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	cancel := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				mw.Synchronize(func() {
					trackLatest := lb.ItemVisible(len(model.items)-1) && len(lb.SelectedIndexes()) <= 1

					model.items = append(model.items, logEntry{time.Now(), "Some new stuff."})
					index := len(model.items) - 1
					model.PublishItemsInserted(index, index)

					if trackLatest {
						lb.EnsureItemVisible(len(model.items) - 1)
					}
				})

			case <-cancel:
				break
			}
		}
	}()

	mw.Show()
	mw.Run()

	cancel <- true
}

type logModel struct {
	walk.ReflectListModelBase
	items []logEntry
}

func (m *logModel) Items() interface{} {
	return m.items
}

type logEntry struct {
	timestamp time.Time
	message   string
}

type widthDPI struct {
	width int // in native pixels
	dpi   int
}

type textWidthDPI struct {
	text  string
	width int // in native pixels
	dpi   int
}

type Styler struct {
	lb                  **walk.ListBox
	canvas              *walk.Canvas
	model               *logModel
	font                *walk.Font
	dpi2StampSize       map[int]walk.Size
	widthDPI2WsPerLine  map[widthDPI]int
	textWidthDPI2Height map[textWidthDPI]int // in native pixels
}

func (s *Styler) ItemHeightDependsOnWidth() bool {
	return true
}

func (s *Styler) DefaultItemHeight() int {
	dpi := (*s.lb).DPI()
	marginV := walk.IntFrom96DPI(marginV96dpi, dpi)

	return s.StampSize().Height + marginV*2
}

const (
	marginH96dpi int = 6
	marginV96dpi int = 2
	lineW96dpi   int = 1
)

func (s *Styler) ItemHeight(index, width int) int {
	dpi := (*s.lb).DPI()
	marginH := walk.IntFrom96DPI(marginH96dpi, dpi)
	marginV := walk.IntFrom96DPI(marginV96dpi, dpi)
	lineW := walk.IntFrom96DPI(lineW96dpi, dpi)

	msg := s.model.items[index].message

	twd := textWidthDPI{msg, width, dpi}

	if height, ok := s.textWidthDPI2Height[twd]; ok {
		return height + marginV*2
	}

	canvas, err := s.Canvas()
	if err != nil {
		return 0
	}

	stampSize := s.StampSize()

	wd := widthDPI{width, dpi}
	wsPerLine, ok := s.widthDPI2WsPerLine[wd]
	if !ok {
		bounds, _, err := canvas.MeasureTextPixels("W", (*s.lb).Font(), walk.Rectangle{Width: 9999999}, walk.TextCalcRect)
		if err != nil {
			return 0
		}
		wsPerLine = (width - marginH*4 - lineW - stampSize.Width) / bounds.Width
		s.widthDPI2WsPerLine[wd] = wsPerLine
	}

	if len(msg) <= wsPerLine {
		s.textWidthDPI2Height[twd] = stampSize.Height
		return stampSize.Height + marginV*2
	}

	bounds, _, err := canvas.MeasureTextPixels(msg, (*s.lb).Font(), walk.Rectangle{Width: width - marginH*4 - lineW - stampSize.Width, Height: 255}, walk.TextEditControl|walk.TextWordbreak|walk.TextEndEllipsis)
	if err != nil {
		return 0
	}

	s.textWidthDPI2Height[twd] = bounds.Height

	return bounds.Height + marginV*2
}

func (s *Styler) StyleItem(style *walk.ListItemStyle) {
	if canvas := style.Canvas(); canvas != nil {
		if style.Index()%2 == 1 && style.BackgroundColor == walk.Color(win.GetSysColor(win.COLOR_WINDOW)) {
			style.BackgroundColor = walk.Color(win.GetSysColor(win.COLOR_BTNFACE))
			if err := style.DrawBackground(); err != nil {
				return
			}
		}

		pen, err := walk.NewCosmeticPen(walk.PenSolid, style.LineColor)
		if err != nil {
			return
		}
		defer pen.Dispose()

		dpi := (*s.lb).DPI()
		marginH := walk.IntFrom96DPI(marginH96dpi, dpi)
		marginV := walk.IntFrom96DPI(marginV96dpi, dpi)
		lineW := walk.IntFrom96DPI(lineW96dpi, dpi)

		b := style.BoundsPixels()
		b.X += marginH
		b.Y += marginV

		item := s.model.items[style.Index()]

		style.DrawText(item.timestamp.Format(time.StampMilli), b, walk.TextEditControl|walk.TextWordbreak)

		stampSize := s.StampSize()

		x := b.X + stampSize.Width + marginH + lineW
		canvas.DrawLinePixels(pen, walk.Point{x, b.Y - marginV}, walk.Point{x, b.Y - marginV + b.Height})

		b.X += stampSize.Width + marginH*2 + lineW
		b.Width -= stampSize.Width + marginH*4 + lineW

		style.DrawText(item.message, b, walk.TextEditControl|walk.TextWordbreak|walk.TextEndEllipsis)
	}
}

func (s *Styler) StampSize() walk.Size {
	dpi := (*s.lb).DPI()

	stampSize, ok := s.dpi2StampSize[dpi]
	if !ok {
		canvas, err := s.Canvas()
		if err != nil {
			return walk.Size{}
		}

		bounds, _, err := canvas.MeasureTextPixels("Jan _2 20:04:05.000", (*s.lb).Font(), walk.Rectangle{Width: 9999999}, walk.TextCalcRect)
		if err != nil {
			return walk.Size{}
		}

		stampSize = bounds.Size()
		s.dpi2StampSize[dpi] = stampSize
	}

	return stampSize
}

func (s *Styler) Canvas() (*walk.Canvas, error) {
	if s.canvas != nil {
		return s.canvas, nil
	}

	canvas, err := (*s.lb).CreateCanvas()
	if err != nil {
		return nil, err
	}
	s.canvas = canvas
	(*s.lb).AddDisposable(canvas)

	return canvas, nil
}
