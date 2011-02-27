// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

import (
	. "walk/winapi/user32"
)

const splitterWindowClass = `\o/ Walk_Splitter_Class \o/`

var splitterWindowClassRegistered bool

type Splitter struct {
	ContainerBase
	handleWidth   int
	mouseDownPos  Point
	draggedHandle *splitterHandle
	widget2Fixed  map[*WidgetBase]bool
	oldClientSize Size
	persistent    bool
}

func NewSplitter(parent Container) (*Splitter, os.Error) {
	ensureRegisteredWindowClass(splitterWindowClass, &splitterWindowClassRegistered)

	layout := NewHBoxLayout()
	s := &Splitter{
		ContainerBase: ContainerBase{
			layout: layout,
		},
		handleWidth:  4,
		widget2Fixed: make(map[*WidgetBase]bool),
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
	return HGrow | VGrow | HShrink | VShrink
}

func (s *Splitter) PreferredSize() Size {
	return s.dialogBaseUnitsToPixels(Size{100, 100})
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
	layout := s.layout.(*BoxLayout)
	return layout.Orientation()
}

func (s *Splitter) SetOrientation(value Orientation) os.Error {
	layout := s.layout.(*BoxLayout)
	return layout.SetOrientation(value)
}

func (s *Splitter) Fixed(widget Widget) bool {
	return s.widget2Fixed[widget.BaseWidget()]
}

func (s *Splitter) SetFixed(widget Widget, fixed bool) os.Error {
	if !s.Children().containsHandle(widget.BaseWidget().hWnd) {
		return newError("unknown widget")
	}

	s.widget2Fixed[widget.BaseWidget()] = fixed

	return nil
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
	orientation := s.Orientation()
	for i := 0; i < count; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}

		var size int
		if orientation == Horizontal {
			size = s.children.At(i).Width()
		} else {
			size = s.children.At(i).Height()
		}

		buf.WriteString(strconv.Itoa(size))
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

	sizeStrs := strings.Split(state, " ", -1)

	if len(sizeStrs) != s.children.Len() {
		return newError("unexpected child count")
	}

	layout := s.layout.(*BoxLayout)

	s.SetSuspended(true)
	defer s.SetSuspended(false)

	for i, widget := range s.children.items {
		size, err := strconv.Atoi(sizeStrs[i])
		if err != nil {
			return err
		}

		layout.SetStretchFactor(widget, size)

		if persistable, ok := widget.(Persistable); ok {
			if err := persistable.RestoreState(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Splitter) onResize() {
	clientSize := s.ClientBounds().Size()

	if s.oldClientSize.Width > 0 || s.oldClientSize.Height > 0 {
		layout := s.layout.(*BoxLayout)

		widgets := s.Children()
		orientation := s.Orientation()

		s.SetSuspended(true)
		defer s.SetSuspended(false)

		var fixedSizeTotal int
		for i := widgets.Len() - 1; i >= 0; i-- {
			widget := widgets.At(i)

			if i%2 == 1 {
				fixedSizeTotal += s.handleWidth
			} else if s.Fixed(widget) {
				fixedSizeTotal += widget.Width()
			}
		}

		for i := widgets.Len() - 1; i >= 0; i-- {
			widget := widgets.At(i)

			var stretch int

			if i%2 == 1 {
				stretch = s.handleWidth
			} else if s.Fixed(widget) {
				if orientation == Horizontal {
					stretch = widget.Width()
				} else {
					stretch = widget.Height()
				}
			} else {
				if orientation == Horizontal {
					stretch = int(float64(widget.Width()) * float64(clientSize.Width-fixedSizeTotal) / float64(s.oldClientSize.Width-fixedSizeTotal))
				} else {
					stretch = int(float64(widget.Height()) * float64(clientSize.Height-fixedSizeTotal) / float64(s.oldClientSize.Height-fixedSizeTotal))
				}
			}

			layout.SetStretchFactor(widget, stretch)
		}
	}

	s.oldClientSize = clientSize
}

func (s *Splitter) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_SIZE, WM_SIZING:
		s.onResize()
	}

	return s.ContainerBase.wndProc(hwnd, msg, wParam, lParam)
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
							log.Println(e)
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
							log.Println(e)
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

						if s.Orientation() == Horizontal {
							bp.Width = bh.X - bp.X
							bn.Width -= (bh.X + bh.Width) - bn.X
							bn.X = bh.X + bh.Width
						} else {
							bp.Height = bh.Y - bp.Y
							bn.Height -= (bh.Y + bh.Height) - bn.Y
							bn.Y = bh.Y + bh.Height
						}

						if e := prev.SetBounds(bp); e != nil {
							log.Println(e)
							return
						}

						if e := next.SetBounds(bn); e != nil {
							log.Println(e)
							return
						}
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
