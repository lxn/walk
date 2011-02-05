// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"log"
	"os"
	"syscall"
)

import (
	. "walk/winapi/user32"
)

const splitterWindowClass = `\o/ Walk_Splitter_Class \o/`

var splitterWndProcPtr uintptr

func splitterWndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	s, ok := widgetsByHWnd[hwnd].(*Splitter)
	if !ok {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return s.wndProc(hwnd, msg, wParam, lParam, 0)
}

type Splitter struct {
	Container
	handleWidth   int
	mouseDownPos  Point
	draggedHandle *splitterHandle
}

func NewSplitter(parent IContainer) (*Splitter, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureRegisteredWindowClass(splitterWindowClass, splitterWndProc, &splitterWndProcPtr)

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr(splitterWindowClass), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	layout := NewHBoxLayout()
	s := &Splitter{
		Container: Container{
			Widget: Widget{
				hWnd:   hWnd,
				parent: parent,
			},
			layout: layout,
		},
		handleWidth: 4,
	}
	layout.container = s

	succeeded := false
	defer func() {
		if !succeeded {
			s.Dispose()
		}
	}()

	s.SetPersistent(true)

	s.children = newObservedWidgetList(s)

	s.SetFont(defaultFont)

	if err := parent.Children().Add(s); err != nil {
		return nil, err
	}

	widgetsByHWnd[hWnd] = s

	succeeded = true

	return s, nil
}

func (s *Splitter) LayoutFlags() LayoutFlags {
	return GrowHorz | GrowVert | ShrinkHorz | ShrinkVert
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

func (s *Splitter) onInsertingWidget(index int, widget IWidget) (err os.Error) {
	return s.Container.onInsertingWidget(index, widget)
}

func (s *Splitter) onInsertedWidget(index int, widget IWidget) (err os.Error) {
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

					bp, e := prev.Bounds()
					if e != nil {
						log.Println(e)
						return
					}

					bn, e := next.Bounds()
					if e != nil {
						log.Println(e)
						return
					}

					if s.Orientation() == Horizontal {
						xh, e := s.draggedHandle.X()
						if e != nil {
							log.Println(e)
							return
						}

						xnew := xh + x - s.mouseDownPos.X
						if xnew < bp.X {
							xnew = bp.X
						} else if xnew >= bn.X+bn.Width-s.handleWidth {
							xnew = bn.X + bn.Width - s.handleWidth
						}

						if e = s.draggedHandle.SetX(xnew); e != nil {
							log.Println(e)
							return
						}
					} else {
						yh, e := s.draggedHandle.Y()
						if e != nil {
							log.Println(e)
							return
						}

						ynew := yh + y - s.mouseDownPos.Y
						if ynew < bp.Y {
							ynew = bp.Y
						} else if ynew >= bn.Y+bn.Height-s.handleWidth {
							ynew = bn.Y + bn.Height - s.handleWidth
						}

						if e = s.draggedHandle.SetY(ynew); e != nil {
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

						prev.BeginUpdate()
						defer prev.Invalidate()
						defer prev.EndUpdate()
						next.BeginUpdate()
						defer next.Invalidate()
						defer next.EndUpdate()

						bh, e := dragHandle.Bounds()
						if e != nil {
							log.Println(e)
							return
						}

						bp, e := prev.Bounds()
						if e != nil {
							log.Println(e)
							return
						}

						bn, e := next.Bounds()
						if e != nil {
							log.Println(e)
							return
						}

						if s.Orientation() == Horizontal {
							bp.Width = bh.X - bp.X
							bn.Width -= (bh.X + bh.Width) - bn.X
							bn.X = bh.X + bh.Width
						} else {
							bp.Height = bh.Y - bp.Y
							bn.Height -= (bh.Y + bh.Height) - bn.Y
							bn.Y = bh.Y + bh.Height
						}

						if e = prev.SetBounds(bp); e != nil {
							log.Println(e)
							return
						}

						if e = next.SetBounds(bn); e != nil {
							log.Println(e)
							return
						}
					}
				})
			}
		}()
	}

	return s.Container.onInsertedWidget(index, widget)
}

func (s *Splitter) onRemovingWidget(index int, widget IWidget) (err os.Error) {
	return s.Container.onRemovingWidget(index, widget)
}

func (s *Splitter) onRemovedWidget(index int, widget IWidget) (err os.Error) {
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

	err = s.Container.onRemovedWidget(index, widget)
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
