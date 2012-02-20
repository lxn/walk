// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type actionChangedHandler interface {
	onActionChanged(action *Action) (err error)
}

var (
	// ISSUE: When pressing enter resp. escape,
	// WM_COMMAND with wParam=1 resp. 2 is sent.
	// Maybe there is more to consider.
	nextActionId uint16             = 3
	actionsById  map[uint16]*Action = make(map[uint16]*Action)
)

type Action struct {
	menu               *Menu
	triggeredPublisher EventPublisher
	changedHandlers    []actionChangedHandler
	text               string
	toolTip            string
	image              *Bitmap
	enabled            bool
	visible            bool
	checkable          bool
	checked            bool
	exclusive          bool
	id                 uint16
}

func NewAction() *Action {
	a := &Action{
		enabled: true,
		id:      nextActionId,
		visible: true,
	}

	actionsById[a.id] = a

	nextActionId++

	return a
}

func (a *Action) Checkable() bool {
	return a.checkable
}

func (a *Action) SetCheckable(value bool) (err error) {
	if value != a.checkable {
		old := a.checkable

		a.checkable = value

		if err = a.raiseChanged(); err != nil {
			a.checkable = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Checked() bool {
	return a.checked
}

func (a *Action) SetChecked(value bool) (err error) {
	if value != a.checked {
		old := a.checked

		a.checked = value

		if err = a.raiseChanged(); err != nil {
			a.checked = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Enabled() bool {
	return a.enabled
}

func (a *Action) SetEnabled(value bool) (err error) {
	if value != a.enabled {
		old := a.enabled

		a.enabled = value

		if err = a.raiseChanged(); err != nil {
			a.enabled = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Exclusive() bool {
	return a.exclusive
}

func (a *Action) SetExclusive(value bool) (err error) {
	if value != a.exclusive {
		old := a.exclusive

		a.exclusive = value

		if err = a.raiseChanged(); err != nil {
			a.exclusive = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Image() *Bitmap {
	return a.image
}

func (a *Action) SetImage(value *Bitmap) (err error) {
	if value != a.image {
		old := a.image

		a.image = value

		if err = a.raiseChanged(); err != nil {
			a.image = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Text() string {
	return a.text
}

func (a *Action) SetText(value string) (err error) {
	if value != a.text {
		old := a.text

		a.text = value

		if err = a.raiseChanged(); err != nil {
			a.text = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) ToolTip() string {
	return a.toolTip
}

func (a *Action) SetToolTip(value string) (err error) {
	if value != a.toolTip {
		old := a.toolTip

		a.toolTip = value

		if err = a.raiseChanged(); err != nil {
			a.toolTip = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Visible() bool {
	return a.visible
}

func (a *Action) SetVisible(value bool) (err error) {
	if value != a.visible {
		old := a.visible

		a.visible = value

		if err = a.raiseChanged(); err != nil {
			a.visible = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Triggered() *Event {
	return a.triggeredPublisher.Event()
}

func (a *Action) raiseTriggered() {
	a.triggeredPublisher.Publish()
}

func (a *Action) addChangedHandler(handler actionChangedHandler) {
	a.changedHandlers = append(a.changedHandlers, handler)
}

func (a *Action) removeChangedHandler(handler actionChangedHandler) {
	for i, h := range a.changedHandlers {
		if h == handler {
			a.changedHandlers = append(a.changedHandlers[:i], a.changedHandlers[i+1:]...)
			break
		}
	}
}

func (a *Action) raiseChanged() (err error) {
	for _, handler := range a.changedHandlers {
		if err = handler.onActionChanged(a); err != nil {
			return
		}
	}

	return
}
