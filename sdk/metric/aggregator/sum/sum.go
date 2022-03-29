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
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

type (
	Config struct{}

	Methods[N number.Any, Traits traits.Any[N], Storage State[N, Traits]] struct{}

	State[N number.Any, Traits traits.Any[N]] struct {
		value N
	}
)

var (
	_ aggregator.Methods[int64, State[int64, traits.Int64], Config]       = Methods[int64, traits.Int64, State[int64, traits.Int64]]{}
	_ aggregator.Methods[float64, State[float64, traits.Float64], Config] = Methods[float64, traits.Float64, State[float64, traits.Float64]]{}

	_ aggregation.Sum = &State[int64, traits.Int64]{}
	_ aggregation.Sum = &State[float64, traits.Float64]{}
)

func (s *State[N, Traits]) Sum() (number.Number, error) {
	var traits Traits
	return traits.ToNumber(s.value), nil
}

func (s *State[N, Traits]) Kind() aggregation.Kind {
	return aggregation.SumKind
}

func (Methods[N, Traits, Storage]) Init(state *State[N, Traits], _ Config) {
	// Note: storage is zero to start
}

func (Methods[N, Traits, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits]) {
	var traits Traits
	if dest == nil {
		traits.SetAtomic(&resetSrc.value, 0)
		return
	}
	dest.value = traits.SwapAtomic(&resetSrc.value, 0)
}

func (Methods[N, Traits, Storage]) Update(state *State[N, Traits], value N) {
	var traits Traits
	traits.AddAtomic(&state.value, value)
}

func (Methods[N, Traits, Storage]) Merge(to, from *State[N, Traits]) {
	to.value += from.value
}

func (Methods[N, Traits, Storage]) Aggregation(state *State[N, Traits]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, Storage]) Storage(aggr aggregation.Aggregation) *State[N, Traits] {
	return aggr.(*State[N, Traits])
}

func (Methods[N, Traits, Storage]) Subtract(valueToModify, operand *State[N, Traits]) error {
	valueToModify.value -= operand.value
	return nil
}
