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

package array // import "go.opentelemetry.io/otel/sdk/metric/aggregator/array"

import (
	"context"
	"math"
	"sort"
	"sync"
	"unsafe"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

type (
	Aggregator struct {
		lock       sync.Mutex
		current    Points
		checkpoint Points
		ckptSum    core.Number
	}

	Points []core.Number
)

var _ export.Aggregator = &Aggregator{}

func New() *Aggregator {
	return &Aggregator{}
}

// Sum returns the sum of the checkpoint.
func (c *Aggregator) Sum() core.Number {
	return c.ckptSum
}

// Count returns the count of the checkpoint.
func (c *Aggregator) Count() int64 {
	return int64(len(c.checkpoint))
}

// Max returns the max of the checkpoint.
func (c *Aggregator) Max() (core.Number, error) {
	return c.checkpoint.Quantile(1)
}

// Min returns the min of the checkpoint.
func (c *Aggregator) Min() (core.Number, error) {
	return c.checkpoint.Quantile(0)
}

// Quantile returns the estimated quantile of the checkpoint.
func (c *Aggregator) Quantile(q float64) (core.Number, error) {
	return c.checkpoint.Quantile(q)
}

func (c *Aggregator) Collect(ctx context.Context, rec export.Record, exp export.Batcher) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, nil
	c.lock.Unlock()

	desc := rec.Descriptor()
	kind := desc.NumberKind()

	c.sort(kind)

	c.ckptSum = core.Number(0)

	for _, v := range c.checkpoint {
		c.ckptSum.AddNumber(kind, v)
	}

	exp.Export(ctx, rec, c)
}

func (c *Aggregator) Update(_ context.Context, number core.Number, rec export.Record) {
	desc := rec.Descriptor()
	kind := desc.NumberKind()

	if kind == core.Float64NumberKind && math.IsNaN(number.AsFloat64()) {
		// TODO warn
		// NOTE: add this to the specification.
		return
	}

	if !desc.Alternate() && number.IsNegative(kind) {
		// TODO warn
		return
	}

	c.lock.Lock()
	c.current = append(c.current, number)
	c.lock.Unlock()
}

func (c *Aggregator) Merge(oa export.Aggregator, desc *export.Descriptor) {
	o, _ := oa.(*Aggregator)
	if o == nil {
		// TODO warn
		return
	}

	c.ckptSum.AddNumber(desc.NumberKind(), o.ckptSum)
	c.checkpoint = combine(c.checkpoint, o.checkpoint, desc.NumberKind())
}

func (c *Aggregator) sort(kind core.NumberKind) {
	switch kind {
	case core.Float64NumberKind:
		sort.Float64s(*(*[]float64)(unsafe.Pointer(&c.checkpoint)))

	case core.Int64NumberKind:
		sort.Sort(&c.checkpoint)

	default:
		// NOTE: This can't happen because the SDK doesn't
		// support uint64-kind metric instruments.
		panic("Impossible case")
	}
}

func combine(a, b Points, kind core.NumberKind) Points {
	result := make(Points, 0, len(a)+len(b))

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

func (p *Points) Len() int {
	return len(*p)
}

func (p *Points) Less(i, j int) bool {
	// Note this is specialized for int64, because float64 is
	// handled by `sort.Float64s` and uint64 numbers never appear
	// in this data.
	return int64((*p)[i]) < int64((*p)[j])
}

func (p *Points) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// Quantile returns the least X such that Pr(x<X)>=q, where X is an
// element of the data set.
func (p *Points) Quantile(q float64) (core.Number, error) {
	if len(*p) == 0 {
		return core.Number(0), aggregator.ErrEmptyDataSet
	}

	if q < 0 || q > 1 {
		return core.Number(0), aggregator.ErrInvalidQuantile
	}

	if q == 0 || len(*p) == 1 {
		return (*p)[0], nil
	} else if q == 1 {
		return (*p)[len(*p)-1], nil
	}

	// Note: There's no interpolation being done here.  There are
	// many definitions for "quantile", some interpolate, some do
	// not.  What is expected?
	position := float64(len(*p)-1) * q
	ceil := int(math.Ceil(position))
	return (*p)[ceil], nil
}
