// Copyright The OpenTelemetry Authors
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

package number

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestNumber(t *testing.T) {
	iNeg := NewInt64Number(-42)
	iZero := NewInt64Number(0)
	iPos := NewInt64Number(42)
	i64Numbers := [3]Number{iNeg, iZero, iPos}

	for idx, i := range []int64{-42, 0, 42} {
		n := i64Numbers[idx]
		if got := n.AsInt64(); got != i {
			t.Errorf("Number %#v (%s) int64 check failed, expected %d, got %d", n, n.Emit(Int64Kind), i, got)
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
			t.Errorf("Number %#v (%s) float64 check failed, expected %f, got %f", n, n.Emit(Int64Kind), f, got)
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
		// n needs to be aligned for 64-bit atomic operations.
		n Number
		// nums needs to be aligned for 64-bit atomic operations.
		nums [3]Number
		kind Kind
		pos  bool
		zero bool
		neg  bool
		cmps [3]int
	}
	testcases := []testcase{
		{
			n:    iNeg,
			kind: Int64Kind,
			pos:  false,
			zero: false,
			neg:  true,
			nums: i64Numbers,
			cmps: cmpsForNeg,
		},
		{
			n:    iZero,
			kind: Int64Kind,
			pos:  false,
			zero: true,
			neg:  false,
			nums: i64Numbers,
			cmps: cmpsForZero,
		},
		{
			n:    iPos,
			kind: Int64Kind,
			pos:  true,
			zero: false,
			neg:  false,
			nums: i64Numbers,
			cmps: cmpsForPos,
		},
		{
			n:    fNeg,
			kind: Float64Kind,
			pos:  false,
			zero: false,
			neg:  true,
			nums: f64Numbers,
			cmps: cmpsForNeg,
		},
		{
			n:    fZero,
			kind: Float64Kind,
			pos:  false,
			zero: true,
			neg:  false,
			nums: f64Numbers,
			cmps: cmpsForZero,
		},
		{
			n:    fPos,
			kind: Float64Kind,
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

func TestNumberZero(t *testing.T) {
	zero := Number(0)
	zerof := NewFloat64Number(0)
	zeroi := NewInt64Number(0)

	if zero != zerof || zero != zeroi {
		t.Errorf("Invalid zero representations")
	}
}

func TestNumberAsInterface(t *testing.T) {
	i64 := NewInt64Number(10)
	f64 := NewFloat64Number(11.11)
	require.Equal(t, int64(10), (&i64).AsInterface(Int64Kind).(int64))
	require.Equal(t, 11.11, (&f64).AsInterface(Float64Kind).(float64))
}

func TestNumberSignChange(t *testing.T) {
	t.Run("Int64", func(t *testing.T) {
		posInt := NewInt64Number(10)
		negInt := NewInt64Number(-10)

		require.Equal(t, posInt, NewNumberSignChange(Int64Kind, negInt))
		require.Equal(t, negInt, NewNumberSignChange(Int64Kind, posInt))
	})

	t.Run("Float64", func(t *testing.T) {
		posFloat := NewFloat64Number(10)
		negFloat := NewFloat64Number(-10)

		require.Equal(t, posFloat, NewNumberSignChange(Float64Kind, negFloat))
		require.Equal(t, negFloat, NewNumberSignChange(Float64Kind, posFloat))
	})

	t.Run("Float64Zero", func(t *testing.T) {
		posFloat := NewFloat64Number(0)
		negFloat := NewFloat64Number(math.Copysign(0, -1))

		require.Equal(t, posFloat, NewNumberSignChange(Float64Kind, negFloat))
		require.Equal(t, negFloat, NewNumberSignChange(Float64Kind, posFloat))
	})

	t.Run("Float64Inf", func(t *testing.T) {
		posFloat := NewFloat64Number(math.Inf(+1))
		negFloat := NewFloat64Number(math.Inf(-1))

		require.Equal(t, posFloat, NewNumberSignChange(Float64Kind, negFloat))
		require.Equal(t, negFloat, NewNumberSignChange(Float64Kind, posFloat))
	})

	t.Run("Float64NaN", func(t *testing.T) {
		posFloat := NewFloat64Number(math.NaN())
		negFloat := NewFloat64Number(math.Copysign(math.NaN(), -1))

		require.Equal(t, posFloat, NewNumberSignChange(Float64Kind, negFloat))
		require.Equal(t, negFloat, NewNumberSignChange(Float64Kind, posFloat))
	})
}
