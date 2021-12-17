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
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

// Aggregator aggregates counter events.
type Aggregator struct {
	// current holds current increments to this counter record
	// current needs to be aligned for 64-bit atomic operations.
	value number.Number
}

// New returns a new counter aggregator implemented by atomic
// operations.  This aggregator implements the aggregation.Sum
// export interface.
func New(cnt int) []Aggregator {
	return make([]Aggregator, cnt)
}

// Sum returns the last-checkpointed sum.  This will never return an
// error.
func (c *Aggregator) Sum() (number.Number, error) {
	return c.value, nil
}

// SynchronizedMove atomically saves the current value into oa and resets the
// current sum to zero.
func (c *Aggregator) SynchronizedMove(oa export.Aggregator, _ *sdkapi.Descriptor) error {
	if oa == nil {
		c.value.SetRawAtomic(0)
		return nil
	}
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}
	o.value = c.value.SwapAtomic(number.Number(0))
	return nil
}

// Update atomically adds to the current value.
func (c *Aggregator) Update(num number.Number, desc *sdkapi.Descriptor) {
	c.value.AddNumberAtomic(desc.NumberKind(), num)
}

// Merge combines two counters by adding their sums.
func (c *Aggregator) Merge(oa export.Aggregator, desc *sdkapi.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}
	c.value.Add(desc.NumberKind(), o.value)
	return nil
}
