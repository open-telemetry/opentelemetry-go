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
	"testing"
	"unsafe"
)

func TestMeasurementValue(t *testing.T) {
	iNeg := NewInt64MeasurementValue(-42)
	iZero := NewInt64MeasurementValue(0)
	iPos := NewInt64MeasurementValue(42)
	i64Values := [3]MeasurementValue{iNeg, iZero, iPos}

	for idx, i := range []int64{-42, 0, 42} {
		v := i64Values[idx]
		if got := v.AsInt64(); got != i {
			t.Errorf("Value %#v (%s) int64 check failed, expected %d, got %d", v, v.Emit(Int64ValueKind), i, got)
		}
	}

	for _, v := range i64Values {
		expected := unsafe.Pointer(&v)
		got := unsafe.Pointer(v.AsRawPtr())
		if expected != got {
			t.Errorf("Getting raw pointer failed, got %v, expected %v", got, expected)
		}
	}

	fNeg := NewFloat64MeasurementValue(-42.)
	fZero := NewFloat64MeasurementValue(0.)
	fPos := NewFloat64MeasurementValue(42.)
	f64Values := [3]MeasurementValue{fNeg, fZero, fPos}

	for idx, f := range []float64{-42., 0., 42.} {
		v := f64Values[idx]
		if got := v.AsFloat64(); got != f {
			t.Errorf("Value %#v (%s) float64 check failed, expected %f, got %f", v, v.Emit(Int64ValueKind), f, got)
		}
	}

	for _, v := range f64Values {
		expected := unsafe.Pointer(&v)
		got := unsafe.Pointer(v.AsRawPtr())
		if expected != got {
			t.Errorf("Getting raw pointer failed, got %v, expected %v", got, expected)
		}
	}

	cmpsForNeg := [3]int{0, -1, -1}
	cmpsForZero := [3]int{1, 0, -1}
	cmpsForPos := [3]int{1, 1, 0}

	type testcase struct {
		v    MeasurementValue
		kind ValueKind
		pos  bool
		zero bool
		neg  bool
		vals [3]MeasurementValue
		cmps [3]int
	}
	testcases := []testcase{
		{
			v:    iNeg,
			kind: Int64ValueKind,
			pos:  false,
			zero: false,
			neg:  true,
			vals: i64Values,
			cmps: cmpsForNeg,
		},
		{
			v:    iZero,
			kind: Int64ValueKind,
			pos:  false,
			zero: true,
			neg:  false,
			vals: i64Values,
			cmps: cmpsForZero,
		},
		{
			v:    iPos,
			kind: Int64ValueKind,
			pos:  true,
			zero: false,
			neg:  false,
			vals: i64Values,
			cmps: cmpsForPos,
		},
		{
			v:    fNeg,
			kind: Float64ValueKind,
			pos:  false,
			zero: false,
			neg:  true,
			vals: f64Values,
			cmps: cmpsForNeg,
		},
		{
			v:    fZero,
			kind: Float64ValueKind,
			pos:  false,
			zero: true,
			neg:  false,
			vals: f64Values,
			cmps: cmpsForZero,
		},
		{
			v:    fPos,
			kind: Float64ValueKind,
			pos:  true,
			zero: false,
			neg:  false,
			vals: f64Values,
			cmps: cmpsForPos,
		},
	}
	for _, tt := range testcases {
		if got := tt.v.IsPositive(tt.kind); got != tt.pos {
			t.Errorf("Value %#v (%s) positive check failed, expected %v, got %v", tt.v, tt.v.Emit(tt.kind), tt.pos, got)
		}
		if got := tt.v.IsZero(tt.kind); got != tt.zero {
			t.Errorf("Value %#v (%s) zero check failed, expected %v, got %v", tt.v, tt.v.Emit(tt.kind), tt.pos, got)
		}
		if got := tt.v.IsNegative(tt.kind); got != tt.neg {
			t.Errorf("Value %#v (%s) negative check failed, expected %v, got %v", tt.v, tt.v.Emit(tt.kind), tt.pos, got)
		}
		for i := 0; i < 3; i++ {
			if got := tt.v.RawCompare(tt.vals[i].AsRaw(), tt.kind); got != tt.cmps[i] {
				t.Errorf("Value %#v (%s) compare check with %#v (%s) failed, expected %d, got %d", tt.v, tt.v.Emit(tt.kind), tt.vals[i], tt.vals[i].Emit(tt.kind), tt.cmps[i], got)
			}
		}
	}
}
