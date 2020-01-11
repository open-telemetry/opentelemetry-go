// Copyright 2019, OpenTelemetry Authors
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

package histogram // import "go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Aggregator aggregates measure events and calculates
	// sum, count and buckets count.
	Aggregator struct {
		current    State
		checkpoint State
		bounds     []float64
		kind       core.NumberKind
	}

	State struct {
		Buckets []core.Number
		Count   core.Number
		Sum     core.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregator.Sum = &Aggregator{}
var _ aggregator.Count = &Aggregator{}
var _ aggregator.Histogram = &Aggregator{}

// New returns a new measure aggregator for computing count, sum and buckets count.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.
func New(desc *export.Descriptor, bounds []float64) *Aggregator {
	return &Aggregator{
		kind: desc.NumberKind(),
		current: State{
			Buckets: make([]core.Number, len(bounds)+1),
		},
		bounds: bounds,
	}
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Sum() (core.Number, error) {
	return c.checkpoint.Sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return int64(c.checkpoint.Count.AsUint64()), nil
}

func (c *Aggregator) Histogram() (State, error) {
	return c.checkpoint, nil
}

// Checkpoint saves the current bucket and resets the current bucket to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Min, Max, Sum, and Count are not consistent with each
// other.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *export.Descriptor) {
	// N.B. There is no atomic operation that can update all three
	// values at once without a memory allocation.
	//
	// This aggregator is intended to trade this correctness for
	// speed.
	//
	// Therefore, atomically swap fields independently, knowing
	// that individually the three parts of this aggregation could
	// be spread across multiple collections in rare cases.

	c.checkpoint.Count.SetUint64(c.current.Count.SwapUint64Atomic(0))
	c.checkpoint.Sum = c.current.Sum.SwapNumberAtomic(core.Number(0))
	c.checkpoint.Buckets = c.current.Buckets
	c.current.Buckets = make([]core.Number, len(c.bounds)+1)
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	kind := desc.NumberKind()

	c.current.Count.AddUint64Atomic(1)
	c.current.Sum.AddNumberAtomic(kind, number)

	for i, boundary := range c.bounds {
		if number.CoerceToFloat64(kind) <= boundary {
			c.current.Buckets[i].AddUint64Atomic(1)
			return nil
		}
	}

	c.current.Buckets[len(c.bounds)].AddUint64Atomic(1)
	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *export.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.checkpoint.Sum.AddNumber(desc.NumberKind(), o.checkpoint.Sum)
	c.checkpoint.Count.AddNumber(core.Uint64NumberKind, o.checkpoint.Count)

	for i := 0; i < len(c.current.Buckets); i++ {
		c.checkpoint.Buckets[i].AddNumber(desc.NumberKind(), o.checkpoint.Buckets[i])
	}
	return nil
}
