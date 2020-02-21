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

package minmaxsumcount // import "go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"

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
		// current has to be aligned for 64-bit atomic operations.
		current state
		// checkpoint has to be aligned for 64-bit atomic operations.
		checkpoint state
		kind       core.NumberKind
	}

	state struct {
		// all fields have to be aligned for 64-bit atomic operations.
		count core.Number
		sum   core.Number
		min   core.Number
		max   core.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregator.MinMaxSumCount = &Aggregator{}

// New returns a new measure aggregator for computing min, max, sum, and
// count.  It does not compute quantile information other than Max.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.  For greater consistency and lower
// performance, consider using Array or DDSketch aggregators.
func New(desc *export.Descriptor) *Aggregator {
	return &Aggregator{
		kind:    desc.NumberKind(),
		current: unsetMinMaxSumCount(desc.NumberKind()),
	}
}

func unsetMinMaxSumCount(kind core.NumberKind) state {
	return state{min: kind.Maximum(), max: kind.Minimum()}
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (core.Number, error) {
	return c.checkpoint.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return int64(c.checkpoint.count.AsUint64()), nil
}

// Min returns the minimum value in the checkpoint.
// The error value aggregator.ErrEmptyDataSet will be returned if
// (due to a race condition) the checkpoint was set prior to
// current.min being computed in Update().
//
// Note: If a measure's recorded values for a given checkpoint are
// all equal to NumberKind.Maximum(), Min() will return ErrEmptyDataSet
func (c *Aggregator) Min() (core.Number, error) {
	if c.checkpoint.min == c.kind.Maximum() {
		return core.Number(0), aggregator.ErrEmptyDataSet
	}
	return c.checkpoint.min, nil
}

// Max returns the maximum value in the checkpoint.
// The error value aggregator.ErrEmptyDataSet will be returned if
// (due to a race condition) the checkpoint was set prior to
// current.max being computed in Update().
//
// Note: If a measure's recorded values for a given checkpoint are
// all equal to NumberKind.Minimum(), Max() will return ErrEmptyDataSet
func (c *Aggregator) Max() (core.Number, error) {
	if c.checkpoint.max == c.kind.Minimum() {
		return core.Number(0), aggregator.ErrEmptyDataSet
	}
	return c.checkpoint.max, nil
}

// Checkpoint saves the current state and resets the current state to
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

	c.checkpoint.count.SetUint64(c.current.count.SwapUint64Atomic(0))
	c.checkpoint.sum = c.current.sum.SwapNumberAtomic(core.Number(0))
	c.checkpoint.max = c.current.max.SwapNumberAtomic(c.kind.Minimum())
	c.checkpoint.min = c.current.min.SwapNumberAtomic(c.kind.Maximum())
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	kind := desc.NumberKind()

	c.current.count.AddUint64Atomic(1)
	c.current.sum.AddNumberAtomic(kind, number)

	for {
		current := c.current.min.AsNumberAtomic()

		if number.CompareNumber(kind, current) >= 0 {
			break
		}
		if c.current.min.CompareAndSwapNumber(current, number) {
			break
		}
	}
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

	if c.checkpoint.min.CompareNumber(desc.NumberKind(), o.checkpoint.min) > 0 {
		c.checkpoint.min.SetNumber(o.checkpoint.min)
	}
	if c.checkpoint.max.CompareNumber(desc.NumberKind(), o.checkpoint.max) < 0 {
		c.checkpoint.max.SetNumber(o.checkpoint.max)
	}
	return nil
}
