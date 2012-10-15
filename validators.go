// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	errNumberOutOfRange  = errors.New("The number is out of the allowed range.")
	errPatternNotMatched = errors.New("The text does not match the required pattern.")
	errSelectionRequired = errors.New("A selection is required.")
)

type Validator interface {
	Validate(v interface{}) error
}

type NumberValidator struct {
	min float64
	max float64
}

func NewNumberValidator(min, max float64) (*NumberValidator, error) {
	if max <= min {
		return nil, errors.New("max <= min")
	}

	return &NumberValidator{min: min, max: max}, nil
}

func (nv *NumberValidator) Min() float64 {
	return nv.min
}

func (nv *NumberValidator) Max() float64 {
	return nv.max
}

func (nv *NumberValidator) Validate(v interface{}) error {
	f64 := v.(float64)

	if f64 < nv.min || f64 > nv.max {
		return errNumberOutOfRange
	}

	return nil
}

type RegexpValidator struct {
	re *regexp.Regexp
}

func NewRegexpValidator(pattern string) (*RegexpValidator, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RegexpValidator{re}, nil
}

func (rv *RegexpValidator) Pattern() string {
	return rv.re.String()
}

func (rv *RegexpValidator) Validate(v interface{}) error {
	var matched bool

	switch val := v.(type) {
	case string:
		matched = rv.re.MatchString(val)

	case []byte:
		matched = rv.re.Match(val)

	case fmt.Stringer:
		matched = rv.re.MatchString(val.String())

	default:
		panic("Unsupported type")
	}

	if !matched {
		return errPatternNotMatched
	}

	return nil
}

type selectionRequiredValidator struct {
}

var selectionRequiredValidatorSingleton Validator = selectionRequiredValidator{}

func SelectionRequiredValidator() Validator {
	return selectionRequiredValidatorSingleton
}

func (selectionRequiredValidator) Validate(v interface{}) error {
	if v == nil {
		// For Widgets like ComboBox nil is passed to indicate "no selection".
		return errSelectionRequired
	}

	return nil
}
