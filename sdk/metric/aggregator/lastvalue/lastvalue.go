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

package lastvalue // import "go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"

import (
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

var ErrNoSubtract = fmt.Errorf("lastvalue subtract not implemented")

type (
	Config struct{}

	Methods[N number.Any, Traits traits.Any[N], Storage State[N, Traits]] struct{}

	State[N number.Any, Traits traits.Any[N]] struct {
		lock      sync.Mutex
		value     N
		timestamp time.Time
	}
)

var (
	_ aggregator.Methods[int64, State[int64, traits.Int64], Config]       = Methods[int64, traits.Int64, State[int64, traits.Int64]]{}
	_ aggregator.Methods[float64, State[float64, traits.Float64], Config] = Methods[float64, traits.Float64, State[float64, traits.Float64]]{}

	_ aggregation.LastValue = &State[int64, traits.Int64]{}
	_ aggregation.LastValue = &State[float64, traits.Float64]{}
)

// LastValue returns the last-recorded lastValue value and the
// corresponding timestamp.  The error value aggregation.ErrNoData
// will be returned if (due to a race condition) the checkpoint was
// computed before the first value was set.
func (lv *State[N, Traits]) LastValue() (number.Number, time.Time, error) {
	var traits Traits
	lv.lock.Lock()
	defer lv.lock.Unlock()
	if lv.timestamp.IsZero() {
		return 0, time.Time{}, aggregation.ErrNoData
	}
	return traits.ToNumber(lv.value), lv.timestamp, nil
}

func (lv *State[N, Traits]) Kind() aggregation.Kind {
	return aggregation.LastValueKind
}

func (Methods[N, Traits, Storage]) Init(state *State[N, Traits], _ Config) {
	// Note: storage is zero to start
}

func (Methods[N, Traits, Storage]) Reset(ptr *State[N, Traits]) {
	ptr.value = 0
	ptr.timestamp = time.Time{}
}

func (Methods[N, Traits, Storage]) HasData(ptr *State[N, Traits]) bool {
	return ptr.timestamp.IsZero()
}

func (Methods[N, Traits, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits]) {
	resetSrc.lock.Lock()
	defer resetSrc.lock.Unlock()

	dest.value = resetSrc.value
	dest.timestamp = resetSrc.timestamp

	resetSrc.value = 0
	resetSrc.timestamp = time.Time{}
}

func (Methods[N, Traits, Storage]) Update(state *State[N, Traits], number N) {
	now := time.Now()

	state.lock.Lock()
	defer state.lock.Unlock()

	state.value = number
	state.timestamp = now
}

func (Methods[N, Traits, Storage]) Merge(to, from *State[N, Traits]) {
	if to.timestamp.After(from.timestamp) {
		return
	}

	to.value = from.value
	to.timestamp = from.timestamp
}

func (Methods[N, Traits, Storage]) Aggregation(state *State[N, Traits]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, Storage]) Storage(aggr aggregation.Aggregation) *State[N, Traits] {
	return aggr.(*State[N, Traits])
}

func (Methods[N, Traits, Storage]) Subtract(valueToModify, operand *State[N, Traits]) error {
	return ErrNoSubtract
}
