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
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

type (
	Config struct {}

	// Aggregator aggregates lastValue events.
	Aggregator[N number.Any, Traits traits.Any[N]] struct {
		lock      sync.Mutex
		value     N
		timestamp time.Time
		traits    Traits
	}
)

var _ aggregator.Aggregator[int64, Aggregator[int64, traits.Int64], Config] = &Aggregator[int64, traits.Int64]{}
var _ aggregator.Aggregator[float64, Aggregator[float64, traits.Float64], Config] = &Aggregator[float64, traits.Float64]{}

// New returns a new lastValue aggregator.  This aggregator retains the
// last value and timestamp that were recorded.
func (a *Aggregator[N, Traits]) Init(_ Config) {
	a.value = 0
	a.timestamp = time.Time{}
}

// LastValue returns the last-recorded lastValue value and the
// corresponding timestamp.  The error value aggregation.ErrNoData
// will be returned if (due to a race condition) the checkpoint was
// computed before the first value was set.
func (g *Aggregator[N, Traits]) LastValue() (number.Number, time.Time, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.timestamp.IsZero() {
		return 0, time.Time{}, aggregation.ErrNoData
	}
	return g.traits.ToNumber(g.value), g.timestamp, nil
}

// SynchronizedMove atomically saves the current value.
func (g *Aggregator[N, Traits]) SynchronizedMove(o *Aggregator[N, Traits]) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if o != nil {
		o.value = g.value
		o.timestamp = g.timestamp
	}
	g.value = 0
	g.timestamp = time.Time{}
}

// Update atomically sets the current "last" value.
func (g *Aggregator[N, Traits]) Update(number N) {
	now := time.Now()

	g.lock.Lock()
	defer g.lock.Unlock()

	g.value = number
	g.timestamp = now
}

// Merge combines state from two aggregators.  The most-recently set
// value is chosen.
func (g *Aggregator[N, Traits]) Merge(o *Aggregator[N, Traits]) {
	if g.timestamp.After(o.timestamp) {
		return
	}

	g.value = o.value
	g.timestamp = o.timestamp
}
