// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errValidationFailed = errors.New("validation failed")
)

type DataBindable interface {
	Widget
	BindingMember() string
	SetBindingMember(member string) error
	BindingValue() interface{}
	SetBindingValue(value interface{}) error
	BindingValueChanged() *Event
}

type Validatable interface {
	Validator() Validator
	SetValidator(v Validator)
}

type ErrorPresenter interface {
	PresentError(err error, widget Widget)
}

type DataBinder struct {
	dataSource                interface{}
	boundWidgets              []DataBindable
	widget2ChangedHandle      map[DataBindable]int
	invalidValueWidgets       map[DataBindable]bool
	errorPresenter            ErrorPresenter
	canSubmitChangedPublisher EventPublisher
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
	for widget, handle := range db.widget2ChangedHandle {
		widget.BindingValueChanged().Detach(handle)
	}

	db.boundWidgets = boundWidgets

	db.widget2ChangedHandle = make(map[DataBindable]int)
	db.invalidValueWidgets = make(map[DataBindable]bool)

	for _, widget := range boundWidgets {
		widget := widget

		changedEvent := widget.BindingValueChanged()
		if changedEvent == nil {
			continue
		}

		db.widget2ChangedHandle[widget] = changedEvent.Attach(func() {
			db.validateWidget(widget)
		})
	}
}

func (db *DataBinder) validateWidget(widget DataBindable) {
	validatable, ok := widget.(Validatable)
	if !ok {
		return
	}

	validator := validatable.Validator()
	if validator == nil {
		return
	}

	if err := validator.Validate(widget.BindingValue()); err != nil {
		if db.invalidValueWidgets[widget] {
			return
		}

		db.invalidValueWidgets[widget] = true

		if len(db.invalidValueWidgets) == 1 {
			db.canSubmitChangedPublisher.Publish()
		}

		if db.errorPresenter != nil {
			db.errorPresenter.PresentError(err, widget)
		}
	} else {
		if !db.invalidValueWidgets[widget] {
			return
		}

		delete(db.invalidValueWidgets, widget)

		if len(db.invalidValueWidgets) == 0 {
			db.canSubmitChangedPublisher.Publish()
		}

		if db.errorPresenter != nil {
			db.errorPresenter.PresentError(nil, widget)
		}
	}
}

func (db *DataBinder) ErrorPresenter() ErrorPresenter {
	return db.errorPresenter
}

func (db *DataBinder) SetErrorPresenter(ep ErrorPresenter) {
	db.errorPresenter = ep
}

func (db *DataBinder) CanSubmit() bool {
	return len(db.invalidValueWidgets) == 0
}

func (db *DataBinder) CanSubmitChanged() *Event {
	return db.canSubmitChangedPublisher.Event()
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

			if err := widget.SetBindingValue(f64); err != nil {
				return err
			}
		} else {
			if err := widget.SetBindingValue(field.Interface()); err != nil {
				return err
			}
		}

		db.validateWidget(widget)
		return nil
	})
}

func (db *DataBinder) Submit() error {
	if !db.CanSubmit() {
		return errValidationFailed
	}

	return db.forEach(func(widget DataBindable, field reflect.Value) error {
		value := widget.BindingValue()
		if value == nil {
			// This happens e.g. if CurrentIndex() of a ComboBox returns -1.
			// FIXME: Should we handle this differently?
			return nil
		}

		if f64, ok := value.(float64); ok {
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
