// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"fmt"
	"reflect"
)

type DataBindable interface {
	BindingMember() string
	SetBindingMember(member string) error
	BindingValue() interface{}
	SetBindingValue(value interface{}) error
	BindingValueChanged() *Event
}

type DataBinder struct {
	dataSource   interface{}
	boundWidgets []DataBindable
}

func NewDataBinder() *DataBinder {
	return new(DataBinder)
}

func (db *DataBinder) DataSource() interface{} {
	return db.dataSource
}

func (db *DataBinder) SetDataSource(dataSource interface{}) {
	db.dataSource = dataSource
}

func (db *DataBinder) BoundWidgets() []DataBindable {
	return db.boundWidgets
}

func (db *DataBinder) SetBoundWidgets(boundWidgets []DataBindable) {
	db.boundWidgets = boundWidgets
}

func (db *DataBinder) Reset() error {
	return db.forEach(func(widget DataBindable, field reflect.Value) error {
		return widget.SetBindingValue(field.Interface())
	})
}

func (db *DataBinder) Submit() error {
	return db.forEach(func(widget DataBindable, field reflect.Value) error {
		field.Set(reflect.ValueOf(widget.BindingValue()))

		return nil
	})
}

func (db *DataBinder) forEach(f func(widget DataBindable, field reflect.Value) error) error {
	p := reflect.ValueOf(db.dataSource)
	if p.Type().Kind() != reflect.Ptr {
		return newError("DataSource must be a pointer to a struct.")
	}

	if p.IsNil() {
		return nil
	}

	s := reflect.Indirect(p)
	if s.Type().Kind() != reflect.Struct {
		return newError("DataSource must be a pointer to a struct.")
	}

	for _, widget := range db.boundWidgets {
		if field := s.FieldByName(widget.BindingMember()); field.IsValid() {
			if err := f(widget, field); err != nil {
				return err
			}
		} else {
			return newError(fmt.Sprintf("Field '%s' not found in struct '%s'.", widget.BindingMember(), s.Type().Name()))
		}
	}

	return nil
}

func validateBindingMemberSyntax(member string) error {
	// FIXME
	return nil
}
