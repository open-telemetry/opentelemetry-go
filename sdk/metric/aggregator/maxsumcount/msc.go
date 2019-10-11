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

package maxsumcount

import (
	"context"
	"math"
	"sync/atomic"

	api "go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/sdk/export"
	"go.opentelemetry.io/sdk/metric/internal"
)

type (
	// Aggregator aggregates measure events, keeping only the max,
	// sum, and count.
	Aggregator struct {
		live state
		save state
	}

	state struct {
		count uint64
		sum   uint64
		max   uint64
	}
)

var _ export.MetricAggregator = &Aggregator{}

// New returns a new measure aggregator for computing max, sum, and count.
func New() *Aggregator {
	return &Aggregator{}
}

// SumAsInt64 returns the accumulated sum as an int64.
func (c *Aggregator) SumAsInt64() int64 {
	return int64(c.save.sum)
}

// SumAsFloat64 returns the accumulated sum as an float64.
func (c *Aggregator) SumAsFloat64() float64 {
	return math.Float64frombits(c.save.sum)
}

// Count returns the accumulated count.
func (c *Aggregator) Count() uint64 {
	return c.save.count
}

// MaxAsInt64 returns the accumulated max as an int64.
func (c *Aggregator) MaxAsInt64() int64 {
	return int64(c.save.max)
}

// MaxAsFloat64 returns the accumulated max as an float64.
func (c *Aggregator) MaxAsFloat64() float64 {
	return math.Float64frombits(c.save.max)
}

// Collect saves the current value (atomically) and exports it.
func (c *Aggregator) Collect(ctx context.Context, rec export.MetricRecord, exp export.MetricBatcher) {

	// N.B. There is no atomic operation that can update all three
	// values at once, so there are races between Update() and
	// Collect().  Therefore, atomically swap fields independently,
	// knowing that individually the three parts of this aggregation
	// could be spread across multiple collections in rare cases.

	c.save.count = atomic.SwapUint64(&c.live.count, 0)
	c.save.sum = atomic.SwapUint64(&c.live.sum, 0)
	c.save.max = atomic.SwapUint64(&c.live.max, 0)

	exp.Export(ctx, rec, c)
}

// Collect updates the current value (atomically) for later export.
func (c *Aggregator) Update(_ context.Context, value api.MeasurementValue, rec export.MetricRecord) {
	descriptor := rec.Descriptor()

	if !descriptor.Alternate() && value.IsNegative(descriptor.ValueKind()) {
		// TODO warn
		return
	}

	atomic.AddUint64(&c.live.count, 1)

	if descriptor.ValueKind() == api.Int64ValueKind {
		internal.NewAtomicInt64(&c.live.sum).Add(value.AsInt64())
	} else {
		internal.NewAtomicFloat64(&c.live.sum).Add(value.AsFloat64())
	}

	if descriptor.ValueKind() == api.Int64ValueKind {
		update := value.AsInt64()
		for {
			current := internal.NewAtomicInt64(&c.live.max).Load()

			if update <= current {
				break
			}

			if atomic.CompareAndSwapUint64(&c.live.max, uint64(current), uint64(update)) {
				break
			}
		}
	} else {
		update := value.AsFloat64()
		for {
			current := internal.NewAtomicFloat64(&c.live.max).Load()

			if update <= current {
				break
			}

			if atomic.CompareAndSwapUint64(&c.live.max, math.Float64bits(current), math.Float64bits(update)) {
				break
			}
		}
	}
}
