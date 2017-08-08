// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

type LinkLabel struct {
	WidgetBase
	textChangedPublisher   EventPublisher
	linkActivatedPublisher LinkLabelLinkEventPublisher
}

func NewLinkLabel(parent Container) (*LinkLabel, error) {
	ll := new(LinkLabel)

	if err := InitWidget(
		ll,
		parent,
		"SysLink",
		win.WS_TABSTOP|win.WS_VISIBLE,
		0); err != nil {
		return nil, err
	}

	ll.SetBackground(nullBrushSingleton)

	ll.MustRegisterProperty("Text", NewProperty(
		func() interface{} {
			return ll.Text()
		},
		func(v interface{}) error {
			return ll.SetText(v.(string))
		},
		ll.textChangedPublisher.Event()))

	return ll, nil
}

func (*LinkLabel) LayoutFlags() LayoutFlags {
	return GrowableVert
}

func (ll *LinkLabel) MinSizeHint() Size {
	var s win.SIZE

	ll.SendMessage(win.LM_GETIDEALSIZE, uintptr(ll.maxSize.Width), uintptr(unsafe.Pointer(&s)))

	return Size{int(s.CX), int(s.CY)}
}

func (ll *LinkLabel) SizeHint() Size {
	return ll.MinSizeHint()
}

func (ll *LinkLabel) Text() string {
	return windowText(ll.hWnd)
}

func (ll *LinkLabel) SetText(value string) error {
	if value == ll.Text() {
		return nil
	}

	if err := setWindowText(ll.hWnd, value); err != nil {
		return err
	}

	return ll.updateParentLayout()
}

func (ll *LinkLabel) LinkActivated() *LinkLabelLinkEvent {
	return ll.linkActivatedPublisher.Event()
}

func (ll *LinkLabel) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_NOTIFY:
		nml := (*win.NMLINK)(unsafe.Pointer(lParam))

		switch nml.Hdr.Code {
		case win.NM_CLICK, win.NM_RETURN:
			ll.linkActivatedPublisher.Publish(ll.linkFromLITEM(&nml.Item))
		}

	case win.WM_SETTEXT:
		ll.textChangedPublisher.Publish()

	case win.WM_SIZE, win.WM_SIZING:
		ll.Invalidate()
	}

	return ll.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

func (ll *LinkLabel) Link(index int) (*LinkLabelLink, error) {
	li := win.LITEM{
		ILink: int32(index),
		Mask:  win.LIF_ITEMID | win.LIF_ITEMINDEX | win.LIF_URL,
	}

	if win.TRUE != ll.SendMessage(win.LM_GETITEM, 0, uintptr(unsafe.Pointer(&li))) {
		return nil, newErr("LM_GETITEM")
	}

	return ll.linkFromLITEM(&li), nil
}

func (ll *LinkLabel) linkFromLITEM(li *win.LITEM) *LinkLabelLink {
	return &LinkLabelLink{
		ll:    ll,
		index: int(li.ILink),
		id:    syscall.UTF16ToString(li.SzID[:]),
		url:   syscall.UTF16ToString(li.SzUrl[:]),
	}
}

type LinkLabelLinkEventHandler func(link *LinkLabelLink)

type LinkLabelLinkEvent struct {
	handlers []LinkLabelLinkEventHandler
}

func (e *LinkLabelLinkEvent) Attach(handler LinkLabelLinkEventHandler) int {
	for i, h := range e.handlers {
		if h == nil {
			e.handlers[i] = handler
			return i
		}
	}

	e.handlers = append(e.handlers, handler)
	return len(e.handlers) - 1
}

func (e *LinkLabelLinkEvent) Detach(handle int) {
	e.handlers[handle] = nil
}

type LinkLabelLinkEventPublisher struct {
	event LinkLabelLinkEvent
}

func (p *LinkLabelLinkEventPublisher) Event() *LinkLabelLinkEvent {
	return &p.event
}

func (p *LinkLabelLinkEventPublisher) Publish(link *LinkLabelLink) {
	for _, handler := range p.event.handlers {
		if handler != nil {
			handler(link)
		}
	}
}

type LinkLabelLink struct {
	ll    *LinkLabel
	index int
	id    string
	url   string
}

func (lll *LinkLabelLink) Index() int {
	return lll.index
}

func (lll *LinkLabelLink) Id() string {
	return lll.id
}

func (lll *LinkLabelLink) SetId(id string) error {
	old := lll.id

	lll.id = id

	if err := lll.update(); err != nil {
		lll.id = old
		return err
	}

	return nil
}

func (lll *LinkLabelLink) URL() string {
	return lll.url
}

func (lll *LinkLabelLink) SetURL(url string) error {
	old := lll.url

	lll.url = url

	if err := lll.update(); err != nil {
		lll.url = old
		return err
	}

	return nil
}

func (lll *LinkLabelLink) update() error {
	li := win.LITEM{
		ILink: int32(lll.index),
		Mask:  win.LIF_ITEMID | win.LIF_ITEMINDEX | win.LIF_URL,
	}

	id := syscall.StringToUTF16(lll.id)
	url := syscall.StringToUTF16(lll.url)
	copy(li.SzID[:], id[:mini(len(id), win.MAX_LINKID_TEXT)])
	copy(li.SzUrl[:], url[:mini(len(url), win.L_MAX_URL_LENGTH)])

	if win.TRUE != lll.ll.SendMessage(win.LM_SETITEM, 0, uintptr(unsafe.Pointer(&li))) {
		return newErr("LM_SETITEM")
	}

	return nil
}
