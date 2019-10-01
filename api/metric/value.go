// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"fmt"
	"math"
)

type valueKind int8

const (
	int64Kind valueKind = iota
	float64Kind
)

// MeasurementValue represents either an integral or a floating point
// value of a measurement. Empty MeasurementValue is treated as
// integral zero.
type MeasurementValue struct {
	raw  uint64
	kind valueKind
}

// MeasurementValueVisitor is an interface used for visitation of
// MeasurementValues.
type MeasurementValueVisitor interface {
	// Int64 gets called for an integral MeasurementValue.
	Int64(int64)
	// Int64 gets called for a floating point MeasurementValue.
	Float64(float64)
	// Unknown gets called on the unknown value. This normally
	// doesn't happen - it is meant for the future compatilibity
	// when new types may get added to the MeasurementValue. The
	// additional types would be handled an in extension interface
	// of the MeasurementValueVisitor.
	Unknown(interface{})
}

// NewInt64MeasurementValue creates an integral MeasurementValue.
func NewInt64MeasurementValue(i int64) MeasurementValue {
	return MeasurementValue{
		raw:  uint64(i),
		kind: int64Kind,
	}
}

// NewFloat64MeasurementValue creates a floating point
// MeasurementValue.
func NewFloat64MeasurementValue(f float64) MeasurementValue {
	return MeasurementValue{
		raw:  math.Float64bits(f),
		kind: float64Kind,
	}
}

// GetInt64 tries to get an integral value from MeasurementValue. The
// use of this function is similar to doing a checked type assertion
// (i, ok := v.(int64)).
func (v MeasurementValue) GetInt64() (int64, bool) {
	if v.kind != int64Kind {
		return 0, false
	}
	return v.int64(), true
}

// GetFloat64 tries to get a floating point value from
// MeasurementValue. The use of this function is similar to doing a
// checked type assertion (i, ok := v.(float64)).
func (v MeasurementValue) GetFloat64() (float64, bool) {
	if v.kind != float64Kind {
		return 0, false
	}
	return v.float64(), true
}

func (v MeasurementValue) int64() int64 {
	return int64(v.raw)
}

func (v MeasurementValue) float64() float64 {
	return math.Float64frombits(v.raw)
}

// Visit uses the visitor to find out the actual value of the
// MeasurementValue.
func (v MeasurementValue) Visit(visitor MeasurementValueVisitor) {
	switch v.kind {
	case int64Kind:
		visitor.Int64(v.int64())
	case float64Kind:
		visitor.Float64(v.float64())
	}
}

type stringEmitter struct {
	s *string
}

var _ MeasurementValueVisitor = stringEmitter{}

func (e stringEmitter) Int64(i int64) {
	*e.s = fmt.Sprintf("%d", i)
}

func (e stringEmitter) Float64(f float64) {
	*e.s = fmt.Sprintf("%f", f)
}

func (e stringEmitter) Unknown(interface{}) {
}

// Emit returns a string representation of the actual value of the
// MeasurementValue. A %d is used for integral values, %f for floating
// point values.
func (v MeasurementValue) Emit() string {
	var s string
	c := stringEmitter{
		s: &s,
	}
	v.Visit(c)
	return s
}

type int64Coercer struct {
	i *int64
}

var _ MeasurementValueVisitor = int64Coercer{}

func (c int64Coercer) Int64(i int64) {
	*c.i = i
}

func (c int64Coercer) Float64(f float64) {
	*c.i = int64(f)
}

func (c int64Coercer) Unknown(interface{}) {
}

// CoerceToInt64 returns the value of the MeasurementValue casted to
// int64. This may result in loss of precision of the actual value was
// a floating point one.
func (v MeasurementValue) CoerceToInt64() int64 {
	var i int64
	c := int64Coercer{
		i: &i,
	}
	v.Visit(c)
	return i
}

type float64Coercer struct {
	f *float64
}

var _ MeasurementValueVisitor = float64Coercer{}

func (c float64Coercer) Int64(i int64) {
	*c.f = float64(i)
}

func (c float64Coercer) Float64(f float64) {
	*c.f = f
}

func (c float64Coercer) Unknown(interface{}) {
}

// CoerceToFloat64 returns the value of the MeasurementValue casted to
// float64. This may result in loss of precision of the actual value
// was an integral one.
func (v MeasurementValue) CoerceToFloat64() float64 {
	var f float64
	c := float64Coercer{
		f: &f,
	}
	v.Visit(c)
	return f
}
