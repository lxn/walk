// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"os"
	"strconv"
	"strings"
)

import . "walk/winapi"

const splitterWindowClass = `\o/ Walk_Splitter_Class \o/`

var splitterWindowClassRegistered bool

type Splitter struct {
	ContainerBase
	handleWidth   int
	mouseDownPos  Point
	draggedHandle *splitterHandle
	persistent    bool
}

func NewSplitter(parent Container) (*Splitter, os.Error) {
	ensureRegisteredWindowClass(splitterWindowClass, &splitterWindowClassRegistered)

	layout := newSplitterLayout(Horizontal)
	s := &Splitter{
		ContainerBase: ContainerBase{
			layout: layout,
		},
		handleWidth: 4,
	}
	s.children = newWidgetList(s)
	layout.container = s

	if err := initChildWidget(
		s,
		parent,
		splitterWindowClass,
		WS_VISIBLE,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	s.SetPersistent(true)

	return s, nil
}

func (s *Splitter) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | ShrinkableVert | GrowableHorz | GrowableVert | GreedyHorz | GreedyVert
}

func (s *Splitter) MinSizeHint() Size {
	return Size{10, 10}
}

func (s *Splitter) SizeHint() Size {
	return Size{100, 100}
}

func (s *Splitter) SetLayout(value Layout) os.Error {
	return newError("not supported")
}

func (s *Splitter) HandleWidth() int {
	return s.handleWidth
}

func (s *Splitter) SetHandleWidth(value int) os.Error {
	if value == s.handleWidth {
		return nil
	}

	if value < 1 {
		return newError("invalid handle width")
	}

	s.handleWidth = value

	return s.layout.Update(false)
}

func (s *Splitter) Orientation() Orientation {
	layout := s.layout.(*splitterLayout)
	return layout.Orientation()
}

func (s *Splitter) SetOrientation(value Orientation) os.Error {
	layout := s.layout.(*splitterLayout)
	return layout.SetOrientation(value)
}

func (s *Splitter) Persistent() bool {
	return s.persistent
}

func (s *Splitter) SetPersistent(value bool) {
	s.persistent = value
}

