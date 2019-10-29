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

package core

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"unsafe"
)

func TestNumber(t *testing.T) {
	iNeg := NewInt64Number(-42)
	iZero := NewInt64Number(0)
	iPos := NewInt64Number(42)
	i64Numbers := [3]Number{iNeg, iZero, iPos}

	for idx, i := range []int64{-42, 0, 42} {
		n := i64Numbers[idx]
		if got := n.AsInt64(); got != i {
			t.Errorf("Number %#v (%s) int64 check failed, expected %d, got %d", n, n.Emit(Int64NumberKind), i, got)
		}
	}

	for _, n := range i64Numbers {
		expected := unsafe.Pointer(&n)
		got := unsafe.Pointer(n.AsRawPtr())
		if expected != got {
			t.Errorf("Getting raw pointer failed, got %v, expected %v", got, expected)
		}
	}

	fNeg := NewFloat64Number(-42.)
	fZero := NewFloat64Number(0.)
	fPos := NewFloat64Number(42.)
	f64Numbers := [3]Number{fNeg, fZero, fPos}

	for idx, f := range []float64{-42., 0., 42.} {
		n := f64Numbers[idx]
		if got := n.AsFloat64(); got != f {
			t.Errorf("Number %#v (%s) float64 check failed, expected %f, got %f", n, n.Emit(Int64NumberKind), f, got)
		}
	}

	for _, n := range f64Numbers {
		expected := unsafe.Pointer(&n)
		got := unsafe.Pointer(n.AsRawPtr())
		if expected != got {
			t.Errorf("Getting raw pointer failed, got %v, expected %v", got, expected)
		}
	}

	cmpsForNeg := [3]int{0, -1, -1}
	cmpsForZero := [3]int{1, 0, -1}
	cmpsForPos := [3]int{1, 1, 0}

	type testcase struct {
		n    Number
		kind NumberKind
		pos  bool
		zero bool
		neg  bool
		nums [3]Number
		cmps [3]int
	}
	testcases := []testcase{
		{
			n:    iNeg,
			kind: Int64NumberKind,
			pos:  false,
			zero: false,
			neg:  true,
			nums: i64Numbers,
			cmps: cmpsForNeg,
		},
		{
			n:    iZero,
			kind: Int64NumberKind,
			pos:  false,
			zero: true,
			neg:  false,
			nums: i64Numbers,
			cmps: cmpsForZero,
		},
		{
			n:    iPos,
			kind: Int64NumberKind,
			pos:  true,
			zero: false,
			neg:  false,
			nums: i64Numbers,
			cmps: cmpsForPos,
		},
		{
			n:    fNeg,
			kind: Float64NumberKind,
			pos:  false,
			zero: false,
			neg:  true,
			nums: f64Numbers,
			cmps: cmpsForNeg,
		},
		{
			n:    fZero,
			kind: Float64NumberKind,
			pos:  false,
			zero: true,
			neg:  false,
			nums: f64Numbers,
			cmps: cmpsForZero,
		},
		{
			n:    fPos,
			kind: Float64NumberKind,
			pos:  true,
			zero: false,
			neg:  false,
			nums: f64Numbers,
			cmps: cmpsForPos,
		},
	}
	for _, tt := range testcases {
		if got := tt.n.IsPositive(tt.kind); got != tt.pos {
			t.Errorf("Number %#v (%s) positive check failed, expected %v, got %v", tt.n, tt.n.Emit(tt.kind), tt.pos, got)
		}
		if got := tt.n.IsZero(tt.kind); got != tt.zero {
			t.Errorf("Number %#v (%s) zero check failed, expected %v, got %v", tt.n, tt.n.Emit(tt.kind), tt.pos, got)
		}
		if got := tt.n.IsNegative(tt.kind); got != tt.neg {
			t.Errorf("Number %#v (%s) negative check failed, expected %v, got %v", tt.n, tt.n.Emit(tt.kind), tt.pos, got)
		}
		for i := 0; i < 3; i++ {
			if got := tt.n.CompareRaw(tt.kind, tt.nums[i].AsRaw()); got != tt.cmps[i] {
				t.Errorf("Number %#v (%s) compare check with %#v (%s) failed, expected %d, got %d", tt.n, tt.n.Emit(tt.kind), tt.nums[i], tt.nums[i].Emit(tt.kind), tt.cmps[i], got)
			}
		}
	}
}

type testFormatCase struct {
	num  Number
	kind NumberKind
	want string
}

func makeInt64Case(val int64) testFormatCase {
	return testFormatCase{
		num:  NewInt64Number(val),
		kind: Int64NumberKind,
		want: fmt.Sprintf("%d", val),
	}
}

func makeFloat64Case(val float64) testFormatCase {
	return testFormatCase{
		num:  NewFloat64Number(val),
		kind: Float64NumberKind,
		want: fmt.Sprintf("%.*g", OutputFloatingPointPrecision, val),
	}
}

var testFormatData = []testFormatCase{
	makeInt64Case(100),
	makeInt64Case(-100),
	makeInt64Case(0),
	makeInt64Case(math.MaxInt8),
	makeInt64Case(math.MaxInt16),
	makeInt64Case(math.MaxInt32),
	makeInt64Case(math.MaxInt64),
	makeInt64Case(math.MinInt8),
	makeInt64Case(math.MinInt16),
	makeInt64Case(math.MinInt32),
	makeInt64Case(math.MinInt64),
	makeFloat64Case(77),
	makeFloat64Case(-77e-77),
	makeFloat64Case(math.MaxFloat64),
	makeFloat64Case(math.MaxFloat32),
	makeFloat64Case(math.SmallestNonzeroFloat32),
	makeFloat64Case(math.SmallestNonzeroFloat64),
	makeFloat64Case(0.123456789123456789123456789123456789),
}

func TestNumberEmit(t *testing.T) {
	for _, data := range testFormatData {
		have := data.num.Emit(data.kind)
		if have != data.want {
			t.Errorf("Invalid Emit() - got %s want %s", have, data.want)
		}
	}
}

func TestNumberEncode(t *testing.T) {
	for _, data := range testFormatData {
		var sb strings.Builder
		var tmp [32]byte
		_, _ = data.num.Encode(data.kind, &sb, tmp[:])
		have := sb.String()
		if have != data.want {
			t.Errorf("Invalid Encode() - got %s want %s", have, data.want)
		}
	}
}
