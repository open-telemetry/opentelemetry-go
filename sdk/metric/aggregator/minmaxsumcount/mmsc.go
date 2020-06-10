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
		lock sync.Mutex
		kind metric.NumberKind
		state
	}

	state struct {
		count int64
		sum   metric.Number
		min   metric.Number
		max   metric.Number
	}
)

var w_ export.Aggregator = &Aggregator{}
var _ aggregator.MinMaxSumCount = &Aggregator{}

// New returns a new aggregator for computing the min, max, sum, and
// count.  It does not compute quantile information other than Min and
// Max.
//
// This type uses a mutex for Update() and Checkpoint() concurrency.
func New(desc *metric.Descriptor) *Aggregator {
	kind := desc.NumberKind()
	return &Aggregator{
		kind:  kind,
		state: emptyState(kind),
	}
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	return c.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return c.count, nil
}

// Min returns the minimum value in the checkpoint.
// The error value aggregator.ErrNoData will be returned
// if there were no measurements recorded during the checkpoint.
func (c *Aggregator) Min() (metric.Number, error) {
	if c.count == 0 {
		return 0, aggregator.ErrNoData
	}
	return c.min, nil
}

// Max returns the maximum value in the checkpoint.
// The error value aggregator.ErrNoData will be returned
// if there were no measurements recorded during the checkpoint.
func (c *Aggregator) Max() (metric.Number, error) {
	if c.count == 0 {
		return 0, aggregator.ErrNoData
	}
	return c.max, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.
func (c *Aggregator) Checkpoint(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.lock.Lock()
	o.state, c.state = c.state, emptyState(c.kind)
	c.lock.Unlock()

	return nil
}

func emptyState(kind metric.NumberKind) state {
	return state{
		count: 0,
		sum:   0,
		min:   kind.Maximum(),
		max:   kind.Minimum(),
	}
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	kind := desc.NumberKind()

	c.lock.Lock()
	defer c.lock.Unlock()
	c.count++
	c.sum.AddNumber(kind, number)
	if number.CompareNumber(kind, c.min) < 0 {
		c.min = number
	}
	if number.CompareNumber(kind, c.max) > 0 {
		c.max = number
	}
	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.count += o.count
	c.sum.AddNumber(desc.NumberKind(), o.sum)

	if c.min.CompareNumber(desc.NumberKind(), o.min) > 0 {
		c.min.SetNumber(o.min)
	}
	if c.max.CompareNumber(desc.NumberKind(), o.max) < 0 {
		c.max.SetNumber(o.max)
	}
	return nil
}