func (s *Splitter) SaveState() os.Error {
	buf := bytes.NewBuffer(nil)

	count := s.children.Len()
	layout := s.Layout().(*splitterLayout)

	for i := 0; i < count; i += 2 {
		if i > 0 {
			buf.WriteString(" ")
		}

		buf.WriteString(strconv.Ftoa64(layout.fractions[i/2], 'f', -1))
	}

	s.putState(buf.String())

	for _, widget := range s.children.items {
		if persistable, ok := widget.(Persistable); ok {
			if err := persistable.SaveState(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Splitter) RestoreState() os.Error {
	state, err := s.getState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	fractionStrs := strings.Split(state, " ", -1)

	layout := s.layout.(*splitterLayout)

	s.SetSuspended(true)
	defer s.SetSuspended(false)

	var fractionsTotal float64
	var fractions []float64

	for i, widget := range s.children.items {
		if i%2 == 0 {
			fraction, err := strconv.Atof64(fractionStrs[i/2+i%2])
			if err != nil {
				return err
			}

			fractionsTotal += fraction
			fractions = append(fractions, fraction)
		}

		if persistable, ok := widget.(Persistable); ok {
			if err := persistable.RestoreState(); err != nil {
				return err
			}
		}
	}

	for i := range fractions {
		fractions[i] = fractions[i] / fractionsTotal
	}

	return layout.SetFractions(fractions)
}

func (s *Splitter) onInsertingWidget(index int, widget Widget) (err os.Error) {
	return s.ContainerBase.onInsertingWidget(index, widget)
}

func (s *Splitter) onInsertedWidget(index int, widget Widget) (err os.Error) {
	_, isHandle := widget.(*splitterHandle)
	if isHandle {
		if s.Orientation() == Horizontal {
			widget.SetCursor(CursorSizeWE())
		} else {
			widget.SetCursor(CursorSizeNS())
		}
	} else if s.children.Len()%2 == 0 {
		defer func() {
			if err != nil {
				return
			}

			var handle *splitterHandle
			handle, err = newSplitterHandle(s)
			if err != nil {
				return
			}

			var handleIndex int
			if index == 0 {
				handleIndex = 1
			} else {
				handleIndex = index
			}
			err = s.children.Insert(handleIndex, handle)
			if err == nil {
				// FIXME: These handlers will be leaked, if widgets get removed.
				handle.MouseDown().Attach(func(x, y int, button MouseButton) {
					s.draggedHandle = handle
					s.mouseDownPos = Point{x, y}
				})

				handle.MouseMove().Attach(func(x, y int, button MouseButton) {
					if s.draggedHandle == nil {
						return
					}

					handleIndex := s.children.Index(s.draggedHandle)
					prev := s.children.At(handleIndex - 1)
					next := s.children.At(handleIndex + 1)

					bp := prev.Bounds()
					bn := next.Bounds()

					if s.Orientation() == Horizontal {
						xh := s.draggedHandle.X()

						xnew := xh + x - s.mouseDownPos.X
						if xnew < bp.X {
							xnew = bp.X
						} else if xnew >= bn.X+bn.Width-s.handleWidth {
							xnew = bn.X + bn.Width - s.handleWidth
						}

						if e := s.draggedHandle.SetX(xnew); e != nil {
							return
						}
					} else {
						yh := s.draggedHandle.Y()

						ynew := yh + y - s.mouseDownPos.Y
						if ynew < bp.Y {
							ynew = bp.Y
						} else if ynew >= bn.Y+bn.Height-s.handleWidth {
							ynew = bn.Y + bn.Height - s.handleWidth
						}

						if e := s.draggedHandle.SetY(ynew); e != nil {
							return
						}
					}
				})

				handle.MouseUp().Attach(func(x, y int, button MouseButton) {
					if s.draggedHandle != nil {
						dragHandle := s.draggedHandle
						s.draggedHandle = nil

						handleIndex := s.children.Index(dragHandle)
						prev := s.children.At(handleIndex - 1)
						next := s.children.At(handleIndex + 1)

						prev.SetSuspended(true)
						defer prev.Invalidate()
						defer prev.SetSuspended(false)
						next.SetSuspended(true)
						defer next.Invalidate()
						defer next.SetSuspended(false)

						bh := dragHandle.Bounds()
						bp := prev.Bounds()
						bn := next.Bounds()

						var sizePrev int
						var sizeNext int

						if s.Orientation() == Horizontal {
							bp.Width = bh.X - bp.X
							bn.Width -= (bh.X + bh.Width) - bn.X
							bn.X = bh.X + bh.Width
							sizePrev = bp.Width
							sizeNext = bn.Width
						} else {
							bp.Height = bh.Y - bp.Y
							bn.Height -= (bh.Y + bh.Height) - bn.Y
							bn.Y = bh.Y + bh.Height
							sizePrev = bp.Height
							sizeNext = bn.Height
						}

						if e := prev.SetBounds(bp); e != nil {
							return
						}

						if e := next.SetBounds(bn); e != nil {
							return
						}

						layout := s.Layout().(*splitterLayout)
						space := float64(layout.spaceForRegularWidgets())
						fractions := layout.fractions
						i := handleIndex - 1
						prevFracIndex := i/2 + i%2
						nextFracIndex := prevFracIndex + 1
						fractions[prevFracIndex] = float64(sizePrev) / space
						fractions[nextFracIndex] = float64(sizeNext) / space
					}
				})
			}
		}()
	}

	return s.ContainerBase.onInsertedWidget(index, widget)
}

func (s *Splitter) onRemovingWidget(index int, widget Widget) (err os.Error) {
	return s.ContainerBase.onRemovingWidget(index, widget)
}

func (s *Splitter) onRemovedWidget(index int, widget Widget) (err os.Error) {
	_, isHandle := widget.(*splitterHandle)
	if isHandle && s.children.Len()%2 == 1 {
		return newError("cannot remove splitter handle")
	}

	if !isHandle && s.children.Len() > 1 {
		defer func() {
			if err != nil {
				return
			}

			var handleIndex int
			if index == 0 {
				handleIndex = 0
			} else {
				handleIndex = index - 1
			}
			err = s.children.RemoveAt(handleIndex)
		}()
	}

	err = s.ContainerBase.onRemovedWidget(index, widget)
	if isHandle && err == nil {
		widget.Dispose()
	}

	return
}

func (s *Splitter) onClearingWidgets() (err os.Error) {
	panic("not implemented")
}

func (s *Splitter) onClearedWidgets() (err os.Error) {
	panic("not implemented")
}
