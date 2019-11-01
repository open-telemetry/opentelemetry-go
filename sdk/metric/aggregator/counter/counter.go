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

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/export"
)

// Aggregator aggregates counter events.
type Aggregator struct {
	// current holds current increments to this counter record
	current core.Number

	// checkpoint is a temporary used during Collect()
	checkpoint core.Number
}

var _ export.MetricAggregator = &Aggregator{}

// New returns a new counter aggregator.  This aggregator computes an
// atomic sum.
func New() *Aggregator {
	return &Aggregator{}
}

// AsNumber returns the accumulated count as an int64.
func (c *Aggregator) AsNumber() core.Number {
	return c.checkpoint.AsNumber()
}

// Collect checkpoints the current value (atomically) and exports it.
func (c *Aggregator) Collect(ctx context.Context, rec export.MetricRecord, exp export.MetricBatcher) {
	c.checkpoint = c.current.SwapNumberAtomic(core.Number(0))

	exp.Export(ctx, rec, c)
}

// Update modifies the current value (atomically) for later export.
func (c *Aggregator) Update(_ context.Context, number core.Number, rec export.MetricRecord) {
	desc := rec.Descriptor()
	kind := desc.NumberKind()
	if !desc.Alternate() && number.IsNegative(kind) {
		// TODO warn
		return
	}

	c.current.AddNumberAtomic(kind, number)
}

func (c *Aggregator) Merge(oa export.MetricAggregator, desc *export.Descriptor) {
	o, _ := oa.(*Aggregator)
	if o == nil {
		// TODO warn
		return
	}
	c.checkpoint.AddNumber(desc.NumberKind(), o.checkpoint)
}
