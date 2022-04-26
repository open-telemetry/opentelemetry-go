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
)

type (
	Methods[N number.Any, Traits number.Traits[N], Storage State[N, Traits]] struct{}

	State[N number.Any, Traits number.Traits[N]] struct {
		value N
	}

	Int64   = State[int64, number.Int64Traits]
	Float64 = State[float64, number.Float64Traits]
)

var (
	_ aggregator.Methods[int64, Int64]     = Methods[int64, number.Int64Traits, Int64]{}
	_ aggregator.Methods[float64, Float64] = Methods[float64, number.Float64Traits, Float64]{}

	_ aggregation.Gauge = &Int64{}
	_ aggregation.Gauge = &Float64{}
)

func NewInt64(x int64) *Int64 {
	return &Int64{value: x}
}

func NewFloat64(x float64) *Float64 {
	return &Float64{value: x}
}

func (lv *State[N, Traits]) Gauge() number.Number {
	var t Traits
	return t.ToNumber(lv.value)
}

func (lv *State[N, Traits]) Kind() aggregation.Kind {
	return aggregation.GaugeKind
}

func (Methods[N, Traits, Storage]) Kind() aggregation.Kind {
	return aggregation.GaugeKind
}

func (Methods[N, Traits, Storage]) Init(state *State[N, Traits], _ aggregator.Config) {
	// Note: storage is zero to start
}

func (Methods[N, Traits, Storage]) Reset(ptr *State[N, Traits]) {
	var t Traits
	t.SetAtomic(&ptr.value, 0)
}

func (Methods[N, Traits, Storage]) HasChange(ptr *State[N, Traits]) bool {
	return ptr.value != 0
}

func (Methods[N, Traits, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits]) {
	var t Traits
	dest.value = t.SwapAtomic(&resetSrc.value, 0)
}

func (Methods[N, Traits, Storage]) Update(state *State[N, Traits], number N) {
	if !aggregator.RangeTest[N, Traits](number, aggregation.GaugeCategory) {
		return
	}
	var t Traits
	t.SetAtomic(&state.value, number)
}

func (Methods[N, Traits, Storage]) Merge(to, from *State[N, Traits]) {
	to.value = from.value
}

func (Methods[N, Traits, Storage]) ToAggregation(state *State[N, Traits]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, Storage]) ToStorage(aggr aggregation.Aggregation) (*State[N, Traits], bool) {
	r, ok := aggr.(*State[N, Traits])
	return r, ok
}

func (Methods[N, Traits, Storage]) SubtractSwap(valueUnmodified, operandToModify *State[N, Traits]) {
	operandToModify.value = valueUnmodified.value - operandToModify.value
}
