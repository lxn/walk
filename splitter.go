// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"bytes"
	"log"
	"strconv"
	"strings"
)

import (
	"github.com/lxn/win"
)

const splitterWindowClass = `\o/ Walk_Splitter_Class \o/`

var splitterHandleDraggingBrush *SolidColorBrush

func init() {
	MustRegisterWindowClass(splitterWindowClass)

	splitterHandleDraggingBrush, _ = NewSolidColorBrush(Color(win.GetSysColor(win.COLOR_BTNSHADOW)))
}

type Splitter struct {
	ContainerBase
	handleWidth   int
	mouseDownPos  Point
	draggedHandle *splitterHandle
	persistent    bool
	removing      bool
}

func newSplitter(parent Container, orientation Orientation) (*Splitter, error) {
	layout := newSplitterLayout(Horizontal)
	s := &Splitter{
		ContainerBase: ContainerBase{
			layout: layout,
		},
		handleWidth: 5,
	}
	s.children = newWidgetList(s)
	layout.container = s

	if err := InitWidget(
		s,
		parent,
		splitterWindowClass,
		win.WS_VISIBLE,
		win.WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			s.Dispose()
		}
	}()

	s.SetBackground(NullBrush())

	if err := s.setOrientation(orientation); err != nil {
		return nil, err
	}

	s.SetPersistent(true)

	succeeded = true

	return s, nil
}

func NewHSplitter(parent Container) (*Splitter, error) {
	return newSplitter(parent, Horizontal)
}

func NewVSplitter(parent Container) (*Splitter, error) {
	return newSplitter(parent, Vertical)
}

func (s *Splitter) LayoutFlags() LayoutFlags {
	return s.layout.LayoutFlags()
}

func (s *Splitter) SizeHint() Size {
	return Size{100, 100}
}

func (s *Splitter) SetLayout(value Layout) error {
	return newError("not supported")
}

func (s *Splitter) HandleWidth() int {
	return s.handleWidth
}

