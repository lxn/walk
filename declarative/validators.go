// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ValidatorRef struct {
	Validator walk.Validator
}

func (vr ValidatorRef) Create() (walk.Validator, error) {
	return vr.Validator, nil
}

type Range struct {
	Min float64
	Max float64
}

func (r Range) Create() (walk.Validator, error) {
	return walk.NewRangeValidator(r.Min, r.Max)
}

type Regexp struct {
	Pattern string
}

func (re Regexp) Create() (walk.Validator, error) {
	return walk.NewRegexpValidator(re.Pattern)
}
