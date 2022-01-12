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
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

type Config struct {}

var _ aggregator.Aggregator[int64, Aggregator[int64, traits.Int64], Config] = &Aggregator[int64, traits.Int64]{}
var _ aggregator.Aggregator[float64, Aggregator[float64, traits.Float64], Config] = &Aggregator[float64, traits.Float64]{}

// Aggregator aggregates counter events.
type Aggregator[N number.Any, Traits traits.Any[N]] struct {
	// current holds current increments to this counter record
	// current needs to be aligned for 64-bit atomic operations.
	value  N
	traits Traits
}

func (c *Aggregator[N, Traits]) Init(_ Config) {
	c.value = 0
}

// Sum returns the last-checkpointed sum.  This will never return an
// error.
func (c *Aggregator[N, Traits]) Sum() (number.Number, error) {
	return c.traits.ToNumber(c.value), nil
}

// SynchronizedMove atomically saves the current value into oa and resets the
// current sum to zero.
func (c *Aggregator[N, Traits]) SynchronizedMove(o *Aggregator[N, Traits]) {
	if o == nil {
		c.traits.SetAtomic(&c.value, 0)
		return
	} 
	o.value = c.traits.SwapAtomic(&c.value, 0)
}

// Update atomically adds to the current value.
func (c *Aggregator[N, Traits]) Update(num N) {
	c.traits.AddAtomic(&c.value, num)
}

// Merge combines two counters by adding their sums.
func (c *Aggregator[N, Traits]) Merge(o *Aggregator[N, Traits]) {
	c.value += o.value
}
