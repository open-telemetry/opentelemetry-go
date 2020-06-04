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
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
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
		self  *Aggregator
		count metric.Number
		sum   metric.Number
		min   metric.Number
		max   metric.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregation.MinMaxSumCount = &state{}

// New returns a new aggregator for computing the min, max, sum, and
// count.  It does not compute quantile information other than Min and
// Max.
//
// This type uses a mutex for Update() and Checkpoint() concurrency.
func New(desc *metric.Descriptor) *Aggregator {
	agg := &Aggregator{
		kind: desc.NumberKind(),
	}
	agg.current = agg.emptyState()
	agg.checkpoint = agg.emptyState()
	return agg
}

// Kind returns aggregation.MinMaxSumCountKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.MinMaxSumCountKind
}

// Checkpoint saves the current state and resets the current state to
// the empty set.
func (c *Aggregator) Checkpoint(desc *metric.Descriptor) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, c.emptyState()
	c.lock.Unlock()
}

func (c *Aggregator) Swap() {
	c.checkpoint, c.current = c.current, c.checkpoint
}

func (c *Aggregator) emptyState() state {
	return state{
		self:  c,
		count: metric.NewUint64Number(0),
		sum:   c.kind.Zero(),
		min:   c.kind.Maximum(),
		max:   c.kind.Minimum(),
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

	c.current.count.AddNumber(metric.Uint64NumberKind, o.checkpoint.count)
	c.current.sum.AddNumber(desc.NumberKind(), o.checkpoint.sum)

	if c.current.min.CompareNumber(desc.NumberKind(), o.checkpoint.min) > 0 {
		c.current.min.SetNumber(o.checkpoint.min)
	}
	if c.current.max.CompareNumber(desc.NumberKind(), o.checkpoint.max) < 0 {
		c.current.max.SetNumber(o.checkpoint.max)
	}
	return nil
}

func (c *Aggregator) CheckpointedValue() aggregation.Aggregation {
	return &c.checkpoint
}

func (c *Aggregator) AccumulatedValue() aggregation.Aggregation {
	return &c.current
}

// Kind returns aggregation.MinMaxSumCountKind.
func (s *state) Kind() aggregation.Kind {
	return aggregation.MinMaxSumCountKind
}

// Sum returns the sum of values in the checkpoint.
func (s *state) Sum() (metric.Number, error) {
	return s.sum, nil
}

// Count returns the number of values in the checkpoint.
func (s *state) Count() (int64, error) {
	return s.count.CoerceToInt64(metric.Uint64NumberKind), nil
}

// Min returns the minimum value in the checkpoint.
// The error value aggregation.ErrNoData will be returned
// if there were no measurements recorded during the checkpoint.
func (s *state) Min() (metric.Number, error) {
	if s.count.IsZero(metric.Uint64NumberKind) {
		return s.self.kind.Zero(), aggregation.ErrNoData
	}
	return s.min, nil
}

// Max returns the maximum value in the checkpoint.
// The error value aggregation.ErrNoData will be returned
// if there were no measurements recorded during the checkpoint.
func (s *state) Max() (metric.Number, error) {
	if s.count.IsZero(metric.Uint64NumberKind) {
		return s.self.kind.Zero(), aggregation.ErrNoData
	}
	return s.max, nil
}
