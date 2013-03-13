// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

type Property struct {
	name                  string
	get                   func() interface{}
	set                   func(v interface{}) error
	validator             Validator
	source                interface{}
	sourceChangedHandle   int
	changedEventPublisher *EventPublisher
	customChangedEvent    *Event
}

func NewProperty(name string, get func() interface{}, set func(v interface{}) error, customChangedEvent *Event) *Property {
	p := &Property{name: name, get: get, set: set}

	if customChangedEvent != nil {
		p.customChangedEvent = customChangedEvent
	} else {
		p.changedEventPublisher = new(EventPublisher)
	}

	return p
}

func (p *Property) Name() string {
	return p.name
}

func (p *Property) Get() interface{} {
	return p.get()
}

func (p *Property) Set(v interface{}) error {
	p.assertNotReadOnly()

	// FIXME: Ugly special case for Visible property
	if p.name != "Visible" && v == p.Get() {
		return nil
	}

	if err := p.set(v); err != nil {
		return err
	}

	if p.customChangedEvent == nil {
		p.changedEventPublisher.Publish()
	}

	return nil
}

func (p *Property) Validator() Validator {
	return p.validator
}

func (p *Property) SetValidator(v Validator) {
	p.validator = v
}

func (p *Property) Source() interface{} {
	return p.source
}

func (p *Property) SetSource(source interface{}) {
	switch source := source.(type) {
	case *Property:
		if source != nil {
			p.assertNotReadOnly()
		}

		for cur := source; cur != nil; cur, _ = cur.source.(*Property) {
			if cur == p {
				panic("source cycle")
			}
		}

		if source != nil {
			p.Set(source.Get())

			p.sourceChangedHandle = source.Changed().Attach(func() {
				p.Set(source.Get())
			})
		}

	case string:
		// nop

	default:
		panic("invalid source type")
	}

	if oldProp, ok := p.source.(*Property); ok {
		oldProp.Changed().Detach(p.sourceChangedHandle)
	}

	p.source = source
}

func (p *Property) Changed() *Event {
	if p.customChangedEvent != nil {
		return p.customChangedEvent
	}

	return p.changedEventPublisher.Event()
}

func (p *Property) ReadOnly() bool {
	return p.set == nil
}

func (p *Property) assertNotReadOnly() {
	if p.ReadOnly() {
		panic("property is read-only")
	}
}