func (s *Splitter) SetHandleWidth(value int) error {
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

func (s *Splitter) setOrientation(value Orientation) error {
	var cursor Cursor
	if value == Horizontal {
		cursor = CursorSizeWE()
	} else {
		cursor = CursorSizeNS()
	}

	for i, w := range s.Children().items {
		if i%2 == 1 {
			w.SetCursor(cursor)
		}
	}

	layout := s.layout.(*splitterLayout)
	return layout.SetOrientation(value)
}

func (s *Splitter) updateMarginsForFocusEffect() {
	var margins Margins
	var parentLayout Layout

	if s.parent != nil {
		if parentLayout = s.parent.Layout(); parentLayout != nil {
			if m := parentLayout.Margins(); m.HNear < 9 || m.HFar < 9 || m.VNear < 9 || m.VFar < 9 {
				parentLayout = nil
			}
		}
	}

	var affected bool
	if FocusEffect != nil {
		for _, w := range s.children.items {
			if w.GraphicsEffects().Contains(FocusEffect) {
				affected = true
				break
			}
		}
	}

	if affected {
		var marginsNeeded bool
		for _, w := range s.children.items {
			switch w.(type) {
			case *splitterHandle, *TabWidget, Container:

			default:
				marginsNeeded = true
				break
			}
		}

		if marginsNeeded {
			margins = Margins{5, 5, 5, 5}
		}
	}

	if parentLayout != nil {
		parentLayout.SetMargins(Margins{9 - margins.HNear, 9 - margins.VNear, 9 - margins.HFar, 9 - margins.VFar})
	}

	s.layout.SetMargins(margins)
}

func (s *Splitter) Persistent() bool {
	return s.persistent
}

func (s *Splitter) SetPersistent(value bool) {
	s.persistent = value
}

func (s *Splitter) SaveState() error {
	buf := bytes.NewBuffer(nil)

	count := s.children.Len()
	layout := s.Layout().(*splitterLayout)

	for i := 0; i < count; i += 2 {
		if i > 0 {
			buf.WriteString(" ")
		}

		item := layout.hwnd2Item[s.children.At(i).Handle()]
		size := item.oldExplicitSize
		if size == 0 {
			size = item.size
		}
		buf.WriteString(strconv.FormatInt(int64(size), 10))
	}

	s.WriteState(buf.String())

	for _, widget := range s.children.items {
		if persistable, ok := widget.(Persistable); ok {
			if err := persistable.SaveState(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Splitter) RestoreState() error {
	childCount := s.children.Len()/2 + 1
	if childCount == 0 {
		return nil
	}

	state, err := s.ReadState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	sizeStrs := strings.Split(state, " ")

	// FIXME: Solve this in a better way.
	if len(sizeStrs) != childCount {
		log.Print("*Splitter.RestoreState: failed due to unexpected child count (FIXME!)")
		return nil
	}

	s.SetSuspended(true)
	defer s.SetSuspended(false)

	layout := s.layout.(*splitterLayout)

	regularSpace := layout.spaceForRegularWidgets()

	for i, widget := range s.children.items {
		if i%2 == 0 {
			j := i/2 + i%2
			s := sizeStrs[j]

			size, err := strconv.Atoi(s)
			if err != nil {
				// OK, we probably got old style settings which were stored as fractions.
				fraction, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}

				size = int(float64(regularSpace) * fraction)
			}

			item := layout.hwnd2Item[widget.Handle()]
			item.size = size
			item.oldExplicitSize = size
		}

		if persistable, ok := widget.(Persistable); ok {
			if err := persistable.RestoreState(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Splitter) Fixed(widget Widget) bool {
	return s.layout.(*splitterLayout).Fixed(widget)
}

func (s *Splitter) SetFixed(widget Widget, fixed bool) error {
	item := s.layout.(*splitterLayout).hwnd2Item[widget.Handle()]
	if item == nil {
		return newErr("unknown widget")
	}

	item.fixed = fixed

	if b := widget.Bounds(); fixed && b.Width == 0 || b.Height == 0 {
		b.Width, b.Height = 100, 100
		widget.SetBounds(b)
		item.size = 100
	}

	return nil
}

func (s *Splitter) onInsertingWidget(index int, widget Widget) (err error) {
	return s.ContainerBase.onInsertingWidget(index, widget)
}

func (s *Splitter) onInsertedWidget(index int, widget Widget) (err error) {
	defer func() {
		if err != nil {
			return
		}

		s.updateMarginsForFocusEffect()
	}()

	_, isHandle := widget.(*splitterHandle)
	if isHandle {
		if s.Orientation() == Horizontal {
			widget.SetCursor(CursorSizeWE())
		} else {
			widget.SetCursor(CursorSizeNS())
		}
	} else {
		layout := s.Layout().(*splitterLayout)
		layout.hwnd2Item[widget.Handle()] = &splitterLayoutItem{stretchFactor: 1}

		if s.children.Len()%2 == 0 {
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
						if button != LeftButton {
							return
						}

						s.draggedHandle = handle
						s.mouseDownPos = Point{x, y}
						handle.SetBackground(splitterHandleDraggingBrush)
					})

					handle.MouseMove().Attach(func(x, y int, button MouseButton) {
						if s.draggedHandle == nil {
							return
						}

						handleIndex := s.children.Index(s.draggedHandle)
						bh := s.draggedHandle.Bounds()

						prev := s.children.At(handleIndex - 1)
						bp := prev.Bounds()
						msep := minSizeEffective(prev)

						next := s.children.At(handleIndex + 1)
						bn := next.Bounds()
						msen := minSizeEffective(next)

						if s.Orientation() == Horizontal {
							xh := s.draggedHandle.X()

							xnew := xh + x - s.mouseDownPos.X
							if xnew < bp.X+msep.Width {
								xnew = bp.X + msep.Width
							} else if xnew >= bn.X+bn.Width-msen.Width-s.handleWidth {
								xnew = bn.X + bn.Width - msen.Width - s.handleWidth
							}

							if e := s.draggedHandle.SetX(xnew); e != nil {
								return
							}
						} else {
							yh := s.draggedHandle.Y()

							ynew := yh + y - s.mouseDownPos.Y
							if ynew < bp.Y+msep.Height {
								ynew = bp.Y + msep.Height
							} else if ynew >= bn.Y+bn.Height-msen.Height-s.handleWidth {
								ynew = bn.Y + bn.Height - msen.Height - s.handleWidth
							}

							if e := s.draggedHandle.SetY(ynew); e != nil {
								return
							}
						}

						rc := bh.toRECT()
						if s.Orientation() == Horizontal {
							rc.Left -= int32(bp.X)
							rc.Right -= int32(bp.X)
						} else {
							rc.Top -= int32(bp.Y)
							rc.Bottom -= int32(bp.Y)
						}
						win.InvalidateRect(prev.Handle(), &rc, true)

						rc = bh.toRECT()
						if s.Orientation() == Horizontal {
							rc.Left -= int32(bn.X)
							rc.Right -= int32(bn.X)
						} else {
							rc.Top -= int32(bn.Y)
							rc.Bottom -= int32(bn.Y)
						}
						win.InvalidateRect(next.Handle(), &rc, true)

						s.draggedHandle.Invalidate()
					})

					handle.MouseUp().Attach(func(x, y int, button MouseButton) {
						if s.draggedHandle == nil {
							return
						}

						dragHandle := s.draggedHandle

						handleIndex := s.children.Index(dragHandle)
						prev := s.children.At(handleIndex - 1)
						next := s.children.At(handleIndex + 1)

						s.draggedHandle = nil
						dragHandle.SetBackground(NullBrush())
						prev.AsWidgetBase().invalidateBorderInParent()
						next.AsWidgetBase().invalidateBorderInParent()

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

						prevItem := layout.hwnd2Item[prev.Handle()]
						prevItem.size = sizePrev
						prevItem.oldExplicitSize = sizePrev

						nextItem := layout.hwnd2Item[next.Handle()]
						nextItem.size = sizeNext
						nextItem.oldExplicitSize = sizeNext
					})
				}
			}()
		}
	}

	return s.ContainerBase.onInsertedWidget(index, widget)
}

func (s *Splitter) onRemovingWidget(index int, widget Widget) (err error) {
	return s.ContainerBase.onRemovingWidget(index, widget)
}

func (s *Splitter) onRemovedWidget(index int, widget Widget) (err error) {
	defer func() {
		if err != nil {
			return
		}

		s.updateMarginsForFocusEffect()
	}()

	_, isHandle := widget.(*splitterHandle)
	if !s.removing && isHandle && s.children.Len()%2 == 1 {
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

			s.removing = true
			handle := s.children.items[handleIndex]
			if err = handle.SetParent(nil); err == nil {
				s.children.items = append(s.children.items[:index], s.children.items[index+1:]...)

				s.layout.Update(true)

				handle.Dispose()
			}

			s.removing = false
		}()
	}

	err = s.ContainerBase.onRemovedWidget(index, widget)

	return
}

func (s *Splitter) onClearingWidgets() (err error) {
	panic("not implemented")
}

func (s *Splitter) onClearedWidgets() (err error) {
	panic("not implemented")
}
