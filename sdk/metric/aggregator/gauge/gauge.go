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

package gauge // import "go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"

import (
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

type (
	Methods[N number.Any, Traits traits.Any[N], Storage State[N, Traits]] struct{}

	State[N number.Any, Traits traits.Any[N]] struct {
		value N
	}
)

var (
	_ aggregator.Methods[int64, State[int64, traits.Int64]]       = Methods[int64, traits.Int64, State[int64, traits.Int64]]{}
	_ aggregator.Methods[float64, State[float64, traits.Float64]] = Methods[float64, traits.Float64, State[float64, traits.Float64]]{}

	_ aggregation.Gauge = &State[int64, traits.Int64]{}
	_ aggregation.Gauge = &State[float64, traits.Float64]{}
)

func NewInt64(value int64) aggregation.Gauge {
	return &State[int64, traits.Int64]{value: value}
}
func NewFloat64(value float64) aggregation.Gauge {
	return &State[float64, traits.Float64]{value: value}
}

func (lv *State[N, Traits]) Gauge() number.Number {
	var traits Traits
	return traits.ToNumber(lv.value)
}

func (lv *State[N, Traits]) Category() aggregation.Category {
	return aggregation.GaugeCategory
}

func (Methods[N, Traits, Storage]) Kind() aggregation.Kind {
	return aggregation.GaugeKind
}

func (Methods[N, Traits, Storage]) Init(state *State[N, Traits], _ aggregator.Config) {
	// Note: storage is zero to start
}

func (Methods[N, Traits, Storage]) Reset(ptr *State[N, Traits]) {
	var traits Traits
	traits.SetAtomic(&ptr.value, 0)
}

func (Methods[N, Traits, Storage]) HasChange(ptr *State[N, Traits]) bool {
	return ptr.value == 0
}

func (Methods[N, Traits, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits]) {
	var traits Traits
	dest.value = traits.SwapAtomic(&resetSrc.value, 0)
}

func (Methods[N, Traits, Storage]) Update(state *State[N, Traits], number N) {
	var traits Traits
	traits.SetAtomic(&state.value, number)
}

func (Methods[N, Traits, Storage]) Merge(to, from *State[N, Traits]) {
	to.value = from.value
}

func (Methods[N, Traits, Storage]) Aggregation(state *State[N, Traits]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, Storage]) Storage(aggr aggregation.Aggregation) *State[N, Traits] {
	return aggr.(*State[N, Traits])
}

func (Methods[N, Traits, Storage]) SubtractSwap(valueUnmodified, operandToModify *State[N, Traits]) {
	operandToModify.value = valueUnmodified.value - operandToModify.value
}
