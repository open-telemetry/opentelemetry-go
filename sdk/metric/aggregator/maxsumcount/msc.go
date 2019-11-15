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

package maxsumcount // import "go.opentelemetry.io/otel/sdk/metric/aggregator/maxsumcount"

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Aggregator aggregates measure events, keeping only the max,
	// sum, and count.
	Aggregator struct {
		current    state
		checkpoint state
	}

	state struct {
		count core.Number
		sum   core.Number
		max   core.Number
	}
)

// TODO: The SDK specification says this type should support Min
// values, see #319.

var _ export.Aggregator = &Aggregator{}
var _ aggregator.MaxSumCount = &Aggregator{}

// New returns a new measure aggregator for computing max, sum, and
// count.  It does not compute quantile information other than Max.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.  For greater consistency and lower
// performance, consider using Array or DDSketch aggregators.
func New() *Aggregator {
	return &Aggregator{}
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (core.Number, error) {
	return c.checkpoint.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return int64(c.checkpoint.count.AsUint64()), nil
}

// Max returns the maximum value in the checkpoint.
func (c *Aggregator) Max() (core.Number, error) {
	return c.checkpoint.max, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Max, Sum, and Count are not consistent with each
// other.
func (c *Aggregator) Checkpoint(ctx context.Context, _ *export.Descriptor) {
	// N.B. There is no atomic operation that can update all three
	// values at once without a memory allocation.
	//
	// This aggregator is intended to trade this correctness for
	// speed.
	//
	// Therefore, atomically swap fields independently, knowing
	// that individually the three parts of this aggregation could
	// be spread across multiple collections in rare cases.

	c.checkpoint.count.SetUint64(c.current.count.SwapUint64Atomic(0))
	c.checkpoint.sum = c.current.sum.SwapNumberAtomic(core.Number(0))
	c.checkpoint.max = c.current.max.SwapNumberAtomic(core.Number(0))
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	kind := desc.NumberKind()

	c.current.count.AddUint64Atomic(1)
	c.current.sum.AddNumberAtomic(kind, number)

	for {
		current := c.current.max.AsNumberAtomic()

		if number.CompareNumber(kind, current) <= 0 {
			break
		}
		if c.current.max.CompareAndSwapNumber(current, number) {
			break
		}
	}
	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *export.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.checkpoint.sum.AddNumber(desc.NumberKind(), o.checkpoint.sum)
	c.checkpoint.count.AddNumber(core.Uint64NumberKind, o.checkpoint.count)

	if c.checkpoint.max.CompareNumber(desc.NumberKind(), o.checkpoint.max) < 0 {
		c.checkpoint.max.SetNumber(o.checkpoint.max)
	}
	return nil
}
