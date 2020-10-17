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

package pull // import "go.opentelemetry.io/otel/sdk/metric/controller/pull"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/registry"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
	"go.opentelemetry.io/otel/sdk/resource"
)

// DefaultCachePeriod determines how long a recently-computed result
// will be returned without gathering metric data again.
const DefaultCachePeriod time.Duration = 10 * time.Second

// Controller manages access to a *sdk.Accumulator and *basic.Processor. Use
// MeterProvider() for obtaining Meters. Use Foreach() for accessing current
// records.
type Controller struct {
	accumulator  *sdk.Accumulator
	checkpointer export.Checkpointer
	provider     *registry.MeterProvider
	period       time.Duration
	lastCollect  time.Time
	clock        controllerTime.Clock
	checkpoint   export.CheckpointSet
}

// New returns a *Controller configured with an export.Checkpointer.
//
// Pull controllers are typically used in an environment where there
// are multiple readers.  It is common, therefore, when configuring a
// basic Processor for use with this controller, to use a
// CumulativeExport strategy and the basic.WithMemory(true) option,
// which ensures that every CheckpointSet includes full state.
func New(checkpointer export.Checkpointer, options ...Option) *Controller {
	config := &Config{
		Resource:    resource.Empty(),
		CachePeriod: DefaultCachePeriod,
	}
	for _, opt := range options {
		opt.Apply(config)
	}
	accum := sdk.NewAccumulator(
		checkpointer,
		sdk.WithResource(config.Resource),
	)
	return &Controller{
		accumulator:  accum,
		checkpointer: checkpointer,
		provider:     registry.NewMeterProvider(accum),
		period:       config.CachePeriod,
		checkpoint:   checkpointer.CheckpointSet(),
		clock:        controllerTime.RealClock{},
	}
}

// SetClock sets the clock used for caching.  For testing purposes.
func (c *Controller) SetClock(clock controllerTime.Clock) {
	c.checkpointer.CheckpointSet().Lock()
	defer c.checkpointer.CheckpointSet().Unlock()
	c.clock = clock
}

// MeterProvider returns a MeterProvider for the implementation managed by
// this controller.
func (c *Controller) MeterProvider() otel.MeterProvider {
	return c.provider
}

// Foreach gives the caller read-locked access to the current
// export.CheckpointSet.
func (c *Controller) ForEach(ks export.ExportKindSelector, f func(export.Record) error) error {
	c.checkpointer.CheckpointSet().RLock()
	defer c.checkpointer.CheckpointSet().RUnlock()

	return c.checkpoint.ForEach(ks, f)
}

// Collect requests a collection.  The collection will be skipped if
// the last collection is aged less than the CachePeriod.
func (c *Controller) Collect(ctx context.Context) error {
	c.checkpointer.CheckpointSet().Lock()
	defer c.checkpointer.CheckpointSet().Unlock()

	if c.period > 0 {
		now := c.clock.Now()
		elapsed := now.Sub(c.lastCollect)

		if elapsed < c.period {
			return nil
		}
		c.lastCollect = now
	}

	c.checkpointer.StartCollection()
	c.accumulator.Collect(ctx)
	err := c.checkpointer.FinishCollection()
	c.checkpoint = c.checkpointer.CheckpointSet()
	return err
}
