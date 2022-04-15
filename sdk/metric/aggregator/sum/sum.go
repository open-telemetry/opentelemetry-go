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
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

type (
	Monotonicity interface {
		category() aggregation.Category
	}

	Monotonic    struct{}
	NonMonotonic struct{}

	Methods[N number.Any, Traits traits.Any[N], M Monotonicity, Storage State[N, Traits, M]] struct{}

	State[N number.Any, Traits traits.Any[N], M Monotonicity] struct {
		value N
	}
)

func (Monotonic) category() aggregation.Category {
	return aggregation.MonotonicSumCategory
}

func (NonMonotonic) category() aggregation.Category {
	return aggregation.NonMonotonicSumCategory
}

var (
	_ aggregator.Methods[int64, State[int64, traits.Int64, Monotonic]]       = Methods[int64, traits.Int64, Monotonic, State[int64, traits.Int64, Monotonic]]{}
	_ aggregator.Methods[float64, State[float64, traits.Float64, Monotonic]] = Methods[float64, traits.Float64, Monotonic, State[float64, traits.Float64, Monotonic]]{}

	_ aggregation.Sum = &State[int64, traits.Int64, Monotonic]{}
	_ aggregation.Sum = &State[float64, traits.Float64, Monotonic]{}

	_ aggregation.Sum = &State[int64, traits.Int64, NonMonotonic]{}
	_ aggregation.Sum = &State[float64, traits.Float64, NonMonotonic]{}
)

func (s *State[N, Traits, M]) Sum() number.Number {
	var traits Traits
	return traits.ToNumber(s.value)
}

func (s *State[N, Traits, M]) Category() aggregation.Category {
	var m M
	return m.category()
}

func (Methods[N, Traits, M, Storage]) Kind() aggregation.Kind {
	return aggregation.SumKind
}

func (Methods[N, Traits, M, Storage]) Init(state *State[N, Traits, M], _ aggregator.Config) {
	// Note: storage is zero to start
}

func (Methods[N, Traits, M, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits, M]) {
	var traits Traits
	dest.value = traits.SwapAtomic(&resetSrc.value, 0)
}

func (Methods[N, Traits, M, Storage]) Reset(ptr *State[N, Traits, M]) {
	ptr.value = 0
}

func (Methods[N, Traits, M, Storage]) HasChange(ptr *State[N, Traits, M]) bool {
	return ptr.value == 0
}

func (Methods[N, Traits, M, Storage]) Update(state *State[N, Traits, M], value N) {
	var traits Traits
	traits.AddAtomic(&state.value, value)
}

func (Methods[N, Traits, M, Storage]) Merge(to, from *State[N, Traits, M]) {
	to.value += from.value
}

func (Methods[N, Traits, M, Storage]) Aggregation(state *State[N, Traits, M]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, M, Storage]) Storage(aggr aggregation.Aggregation) *State[N, Traits, M] {
	return aggr.(*State[N, Traits, M])
}

func (Methods[N, Traits, M, Storage]) SubtractSwap(newValue, oldValueModified *State[N, Traits, M]) {
	oldValueModified.value = newValue.value - oldValueModified.value
}
