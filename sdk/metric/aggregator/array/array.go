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

package array // import "go.opentelemetry.io/otel/sdk/metric/aggregator/array"

import (
	"context"
	"math"
	"sort"
	"sync"
	"unsafe"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Aggregator aggregates events that form a distribution, keeping
	// an array with the exact set of values.
	Aggregator struct {
		// ckptSum needs to be aligned for 64-bit atomic operations.
		ckptSum    metric.Number
		lock       sync.Mutex
		current    points
		checkpoint points
	}

	points []metric.Number
)

var _ export.Aggregator = &Aggregator{}
var _ aggregator.MinMaxSumCount = &Aggregator{}
var _ aggregator.Distribution = &Aggregator{}
var _ aggregator.Points = &Aggregator{}

// New returns a new array aggregator, which aggregates recorded
// measurements by storing them in an array.  This type uses a mutex
// for Update() and Checkpoint() concurrency.
func New() *Aggregator {
	return &Aggregator{}
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	return c.ckptSum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return int64(len(c.checkpoint)), nil
}

// Max returns the maximum value in the checkpoint.
func (c *Aggregator) Max() (metric.Number, error) {
	return c.checkpoint.Quantile(1)
}

// Min returns the mininum value in the checkpoint.
func (c *Aggregator) Min() (metric.Number, error) {
	return c.checkpoint.Quantile(0)
}

// Quantile returns the estimated quantile of data in the checkpoint.
// It is an error if `q` is less than 0 or greated than 1.
func (c *Aggregator) Quantile(q float64) (metric.Number, error) {
	return c.checkpoint.Quantile(q)
}

// Points returns access to the raw data set.
func (c *Aggregator) Points() ([]metric.Number, error) {
	return c.checkpoint, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set, taking a lock to prevent concurrent Update() calls.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *metric.Descriptor) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, nil
	c.lock.Unlock()

	kind := desc.NumberKind()

	// TODO: This sort should be done lazily, only when quantiles
	// are requested.  The SDK specification says you can use this
	// aggregator to simply list values in the order they were
	// received as an alternative to requesting quantile information.
	c.sort(kind)

	c.ckptSum = metric.Number(0)

	for _, v := range c.checkpoint {
		c.ckptSum.AddNumber(kind, v)
	}
}

// Update adds the recorded measurement to the current data set.
// Update takes a lock to prevent concurrent Update() and Checkpoint()
// calls.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	c.lock.Lock()
	c.current = append(c.current, number)
	c.lock.Unlock()
	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.ckptSum.AddNumber(desc.NumberKind(), o.ckptSum)
	c.checkpoint = combine(c.checkpoint, o.checkpoint, desc.NumberKind())
	return nil
}

func (c *Aggregator) sort(kind metric.NumberKind) {
	switch kind {
	case metric.Float64NumberKind:
		sort.Float64s(*(*[]float64)(unsafe.Pointer(&c.checkpoint)))

	case metric.Int64NumberKind:
		sort.Sort(&c.checkpoint)

	default:
		// NOTE: This can't happen because the SDK doesn't
		// support uint64-kind metric instruments.
		panic("Impossible case")
	}
}

func combine(a, b points, kind metric.NumberKind) points {
	result := make(points, 0, len(a)+len(b))

	for len(a) != 0 && len(b) != 0 {
		if a[0].CompareNumber(kind, b[0]) < 0 {
			result = append(result, a[0])
			a = a[1:]
		} else {
			result = append(result, b[0])
			b = b[1:]
		}
	}
	result = append(result, a...)
	result = append(result, b...)
	return result
}

func (p *points) Len() int {
	return len(*p)
}

func (p *points) Less(i, j int) bool {
	// Note this is specialized for int64, because float64 is
	// handled by `sort.Float64s` and uint64 numbers never appear
	// in this data.
	return int64((*p)[i]) < int64((*p)[j])
}

func (p *points) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// Quantile returns the least X such that Pr(x<X)>=q, where X is an
// element of the data set.  This uses the "Nearest-Rank" definition
// of a quantile.
func (p *points) Quantile(q float64) (metric.Number, error) {
	if len(*p) == 0 {
		return metric.Number(0), aggregator.ErrNoData
	}

	if q < 0 || q > 1 {
		return metric.Number(0), aggregator.ErrInvalidQuantile
	}

	if q == 0 || len(*p) == 1 {
		return (*p)[0], nil
	} else if q == 1 {
		return (*p)[len(*p)-1], nil
	}

	position := float64(len(*p)-1) * q
	ceil := int(math.Ceil(position))
	return (*p)[ceil], nil
}
