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

package counter // import "go.opentelemetry.io/otel/sdk/metric/aggregator/counter"

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

// Aggregator aggregates counter events.
type Aggregator struct {
	// current holds current increments to this counter record
	// current needs to be aligned for 64-bit atomic operations.
	current core.Number

	// checkpoint is a temporary used during Checkpoint()
	// checkpoint needs to be aligned for 64-bit atomic operations.
	checkpoint core.Number
}

var _ export.Aggregator = &Aggregator{}
var _ aggregator.Sum = &Aggregator{}

// New returns a new counter aggregator implemented by atomic
// operations.  This aggregator implements the aggregator.Sum
// export interface.
func New() *Aggregator {
	return &Aggregator{}
}

// Sum returns the last-checkpointed sum.  This will never return an
// error.
func (c *Aggregator) Sum() (core.Number, error) {
	return c.checkpoint, nil
}

// Checkpoint atomically saves the current value and resets the
// current sum to zero.
func (c *Aggregator) Checkpoint(ctx context.Context, _ *export.Descriptor) {
	c.checkpoint = c.current.SwapNumberAtomic(core.Number(0))
}

// Update atomically adds to the current value.
func (c *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	c.current.AddNumberAtomic(desc.NumberKind(), number)
	return nil
}

// Merge combines two counters by adding their sums.
func (c *Aggregator) Merge(oa export.Aggregator, desc *export.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}
	c.checkpoint.AddNumber(desc.NumberKind(), o.checkpoint)
	return nil
}
