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
	"go.opentelemetry.io/otel/sdk/internal"
)

type (
	// Aggregator aggregates measure events, keeping only the max,
	// sum, and count.
	Aggregator struct {
		// states has to be aligned for 64-bit atomic operations.
		states [2]state
		lock   internal.StateLocker
		kind   core.NumberKind
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
// This aggregator uses the StateLocker pattern to guarantee
// the count, sum, min and max are consistent within a checkpoint
func New(desc *export.Descriptor) *Aggregator {
	kind := desc.NumberKind()
	return &Aggregator{
		kind: kind,
		states: [2]state{
			{
				count: core.NewUint64Number(0),
				sum:   kind.Zero(),
				min:   kind.Maximum(),
				max:   kind.Minimum(),
			},
			{
				count: core.NewUint64Number(0),
				sum:   kind.Zero(),
				min:   kind.Maximum(),
				max:   kind.Minimum(),
			},
		},
	}
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (core.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint().sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint().count.CoerceToInt64(core.Uint64NumberKind), nil
}

// Min returns the minimum value in the checkpoint.
// The error value aggregator.ErrEmptyDataSet will be returned
// if there were no measurements recorded during the checkpoint.
func (c *Aggregator) Min() (core.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.checkpoint().count.IsZero(core.Uint64NumberKind) {
		return c.kind.Zero(), aggregator.ErrEmptyDataSet
	}
	return c.checkpoint().min, nil
}

// Max returns the maximum value in the checkpoint.
// The error value aggregator.ErrEmptyDataSet will be returned
// if there were no measurements recorded during the checkpoint.
func (c *Aggregator) Max() (core.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.checkpoint().count.IsZero(core.Uint64NumberKind) {
		return c.kind.Zero(), aggregator.ErrEmptyDataSet
	}
	return c.checkpoint().max, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *export.Descriptor) {
	c.lock.SwapActiveState(c.resetCheckpoint)
}

// checkpoint returns the "cold" state, i.e. state collected prior to the
// most recent Checkpoint() call
func (c *Aggregator) checkpoint() *state {
	return &c.states[c.lock.ColdIdx()]
}

func (c *Aggregator) resetCheckpoint() {
	checkpoint := c.checkpoint()

	checkpoint.count.SetUint64(0)
	checkpoint.sum.SetNumber(c.kind.Zero())
	checkpoint.min.SetNumber(c.kind.Maximum())
	checkpoint.max.SetNumber(c.kind.Minimum())
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	kind := desc.NumberKind()

	cIdx := c.lock.Start()
	defer c.lock.End(cIdx)

	current := &c.states[cIdx]
	current.count.AddUint64Atomic(1)
	current.sum.AddNumberAtomic(kind, number)

	for {
		cmin := current.min.AsNumberAtomic()

		if number.CompareNumber(kind, cmin) >= 0 {
			break
		}
		if current.min.CompareAndSwapNumber(cmin, number) {
			break
		}
	}
	for {
		cmax := current.max.AsNumberAtomic()

		if number.CompareNumber(kind, cmax) <= 0 {
			break
		}
		if current.max.CompareAndSwapNumber(cmax, number) {
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

	// Lock() synchronizes Merge() and Checkpoint() to ensure all operations of
	// Merge() are performed on the same state.
	c.lock.Lock()
	defer c.lock.Unlock()

	current := c.checkpoint()
	ocheckpoint := o.checkpoint()

	current.count.AddNumber(core.Uint64NumberKind, ocheckpoint.count)
	current.sum.AddNumber(desc.NumberKind(), ocheckpoint.sum)

	if current.min.CompareNumber(desc.NumberKind(), ocheckpoint.min) > 0 {
		current.min.SetNumber(ocheckpoint.min)
	}
	if current.max.CompareNumber(desc.NumberKind(), ocheckpoint.max) < 0 {
		current.max.SetNumber(ocheckpoint.max)
	}
	return nil
}
