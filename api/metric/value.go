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
)

type MeasurementValue struct {
	i interface{}
}

type MeasurementValueVisitor interface {
	Int(int64)
	Float(float64)
}

func NewIntMeasurementValue(i int64) MeasurementValue {
	return MeasurementValue{
		i: i,
	}
}

func NewFloatMeasurementValue(f float64) MeasurementValue {
	return MeasurementValue{
		i: f,
	}
}

func (v MeasurementValue) GetInt() (int64, bool) {
	i, ok := v.i.(int64)
	return i, ok
}

func (v MeasurementValue) GetFloat() (float64, bool) {
	f, ok := v.i.(float64)
	return f, ok
}

func (v MeasurementValue) Visit(visitor MeasurementValueVisitor) {
	switch rv := v.i.(type) {
	case int64:
		visitor.Int(rv)
	case float64:
		visitor.Float(rv)
	default:
		panic(fmt.Sprintf("invalid measurement value: %#v", rv))
	}
}

func (v MeasurementValue) Emit() string {
	return fmt.Sprintf("%v", v.i)
}

type intCoercer struct {
	i *int64
}

var _ MeasurementValueVisitor = intCoercer{}

func (c intCoercer) Int(i int64) {
	*c.i = i
}

func (c intCoercer) Float(f float64) {
	*c.i = int64(f)
}

func (v MeasurementValue) CoerceToInt() int64 {
	var i int64
	c := intCoercer{
		i: &i,
	}
	v.Visit(c)
	return i
}

type floatCoercer struct {
	f *float64
}

var _ MeasurementValueVisitor = floatCoercer{}

func (c floatCoercer) Int(i int64) {
	*c.f = float64(i)
}

func (c floatCoercer) Float(f float64) {
	*c.f = f
}

func (v MeasurementValue) CoerceToFloat() float64 {
	var f float64
	c := floatCoercer{
		f: &f,
	}
	v.Visit(c)
	return f
}
