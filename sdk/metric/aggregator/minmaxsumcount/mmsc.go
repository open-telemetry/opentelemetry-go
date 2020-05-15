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

package minmaxsumcount // import "go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Aggregator aggregates events that form a distribution,
	// keeping only the min, max, sum, and count.
	Aggregator struct {
		lock       sync.Mutex
		current    state
		checkpoint state
		kind       metric.NumberKind
	}

	state struct {
		count metric.Number
		sum   metric.Number
		min   metric.Number
		max   metric.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregator.MinMaxSumCount = &Aggregator{}

// New returns a new aggregator for computing the min, max, sum, and
// count.  It does not compute quantile information other than Min and
// Max.
//
// This type uses a mutex for Update() and Checkpoint() concurrency.
func New(desc *metric.Descriptor) *Aggregator {
	kind := desc.NumberKind()
	return &Aggregator{
		kind: kind,
		current: state{
			count: metric.NewUint64Number(0),
			sum:   kind.Zero(),
			min:   kind.Maximum(),
			max:   kind.Minimum(),
		},
	}
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint.count.CoerceToInt64(metric.Uint64NumberKind), nil
}

// Min returns the minimum value in the checkpoint.
// The error value aggregator.ErrNoData will be returned
// if there were no measurements recorded during the checkpoint.
func (c *Aggregator) Min() (metric.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.checkpoint.count.IsZero(metric.Uint64NumberKind) {
		return c.kind.Zero(), aggregator.ErrNoData
	}
	return c.checkpoint.min, nil
}

// Max returns the maximum value in the checkpoint.
// The error value aggregator.ErrNoData will be returned
// if there were no measurements recorded during the checkpoint.
func (c *Aggregator) Max() (metric.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.checkpoint.count.IsZero(metric.Uint64NumberKind) {
		return c.kind.Zero(), aggregator.ErrNoData
	}
	return c.checkpoint.max, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *metric.Descriptor) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, c.emptyState()
	c.lock.Unlock()
}

func (c *Aggregator) emptyState() state {
	kind := c.kind
	return state{
		count: metric.NewUint64Number(0),
		sum:   kind.Zero(),
		min:   kind.Maximum(),
		max:   kind.Minimum(),
	}
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	kind := desc.NumberKind()

	c.lock.Lock()
	defer c.lock.Unlock()
	c.current.count.AddInt64(1)
	c.current.sum.AddNumber(kind, number)
	if number.CompareNumber(kind, c.current.min) < 0 {
		c.current.min = number
	}
	if number.CompareNumber(kind, c.current.max) > 0 {
		c.current.max = number
	}
	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.checkpoint.count.AddNumber(metric.Uint64NumberKind, o.checkpoint.count)
	c.checkpoint.sum.AddNumber(desc.NumberKind(), o.checkpoint.sum)

	if c.checkpoint.min.CompareNumber(desc.NumberKind(), o.checkpoint.min) > 0 {
		c.checkpoint.min.SetNumber(o.checkpoint.min)
	}
	if c.checkpoint.max.CompareNumber(desc.NumberKind(), o.checkpoint.max) < 0 {
		c.checkpoint.max.SetNumber(o.checkpoint.max)
	}
	return nil
}
