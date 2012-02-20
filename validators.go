// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"strconv"
)

type ValidationStatus uint

const (
	Invalid ValidationStatus = iota
	Partial
	Valid
)

type Validator interface {
	Validate(s string) ValidationStatus
}

type NumberValidator struct {
	decimals int
	minValue float64
	maxValue float64
}

func NewNumberValidator() *NumberValidator {
	return &NumberValidator{}
}

func (nv *NumberValidator) Validate(s string) ValidationStatus {
	num, err := strconv.Atof64(s)
	if err != nil {
		return Invalid
	}

	if num < nv.minValue {
		return Partial
	}

	if num > nv.maxValue {
		return Invalid
	}

	str := strconv.Ftoa64(num, 'f', nv.decimals)

	if s != str {
		return Invalid
	}

	return Valid
}

func (nv *NumberValidator) Decimals() int {
	return nv.decimals
}

func (nv *NumberValidator) SetDecimals(value int) os.Error {
	if value < 0 {
		return newError("invalid value")
	}

	nv.decimals = value

	return nil
}

func (nv *NumberValidator) MinValue() float64 {
	return nv.minValue
}

func (nv *NumberValidator) MaxValue() float64 {
	return nv.maxValue
}

func (nv *NumberValidator) SetRange(min, max float64) os.Error {
	if min > max {
		return newError("invalid range")
	}

	nv.minValue = min
	nv.maxValue = max

	return nil
}
