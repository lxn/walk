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
		if f64, ok := widget.BindingValue().(float64); ok {
			switch v := field.Interface().(type) {
			case float32:
				f64 = float64(v)

			case float64:
				f64 = v

			case int:
				f64 = float64(v)

			case int8:
				f64 = float64(v)

			case int16:
				f64 = float64(v)

			case int32:
				f64 = float64(v)

			case int64:
				f64 = float64(v)

			case uint:
				f64 = float64(v)

			case uint8:
				f64 = float64(v)

			case uint16:
				f64 = float64(v)

			case uint32:
				f64 = float64(v)

			case uint64:
				f64 = float64(v)

			case uintptr:
				f64 = float64(v)

			default:
				return newError(fmt.Sprintf("Field '%s': Can't convert %s to float64.", widget.BindingMember(), field.Type().Name()))
			}

			return widget.SetBindingValue(f64)
		}

		return widget.SetBindingValue(field.Interface())
	})
}

func (db *DataBinder) Submit() error {
	return db.forEach(func(widget DataBindable, field reflect.Value) error {
		if f64, ok := widget.BindingValue().(float64); ok {
			switch field.Kind() {
			case reflect.Float32, reflect.Float64:
				field.SetFloat(f64)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.SetInt(int64(f64))

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				field.SetUint(uint64(f64))

			default:
				return newError(fmt.Sprintf("Field '%s': Can't convert float64 to %s.", widget.BindingMember(), field.Type().Name()))
			}

			return nil
		}

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
