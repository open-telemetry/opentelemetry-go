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

package sum // import "go.opentelemetry.io/otel/sdk/metric/aggregator/sum"

import (
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

type (
	Monotonicity interface {
		kind() aggregation.Kind
		category() aggregation.Category
	}

	Monotonic    struct{}
	NonMonotonic struct{}

	Methods[N number.Any, Traits number.Traits[N], M Monotonicity, Storage State[N, Traits, M]] struct{}

	State[N number.Any, Traits number.Traits[N], M Monotonicity] struct {
		value N
	}

	MonotonicInt64    = State[int64, number.Int64Traits, Monotonic]
	NonMonotonicInt64 = State[int64, number.Int64Traits, NonMonotonic]

	MonotonicFloat64    = State[float64, number.Float64Traits, Monotonic]
	NonMonotonicFloat64 = State[float64, number.Float64Traits, NonMonotonic]
)

func NewMonotonicInt64(x int64) *MonotonicInt64 {
	return &MonotonicInt64{value: x}
}

func NewNonMonotonicInt64(x int64) *NonMonotonicInt64 {
	return &NonMonotonicInt64{value: x}
}

func NewMonotonicFloat64(x float64) *MonotonicFloat64 {
	return &MonotonicFloat64{value: x}
}

func NewNonMonotonicFloat64(x float64) *NonMonotonicFloat64 {
	return &NonMonotonicFloat64{value: x}
}

func (Monotonic) kind() aggregation.Kind {
	return aggregation.MonotonicSumKind
}

func (NonMonotonic) kind() aggregation.Kind {
	return aggregation.NonMonotonicSumKind
}

func (Monotonic) category() aggregation.Category {
	return aggregation.MonotonicSumCategory
}

func (NonMonotonic) category() aggregation.Category {
	return aggregation.NonMonotonicSumCategory
}

var (
	_ aggregator.Methods[int64, MonotonicInt64]        = Methods[int64, number.Int64Traits, Monotonic, MonotonicInt64]{}
	_ aggregator.Methods[float64, MonotonicFloat64]    = Methods[float64, number.Float64Traits, Monotonic, MonotonicFloat64]{}
	_ aggregator.Methods[int64, NonMonotonicInt64]     = Methods[int64, number.Int64Traits, NonMonotonic, NonMonotonicInt64]{}
	_ aggregator.Methods[float64, NonMonotonicFloat64] = Methods[float64, number.Float64Traits, NonMonotonic, NonMonotonicFloat64]{}

	_ aggregation.Sum = &MonotonicInt64{}
	_ aggregation.Sum = &MonotonicFloat64{}
	_ aggregation.Sum = &NonMonotonicInt64{}
	_ aggregation.Sum = &NonMonotonicFloat64{}
)

func (s *State[N, Traits, M]) Sum() number.Number {
	var t Traits
	return t.ToNumber(s.value)
}

func (lv *State[N, Traits, M]) Kind() aggregation.Kind {
	var m M
	return m.kind()
}

func (lv *State[N, Traits, M]) IsMonotonic() bool {
	var m M
	return m.kind() == aggregation.MonotonicSumKind
}

func (Methods[N, Traits, M, Storage]) Kind() aggregation.Kind {
	var m M
	return m.kind()
}

func (Methods[N, Traits, M, Storage]) Init(state *State[N, Traits, M], _ aggregator.Config) {
	// Note: storage is zero to start
}

func (Methods[N, Traits, M, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits, M]) {
	var t Traits
	dest.value = t.SwapAtomic(&resetSrc.value, 0)
}

func (Methods[N, Traits, M, Storage]) Reset(ptr *State[N, Traits, M]) {
	ptr.value = 0
}

func (Methods[N, Traits, M, Storage]) HasChange(ptr *State[N, Traits, M]) bool {
	return ptr.value != 0
}

func (Methods[N, Traits, M, Storage]) Update(state *State[N, Traits, M], value N) {
	var m M
	if !aggregator.RangeTest[N, Traits](value, m.category()) {
		return
	}

	var t Traits
	t.AddAtomic(&state.value, value)
}

func (Methods[N, Traits, M, Storage]) Merge(to, from *State[N, Traits, M]) {
	to.value += from.value
}

func (Methods[N, Traits, M, Storage]) ToAggregation(state *State[N, Traits, M]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, M, Storage]) ToStorage(aggr aggregation.Aggregation) (*State[N, Traits, M], bool) {
	r, ok := aggr.(*State[N, Traits, M])
	return r, ok
}

func (Methods[N, Traits, M, Storage]) SubtractSwap(newvalue, oldvalueModified *State[N, Traits, M]) {
	oldvalueModified.value = newvalue.value - oldvalueModified.value
}
