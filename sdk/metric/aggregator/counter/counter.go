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

package counter

import (
	"context"
	"math"
	"sync/atomic"

	api "go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/sdk/export"
	"go.opentelemetry.io/sdk/metric/internal"
)

type (
	// Aggregator aggregates counter events.
	Aggregator struct {
		// live holds current increments to this counter record
		live uint64

		// save is a temporary used during Collect()
		save uint64
	}
)

var _ export.MetricAggregator = &Aggregator{}

// New returns a new counter aggregator.  This aggregator computes an
// atomic sum.
func New() *Aggregator {
	return &Aggregator{}
}

// AsInt64 returns the accumulated count as an int64.
func (c *Aggregator) AsInt64() int64 {
	return int64(c.save)
}

// AsFloat64 returns the accumulated count as an float64.
func (c *Aggregator) AsFloat64() float64 {
	return math.Float64frombits(c.save)
}

// Collect saves the current value (atomically) and exports it.
func (c *Aggregator) Collect(ctx context.Context, rec export.MetricRecord, exp export.MetricBatcher) {
	c.save = atomic.SwapUint64(&c.live, 0)

	if c.save == 0 {
		return
	}

	exp.Export(ctx, rec, c)
}

// Collect updates the current value (atomically) for later export.
func (c *Aggregator) Update(_ context.Context, value api.MeasurementValue, rec export.MetricRecord) {
	// TODO: See https://github.com/open-telemetry/opentelemetry-go/issues/196
	// Assumption is that `value` has a type corresponding to
	// `descriptor.ValueKind()`, which is not enforced for
	// RecordBatch measurements.  The AsInt64(), AsFloat64(), and
	// IsNegative() tests here are incorrect if the kind is a
	// variable.
	descriptor := rec.Descriptor()

	if !descriptor.Alternate() && value.IsNegative(descriptor.ValueKind()) {
		// TODO warn
		return
	}

	if descriptor.ValueKind() == api.Int64ValueKind {
		internal.NewAtomicInt64(&c.live).Add(value.AsInt64())
	} else {
		internal.NewAtomicFloat64(&c.live).Add(value.AsFloat64())
	}
}
